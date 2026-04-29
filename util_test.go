package jsonpointer

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestParseJsonPointer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		pointer string
		want    Path
	}{
		{name: "returns path without escaped characters parsed into array", pointer: "/foo/bar", want: Path{"foo", "bar"}},
		{name: "trailing slashes result into empty string elements", pointer: "/foo///", want: Path{"foo", "", "", ""}},
		{name: "for root path returns empty array", pointer: "", want: Path{}},
		{name: "slash path returns single empty string", pointer: "/", want: Path{""}},
		{name: "un-escapes special characters", pointer: "/a~0b/c~1d/1", want: Path{"a~b", "c/d", "1"}},
		{name: "keeps permissive parsing for non-valid escape sequences", pointer: "/~2/foo~", want: Path{"~2", "foo~"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := parseJSONPointer(tc.pointer)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("parseJSONPointer() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestFormatJsonPointer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		path Path
		want string
	}{
		{name: "returns path without escaped characters parsed into array", path: Path{"foo", "bar"}, want: "/foo/bar"},
		{name: "empty string elements add trailing slashes", path: Path{"foo", "", "", ""}, want: "/foo///"},
		{name: "array with single empty string results into root element", path: Path{}, want: ""},
		{name: "two empty strings result in a single slash", path: Path{""}, want: "/"},
		{name: "escapes special characters", path: Path{"a~b", "c/d", "1"}, want: "/a~0b/c~1d/1"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := formatJSONPointer(tc.path)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestEscapeComponent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		component string
		want      string
	}{
		{name: "string without escaped characters as is", component: "foobar", want: "foobar"},
		{name: "replaces special characters", component: "foo~/", want: "foo~0~1"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := escapeComponent(tc.component)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestUnescapeComponent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		component string
		want      string
	}{
		{name: "string without escaped characters as is", component: "foobar", want: "foobar"},
		{name: "unescapes slash and tilde", component: "foo~0~1", want: "foo~/"},
		{name: "unescapes slash inside component", component: "fo~1o", want: "fo/o"},
		{name: "unescapes tilde inside component", component: "fo~0o", want: "fo~o"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := unescapeComponent(tc.component)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestIsChild(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		parent Path
		child  Path
		want   bool
	}{
		{name: "returns false if parent path is longer than child path", parent: Path{"", "foo", "bar", "baz"}, child: Path{"", "foo"}, want: false},
		{name: "returns true for real child", parent: Path{"", "foo"}, child: Path{"", "foo", "bar", "baz"}, want: true},
		{name: "returns false for different root steps", parent: Path{"", "foo"}, child: Path{"", "foo2", "bar", "baz"}, want: false},
		{name: "returns false for adjacent paths", parent: Path{"", "foo", "baz"}, child: Path{"", "foo", "bar"}, want: false},
		{name: "returns false for two roots", parent: Path{""}, child: Path{""}, want: false},
		{name: "always returns true when parent is root and child is not", parent: Path{""}, child: Path{"", "a", "b", "c", "1", "2", "3"}, want: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := IsChild(tc.parent, tc.child)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestParent(t *testing.T) {
	t.Parallel()

	t.Run("returns parent path", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			path Path
			want Path
		}{
			{name: "three components", path: Path{"foo", "bar", "baz"}, want: Path{"foo", "bar"}},
			{name: "two components", path: Path{"foo", "bar"}, want: Path{"foo"}},
			{name: "one component", path: Path{"foo"}, want: Path{}},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				got, err := Parent(tc.path)
				assert.NoError(t, err)
				if diff := cmp.Diff(tc.want, got); diff != "" {
					t.Errorf("Parent() mismatch (-want +got):\n%s", diff)
				}
			})
		}
	})

	t.Run("returns no parent error when path has no parent", func(t *testing.T) {
		t.Parallel()

		_, err := Parent(Path{})
		assert.ErrorIs(t, err, ErrNoParent)
	})
}

func TestIsValidIndex(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		index string
		want  bool
	}{
		{name: "zero", index: "0", want: true},
		{name: "single digit", index: "5", want: true},
		{name: "multiple digits", index: "10", want: true},
		{name: "array end marker", index: "-", want: true},
		{name: "leading zero", index: "01", want: false},
		{name: "letters", index: "abc", want: false},
		{name: "decimal", index: "1.5", want: false},
		{name: "negative one", index: "-1", want: false},
		{name: "negative many", index: "-5", want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := IsValidIndex(tc.index)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestIsRoot(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		path Path
		want bool
	}{
		{name: "empty path is root", path: Path{}, want: true},
		{name: "single component path is not root", path: Path{"foo"}, want: false},
		{name: "multi component path is not root", path: Path{"foo", "bar"}, want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := IsRoot(tc.path)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestIsPathEqual(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		left  Path
		right Path
		want  bool
	}{
		{name: "empty paths", left: Path{}, right: Path{}, want: true},
		{name: "single component paths", left: Path{"foo"}, right: Path{"foo"}, want: true},
		{name: "multi component paths", left: Path{"foo", "bar"}, right: Path{"foo", "bar"}, want: true},
		{name: "numeric component paths", left: Path{"foo", "0"}, right: Path{"foo", "0"}, want: true},
		{name: "different components", left: Path{"foo"}, right: Path{"bar"}, want: false},
		{name: "different lengths", left: Path{"foo", "bar"}, right: Path{"foo"}, want: false},
		{name: "different numeric components", left: Path{"foo", "0"}, right: Path{"foo", "1"}, want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := IsPathEqual(tc.left, tc.right)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestIsInteger(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{name: "zero", value: "0", want: true},
		{name: "digits", value: "123", want: true},
		{name: "many digits", value: "999", want: true},
		{name: "letters", value: "abc", want: false},
		{name: "decimal", value: "1.5", want: false},
		{name: "negative", value: "-1", want: false},
		{name: "empty", value: "", want: false},
		{name: "leading zero", value: "01", want: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := IsInteger(tc.value)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestFastAtoi(t *testing.T) {
	t.Parallel()

	t.Run("valid non-negative integers", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name  string
			input string
			want  int
		}{
			{name: "zero", input: "0", want: 0},
			{name: "one", input: "1", want: 1},
			{name: "hundreds", input: "123", want: 123},
			{name: "nines", input: "999", want: 999},
			{name: "thousand", input: "1000", want: 1000},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				got := fastAtoi(tc.input)
				assert.Equal(t, tc.want, got)
			})
		}
	})

	t.Run("invalid inputs return -1", func(t *testing.T) {
		t.Parallel()

		tests := []string{"", "abc", "-1", "01", "00", "123abc", "12.34", " 123", "123 "}
		for _, input := range tests {
			t.Run(input, func(t *testing.T) {
				t.Parallel()

				got := fastAtoi(input)
				assert.Equal(t, -1, got)
			})
		}
	})

	t.Run("overflow detection", func(t *testing.T) {
		t.Parallel()

		got := fastAtoi("99999999999999999999999999999999")
		assert.Equal(t, -1, got)
	})
}
