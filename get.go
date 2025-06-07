package jsonpointer

import (
	"reflect"
)

// fastGet implements ultra-fast path that avoids token allocation entirely.
// direct access without intermediate token creation.
func fastGet(val any, step any) (any, bool) {
	switch v := val.(type) {
	case map[string]any:
		// Most common case: map[string]any - direct string key access
		if keyStr, ok := step.(string); ok {
			result, exists := v[keyStr]
			return result, exists
		}
		return nil, false

	case *map[string]any:
		// Pointer to map optimization
		if v == nil {
			return nil, false
		}
		if keyStr, ok := step.(string); ok {
			result, exists := (*v)[keyStr]
			return result, exists
		}
		return nil, false

	case []any:
		// Array access - need to compute index
		keyStr := componentToString(step)
		if keyStr == "-" {
			return nil, false // array end marker
		}
		index := fastAtoi(keyStr)
		if index < 0 || index >= len(v) {
			return nil, false // invalid or out of bounds
		}
		return v[index], true

	case *[]any:
		// Pointer to slice optimization
		if v == nil {
			return nil, false
		}
		keyStr := componentToString(step)
		if keyStr == "-" {
			return nil, false // array end marker
		}
		index := fastAtoi(keyStr)
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

// get retrieves value at JSON pointer path, returns nil if not found.
// uses precomputed tokens for faster array index access.
// TypeScript original code:
//
//	export const get = (val: unknown, path: Path): unknown | undefined => {
//	  const pathLength = path.length;
//	  let key: string | number;
//	  if (!pathLength) return val;
//	  for (let i = 0; i < pathLength; i++) {
//	    key = path[i];
//	    if (val instanceof Array) {
//	      if (key === '-') return undefined;
//	      const key2 = ~~key;
//	      if ('' + key2 !== key) return undefined;
//	      key = key2;
//	      if (key < 0) return undefined;
//	      val = val[key];
//	    } else if (typeof val === 'object') {
//	      if (!val || !has(val as object, key as string)) return undefined;
//	      val = (val as any)[key];
//	    } else return undefined;
//	  }
//	  return val;
//	};
//
// getTokenAtIndex computes an internalToken for a specific path step without allocating a slice.
// avoid allocating entire tokens slice, compute on-demand.
func getTokenAtIndex(path Path, index int) internalToken {
	if index >= len(path) {
		return internalToken{}
	}

	step := path[index]
	keyStr := componentToString(step)
	return internalToken{
		key:   keyStr,
		index: fastAtoi(keyStr),
	}
}

// tryArrayAccess attempts array access using type assertions for performance.
// It prioritizes fast type assertions over reflection.
func tryArrayAccess(current any, token internalToken) (any, bool) {
	// Fast type assertion path
	switch arr := current.(type) {
	case []any:
		if token.key == "-" {
			return nil, true // array end marker
		}
		if token.index < 0 || token.index >= len(arr) {
			return nil, true // invalid or out of bounds
		}
		return arr[token.index], true

	case *[]any:
		if arr == nil {
			return nil, true
		}
		if token.key == "-" {
			return nil, true // array end marker
		}
		if token.index < 0 || token.index >= len(*arr) {
			return nil, true // invalid or out of bounds
		}
		return (*arr)[token.index], true

	case []string:
		if token.key == "-" {
			return nil, true // array end marker
		}
		if token.index < 0 || token.index >= len(arr) {
			return nil, true // invalid or out of bounds
		}
		return arr[token.index], true

	case []int:
		if token.key == "-" {
			return nil, true // array end marker
		}
		if token.index < 0 || token.index >= len(arr) {
			return nil, true // invalid or out of bounds
		}
		return arr[token.index], true

	case []float64:
		if token.key == "-" {
			return nil, true // array end marker
		}
		if token.index < 0 || token.index >= len(arr) {
			return nil, true // invalid or out of bounds
		}
		return arr[token.index], true

	default:
		// Fallback to reflection for other array types
		if !isArray(current) {
			return nil, false
		}

		if token.key == "-" {
			return nil, true // array end marker
		}
		if token.index < 0 {
			return nil, true // invalid array index
		}

		arrayVal := reflect.ValueOf(current)
		if token.index >= arrayVal.Len() {
			return nil, true
		}
		return arrayVal.Index(token.index).Interface(), true
	}
}

// tryObjectAccess attempts object access using type assertions for performance.
// It prioritizes fast type assertions over reflection.
func tryObjectAccess(current any, token internalToken) (any, bool) {
	// Fast type assertion path
	switch obj := current.(type) {
	case map[string]any:
		result, exists := obj[token.key]
		if !exists {
			return nil, true // Key doesn't exist
		}
		return result, true

	case *map[string]any:
		if obj == nil {
			return nil, true
		}
		result, exists := (*obj)[token.key]
		if !exists {
			return nil, true // Key doesn't exist
		}
		return result, true

	case map[string]string:
		result, exists := obj[token.key]
		if !exists {
			return nil, true // Key doesn't exist
		}
		return result, true

	case map[string]int:
		result, exists := obj[token.key]
		if !exists {
			return nil, true // Key doesn't exist
		}
		return result, true

	case map[string]float64:
		result, exists := obj[token.key]
		if !exists {
			return nil, true // Key doesn't exist
		}
		return result, true

	default:
		// Fallback to reflection for other object types
		if !isObject(current) {
			return nil, false
		}

		objVal := reflect.ValueOf(current)
		if objVal.Kind() == reflect.Map {
			mapKey := reflect.ValueOf(token.key)
			mapVal := objVal.MapIndex(mapKey)
			if !mapVal.IsValid() {
				return nil, true // Key doesn't exist
			}
			return mapVal.Interface(), true
		} else {
			// Handle struct fields
			if objVal.Kind() == reflect.Ptr {
				objVal = objVal.Elem()
			}
			if objVal.Kind() != reflect.Struct {
				return nil, false
			}

			fieldVal := objVal.FieldByName(token.key)
			if !fieldVal.IsValid() || !fieldVal.CanInterface() {
				return nil, true
			}
			return fieldVal.Interface(), true
		}
	}
}

func get(val any, path Path) any {
	pathLength := len(path)
	if pathLength == 0 {
		return val
	}

	// zero-allocation fast path for common cases
	current := val
	fastPathDepth := 0

	// Ultra-fast path - direct access without token creation
	for i := 0; i < pathLength; i++ {
		step := path[i]

		// Try direct fast path first (zero allocations for map[string]any)
		if result, ok := fastGet(current, step); ok {
			current = result
			fastPathDepth = i + 1
		} else {
			// Direct fast path failed, break to reflection fallback
			break
		}
	}

	// Optimized type assertion fallback for remaining path (if any)
	if fastPathDepth < pathLength {
		// Use optimized type assertions for the remaining tokens
		for i := fastPathDepth; i < pathLength; i++ {
			// compute token on-demand only when needed
			token := getTokenAtIndex(path, i)

			if current == nil {
				return nil
			}

			// Try optimized array access first
			if result, handled := tryArrayAccess(current, token); handled {
				current = result
				continue
			}

			// Try optimized object access
			if result, handled := tryObjectAccess(current, token); handled {
				current = result
				continue
			}

			// Neither array nor object, can't traverse further
			return nil
		}
	}

	return current
}

// Helper function to check if value is an array (slice)
func isArray(val any) bool {
	if val == nil {
		return false
	}
	return reflect.TypeOf(val).Kind() == reflect.Slice
}

// Helper function to check if value is an object (map or struct)
func isObject(val any) bool {
	if val == nil {
		return false
	}
	kind := reflect.TypeOf(val).Kind()
	return kind == reflect.Map || kind == reflect.Struct || kind == reflect.Ptr
}
