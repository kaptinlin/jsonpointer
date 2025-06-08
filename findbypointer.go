package jsonpointer

import (
	"reflect"
	"strconv"
	"strings"
)

// findByPointer optimized string-based find operation.
// Direct string parsing without path array allocation for better performance.
//
// TypeScript original code from findByPointer/v5.ts:
//
//	export const findByPointer = (pointer: string, val: unknown): Reference => {
//	  if (!pointer) return {val};
//	  let obj: Reference['obj'];
//	  let key: Reference['key'];
//	  let indexOfSlash = 0;
//	  let indexAfterSlash = 1;
//	  while (indexOfSlash > -1) {
//	    indexOfSlash = pointer.indexOf('/', indexAfterSlash);
//	    key = indexOfSlash > -1 ? pointer.substring(indexAfterSlash, indexOfSlash) : pointer.substring(indexAfterSlash);
//	    indexAfterSlash = indexOfSlash + 1;
//	    obj = val;
//	    if (isArray(obj)) {
//	      const length = obj.length;
//	      if (key === '-') key = length;
//	      else {
//	        const key2 = ~~key;
//	        if ('' + key2 !== key) throw new Error('INVALID_INDEX');
//	        key = key2;
//	        if (key < 0) throw 'INVALID_INDEX';
//	      }
//	      val = obj[key];
//	    } else if (typeof obj === 'object' && !!obj) {
//	      key = unescapeComponent(key);
//	      val = has(obj, key) ? (obj as any)[key] : undefined;
//	    } else throw 'NOT_FOUND';
//	  }
//	  return {val, obj, key};
//	};
func findByPointer(pointer string, val any) (*Reference, error) {
	if pointer == "" {
		return &Reference{Val: val}, nil
	}

	var obj any
	var key any
	indexOfSlash := 0
	indexAfterSlash := 1

	for indexOfSlash > -1 {
		// Find next slash or end of string
		indexOfSlash = strings.Index(pointer[indexAfterSlash:], "/")
		if indexOfSlash > -1 {
			indexOfSlash += indexAfterSlash // Adjust for substring offset
		}

		// Extract key substring
		var keyStr string
		if indexOfSlash > -1 {
			keyStr = pointer[indexAfterSlash:indexOfSlash]
		} else {
			keyStr = pointer[indexAfterSlash:]
		}

		indexAfterSlash = indexOfSlash + 1
		obj = val

		switch {
		case isArrayPointer(obj):
			// Handle array access
			arrayVal := reflect.ValueOf(obj)
			length := arrayVal.Len()

			if keyStr == "-" {
				// Array end marker: key becomes array length
				key = length
				val = nil // undefined in TypeScript
			} else {
				// Convert key to integer (~~key behavior in TypeScript)
				keyInt, err := strconv.Atoi(keyStr)
				if err != nil {
					return nil, ErrInvalidIndex
				}
				// Check if string representation matches parsed value
				if strconv.Itoa(keyInt) != keyStr {
					return nil, ErrInvalidIndex
				}
				if keyInt < 0 {
					return nil, ErrInvalidIndex
				}

				key = keyInt

				// Get array value if index is valid
				if keyInt < length {
					val = arrayVal.Index(keyInt).Interface()
				} else {
					val = nil // undefined in TypeScript
				}
			}
		case isObjectPointer(obj) && obj != nil:
			// Handle object/map access
			// Unescape the key component
			keyStr = UnescapeComponent(keyStr)
			key = keyStr

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
				// 使用优化的 struct 字段查找处理结构体
				if structField(keyStr, &objVal) {
					val = objVal.Interface()
				} else {
					val = nil // 字段未找到
				}
			}
		default:
			// Not an array or object, can't traverse further
			return nil, ErrNotFound
		}
	}

	return &Reference{
		Val: val,
		Obj: obj,
		Key: key,
	}, nil
}

// Helper function to check if value is an array (slice) for pointer operations
func isArrayPointer(val any) bool {
	if val == nil {
		return false
	}
	return reflect.TypeOf(val).Kind() == reflect.Slice
}

// Helper function to check if value is an object (map or struct) for pointer operations
func isObjectPointer(val any) bool {
	if val == nil {
		return false
	}
	kind := reflect.TypeOf(val).Kind()
	return kind == reflect.Map || kind == reflect.Struct || kind == reflect.Ptr
}
