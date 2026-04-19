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
		assert.ErrorIs(t, err, ErrPointerInvalid)
	})

	t.Run("invalid pointer too long", func(t *testing.T) {
		// Create a pointer longer than 1024 characters
		longPointer := "/" + strings.Repeat("a", 1024)
		err := Validate(longPointer)
		assert.ErrorIs(t, err, ErrPointerTooLong)
	})

	t.Run("valid pointer exactly 1024 characters", func(t *testing.T) {
		// Create a pointer exactly 1024 characters (including leading slash)
		exactPointer := "/" + strings.Repeat("a", 1023)
		err := Validate(exactPointer)
		assert.NoError(t, err)
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

	t.Run("valid path with string numbers", func(t *testing.T) {
		err := ValidatePath(Path{"0", "1", "2"})
		assert.NoError(t, err)
	})

	t.Run("valid path with mixed string types", func(t *testing.T) {
		err := ValidatePath(Path{"foo", "0", "bar", "1"})
		assert.NoError(t, err)
	})

	t.Run("invalid path - too long", func(t *testing.T) {
		// Create a path with more than 256 elements
		longPath := make(Path, 257)
		for i := range longPath {
			longPath[i] = "step"
		}
		err := ValidatePath(longPath)
		assert.ErrorIs(t, err, ErrPathTooLong)
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
		complexPath := Path{"users", "0", "profile", "settings", "notifications", "email", "enabled"}
		err := ValidatePath(complexPath)
		assert.NoError(t, err)
	})
}
