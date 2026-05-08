# 0017 - Dependency Policy

## Status

Accepted

## Context

This project touches private mailbox data. Dependencies should be boring, maintained, and justified.

## Decision

Use production-ready libraries only. Keep dependencies pinned by lockfiles. Run `go mod tidy`, `govulncheck` when available, `npm audit`, and local hooks before release.

Approved v1 choices include:

- Go: `chi`, `envconfig`, `go-playground/validator`, `prometheus/client_golang`, `emersion/go-message`, `modernc.org/sqlite`.
- Frontend: `vite`, `react`, `tailwindcss`, `zod`, `@tanstack/react-query`, `lucide-react`, `vitest`, `playwright`.

## Consequences

- The dependency graph stays understandable.
- Future replacements need an ADR when they change runtime behavior or security posture.
- Optional native tools are documented as capabilities rather than hidden requirements.

## Alternatives Considered

- Custom parsers/search engines: rejected where maintained libraries or local tools exist.
- Large UI frameworks: rejected to keep the first-load payload small.

