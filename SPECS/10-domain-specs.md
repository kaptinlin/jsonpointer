# JSON Pointer Domain Specs

## Overview

This spec defines the domain rules for pointers, tokens, traversal, and errors.
It describes the semantics exposed by the package regardless of private
optimization strategy.

## Pointer and Token Model

- `Pointer` is an immutable ordered list of raw string tokens.
- The root value is represented by `Root()` and the empty pointer string.
- `Parse` accepts the empty string or strings beginning with `/`.
- `Parse` unescapes `~0` to `~` and `~1` to `/`.
- `Parse` rejects malformed escape sequences such as `~`, `~2`, and `~x`.
- `FromTokens` accepts raw token strings and does not interpret `~2` as syntax.
- `Pointer.String` returns the canonical RFC 6901 string form.
- `Pointer.Tokens` returns a detached copy.
- `Pointer.Parent` returns `ErrNoParent` for root.
- `Pointer.Child` returns a new pointer and does not mutate the receiver.

> **Why**: pointer-string syntax and raw token data are different concepts.
> Keeping them separate prevents malformed syntax from becoming a successful
> lookup.
>
> **Rejected**: a public mutable `Path` slice as the central pointer model.

## Traversal Semantics

### Supported Containers

Read traversal supports:

- maps
- slices and arrays
- string-keyed typed maps
- pointers to supported values
- interface values that wrap supported values

Root traversal returns the current value. Traversal never mutates the document.

Struct values are not field-selected. Callers that want to use JSON Pointer over
struct fields should first project those structs into a JSON document shape.

### Map Access

Direct `map[string]any` traversal uses the token as the map key. Reflective map
traversal converts string tokens to the map key type when Go permits that
conversion.

> **Why**: map, slice, array, pointer, and interface traversal are container
> mechanics. Struct field selection is a naming and projection contract that
> belongs outside a small JSON Pointer kernel.
>
> **Rejected**: cloning `encoding/json` field-selection rules, name providers,
> or custom lookup interfaces.

### Pointer Handling

- Traversal dereferences pointers and unwraps non-nil interfaces until it reaches
  a concrete value.
- If a required pointer is nil, traversal returns `ErrNilPointer`.
- Nil interfaces behave like nil values and return `ErrNotTraversable` when
  traversal must continue.
- Pointer and interface wrapper cycles return `ErrNotTraversable` instead of
  looping or relying on an arbitrary dereference limit.
- Root traversal consumes no token, so it returns the original document without
  dereferencing it, including a self-referential wrapper.
- Nil is not silently treated as missing data once traversal has committed to
  dereferencing a pointer.

### Reference Context

For non-root references, `Reference.Parent` returns the concrete container that
consumed the final token after pointer dereferencing and interface unwrapping.
Root references have no parent.

> **Why**: parent context should describe the JSON-shaped traversal point, not an
> incidental Go wrapper used to reach it.

## Array Semantics

- Array indexes must be non-negative base-10 integers.
- Leading zeros are rejected except for `0` itself.
- `-` refers to the nonexistent member after the last array element, so read
  APIs treat it as out of bounds.
- A non-numeric index returns `ErrInvalidIndex`.
- A canonical decimal index may have any length; if it cannot identify an
  existing element, traversal returns `ErrIndexOutOfBounds`.
- An index greater than or equal to the collection length returns
  `ErrIndexOutOfBounds`.
- A readable index must be strictly less than the collection length.

> **Why**: the library implements read traversal, not insertion. Preserving the
> RFC 6901 `-` token while rejecting reads through it keeps the semantics
> explicit.
>
> **Rejected**: treating `-` as a readable alias for the last element.

## Error Semantics

### Pointer Construction Errors

- Invalid pointer syntax: `ErrInvalidPointer`

### Traversal Errors

- Missing map entry: `ErrKeyNotFound`
- Nil pointer during dereference: `ErrNilPointer`
- Invalid array token: `ErrInvalidIndex`
- Array bounds failure: `ErrIndexOutOfBounds`
- Unsupported or non-traversable value: `ErrNotTraversable`

Traversal failures wrap the sentinel in `*Error` when pointer context is known.
`errors.Is` checks the class. `errors.As` exposes the requested pointer, failing
token, and zero-based token depth.

## Forbidden

- Do not read through the `-` array marker.
- Do not add struct field selection to traversal.
- Do not collapse key, index, and nil-pointer failures into one generic
  error.
- Do not silently accept malformed pointer-string escapes.
- Do not reintroduce a public mutable path slice as the central model.

## Acceptance Criteria

- [ ] Root, token, and pointer semantics are RFC 6901 compatible.
- [ ] Invalid pointer strings fail before traversal.
- [ ] Raw-token construction supports literal strings that look like malformed
      escapes.
- [ ] Traversal supports JSON-shaped Go containers without marshal/unmarshal
      round-trips.
- [ ] Pointer and interface wrapper chains terminate with either a concrete
      container or a structured traversal error.
- [ ] Array read behavior rejects invalid or out-of-range indexes explicitly.
- [ ] Error docs match public behavior.
