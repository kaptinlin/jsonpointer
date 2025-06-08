package jsonpointer

import (
	"reflect"
	"strconv"
)

// find locates a target in document specified by JSON Pointer.
// Returns the object containing the target and key used to reference that object.
//
// Throws ErrNotFound if pointer does not result into a value in the middle
// of the path. If the last element of the path does not result into a value, the
// lookup succeeds with `val` set to nil. It can be used to discriminate
// missing values, because nil is not a valid JSON value.
//
// TypeScript original code:
//
//	export const find = (val: unknown, path: Path): Reference => {
//	  const pathLength = path.length;
//	  if (!pathLength) return {val};
//	  let obj: Reference['obj'];
//	  let key: Reference['key'];
//	  for (let i = 0; i < pathLength; i++) {
//	    obj = val;
//	    key = path[i];
//	    if (isArray(obj)) {
//	      const length = obj.length;
//	      if (key === '-') key = length;
//	      else {
//	        if (typeof key === 'string') {
//	          const key2 = ~~key;
//	          if ('' + key2 !== key) throw new Error('INVALID_INDEX');
//	          key = key2;
//	          if (key < 0) throw new Error('INVALID_INDEX');
//	        }
//	      }
//	      val = obj[key];
//	    } else if (typeof obj === 'object' && !!obj) {
//	      val = has(obj, key as string) ? (obj as any)[key] : undefined;
//	    } else throw new Error('NOT_FOUND');
//	  }
//	  const ref: Reference = {val, obj, key};
//	  return ref;
//	};
func find(val any, path Path) (*Reference, error) {
	pathLength := len(path)
	if pathLength == 0 {
		return &Reference{Val: val}, nil
	}

	var obj any
	var key any

	for i := 0; i < pathLength; i++ {
		obj = val
		key = path[i]

		switch {
		case isArrayValue(obj):
			// Handle array access - prioritize type assertions for performance
			var length int
			var getArrayValue func(int) any

			// Fast type assertion path
			switch arr := obj.(type) {
			case []any:
				length = len(arr)
				getArrayValue = func(index int) any { return arr[index] }
			case *[]any:
				if arr == nil {
					return nil, ErrNotFound
				}
				length = len(*arr)
				getArrayValue = func(index int) any { return (*arr)[index] }
			case []string:
				length = len(arr)
				getArrayValue = func(index int) any { return arr[index] }
			case []int:
				length = len(arr)
				getArrayValue = func(index int) any { return arr[index] }
			case []float64:
				length = len(arr)
				getArrayValue = func(index int) any { return arr[index] }
			default:
				// Fallback to reflection for other array types
				arrayVal := reflect.ValueOf(obj)
				length = arrayVal.Len()
				getArrayValue = func(index int) any { return arrayVal.Index(index).Interface() }
			}

			if keyStr, ok := key.(string); ok && keyStr == "-" {
				// Array end marker: key becomes array length
				key = length
				val = nil // undefined in TypeScript
			} else {
				// Convert key to integer
				var keyInt int
				switch k := key.(type) {
				case int:
					keyInt = k
				case string:
					// Parse string to int
					parsed, err := strconv.Atoi(k)
					if err != nil {
						return nil, ErrInvalidIndex
					}
					// Check if string representation matches parsed value (~~key behavior)
					if strconv.Itoa(parsed) != k {
						return nil, ErrInvalidIndex
					}
					keyInt = parsed
				default:
					return nil, ErrInvalidIndex
				}

				// Check for negative index
				if keyInt < 0 {
					return nil, ErrInvalidIndex
				}

				// Update key to the integer value
				key = keyInt

				// Get array value if index is valid
				if keyInt < length {
					val = getArrayValue(keyInt)
				} else {
					val = nil // undefined in TypeScript
				}
			}
		case isObjectValue(obj) && obj != nil:
			// Handle object/map access - prioritize type assertions for performance
			keyStr, ok := key.(string)
			if !ok {
				return nil, ErrNotFound
			}

			// Fast type assertion path
			switch objMap := obj.(type) {
			case map[string]any:
				if result, exists := objMap[keyStr]; exists {
					val = result
				} else {
					val = nil // undefined in TypeScript
				}
			case *map[string]any:
				if objMap == nil {
					return nil, ErrNotFound
				}
				if result, exists := (*objMap)[keyStr]; exists {
					val = result
				} else {
					val = nil // undefined in TypeScript
				}
			case map[string]string:
				if result, exists := objMap[keyStr]; exists {
					val = result
				} else {
					val = nil // undefined in TypeScript
				}
			case map[string]int:
				if result, exists := objMap[keyStr]; exists {
					val = result
				} else {
					val = nil // undefined in TypeScript
				}
			case map[string]float64:
				if result, exists := objMap[keyStr]; exists {
					val = result
				} else {
					val = nil // undefined in TypeScript
				}
			default:
				// Fallback to reflection for other object types
				objVal := reflect.ValueOf(obj)
				if objVal.Kind() == reflect.Map {
					// Handle map
					mapKey := reflect.ValueOf(keyStr)
					mapVal := objVal.MapIndex(mapKey)
					if mapVal.IsValid() {
						val = mapVal.Interface()
					} else {
						val = nil // undefined in TypeScript
					}
				} else {
					// Handle struct using our optimized struct field lookup
					if structField(keyStr, &objVal) {
						val = objVal.Interface()
					} else {
						val = nil // Field not found
					}
				}
			}
		default:
			// Not an array or object, can't traverse further
			return nil, ErrNotFound
		}
	}

	ref := &Reference{
		Val: val,
		Obj: obj,
		Key: key,
	}
	return ref, nil
}

// isArrayValue checks if value is an array (slice).
func isArrayValue(val any) bool {
	if val == nil {
		return false
	}
	return reflect.TypeOf(val).Kind() == reflect.Slice
}

// isObjectValue checks if value is an object (map or struct).
func isObjectValue(val any) bool {
	if val == nil {
		return false
	}
	kind := reflect.TypeOf(val).Kind()
	return kind == reflect.Map || kind == reflect.Struct || kind == reflect.Ptr
}
