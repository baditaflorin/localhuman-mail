# 0001 - Deployment Mode

## Status

Accepted

## Context

localhuman-mail needs a public frontend and private access to a user's mailbox. Browsers cannot connect directly to IMAP servers with raw TLS sockets, cannot safely keep mailbox credentials in a static bundle, and should not publish private mailbox artifacts to a public repository or GitHub Release.

## Decision

Use Mode C: GitHub Pages frontend plus a user-controlled Docker backend.

The frontend is static and served from GitHub Pages at https://baditaflorin.github.io/localhuman-mail/. The backend is an API-only Go service packaged as a Docker image. It owns IMAP fetch, mailbox parsing, local indexes, encryption hooks, local LLM calls, and private runtime data.

## Consequences

- GitHub Pages remains the public user interface.
- Mailbox credentials and mailbox contents never enter the frontend bundle or public repo.
- Docker deployment is required for real mailbox import and AI assist.
- CORS must allow the Pages origin and local development origins.
- The backend must expose health, readiness, metrics, and a documented API.

## Alternatives Considered

- Mode A: rejected because browser-only IMAP is not feasible with standard browser APIs.
- Mode B: rejected because private mailbox artifacts must not be committed or released publicly.
- Hosted SaaS: rejected because the project value proposition is local-first privacy.

