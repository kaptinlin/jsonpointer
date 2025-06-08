package benchmarks

import (
	"testing"

	"github.com/kaptinlin/jsonpointer"
)

// BenchmarkParseJsonPointer benchmarks the Parse function.
// Maps to: __bench__/parseJsonPointer.ts
func BenchmarkParseJsonPointer(b *testing.B) {
	b.Run("root_pointer", func(b *testing.B) {
		pointer := ""
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result := jsonpointer.Parse(pointer)
			if len(result) != 0 {
				b.Fatal("expected empty path")
			}
		}
	})

	b.Run("simple_pointer", func(b *testing.B) {
		pointer := "/foo"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result := jsonpointer.Parse(pointer)
			if len(result) != 1 {
				b.Fatal("expected single step path")
			}
		}
	})

	b.Run("nested_pointer", func(b *testing.B) {
		pointer := "/foo/bar/baz"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result := jsonpointer.Parse(pointer)
			if len(result) != 3 {
				b.Fatal("expected three step path")
			}
		}
	})

	b.Run("deep_nested_pointer", func(b *testing.B) {
		pointer := "/users/0/profile/settings/notifications/email/enabled"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result := jsonpointer.Parse(pointer)
			if len(result) != 7 {
				b.Fatal("expected seven step path")
			}
		}
	})

	b.Run("escaped_characters", func(b *testing.B) {
		pointer := "/foo~1bar/baz~0qux"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result := jsonpointer.Parse(pointer)
			if len(result) != 2 {
				b.Fatal("expected two step path")
			}
		}
	})

	b.Run("complex_escaped", func(b *testing.B) {
		pointer := "/a~1b/c~0d/e~1f~0g"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result := jsonpointer.Parse(pointer)
			if len(result) != 3 {
				b.Fatal("expected three step path")
			}
		}
	})

	b.Run("array_indices", func(b *testing.B) {
		pointer := "/0/1/2/3/4/5/6/7/8/9"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result := jsonpointer.Parse(pointer)
			if len(result) != 10 {
				b.Fatal("expected ten step path")
			}
		}
	})

	b.Run("array_end_marker", func(b *testing.B) {
		pointer := "/users/-"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result := jsonpointer.Parse(pointer)
			if len(result) != 2 {
				b.Fatal("expected two step path")
			}
		}
	})
}

// BenchmarkFormatJsonPointer benchmarks the formatJsonPointer function.
// Maps to: __bench__/parseJsonPointer.ts formatJsonPointer benchmarks
func BenchmarkFormatJsonPointer(b *testing.B) {
	b.Run("root_path", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result := jsonpointer.Format()
			if result != "" {
				b.Fatal("expected empty string")
			}
		}
	})

	b.Run("simple_path", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result := jsonpointer.Format("foo")
			if result != "/foo" {
				b.Fatal("expected '/foo'")
			}
		}
	})

	b.Run("nested_path", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result := jsonpointer.Format("foo", "bar", "baz")
			if result != "/foo/bar/baz" {
				b.Fatal("expected '/foo/bar/baz'")
			}
		}
	})

	b.Run("deep_nested_path", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result := jsonpointer.Format("users", 0, "profile", "settings", "notifications", "email", "enabled")
			if result != "/users/0/profile/settings/notifications/email/enabled" {
				b.Fatal("expected deep nested pointer")
			}
		}
	})

	b.Run("path_with_special_chars", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result := jsonpointer.Format("foo/bar", "baz~qux")
			if result != "/foo~1bar/baz~0qux" {
				b.Fatal("expected escaped pointer")
			}
		}
	})

	b.Run("array_indices", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result := jsonpointer.Format(0, 1, 2, 3, 4, 5, 6, 7, 8, 9)
			if result != "/0/1/2/3/4/5/6/7/8/9" {
				b.Fatal("expected array indices pointer")
			}
		}
	})

	b.Run("array_end_marker", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result := jsonpointer.Format("users", "-")
			if result != "/users/-" {
				b.Fatal("expected array end marker pointer")
			}
		}
	})
}

