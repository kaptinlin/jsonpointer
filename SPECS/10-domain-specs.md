# JSON Pointer Domain Specs

## Overview

This spec defines the domain rules for paths, pointers, traversal, and errors. It describes the semantics that the package exposes regardless of how the implementation is optimized internally.

## Path and Pointer Model

- `Path` is an ordered slice of string tokens.
- The root value is represented by an empty `Path` and an empty pointer string.
- `Parse` splits a pointer on `/` and unescapes `~0` to `~` and `~1` to `/`.
- `Format` joins path tokens with `/` and escapes only `~` and `/`.
- `IsRoot`, `Parent`, `IsChild`, `IsPathEqual`, `IsValidIndex`, and `IsInteger` operate on string tokens, not on a separate token type.
- `Parent` returns a detached parent path so mutating it does not mutate the input path.

> **Why**: a string-first public model matches RFC 6901 and keeps the API simple for callers that already manipulate path segments as strings.
>
> **Rejected**: exposing a second public token type with pre-parsed numeric indexes.

## Traversal Semantics

### Supported Containers

Read traversal supports:

- maps
- slices and arrays
- structs
- pointers to any supported container
- interface values that wrap any supported container

Empty path traversal returns the current value. Empty pointer traversal returns the root value.

### Struct and Map Access

- Struct lookup uses the JSON tag name when present.
- A JSON tag with an empty name falls back to the Go field name.
- Unexported fields and exact `json:"-"` fields are not addressable.
- A tag such as `json:"-,"` uses `-` as a literal field name, matching `encoding/json`.
- Reflective map traversal converts the string path token to the map key type when Go allows the conversion.

> **Why**: callers should be able to traverse ordinary Go values directly instead of pre-normalizing everything into generic JSON-shaped containers.
>
> **Rejected**: supporting only `map[string]any` and `[]any`.

### Pointer Handling

- Traversal dereferences pointers and unwraps non-nil interfaces until it reaches a concrete value.
- If a required pointer in that chain is nil, traversal returns `ErrNilPointer`.
- Nil interfaces behave like nil values and return `ErrNotFound` when traversal must continue.
- Nil is not silently treated as missing data once traversal has committed to dereferencing a pointer.

## Array Semantics

- Array indexes must be non-negative base-10 integers.
- Leading zeros are rejected except for `0` itself.
- `-` is the array end marker, but read APIs treat it as out of bounds.
- A non-numeric index returns `ErrInvalidIndex`.
- An index greater than the collection length returns `ErrIndexOutOfBounds`.
- A readable index must be strictly less than the collection length.

> **Why**: the library implements read traversal, not insertion. Preserving the RFC 6901 `-` token while rejecting reads through it keeps the semantics explicit.
>
> **Rejected**: treating `-` as a readable alias for the last element.

## Error Semantics

### Traversal Errors

- Missing map entry: `ErrKeyNotFound`
- Missing struct field: `ErrFieldNotFound`
- Nil pointer during dereference: `ErrNilPointer`
- Invalid array token: `ErrInvalidIndex`
- Array bounds failure: `ErrIndexOutOfBounds`
- Unsupported or non-traversable value: `ErrNotFound`

### Validation Errors

- Invalid pointer syntax: `ErrPointerInvalid`
- Pointer longer than `MaxPointerLength`: `ErrPointerTooLong`
- Path longer than `MaxPathLength`: `ErrPathTooLong`

`ValidatePath` currently enforces path length only. Because `Path` is already typed as `[]string`, current validation does not emit `ErrInvalidPath` or `ErrInvalidPathStep`.

> **Why**: the public `Path` type already constrains path step shape at compile time, so runtime validation focuses on length rather than re-checking every element.
>
> **Rejected**: reflection-heavy runtime validation of every `Path` entry.

## Forbidden

- Do not read through the `-` array marker.
- Do not expose unexported struct fields or `json:"-"` fields through traversal.
- Do not collapse key, field, index, and nil-pointer failures into one generic error.
- Do not document `ErrInvalidPath` or `ErrInvalidPathStep` as if current public APIs return them.

## Acceptance Criteria

- [ ] Root, path, and pointer semantics stay string-based and RFC 6901 compatible.
- [ ] Traversal supports native Go containers without marshal/unmarshal round-trips.
- [ ] Array read behavior rejects invalid or out-of-range indexes explicitly.
- [ ] Error docs match the current public behavior.

**Origin:** Migrated from `CLAUDE.md`.
