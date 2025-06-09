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
| | this (Get) | 0.55 | 0 | 0 |
| | woodsbury | 1.10 | 0 | 0 |
| | BragdonD | 2.48 | 0 | 0 |
| | dolmen-go | 13.85 | 16 | 1 |
| | this (Find) | 15.67 | 48 | 1 |
| | go-openapi | 15.80 | 24 | 1 |
| **Shallow (`/name`)** | | | | |
| | this (Get) | 9.28 | 0 | 0 |
| | dolmen-go | 27.44 | 16 | 1 |
| | woodsbury | 30.49 | 16 | 1 |
| | BragdonD | 50.66 | 32 | 1 |
| | this (Find) | 67.06 | 80 | 3 |
| | go-openapi | 114.7 | 124 | 5 |
| **Deep (`/profile/settings/theme`)** | | | | |
| | this (Get) | 29.54 | 0 | 0 |
| | woodsbury | 51.58 | 0 | 0 |
| | dolmen-go | 57.32 | 16 | 1 |
| | BragdonD | 128.8 | 64 | 1 |
| | go-openapi | 143.0 | 192 | 5 |
| | this (Find) | 191.0 | 144 | 7 |

### Parser Performance

| Library | ns/op | B/op | allocs/op |
|---------|-------|------|-----------|
| this | 45.11 | 48 | 1 |
| BragdonD | 46.71 | 64 | 1 |
| go-openapi | 73.95 | 112 | 2 |

### Data Structure Access

| Scenario | this (Get) | this (Find) |
|----------|------------|-------------|
| **Struct** | 78.43 ns/op, 112 B/op | 66.41 ns/op, 160 B/op |
| **Map** | 8.54 ns/op, 0 B/op | 65.77 ns/op, 80 B/op |
| **Not Found** | 48.85 ns/op, 64 B/op | 48.85 ns/op, 64 B/op |

## API Comparison

This implementation provides two functions:
- **Get**: Direct value retrieval
- **Find**: Returns reference object with metadata

| Scenario | Get | Find |
|----------|-----|------|
| Root | 0.55 ns/op, 0 allocs | 15.67 ns/op, 1 alloc |
| Shallow | 9.28 ns/op, 0 allocs | 67.06 ns/op, 3 allocs |
| Deep | 29.54 ns/op, 0 allocs | 191.0 ns/op, 7 allocs |

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
