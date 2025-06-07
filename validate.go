package jsonpointer

import "reflect"

// validateJsonPointer validates a JSON Pointer string or Path.
//
// TypeScript original code from validate.ts:
//
//	export const validateJsonPointer = (pointer: string | Path | unknown) => {
//	  if (typeof pointer === 'string') {
//	    if (pointer) {
//	      if (pointer[0] !== '/') throw new Error('POINTER_INVALID');
//	      if (pointer.length > 1024) throw new Error('POINTER_TOO_LONG');
//	    }
//	  } else validatePath(pointer);
//	};
func validateJsonPointer(pointer any) error {
	if str, ok := pointer.(string); ok {
		// Handle string pointer
		if str != "" {
			if len(str) == 0 || str[0] != '/' {
				return ErrPointerInvalid
			}
			if len(str) > 1024 {
				return ErrPointerTooLong
			}
		}
	} else {
		// Validate as path
		return validatePath(pointer)
	}
	return nil
}

// validatePath validates a path array.
//
// TypeScript original code from validate.ts:
//
//	export const validatePath = (path: Path | unknown) => {
//	  if (!isArray(path)) throw new Error('Invalid path.');
//	  if (path.length > 256) throw new Error('Path too long.');
//	  for (const step of path) {
//	    switch (typeof step) {
//	      case 'string':
//	      case 'number':
//	        continue;
//	      default:
//	        throw new Error('Invalid path step.');
//	    }
//	  }
//	};
func validatePath(path any) error {
	// Check if path is an array (slice)
	val := reflect.ValueOf(path)
	if val.Kind() != reflect.Slice {
		return ErrInvalidPath
	}

	// Check length
	length := val.Len()
	if length > 256 {
		return ErrPathTooLong
	}

	// Validate each step
	for i := 0; i < length; i++ {
		step := val.Index(i).Interface()

		// Check if step is string or number
		switch step.(type) {
		case string:
			// Valid
		case int, int8, int16, int32, int64:
			// Valid integers
		case uint, uint8, uint16, uint32, uint64:
			// Valid unsigned integers
		case float32, float64:
			// Valid floats (numbers in TypeScript)
		default:
			return ErrInvalidPathStep
		}
	}

	return nil
}
