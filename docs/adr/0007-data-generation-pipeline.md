# 0007 - Data Generation Pipeline

## Status

Accepted

## Context

Mode C does not use a public static data-generation pipeline. Private mailbox data must stay in the backend runtime store.

## Decision

Do not create a Mode B data pipeline in v1. Use backend import endpoints and local runtime storage instead.

## Consequences

- No mailbox artifacts are committed to the repo or attached to releases.
- `make data` is intentionally omitted because Mode C has no static data artifact contract.
- Future demo datasets must be synthetic and safe to commit.

## Alternatives Considered

- Commit anonymized mailbox data: rejected because anonymization risk is not worth the v1 value.
- Upload encrypted mailbox artifacts to Releases: rejected because frontend-held secrets are forbidden.
- Synthetic demo data: used through an API endpoint instead of a pipeline.

