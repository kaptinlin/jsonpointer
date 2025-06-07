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

    ref, err := jsonpointer.FindByPointer("/foo/bar", doc)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(ref.Val) // 123
}
```

### Find by Path Array

Use an array of steps to navigate to a value:

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

    path := jsonpointer.ParseJsonPointer("/foo/bar")
    ref, err := jsonpointer.Find(doc, path)
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

    // Get existing value
    name := jsonpointer.Get(doc, jsonpointer.Path{"users", 0, "name"})
    fmt.Println(name) // Alice

    // Get non-existing value
    missing := jsonpointer.Get(doc, jsonpointer.Path{"users", 5, "name"})
    fmt.Println(missing) // <nil>
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
    // Parse JSON Pointer string to path
    path := jsonpointer.ParseJsonPointer("/f~0o~1o/bar/1/baz")
    fmt.Printf("%+v\n", path)
    // [f~o/o bar 1 baz]

    // Format path array to JSON Pointer string
    pointer := jsonpointer.FormatJsonPointer(jsonpointer.Path{"f~o/o", "bar", "1", "baz"})
    fmt.Println(pointer)
    // /f~0o~1o/bar/1/baz
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
    unescaped := jsonpointer.UnescapeComponent("~0~1")
    fmt.Println(unescaped) // ~/

    // Escape component
    escaped := jsonpointer.EscapeComponent("~/")
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

    // Access array element
    ref, err := jsonpointer.Find(doc, jsonpointer.Path{"items", 1})
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(ref.Val) // 2

    // Array end marker "-" points to next index
    ref, err = jsonpointer.Find(doc, jsonpointer.Path{"items", "-"})
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(ref.Key) // 3 (next available index)
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
    err := jsonpointer.ValidateJsonPointer("/foo/bar")
    fmt.Println(err) // <nil>

    // Invalid JSON Pointer
    err = jsonpointer.ValidateJsonPointer("foo/bar") // missing leading slash
    fmt.Println(err) // error message
}
```

## Acknowledgments

This project is a Go port of the excellent [jsonjoy-com/json-pointer](https://github.com/jsonjoy-com/json-pointer) TypeScript implementation. We've adapted the core algorithms and added Go-specific performance optimizations while maintaining full RFC 6901 compatibility.

Special thanks to the original json-pointer project for providing a solid foundation and comprehensive test cases that enabled this high-quality Go implementation.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 
