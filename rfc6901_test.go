package jsonpointer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRFC6901Examples(t *testing.T) {
	doc := map[string]any{
		"foo":  []any{"bar", "baz"},
		"":     0,
		"a/b":  1,
		"c%d":  2,
		"e^f":  3,
		"g|h":  4,
		"i\\j": 5,
		"k\"l": 6,
		" ":    7,
		"m~n":  8,
	}

	tests := []struct {
		pointer string
		want    any
	}{
		{pointer: "", want: doc},
		{pointer: "/foo", want: []any{"bar", "baz"}},
		{pointer: "/foo/0", want: "bar"},
		{pointer: "/", want: 0},
		{pointer: "/a~1b", want: 1},
		{pointer: "/c%d", want: 2},
		{pointer: "/e^f", want: 3},
		{pointer: "/g|h", want: 4},
		{pointer: "/i\\j", want: 5},
		{pointer: "/k\"l", want: 6},
		{pointer: "/ ", want: 7},
		{pointer: "/m~0n", want: 8},
	}

	for _, tt := range tests {
		t.Run(tt.pointer, func(t *testing.T) {
			p, err := Parse(tt.pointer)
			require.NoError(t, err)
			assert.Equal(t, tt.pointer, p.String())

			got, err := p.Value(doc)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
