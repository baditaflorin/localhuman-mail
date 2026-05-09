# 0062 - Output Pathway Coverage Policy

## Status

Accepted

## Context

The app displayed useful results but did not let users take them out except by manual selection.

## Decision

Provide explicit copy controls for body, draft, and JSON; JSON and CSV downloads for the current list; a versioned state snapshot export/import; hash-based share links for small state snapshots; and print support for the selected message/draft. Screenshot and embed output are out of scope.

## Consequences

- Export formats become part of the tested product contract.
- Snapshot files carry a schema version and are validated on import.
- Share links must stay small and fail with an actionable message if too large.

## Alternatives Considered

- Backend export endpoints: rejected for Phase 3 because the frontend already has the current list and draft state.
- Hosted short-link service: rejected because it would introduce a public data surface.
