# 0011 - Logging Strategy

## Status

Accepted

## Context

Backend logs need to be parseable in Docker while the frontend should avoid noisy production console output.

## Decision

Use Go `slog` JSON logs to stdout in the backend. Include method, path, status, duration, and request id for HTTP requests. Use minimal frontend console logging only in development builds.

## Consequences

- Docker log collectors can parse backend logs.
- Production frontend does not leak sensitive data through console output.
- Errors returned to users are clear but do not include secrets or raw stack traces.

## Alternatives Considered

- Text logs: rejected because JSON is easier for production pipelines.
- Zap or Zerolog: viable, but stdlib `slog` is sufficient for v1.

