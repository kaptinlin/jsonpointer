package jsonpointer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetByPointer(t *testing.T) {
	t.Parallel()

	t.Run("returns root value for empty pointer", func(t *testing.T) {
		t.Parallel()

		doc := map[string]any{"foo": "bar"}

		val, err := GetByPointer(doc, "")
		require.NoError(t, err)

		got, ok := val.(map[string]any)
		require.True(t, ok)
		assert.Equal(t, "bar", got["foo"])
	})

	t.Run("returns nested value for escaped pointer", func(t *testing.T) {
		t.Parallel()

		doc := map[string]any{
			"foo/bar": map[string]any{
				"tilde~key": "value",
			},
		}

		val, err := GetByPointer(doc, "/foo~1bar/tilde~0key")
		require.NoError(t, err)
		assert.Equal(t, "value", val)
	})

	t.Run("returns key not found for pointer without leading slash", func(t *testing.T) {
		t.Parallel()

		val, err := GetByPointer(map[string]any{"foo": "bar"}, "foo")
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrKeyNotFound)
		assert.Nil(t, val)
	})
}

func TestReferenceTypePredicates(t *testing.T) {
	t.Parallel()

	t.Run("identifies array references", func(t *testing.T) {
		t.Parallel()

		ref := Reference{Val: "value", Obj: []any{"value"}, Key: "0"}
		assert.True(t, IsArrayReference(ref))
		assert.False(t, IsObjectReference(ref))
	})

	t.Run("rejects non numeric array keys", func(t *testing.T) {
		t.Parallel()

		ref := Reference{Val: "value", Obj: []any{"value"}, Key: "first"}
		assert.False(t, IsArrayReference(ref))
	})

	t.Run("identifies object references", func(t *testing.T) {
		t.Parallel()

		ref := Reference{Val: "value", Obj: map[string]any{"name": "value"}, Key: "name"}
		assert.True(t, IsObjectReference(ref))
		assert.False(t, IsArrayReference(ref))
	})

	t.Run("rejects references without parent context", func(t *testing.T) {
		t.Parallel()

		assert.False(t, IsArrayReference(Reference{Val: "value"}))
		assert.False(t, IsObjectReference(Reference{Val: "value"}))
	})
}

func TestPointerObjectTraversal(t *testing.T) {
	t.Parallel()

	t.Run("get supports pointer to map", func(t *testing.T) {
		t.Parallel()

		doc := map[string]any{"foo": "bar"}

		val, err := Get(&doc, "foo")
		require.NoError(t, err)
		assert.Equal(t, "bar", val)
	})

	t.Run("find supports pointer to map", func(t *testing.T) {
		t.Parallel()

		doc := map[string]any{"foo": "bar"}

		ref, err := Find(&doc, "foo")
		require.NoError(t, err)
		assert.Equal(t, "bar", ref.Val)
		assert.Equal(t, "foo", ref.Key)
		assert.Equal(t, &doc, ref.Obj)
	})

	t.Run("find by pointer supports pointer to map", func(t *testing.T) {
		t.Parallel()

		doc := map[string]any{"foo": "bar"}

		ref, err := FindByPointer(&doc, "/foo")
		require.NoError(t, err)
		assert.Equal(t, "bar", ref.Val)
		assert.Equal(t, "foo", ref.Key)
		assert.Equal(t, &doc, ref.Obj)
	})

	t.Run("get returns nil pointer error for nil map pointer", func(t *testing.T) {
		t.Parallel()

		var doc *map[string]any

		val, err := Get(doc, "foo")
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrNilPointer)
		assert.Nil(t, val)
	})

	t.Run("find returns nil pointer error for nil map pointer", func(t *testing.T) {
		t.Parallel()

		var doc *map[string]any

		ref, err := Find(doc, "foo")
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrNilPointer)
		assert.Nil(t, ref)
	})

	t.Run("find by pointer returns nil pointer error for nil map pointer", func(t *testing.T) {
		t.Parallel()

		var doc *map[string]any

		ref, err := FindByPointer(doc, "/foo")
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrNilPointer)
		assert.Nil(t, ref)
	})
}

type publicAPICoverageContainer struct {
	Items []int `json:"items"`
}

func TestCollectionTraversalCoverage(t *testing.T) {
	t.Parallel()

	t.Run("get supports typed slices through reflection", func(t *testing.T) {
		t.Parallel()

		doc := []string{"zero", "one", "two"}

		val, err := Get(doc, "1")
		require.NoError(t, err)
		assert.Equal(t, "one", val)
	})

	t.Run("find supports typed arrays through reflection", func(t *testing.T) {
		t.Parallel()

		doc := [2]string{"zero", "one"}

		ref, err := Find(doc, "1")
		require.NoError(t, err)
		assert.Equal(t, "one", ref.Val)
		assert.Equal(t, "1", ref.Key)

		got, ok := ref.Obj.([2]string)
		require.True(t, ok)
		assert.Equal(t, "zero", got[0])
		assert.Equal(t, "one", got[1])
	})

	t.Run("get returns invalid index for typed slices", func(t *testing.T) {
		t.Parallel()

		_, err := Get([]int{1, 2, 3}, "abc")
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidIndex)
	})

	t.Run("find returns field not found for missing struct field", func(t *testing.T) {
		t.Parallel()

		_, err := Find(publicAPICoverageContainer{Items: []int{1, 2, 3}}, "missing")
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrFieldNotFound)
	})
}
