package jsonpointer

import (
	"reflect"
	"slices"
	"strings"
	"sync"
	"unicode"
)

type structFields map[string][]int

type structFieldInfo struct {
	name   string
	index  []int
	tagged bool
}

type structScan struct {
	typ   reflect.Type
	index []int
}

var structFieldsCache sync.Map

func structField(value reflect.Value, token string) (reflect.Value, error) {
	fields := getStructFields(value.Type())
	index, ok := fields[token]
	if !ok {
		return reflect.Value{}, ErrFieldNotFound
	}

	field, err := fieldByIndex(value, index)
	if err != nil {
		return reflect.Value{}, err
	}
	if !field.CanInterface() {
		return reflect.Value{}, ErrFieldNotFound
	}
	return field, nil
}

func getStructFields(t reflect.Type) structFields {
	if cached, ok := structFieldsCache.Load(t); ok {
		return cached.(structFields)
	}

	fields := buildStructFields(t)
	cached, _ := structFieldsCache.LoadOrStore(t, fields)
	return cached.(structFields)
}

func buildStructFields(t reflect.Type) structFields {
	candidates := collectStructFields(t)
	slices.SortFunc(candidates, compareStructFields)

	fields := make(structFields)
	for i := 0; i < len(candidates); {
		j := i + 1
		for j < len(candidates) && candidates[j].name == candidates[i].name {
			j++
		}
		if field, ok := dominantField(candidates[i:j]); ok {
			fields[field.name] = slices.Clone(field.index)
		}
		i = j
	}
	return fields
}

func collectStructFields(t reflect.Type) []structFieldInfo {
	var fields []structFieldInfo
	current := []structScan{{typ: t}}
	visited := map[reflect.Type]bool{}

	for len(current) > 0 {
		var next []structScan
		for _, scan := range current {
			if visited[scan.typ] {
				continue
			}
			visited[scan.typ] = true

			for i := range scan.typ.NumField() {
				field := scan.typ.Field(i)
				fieldType := field.Type
				if fieldType.Kind() == reflect.Pointer {
					fieldType = fieldType.Elem()
				}

				if !field.IsExported() && !isUnexportedEmbeddedStruct(&field, fieldType) {
					continue
				}

				tag := field.Tag.Get("json")
				if tag == "-" {
					continue
				}

				name, _, _ := strings.Cut(tag, ",")
				if !validTag(name) {
					name = ""
				}

				index := slices.Concat(scan.index, []int{i})
				tagged := name != ""
				if tagged || !field.Anonymous || fieldType.Kind() != reflect.Struct {
					if name == "" {
						name = field.Name
					}
					fields = append(fields, structFieldInfo{
						name:   name,
						index:  index,
						tagged: tagged,
					})
				}

				if !tagged && field.Anonymous && fieldType.Kind() == reflect.Struct {
					next = append(next, structScan{typ: fieldType, index: index})
				}
			}
		}
		current = next
	}

	return fields
}

func isUnexportedEmbeddedStruct(field *reflect.StructField, fieldType reflect.Type) bool {
	return field.Anonymous && !field.IsExported() && fieldType.Kind() == reflect.Struct
}

func compareStructFields(a, b structFieldInfo) int {
	if a.name != b.name {
		return strings.Compare(a.name, b.name)
	}
	if len(a.index) != len(b.index) {
		return len(a.index) - len(b.index)
	}
	if a.tagged != b.tagged {
		if a.tagged {
			return -1
		}
		return 1
	}
	return slices.Compare(a.index, b.index)
}

func dominantField(fields []structFieldInfo) (structFieldInfo, bool) {
	if len(fields) > 1 &&
		len(fields[0].index) == len(fields[1].index) &&
		fields[0].tagged == fields[1].tagged {
		return structFieldInfo{}, false
	}
	return fields[0], true
}

func fieldByIndex(value reflect.Value, index []int) (reflect.Value, error) {
	for depth, fieldIndex := range index {
		for value.Kind() == reflect.Pointer {
			if value.IsNil() {
				return reflect.Value{}, ErrNilPointer
			}
			value = value.Elem()
		}

		value = value.Field(fieldIndex)
		if depth < len(index)-1 && value.Kind() == reflect.Pointer && value.IsNil() {
			return reflect.Value{}, ErrNilPointer
		}
	}
	return value, nil
}

func validTag(tag string) bool {
	if tag == "" {
		return true
	}
	for _, r := range tag {
		switch {
		case strings.ContainsRune("!#$%&()*+-./:<=>?@[]^_{|}~ ", r):
		case unicode.IsLetter(r), unicode.IsDigit(r):
		default:
			return false
		}
	}
	return true
}
