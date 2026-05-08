# 0012 - Metrics and Observability

## Status

Accepted

## Context

Mode C requires scrape-ready Prometheus metrics while keeping user mailbox contents out of telemetry.

## Decision

Expose `/metrics` using `prometheus/client_golang`.

Metrics:

- HTTP request count by method, route, and status.
- HTTP request duration histogram.
- Imported message counter.
- Search request counter.
- AI assist request counter.

Do not include message subjects, senders, recipients, prompts, or query strings in metrics labels.

## Consequences

- Prometheus can scrape the backend.
- Metrics are useful without leaking mailbox content.
- nginx blocks public access to `/metrics` in production deployment.

## Alternatives Considered

- No metrics: rejected by Mode C requirements.
- Client analytics: rejected for v1; the Pages frontend has no analytics.
- High-cardinality labels: rejected to avoid PII leaks and Prometheus pressure.

