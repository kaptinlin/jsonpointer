package jsonpointer

import (
	"reflect"
)

func fastGet(val any, step string) (any, bool) {
	switch v := val.(type) {
	case map[string]any:
		result, exists := v[step]
		return result, exists

	case *map[string]any:
		if v == nil {
			return nil, false
		}
		result, exists := (*v)[step]
		return result, exists

	case []any:
		return fastSliceGet(v, step)

	case *[]any:
		if v == nil {
			return nil, false
		}
		return fastSliceGet(*v, step)

	case *any:
		if v == nil {
			return nil, false
		}
		return fastGet(*v, step)

	default:
		return nil, false
	}
}

func fastSliceGet(values []any, step string) (any, bool) {
	if step == "-" {
		return nil, false
	}
	index := fastAtoi(step)
	if index < 0 || index >= len(values) {
		return nil, false
	}
	return values[index], true
}

func tryArrayAccess(current any, key string) (any, bool, error) {
	switch arr := current.(type) {
	case []any:
		index, err := validateAndAccessArray(key, len(arr))
		if err != nil {
			return nil, true, err
		}
		return arr[index], true, nil

	case *[]any:
		if arr == nil {
			return nil, true, ErrNilPointer
		}
		index, err := validateAndAccessArray(key, len(*arr))
		if err != nil {
			return nil, true, err
		}
		return (*arr)[index], true, nil

	default:
		arrayVal, err := derefValue(reflect.ValueOf(current))
		if err != nil {
			return nil, true, err
		}

		if arrayVal.Kind() != reflect.Slice && arrayVal.Kind() != reflect.Array {
			return nil, false, nil
		}

		index, err := validateAndAccessArray(key, arrayVal.Len())
		if err != nil {
			return nil, true, err
		}
		return arrayVal.Index(index).Interface(), true, nil
	}
}

func tryObjectAccess(current any, key string) (any, bool, error) {
	switch obj := current.(type) {
	case map[string]any:
		result, exists := obj[key]
		if !exists {
			return nil, true, ErrKeyNotFound
		}
		return result, true, nil

	case *map[string]any:
		if obj == nil {
			return nil, true, ErrNilPointer
		}
		result, exists := (*obj)[key]
		if !exists {
			return nil, true, ErrKeyNotFound
		}
		return result, true, nil

	default:
		objVal, err := derefValue(reflect.ValueOf(current))
		if err != nil {
			return nil, false, err
		}

		switch objVal.Kind() {
		case reflect.Map:
			mapEntry, err := mapValueByPathKey(objVal, key)
			if err != nil {
				return nil, true, err
			}
			return mapEntry.Interface(), true, nil

		case reflect.Struct:
			if !structField(key, &objVal) {
				return nil, true, ErrFieldNotFound
			}
			return objVal.Interface(), true, nil

		default:
			return nil, false, nil
		}
	}
}

func get(val any, path Path) (any, error) {
	pathLength := len(path)
	if pathLength == 0 {
		return val, nil
	}

	current := val
	fastPathDepth := 0

	for i := range pathLength {
		step := path[i]

		if result, ok := fastGet(current, step); ok {
			current = result
			fastPathDepth = i + 1
		} else {
			break
		}
	}

	for i := fastPathDepth; i < pathLength; i++ {
		step := path[i]

		if current == nil {
			return nil, ErrNotFound
		}

		result, handled, err := tryArrayAccess(current, step)
		if err != nil {
			return nil, err
		}
		if handled {
			current = result
			continue
		}

		result, handled, err = tryObjectAccess(current, step)
		if err != nil {
			return nil, err
		}
		if handled {
			current = result
			continue
		}

		return nil, ErrNotFound
	}

	return current, nil
}
