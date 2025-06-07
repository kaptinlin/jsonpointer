# Go JSON Pointer Performance Benchmarks

Performance comparison tests for popular Go JSON Pointer libraries.

## Libraries Under Test

1. **Our Implementation** - `github.com/kaptinlin/jsonpointer`
2. **go-openapi** - `github.com/go-openapi/jsonpointer`
3. **BragdonD** - `github.com/bragdond/jsonpointer-go`
4. **woodsbury** - `github.com/woodsbury/jsonpointer`
5. **dolmen-go** - `github.com/dolmen-go/jsonptr`

## Usage

### Running Benchmarks

```bash
cd benchmarks
go test -bench=. -benchmem -count=3
```

### Running Specific Benchmarks

```bash
# Test only our implementation
go test -bench=BenchmarkOur -benchmem

# Test specific scenarios
go test -bench=Deep -benchmem
go test -bench=Parse -benchmem
```

### Benchmark Output

The benchmarks measure:
- **ns/op**: Nanoseconds per operation
- **B/op**: Bytes allocated per operation  
- **allocs/op**: Memory allocations per operation

### Test Scenarios

- **Root**: Access root path (`""`)
- **Shallow**: Access top-level field (`"/name"`)
- **Deep**: Access nested field (`"/profile/settings/theme"`)
- **Parse**: Parse JSON Pointer string to path
- **Medium**: Access data in larger dataset
- **NotFound**: Handle missing paths
- **Precompiled**: Reuse parsed paths