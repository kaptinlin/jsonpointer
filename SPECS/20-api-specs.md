# JSON Pointer API Specs

## Overview

This spec defines the exported API surface and the contracts callers can rely on. It covers public functions, exported types, validation limits, and compatibility notes that are visible outside the package.

## Traversal APIs

| API | Contract |
| --- | --- |
| `Get(doc any, path ...string) (any, error)` | Returns the value at `path`. An empty path returns `doc`. |
| `Find(doc any, path ...string) (*Reference, error)` | Returns the value plus its parent object and key. An empty path returns `&Reference{Val: doc}`. |
| `GetByPointer(doc any, pointer string) (any, error)` | Uses pointer parsing semantics and then traverses the result. An empty pointer returns `doc`. |
| `FindByPointer(doc any, pointer string) (*Reference, error)` | Traverses directly from the pointer string and returns a `Reference`. An empty pointer returns `&Reference{Val: doc}`. |

`GetByPointer` and `FindByPointer` do not call `Validate` automatically. Callers that need strict pointer syntax checks must call `Validate` explicitly.

> **Why**: keeping pointer validation explicit preserves the fast traversal path and avoids scanning the same pointer twice.
>
> **Rejected**: implicit validation inside every pointer-based read.

## Reference Types

### `Reference`

`Reference` is the generic traversal result:

- `Val`: the resolved value
- `Obj`: the parent container when a parent exists
- `Key`: the last traversed key or index, encoded as a string

### `ArrayReference[T]`

A typed array reference contains:

- `Val *T`
- `Obj []T`
- `Key int`

### `ObjectReference[T]`

A typed object reference contains:

- `Val T`
- `Obj map[string]T`
- `Key string`

### Type Predicates

- `IsArrayReference(ref Reference) bool`
- `IsObjectReference(ref Reference) bool`

The predicates inspect `Reference.Obj` and `Reference.Key`; `Find` and `FindByPointer` still return `Reference`, not the generic typed wrappers.

## Pointer and Path Utilities

| API | Contract |
| --- | --- |
| `Parse(pointer string) Path` | Parses a pointer into path tokens without returning an error. |
| `Format(path ...string) string` | Formats path tokens into a pointer string. |
| `Escape(component string) string` | Escapes `~` and `/` for a single path token. |
| `Unescape(component string) string` | Reverses `Escape` semantics for a single path token. |
| `IsChild(parent, child Path) bool` | Returns whether `child` has `parent` as a strict prefix. |
| `IsPathEqual(p1, p2 Path) bool` | Returns whether two paths are identical. |
| `IsRoot(path Path) bool` | Returns whether the path points at the root value. |
| `Parent(path Path) (Path, error)` | Returns the parent path or `ErrNoParent` for the root path. |
| `IsValidIndex(index string) bool` | Returns whether the token is a valid array index or `-`. |
| `IsInteger(str string) bool` | Returns whether the string is composed only of decimal digits. |

## Validation Surface

### Functions

- `Validate(pointer string) error`
- `ValidatePath(path Path) error`

### Constants

- `MaxPointerLength = 1024`
- `MaxPathLength = 256`

`Validate` accepts the empty string, requires a leading `/` for non-empty pointers, enforces the pointer length limit, and rejects invalid `~` escape sequences.

`ValidatePath` currently checks only the path length limit.

## Exported Errors

### TypeScript-Compatibility Errors

- `ErrInvalidIndex`
- `ErrNotFound`
- `ErrNoParent`
- `ErrPointerInvalid`
- `ErrPointerTooLong`
- `ErrInvalidPath`
- `ErrPathTooLong`
- `ErrInvalidPathStep`

### Go-Specific Errors

- `ErrIndexOutOfBounds`
- `ErrNilPointer`
- `ErrFieldNotFound`
- `ErrKeyNotFound`

Errors are part of the public compatibility surface and should remain stable.

> **Why**: callers often use `errors.Is` against these sentinels, so changing them is an API break even when function signatures stay the same.
>
> **Rejected**: hiding traversal failures behind formatted string errors.

## Forbidden

- Do not add implicit pointer validation to traversal APIs without updating this spec and benchmarks.
- Do not change exported error values casually.
- Do not return typed `ArrayReference[T]` or `ObjectReference[T]` from `Find` or `FindByPointer`.
- Do not document unimplemented behavior as if the current API guarantees it.

## Acceptance Criteria

- [ ] The exported function list here matches the package surface.
- [ ] Validation limits match `MaxPointerLength` and `MaxPathLength`.
- [ ] Compatibility notes cover the non-validating pointer traversal behavior.
- [ ] Error docs remain stable enough for `errors.Is` users.

**Origin:** Migrated from `CLAUDE.md`.
