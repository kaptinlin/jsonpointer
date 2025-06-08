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

		if val == nil {
			return nil, ErrNotFound
		}

		switch obj := obj.(type) {
		// Fast path for []any (most common case)
		case []any:
			keyInt, isEndMarker, err := parseArrayKey(key)
			if err != nil {
				return nil, err
			}
			if isEndMarker {
				key = len(obj)
				val = nil
			} else {
				key = keyInt
				if keyInt < len(obj) {
					val = obj[keyInt]
				} else {
					val = nil
				}
			}

		// Fast path for map[string]any (most common case)
		case map[string]any:
			keyStr, ok := key.(string)
			if !ok {
				return nil, ErrNotFound
			}
			if result, exists := obj[keyStr]; exists {
				val = result
			} else {
				val = nil
			}

		// Handle other specific types without conversion
		case []string:
			keyInt, isEndMarker, err := parseArrayKey(key)
			if err != nil {
				return nil, err
			}
			if isEndMarker {
				key = len(obj)
				val = nil
			} else {
				key = keyInt
				if keyInt < len(obj) {
					val = obj[keyInt]
				} else {
					val = nil
				}
			}

		case []int:
			keyInt, isEndMarker, err := parseArrayKey(key)
			if err != nil {
				return nil, err
			}
			if isEndMarker {
				key = len(obj)
				val = nil
			} else {
				key = keyInt
				if keyInt < len(obj) {
					val = obj[keyInt]
				} else {
					val = nil
				}
			}

		case []float64:
			keyInt, isEndMarker, err := parseArrayKey(key)
			if err != nil {
				return nil, err
			}
			if isEndMarker {
				key = len(obj)
				val = nil
			} else {
				key = keyInt
				if keyInt < len(obj) {
					val = obj[keyInt]
				} else {
					val = nil
				}
			}

		case map[string]string:
			keyStr, ok := key.(string)
			if !ok {
				return nil, ErrNotFound
			}
			if result, exists := obj[keyStr]; exists {
				val = result
			} else {
				val = nil
			}

		case map[string]int:
			keyStr, ok := key.(string)
			if !ok {
				return nil, ErrNotFound
			}
			if result, exists := obj[keyStr]; exists {
				val = result
			} else {
				val = nil
			}

		case map[string]float64:
			keyStr, ok := key.(string)
			if !ok {
				return nil, ErrNotFound
			}
			if result, exists := obj[keyStr]; exists {
				val = result
			} else {
				val = nil
			}

		default:
			// Handle reflection-based access for other types
			objVal := reflect.ValueOf(obj)
			switch objVal.Kind() {
			case reflect.Slice:
				keyInt, isEndMarker, err := parseArrayKey(key)
				if err != nil {
					return nil, err
				}
				if isEndMarker {
					key = objVal.Len()
					val = nil
				} else {
					key = keyInt
					if keyInt < objVal.Len() {
						val = objVal.Index(keyInt).Interface()
					} else {
						val = nil
					}
				}

			case reflect.Map:
				keyStr, ok := key.(string)
				if !ok {
					return nil, ErrNotFound
				}
				mapKey := reflect.ValueOf(keyStr)
				mapVal := objVal.MapIndex(mapKey)
				if mapVal.IsValid() {
					val = mapVal.Interface()
				} else {
					val = nil
				}

			case reflect.Struct:
				keyStr, ok := key.(string)
				if !ok {
					return nil, ErrNotFound
				}
				if structField(keyStr, &objVal) {
					val = objVal.Interface()
				} else {
					val = nil
				}

			case reflect.Ptr:
				if objVal.IsNil() {
					return nil, ErrNotFound
				}
				// Dereference pointer and retry
				val = objVal.Elem().Interface()
				i-- // Retry with dereferenced value
				continue

			case reflect.Array:
				// Handle arrays similar to slices
				keyInt, isEndMarker, err := parseArrayKey(key)
				if err != nil {
					return nil, err
				}
				if isEndMarker {
					key = objVal.Len()
					val = nil
				} else {
					key = keyInt
					if keyInt < objVal.Len() {
						val = objVal.Index(keyInt).Interface()
					} else {
						val = nil
					}
				}

			case reflect.Interface:
				// Dereference interface and retry
				if objVal.IsNil() {
					return nil, ErrNotFound
				}
				val = objVal.Elem().Interface()
				i-- // Retry with dereferenced value
				continue

			// All primitive types and non-traversable types
			case reflect.Invalid, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16,
				reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16,
				reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32,
				reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.Chan,
				reflect.Func, reflect.String, reflect.UnsafePointer:
				return nil, ErrNotFound
			}
		}
	}

	return &Reference{
		Val: val,
		Obj: obj,
		Key: key,
	}, nil
}

// parseArrayKey converts a key to array index with clear return values
func parseArrayKey(key any) (index int, isEndMarker bool, err error) {
	switch k := key.(type) {
	case int:
		if k < 0 {
			return 0, false, ErrInvalidIndex
		}
		return k, false, nil
	case string:
		if k == "-" {
			return 0, true, nil
		}
		parsed, parseErr := strconv.Atoi(k)
		if parseErr != nil || parsed < 0 {
			return 0, false, ErrInvalidIndex
		}
		// Check if string representation matches parsed value
		if strconv.Itoa(parsed) != k {
			return 0, false, ErrInvalidIndex
		}
		return parsed, false, nil
	default:
		return 0, false, ErrInvalidIndex
	}
}
