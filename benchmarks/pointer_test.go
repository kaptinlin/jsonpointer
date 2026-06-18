package benchmarks

import (
	"testing"

	"github.com/kaptinlin/jsonpointer"
)

var (
	benchValue any
	benchRef   jsonpointer.Reference
	benchPtr   jsonpointer.Pointer
	benchText  string
)

func BenchmarkParse(b *testing.B) {
	tests := []struct {
		name    string
		pointer string
	}{
		{name: "root", pointer: ""},
		{name: "shallow", pointer: "/name"},
		{name: "deep", pointer: "/users/0/profile/settings/notifications/email"},
		{name: "escaped", pointer: "/foo~1bar/tilde~0key"},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			for b.Loop() {
				p, err := jsonpointer.Parse(tt.pointer)
				if err != nil {
					b.Fatal(err)
				}
				benchPtr = p
			}
		})
	}
}

func BenchmarkFromTokens(b *testing.B) {
	for b.Loop() {
		p := jsonpointer.FromTokens("users", "0", "profile", "settings", "notifications", "email")
		benchPtr = p
	}
}

func BenchmarkPointerString(b *testing.B) {
	p := mustPointer(b, "/foo~1bar/tilde~0key")
	for b.Loop() {
		benchText = p.String()
	}
}

func BenchmarkPointerValue(b *testing.B) {
	doc := benchmarkDoc()
	tests := []struct {
		name    string
		pointer string
	}{
		{name: "root", pointer: ""},
		{name: "shallow", pointer: "/metadata/version"},
		{name: "deep", pointer: "/users/0/profile/settings/notifications/email"},
		{name: "escaped", pointer: "/foo~1bar/tilde~0key"},
	}

	for _, tt := range tests {
		p := mustPointer(b, tt.pointer)
		b.Run(tt.name, func(b *testing.B) {
			for b.Loop() {
				value, err := p.Value(doc)
				if err != nil {
					b.Fatal(err)
				}
				benchValue = value
			}
		})
	}
}

func BenchmarkPointerReference(b *testing.B) {
	doc := benchmarkDoc()
	p := mustPointer(b, "/users/0/profile/settings/notifications/email")

	for b.Loop() {
		ref, err := p.Reference(doc)
		if err != nil {
			b.Fatal(err)
		}
		benchRef = ref
	}
}

func BenchmarkOneShotValue(b *testing.B) {
	doc := benchmarkDoc()
	for b.Loop() {
		value, err := jsonpointer.Value(doc, "/users/0/profile/settings/notifications/email")
		if err != nil {
			b.Fatal(err)
		}
		benchValue = value
	}
}

func BenchmarkTraversalErrors(b *testing.B) {
	doc := benchmarkDoc()
	p := mustPointer(b, "/users/9/name")

	for b.Loop() {
		_, err := p.Value(doc)
		if err == nil {
			b.Fatal("expected error")
		}
	}
}

func BenchmarkTypedContainerTraversal(b *testing.B) {
	type mapKey string

	doc := map[string]any{
		"groups": []map[mapKey]string{
			{"theme": "dark"},
		},
	}
	p := mustPointer(b, "/groups/0/theme")

	for b.Loop() {
		value, err := p.Value(doc)
		if err != nil {
			b.Fatal(err)
		}
		benchValue = value
	}
}

func benchmarkDoc() map[string]any {
	return map[string]any{
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
					},
				},
			},
		},
		"metadata": map[string]any{
			"version": "1.0",
		},
		"foo/bar": map[string]any{
			"tilde~key": "ready",
		},
	}
}

func mustPointer(tb testing.TB, pointer string) jsonpointer.Pointer {
	tb.Helper()

	p, err := jsonpointer.Parse(pointer)
	if err != nil {
		tb.Fatal(err)
	}
	return p
}
