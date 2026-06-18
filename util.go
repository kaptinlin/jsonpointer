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

func fastAtoi(s string) int {
	if len(s) == 0 {
		return -1
	}
	if s == "0" {
		return 0
	}
	if s[0] == '0' {
		return -1
	}

	var n int
	for i := range len(s) {
		c := s[i]
		if c < '0' || c > '9' {
			return -1
		}
		next := n*10 + int(c-'0')
		if next < n {
			return -1
		}
		n = next
	}
	return n
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

// IsArrayIndex reports whether token is a readable JSON Pointer array index.
func IsArrayIndex(token string) bool {
	return fastAtoi(token) >= 0
}

func validateAndAccessArray(token string, length int) (int, error) {
	if token == "-" {
		return -1, ErrIndexOutOfBounds
	}
	index := fastAtoi(token)
	if index < 0 {
		return -1, ErrInvalidIndex
	}
	if index >= length {
		return -1, ErrIndexOutOfBounds
	}
	return index, nil
}
