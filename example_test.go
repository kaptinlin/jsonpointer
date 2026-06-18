package jsonpointer_test

import (
	"fmt"
	"log"

	"github.com/kaptinlin/jsonpointer"
)

func ExamplePointer_Value() {
	doc := map[string]any{
		"users": []any{
			map[string]any{"name": "Alice"},
		},
	}

	p, err := jsonpointer.Parse("/users/0/name")
	if err != nil {
		log.Fatal(err)
	}
	name, err := p.Value(doc)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(name)
	// Output: Alice
}

func ExamplePointer_Reference() {
	doc := map[string]any{
		"foo/bar": map[string]any{
			"tilde~key": "ready",
		},
	}

	p, err := jsonpointer.Parse("/foo~1bar/tilde~0key")
	if err != nil {
		log.Fatal(err)
	}
	ref, err := p.Reference(doc)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(ref.Value())
	fmt.Println(ref.Token())
	// Output:
	// ready
	// tilde~key
}

func ExampleFromTokens() {
	p := jsonpointer.FromTokens("foo/bar", "tilde~key")

	fmt.Println(p.String())
	fmt.Println(p.Tokens())
	// Output:
	// /foo~1bar/tilde~0key
	// [foo/bar tilde~key]
}

func ExampleValue() {
	doc := map[string]any{"name": "Alice"}

	name, err := jsonpointer.Value(doc, "/name")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(name)
	// Output: Alice
}
