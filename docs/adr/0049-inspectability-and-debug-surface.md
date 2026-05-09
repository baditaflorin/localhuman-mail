# 0049 - Inspectability and Debug Surface

## Status

Accepted

## Context

Power users and maintainers need to understand why the importer made a decision.

## Decision

Expose parser warnings, confidence reasons, shape, source hash, parser version, and schema version through the existing message API. The frontend reveals extra details when `?debug=1` is present.

## Consequences

- Support can diagnose fixture failures from exported/API data.
- Users can trust low-confidence warnings because reasoning is visible.

## Alternatives Considered

- Hidden logs only: rejected because users need visible confidence.

