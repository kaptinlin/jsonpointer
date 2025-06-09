# Go JSON Pointer Performance Benchmarks

Performance comparison for Go JSON Pointer libraries.

## Libraries

| Library | Package |
|---------|---------|
| **this** | `github.com/kaptinlin/jsonpointer` |
| **go-openapi** | `github.com/go-openapi/jsonpointer` |
| **BragdonD** | `github.com/bragdond/jsonpointer-go` |
| **woodsbury** | `github.com/woodsbury/jsonpointer` |
| **dolmen-go** | `github.com/dolmen-go/jsonptr` |

## Environment

- **Platform**: Apple M3, macOS (darwin/arm64)
- **Test Method**: `go test -bench=. -benchmem`

## Results

### Basic Operations

| Operation | Library | ns/op | B/op | allocs/op |
|-----------|---------|-------|------|-----------|
| **Root (`""`)** | | | | |
| | this (Get) | 0.31 | 0 | 0 |
| | woodsbury | 1.11 | 0 | 0 |
| | BragdonD | 2.46 | 0 | 0 |
| | dolmen-go | 13.79 | 16 | 1 |
| | this (Find) | 16.60 | 48 | 1 |
| | go-openapi | 16.37 | 24 | 1 |
| **Shallow (`/name`)** | | | | |
| | this (Get) | 8.50 | 0 | 0 |
| | dolmen-go | 26.82 | 16 | 1 |
| | woodsbury | 29.23 | 16 | 1 |
| | BragdonD | 50.24 | 32 | 1 |
| | this (Find) | 63.97 | 80 | 3 |
| | go-openapi | 116.7 | 124 | 5 |
| **Deep (`/profile/settings/theme`)** | | | | |
| | this (Get) | 30.51 | 0 | 0 |
| | woodsbury | 50.87 | 0 | 0 |
| | dolmen-go | 64.33 | 16 | 1 |
| | BragdonD | 125.5 | 64 | 1 |
| | go-openapi | 140.7 | 192 | 5 |
| | this (Find) | 169.9 | 144 | 7 |

### Parser Performance

| Library | ns/op | B/op | allocs/op |
|---------|-------|------|-----------|
| this | 43.98 | 48 | 1 |
| BragdonD | 46.18 | 64 | 1 |
| go-openapi | 66.92 | 112 | 2 |

### Data Structure Access

| Scenario | this (Get) | this (Find) |
|----------|------------|-------------|
| **Struct** | 81.45 ns/op, 112 B/op | 66.71 ns/op, 160 B/op |
| **Map** | 8.20 ns/op, 0 B/op | 63.59 ns/op, 80 B/op |
| **Not Found** | 25.31 ns/op, 0 B/op | 48.57 ns/op, 64 B/op |

## API Comparison

This implementation provides two functions:
- **Get**: Direct value retrieval
- **Find**: Returns reference object with metadata

| Scenario | Get | Find |
|----------|-----|------|
| Root | 0.31 ns/op, 0 allocs | 16.60 ns/op, 1 alloc |
| Shallow | 8.50 ns/op, 0 allocs | 63.97 ns/op, 3 allocs |
| Deep | 30.51 ns/op, 0 allocs | 169.9 ns/op, 7 allocs |

## Memory Allocation Patterns

### Zero Allocation
- this (Get): All operations
- woodsbury: Root and deep access
- BragdonD: Root access only

### Single Allocation
- dolmen-go: 16B consistent
- BragdonD: Variable allocation
- this (Find): 48B minimum

### Multi Allocation
- go-openapi: 1-5 allocations
- this (Find): 1-7 allocations for complex operations

## Usage

```bash
cd benchmarks
go test -bench=. -benchmem
```

### Test Scenarios

- **Root**: Document root (`""`)
- **Shallow**: Top-level field (`"/name"`)
- **Deep**: Nested field (`"/profile/settings/theme"`)
- **Struct**: Go struct field access
- **Arrays**: Numeric index and end-marker access
- **Edge Cases**: Missing keys, invalid paths, escaped characters
