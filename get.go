package jsonpointer

import "reflect"

func resolveValue(doc any, pointer Pointer) (any, error) {
	current := doc
	for depth, token := range pointer.tokens {
		next, err := step(current, token)
		if err != nil {
			return nil, newError(err, pointer, depth)
		}
		current = next
	}
	return current, nil
}

func resolveReference(doc any, pointer Pointer) (Reference, error) {
	if pointer.IsRoot() {
		return Reference{value: doc, pointer: pointer}, nil
	}

	current := doc
	var parent any
	var token string
	for depth, stepToken := range pointer.tokens {
		parent = current
		token = stepToken

		next, err := step(current, stepToken)
		if err != nil {
			return Reference{}, newError(err, pointer, depth)
		}
		current = next
	}

	return Reference{
		value:   current,
		parent:  parent,
		token:   token,
		pointer: pointer,
	}, nil
}

func step(current any, token string) (any, error) {
	if current == nil {
		return nil, ErrNotFound
	}

	switch value := current.(type) {
	case map[string]any:
		return stringMapValue(value, token)
	case *map[string]any:
		if value == nil {
			return nil, ErrNilPointer
		}
		return stringMapValue(*value, token)
	case []any:
		return sliceValue(value, token)
	case *[]any:
		if value == nil {
			return nil, ErrNilPointer
		}
		return sliceValue(*value, token)
	}

	value, err := derefValue(reflect.ValueOf(current))
	if err != nil {
		return nil, err
	}

	switch value.Kind() {
	case reflect.Slice, reflect.Array:
		index, err := validateAndAccessArray(token, value.Len())
		if err != nil {
			return nil, err
		}
		return value.Index(index).Interface(), nil
	case reflect.Map:
		result, err := mapValueByToken(value, token)
		if err != nil {
			return nil, err
		}
		return result.Interface(), nil
	default:
		return nil, ErrNotFound
	}
}

func sliceValue(values []any, token string) (any, error) {
	index, err := validateAndAccessArray(token, len(values))
	if err != nil {
		return nil, err
	}
	return values[index], nil
}

func stringMapValue(values map[string]any, token string) (any, error) {
	result, ok := values[token]
	if !ok {
		return nil, ErrKeyNotFound
	}
	return result, nil
}
