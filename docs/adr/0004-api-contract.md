# 0004 - API Contract

## Status

Accepted

## Context

Mode C requires a runtime API. The frontend must not use hand-written, drifting assumptions about response payloads.

## Decision

Expose REST/JSON endpoints under `/api/v1` and document them in `api/openapi.yaml`.

Initial endpoints:

- `GET /healthz`
- `GET /readyz`
- `GET /metrics`
- `GET /api/v1/version`
- `GET /api/v1/capabilities`
- `GET /api/v1/messages`
- `GET /api/v1/messages/{id}`
- `GET /api/v1/search?q=...`
- `POST /api/v1/import/eml`
- `POST /api/v1/import/demo`
- `POST /api/v1/assist/reply`

The frontend validates payloads with Zod. Generated API clients can be added later with `openapi-typescript` and `openapi-fetch`.

## Consequences

- The frontend can point to local or deployed backends using build-time configuration.
- API drift is visible in `api/openapi.yaml`.
- The public Pages app can run in demo mode until the user connects a backend URL.

## Alternatives Considered

- GraphQL: rejected as unnecessary for v1.
- gRPC-Web: rejected because static Pages deployment and browser debugging are simpler with REST/JSON.
- WebSockets: deferred until real-time mailbox sync is a v1 requirement.

