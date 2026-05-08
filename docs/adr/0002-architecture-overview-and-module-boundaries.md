# 0002 - Architecture Overview and Module Boundaries

## Status

Accepted

## Context

The project combines a static web app, a private local backend, mailbox ingestion, search, and local AI integrations. Boundaries need to be explicit so mailbox contents do not cross into public artifacts.

## Decision

Split the system into these modules:

- `frontend/`: Vite React TypeScript app, built into `docs/` for GitHub Pages.
- `cmd/server/`: Go API server entry point.
- `internal/api/`: HTTP routing, handlers, middleware, OpenAPI-facing DTOs.
- `internal/config/`: environment configuration.
- `internal/mailbox/`: EML, Maildir, and PST import orchestration.
- `internal/search/`: local full-text search and index boundary.
- `internal/ai/`: local LLM and embedding adapters.
- `internal/security/`: encryption and secret-safe helpers.
- `internal/metrics/`: Prometheus metrics.
- `deploy/`: Docker Compose and nginx deployment files.

## Consequences

- Frontend code never reads mailbox files directly.
- Backend internals stay behind stable REST/JSON endpoints.
- Search and AI providers can be swapped without changing the UI.
- v1 can ship with safe local fallbacks while exposing clear extension points for Tantivy and sentence-transformers.

## Alternatives Considered

- Monolithic frontend-only app: rejected by ADR 0001.
- Backend serving the frontend: rejected because GitHub Pages is the canonical frontend host.
- Multiple backend services: rejected for v1 operational simplicity.

