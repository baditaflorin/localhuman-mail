# 0015 - Deployment Topology

## Status

Accepted

## Context

The frontend is hosted by GitHub Pages. The backend runs on a user-controlled server or workstation.

## Decision

Use GitHub Pages for the frontend and Docker Compose for the backend. Production Compose pulls `ghcr.io/baditaflorin/localhuman-mail:latest` and runs it behind nginx. nginx exposes public host port `25342`, terminates TLS, routes to internal `app:8080`, blocks public `/metrics`, and allows CORS from the Pages origin.

## Consequences

- Frontend and backend deploy independently.
- Backend can access local runtime storage and local AI services.
- Server setup is documented in `deploy/README.md`.

## Alternatives Considered

- Single binary on host: rejected because Docker provides repeatable deployment.
- Backend serving frontend: rejected by ADR 0010.
- Kubernetes: rejected as too heavy for v1.

