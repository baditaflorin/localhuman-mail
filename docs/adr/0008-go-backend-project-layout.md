# 0008 - Go Backend Project Layout

## Status

Accepted

## Context

The backend needs a conventional Go layout that supports a single API server now and future local CLI helpers.

## Decision

Use the golang-standards-inspired layout:

- `cmd/server/`
- `internal/`
- `pkg/`
- `api/`
- `configs/`
- `scripts/`
- `test/`

One concern should live in one file when practical. Backend files should stay below 300 lines unless there is a clear reason.

## Consequences

- The runtime API has a clear entry point.
- Internal packages are protected from accidental public import contracts.
- Future generator or migration commands can be added under `cmd/`.

## Alternatives Considered

- Flat package layout: rejected because the project crosses multiple concerns.
- Hexagonal architecture framework: rejected as too heavy for v1.

