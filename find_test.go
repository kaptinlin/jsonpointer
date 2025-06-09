package jsonpointer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFind tests the find function.
func TestFind(t *testing.T) {
	t.Run("can find root", func(t *testing.T) {
		doc := map[string]any{"foo": "bar"}
		ref, err := Find(doc)
		assert.NoError(t, err)
		assert.Equal(t, doc, ref.Val)
		assert.Nil(t, ref.Obj)
		assert.Equal(t, "", ref.Key)
	})

	t.Run("can find object key", func(t *testing.T) {
		doc := map[string]any{"foo": "bar"}
		ref, err := Find(doc, "foo")
		assert.NoError(t, err)
		assert.Equal(t, "bar", ref.Val)
		assert.Equal(t, doc, ref.Obj)
		assert.Equal(t, "foo", ref.Key)
	})

	t.Run("can find nested object key", func(t *testing.T) {
		doc := map[string]any{
			"a": map[string]any{
				"b": map[string]any{
					"c": "value",
				},
			},
		}
		ref, err := Find(doc, "a", "b", "c")
		assert.NoError(t, err)
		assert.Equal(t, "value", ref.Val)
		assert.Equal(t, "c", ref.Key)
	})

	t.Run("can find array element", func(t *testing.T) {
		doc := map[string]any{
			"a": map[string]any{
				"b": []any{1, 2, 3},
			},
		}
		res, err := Find(doc, "a", "b", "1")
		assert.NoError(t, err)

		assert.Equal(t, 2, res.Val)
		assert.Equal(t, "1", res.Key)
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
		assert.Equal(t, "3", ref.Key) // Array length as string

		// Test type guards
		assert.True(t, IsArrayReference(*ref))
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
		ref, err := Find(doc, "a", "b", "3")
		assert.NoError(t, err)

		assert.Nil(t, ref.Val) // undefined in TypeScript
		assert.Equal(t, []any{1, 2, 3}, ref.Obj)
		assert.Equal(t, "3", ref.Key) // Index as string

		// Test type guards
		assert.True(t, IsArrayReference(*ref))
	})

	t.Run("throws for missing object key", func(t *testing.T) {
		doc := map[string]any{"foo": 123}
		// path := ParseJsonPointer("/bar")
		_, err := Find(doc, "bar")
		assert.Error(t, err)
		assert.Equal(t, ErrKeyNotFound, err)
	})

	t.Run("can reference missing array key within bounds", func(t *testing.T) {
		doc := map[string]any{
			"foo": 123,
			"bar": []any{1, 2, 3},
		}
		// path := ParseJsonPointer("/bar/3")
		ref, err := Find(doc, "bar", "3")
		assert.NoError(t, err)

		assert.Nil(t, ref.Val) // undefined in TypeScript
		assert.Equal(t, []any{1, 2, 3}, ref.Obj)
		assert.Equal(t, "3", ref.Key) // Index as string
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
		assert.Equal(t, "3", ref.Key) // Array length as string
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
		val, err := Get(doc, "foo")
		assert.NoError(t, err)
		assert.Equal(t, "bar", val)
	})

	t.Run("missing key returns error", func(t *testing.T) {
		doc := map[string]any{"foo": "bar"}
		val, err := Get(doc, "missing")
		assert.Error(t, err)
		assert.Equal(t, ErrKeyNotFound, err)
		assert.Nil(t, val)
	})

	t.Run("array access", func(t *testing.T) {
		doc := []any{1, 2, 3}
		val, err := Get(doc, "1")
		assert.NoError(t, err)
		assert.Equal(t, 2, val)
	})

	t.Run("invalid array index returns error", func(t *testing.T) {
		doc := []any{1, 2, 3}
		val, err := Get(doc, "5")
		assert.Error(t, err)
		assert.Nil(t, val)
	})

	t.Run("array end marker returns nil", func(t *testing.T) {
		doc := []any{1, 2, 3}
		val, err := Get(doc, "-")
		assert.NoError(t, err)
		assert.Nil(t, val)
	})

	t.Run("nested access", func(t *testing.T) {
		doc := map[string]any{
			"users": []any{
				map[string]any{"name": "Alice"},
			},
		}
		val, err := Get(doc, "users", "0", "name")
		assert.NoError(t, err)
		assert.Equal(t, "Alice", val)
	})
}
