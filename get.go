package jsonpointer

import (
	"reflect"
	"strconv"
)

// fastGet implements ultra-fast path that avoids token allocation entirely.
// Optimized for string-only Path - direct access without intermediate token creation.
func fastGet(val any, step string) (any, bool) {
	switch v := val.(type) {
	case map[string]any:
		// Most common case: map[string]any - direct string key access
		result, exists := v[step]
		return result, exists

	case *map[string]any:
		// Pointer to map optimization
		if v == nil {
			return nil, false
		}
		result, exists := (*v)[step]
		return result, exists

	case []any:
		// Array access - parse index from string
		if step == "-" {
			return nil, false // array end marker
		}
		index := fastAtoi(step)
		if index < 0 || index >= len(v) {
			return nil, false // invalid or out of bounds
		}
		return v[index], true

	case *[]any:
		// Pointer to slice optimization
		if v == nil {
			return nil, false
		}
		if step == "-" {
			return nil, false // array end marker
		}
		index := fastAtoi(step)
		if index < 0 || index >= len(*v) {
			return nil, false // invalid or out of bounds
		}
		return (*v)[index], true

	case *any:
		// Interface pointer - recurse once
		if v == nil {
			return nil, false
		}
		return fastGet(*v, step)

	default:
		// Fast path failed, need reflection fallback
		return nil, false
	}
}

// getTokenAtIndex computes an internalToken for a specific path step without allocating a slice.
// Optimized for string-only Path - avoid allocating entire tokens slice, compute on-demand.
func getTokenAtIndex(path Path, index int) internalToken {
	if index >= len(path) {
		return internalToken{}
	}

	step := path[index] // step is already a string
	return internalToken{
		key:   step,
		index: fastAtoi(step),
	}
}

// tryArrayAccess attempts array access using type assertions for performance.
// Enhanced to handle all slice types efficiently.
func tryArrayAccess(current any, token internalToken) (any, bool, error) {
	// Fast type assertion path for common slice types
	switch arr := current.(type) {
	case []any:
		if token.key == "-" {
			return nil, true, nil // array end marker
		}
		if token.index < 0 || strconv.Itoa(token.index) != token.key {
			return nil, true, ErrInvalidIndex
		}
		if token.index < len(arr) {
			return arr[token.index], true, nil
		} else if token.index == len(arr) {
			// Allow pointing to one past array end (JSON Pointer spec)
			return nil, true, nil
		} else {
			return nil, true, ErrIndexOutOfBounds
		}

	case *[]any:
		if arr == nil {
			return nil, true, ErrNilPointer
		}
		if token.key == "-" {
			return nil, true, nil // array end marker
		}
		if token.index < 0 || strconv.Itoa(token.index) != token.key {
			return nil, true, ErrInvalidIndex
		}
		if token.index < len(*arr) {
			return (*arr)[token.index], true, nil
		} else if token.index == len(*arr) {
			// Allow pointing to one past array end (JSON Pointer spec)
			return nil, true, nil
		} else {
			return nil, true, ErrIndexOutOfBounds
		}

	case []string:
		if token.key == "-" {
			return nil, true, nil // array end marker
		}
		if token.index < 0 || strconv.Itoa(token.index) != token.key {
			return nil, true, ErrInvalidIndex
		}
		if token.index < len(arr) {
			return arr[token.index], true, nil
		} else if token.index == len(arr) {
			// Allow pointing to one past array end (JSON Pointer spec)
			return nil, true, nil
		} else {
			return nil, true, ErrIndexOutOfBounds
		}

	case []int:
		if token.key == "-" {
			return nil, true, nil // array end marker
		}
		if token.index < 0 || strconv.Itoa(token.index) != token.key {
			return nil, true, ErrInvalidIndex
		}
		if token.index < len(arr) {
			return arr[token.index], true, nil
		} else if token.index == len(arr) {
			// Allow pointing to one past array end (JSON Pointer spec)
			return nil, true, nil
		} else {
			return nil, true, ErrIndexOutOfBounds
		}

	case []float64:
		if token.key == "-" {
			return nil, true, nil // array end marker
		}
		if token.index < 0 || strconv.Itoa(token.index) != token.key {
			return nil, true, ErrInvalidIndex
		}
		if token.index < len(arr) {
			return arr[token.index], true, nil
		} else if token.index == len(arr) {
			// Allow pointing to one past array end (JSON Pointer spec)
			return nil, true, nil
		} else {
			return nil, true, ErrIndexOutOfBounds
		}

	default:
		// Fallback to reflection for other array types (like []User)
		if !isArray(current) {
			return nil, false, nil
		}

		if token.key == "-" {
			return nil, true, nil // array end marker
		}
		if token.index < 0 || strconv.Itoa(token.index) != token.key {
			return nil, true, ErrInvalidIndex
		}

		arrayVal := reflect.ValueOf(current)
		if token.index < arrayVal.Len() {
			return arrayVal.Index(token.index).Interface(), true, nil
		} else if token.index == arrayVal.Len() {
			// Allow pointing to one past array end (JSON Pointer spec)
			return nil, true, nil
		} else {
			return nil, true, ErrIndexOutOfBounds
		}
	}
}

