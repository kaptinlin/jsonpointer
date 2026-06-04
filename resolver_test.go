package jsonpointer

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPointerValueTraversal(t *testing.T) {
	type mapKey string
	type profile struct {
		Email string `json:"email"`
	}
	type user struct {
		Name    string  `json:"name"`
		Profile profile `json:"profile"`
	}

	doc := map[string]any{
		"users": []any{
			map[string]any{"name": "Alice"},
			user{Name: "Bob", Profile: profile{Email: "bob@example.com"}},
		},
		"typedSlice": []int{10, 20, 30},
		"typedMap":   map[mapKey]string{"one": "typed"},
		"array":      [2]string{"zero", "one"},
		"empty":      map[string]any{"": "empty-key"},
	}

	tests := []struct {
		name    string
		pointer string
		want    any
	}{
		{name: "map and slice", pointer: "/users/0/name", want: "Alice"},
		{name: "struct tag", pointer: "/users/1/profile/email", want: "bob@example.com"},
		{name: "typed slice", pointer: "/typedSlice/2", want: 30},
		{name: "typed map", pointer: "/typedMap/one", want: "typed"},
		{name: "array", pointer: "/array/1", want: "one"},
		{name: "empty token", pointer: "/empty/", want: "empty-key"},
		{name: "root", pointer: "", want: doc},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mustPointer(t, tt.pointer).Value(doc)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestOneShotValueAndReferenceOf(t *testing.T) {
	doc := map[string]any{
		"foo/bar": map[string]any{"tilde~key": "ready"},
	}

	got, err := Value(doc, "/foo~1bar/tilde~0key")
	require.NoError(t, err)
	assert.Equal(t, "ready", got)

	ref, err := ReferenceOf(doc, "/foo~1bar/tilde~0key")
	require.NoError(t, err)
	assert.Equal(t, "ready", ref.Value())
	assert.Equal(t, "tilde~key", ref.Token())
	assert.Equal(t, "/foo~1bar/tilde~0key", ref.Pointer().String())

	parent, ok := ref.Parent()
	require.True(t, ok)
	assert.Equal(t, map[string]any{"tilde~key": "ready"}, parent)

	_, err = Value(doc, "/~2")
	require.ErrorIs(t, err, ErrInvalidPointer)
}

func TestPointerReferenceRoot(t *testing.T) {
	doc := map[string]any{"name": "Alice"}
	ref, err := Root().Reference(doc)
	require.NoError(t, err)

	assert.Equal(t, doc, ref.Value())
	assert.Equal(t, "", ref.Token())
	assert.True(t, ref.Pointer().IsRoot())

	parent, ok := ref.Parent()
	assert.False(t, ok)
	assert.Nil(t, parent)
}

func TestPointerValueErrors(t *testing.T) {
	type user struct {
		Name string `json:"name"`
	}

	var nilUser *user
	tests := []struct {
		name    string
		doc     any
		pointer string
		wantErr error
		token   string
		depth   int
	}{
		{name: "missing key", doc: map[string]any{}, pointer: "/missing", wantErr: ErrKeyNotFound, token: "missing", depth: 0},
		{name: "invalid index", doc: []any{"zero"}, pointer: "/abc", wantErr: ErrInvalidIndex, token: "abc", depth: 0},
		{name: "leading zero index", doc: []any{"zero"}, pointer: "/01", wantErr: ErrInvalidIndex, token: "01", depth: 0},
		{name: "array end marker", doc: []any{"zero"}, pointer: "/-", wantErr: ErrIndexOutOfBounds, token: "-", depth: 0},
		{name: "out of bounds", doc: []any{"zero"}, pointer: "/1", wantErr: ErrIndexOutOfBounds, token: "1", depth: 0},
		{name: "missing field", doc: user{}, pointer: "/email", wantErr: ErrFieldNotFound, token: "email", depth: 0},
		{name: "nil pointer", doc: nilUser, pointer: "/name", wantErr: ErrNilPointer, token: "name", depth: 0},
		{name: "scalar", doc: "value", pointer: "/name", wantErr: ErrNotFound, token: "name", depth: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := mustPointer(t, tt.pointer)
			_, err := p.Value(tt.doc)
			require.ErrorIs(t, err, tt.wantErr)

			var pointerErr *Error
			require.True(t, errors.As(err, &pointerErr))
			assert.Equal(t, tt.token, pointerErr.Token())
			assert.Equal(t, tt.depth, pointerErr.Depth())
			assert.Equal(t, p.String(), pointerErr.Pointer().String())
		})
	}
}
