# JSON Pointer Coding Standards

## Overview

This spec defines the standards for changing code, docs, tests, and
configuration in this repository. The rules favor RFC 6901 correctness, readable
Go code, and lightweight maintenance.

## Compatibility Standards

- Keep RFC 6901 syntax, escaping, and array token semantics intact.
- Keep raw-token construction distinct from pointer-string parsing.
- Preserve `errors.Is` sentinel checks for documented error classes.
- Public behavior changes must update the relevant `SPECS/*.md` file and
  matching tests.

> **Why**: the library's value is stable pointer semantics with a small Go-native
> traversal layer.
>
> **Rejected**: convenience-only behavior that quietly accepts malformed pointer
> syntax.

## Documentation Standards

- `SPECS/` is the canonical home for design and contract documentation.
- Root `CLAUDE.md` stays short and operational.
- Comments and docs must be written in English.
- Exported names need godoc comments that start with the symbol name.
- Performance-sensitive functions should document complexity or hot-path intent
  when that context would otherwise be lost.

> **Why**: readers need to understand both the current rule and the compatibility
> reason behind it without hunting through old docs.
>
> **Rejected**: scattering design intent across comments, plans, and one-off
> markdown files.

## Testing and Validation Standards

- Run `task test` and `task lint` for code changes.
- Run `task specs-check` for markdown changes.
- Run `task yamllint` for YAML changes.
- Use benchmarks when changing hot traversal code.
- Keep documentation-structure checks in `lefthook.yml` or lint rules, not in
  `_test.go` files.

> **Why**: documentation layout is a repository invariant, not runtime behavior.
> Commit-time validation is the right enforcement point.
>
> **Rejected**: Go tests that read `SPECS/`, `README.md`, `CLAUDE.md`, or similar
> docs just to enforce layout.

## Change Workflow

- Work on `main`.
- Keep one logical change per commit.
- Use Conventional Commits.
- Fix lint and test failures before committing.
- Keep changes focused; delete dead code instead of leaving compatibility
  scaffolding behind.

## Forbidden

- Do not bypass RFC 6901 checks for behavior changes.
- Do not add docs-only layout tests under `*_test.go`.
- Do not leave design rules outside `SPECS/`.
- Do not preserve old public surface only as compatibility scaffolding.
- Do not commit with failing lint, test, or doc-validation results.

## Acceptance Criteria

- [ ] Behavior changes are checked against RFC 6901 and relevant Go-native
      traversal specs.
- [ ] The relevant spec is updated when public contracts change.
- [ ] Code, markdown, and YAML checks all pass for touched files.
- [ ] Documentation invariants stay enforced outside Go test binaries.