// tryObjectAccess attempts object access using type assertions for performance.
// Enhanced with proper struct field handling.
func tryObjectAccess(current any, token internalToken) (any, bool, error) {
	// Fast type assertion path for common map types
	switch obj := current.(type) {
	case map[string]any:
		result, exists := obj[token.key]
		if !exists {
			return nil, true, nil // Key doesn't exist
		}
		return result, true, nil

	case *map[string]any:
		if obj == nil {
			return nil, true, ErrNilPointer
		}
		result, exists := (*obj)[token.key]
		if !exists {
			return nil, true, nil // Key doesn't exist
		}
		return result, true, nil

	case map[string]string:
		result, exists := obj[token.key]
		if !exists {
			return nil, true, nil // Key doesn't exist
		}
		return result, true, nil

	case map[string]int:
		result, exists := obj[token.key]
		if !exists {
			return nil, true, nil // Key doesn't exist
		}
		return result, true, nil

	case map[string]float64:
		result, exists := obj[token.key]
		if !exists {
			return nil, true, nil // Key doesn't exist
		}
		return result, true, nil

	default:
		// Fallback to reflection for other object types
		objVal := reflect.ValueOf(current)

		// Handle pointer dereferencing
		for objVal.Kind() == reflect.Ptr {
			if objVal.IsNil() {
				return nil, false, ErrNilPointer
			}
			objVal = objVal.Elem()
		}

		switch objVal.Kind() {
		case reflect.Map:
			mapKey := reflect.ValueOf(token.key)
			mapVal := objVal.MapIndex(mapKey)
			if !mapVal.IsValid() {
				return nil, true, nil // Key doesn't exist
			}
			return mapVal.Interface(), true, nil
		case reflect.Struct:
			// Handle struct fields using optimized struct field lookup
			if field := findStructField(objVal, token.key); field.IsValid() {
				return field.Interface(), true, nil
			}
			return nil, true, nil // Field not found
		case reflect.Invalid, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
			reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.Array,
			reflect.Chan, reflect.Func, reflect.Interface, reflect.Ptr, reflect.Slice, reflect.String, reflect.UnsafePointer:
			// Handle all other reflect.Kind types not supported for JSON Pointer traversal
			return nil, false, nil
		}
		// This should never be reached due to exhaustive case coverage
		return nil, false, nil
	}
}

// get retrieves value at JSON pointer path, returns error if path cannot be traversed.
// Optimized for zero-allocation string-only paths with layered fallback strategy.
func get(val any, path Path) (any, error) {
	pathLength := len(path)
	if pathLength == 0 {
		return val, nil
	}

	// Zero-allocation fast path for common cases
	current := val
	fastPathDepth := 0

	// Ultra-fast path - direct access without token creation
	for i := 0; i < pathLength; i++ {
		step := path[i] // step is already a string

		// Try direct fast path first (zero allocations for map[string]any)
		if result, ok := fastGet(current, step); ok {
			current = result
			fastPathDepth = i + 1
		} else {
			// Direct fast path failed, break to optimized type assertion fallback
			break
		}
	}

	// Optimized type assertion fallback for remaining path (if any)
	if fastPathDepth < pathLength {
		// Use optimized type assertions for the remaining tokens
		for i := fastPathDepth; i < pathLength; i++ {
			// Compute token on-demand only when needed
			token := getTokenAtIndex(path, i)

			if current == nil {
				return nil, ErrNotFound
			}

			// Try optimized array access first
			if result, handled, err := tryArrayAccess(current, token); err != nil {
				return nil, err
			} else if handled {
				current = result
				continue
			}

			// Try optimized object access
			if result, handled, err := tryObjectAccess(current, token); err != nil {
				return nil, err
			} else if handled {
				current = result
				continue
			}

			// Neither array nor object, can't traverse further
			return nil, ErrNotFound
		}
	}

	return current, nil
}

// findStructField finds a struct field by JSON tag or field name.
// Returns the field value if found, invalid reflect.Value otherwise.
func findStructField(structVal reflect.Value, key string) reflect.Value {
	structType := structVal.Type()
	numFields := structType.NumField()

	// First pass: look for exact JSON tag match
	for i := 0; i < numFields; i++ {
		field := structType.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		// Check JSON tag
		if jsonTag := field.Tag.Get("json"); jsonTag != "" {
			tagName := jsonTag
			// Find comma to extract just the field name part
			for j, r := range jsonTag {
				if r == ',' {
					tagName = jsonTag[:j]
					break
				}
			}
			if tagName == key {
				return structVal.Field(i)
			}
			if tagName == "-" {
				continue // Explicitly ignored field
			}
		}
	}

	// Second pass: look for field name match (if no JSON tag found)
	for i := 0; i < numFields; i++ {
		field := structType.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		// Skip if has JSON tag (already checked above)
		if field.Tag.Get("json") != "" {
			continue
		}

		// Match field name
		if field.Name == key {
			return structVal.Field(i)
		}
	}

	return reflect.Value{} // Not found
}

// Helper function to check if value is an array (slice)
func isArray(val any) bool {
	if val == nil {
		return false
	}
	return reflect.TypeOf(val).Kind() == reflect.Slice
}
