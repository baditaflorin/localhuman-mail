# 0063 - Half-Baked Feature Triage Decisions

## Status

Accepted

## Context

Phase 3 requires every visible feature to be finished, hidden, or deleted.

## Decision

| Feature | Decision | Rationale |
| --- | --- | --- |
| EML import | Finish | Core workflow; add batch/drag/paste/progress. |
| Demo button | Finish/clarify | Keep backend seed and offline sample, label behavior honestly. |
| Draft textarea | Finish | Add persistence and copy/export. |
| Backend URL as setting | Finish | Add reset, clear UI state, and state export/import. |
| Capability detection | Keep and reword | Detection is useful but not a PST/Tantivy execution workflow. |
| URL/screenshot/folder import | Hide/delete from claims | Not shipped and not in scope. |

## Consequences

- No placeholder settings or controls ship.
- Documentation must use "detected" for optional native tools unless an execution path exists.

## Alternatives Considered

- Hide import until batch import exists: rejected because single import already works and remains useful.
