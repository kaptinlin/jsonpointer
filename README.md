# JSON Pointer

[![Go Reference](https://pkg.go.dev/badge/github.com/kaptinlin/jsonpointer.svg)](https://pkg.go.dev/github.com/kaptinlin/jsonpointer)
[![Go Report Card](https://goreportcard.com/badge/github.com/kaptinlin/jsonpointer)](https://goreportcard.com/report/github.com/kaptinlin/jsonpointer)

Fast implementation of [JSON Pointer (RFC 6901)][json-pointer] specification in Go.

[json-pointer]: https://tools.ietf.org/html/rfc6901

## Installation

```bash
go get github.com/kaptinlin/jsonpointer
```

## Usage

### Basic Operations

Find a value in a JSON object using a JSON Pointer string:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/kaptinlin/jsonpointer"
)

func main() {
    doc := map[string]any{
        "foo": map[string]any{
            "bar": 123,
        },
    }

    ref, err := jsonpointer.FindByPointer(doc, "/foo/bar")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(ref.Val) // 123
}
```

### Find by Path Components

Use variadic arguments to navigate to a value:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/kaptinlin/jsonpointer"
)

func main() {
    doc := map[string]any{
        "foo": map[string]any{
            "bar": 123,
        },
    }

    ref, err := jsonpointer.Find(doc, "foo", "bar")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Value: %v, Object: %v, Key: %v\n", ref.Val, ref.Obj, ref.Key)
    // Value: 123, Object: map[bar:123], Key: bar
}
```

### Safe Get Operations

Get values without error handling (returns nil if path doesn't exist):

```go
package main

import (
    "fmt"
    
    "github.com/kaptinlin/jsonpointer"
)

func main() {
    doc := map[string]any{
        "users": []any{
            map[string]any{"name": "Alice", "age": 30},
            map[string]any{"name": "Bob", "age": 25},
        },
    }

    // Get existing value using variadic arguments
    name := jsonpointer.Get(doc, "users", 0, "name")
    fmt.Println(name) // Alice

    // Get non-existing value
    missing := jsonpointer.Get(doc, "users", 5, "name")
    fmt.Println(missing) // <nil>
    
    // Get using JSON Pointer string
    age := jsonpointer.GetByPointer(doc, "/users/1/age")
    fmt.Println(age) // 25
}
```

### Path Manipulation

Convert between JSON Pointer strings and path arrays:

```go
package main

import (
    "fmt"
    
    "github.com/kaptinlin/jsonpointer"
)

func main() {
    // Parse JSON Pointer string to path array
    path := jsonpointer.Parse("/f~0o~1o/bar/1/baz")
    fmt.Printf("%+v\n", path)
    // [f~o/o bar 1 baz]

    // Format path components to JSON Pointer string
    pointer := jsonpointer.Format("f~o/o", "bar", "1", "baz")
    fmt.Println(pointer)
    // /f~0o~1o/bar/1/baz
    
    // Performance tip: For repeated access to the same path,
    // pre-parse the pointer once and reuse the path
    userNamePath := jsonpointer.Parse("/users/0/name")
    
    // Efficient repeated access
    for _, data := range datasets {
        name := jsonpointer.Get(data, userNamePath...)
        fmt.Println(name)
    }
}
```

### Component Encoding/Decoding

Encode and decode individual path components:

```go
package main

import (
    "fmt"
    
    "github.com/kaptinlin/jsonpointer"
)

func main() {
    // Unescape component
    unescaped := jsonpointer.Unescape("~0~1")
    fmt.Println(unescaped) // ~/

    // Escape component
    escaped := jsonpointer.Escape("~/")
    fmt.Println(escaped) // ~0~1
}
```

### Array Operations

Working with arrays and array indices:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/kaptinlin/jsonpointer"
)

func main() {
    doc := map[string]any{
        "items": []any{1, 2, 3},
    }

    // Access array element using variadic arguments
    ref, err := jsonpointer.Find(doc, "items", 1)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(ref.Val) // 2

    // Array end marker "-" points to next index
    ref, err = jsonpointer.Find(doc, "items", "-")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(ref.Key) // 3 (next available index)
    
    // Using JSON Pointer string
    value := jsonpointer.GetByPointer(doc, "/items/0")
    fmt.Println(value) // 1
}
```

### Struct Operations

Working with Go structs and JSON tags:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/kaptinlin/jsonpointer"
)

type User struct {
    Name    string `json:"name"`
    Age     int    `json:"age"`
    Email   string // No JSON tag, uses field name
    private string // Private field, ignored
    Ignored string `json:"-"` // Explicitly ignored
}

type Profile struct {
    User     *User  `json:"user"` // Pointer to struct
    Location string `json:"location"`
}

func main() {
    profile := Profile{
        User: &User{ // Pointer to struct
            Name:    "Alice",
            Age:     30,
            Email:   "alice@example.com",
            private: "secret",
            Ignored: "ignored",
        },
        Location: "New York",
    }

    // JSON tag access using variadic arguments
    name := jsonpointer.Get(profile, "user", "name")
    fmt.Println(name) // Alice

    // Field name access (no JSON tag)
    email := jsonpointer.Get(profile, "user", "Email")
    fmt.Println(email) // alice@example.com

    // Private fields are ignored
    private := jsonpointer.Get(profile, "user", "private")
    fmt.Println(private) // <nil>

    // json:"-" fields are ignored  
    ignored := jsonpointer.Get(profile, "user", "Ignored")
    fmt.Println(ignored) // <nil>

    // Nested struct navigation
    age := jsonpointer.Get(profile, "user", "age")
    fmt.Println(age) // 30

    // JSON Pointer syntax
    ref, err := jsonpointer.FindByPointer(profile, "/user/name")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(ref.Val) // Alice

    // Mixed struct and map data
    data := map[string]any{
        "profile": profile,
        "meta":    map[string]any{"version": "1.0"},
        "users":   []User{{Name: "Bob", Age: 25}},
    }
    
    // Access struct in map
    location := jsonpointer.Get(data, "profile", "location")
    fmt.Println(location) // New York
    
    // Access array of structs
    userName := jsonpointer.Get(data, "users", 0, "name")
    fmt.Println(userName) // Bob
}
```

### Validation

Validate JSON Pointer strings:

```go
package main

import (
    "fmt"
    
    "github.com/kaptinlin/jsonpointer"
)

func main() {
    // Valid JSON Pointer
    err := jsonpointer.Validate("/foo/bar")
    if err != nil {
        fmt.Printf("Invalid pointer: %v\n", err)
    } else {
        fmt.Println("Valid pointer")
    }

    // Invalid JSON Pointer
    err = jsonpointer.Validate("foo/bar") // missing leading slash
    if err != nil {
        fmt.Printf("Invalid pointer: %v\n", err)
    } else {
        fmt.Println("Valid pointer")
    }
}
```

### Performance Results

| Operation | Library | ns/op | B/op | allocs/op |
|-----------|---------|-------|------|-----------|
| Shallow (`/name`) | dolmen-go | 25.40 | 16 | 1 |
| | woodsbury | 28.32 | 16 | 1 |
| | BragdonD | 51.74 | 32 | 1 |
| | This impl (Find) | 72.77 | 96 | 4 |
| | This impl (Get) | 8.213 | 0 | 0 |
| | go-openapi | 111.4 | 124 | 5 |
| Deep (`/profile/settings/theme`) | dolmen-go | 55.02 | 16 | 1 |
| | woodsbury | 57.06 | 0 | 0 |
| | BragdonD | 121.6 | 64 | 1 |
| | This impl (Get) | 26.43 | 0 | 0 |
| | go-openapi | 133.4 | 192 | 5 |
| | This impl (Find) | 194.5 | 192 | 10 |
| Parse | BragdonD | 43.15 | 64 | 1 |
| | This impl | 78.95 | 96 | 4 |
| | go-openapi | 102.0 | 112 | 2 |

## Acknowledgments

This project is a Go port of the excellent [jsonjoy-com/json-pointer](https://github.com/jsonjoy-com/json-pointer) TypeScript implementation. We've adapted the core algorithms and added Go-specific performance optimizations while maintaining full RFC 6901 compatibility.

Special thanks to the original json-pointer project for providing a solid foundation and comprehensive test cases that enabled this high-quality Go implementation.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 
