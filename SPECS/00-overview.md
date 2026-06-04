# JSON Pointer Library Overview

## Overview

This library provides read-only JSON Pointer traversal and pointer utilities for
Go. The public API is centered on an immutable `Pointer` value: parse or build a
pointer, then use it to resolve values or references.

The package follows RFC 6901 for pointer string syntax and escaping while
supporting Go-native traversal over decoded JSON data and ordinary Go values.

## Product Goals

- Keep one obvious public model: `Pointer`.
- Parse pointer strings strictly so malformed syntax cannot become lookup data.
- Preserve raw-token construction for callers that already have token strings.
- Keep read APIs panic-free and error-returning.
- Preserve fast common paths for `map[string]any` and `[]any`.
- Stay small: pointer traversal, reference context, errors, and token helpers
  only.

## Scope

### In Scope

- Strict pointer parsing through `Parse`.
- Raw-token pointer construction through `FromTokens`.
- Value traversal through `Pointer.Value` and `Value`.
- Reference traversal through `Pointer.Reference` and `ReferenceOf`.
- Pointer formatting, token copying, parent and child derivation.
- Escape, unescape, index, and integer token helpers.
- Native Go traversal over maps, slices, arrays, structs, pointers, and
  interface values.

### Out of Scope

- Mutation APIs such as set, delete, append, patch, or merge helpers.
- Schema validation, JSON Patch, JSONPath, filters, wildcards, or query DSLs.
- Stateful compiled-pointer caches beyond the immutable `Pointer` value.
- Pluggable resolver interfaces or traversal policy systems.
- CLI diagnostics, repair suggestions, file IO, or product-specific envelopes.

> **Why**: the package is most useful as a tiny, predictable dependency for JSON
> Pointer reads and pointer utilities.
>
> **Rejected**: turning the library into a broader document-editing or querying
> toolkit.

## Compatibility Contract

- RFC 6901 controls pointer string syntax, escaping, and array token rules.
- Malformed pointer strings are rejected by `Parse`, `Value`, and
  `ReferenceOf`.
- Go-specific traversal extensions are allowed only when they do not change JSON
  Pointer meaning.
- Errors remain inspectable through `errors.Is`; traversal context is available
  through `errors.As` with `*Error`.

This API intentionally does not preserve the older mutable `Path` model or
permissive pointer-string traversal.

## Canonical Specifications

- `SPECS/10-domain-specs.md` defines pointer, token, traversal, and error
  semantics.
- `SPECS/20-api-specs.md` defines the exported surface.
- `SPECS/40-architecture-specs.md` defines package structure and performance
  strategy.
- `SPECS/50-coding-standards.md` defines contribution, documentation, and
  testing rules.

## Forbidden

- Do not add mutation APIs without a new spec that expands scope.
- Do not add permissive pointer parsing or implicit repair behavior.
- Do not add public traversal interfaces before real consumers need them.
- Do not duplicate traversal semantics across separate public paths.
- Do not trade the fast common read path for convenience-only behavior.

## Acceptance Criteria

- [ ] Public code centers on `Pointer`.
- [ ] Invalid pointer strings cannot produce usable pointers.
- [ ] Raw-token construction keeps literal token data distinct from pointer
      syntax.
- [ ] Scope remains limited to pointer traversal, reference context, errors, and
      token helpers.
- [ ] The canonical design guidance lives under `SPECS/`.
