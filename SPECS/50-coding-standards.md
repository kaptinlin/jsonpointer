# JSON Pointer Coding Standards

## Overview

This spec defines the standards for changing code, docs, tests, and configuration in this repository. The rules favor TypeScript compatibility, readable Go code, and lightweight maintenance.

## Compatibility Standards

- Check the TypeScript reference implementation before changing public behavior.
- Keep RFC 6901 semantics intact unless the spec explicitly documents a Go-specific extension.
- Public behavior changes must update the relevant `SPECS/*.md` file and the matching tests.

> **Why**: the library's value is stable parity with the reference implementation, not novel semantics.
>
> **Rejected**: convenience-only Go changes that quietly drift from the documented contracts.

## Documentation Standards

- `SPECS/` is the canonical home for design and contract documentation.
- Root `CLAUDE.md` stays short and operational.
- Comments and docs must be written in English.
- Exported or behavior-heavy functions should cite the TypeScript source when the mapping is non-obvious.
- Performance-sensitive functions should document complexity or hot-path intent when that context would otherwise be lost.

> **Why**: readers need to understand both the current rule and the compatibility reason behind it without hunting through old docs.
>
> **Rejected**: scattering design intent across comments, plans, and one-off markdown files.
>
> **Status**: TypeScript-source references are not yet implemented consistently in `jsonpointer.go`.

## Testing and Validation Standards

- Run `task test` and `task lint` for code changes.
- Run `task specs-check` for markdown changes.
- Run `task yamllint` for YAML changes.
- Use benchmarks when changing hot traversal code.
- Keep documentation-structure checks in `lefthook.yml` or lint rules, not in `_test.go` files.

> **Why**: documentation layout is a repository invariant, not runtime behavior. Commit-time validation is the right enforcement point.
>
> **Rejected**: Go tests that read `SPECS/`, `README.md`, `CLAUDE.md`, or similar docs just to enforce layout.

## Change Workflow

- Work on `main`.
- Keep one logical change per commit.
- Use Conventional Commits.
- Fix lint and test failures before committing.
- Keep changes focused; delete dead code instead of leaving compatibility scaffolding behind.

## Forbidden

- Do not bypass TypeScript compatibility checks for behavior changes.
- Do not add docs-only layout tests under `*_test.go`.
- Do not leave design rules outside `SPECS/`.
- Do not commit with failing lint, test, or doc-validation results.

## Acceptance Criteria

- [ ] Behavior changes are checked against the reference implementation.
- [ ] The relevant spec is updated when public contracts change.
- [ ] Code, markdown, and YAML checks all pass for touched files.
- [ ] Documentation invariants stay enforced outside Go test binaries.

**Origin:** Migrated from `CLAUDE.md`.
