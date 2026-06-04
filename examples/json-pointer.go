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
	namePointer, err := jsonpointer.Parse(data.namePath)
	if err != nil {
		return err
	}
	name, err := namePointer.Value(data.doc)
	if err != nil {
		return err
	}
	escapePointer, err := jsonpointer.Parse(data.escapePath)
	if err != nil {
		return err
	}
	ref, err := escapePointer.Reference(data.doc)
	if err != nil {
		return err
	}
	emailPointer, err := jsonpointer.FromTokens("email")
	if err != nil {
		return err
	}
	email, err := emailPointer.Value(&data.user)
	if err != nil {
		return err
	}

	output := fmt.Sprintf(
		"name: %v\nescaped value: %v\nescaped key: %v\nstruct email: %v\ntokens: %v\npointer: %v\n",
		name,
		ref.Value(),
		ref.Token(),
		email,
		namePointer.Tokens(),
		namePointer.String(),
	)
	_, err = io.WriteString(out, output)
	return err
}
