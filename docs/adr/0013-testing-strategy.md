# 0013 - Testing Strategy

## Status

Accepted

## Context

Checks must run locally through Make and git hooks. There are no GitHub Actions.

## Decision

Use:

- Go unit tests colocated with backend source.
- Vitest unit tests colocated with frontend logic.
- Playwright smoke/e2e tests against a local static server.
- `scripts/smoke.sh` to build, serve, and verify a happy path.

## Consequences

- Contributors can run the same checks as hooks.
- The Pages build is validated before push.
- Backend integration tests can be added under `test/integration/` when external IMAP fixtures exist.

## Alternatives Considered

- Remote CI: rejected by project constraints.
- Manual-only testing: rejected because Pages publishing needs repeatable validation.

