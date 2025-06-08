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
//	path := jsonpointer.Parse("/users/0/name")
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
//	err = jsonpointer.Validate("/users/0/name")
package jsonpointer

// Get retrieves a value from document using path components (never returns errors, returns nil for not found).
func Get(doc any, path ...any) any {
	if len(path) == 0 {
		return doc
	}
	return get(doc, Path(path))
}

// Find locates a reference in document using path components (returns errors for invalid operations).
func Find(doc any, path ...any) (*Reference, error) {
	if len(path) == 0 {
		return &Reference{Val: doc}, nil
	}
	return find(doc, Path(path))
}

// GetByPointer retrieves a value from document using JSON Pointer string (never returns errors).
func GetByPointer(doc any, pointer string) any {
	path := Parse(pointer)
	return get(doc, path)
}

// FindByPointer locates a reference in document using JSON Pointer string.
func FindByPointer(doc any, pointer string) (*Reference, error) {
	return findByPointer(pointer, doc)
}

// Parse parses a JSON Pointer string to a path array.
func Parse(pointer string) Path {
	return parseJsonPointer(pointer)
}

// Format formats path components into a JSON Pointer string.
func Format(path ...any) string {
	return formatJsonPointer(Path(path))
}

// Escape escapes special characters in a path component.
func Escape(component string) string {
	return escapeComponent(component)
}

// Unescape unescapes special characters in a path component.
func Unescape(component string) string {
	return unescapeComponent(component)
}

// Validate validates a JSON Pointer string or Path.
func Validate(pointer any) error {
	return validateJsonPointer(pointer)
}

// ValidatePath validates a path array.
func ValidatePath(path any) error {
	return validatePath(path)
}
