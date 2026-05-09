# 0068 - Persistence Schema And Migration Policy

## Status

Accepted

## Context

Phase 3 adds local UI state and state snapshots. These need to survive reloads and bad imports.

## Decision

Use schema version `localhuman.uiState.v1` for localStorage and exported snapshots. Validate external state files with zod. Unknown or old shapes migrate to safe defaults where possible; invalid external files are rejected with an actionable error and never partially applied.

## Consequences

- Export/import round-trip becomes deterministic for supported UI fields.
- Future version bumps must include a migration note.
- Tests cover valid, invalid, and minimal legacy shapes.

## Alternatives Considered

- Ad hoc JSON parsing: rejected because corrupted files would crash or silently drop data.
