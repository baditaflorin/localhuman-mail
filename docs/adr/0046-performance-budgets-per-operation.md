# 0046 - Performance Budgets Per Operation

## Status

Accepted

## Context

Huge messages and nested MIME can make the app feel frozen.

## Decision

Phase 2 budgets:

- EML upload cap: 25MB.
- Indexed body cap: 2MB normalized text.
- Median fixture import: under 300ms.
- Worst fixture import: under 5s.
- Operations over 300ms must expose progress in the UI; operations over 5s must be cancellable in a later interaction pass.

## Consequences

- Parser code must avoid unbounded body growth.
- Performance fixtures track real regressions.

## Alternatives Considered

- Unlimited body indexing: rejected for memory and UI responsiveness.

