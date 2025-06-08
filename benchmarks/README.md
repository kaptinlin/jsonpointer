# Go JSON Pointer Performance Benchmarks

Performance comparison tests for Go JSON Pointer libraries.

## Libraries Under Test

1. **This Implementation** - `github.com/kaptinlin/jsonpointer`
2. **go-openapi** - `github.com/go-openapi/jsonpointer`
3. **BragdonD** - `github.com/bragdond/jsonpointer-go`
4. **woodsbury** - `github.com/woodsbury/jsonpointer`
5. **dolmen-go** - `github.com/dolmen-go/jsonptr`

## Benchmark Results (Apple M3)

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
- **Parse**: Parse JSON Pointer string to path
- **Struct**: Direct struct field access with caching
