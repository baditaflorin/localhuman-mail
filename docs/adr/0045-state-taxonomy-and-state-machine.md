# 0045 - State Taxonomy and State Machine

## Status

Accepted

## Context

Users should never land in a half-imported or unexplained state.

## Decision

Use the states documented in `docs/phase2-substance/states.md`. Recoverable failures keep existing mailbox state intact and include a next step.

## Consequences

- UI and API states are named consistently.
- Import failures do not erase useful state.

## Alternatives Considered

- Let each handler invent states: rejected because it leads to stuck states.

