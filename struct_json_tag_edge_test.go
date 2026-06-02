package jsonpointer

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStructJSONDashTagName(t *testing.T) {
	t.Parallel()

	docType := reflect.StructOf([]reflect.StructField{
		{Name: "Hidden", Type: reflect.TypeFor[string](), Tag: `json:"-"`},
		{Name: "Dash", Type: reflect.TypeFor[string](), Tag: `json:"-,"`},
	})
	value := reflect.New(docType).Elem()
	value.Field(0).SetString("hidden")
	value.Field(1).SetString("dash")
	doc := value.Interface()

	got, err := Get(doc, "-")
	require.NoError(t, err)
	assert.Equal(t, "dash", got)

	ref, err := FindByPointer(doc, "/-")
	require.NoError(t, err)
	assert.Equal(t, "dash", ref.Val)
	assert.Equal(t, "-", ref.Key)

	_, err = Get(doc, "Hidden")
	assert.ErrorIs(t, err, ErrFieldNotFound)
}
