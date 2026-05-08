# 0009 - Configuration and Secrets Management

## Status

Accepted

## Context

The backend needs runtime configuration for ports, CORS origins, storage paths, and local AI endpoints. Secrets must never be committed.

## Decision

Use environment variables loaded by `envconfig`, document placeholders in `.env.example`, and keep real `.env` files gitignored.

## Consequences

- Docker Compose can supply runtime configuration without baking secrets into images.
- Frontend build variables are public and must never contain secrets.
- gitleaks runs in the pre-commit hook.

## Alternatives Considered

- Config files with secrets: rejected because they are easier to commit accidentally.
- Viper-only config: rejected because env vars are sufficient for v1.
- Frontend-stored encrypted secrets: rejected because encrypted secrets in the frontend still violate the project constraints.

