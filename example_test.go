package jsonpointer_test

import (
	"fmt"
	"log"

	"github.com/kaptinlin/jsonpointer"
)

func ExampleGet() {
	doc := map[string]any{
		"users": []any{
			map[string]any{"name": "Alice"},
		},
	}

	name, err := jsonpointer.Get(doc, "users", "0", "name")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(name)
	// Output: Alice
}

func ExampleGetByPointer() {
	doc := map[string]any{
		"users": []any{
			map[string]any{"name": "Alice"},
		},
	}

	name, err := jsonpointer.GetByPointer(doc, "/users/0/name")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(name)
	// Output: Alice
}

func ExampleFindByPointer() {
	doc := map[string]any{
		"foo/bar": map[string]any{
			"tilde~key": "ready",
		},
	}

	ref, err := jsonpointer.FindByPointer(doc, "/foo~1bar/tilde~0key")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(ref.Val)
	fmt.Println(ref.Key)
	// Output:
	// ready
	// tilde~key
}

func ExampleGet_struct() {
	type user struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	val, err := jsonpointer.Get(&user{Name: "Alice", Email: "alice@example.com"}, "email")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(val)
	// Output: alice@example.com
}

func ExampleParse() {
	path := jsonpointer.Parse("/foo~1bar/tilde~0key")
	fmt.Println(path)
	fmt.Println(jsonpointer.Format(path...))
	// Output:
	// [foo/bar tilde~key]
	// /foo~1bar/tilde~0key
}

func ExampleValidate() {
	fmt.Println(jsonpointer.Validate("/users/0/name") == nil)
	fmt.Println(jsonpointer.Validate("users/0/name") == nil)
	// Output:
	// true
	// false
}
