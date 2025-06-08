package jsonpointer

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestValidate tests JSON Pointer string validation.
func TestValidate(t *testing.T) {
	t.Run("valid empty string", func(t *testing.T) {
		err := Validate("")
		assert.NoError(t, err)
	})

	t.Run("valid root pointer", func(t *testing.T) {
		err := Validate("/")
		assert.NoError(t, err)
	})

	t.Run("valid simple pointer", func(t *testing.T) {
		err := Validate("/foo")
		assert.NoError(t, err)
	})

	t.Run("valid nested pointer", func(t *testing.T) {
		err := Validate("/foo/bar/baz")
		assert.NoError(t, err)
	})

	t.Run("valid pointer with escaped characters", func(t *testing.T) {
		err := Validate("/foo~0bar/baz~1qux")
		assert.NoError(t, err)
	})

	t.Run("invalid pointer without leading slash", func(t *testing.T) {
		err := Validate("foo/bar")
		assert.Error(t, err)
		assert.Equal(t, "pointer invalid", err.Error())
	})

	t.Run("invalid pointer too long", func(t *testing.T) {
		// Create a pointer longer than 1024 characters
		longPointer := "/" + strings.Repeat("a", 1024)
		err := Validate(longPointer)
		assert.Error(t, err)
		assert.Equal(t, "pointer too long", err.Error())
	})

	t.Run("valid pointer exactly 1024 characters", func(t *testing.T) {
		// Create a pointer exactly 1024 characters (including leading slash)
		exactPointer := "/" + strings.Repeat("a", 1023)
		err := Validate(exactPointer)
		assert.NoError(t, err)
	})

	t.Run("validates path when not string", func(t *testing.T) {
		// Valid path
		err := Validate(Path{"foo", "bar"})
		assert.NoError(t, err)

		// Invalid path (not a slice)
		err = Validate(123)
		assert.Error(t, err)
		assert.Equal(t, "invalid path", err.Error())
	})
}

// TestValidatePath tests path array validation.
func TestValidatePath(t *testing.T) {
	t.Run("valid empty path", func(t *testing.T) {
		err := ValidatePath(Path{})
		assert.NoError(t, err)
	})

	t.Run("valid path with strings", func(t *testing.T) {
		err := ValidatePath(Path{"foo", "bar", "baz"})
		assert.NoError(t, err)
	})

	t.Run("valid path with numbers", func(t *testing.T) {
		err := ValidatePath(Path{0, 1, 2})
		assert.NoError(t, err)
	})

	t.Run("valid path with mixed types", func(t *testing.T) {
		err := ValidatePath(Path{"foo", 0, "bar", 1})
		assert.NoError(t, err)
	})

	t.Run("valid path with different number types", func(t *testing.T) {
		err := ValidatePath(Path{
			int(1),
			int8(2),
			int16(3),
			int32(4),
			int64(5),
			uint(6),
			uint8(7),
			uint16(8),
			uint32(9),
			uint64(10),
			float32(11.0),
			float64(12.0),
		})
		assert.NoError(t, err)
	})

	t.Run("invalid path - not a slice", func(t *testing.T) {
		err := ValidatePath("not a slice")
		assert.Error(t, err)
		assert.Equal(t, "invalid path", err.Error())
	})

	t.Run("invalid path - not a slice (number)", func(t *testing.T) {
		err := ValidatePath(123)
		assert.Error(t, err)
		assert.Equal(t, "invalid path", err.Error())
	})

	t.Run("invalid path - not a slice (map)", func(t *testing.T) {
		err := ValidatePath(map[string]any{"foo": "bar"})
		assert.Error(t, err)
		assert.Equal(t, "invalid path", err.Error())
	})

	t.Run("invalid path - too long", func(t *testing.T) {
		// Create a path with more than 256 elements
		longPath := make(Path, 257)
		for i := range longPath {
			longPath[i] = "step"
		}
		err := ValidatePath(longPath)
		assert.Error(t, err)
		assert.Equal(t, "path too long", err.Error())
	})

	t.Run("valid path - exactly 256 elements", func(t *testing.T) {
		// Create a path with exactly 256 elements
		exactPath := make(Path, 256)
		for i := range exactPath {
			exactPath[i] = "step"
		}
		err := ValidatePath(exactPath)
		assert.NoError(t, err)
	})

	t.Run("invalid path step - boolean", func(t *testing.T) {
		err := ValidatePath(Path{"foo", true, "bar"})
		assert.Error(t, err)
		assert.Equal(t, "invalid path step", err.Error())
	})

	t.Run("invalid path step - nil", func(t *testing.T) {
		err := ValidatePath(Path{"foo", nil, "bar"})
		assert.Error(t, err)
		assert.Equal(t, "invalid path step", err.Error())
	})

	t.Run("invalid path step - slice", func(t *testing.T) {
		err := ValidatePath(Path{"foo", []string{"nested"}, "bar"})
		assert.Error(t, err)
		assert.Equal(t, "invalid path step", err.Error())
	})

	t.Run("invalid path step - map", func(t *testing.T) {
		err := ValidatePath(Path{"foo", map[string]any{"nested": "value"}, "bar"})
		assert.Error(t, err)
		assert.Equal(t, "invalid path step", err.Error())
	})

	t.Run("works with regular slice", func(t *testing.T) {
		// Test with []any slice
		regularSlice := []any{"foo", "bar", 0, 1}
		err := ValidatePath(regularSlice)
		assert.NoError(t, err)
	})

	t.Run("works with string slice", func(t *testing.T) {
		// Test with []string slice
		stringSlice := []string{"foo", "bar", "baz"}
		err := ValidatePath(stringSlice)
		assert.NoError(t, err)
	})

	t.Run("works with int slice", func(t *testing.T) {
		// Test with []int slice
		intSlice := []int{0, 1, 2, 3}
		err := ValidatePath(intSlice)
		assert.NoError(t, err)
	})
}

// TestValidateEdgeCases tests edge cases and integration scenarios.
func TestValidateEdgeCases(t *testing.T) {
	t.Run("validate pointer with unicode characters", func(t *testing.T) {
		err := Validate("/café/naïve/résumé")
		assert.NoError(t, err)
	})

	t.Run("validate path with unicode strings", func(t *testing.T) {
		err := ValidatePath(Path{"café", "naïve", "résumé"})
		assert.NoError(t, err)
	})

	t.Run("validate pointer with numbers as strings", func(t *testing.T) {
		err := Validate("/0/1/2")
		assert.NoError(t, err)
	})

	t.Run("validate path with string numbers", func(t *testing.T) {
		err := ValidatePath(Path{"0", "1", "2"})
		assert.NoError(t, err)
	})

	t.Run("validate complex nested pointer", func(t *testing.T) {
		complexPointer := "/users/0/profile/settings/notifications/email/enabled"
		err := Validate(complexPointer)
		assert.NoError(t, err)
	})

	t.Run("validate equivalent complex path", func(t *testing.T) {
		complexPath := Path{"users", 0, "profile", "settings", "notifications", "email", "enabled"}
		err := ValidatePath(complexPath)
		assert.NoError(t, err)
	})
}
