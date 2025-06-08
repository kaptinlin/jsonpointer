package jsonpointer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFind tests the find function comprehensively.
// Maps to: find.spec.ts + testFindRef.ts
func TestFind(t *testing.T) {
	t.Run("can find number root", func(t *testing.T) {
		res, err := Find(123)
		assert.NoError(t, err)
		assert.Equal(t, 123, res.Val)
	})

	t.Run("can find string root", func(t *testing.T) {
		res, err := Find("foo")
		assert.NoError(t, err)
		assert.Equal(t, "foo", res.Val)
	})

	t.Run("can find key in object", func(t *testing.T) {
		data := map[string]any{"foo": "bar"}
		res, err := Find(data, "foo")
		assert.NoError(t, err)
		assert.Equal(t, "bar", res.Val)
	})

	t.Run("returns container object and key", func(t *testing.T) {
		data := map[string]any{
			"foo": map[string]any{
				"bar": map[string]any{
					"baz": "qux",
					"a":   1,
				},
			},
		}
		res, err := Find(data, "foo", "bar", "baz")
		assert.NoError(t, err)

		expected := &Reference{
			Val: "qux",
			Obj: map[string]any{"baz": "qux", "a": 1},
			Key: "baz",
		}
		assert.Equal(t, expected.Val, res.Val)
		assert.Equal(t, expected.Key, res.Key)
		// Check object content without exact map comparison
		objMap, ok := res.Obj.(map[string]any)
		assert.True(t, ok)
		assert.Equal(t, "qux", objMap["baz"])
		assert.Equal(t, 1, objMap["a"])
	})

	t.Run("can reference simple object key", func(t *testing.T) {
		doc := map[string]any{"a": 123}
		res, err := Find(doc, "a")
		assert.NoError(t, err)

		assert.Equal(t, 123, res.Val)
		assert.Equal(t, "a", res.Key)
		assert.Equal(t, doc, res.Obj)
	})

	t.Run("throws when referencing missing key with multiple steps", func(t *testing.T) {
		doc := map[string]any{"a": 123}
		_, err := Find(doc, "b", "c")
		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("can reference array element", func(t *testing.T) {
		doc := map[string]any{
			"a": map[string]any{
				"b": []any{1, 2, 3},
			},
		}
		res, err := Find(doc, "a", "b", 1)
		assert.NoError(t, err)

		assert.Equal(t, 2, res.Val)
		assert.Equal(t, 1, res.Key)
		assert.Equal(t, []any{1, 2, 3}, res.Obj)
	})

	t.Run("can reference end of array", func(t *testing.T) {
		doc := map[string]any{
			"a": map[string]any{
				"b": []any{1, 2, 3},
			},
		}
		// path := ParseJsonPointer("/a/b/-")
		ref, err := Find(doc, "a", "b", "-")
		assert.NoError(t, err)

		assert.Nil(t, ref.Val) // undefined in TypeScript
		assert.Equal(t, []any{1, 2, 3}, ref.Obj)
		assert.Equal(t, 3, ref.Key)

		// Test type guards
		assert.True(t, IsArrayReference(*ref))
		if IsArrayReference(*ref) {
			// In TypeScript this would be checked with generic types,
			// but in Go we work with the general Reference type
			arrayObj, ok := ref.Obj.([]any)
			assert.True(t, ok)
			keyInt, ok := ref.Key.(int)
			assert.True(t, ok)
			assert.Equal(t, len(arrayObj), keyInt) // isArrayEnd equivalent
		}
	})

	t.Run("throws when pointing past array boundary", func(t *testing.T) {
		doc := map[string]any{
			"a": map[string]any{
				"b": []any{1, 2, 3},
			},
		}
		// path := ParseJsonPointer("/a/b/-1")
		_, err := Find(doc, "a", "b", "-1")
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidIndex, err)
	})

	t.Run("can point one element past array boundary", func(t *testing.T) {
		doc := map[string]any{
			"a": map[string]any{
				"b": []any{1, 2, 3},
			},
		}
		// path := ParseJsonPointer("/a/b/3")
		ref, err := Find(doc, "a", "b", 3)
		assert.NoError(t, err)

		assert.Nil(t, ref.Val) // undefined in TypeScript
		assert.Equal(t, []any{1, 2, 3}, ref.Obj)
		assert.Equal(t, 3, ref.Key)

		// Test type guards
		assert.True(t, IsArrayReference(*ref))
		if IsArrayReference(*ref) {
			arrayObj, ok := ref.Obj.([]any)
			assert.True(t, ok)
			keyInt, ok := ref.Key.(int)
			assert.True(t, ok)
			assert.Equal(t, len(arrayObj), keyInt) // isArrayEnd equivalent
		}
	})

	t.Run("can reference missing object key", func(t *testing.T) {
		doc := map[string]any{"foo": 123}
		// path := ParseJsonPointer("/bar")
		ref, err := Find(doc, "bar")
		assert.NoError(t, err)

		assert.Nil(t, ref.Val) // undefined in TypeScript
		assert.Equal(t, doc, ref.Obj)
		assert.Equal(t, "bar", ref.Key)
	})

	t.Run("can reference missing array key within bounds", func(t *testing.T) {
		doc := map[string]any{
			"foo": 123,
			"bar": []any{1, 2, 3},
		}
		// path := ParseJsonPointer("/bar/3")
		ref, err := Find(doc, "bar", 3)
		assert.NoError(t, err)

		assert.Nil(t, ref.Val) // undefined in TypeScript
		assert.Equal(t, []any{1, 2, 3}, ref.Obj)
		assert.Equal(t, 3, ref.Key)
	})
}

// TestFindByPointer tests the optimized findByPointer function.
func TestFindByPointer(t *testing.T) {
	t.Run("works with basic object", func(t *testing.T) {
		doc := map[string]any{"foo": "bar"}
		ref, err := FindByPointer(doc, "/foo")
		assert.NoError(t, err)
		assert.Equal(t, "bar", ref.Val)
		assert.Equal(t, "foo", ref.Key)
		assert.Equal(t, doc, ref.Obj)
	})

	t.Run("works with nested object", func(t *testing.T) {
		doc := map[string]any{
			"users": []any{
				map[string]any{"name": "Alice", "age": 30},
				map[string]any{"name": "Bob", "age": 25},
			},
		}
		ref, err := FindByPointer(doc, "/users/0/name")
		assert.NoError(t, err)
		assert.Equal(t, "Alice", ref.Val)
		assert.Equal(t, "name", ref.Key)
	})

	t.Run("handles array end marker", func(t *testing.T) {
		doc := map[string]any{"arr": []any{1, 2, 3}}
		ref, err := FindByPointer(doc, "/arr/-")
		assert.NoError(t, err)
		assert.Nil(t, ref.Val)
		assert.Equal(t, 3, ref.Key)
		assert.Equal(t, []any{1, 2, 3}, ref.Obj)
	})

	t.Run("throws for invalid array index", func(t *testing.T) {
		doc := map[string]any{"arr": []any{1, 2, 3}}
		_, err := FindByPointer(doc, "/arr/abc")
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidIndex, err)
	})

	t.Run("throws for not found", func(t *testing.T) {
		doc := map[string]any{"foo": "bar"}
		_, err := FindByPointer(doc, "/foo/bar")
		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("handles escaped characters", func(t *testing.T) {
		doc := map[string]any{"foo/bar": "value", "foo~bar": "value2"}
		ref1, err := FindByPointer(doc, "/foo~1bar")
		assert.NoError(t, err)
		assert.Equal(t, "value", ref1.Val)

		ref2, err := FindByPointer(doc, "/foo~0bar")
		assert.NoError(t, err)
		assert.Equal(t, "value2", ref2.Val)
	})
}

// TestGet tests the get function that never throws errors.
func TestGet(t *testing.T) {
	t.Run("basic object access", func(t *testing.T) {
		doc := map[string]any{"foo": "bar"}
		val := Get(doc, "foo")
		assert.Equal(t, "bar", val)
	})

	t.Run("missing key returns nil", func(t *testing.T) {
		doc := map[string]any{"foo": "bar"}
		val := Get(doc, "missing")
		assert.Nil(t, val)
	})

	t.Run("array access", func(t *testing.T) {
		doc := []any{1, 2, 3}
		val := Get(doc, 1)
		assert.Equal(t, 2, val)
	})

	t.Run("invalid array index returns nil", func(t *testing.T) {
		doc := []any{1, 2, 3}
		val := Get(doc, 5)
		assert.Nil(t, val)
	})

	t.Run("array end marker returns nil", func(t *testing.T) {
		doc := []any{1, 2, 3}
		val := Get(doc, "-")
		assert.Nil(t, val)
	})

	t.Run("nested access", func(t *testing.T) {
		doc := map[string]any{
			"users": []any{
				map[string]any{"name": "Alice"},
			},
		}
		val := Get(doc, "users", 0, "name")
		assert.Equal(t, "Alice", val)
	})
}