// BenchmarkParseFormatRoundtrip benchmarks parse->format roundtrip operations.
func BenchmarkParseFormatRoundtrip(b *testing.B) {
	testPointers := []string{
		"",
		"/foo",
		"/foo/bar",
		"/foo/bar/baz",
		"/users/0/profile/settings/notifications/email",
		"/foo~1bar/baz~0qux",
		"/a~1b/c~0d/e~1f~0g",
		"/0/1/2/3/4/5/6/7/8/9",
		"/users/-",
	}

	b.Run("parse_format_roundtrip", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			pointer := testPointers[i%len(testPointers)]
			path := jsonpointer.Parse(pointer)
			result := jsonpointer.Format(path...)
			if result != pointer {
				b.Fatalf("roundtrip failed: %s != %s", result, pointer)
			}
		}
	})
}

// BenchmarkEscapeComponent benchmarks the escapeComponent function.
func BenchmarkEscapeComponent(b *testing.B) {
	b.Run("no_escape_needed", func(b *testing.B) {
		component := "simple"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result := jsonpointer.Escape(component)
			if result != component {
				b.Fatal("expected no changes")
			}
		}
	})

	b.Run("escape_slash", func(b *testing.B) {
		component := "foo/bar"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result := jsonpointer.Escape(component)
			if result != "foo~1bar" {
				b.Fatal("expected escaped slash")
			}
		}
	})

	b.Run("escape_tilde", func(b *testing.B) {
		component := "foo~bar"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result := jsonpointer.Escape(component)
			if result != "foo~0bar" {
				b.Fatal("expected escaped tilde")
			}
		}
	})

	b.Run("escape_both", func(b *testing.B) {
		component := "foo~bar/baz"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result := jsonpointer.Escape(component)
			if result != "foo~0bar~1baz" {
				b.Fatal("expected both escaped")
			}
		}
	})

	b.Run("complex_escaping", func(b *testing.B) {
		component := "~foo~/bar~/"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result := jsonpointer.Escape(component)
			if result != "~0foo~0~1bar~0~1" {
				b.Fatal("expected complex escaping")
			}
		}
	})
}

// BenchmarkUnescapeComponent benchmarks the unescapeComponent function.
func BenchmarkUnescapeComponent(b *testing.B) {
	b.Run("no_unescape_needed", func(b *testing.B) {
		component := "simple"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result := jsonpointer.Unescape(component)
			if result != component {
				b.Fatal("expected no changes")
			}
		}
	})

	b.Run("unescape_slash", func(b *testing.B) {
		component := "foo~1bar"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result := jsonpointer.Unescape(component)
			if result != "foo/bar" {
				b.Fatal("expected unescaped slash")
			}
		}
	})

	b.Run("unescape_tilde", func(b *testing.B) {
		component := "foo~0bar"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result := jsonpointer.Unescape(component)
			if result != "foo~bar" {
				b.Fatal("expected unescaped tilde")
			}
		}
	})

	b.Run("unescape_both", func(b *testing.B) {
		component := "foo~0bar~1baz"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result := jsonpointer.Unescape(component)
			if result != "foo~bar/baz" {
				b.Fatal("expected both unescaped")
			}
		}
	})

	b.Run("complex_unescaping", func(b *testing.B) {
		component := "~0foo~0~1bar~0~1"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result := jsonpointer.Unescape(component)
			if result != "~foo~/bar~/" {
				b.Fatal("expected complex unescaping")
			}
		}
	})
}

// BenchmarkEscapeUnescapeRoundtrip benchmarks escape->unescape roundtrip operations.
func BenchmarkEscapeUnescapeRoundtrip(b *testing.B) {
	testComponents := []string{
		"simple",
		"foo/bar",
		"foo~bar",
		"foo~bar/baz",
		"~foo~/bar~/",
		"complex/path~with/both~chars",
		"",
		"~",
		"/",
		"~/",
	}

	b.Run("escape_unescape_roundtrip", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			component := testComponents[i%len(testComponents)]
			escaped := jsonpointer.Escape(component)
			result := jsonpointer.Unescape(escaped)
			if result != component {
				b.Fatalf("roundtrip failed: %s != %s", result, component)
			}
		}
	})
}
