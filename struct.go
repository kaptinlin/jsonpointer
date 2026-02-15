package jsonpointer

import (
	"reflect"
	"strings"
	"sync"
)

// structFields caches field mapping for struct types.
type structFields map[string]int

// structFieldsCache is a global cache that stores field mapping for each struct type.
var structFieldsCache sync.Map

// structField looks up the specified field in a struct and updates value to point to that field if found.
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

// getStructFields gets field mapping for struct type with caching.
func getStructFields(t reflect.Type) structFields {
	if cached, ok := structFieldsCache.Load(t); ok {
		return cached.(structFields)
	}

	fields := make(structFields)
	numField := t.NumField()

	for i := range numField {
		field := t.Field(i)

		if !field.IsExported() {
			continue
		}

		name := getFieldName(field)
		if name == "-" {
			continue
		}

		fields[name] = i
	}

	structFieldsCache.Store(t, fields)
	return fields
}

// getFieldName gets the JSON name of field, supports basic JSON tags.
func getFieldName(field reflect.StructField) string {
	tag := field.Tag.Get("json")
	if tag == "" {
		return field.Name
	}

	name, _, _ := strings.Cut(tag, ",")
	if name != "" {
		return name
	}

	return field.Name
}
