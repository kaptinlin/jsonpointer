package jsonpointer

import (
	"reflect"
	"strings"
)

func parsePointer(pointer string) ([]string, error) {
	if err := validatePointerString(pointer); err != nil {
		return nil, err
	}
	if pointer == "" {
		return nil, nil
	}

	tokens := make([]string, 0, strings.Count(pointer, "/"))
	for segment := range strings.SplitSeq(pointer[1:], "/") {
		token, err := unescapeToken(segment)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}

func formatPointer(tokens []string) string {
	if len(tokens) == 0 {
		return ""
	}

	var b strings.Builder
	b.Grow(pointerLength(tokens))
	for _, token := range tokens {
		b.WriteByte('/')
		b.WriteString(escapeToken(token))
	}
	return b.String()
}

func pointerLength(tokens []string) int {
	if len(tokens) == 0 {
		return 0
	}

	length := len(tokens)
	for _, token := range tokens {
		for i := range len(token) {
			switch token[i] {
			case '~', '/':
				length += 2
			default:
				length++
			}
		}
	}
	return length
}

func escapeToken(token string) string {
	if !strings.Contains(token, "/") && !strings.Contains(token, "~") {
		return token
	}

	var b strings.Builder
	b.Grow(len(token) * 2)
	for i := range len(token) {
		switch token[i] {
		case '~':
			b.WriteString("~0")
		case '/':
			b.WriteString("~1")
		default:
			b.WriteByte(token[i])
		}
	}
	return b.String()
}

func unescapeToken(encoded string) (string, error) {
	if !strings.Contains(encoded, "~") {
		return encoded, nil
	}

	var b strings.Builder
	b.Grow(len(encoded))
	for i := 0; i < len(encoded); i++ {
		if encoded[i] != '~' {
			b.WriteByte(encoded[i])
			continue
		}
		if i+1 >= len(encoded) {
			return "", ErrInvalidPointer
		}
		switch encoded[i+1] {
		case '0':
			b.WriteByte('~')
		case '1':
			b.WriteByte('/')
		default:
			return "", ErrInvalidPointer
		}
		i++
	}
	return b.String(), nil
}

type arrayIndexState uint8

const (
	arrayIndexInvalid arrayIndexState = iota
	arrayIndexValue
	arrayIndexOverflow
	arrayIndexDash
)

func parseArrayIndex(token string) (int, arrayIndexState) {
	if token == "-" {
		return 0, arrayIndexDash
	}
	if len(token) == 0 {
		return 0, arrayIndexInvalid
	}
	if token == "0" {
		return 0, arrayIndexValue
	}
	if token[0] == '0' {
		return 0, arrayIndexInvalid
	}

	const maxInt = int(^uint(0) >> 1)

	var n int
	for i := range len(token) {
		c := token[i]
		if c < '0' || c > '9' {
			return 0, arrayIndexInvalid
		}
		digit := int(c - '0')
		if n > (maxInt-digit)/10 {
			return 0, arrayIndexOverflow
		}
		n = n*10 + digit
	}
	return n, arrayIndexValue
}

func derefValue(v reflect.Value) (reflect.Value, error) {
	for {
		if !v.IsValid() {
			return reflect.Value{}, ErrNotFound
		}
		switch v.Kind() {
		case reflect.Pointer:
			if v.IsNil() {
				return reflect.Value{}, ErrNilPointer
			}
			v = v.Elem()
		case reflect.Interface:
			if v.IsNil() {
				return reflect.Value{}, ErrNotFound
			}
			v = v.Elem()
		default:
			return v, nil
		}
	}
}

func mapValueByToken(mapVal reflect.Value, token string) (reflect.Value, error) {
	mapKey := reflect.ValueOf(token)
	mapKeyType := mapVal.Type().Key()
	stringType := reflect.TypeFor[string]()

	switch {
	case stringType.AssignableTo(mapKeyType):
	case stringType.ConvertibleTo(mapKeyType):
		mapKey = mapKey.Convert(mapKeyType)
	default:
		return reflect.Value{}, ErrNotFound
	}

	mapEntry := mapVal.MapIndex(mapKey)
	if !mapEntry.IsValid() {
		return reflect.Value{}, ErrKeyNotFound
	}
	return mapEntry, nil
}

// IsArrayIndex reports whether token is a canonical, representable array index.
func IsArrayIndex(token string) bool {
	_, state := parseArrayIndex(token)
	return state == arrayIndexValue
}

func validateAndAccessArray(token string, length int) (int, error) {
	index, state := parseArrayIndex(token)
	switch state {
	case arrayIndexInvalid:
		return -1, ErrInvalidIndex
	case arrayIndexOverflow, arrayIndexDash:
		return -1, ErrIndexOutOfBounds
	case arrayIndexValue:
	default:
		return -1, ErrInvalidIndex
	}
	if index >= length {
		return -1, ErrIndexOutOfBounds
	}
	return index, nil
}
