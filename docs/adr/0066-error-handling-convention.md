# 0066 - Error Handling Convention

## Status

Accepted

## Context

The backend now returns actionable import errors, but frontend toasts and batch flows need the same domain-language convention.

## Decision

User-facing errors use what/why/now-what. Batch import stores errors per file. Low-level exceptions are narrowed at the boundary and converted to this shape before reaching UI state.

## Consequences

- No raw `Error: undefined` messages in production UI.
- Batch import can partially succeed without hiding failures.
- Tests assert visible error text for invalid state/input.

## Alternatives Considered

- Toast-only fallback strings: rejected because they lose file and recovery context.
