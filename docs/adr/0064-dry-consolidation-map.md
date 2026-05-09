# 0064 - DRY Consolidation Map

## Status

Accepted

## Context

The frontend workbench mixed API calls, file handling, copy/download, state, and rendering. Error formatting and browser-output logic were likely to duplicate further during Phase 3.

## Decision

Extract authored frontend utilities for upload/import orchestration, copy/download/export, state snapshots, and share hash handling. Keep frontend and backend demo data separate because one supports offline Pages and one seeds the backend.

## Consequences

- `App.tsx` remains the composition surface, not the owner of every workflow detail.
- Boundary helpers become directly unit-testable.
- Accepted demo duplication is documented rather than abstracted into a build-time coupling.

## Alternatives Considered

- Shared generated demo JSON: rejected because it adds build plumbing for low value.
