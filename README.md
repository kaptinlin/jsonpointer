# JSON Pointer

[![Go Module](https://img.shields.io/badge/go-module-blue.svg)](https://golang.org/)
[![Go Reference](https://pkg.go.dev/badge/github.com/kaptinlin/jsonpointer.svg)](https://pkg.go.dev/github.com/kaptinlin/jsonpointer)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A read-only JSON Pointer (RFC 6901) library for Go built around one strict,
immutable `Pointer` value.

## Features

- **Strict RFC 6901 parsing**: malformed pointer strings return errors instead
  of becoming lookup keys.
- **Raw token construction**: build pointers from literal Go strings without
  confusing token text with pointer-string syntax.
- **Go-native traversal**: read maps, slices, arrays, structs, pointers, and
  interface-wrapped values.
- **Explicit errors**: keep `errors.Is` sentinel checks and use `errors.As` for
  pointer, token, and depth context.
- **Small API**: parse or build a `Pointer`, then ask it for a value or
  reference.
- **Fast common paths**: optimize decoded JSON shapes such as `map[string]any`
  and `[]any` before reflective fallbacks.

## Installation

```bash
go get github.com/kaptinlin/jsonpointer
```

Requires the Go version declared in `go.mod`.

## Quick Start

```go
package main

import (
	"fmt"
	"log"

	"github.com/kaptinlin/jsonpointer"
)

func main() {
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
}
```

## Core APIs

| API | Description |
| --- | --- |
| `Parse(pointer string) (Pointer, error)` | Parse a strict JSON Pointer string |
| `FromTokens(tokens ...string) (Pointer, error)` | Build a pointer from raw token strings |
| `Root() Pointer` | Return the root pointer |
| `(Pointer).String() string` | Format the canonical JSON Pointer string |
| `(Pointer).Tokens() []string` | Return a copy of the raw tokens |
| `(Pointer).Parent() (Pointer, error)` | Return the parent pointer |
| `(Pointer).Child(tokens ...string) (Pointer, error)` | Return a child pointer |
| `(Pointer).Value(doc any) (any, error)` | Resolve a value |
| `(Pointer).Reference(doc any) (Reference, error)` | Resolve a value with parent context |
| `Value(doc any, pointer string) (any, error)` | Strict one-shot value lookup |
| `ReferenceOf(doc any, pointer string) (Reference, error)` | Strict one-shot reference lookup |
| `EscapeToken(token string) string` | Escape one raw token |
| `UnescapeToken(encoded string) (string, error)` | Decode one escaped token strictly |
| `IsValidIndex(token string) bool` | Report whether a token is an array index or `-` |
| `IsInteger(str string) bool` | Report whether a string contains only decimal digits |

`Parse("/~2")` returns `ErrInvalidPointer`. `FromTokens("~2")` succeeds because
`"~2"` is literal token data, not pointer-string syntax.

## Reference Results

`Pointer.Reference` and `ReferenceOf` return a `Reference` with named accessors:

| Method | Meaning |
| --- | --- |
| `Value() any` | The resolved value |
| `Parent() (any, bool)` | The parent container when one exists |
| `Token() string` | The final token used to reach the value |
| `Pointer() Pointer` | The pointer used for traversal |

Root references have no parent and return `(nil, false)` from `Parent`.

## Examples

### Build from raw tokens

```go
p, err := jsonpointer.FromTokens("foo/bar", "tilde~key")
if err != nil {
	log.Fatal(err)
}
fmt.Println(p.String())
```

Output:

```text
/foo~1bar/tilde~0key
```

### Read with parent context

```go
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
```

### Traverse structs and pointers

```go
type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

p, err := jsonpointer.FromTokens("email")
if err != nil {
	log.Fatal(err)
}

email, err := p.Value(&User{Name: "Alice", Email: "alice@example.com"})
if err != nil {
	log.Fatal(err)
}
fmt.Println(email)
```

Struct traversal uses exported fields, JSON tag names, exact `json:"-"` hiding,
`json:"-,"` as the literal `-` name, and `encoding/json`-style dominance for
embedded fields.

## Error Handling

Common sentinel errors include:

- `ErrInvalidPointer`
- `ErrPointerTooLong`
- `ErrPathTooLong`
- `ErrKeyNotFound`
- `ErrFieldNotFound`
- `ErrInvalidIndex`
- `ErrIndexOutOfBounds`
- `ErrNilPointer`
- `ErrNotFound`
- `ErrNoParent`

Use `errors.Is` for error classes and `errors.As` when traversal context matters:

```go
value, err := jsonpointer.Value(doc, "/users/1/name")
if errors.Is(err, jsonpointer.ErrIndexOutOfBounds) {
	var pointerErr *jsonpointer.Error
	if errors.As(err, &pointerErr) {
		fmt.Println(pointerErr.Pointer(), pointerErr.Token(), pointerErr.Depth())
	}
}
_ = value
```

## Performance

The package optimizes common `map[string]any` and `[]any` reads and falls back to
reflection for typed Go values. Reuse parsed pointers when resolving the same
location repeatedly.

See [benchmarks/README.md](benchmarks/README.md) for benchmark coverage.

Run benchmarks with:

```bash
task bench
```

## Development

```bash
task test          # Run package tests with the race detector
task lint          # Run golangci-lint and tidy checks
task yamllint      # Lint YAML files
task bench         # Run benchmark suites
```

Run the demo program with:

```bash
go run ./examples
```

For development workflow and package contracts, see [AGENTS.md](AGENTS.md) and
[`SPECS/`](SPECS/).

## Contributing

Contributions are welcome. Keep `README.md`, `example_test.go`, and the relevant
`SPECS/` documents aligned when public behavior changes.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file
for details.
