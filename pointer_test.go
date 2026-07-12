package jsonpointer

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseStrict(t *testing.T) {
	tests := []struct {
		name    string
		pointer string
		want    []string
		wantErr error
	}{
		{name: "root", pointer: "", want: nil},
		{name: "empty token", pointer: "/", want: []string{""}},
		{name: "nested", pointer: "/users/0/name", want: []string{"users", "0", "name"}},
		{name: "escaped", pointer: "/foo~1bar/tilde~0key", want: []string{"foo/bar", "tilde~key"}},
		{name: "missing leading slash", pointer: "users/0/name", wantErr: ErrInvalidPointer},
		{name: "bad escape", pointer: "/~2", wantErr: ErrInvalidPointer},
		{name: "trailing tilde", pointer: "/foo~", wantErr: ErrInvalidPointer},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.pointer)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				assert.True(t, got.IsRoot())
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.pointer, got.String())
			assert.Equal(t, tt.want, got.Tokens())
		})
	}
}

func TestFromTokensBuildsImmutablePointer(t *testing.T) {
	tokens := []string{"users", "~2", "name/first"}
	p := FromTokens(tokens...)

	tokens[0] = "mutated"
	got := p.Tokens()
	got[1] = "mutated"

	assert.Equal(t, []string{"users", "~2", "name/first"}, p.Tokens())
	assert.Equal(t, "/users/~02/name~1first", p.String())

	parsed, err := Parse(p.String())
	require.NoError(t, err)
	assert.Equal(t, p.Tokens(), parsed.Tokens())
}

func TestPointerConstructionHasNoLengthPolicy(t *testing.T) {
	token := strings.Repeat("a", 2048)

	parsed, err := Parse("/" + token)
	require.NoError(t, err)
	assert.Equal(t, []string{token}, parsed.Tokens())

	fromTokens := FromTokens(token)
	assert.Equal(t, []string{token}, fromTokens.Tokens())
	assert.Equal(t, "/"+token, fromTokens.String())
}

func TestPointerParentChildAndRoot(t *testing.T) {
	root := Root()
	assert.True(t, root.IsRoot())
	assert.Equal(t, "", root.String())

	_, err := root.Parent()
	require.ErrorIs(t, err, ErrNoParent)
	var pointerErr *Error
	assert.False(t, errors.As(err, &pointerErr))

	child := root.Child("users", "0")
	assert.Equal(t, "/users/0", child.String())

	parent, err := child.Parent()
	require.NoError(t, err)
	assert.Equal(t, "/users", parent.String())
	assert.Equal(t, "", root.String())
}

func TestTokenEscaping(t *testing.T) {
	tests := []struct {
		token   string
		encoded string
	}{
		{token: "", encoded: ""},
		{token: "name", encoded: "name"},
		{token: "foo/bar", encoded: "foo~1bar"},
		{token: "tilde~key", encoded: "tilde~0key"},
		{token: "~/", encoded: "~0~1"},
	}

	for _, tt := range tests {
		t.Run(tt.token, func(t *testing.T) {
			assert.Equal(t, tt.encoded, EscapeToken(tt.token))
			got, err := UnescapeToken(tt.encoded)
			require.NoError(t, err)
			assert.Equal(t, tt.token, got)
		})
	}

	_, err := UnescapeToken("~2")
	require.ErrorIs(t, err, ErrInvalidPointer)
}

func TestIndexHelpers(t *testing.T) {
	hugeIndex := "9999999999999999999999999999999999999999"
	valid := []string{"0", "1", "10", hugeIndex}
	for _, token := range valid {
		assert.True(t, IsArrayIndex(token), token)
	}

	invalid := []string{"", "01", "-1", "+1", "1.2", "abc", "-", "\uFF11"}
	for _, token := range invalid {
		assert.False(t, IsArrayIndex(token), token)
	}
}

func TestTraversalErrorContext(t *testing.T) {
	p := mustPointer(t, "/users/1/name")
	_, err := p.Value(map[string]any{"users": []any{map[string]any{"name": "Alice"}}})
	require.ErrorIs(t, err, ErrIndexOutOfBounds)

	var pointerErr *Error
	require.True(t, errors.As(err, &pointerErr))
	assert.Equal(t, p.String(), pointerErr.Pointer().String())
	assert.Equal(t, "1", pointerErr.Token())
	assert.Equal(t, 1, pointerErr.Depth())
}

func TestTraversalErrorIncludesEmptyToken(t *testing.T) {
	p := mustPointer(t, "/empty/")
	_, err := p.Value(map[string]any{"empty": map[string]any{}})
	require.ErrorIs(t, err, ErrKeyNotFound)

	var pointerErr *Error
	require.True(t, errors.As(err, &pointerErr))
	assert.Equal(t, "", pointerErr.Token())
	assert.Equal(t, 1, pointerErr.Depth())
}

func TestParseErrorHasNoTraversalContext(t *testing.T) {
	_, err := Parse("/~2")
	require.ErrorIs(t, err, ErrInvalidPointer)

	var pointerErr *Error
	assert.False(t, errors.As(err, &pointerErr))
}

func mustPointer(t *testing.T, pointer string) Pointer {
	t.Helper()

	p, err := Parse(pointer)
	require.NoError(t, err)
	return p
}
