package jsonpointer

import (
	"strings"
	"testing"
)

// FuzzParseJsonPointer performs fuzz testing on JSON Pointer parsing to discover edge cases
// and ensure parser robustness with various inputs including malformed, unicode, and edge case strings.
func FuzzParseJsonPointer(f *testing.F) {
	// Seed corpus with known valid and edge case JSON Pointers
	// These provide good starting points for the fuzzer to mutate
	f.Add("")                                     // root pointer
	f.Add("/")                                    // root with slash
	f.Add("/foo")                                 // simple pointer
	f.Add("/foo/bar")                             // nested pointer
	f.Add("/users/0/name")                        // array access
	f.Add("/~0~1")                                // escaped characters (~ and /)
	f.Add("/foo%20bar")                           // URL encoded characters
	f.Add("/中文路径")                                // Unicode characters
	f.Add("/emoji/😀")                             // Emoji
	f.Add("/special/chars!@#$%^&*()")             // Special characters
	f.Add("/array/-")                             // Array end marker
	f.Add("/deeply/nested/path/with/many/levels") // Deep nesting
	f.Add("/with spaces/and\ttabs")               // Whitespace characters
	f.Add("/with\nline\rbreaks")                  // Line breaks
	f.Add("/~0tilde~1slash")                      // Mixed escape sequences

	f.Fuzz(func(t *testing.T, pointer string) {
		// Parse the fuzzed pointer string
		path := Parse(pointer)

		// Property test 1: Roundtrip consistency for valid JSON Pointers
		// Valid JSON Pointers should parse and format back to equivalent form
		if len(pointer) > 0 && pointer[0] == '/' {
			formatted := Format(path...)
			parsed2 := Parse(formatted)

			// The two parsed paths should be identical
			if !IsPathEqual(path, parsed2) {
				t.Errorf("roundtrip consistency failed: %q -> %v -> %q -> %v",
					pointer, path, formatted, parsed2)
			}
		}

		// Property test 2: Basic validation - parsing should not panic
		// This is the most important property for fuzz testing
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("parsing panicked with pointer %q: %v", pointer, r)
			}
		}()

		// Property test 3: Validate that parsing is deterministic
		// Multiple calls to Parse with same input should return identical results
		path2 := Parse(pointer)
		if !IsPathEqual(path, path2) {
			t.Errorf("parsing not deterministic: %q -> %v vs %v", pointer, path, path2)
		}

		// Property test 4: Format should produce valid JSON Pointer syntax
		// Formatted result should start with "/" for non-empty paths
		if len(path) > 0 {
			formatted := Format(path...)
			if !strings.HasPrefix(formatted, "/") {
				t.Errorf("formatted pointer doesn't start with '/': %q -> %v -> %q", pointer, path, formatted)
			}
		}

		// Property test 5: Escape/unescape should be inverses
		// For each component, escape then unescape should return original
		for i, component := range path {
			escaped := Escape(component)
			unescaped := Unescape(escaped)
			if unescaped != component {
				t.Errorf("escape/unescape not inverse for component %d: %q -> %q -> %q",
					i, component, escaped, unescaped)
			}
		}
	})
}

// FuzzValidateJsonPointer performs fuzz testing on JSON Pointer validation
// to ensure the validator correctly identifies valid and invalid pointers.
func FuzzValidateJsonPointer(f *testing.F) {
	// Seed with various valid and invalid patterns
	f.Add("")      // valid: root
	f.Add("/foo")  // valid: simple
	f.Add("foo")   // invalid: no leading slash
	f.Add("/~0~1") // valid: escaped
	f.Add("/~2")   // invalid: bad escape
	f.Add("/foo~") // invalid: incomplete escape

	f.Fuzz(func(t *testing.T, pointer string) {
		err := Validate(pointer)

		// Property test: Valid pointers according to Validate should parse without issues
		if err == nil {
			// If validation passes, parsing should not panic
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("parser panicked on validated pointer %q: %v", pointer, r)
				}
			}()

			path := Parse(pointer)

			// Valid empty string should result in empty path
			if pointer == "" && len(path) != 0 {
				t.Errorf("empty pointer should parse to empty path, got %v", path)
			}

			// Valid non-empty pointer should start with "/"
			if pointer != "" && pointer[0] != '/' {
				t.Errorf("validator allowed invalid pointer without leading slash: %q", pointer)
			}
		}

		// Property test: Invalid pointers should not parse to valid JSON Pointer format
		if err != nil && len(pointer) > 0 {
			// Most invalid pointers should either not start with "/" or have invalid escapes
			hasValidPrefix := pointer[0] == '/'
			hasInvalidEscape := strings.Contains(pointer, "~") &&
				!strings.Contains(pointer, "~0") && !strings.Contains(pointer, "~1")

			if hasValidPrefix && !hasInvalidEscape {
				// Might be an edge case where our validation is too strict
				t.Logf("validation rejected potentially valid pointer: %q (error: %v)", pointer, err)
			}
		}
	})
}

// FuzzEscapeUnescape performs fuzz testing on escape/unescape functions
// to ensure they correctly handle all character combinations.
func FuzzEscapeUnescape(f *testing.F) {
	// Seed with various string patterns
	f.Add("")           // empty string
	f.Add("simple")     // no special chars
	f.Add("with/slash") // contains slash
	f.Add("with~tilde") // contains tilde
	f.Add("~/both")     // contains both
	f.Add("~0~1")       // already escaped
	f.Add("unicode中文")  // unicode
	f.Add("emoji😀test") // emoji

	f.Fuzz(func(t *testing.T, input string) {
		// Property test: Escape then unescape should return original
		escaped := Escape(input)
		unescaped := Unescape(escaped)

		if unescaped != input {
			t.Errorf("escape/unescape roundtrip failed: %q -> %q -> %q", input, escaped, unescaped)
		}

		// Property test: Escaped string should not contain raw tildes or slashes
		// (except as part of escape sequences ~0 and ~1)
		if strings.Contains(input, "~") || strings.Contains(input, "/") {
			// Escaped version should have replaced them
			if strings.Contains(escaped, "/") && !strings.Contains(escaped, "~1") {
				t.Errorf("escape failed to replace '/': input %q -> escaped %q", input, escaped)
			}

			// Count raw tildes (not part of ~0 or ~1)
			rawTildes := 0
			for i, r := range escaped {
				if r == '~' {
					// Check if it's followed by 0 or 1
					if i+1 < len(escaped) && (escaped[i+1] == '0' || escaped[i+1] == '1') {
						// This is a valid escape sequence, skip the next char
						continue
					}
					rawTildes++
				}
			}

			if rawTildes > 0 {
				t.Errorf("escape left raw tildes: input %q -> escaped %q", input, escaped)
			}
		}

		// Property test: Double escaping should be idempotent after first unescape
		doubleEscaped := Escape(escaped)
		if Unescape(doubleEscaped) != escaped {
			t.Errorf("double escape not handled correctly: %q -> %q -> %q",
				input, escaped, doubleEscaped)
		}
	})
}
