# JSON Pointer Library Overview

## Overview

This library provides read-only JSON Pointer traversal and pointer utilities for Go. It targets behavioral compatibility with the TypeScript `jsonjoy-com/json-pointer` implementation and RFC 6901 while supporting native Go values such as typed slices, arrays, structs, maps, and pointers.

## Product Goals

- Match the reference implementation for pointer syntax, escaping, and traversal semantics.
- Keep read APIs panic-free and error-returning.
- Preserve zero-allocation hot paths for common `map[string]any` and `[]any` workloads.
- Stay small: pointer traversal and validation only, not document mutation.

## Scope

### In Scope

- Read traversal through `Get`, `Find`, `GetByPointer`, and `FindByPointer`.
- Pointer parsing, formatting, escaping, unescaping, and validation helpers.
- Path utilities such as root checks, parent derivation, and index predicates.
- Native Go traversal over maps, slices, arrays, structs, and pointers to those values.

### Out of Scope

- Mutation APIs such as set, delete, append, patch, or merge helpers.
- Schema validation, JSON Patch, or JSONPath features.
- Stateful compiled-pointer objects or caching layers beyond internal field metadata.

> **Why**: the package is most useful as a tiny, predictable dependency for JSON Pointer reads and pointer utilities.
>
> **Rejected**: turning the library into a broader document-editing toolkit.

## Compatibility Contract

- The TypeScript repository is the behavioral baseline for pointer semantics and public API naming.
- RFC 6901 controls pointer string syntax, escaping, and array index rules.
- Go-specific extensions are allowed only when they make native Go data traversable without changing JSON Pointer meaning.
- When current Go behavior intentionally differs for performance or Go ergonomics, the difference must be documented in `SPECS/20-api-specs.md` or `SPECS/40-architecture-specs.md`.

> **Why**: parity is the library's core value, but Go callers still need to traverse ordinary Go values without marshaling everything into `map[string]any` first.
>
> **Rejected**: strict JSON-only traversal that ignores structs, typed slices, or non-string Go map keys.

## Canonical Specifications

- `SPECS/10-domain-specs.md` defines pointer, path, traversal, and error semantics.
- `SPECS/20-api-specs.md` defines the exported surface and compatibility notes.
- `SPECS/40-architecture-specs.md` defines package structure and performance strategy.
- `SPECS/50-coding-standards.md` defines contribution, documentation, and testing rules.

## Forbidden

- Do not add mutation APIs without a new spec that expands scope.
- Do not change traversal or error behavior without updating the relevant spec and tests.
- Do not duplicate design rules outside `SPECS/`; `CLAUDE.md` stays operational.
- Do not trade away TypeScript compatibility for convenience-only Go behavior.

## Acceptance Criteria

- [ ] Scope remains limited to pointer traversal, validation, and helper utilities.
- [ ] TypeScript compatibility remains the default decision rule for behavioral changes.
- [ ] Any Go-specific extension is documented in the domain or API specs.
- [ ] The canonical design guidance lives under `SPECS/`.

**Origin:** Migrated from `CLAUDE.md`.
