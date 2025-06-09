package jsonpointer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestParseJsonPointer tests JSON Pointer string parsing.
// Maps to: util.parseJsonPointer.spec.ts
func TestParseJsonPointer(t *testing.T) {
	t.Run("returns path without escaped characters parsed into array", func(t *testing.T) {
		res := parseJsonPointer("/foo/bar")
		expected := Path{"foo", "bar"}
		assert.True(t, IsPathEqual(res, expected), "Expected %v, got %v", expected, res)
	})

	t.Run("trailing slashes result into empty string elements", func(t *testing.T) {
		res := parseJsonPointer("/foo///")
		expected := Path{"foo", "", "", ""}
		assert.True(t, IsPathEqual(res, expected), "Expected %v, got %v", expected, res)
	})

	t.Run("for root path returns empty array", func(t *testing.T) {
		res := parseJsonPointer("")
		expected := Path{}
		assert.True(t, IsPathEqual(res, expected), "Expected %v, got %v", expected, res)
	})

	t.Run("slash path \"/\" return single empty string", func(t *testing.T) {
		res := parseJsonPointer("/")
		expected := Path{""}
		assert.True(t, IsPathEqual(res, expected), "Expected %v, got %v", expected, res)
	})

	t.Run("un-escapes special characters", func(t *testing.T) {
		res := parseJsonPointer("/a~0b/c~1d/1")
		expected := Path{"a~b", "c/d", "1"}
		assert.True(t, IsPathEqual(res, expected), "Expected %v, got %v", expected, res)
	})
}

// TestFormatJsonPointer tests path array formatting to JSON Pointer string.
// Maps to: util.formatJsonPointer.spec.ts
func TestFormatJsonPointer(t *testing.T) {
	t.Run("returns path without escaped characters parsed into array", func(t *testing.T) {
		res := formatJsonPointer(Path{"foo", "bar"})
		expected := "/foo/bar"
		assert.Equal(t, expected, res)
	})

	t.Run("empty string elements add trailing slashes", func(t *testing.T) {
		res := formatJsonPointer(Path{"foo", "", "", ""})
		expected := "/foo///"
		assert.Equal(t, expected, res)
	})

	t.Run("array with single empty string results into root element", func(t *testing.T) {
		res := formatJsonPointer(Path{})
		expected := ""
		assert.Equal(t, expected, res)
	})

	t.Run("two empty strings result in a single slash \"/\"", func(t *testing.T) {
		res := formatJsonPointer(Path{""})
		expected := "/"
		assert.Equal(t, expected, res)
	})

	t.Run("escapes special characters", func(t *testing.T) {
		res := formatJsonPointer(Path{"a~b", "c/d", "1"})
		expected := "/a~0b/c~1d/1"
		assert.Equal(t, expected, res)
	})
}

// TestEscapeComponent tests path component escaping.
// Maps to: util.escapeComponent.spec.ts
func TestEscapeComponent(t *testing.T) {
	t.Run("string without escaped characters as is", func(t *testing.T) {
		res := escapeComponent("foobar")
		expected := "foobar"
		assert.Equal(t, expected, res)
	})

	t.Run("replaces special characters", func(t *testing.T) {
		res := escapeComponent("foo~/")
		expected := "foo~0~1"
		assert.Equal(t, expected, res)
	})
}

// TestUnescapeComponent tests path component unescaping.
// Maps to: util.unescapeComponent.spec.ts
func TestUnescapeComponent(t *testing.T) {
	t.Run("string without escaped characters as is", func(t *testing.T) {
		res := unescapeComponent("foobar")
		expected := "foobar"
		assert.Equal(t, expected, res)
	})

	t.Run("replaces special characters", func(t *testing.T) {
		tests := []struct {
			input    string
			expected string
		}{
			{"foo~0~1", "foo~/"},
			{"fo~1o", "fo/o"},
			{"fo~0o", "fo~o"},
		}

		for _, test := range tests {
			res := unescapeComponent(test.input)
			assert.Equal(t, test.expected, res, "unescapeComponent(%s)", test.input)
		}
	})
}

// TestIsChild tests parent-child path relationship checking.
// Maps to: util.isChild.spec.ts
func TestIsChild(t *testing.T) {
	t.Run("returns false if parent path is longer than child path", func(t *testing.T) {
		res := IsChild(Path{"", "foo", "bar", "baz"}, Path{"", "foo"})
		assert.False(t, res)
	})

	t.Run("returns true for real child", func(t *testing.T) {
		res := IsChild(Path{"", "foo"}, Path{"", "foo", "bar", "baz"})
		assert.True(t, res)
	})

	t.Run("returns false for different root steps", func(t *testing.T) {
		res := IsChild(Path{"", "foo"}, Path{"", "foo2", "bar", "baz"})
		assert.False(t, res)
	})

	t.Run("returns false for adjacent paths", func(t *testing.T) {
		res := IsChild(Path{"", "foo", "baz"}, Path{"", "foo", "bar"})
		assert.False(t, res)
	})

	t.Run("returns false for two roots", func(t *testing.T) {
		res := IsChild(Path{""}, Path{""})
		assert.False(t, res)
	})

	t.Run("always returns true when parent is root and child is not", func(t *testing.T) {
		res := IsChild(Path{""}, Path{"", "a", "b", "c", "1", "2", "3"})
		assert.True(t, res)
	})
}

