# Go JSON Pointer Performance Benchmarks

Performance comparison tests for Go JSON Pointer libraries.

## Libraries Under Test

1. **This Implementation** - `github.com/kaptinlin/jsonpointer`
2. **go-openapi** - `github.com/go-openapi/jsonpointer`
3. **BragdonD** - `github.com/bragdond/jsonpointer-go`
4. **woodsbury** - `github.com/woodsbury/jsonpointer`
5. **dolmen-go** - `github.com/dolmen-go/jsonptr`

## Benchmark Results (Apple M3)

### Core Performance Comparison

| Operation | Library | ns/op | B/op | allocs/op |
|-----------|---------|-------|------|-----------|
| **Root (`""`)** | This impl (Get) | 0.27 | 0 | 0 |
| | woodsbury | 1.08 | 0 | 0 |
| | BragdonD | 2.45 | 0 | 0 |
| | dolmen-go | 13.74 | 16 | 1 |
| | This impl (Find) | 14.73 | 48 | 1 |
| | go-openapi | 15.68 | 24 | 1 |
| **Shallow (`/name`)** | This impl (Get) | 8.09 | 0 | 0 |
| | dolmen-go | 29.36 | 16 | 1 |
| | woodsbury | 28.33 | 16 | 1 |
| | This impl (Find) | 62.18 | 80 | 3 |
| | BragdonD | 73.72 | 32 | 1 |
| | go-openapi | 112.5 | 124 | 5 |
| **Deep (`/profile/settings/theme`)** | This impl (Get) | 26.59 | 0 | 0 |
| | woodsbury | 57.04 | 0 | 0 |
| | dolmen-go | 57.60 | 16 | 1 |
| | BragdonD | 125.6 | 64 | 1 |
| | go-openapi | 134.7 | 192 | 5 |
| | This impl (Find) | 166.5 | 144 | 7 |
| **Parse** | This impl | 43.06 | 48 | 1 |
| | BragdonD | 42.09 | 64 | 1 |
| | go-openapi | 64.15 | 112 | 2 |

### Performance Analysis

#### Get vs Find Function Comparison

This implementation provides two distinct APIs:

- **Get function**: Returns values directly with zero memory allocation
- **Find function**: Returns Reference objects containing value, parent object, and key information

| Scenario | Get Performance | Find Performance | Difference |
|----------|----------------|------------------|------------|
| Root | 0.27 ns/op, 0 allocs | 14.73 ns/op, 1 alloc | 55x |
| Shallow | 8.09 ns/op, 0 allocs | 62.18 ns/op, 3 allocs | 7.7x |
| Deep | 26.59 ns/op, 0 allocs | 166.5 ns/op, 7 allocs | 6.3x |

#### Memory Allocation Patterns

- **Zero-allocation operations**: This implementation's Get function performs all operations without memory allocation
- **Allocation comparison**: Other libraries typically require 1-5 memory allocations per operation
- **Reference objects**: Find operations create Reference structures containing additional metadata

### Implementation Characteristics

| Library | Key Features |
|---------|-------------|
| This impl | Dual API (Get/Find), zero-allocation Get operations, string-only paths |
| go-openapi | Established library, JSON Schema integration |
| BragdonD | Simple interface, competitive performance |
| woodsbury | Good performance for deep access, zero allocation in some cases |
| dolmen-go | Consistent allocation pattern across operations |

## Usage

### Running Benchmarks

```bash
cd benchmarks
go test -bench=. -benchmem -count=3
```

### Test Scenarios

- **Root**: Access root path (`""`)
- **Shallow**: Access top-level field (`"/name"`)
- **Deep**: Access nested field (`"/profile/settings/theme"`)
- **Parse**: Parse JSON Pointer string to path components
- **Struct**: Direct struct field access with various patterns
- **Mixed**: Combined data type access patterns

### Complete Benchmark Output

```
BenchmarkOur_Get_Root-8          1000000000    0.27 ns/op     0 B/op    0 allocs/op
BenchmarkOur_Get_Shallow-8        144487514    8.09 ns/op     0 B/op    0 allocs/op  
BenchmarkOur_Get_Deep-8            44451303   26.59 ns/op     0 B/op    0 allocs/op
BenchmarkOur_Find_Root-8           76606346   14.73 ns/op    48 B/op    1 allocs/op
BenchmarkOur_Find_Shallow-8        18808678   62.18 ns/op    80 B/op    3 allocs/op
BenchmarkOur_Find_Deep-8            7079713  166.5 ns/op   144 B/op    7 allocs/op
BenchmarkOur_Parse-8               26890128   43.06 ns/op    48 B/op    1 allocs/op
```

## Technical Notes

- All benchmarks run on Apple M3 processor
- Results represent typical performance characteristics
- Memory allocation patterns vary based on operation complexity
- Performance may vary depending on JSON structure and access patterns
- Each library implements RFC 6901 JSON Pointer specification with different optimization strategies
