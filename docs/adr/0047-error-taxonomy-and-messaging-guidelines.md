# 0047 - Error Taxonomy and Messaging Guidelines

## Status

Accepted

## Context

Generic parse errors force users to guess what went wrong.

## Decision

Every import/search/draft error uses:

- What failed.
- Why, in mailbox terms.
- Now what, as a next step.

Errors are classified as `recoverable_input`, `recoverable_runtime`, or `fatal_runtime`.

## Consequences

- Users can retry with better input.
- Existing mailbox state remains safe on recoverable errors.

## Alternatives Considered

- Raw Go errors: rejected because they are not domain-aware.

