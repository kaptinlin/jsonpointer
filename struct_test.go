package jsonpointer

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStructFieldRules(t *testing.T) {
	type sample struct {
		Name    string `json:"name"`
		Alias   string `json:",omitempty"`
		Hidden  string `json:"-"`
		private string
	}

	doc := sample{
		Name:    "Ada",
		Alias:   "Lovelace",
		Hidden:  "hidden",
		private: "private",
	}

	tests := []struct {
		name    string
		pointer string
		want    any
	}{
		{name: "tag name", pointer: "/name", want: "Ada"},
		{name: "empty tag name falls back", pointer: "/Alias", want: "Lovelace"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mustPointer(t, tt.pointer).Value(doc)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}

	_, err := mustPointer(t, "/Hidden").Value(doc)
	require.ErrorIs(t, err, ErrFieldNotFound)

	_, err = mustPointer(t, "/private").Value(doc)
	require.ErrorIs(t, err, ErrFieldNotFound)
}

func TestStructDashLiteralTag(t *testing.T) {
	typ := reflect.StructOf([]reflect.StructField{
		{
			Name: "Dash",
			Type: reflect.TypeFor[string](),
			Tag:  reflect.StructTag(`json:"-,"`),
		},
	})
	doc := reflect.New(typ).Elem()
	doc.Field(0).SetString("dash")

	got, err := mustPointer(t, "/-").Value(doc.Interface())
	require.NoError(t, err)
	assert.Equal(t, "dash", got)
}

func TestStructEmbeddedDominance(t *testing.T) {
	type embeddedA struct {
		Name string
	}
	type embeddedB struct {
		Name string
	}
	type conflict struct {
		embeddedA
		embeddedB
	}
	type topLevel struct {
		embeddedA
		Name string
	}
	type taggedDominates struct {
		Name  string
		Alias string `json:"Name"`
	}

	_, err := mustPointer(t, "/Name").Value(conflict{
		embeddedA: embeddedA{Name: "a"},
		embeddedB: embeddedB{Name: "b"},
	})
	require.ErrorIs(t, err, ErrFieldNotFound)

	got, err := mustPointer(t, "/Name").Value(topLevel{
		embeddedA: embeddedA{Name: "embedded"},
		Name:      "top",
	})
	require.NoError(t, err)
	assert.Equal(t, "top", got)

	got, err = mustPointer(t, "/Name").Value(taggedDominates{
		Name:  "field",
		Alias: "tagged",
	})
	require.NoError(t, err)
	assert.Equal(t, "tagged", got)
}

func TestStructNilEmbeddedPointer(t *testing.T) {
	type embedded struct {
		Name string `json:"name"`
	}
	type holder struct {
		*embedded
	}

	_, err := mustPointer(t, "/name").Value(holder{})
	require.ErrorIs(t, err, ErrNilPointer)
}
