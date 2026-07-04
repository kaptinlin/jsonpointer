package jsonpointer

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPointerValueTraversal(t *testing.T) {
	type mapKey string

	mapPointer := map[string]any{"name": "Pointer"}
	doc := map[string]any{
		"users": []any{
			map[string]any{"name": "Alice"},
		},
		"typedSlice": []int{10, 20, 30},
		"typedMap":   map[mapKey]string{"one": "typed"},
		"array":      [2]string{"zero", "one"},
		"mapPointer": &mapPointer,
		"empty":      map[string]any{"": "empty-key"},
		"dash":       map[string]any{"-": "dash-key"},
	}

	tests := []struct {
		name    string
		pointer string
		want    any
	}{
		{name: "map and slice", pointer: "/users/0/name", want: "Alice"},
		{name: "typed slice", pointer: "/typedSlice/2", want: 30},
		{name: "typed map", pointer: "/typedMap/one", want: "typed"},
		{name: "array", pointer: "/array/1", want: "one"},
		{name: "pointer to map", pointer: "/mapPointer/name", want: "Pointer"},
		{name: "empty token", pointer: "/empty/", want: "empty-key"},
		{name: "dash map key", pointer: "/dash/-", want: "dash-key"},
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

func TestPointerReferenceParentDereferencesMapPointer(t *testing.T) {
	profile := map[string]any{"name": "Ada"}
	doc := map[string]any{"profile": &profile}

	ref, err := FromTokens("profile", "name").Reference(doc)
	require.NoError(t, err)
	assert.Equal(t, "Ada", ref.Value())
	assert.Equal(t, "name", ref.Token())
	assert.Equal(t, "/profile/name", ref.Pointer().String())

	parent, ok := ref.Parent()
	require.True(t, ok)
	assert.IsType(t, map[string]any{}, parent)
	assert.Equal(t, profile, parent)
}

func TestPointerReferenceParentUsesContainerThatConsumesFinalToken(t *testing.T) {
	pointerSlice := []any{"zero", "one"}
	typedSlice := []string{"zero", "one"}
	interfaceWrapped := any([]any{"zero", "one"})

	tests := []struct {
		name       string
		doc        any
		pointer    Pointer
		wantValue  any
		wantParent any
	}{
		{
			name:       "slice pointer",
			doc:        map[string]any{"items": &pointerSlice},
			pointer:    FromTokens("items", "1"),
			wantValue:  "one",
			wantParent: pointerSlice,
		},
		{
			name:       "typed slice",
			doc:        map[string]any{"items": typedSlice},
			pointer:    FromTokens("items", "1"),
			wantValue:  "one",
			wantParent: typedSlice,
		},
		{
			name:       "interface wrapped container",
			doc:        map[string]any{"items": interfaceWrapped},
			pointer:    FromTokens("items", "1"),
			wantValue:  "one",
			wantParent: []any{"zero", "one"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ref, err := tt.pointer.Reference(tt.doc)
			require.NoError(t, err)
			assert.Equal(t, tt.wantValue, ref.Value())
			assert.Equal(t, "1", ref.Token())
			assert.Equal(t, tt.pointer.String(), ref.Pointer().String())

			parent, ok := ref.Parent()
			require.True(t, ok)
			assert.IsType(t, tt.wantParent, parent)
			assert.Equal(t, tt.wantParent, parent)
		})
	}
}

func TestPointerValueErrors(t *testing.T) {
	type user struct {
		Name string
	}

	var nilUser *user
	hugeIndex := "9999999999999999999999999999999999999999"
	nonASCIIIndex := "\uFF11"
	tests := []struct {
		name    string
		doc     any
		pointer string
		wantErr error
		token   string
		depth   int
	}{
		{name: "missing key", doc: map[string]any{}, pointer: "/missing", wantErr: ErrKeyNotFound, token: "missing", depth: 0},
		{name: "empty index", doc: []any{"zero"}, pointer: "/", wantErr: ErrInvalidIndex, token: "", depth: 0},
		{name: "invalid index", doc: []any{"zero"}, pointer: "/abc", wantErr: ErrInvalidIndex, token: "abc", depth: 0},
		{name: "leading zero index", doc: []any{"zero"}, pointer: "/01", wantErr: ErrInvalidIndex, token: "01", depth: 0},
		{name: "signed index", doc: []any{"zero"}, pointer: "/+1", wantErr: ErrInvalidIndex, token: "+1", depth: 0},
		{name: "decimal index", doc: []any{"zero"}, pointer: "/1.2", wantErr: ErrInvalidIndex, token: "1.2", depth: 0},
		{name: "non-ascii index", doc: []any{"zero"}, pointer: "/" + nonASCIIIndex, wantErr: ErrInvalidIndex, token: nonASCIIIndex, depth: 0},
		{name: "array end marker", doc: []any{"zero"}, pointer: "/-", wantErr: ErrIndexOutOfBounds, token: "-", depth: 0},
		{name: "out of bounds", doc: []any{"zero"}, pointer: "/1", wantErr: ErrIndexOutOfBounds, token: "1", depth: 0},
		{name: "overflow index", doc: []any{"zero"}, pointer: "/" + hugeIndex, wantErr: ErrIndexOutOfBounds, token: hugeIndex, depth: 0},
		{name: "struct", doc: user{Name: "Ada"}, pointer: "/Name", wantErr: ErrNotFound, token: "Name", depth: 0},
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
