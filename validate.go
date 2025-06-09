package jsonpointer

import (
	"reflect"
	"strings"
)

// validateJsonPointer validates a JSON Pointer string or Path.
// Returns an error if the pointer is invalid according to RFC 6901.
func validateJsonPointer(pointer any) error {
	switch p := pointer.(type) {
	case string:
		return validatePointerString(p)
	case Path:
		return validatePath(p)
	case []string:
		return validatePath(Path(p))
	default:
		return ErrPointerInvalid
	}
}

// validatePointerString validates a JSON Pointer string.
func validatePointerString(pointer string) error {
	// Empty string is valid (root pointer)
	if pointer == "" {
		return nil
	}

	// Must start with "/"
	if !strings.HasPrefix(pointer, "/") {
		return ErrPointerInvalid
	}

	// Check length limit
	if len(pointer) > 1024 {
		return ErrPointerTooLong
	}

	// Validate escape sequences
	for i := 0; i < len(pointer); i++ {
		if pointer[i] == '~' {
			if i+1 >= len(pointer) {
				return ErrPointerInvalid // Invalid escape at end
			}
			next := pointer[i+1]
			if next != '0' && next != '1' {
				return ErrPointerInvalid // Invalid escape sequence
			}
			i++ // Skip the next character
		}
	}

	return nil
}

// validatePath validates a path array.
// Returns an error if the path contains invalid components.
func validatePath(path any) error {
	// Check if path is a string slice
	val := reflect.ValueOf(path)
	if val.Kind() != reflect.Slice {
		return ErrInvalidPath
	}

	// Check length
	length := val.Len()
	if length > 256 {
		return ErrPathTooLong
	}

	// Validate each step - all must be strings
	for i := 0; i < length; i++ {
		step := val.Index(i).Interface()

		// Check if step is string
		if _, ok := step.(string); !ok {
			return ErrInvalidPathStep
		}
	}

	return nil
}
