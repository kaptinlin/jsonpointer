package jsonpointer

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFind(t *testing.T) {
	t.Parallel()

	t.Run("can find root", func(t *testing.T) {
		t.Parallel()

		doc := map[string]any{"foo": "bar"}

		ref, err := Find(doc)
		require.NoError(t, err)
		require.NotNil(t, ref)

		got, ok := ref.Val.(map[string]any)
		require.True(t, ok)
		if diff := cmp.Diff(doc, got); diff != "" {
			t.Errorf("Find() root mismatch (-want +got):\n%s", diff)
		}
		assert.Nil(t, ref.Obj)
		assert.Equal(t, "", ref.Key)
	})

	t.Run("can find object key", func(t *testing.T) {
		t.Parallel()

		doc := map[string]any{"foo": "bar"}

		ref, err := Find(doc, "foo")
		require.NoError(t, err)
		require.NotNil(t, ref)

		assert.Equal(t, "bar", ref.Val)
		gotObj, ok := ref.Obj.(map[string]any)
		require.True(t, ok)
		if diff := cmp.Diff(doc, gotObj); diff != "" {
			t.Errorf("Find() parent mismatch (-want +got):\n%s", diff)
		}
		assert.Equal(t, "foo", ref.Key)
	})

	t.Run("can find nested object key", func(t *testing.T) {
		t.Parallel()

		doc := map[string]any{
			"a": map[string]any{
				"b": map[string]any{
					"c": "value",
				},
			},
		}

		ref, err := Find(doc, "a", "b", "c")
		require.NoError(t, err)
		require.NotNil(t, ref)

		assert.Equal(t, "value", ref.Val)
		assert.Equal(t, "c", ref.Key)
	})

	t.Run("can find array element", func(t *testing.T) {
		t.Parallel()

		doc := map[string]any{
			"a": map[string]any{
				"b": []any{1, 2, 3},
			},
		}

		ref, err := Find(doc, "a", "b", "1")
		require.NoError(t, err)
		require.NotNil(t, ref)

		assert.Equal(t, 2, ref.Val)
		assert.Equal(t, "1", ref.Key)
		gotObj, ok := ref.Obj.([]any)
		require.True(t, ok)
		if diff := cmp.Diff([]any{1, 2, 3}, gotObj); diff != "" {
			t.Errorf("Find() array parent mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("array end marker returns error", func(t *testing.T) {
		t.Parallel()

		doc := map[string]any{
			"a": map[string]any{
				"b": []any{1, 2, 3},
			},
		}

		_, err := Find(doc, "a", "b", "-")
		assert.ErrorIs(t, err, ErrIndexOutOfBounds)
	})

	t.Run("negative array index returns error", func(t *testing.T) {
		t.Parallel()

		doc := map[string]any{
			"a": map[string]any{
				"b": []any{1, 2, 3},
			},
		}

		_, err := Find(doc, "a", "b", "-1")
		assert.ErrorIs(t, err, ErrInvalidIndex)
	})

	t.Run("index at array length returns error", func(t *testing.T) {
		t.Parallel()

		doc := map[string]any{
			"a": map[string]any{
				"b": []any{1, 2, 3},
			},
		}

		_, err := Find(doc, "a", "b", "3")
		assert.ErrorIs(t, err, ErrIndexOutOfBounds)
	})

	t.Run("throws for missing object key", func(t *testing.T) {
		t.Parallel()

		doc := map[string]any{"foo": 123}

		_, err := Find(doc, "bar")
		assert.ErrorIs(t, err, ErrKeyNotFound)
	})

	t.Run("array index at length returns error", func(t *testing.T) {
		t.Parallel()

		doc := map[string]any{
			"foo": 123,
			"bar": []any{1, 2, 3},
		}

		_, err := Find(doc, "bar", "3")
		assert.ErrorIs(t, err, ErrIndexOutOfBounds)
	})
}

func TestFindByPointer(t *testing.T) {
	t.Parallel()

	t.Run("works with basic object", func(t *testing.T) {
		t.Parallel()

		doc := map[string]any{"foo": "bar"}

		ref, err := FindByPointer(doc, "/foo")
		require.NoError(t, err)
		require.NotNil(t, ref)

		assert.Equal(t, "bar", ref.Val)
		assert.Equal(t, "foo", ref.Key)
		gotObj, ok := ref.Obj.(map[string]any)
		require.True(t, ok)
		if diff := cmp.Diff(doc, gotObj); diff != "" {
			t.Errorf("FindByPointer() parent mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("works with nested object", func(t *testing.T) {
		t.Parallel()

		doc := map[string]any{
			"users": []any{
				map[string]any{"name": "Alice", "age": 30},
				map[string]any{"name": "Bob", "age": 25},
			},
		}

		ref, err := FindByPointer(doc, "/users/0/name")
		require.NoError(t, err)
		require.NotNil(t, ref)

		assert.Equal(t, "Alice", ref.Val)
		assert.Equal(t, "name", ref.Key)
	})

	t.Run("array end marker returns error", func(t *testing.T) {
		t.Parallel()

		doc := map[string]any{"arr": []any{1, 2, 3}}

		_, err := FindByPointer(doc, "/arr/-")
		assert.ErrorIs(t, err, ErrIndexOutOfBounds)
	})

	t.Run("throws for invalid array index", func(t *testing.T) {
		t.Parallel()

		doc := map[string]any{"arr": []any{1, 2, 3}}

		_, err := FindByPointer(doc, "/arr/abc")
		assert.ErrorIs(t, err, ErrInvalidIndex)
	})

	t.Run("throws for missing map key", func(t *testing.T) {
		t.Parallel()

		doc := map[string]any{"foo": "bar"}

		_, err := FindByPointer(doc, "/missing")
		assert.ErrorIs(t, err, ErrKeyNotFound)
	})

	t.Run("throws for not found", func(t *testing.T) {
		t.Parallel()

		doc := map[string]any{"foo": "bar"}

		_, err := FindByPointer(doc, "/foo/bar")
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("handles escaped characters", func(t *testing.T) {
		t.Parallel()

		doc := map[string]any{"foo/bar": "value", "foo~bar": "value2"}

		ref1, err := FindByPointer(doc, "/foo~1bar")
		require.NoError(t, err)
		require.NotNil(t, ref1)
		assert.Equal(t, "value", ref1.Val)

		ref2, err := FindByPointer(doc, "/foo~0bar")
		require.NoError(t, err)
		require.NotNil(t, ref2)
		assert.Equal(t, "value2", ref2.Val)
	})

	t.Run("handles trailing empty pointer segments", func(t *testing.T) {
		t.Parallel()

		doc := map[string]any{
			"foo": map[string]any{
				"": map[string]any{
					"": "value",
				},
			},
		}

		ref, err := FindByPointer(doc, "/foo//")
		require.NoError(t, err)
		require.NotNil(t, ref)
		assert.Equal(t, "value", ref.Val)
		assert.Equal(t, "", ref.Key)
	})
}

func TestGet(t *testing.T) {
	t.Parallel()

	t.Run("empty path returns root value", func(t *testing.T) {
		t.Parallel()

		doc := map[string]any{"foo": "bar"}

		val, err := Get(doc)
		require.NoError(t, err)

		got, ok := val.(map[string]any)
		require.True(t, ok)
		if diff := cmp.Diff(doc, got); diff != "" {
			t.Errorf("Get() root mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("basic object access", func(t *testing.T) {
		t.Parallel()

		doc := map[string]any{"foo": "bar"}

		val, err := Get(doc, "foo")
		require.NoError(t, err)
		assert.Equal(t, "bar", val)
	})

	t.Run("missing key returns error", func(t *testing.T) {
		t.Parallel()

		doc := map[string]any{"foo": "bar"}

		val, err := Get(doc, "missing")
		assert.ErrorIs(t, err, ErrKeyNotFound)
		assert.Nil(t, val)
	})

	t.Run("array access", func(t *testing.T) {
		t.Parallel()

		doc := []any{1, 2, 3}

		val, err := Get(doc, "1")
		require.NoError(t, err)
		assert.Equal(t, 2, val)
	})

	t.Run("invalid array index returns error", func(t *testing.T) {
		t.Parallel()

		doc := []any{1, 2, 3}

		val, err := Get(doc, "5")
		assert.ErrorIs(t, err, ErrIndexOutOfBounds)
		assert.Nil(t, val)
	})

	t.Run("array end marker returns error", func(t *testing.T) {
		t.Parallel()

		doc := []any{1, 2, 3}

		_, err := Get(doc, "-")
		assert.ErrorIs(t, err, ErrIndexOutOfBounds)
	})

	t.Run("nested access", func(t *testing.T) {
		t.Parallel()

		doc := map[string]any{
			"users": []any{
				map[string]any{"name": "Alice"},
			},
		}

		val, err := Get(doc, "users", "0", "name")
		require.NoError(t, err)
		assert.Equal(t, "Alice", val)
	})
}

func TestReflectiveMapKeyHandling(t *testing.T) {
	t.Parallel()

	t.Run("get returns error for non-string map keys instead of panicking", func(t *testing.T) {
		t.Parallel()

		val, err := Get(map[int]string{1: "one"}, "1")
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
		assert.Nil(t, val)
	})

	t.Run("find returns error for non-string map keys instead of panicking", func(t *testing.T) {
		t.Parallel()

		ref, err := Find(map[int]string{1: "one"}, "1")
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
		assert.Nil(t, ref)
	})

	t.Run("find by pointer returns error for non-string map keys instead of panicking", func(t *testing.T) {
		t.Parallel()

		ref, err := FindByPointer(map[int]string{1: "one"}, "/1")
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
		assert.Nil(t, ref)
	})

	t.Run("string alias map keys remain accessible", func(t *testing.T) {
		t.Parallel()

		type alias string

		doc := map[alias]string{"name": "Alice"}

		val, err := Get(doc, "name")
		require.NoError(t, err)
		assert.Equal(t, "Alice", val)

		ref, err := Find(doc, "name")
		require.NoError(t, err)
		require.NotNil(t, ref)
		assert.Equal(t, "Alice", ref.Val)

		ref, err = FindByPointer(doc, "/name")
		require.NoError(t, err)
		require.NotNil(t, ref)
		assert.Equal(t, "Alice", ref.Val)
	})
}
