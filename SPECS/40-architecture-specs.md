# JSON Pointer Architecture Specs

## Overview

This spec defines how the package is organized and where performance-sensitive behavior belongs. The architecture optimizes common JSON-shaped data first and falls back to reflection only when required for broader Go compatibility.

## Package Layout

- `jsonpointer.go`: public API entry points
- `types.go`: exported types and type predicates
- `errors.go`: exported sentinel errors
- `get.go`: value traversal
- `find.go`: reference traversal with parent tracking
- `find_pointer.go`: direct pointer-string traversal
- `util.go`: escape, unescape, parse, format, and path helpers
- `validate.go`: validation limits and validation helpers
- `struct.go`: cached struct field lookup
- `*_test.go`, `fuzz_test.go`, and `benchmarks/`: verification and performance coverage

## Traversal Strategy

### `Get`

`Get` follows a layered approach:

1. `fastGet` handles common `map[string]any`, `[]any`, and pointer variants without token allocation.
2. `tryArrayAccess` and `tryObjectAccess` extend support to typed collections and structs.
3. Reflection is the final fallback for generic Go values.

### `Find`

`Find` keeps the last parent container and key while traversing so it can return a `Reference` without re-walking the document.

### `FindByPointer`

`FindByPointer` scans the pointer string directly instead of allocating a `Path` first. It preserves the same traversal semantics while avoiding an intermediate slice on the hot path.

> **Why**: most real workloads spend time in repeated reads over JSON-like maps and slices. Optimizing those cases first yields the best payoff without dropping generic Go traversal.
>
> **Rejected**: a reflection-only traversal implementation or a compiled-pointer object layer for every call site.

## Reflection Rules

- Pointer dereferencing is centralized through `derefValue`.
- Reflective map access converts the string token to the map key type when conversion is legal.
- Struct field lookup is centralized through `structField` and backed by a `sync.Map` cache.
- Struct caches store the resolved external field name, respecting JSON tags and excluding hidden fields.

## Performance Rules

- Hot-path improvements must avoid new allocations in the common `map[string]any` and `[]any` cases.
- Reflection should be a fallback, not the first choice.
- Significant traversal changes should be checked against the benchmark suite.
- Validation cost should remain explicit rather than silently folded into pointer traversal.

> **Why**: the package sells both correctness and speed. Hidden extra work in the hot path is a regression even when tests still pass.
>
> **Rejected**: convenience changes that shift steady-state read cost onto every caller.

## Forbidden

- Do not add a second traversal pipeline that duplicates existing logic without a measurable win.
- Do not move validation into hot traversal paths unless the API contract changes.
- Do not introduce caches beyond field metadata without clear invalidation-free semantics.
- Do not let architecture docs drift outside `SPECS/`.

## Acceptance Criteria

- [ ] Package responsibilities stay separated by concern.
- [ ] Common traversal paths remain optimized before reflection fallback.
- [ ] Struct metadata caching remains the only long-lived cache.
- [ ] Pointer-string traversal keeps its direct-scan architecture.

**Origin:** Migrated from `CLAUDE.md`.
