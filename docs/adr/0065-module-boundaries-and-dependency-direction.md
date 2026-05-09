# 0065 - Module Boundaries And Dependency Direction

## Status

Accepted

## Context

The current app is small, but Phase 3 adds enough workflow logic that boundaries matter.

## Decision

Frontend dependency direction is: `App` composition -> `features/mail/*` workflow utilities -> `api`/`lib` primitives. `features/mail` may depend on OpenAPI types and browser primitives; `api` and `lib` must not depend on React components.

Backend boundaries remain unchanged for Phase 3: `api` -> `mailbox`/`search`/`ai`; `mailbox` owns parsing/inference; `search` owns persistence.

## Consequences

- UI controls can call feature helpers but helpers do not render UI.
- Parser refactoring is deferred because Phase 3 explicitly avoids engine changes.

## Alternatives Considered

- Adopt a full frontend architecture framework: rejected as too heavy for this codebase.
