// Package main demonstrates the jsonpointer API.
package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/kaptinlin/jsonpointer"
)

type user struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type exampleData struct {
	doc        map[string]any
	user       user
	namePath   string
	escapePath string
}

func main() {
	if err := run(os.Stdout); err != nil {
		log.Fatal(err)
	}
}

func run(out io.Writer) error {
	return writeExample(out, defaultExampleData())
}

func defaultExampleData() exampleData {
	return exampleData{
		doc: map[string]any{
			"users": []any{
				map[string]any{
					"name":  "Alice",
					"email": "alice@example.com",
				},
			},
			"foo/bar": map[string]any{
				"tilde~key": "ready",
			},
		},
		user:       user{Name: "Bob", Email: "bob@example.com"},
		namePath:   "/users/0/name",
		escapePath: "/foo~1bar/tilde~0key",
	}
}

func writeExample(out io.Writer, data exampleData) error {
	name, err := jsonpointer.GetByPointer(data.doc, data.namePath)
	if err != nil {
		return err
	}
	ref, err := jsonpointer.FindByPointer(data.doc, data.escapePath)
	if err != nil {
		return err
	}
	email, err := jsonpointer.Get(&data.user, "email")
	if err != nil {
		return err
	}

	path := jsonpointer.Parse(data.namePath)
	output := fmt.Sprintf(
		"name: %v\nescaped value: %v\nescaped key: %v\nstruct email: %v\npath: %v\npointer: %v\nvalid pointer: %v\n",
		name,
		ref.Val,
		ref.Key,
		email,
		path,
		jsonpointer.Format(path...),
		jsonpointer.Validate(data.namePath) == nil,
	)
	_, err = io.WriteString(out, output)
	return err
}
