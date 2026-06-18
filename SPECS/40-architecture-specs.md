# JSON Pointer Architecture Specs

## Overview

This spec defines how the package is organized and where performance-sensitive
behavior belongs. The architecture optimizes common JSON-shaped data first and
falls back to reflection only for typed container mechanics.

## Package Layout

- `jsonpointer.go`: public parsing, one-shot traversal, and token helpers
- `types.go`: `Pointer`, `Reference`, and their methods
- `errors.go`: exported sentinel errors and structured `Error`
- `get.go`: shared value and reference resolver
- `util.go`: parse, format, escape, unescape, map, index, and reflection helpers
- `validate.go`: pointer-string syntax validation helpers
- `*_test.go`, `fuzz_test.go`, `examples/`, and `benchmarks/`: verification and
  performance coverage

## Traversal Strategy

`Pointer.Value` and `Pointer.Reference` share the same internal stepping
behavior. Each token is resolved once through the same container rules:

1. Common `map[string]any`, `*map[string]any`, `[]any`, and `*[]any` values are
   handled directly.
2. Pointer and interface values are dereferenced consistently.
3. Typed slices, arrays, and maps use reflective fallback.
4. Structs and other non-container values are not traversable.

`Reference` keeps the last parent container and token while using the same step
function as value traversal.

> **Why**: value lookup and reference lookup should not drift. Fast paths are
> implementation details inside one semantic path.
>
> **Rejected**: a separate direct pointer-string traversal pipeline.

## Pointer Construction Strategy

- `Parse` validates pointer-string syntax before token construction.
- `FromTokens` copies raw tokens without interpreting pointer-string syntax.
- `Pointer.String` formats from raw tokens.
- `Pointer.Tokens` returns a copy.

Strict pointer parsing is a boundary cost. Reusing parsed pointers avoids paying
that cost on repeated reads.

## Reflection Rules

- Pointer dereferencing and interface unwrapping are centralized through
  `derefValue`.
- Reflective map access converts the string token to the map key type when
  conversion is legal.
- Structs intentionally fall through as non-traversable values.

## Performance Rules

- Hot-path improvements must avoid new allocations in successful common
  `map[string]any` and `[]any` traversal.
- Reflection should be a container fallback, not the first choice.
- Reference context and structured errors should cost only when requested or
  when an error occurs.
- Significant traversal changes should be checked against the benchmark suite.
- Benchmark numbers are instruments, not product goals.

> **Why**: the package sells both correctness and speed. Hidden extra work in the
> hot path is a regression even when tests still pass.
>
> **Rejected**: convenience changes that shift steady-state read cost onto every
> caller.

## Forbidden

- Do not add a second traversal pipeline that duplicates resolver semantics.
- Do not parse pointer strings inside `Pointer.Value` or `Pointer.Reference`.
- Do not introduce public resolver interfaces or plugin hooks without real
  consumers.
- Do not introduce caches beyond immutable pointer values without clear
  invalidation-free semantics.
- Do not let architecture docs drift outside `SPECS/`.

## Acceptance Criteria

- [ ] Package responsibilities stay separated by concern.
- [ ] `Value` and `Reference` share traversal semantics.
- [ ] Common traversal paths remain optimized before typed-container fallback.
- [ ] No struct metadata cache or field-selection layer exists.
- [ ] Pointer-string parsing stays at construction and one-shot helper
      boundaries.
