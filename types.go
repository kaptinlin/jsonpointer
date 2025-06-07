// Package jsonpointer provides JSON Pointer (RFC 6901) implementation for Go.
// This is a direct port of the TypeScript json-pointer library with identical behavior,
// using modern Go generics for type safety and performance.
package jsonpointer

import (
	"reflect"
	"strconv"
)

// PathStep represents a single step in a JSON Pointer path (string or number).
// TypeScript original code:
// export type PathStep = string | number;
type PathStep = any

// Path represents a JSON Pointer path as array of steps.
// TypeScript original code:
// export type Path = readonly PathStep[];
type Path []PathStep

// internalToken represents a single token in a JSON Pointer path with precomputed data.
// This is used internally for performance optimization, not exposed in the API.
type internalToken struct {
	key   string // original key string
	index int    // precomputed array index, -1 if not a valid array index
}

// Reference represents a found reference with context.
// TypeScript original code:
//
//	export interface Reference {
//	  readonly val: unknown;
//	  readonly obj?: unknown | object | unknown[];
//	  readonly key?: string | number;
//	}
type Reference struct {
	Val any `json:"val"`
	Obj any `json:"obj,omitempty"`
	Key any `json:"key,omitempty"`
}

// ArrayReference represents a reference to an array element.
// TypeScript original code:
//
//	export interface ArrayReference<T = unknown> {
//	  readonly val: undefined | T;
//	  readonly obj: T[];
//	  readonly key: number;
//	}
type ArrayReference[T any] struct {
	// Use pointer for undefined | T semantics (nil = undefined)
	Val *T  `json:"val"`
	Obj []T `json:"obj"`
	Key int `json:"key"`
}

// ObjectReference represents a reference to an object property.
// TypeScript original code:
//
//	export interface ObjectReference<T = unknown> {
//	  readonly val: T;
//	  readonly obj: Record<string, T>;
//	  readonly key: string;
//	}
type ObjectReference[T any] struct {
	Val T            `json:"val"`
	Obj map[string]T `json:"obj"`
	Key string       `json:"key"`
}

// isArrayReference checks if a Reference points to an array element.
// TypeScript original code:
// export const isArrayReference = <T = unknown>(ref: Reference): ref is ArrayReference<T> =>
//
//	isArray(ref.obj) && typeof ref.key === 'number';
func isArrayReference(ref Reference) bool {
	if ref.Obj == nil || ref.Key == nil {
		return false
	}

	// Check if obj is a slice/array
	objType := reflect.TypeOf(ref.Obj)
	if objType.Kind() != reflect.Slice {
		return false
	}

	// Check if key is numeric
	switch key := ref.Key.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return true
	case string:
		_, err := strconv.Atoi(key)
		return err == nil
	default:
		return false
	}
}

// isArrayEnd checks if an array reference points to the end of the array.
// TypeScript original code:
// export const isArrayEnd = (ref: ArrayReference): boolean => ref.obj.length === ref.key;
func isArrayEnd[T any](ref ArrayReference[T]) bool {
	return len(ref.Obj) == ref.Key
}

// isObjectReference checks if a Reference points to an object property.
// TypeScript original code:
// export const isObjectReference = <T = unknown>(ref: Reference): ref is ObjectReference<T> =>
//
//	typeof ref.obj === 'object' && typeof ref.key === 'string';
func isObjectReference(ref Reference) bool {
	if ref.Obj == nil || ref.Key == nil {
		return false
	}

	// Check if obj is a map with string keys
	objType := reflect.TypeOf(ref.Obj)
	if objType.Kind() != reflect.Map || objType.Key().Kind() != reflect.String {
		return false
	}

	// Check if key is a string
	_, keyIsString := ref.Key.(string)
	return keyIsString
}
