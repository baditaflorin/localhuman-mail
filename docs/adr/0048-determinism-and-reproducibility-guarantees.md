# 0048 - Determinism and Reproducibility Guarantees

## Status

Accepted

## Context

Same input must yield same normalized output and stable IDs.

## Decision

Use source SHA-256, deterministic missing-date fallback, sorted metadata where applicable, and no import-time wall clock in normalized message output. Runtime storage may have created-at timestamps, but API normalized fields stay deterministic.

## Consequences

- Fixture tests can compare canonical JSON.
- Duplicate import is stable and dedupes.

## Alternatives Considered

- Current-time fallback for missing dates: rejected as silent wrongness.