// TestParent tests parent path extraction.
// Maps to: util.parent.spec.ts
func TestParent(t *testing.T) {
	t.Run("returns parent path", func(t *testing.T) {
		tests := []struct {
			input    Path
			expected Path
		}{
			{Path{"foo", "bar", "baz"}, Path{"foo", "bar"}},
			{Path{"foo", "bar"}, Path{"foo"}},
			{Path{"foo"}, Path{}},
		}

		for _, test := range tests {
			res, err := Parent(test.input)
			assert.NoError(t, err, "Parent(%v)", test.input)
			assert.True(t, IsPathEqual(res, test.expected), "Parent(%v): expected %v, got %v", test.input, test.expected, res)
		}
	})

	t.Run("throws when path has no parent", func(t *testing.T) {
		_, err := Parent(Path{})
		assert.Error(t, err)
		assert.Equal(t, ErrNoParent, err)
	})
}

// TestToPath tests path conversion utilities.
func TestToPath(t *testing.T) {
	t.Run("converts string pointer to path", func(t *testing.T) {
		res := ToPath("/foo/bar")
		expected := Path{"foo", "bar"}
		assert.True(t, IsPathEqual(res, expected))
	})

	t.Run("returns path as-is", func(t *testing.T) {
		input := Path{"foo", "bar"}
		res := ToPath(input)
		assert.True(t, IsPathEqual(res, input))
	})

	t.Run("converts string slice to path", func(t *testing.T) {
		input := []string{"foo", "bar"}
		res := ToPath(input)
		expected := Path{"foo", "bar"}
		assert.True(t, IsPathEqual(res, expected))
	})
}

// TestIsValidIndex tests array index validation.
func TestIsValidIndex(t *testing.T) {
	t.Run("valid string indices", func(t *testing.T) {
		assert.True(t, IsValidIndex("0"))
		assert.True(t, IsValidIndex("5"))
		assert.True(t, IsValidIndex("10"))
		assert.True(t, IsValidIndex("-")) // Array end marker
	})

	t.Run("invalid string indices", func(t *testing.T) {
		assert.False(t, IsValidIndex("01")) // Leading zero
		assert.False(t, IsValidIndex("abc"))
		assert.False(t, IsValidIndex("1.5"))
		assert.False(t, IsValidIndex("-1")) // Negative
		assert.False(t, IsValidIndex("-5")) // Negative
	})
}

// TestIsRoot tests root path detection.
func TestIsRoot(t *testing.T) {
	t.Run("empty path is root", func(t *testing.T) {
		assert.True(t, IsRoot(Path{}))
	})

	t.Run("non-empty path is not root", func(t *testing.T) {
		assert.False(t, IsRoot(Path{"foo"}))
		assert.False(t, IsRoot(Path{"foo", "bar"}))
	})
}

// TestIsPathEqual tests path equality checking.
func TestIsPathEqual(t *testing.T) {
	t.Run("equal paths", func(t *testing.T) {
		assert.True(t, IsPathEqual(Path{}, Path{}))
		assert.True(t, IsPathEqual(Path{"foo"}, Path{"foo"}))
		assert.True(t, IsPathEqual(Path{"foo", "bar"}, Path{"foo", "bar"}))
		assert.True(t, IsPathEqual(Path{"foo", "0"}, Path{"foo", "0"}))
	})

	t.Run("unequal paths", func(t *testing.T) {
		assert.False(t, IsPathEqual(Path{"foo"}, Path{"bar"}))
		assert.False(t, IsPathEqual(Path{"foo", "bar"}, Path{"foo"}))
		assert.False(t, IsPathEqual(Path{"foo", "0"}, Path{"foo", "1"}))
	})
}

// TestIsInteger tests integer string validation.
func TestIsInteger(t *testing.T) {
	t.Run("valid integers", func(t *testing.T) {
		assert.True(t, IsInteger("0"))
		assert.True(t, IsInteger("123"))
		assert.True(t, IsInteger("999"))
	})

	t.Run("invalid integers", func(t *testing.T) {
		assert.False(t, IsInteger("abc"))
		assert.False(t, IsInteger("1.5"))
		assert.False(t, IsInteger("-1")) // Negative
		assert.False(t, IsInteger(""))   // Empty
		// Note: "01" is valid for IsInteger (it only checks digits),
		// but invalid for IsValidIndex (which checks format)
		assert.True(t, IsInteger("01")) // Valid digits, but not valid index format
	})
}

// TestFastAtoi tests the fastAtoi function for array index optimization.
func TestFastAtoi(t *testing.T) {
	t.Run("valid non-negative integers", func(t *testing.T) {
		testCases := []struct {
			input    string
			expected int
		}{
			{"0", 0},
			{"1", 1},
			{"123", 123},
			{"999", 999},
			{"1000", 1000},
		}

		for _, tc := range testCases {
			result := fastAtoi(tc.input)
			assert.Equal(t, tc.expected, result, "fastAtoi(%q) should return %d, got %d", tc.input, tc.expected, result)
		}
	})

	t.Run("invalid inputs return -1", func(t *testing.T) {
		testCases := []string{
			"",       // empty string
			"abc",    // non-numeric
			"-1",     // negative
			"01",     // leading zero
			"00",     // multiple leading zeros
			"123abc", // mixed
			"12.34",  // decimal
			" 123",   // leading space
			"123 ",   // trailing space
		}

		for _, input := range testCases {
			result := fastAtoi(input)
			assert.Equal(t, -1, result, "fastAtoi(%q) should return -1, got %d", input, result)
		}
	})

	t.Run("overflow detection", func(t *testing.T) {
		// Test a very large number that would cause overflow
		largeNumber := "99999999999999999999999999999999"
		result := fastAtoi(largeNumber)
		assert.Equal(t, -1, result, "fastAtoi should detect overflow and return -1")
	})
}
