package jsonpointer

import (
	"cmp"
	"reflect"
	"strings"
	"sync"
)

type structFields map[string]int

var structFieldsCache sync.Map

// structField looks up the specified field in a struct and updates value to point to that field if found.
// Returns true if the field exists and is accessible, false otherwise.
func structField(field string, value *reflect.Value) bool {
	for value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return false
		}
		*value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return false
	}

	fields := getStructFields(value.Type())
	fieldIndex, ok := fields[field]
	if !ok {
		return false
	}

	*value = value.Field(fieldIndex)
	return true
}

func getStructFields(t reflect.Type) structFields {
	if cached, ok := structFieldsCache.Load(t); ok {
		return cached.(structFields)
	}

	fields := make(structFields)
	for structFieldInfo := range t.Fields() {
		if !structFieldInfo.IsExported() {
			continue
		}

		name := getFieldName(&structFieldInfo)
		if name == "-" {
			continue
		}

		fields[name] = structFieldInfo.Index[0]
	}

	structFieldsCache.Store(t, fields)
	return fields
}

func getFieldName(field *reflect.StructField) string {
	tag := field.Tag.Get("json")
	if tag == "" {
		return field.Name
	}

	name, _, _ := strings.Cut(tag, ",")
	return cmp.Or(name, field.Name)
}
