package jsonpointer_test

import (
	"testing"

	"github.com/kaptinlin/jsonpointer"
)

// BenchmarkFind benchmarks the Find function with various scenarios.
// Maps to: __bench__/find.js
func BenchmarkFind(b *testing.B) {
	// Complex nested structure for realistic benchmarking
	doc := map[string]any{
		"users": []any{
			map[string]any{
				"id":   1,
				"name": "Alice",
				"profile": map[string]any{
					"email": "alice@example.com",
					"settings": map[string]any{
						"notifications": map[string]any{
							"email": true,
							"sms":   false,
						},
						"privacy": map[string]any{
							"public":  true,
							"friends": false,
						},
					},
				},
			},
			map[string]any{
				"id":   2,
				"name": "Bob",
				"profile": map[string]any{
					"email": "bob@example.com",
					"settings": map[string]any{
						"notifications": map[string]any{
							"email": false,
							"sms":   true,
						},
					},
				},
			},
		},
		"metadata": map[string]any{
			"version": "1.0",
			"created": "2023-01-01",
		},
	}

	b.Run("root", func(b *testing.B) {
		path := jsonpointer.Path{}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := jsonpointer.Find(doc, path)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("simple_object_property", func(b *testing.B) {
		path := jsonpointer.Path{"metadata", "version"}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := jsonpointer.Find(doc, path)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("array_element", func(b *testing.B) {
		path := jsonpointer.Path{"users", 0, "name"}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := jsonpointer.Find(doc, path)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("deep_nested_property", func(b *testing.B) {
		path := jsonpointer.Path{"users", 0, "profile", "settings", "notifications", "email"}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := jsonpointer.Find(doc, path)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("array_end_marker", func(b *testing.B) {
		path := jsonpointer.Path{"users", "-"}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := jsonpointer.Find(doc, path)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("missing_property", func(b *testing.B) {
		path := jsonpointer.Path{"nonexistent", "property"}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := jsonpointer.Find(doc, path)
			if err == nil {
				b.Fatal("expected error")
			}
		}
	})
}

// BenchmarkFindByPointer benchmarks the optimized FindByPointer function.
// Maps to: __bench__/find.js findByPointer benchmarks
func BenchmarkFindByPointer(b *testing.B) {
	// Same complex nested structure
	doc := map[string]any{
		"users": []any{
			map[string]any{
				"id":   1,
				"name": "Alice",
				"profile": map[string]any{
					"email": "alice@example.com",
					"settings": map[string]any{
						"notifications": map[string]any{
							"email": true,
							"sms":   false,
						},
						"privacy": map[string]any{
							"public":  true,
							"friends": false,
						},
					},
				},
			},
			map[string]any{
				"id":   2,
				"name": "Bob",
				"profile": map[string]any{
					"email": "bob@example.com",
				},
			},
		},
		"metadata": map[string]any{
			"version": "1.0",
			"created": "2023-01-01",
		},
	}

	b.Run("root", func(b *testing.B) {
		pointer := ""
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := jsonpointer.FindByPointer(pointer, doc)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("simple_object_property", func(b *testing.B) {
		pointer := "/metadata/version"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := jsonpointer.FindByPointer(pointer, doc)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("array_element", func(b *testing.B) {
		pointer := "/users/0/name"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := jsonpointer.FindByPointer(pointer, doc)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("deep_nested_property", func(b *testing.B) {
		pointer := "/users/0/profile/settings/notifications/email"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := jsonpointer.FindByPointer(pointer, doc)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("array_end_marker", func(b *testing.B) {
		pointer := "/users/-"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := jsonpointer.FindByPointer(pointer, doc)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("escaped_characters", func(b *testing.B) {
		docWithEscaped := map[string]any{
			"foo/bar": "value1",
			"foo~bar": "value2",
		}
		pointer := "/foo~1bar"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := jsonpointer.FindByPointer(pointer, docWithEscaped)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("missing_property", func(b *testing.B) {
		pointer := "/nonexistent/property"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := jsonpointer.FindByPointer(pointer, doc)
			if err == nil {
				b.Fatal("expected error")
			}
		}
	})
}

// BenchmarkFindVsFindByPointer compares performance between Find and FindByPointer.
func BenchmarkFindVsFindByPointer(b *testing.B) {
	doc := map[string]any{
		"users": []any{
			map[string]any{
				"profile": map[string]any{
					"settings": map[string]any{
						"notifications": map[string]any{
							"email": true,
						},
					},
				},
			},
		},
	}

	path := jsonpointer.Path{"users", 0, "profile", "settings", "notifications", "email"}
	pointer := "/users/0/profile/settings/notifications/email"

	b.Run("Find", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := jsonpointer.Find(doc, path)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("FindByPointer", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := jsonpointer.FindByPointer(pointer, doc)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkGet benchmarks the Get function that never returns errors.
func BenchmarkGet(b *testing.B) {
	doc := map[string]any{
		"users": []any{
			map[string]any{
				"name": "Alice",
				"profile": map[string]any{
					"email": "alice@example.com",
				},
			},
		},
		"metadata": map[string]any{
			"version": "1.0",
		},
	}

	b.Run("simple_property", func(b *testing.B) {
		path := jsonpointer.Path{"metadata", "version"}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result := jsonpointer.Get(doc, path)
			if result == nil {
				b.Fatal("expected non-nil result")
			}
		}
	})

	b.Run("nested_property", func(b *testing.B) {
		path := jsonpointer.Path{"users", 0, "profile", "email"}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result := jsonpointer.Get(doc, path)
			if result == nil {
				b.Fatal("expected non-nil result")
			}
		}
	})

	b.Run("missing_property", func(b *testing.B) {
		path := jsonpointer.Path{"nonexistent", "property"}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result := jsonpointer.Get(doc, path)
			if result != nil {
				b.Fatal("expected nil result")
			}
		}
	})
}

// BenchmarkTypeGuards benchmarks the type guard functions.
func BenchmarkTypeGuards(b *testing.B) {
	arrayRef := jsonpointer.Reference{
		Val: "value",
		Obj: []any{1, 2, 3},
		Key: 1,
	}

	objectRef := jsonpointer.Reference{
		Val: "value",
		Obj: map[string]any{"key": "value"},
		Key: "key",
	}

	b.Run("IsArrayReference", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if !jsonpointer.IsArrayReference(arrayRef) {
				b.Fatal("expected true")
			}
		}
	})

	b.Run("IsObjectReference", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if !jsonpointer.IsObjectReference(objectRef) {
				b.Fatal("expected true")
			}
		}
	})
}
