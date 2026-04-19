# CLAUDE.md

This repository contains a Go implementation of JSON Pointer (RFC 6901) with behavior anchored to the TypeScript `jsonjoy-com/json-pointer` project.

## Commands

```bash
task test
task lint
task verify
task markdownlint
task yamllint
task bench
```

## Working Rules

- Work on `main`.
- Keep traversal behavior aligned with the TypeScript reference implementation.
- Update `SPECS/` when changing domain rules, public API contracts, architecture, or coding standards.
- Keep comments and docs in English.
- Run `task lint` and `task test` after code changes.
- Run `task markdownlint` and `task yamllint` after docs or YAML changes.

## SPECS Index

- [SPECS/00-overview.md](SPECS/00-overview.md) — scope, goals, and compatibility boundaries
- [SPECS/10-domain-specs.md](SPECS/10-domain-specs.md) — pointer, path, traversal, and error semantics
- [SPECS/20-api-specs.md](SPECS/20-api-specs.md) — exported functions, types, errors, and validation limits
- [SPECS/40-architecture-specs.md](SPECS/40-architecture-specs.md) — package layout and performance strategy
- [SPECS/50-coding-standards.md](SPECS/50-coding-standards.md) — contribution, documentation, and validation rules
