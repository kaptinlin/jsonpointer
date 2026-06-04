# JSON Pointer Benchmarks

Focused benchmark coverage for the `Pointer`-centered API.

## Covered Paths

- strict `Parse`
- `FromTokens`
- canonical `Pointer.String`
- parsed `Pointer.Value`
- parsed `Pointer.Reference`
- one-shot `Value`
- traversal error wrapping
- typed struct fallback

## Usage

From this directory:

```bash
go test -bench=. -benchmem
```

From the repository root:

```bash
cd benchmarks
go test -bench=. -benchmem
```

## Reading Results

Use benchmark numbers as regression instruments, not product goals. The durable
performance requirement is that repeated reads through a parsed `Pointer` over
common decoded JSON values stay direct and allocation-conscious, while one-shot
helpers naturally include parse cost.
