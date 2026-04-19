// Package main demonstrates the jsonpointer API.
package main

import (
	"fmt"
	"log"

	"github.com/kaptinlin/jsonpointer"
)

type user struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	doc := map[string]any{
		"users": []any{
			map[string]any{
				"name":  "Alice",
				"email": "alice@example.com",
			},
		},
		"foo/bar": map[string]any{
			"tilde~key": "ready",
		},
	}

	name, err := jsonpointer.GetByPointer(doc, "/users/0/name")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("name:", name)

	ref, err := jsonpointer.FindByPointer(doc, "/foo~1bar/tilde~0key")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("escaped value:", ref.Val)
	fmt.Println("escaped key:", ref.Key)

	email, err := jsonpointer.Get(&user{Name: "Bob", Email: "bob@example.com"}, "email")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("struct email:", email)

	path := jsonpointer.Parse("/users/0/name")
	fmt.Println("path:", path)
	fmt.Println("pointer:", jsonpointer.Format(path...))
	fmt.Println("valid pointer:", jsonpointer.Validate("/users/0/name") == nil)
}
