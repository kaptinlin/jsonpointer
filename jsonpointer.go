// Package jsonpointer provides JSON Pointer (RFC 6901) implementation for Go.
// This is a direct port of the TypeScript json-pointer library with identical behavior,
// using modern Go generics for type safety and performance.
//
// This package implements helper functions for JSON Pointer (RFC 6901) specification.
// https://tools.ietf.org/html/rfc6901
//
// TypeScript original source: https://github.com/jsonjoy-com/json-pointer
//
// Usage:
//
//	import "github.com/kaptinlin/jsonpointer"
//
//	// Parse JSON Pointer string to path
//	path := jsonpointer.ParseJsonPointer("/users/0/name")
//
//	// Find value with error handling
//	ref, err := jsonpointer.Find(data, path)
//	if err != nil {
//		// Handle error
//	}
//
//	// Get value without errors (returns nil for not found)
//	value := jsonpointer.Get(data, path)
//
//	// Validate JSON Pointer
//	err = jsonpointer.ValidateJsonPointer("/users/0/name")
//
// Breaking Change Notice:
// This version is a complete rewrite using modern Go generics with zero backward compatibility.
// All function signatures use 'any' instead of 'interface{}' and follow TypeScript API exactly.
package jsonpointer

// IsArrayReference checks if a Reference points to an array element.
func IsArrayReference(ref Reference) bool {
	return isArrayReference(ref)
}

// IsArrayEnd checks if an array reference points to the end of the array.
func IsArrayEnd[T any](ref ArrayReference[T]) bool {
	return isArrayEnd(ref)
}

// IsObjectReference checks if a Reference points to an object property.
func IsObjectReference(ref Reference) bool {
	return isObjectReference(ref)
}

// ValidateJsonPointer validates a JSON Pointer string or Path.
func ValidateJsonPointer(pointer any) error {
	return validateJsonPointer(pointer)
}

// ValidatePath validates a path array.
func ValidatePath(path any) error {
	return validatePath(path)
}

// Get retrieves a value from object using path (never returns errors, returns nil for not found).
func Get(val any, path Path) any {
	return get(val, path)
}

// Find locates a reference in object using path (returns errors for invalid operations).
func Find(val any, path Path) (*Reference, error) {
	return find(val, path)
}

// FindByPointer optimized find operation using direct string parsing.
func FindByPointer(pointer string, val any) (*Reference, error) {
	return findByPointer(pointer, val)
}
