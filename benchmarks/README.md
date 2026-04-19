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
| | this (Get) | 0.41 | 0 | 0 |
| | woodsbury | 2.03 | 0 | 0 |
| | BragdonD | 4.23 | 0 | 0 |
| | this (Find) | 29.18 | 48 | 1 |
| | go-openapi | 30.38 | 24 | 1 |
| | dolmen-go | 125.4 | 16 | 1 |
| **Shallow (`/name`)** | | | | |
| | this (Get) | 13.96 | 0 | 0 |
| | woodsbury | 47.74 | 16 | 1 |
| | BragdonD | 72.60 | 32 | 1 |
| | dolmen-go | 88.89 | 16 | 1 |
| | this (Find) | 127.9 | 80 | 3 |
| | go-openapi | 347.3 | 212 | 10 |
| **Deep (`/profile/settings/theme`)** | | | | |
| | this (Get) | 62.98 | 0 | 0 |
| | woodsbury | 92.48 | 0 | 0 |
| | dolmen-go | 156.0 | 16 | 1 |
| | BragdonD | 218.5 | 64 | 1 |
| | this (Find) | 280.0 | 144 | 7 |
| | go-openapi | 379.9 | 288 | 10 |

### Parser Performance

| Library | ns/op | B/op | allocs/op |
|---------|-------|------|-----------|
| BragdonD | 69.84 | 64 | 1 |
| go-openapi | 117.9 | 112 | 2 |
| this | 176.1 | 48 | 1 |

### Data Structure Access

| Scenario | this (Get) | this (Find) |
|----------|------------|-------------|
| **Struct** | 91.52 ns/op, 112 B/op | 149.7 ns/op, 160 B/op |
| **Map** | 16.87 ns/op, 0 B/op | 138.9 ns/op, 80 B/op |
| **Not Found** | 172.6 ns/op, 16 B/op | 172.6 ns/op, 16 B/op |

## API Comparison

This implementation provides two functions:

- **Get**: Direct value retrieval (zero allocations)
- **Find**: Returns reference object with metadata

| Scenario | Get | Find |
|----------|-----|------|
| Root | 0.41 ns/op, 0 allocs | 29.18 ns/op, 1 alloc |
| Shallow | 13.96 ns/op, 0 allocs | 127.9 ns/op, 3 allocs |
| Deep | 62.98 ns/op, 0 allocs | 280.0 ns/op, 7 allocs |

## Memory Allocation Patterns

### Zero Allocation

- **this (Get)**: All operations (0 allocs)
- **woodsbury**: Root and deep access (0 allocs)
- **BragdonD**: Root access only (0 allocs)

### Single Allocation

- **dolmen-go**: 16B consistent (1 alloc)
- **BragdonD**: 32-64B variable allocation (1 alloc)
- **this (Find)**: 48B minimum (1 alloc for root)

### Multi Allocation

- **go-openapi**: 1-10 allocations, 24-288B
- **this (Find)**: 1-7 allocations for complex operations (48-144B)

## Key Performance Insights

### Generic Types Outperform Specialized Types

Benchmark testing reveals that **generic types (`[]any`, `map[string]any`) are significantly faster** than specialized types (`[]string`, `[]int`, `map[string]int`, etc.):

| Type Comparison | Generic | Specialized | Performance Gain |
|-----------------|---------|-------------|------------------|
| `[]string` access | 72.16 ns/op, 24B | 74.73 ns/op, 40B | **~1% faster, 40% less memory** |
| `[]int` access | 88.48 ns/op, 24B | 98.90 ns/op, 32B | **11% faster, 25% less memory** |
| `[]float64` access | 114.2 ns/op, 24B | 133.5 ns/op, 32B | **14% faster, 25% less memory** |
| `map[string]string` | 57.98 ns/op, 0B | 349.6 ns/op, 32B | **503% faster, zero allocs** |
| `map[string]int` | 69.23 ns/op, 0B | 255.6 ns/op, 24B | **269% faster, zero allocs** |
| `map[string]float64` | 87.56 ns/op, 0B | 357.3 ns/op, 24B | **308% faster, zero allocs** |
| Nested structures | 104.3 ns/op, 0B | 1035 ns/op, 56B | **893% faster, zero allocs** |

**Why generic types are faster:**

- Generic types (`map[string]any`, `[]any`) hit the **ultra-fast path** with direct type assertions
- Specialized types require additional reflection overhead for type conversion
- Generic types dominate real-world JSON parsing (99%+ of use cases)

**Recommendation:** Use standard `json.Unmarshal` which produces generic types for optimal performance.

## Benchmark Files

The benchmark suite is organized into focused test files:

- **`compare_test.go`** - Cross-library performance comparison (46 benchmarks)
  - Compares 5 implementations: this, go-openapi, BragdonD, woodsbury, dolmen-go
  - Tests Get, Find, Parse operations across all libraries
  - Includes struct vs map access patterns

- **`get_test.go`** - Core Get/Find operations (5 benchmarks)
  - Find, FindByPointer, Get functions
  - Type guard functions (IsArrayReference, IsObjectReference)
  - Comparative analysis between Find and FindByPointer

- **`parse_test.go`** - Parser and formatter operations (6 benchmarks)
  - Parse, Format, Escape, Unescape functions
  - Roundtrip validation tests
  - Special character handling

- **`types_test.go`** - Type-specific performance (17 benchmarks)
  - Specialized types: `[]string`, `[]int`, `[]float64`, `map[string]string`, etc.
  - Generic types: `[]any`, `map[string]any`
  - Performance comparison showing generic types are 74%-264% faster
  - Nested structure access patterns

## Usage

```bash
cd benchmarks
go test -bench=. -benchmem
```

### Run Specific Benchmark Files

```bash
# Library comparison
go test -bench=. -benchmem compare_test.go

# Get/Find operations
go test -bench=. -benchmem get_test.go

# Parser operations
go test -bench=. -benchmem parse_test.go

# Type-specific performance
go test -bench=. -benchmem types_test.go
```

### Test Scenarios

- **Root**: Document root (`""`)
- **Shallow**: Top-level field (`"/name"`)
- **Deep**: Nested field (`"/profile/settings/theme"`)
- **Struct**: Go struct field access
- **Arrays**: Numeric index and end-marker access
- **Edge Cases**: Missing keys, invalid paths, escaped characters
- **Type Variants**: Specialized vs generic type performance
