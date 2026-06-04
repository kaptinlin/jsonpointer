# JSON Pointer API Specs

## Overview

This spec defines the exported API surface and the contracts callers can rely
on. The surface centers on `Pointer`, an immutable JSON Pointer value.

## Pointer Construction

| API | Contract |
| --- | --- |
| `Parse(pointer string) (Pointer, error)` | Parses strict RFC 6901 pointer syntax. Empty string returns root. |
| `FromTokens(tokens ...string) (Pointer, error)` | Builds a pointer from raw token strings and copies the input. |
| `Root() Pointer` | Returns the root pointer. |

`Parse` rejects malformed pointer syntax. `FromTokens` accepts literal token data
such as `"~2"` because raw tokens are not encoded pointer syntax.

## Pointer Methods

| API | Contract |
| --- | --- |
| `(Pointer).String() string` | Returns the canonical pointer string. |
| `(Pointer).Tokens() []string` | Returns a detached token copy. |
| `(Pointer).IsRoot() bool` | Reports whether the pointer targets the root. |
| `(Pointer).Parent() (Pointer, error)` | Returns the parent pointer or `ErrNoParent` for root. |
| `(Pointer).Child(tokens ...string) (Pointer, error)` | Returns a new pointer with tokens appended. |
| `(Pointer).Value(doc any) (any, error)` | Resolves the value at the pointer. |
| `(Pointer).Reference(doc any) (Reference, error)` | Resolves the value plus parent context. |

Pointer values are immutable from the caller's perspective. Methods that expose
or append tokens do not share mutable caller-owned slices.

## One-Shot Helpers

| API | Contract |
| --- | --- |
| `Value(doc any, pointer string) (any, error)` | Parses strictly, then resolves a value. |
| `ReferenceOf(doc any, pointer string) (Reference, error)` | Parses strictly, then resolves a reference. |

These helpers are conveniences over `Parse`; they are not a second traversal
model.

## Reference Type

`Reference` exposes named facts through methods:

- `Value() any`
- `Parent() (any, bool)`
- `Token() string`
- `Pointer() Pointer`

Root references have no parent and return `(nil, false)` from `Parent`.

## Token Utilities

| API | Contract |
| --- | --- |
| `EscapeToken(token string) string` | Escapes `~` and `/` for one raw token. |
| `UnescapeToken(encoded string) (string, error)` | Decodes `~0` and `~1`; rejects malformed escapes. |
| `IsValidIndex(token string) bool` | Returns whether the token is a valid array index or `-`. |
| `IsInteger(str string) bool` | Returns whether the string is composed only of decimal digits. |

## Validation Limits

### Constants

- `MaxPointerLength = 1024`
- `MaxPathLength = 256`

`MaxPointerLength` applies to pointer strings and canonical strings built from
raw tokens. `MaxPathLength` applies to token count.

## Exported Errors

### Pointer Construction Errors

- `ErrInvalidPointer`
- `ErrPointerTooLong`
- `ErrPathTooLong`

### Traversal Errors

- `ErrInvalidIndex`
- `ErrIndexOutOfBounds`
- `ErrKeyNotFound`
- `ErrFieldNotFound`
- `ErrNilPointer`
- `ErrNotFound`
- `ErrNoParent`

### Structured Error

`Error` wraps parse or traversal failures when pointer context is known:

- `Error() string`
- `Unwrap() error`
- `Pointer() Pointer`
- `Token() string`
- `Depth() int`

Callers should use `errors.Is` for class checks and `errors.As` with `*Error`
for pointer context. Error messages are for people and should not be parsed.

## Forbidden

- Do not add permissive parsing helpers that turn malformed pointer strings into
  lookup tokens.
- Do not expose a mutable `Path` as the central public API.
- Do not expose raw `Reference` fields.
- Do not return typed `ArrayReference[T]` or `ObjectReference[T]` shapes unless
  traversal APIs actually produce those values.
- Do not document unimplemented behavior as if the current API guarantees it.

## Acceptance Criteria

- [ ] The exported function list here matches the package surface.
- [ ] Validation limits match `MaxPointerLength` and `MaxPathLength`.
- [ ] `Parse` and one-shot helpers reject invalid pointer syntax.
- [ ] Error docs remain stable enough for `errors.Is` and `errors.As` users.
