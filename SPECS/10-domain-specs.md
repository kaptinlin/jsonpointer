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
- structs
- pointers to supported values
- interface values that wrap supported values

Root traversal returns the current value. Traversal never mutates the document.

### Struct and Map Access

- Struct lookup uses JSON tag names when present.
- A JSON tag with an empty name falls back to the Go field name.
- Unexported fields and exact `json:"-"` fields are not addressable.
- A tag such as `json:"-,"` uses `-` as a literal field name.
- Embedded fields follow `encoding/json`-style dominance: shallower fields win,
  tagged fields win at the same depth, and ambiguous same-depth candidates are
  hidden.
- Reflective map traversal converts string tokens to the map key type when Go
  permits that conversion.

> **Why**: Go-native traversal is useful only when its projection rules are
> predictable. Struct support should feel like JSON field selection, not a local
> tag dialect.

### Pointer Handling

- Traversal dereferences pointers and unwraps non-nil interfaces until it reaches
  a concrete value.
- If a required pointer is nil, traversal returns `ErrNilPointer`.
- Nil interfaces behave like nil values and return `ErrNotFound` when traversal
  must continue.
- Nil is not silently treated as missing data once traversal has committed to
  dereferencing a pointer.

## Array Semantics

- Array indexes must be non-negative base-10 integers.
- Leading zeros are rejected except for `0` itself.
- `-` is the array end marker, but read APIs treat it as out of bounds.
- A non-numeric index returns `ErrInvalidIndex`.
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
- Pointer string longer than `MaxPointerLength`: `ErrPointerTooLong`
- Pointer containing more than `MaxPathLength` tokens: `ErrPathTooLong`

### Traversal Errors

- Missing map entry: `ErrKeyNotFound`
- Missing struct field: `ErrFieldNotFound`
- Nil pointer during dereference: `ErrNilPointer`
- Invalid array token: `ErrInvalidIndex`
- Array bounds failure: `ErrIndexOutOfBounds`
- Unsupported or non-traversable value: `ErrNotFound`

Traversal failures wrap the sentinel in `*Error` when pointer context is known.
`errors.Is` checks the class. `errors.As` exposes the requested pointer, failing
token, and zero-based token depth.

## Forbidden

- Do not read through the `-` array marker.
- Do not expose unexported struct fields or `json:"-"` fields through traversal.
- Do not collapse key, field, index, and nil-pointer failures into one generic
  error.
- Do not silently accept malformed pointer-string escapes.
- Do not reintroduce a public mutable path slice as the central model.

## Acceptance Criteria

- [ ] Root, token, and pointer semantics are RFC 6901 compatible.
- [ ] Invalid pointer strings fail before traversal.
- [ ] Raw-token construction supports literal strings that look like malformed
      escapes.
- [ ] Traversal supports native Go containers without marshal/unmarshal
      round-trips.
- [ ] Array read behavior rejects invalid or out-of-range indexes explicitly.
- [ ] Error docs match public behavior.
