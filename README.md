# localhuman-mail

![license](https://img.shields.io/badge/license-MIT-172033)
![pages](https://img.shields.io/badge/GitHub%20Pages-live-276749)
![mode](https://img.shields.io/badge/deployment-Mode%20C-b45309)

Live app: https://baditaflorin.github.io/localhuman-mail/

Repository: https://github.com/baditaflorin/localhuman-mail

Support development: https://www.paypal.com/paypalme/florinbadita

localhuman-mail is a privacy-first email workbench that gives you fast search, local mailbox import, and local-AI reply drafting without handing a hosted service your inbox.

![localhuman-mail screenshot](docs/screenshot.svg)

## Quickstart

```bash
git clone https://github.com/baditaflorin/localhuman-mail.git
cd localhuman-mail
npm --prefix frontend install
make install-hooks
make dev
```

The frontend runs with Vite in development and builds into `docs/` for GitHub Pages. The backend listens on `http://localhost:8080` by default.

## What Works In V1

- GitHub Pages frontend with repo, PayPal, version, and commit metadata visible in the UI.
- Go backend with `/healthz`, `/readyz`, `/metrics`, `/api/v1/version`, mailbox import, search, and reply assist.
- EML import, synthetic demo import, SQLite-backed local message storage, and local search.
- Local LLM reply assist through Ollama, with a local template fallback.
- Capability detection for `readpst`, libmagic via `file`, Tantivy CLI, sentence-transformers, `age`, SQLite, and Ollama.
- Docker Compose production topology with nginx on public host port `25342`.

## Architecture

```mermaid
C4Context
  title localhuman-mail context
  Person(user, "User", "Runs the backend and opens the Pages UI")
  System_Boundary(pages, "GitHub Pages") {
    System(frontend, "React frontend", "Static app from docs/")
  }
  System_Boundary(local, "User-controlled host") {
    System(api, "Go backend", "Mailbox import, search, AI assist")
    SystemDb(store, "SQLite runtime store", "Private mailbox index")
    System_Ext(ollama, "Ollama", "Optional local LLM")
  }
  System_Ext(repo, "GitHub repository", "Source, stars, Pages")
  Rel(user, frontend, "Uses", "HTTPS")
  Rel(frontend, api, "Calls", "REST/JSON")
  Rel(api, store, "Reads/writes", "local file")
  Rel(api, ollama, "Generates drafts", "local HTTP")
  Rel(frontend, repo, "Links to", "HTTPS")
```

More detail: docs/architecture.md

ADRs: docs/adr/

Deployment: deploy/README.md

API: docs/api.md and api/openapi.yaml

## Commands

```bash
make help
make build
make test
make smoke
make docker-build
```

## Configuration

Copy `.env.example` for local backend values. Never commit real `.env` files, mailbox data, private keys, or credentials.

Frontend build-time variables are public. Do not put secrets in `VITE_*` variables.

## GitHub Pages Publishing

Pages is configured from `main` branch `/docs` folder.

Live URL: https://baditaflorin.github.io/localhuman-mail/

`make build` regenerates the frontend into `docs/`, including hashed assets and `404.html`.

## Security

No mailbox contents, credentials, tokens, or private keys belong in git. The frontend stores only non-sensitive UI settings. Mailbox contents stay in the backend runtime volume.

Security policy: SECURITY.md

