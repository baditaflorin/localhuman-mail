# 0016 - Local Git Hooks

## Status

Accepted

## Context

No GitHub Actions are allowed, so local hooks are the enforcement point for formatting, linting, tests, builds, smoke tests, secrets scanning, and commit message shape.

## Decision

Use a plain `.githooks/` directory wired by `make install-hooks`.

Hooks:

- `pre-commit`: formatting/lint/type checks and `gitleaks protect --staged`.
- `commit-msg`: Conventional Commits validation.
- `pre-push`: `make test`, `make build`, and `make smoke`.
- `post-merge` and `post-checkout`: run lightweight dependency/codegen refresh hooks.

## Consequences

- Hooks work without adding another package manager.
- Contributors must install hooks locally.
- Missing optional tools produce clear setup messages.

## Alternatives Considered

- Lefthook: viable, but plain shell hooks reduce dependencies.
- Husky: rejected because hooks need to cover Go and Docker workflows too.

