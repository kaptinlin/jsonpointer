package jsonpointer

import (
	"testing"

	"github.com/google/go-cmp/cmp"
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

	t.Run("typed array reference preserves value and numeric key", func(t *testing.T) {
		t.Parallel()

		value := 2
		ref := ArrayReference[int]{Val: &value, Obj: []int{1, 2, 3}, Key: 1}

		assert.Equal(t, 2, *ref.Val)
		assert.Equal(t, 1, ref.Key)
		if diff := cmp.Diff([]int{1, 2, 3}, ref.Obj); diff != "" {
			t.Errorf("ArrayReference obj mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("typed object reference preserves value and key", func(t *testing.T) {
		t.Parallel()

		ref := ObjectReference[string]{Val: "value", Obj: map[string]string{"name": "value"}, Key: "name"}

		assert.Equal(t, "value", ref.Val)
		assert.Equal(t, "name", ref.Key)
		if diff := cmp.Diff(map[string]string{"name": "value"}, ref.Obj); diff != "" {
			t.Errorf("ObjectReference obj mismatch (-want +got):\n%s", diff)
		}
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
		require.NotNil(t, ref)
		assert.Equal(t, "bar", ref.Val)
		assert.Equal(t, "foo", ref.Key)

		gotObj, ok := ref.Obj.(*map[string]any)
		require.True(t, ok)
		if diff := cmp.Diff(&doc, gotObj); diff != "" {
			t.Errorf("Find() pointer parent mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("find by pointer supports pointer to map", func(t *testing.T) {
		t.Parallel()

		doc := map[string]any{"foo": "bar"}

		ref, err := FindByPointer(&doc, "/foo")
		require.NoError(t, err)
		require.NotNil(t, ref)
		assert.Equal(t, "bar", ref.Val)
		assert.Equal(t, "foo", ref.Key)

		gotObj, ok := ref.Obj.(*map[string]any)
		require.True(t, ok)
		if diff := cmp.Diff(&doc, gotObj); diff != "" {
			t.Errorf("FindByPointer() pointer parent mismatch (-want +got):\n%s", diff)
		}
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

type ignoredFieldCoverageContainer struct {
	Visible string `json:"visible"`
	Hidden  string `json:"-"`
}

func TestIgnoredStructFieldCoverage(t *testing.T) {
	t.Parallel()

	doc := ignoredFieldCoverageContainer{Visible: "ok", Hidden: "secret"}

	t.Run("get returns field not found for ignored struct field", func(t *testing.T) {
		t.Parallel()

		val, err := Get(doc, "Hidden")
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrFieldNotFound)
		assert.Nil(t, val)
	})

	t.Run("find returns field not found for ignored struct field", func(t *testing.T) {
		t.Parallel()

		ref, err := Find(doc, "Hidden")
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrFieldNotFound)
		assert.Nil(t, ref)
	})

	t.Run("find by pointer returns field not found for ignored struct field", func(t *testing.T) {
		t.Parallel()

		ref, err := FindByPointer(doc, "/Hidden")
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrFieldNotFound)
		assert.Nil(t, ref)
	})
}

func TestPointerSliceTraversal(t *testing.T) {
	t.Parallel()

	t.Run("find by pointer returns root reference for empty pointer", func(t *testing.T) {
		t.Parallel()

		doc := map[string]any{"foo": "bar"}

		ref, err := FindByPointer(doc, "")
		require.NoError(t, err)
		require.NotNil(t, ref)

		got, ok := ref.Val.(map[string]any)
		require.True(t, ok)
		assert.Equal(t, "bar", got["foo"])
		assert.Nil(t, ref.Obj)
		assert.Empty(t, ref.Key)
	})

	t.Run("get supports pointer to slice", func(t *testing.T) {
		t.Parallel()

		doc := []any{"zero", "one"}

		val, err := Get(&doc, "1")
		require.NoError(t, err)
		assert.Equal(t, "one", val)
	})

	t.Run("find supports pointer to slice", func(t *testing.T) {
		t.Parallel()

		doc := []any{"zero", "one"}

		ref, err := Find(&doc, "1")
		require.NoError(t, err)
		require.NotNil(t, ref)
		assert.Equal(t, "one", ref.Val)
		assert.Equal(t, "1", ref.Key)
		assert.Same(t, &doc, ref.Obj)
	})

	t.Run("find by pointer supports pointer to slice", func(t *testing.T) {
		t.Parallel()

		doc := []any{"zero", "one"}

		ref, err := FindByPointer(&doc, "/1")
		require.NoError(t, err)
		require.NotNil(t, ref)
		assert.Equal(t, "one", ref.Val)
		assert.Equal(t, "1", ref.Key)
		assert.Same(t, &doc, ref.Obj)
	})

	t.Run("get supports pointer to interface wrapped map", func(t *testing.T) {
		t.Parallel()

		var doc any = map[string]any{"foo": "bar"}

		val, err := Get(&doc, "foo")
		require.NoError(t, err)
		assert.Equal(t, "bar", val)
	})

	t.Run("get returns nil pointer error for nil slice pointer", func(t *testing.T) {
		t.Parallel()

		var doc *[]any

		val, err := Get(doc, "0")
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrNilPointer)
		assert.Nil(t, val)
	})

	t.Run("find returns nil pointer error for nil slice pointer", func(t *testing.T) {
		t.Parallel()

		var doc *[]any

		ref, err := Find(doc, "0")
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrNilPointer)
		assert.Nil(t, ref)
	})

	t.Run("find by pointer returns nil pointer error for nil slice pointer", func(t *testing.T) {
		t.Parallel()

		var doc *[]any

		ref, err := FindByPointer(doc, "/0")
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrNilPointer)
		assert.Nil(t, ref)
	})
}

type nestedPointerCoverageContainer struct {
	Child **publicAPICoverageContainer `json:"child"`
}

func TestNilPointerChainTraversal(t *testing.T) {
	t.Parallel()

	var child *publicAPICoverageContainer
	doc := nestedPointerCoverageContainer{Child: &child}

	t.Run("get returns nil pointer error inside pointer chain", func(t *testing.T) {
		t.Parallel()

		val, err := Get(doc, "child", "items")
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrNilPointer)
		assert.Nil(t, val)
	})

	t.Run("find returns nil pointer error inside pointer chain", func(t *testing.T) {
		t.Parallel()

		ref, err := Find(doc, "child", "items")
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrNilPointer)
		assert.Nil(t, ref)
	})

	t.Run("find by pointer returns nil pointer error inside pointer chain", func(t *testing.T) {
		t.Parallel()

		ref, err := FindByPointer(doc, "/child/items")
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrNilPointer)
		assert.Nil(t, ref)
	})
}

func TestNilValueTraversal(t *testing.T) {
	t.Parallel()

	doc := map[string]any{"foo": nil}

	t.Run("get returns not found when traversal continues through nil", func(t *testing.T) {
		t.Parallel()

		val, err := Get(doc, "foo", "bar")
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
		assert.Nil(t, val)
	})

	t.Run("find returns not found when traversal continues through nil", func(t *testing.T) {
		t.Parallel()

		ref, err := Find(doc, "foo", "bar")
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
		assert.Nil(t, ref)
	})

	t.Run("find by pointer returns not found when traversal continues through nil", func(t *testing.T) {
		t.Parallel()

		ref, err := FindByPointer(doc, "/foo/bar")
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
		assert.Nil(t, ref)
	})
}
