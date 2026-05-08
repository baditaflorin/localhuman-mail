# Architecture

Live frontend: https://baditaflorin.github.io/localhuman-mail/

Repository: https://github.com/baditaflorin/localhuman-mail

## Context

```mermaid
C4Context
  title localhuman-mail context
  Person(user, "User", "Owns mailbox credentials and runtime data")
  System_Boundary(pages, "GitHub Pages") {
    System(frontend, "Static frontend", "React, TypeScript, Vite")
  }
  System_Boundary(host, "User-controlled host") {
    System(api, "Go backend", "REST API")
    SystemDb(sqlite, "SQLite", "Private messages and search data")
    System_Ext(tools, "Local tools", "readpst, file/libmagic, age, Tantivy CLI")
    System_Ext(ollama, "Ollama", "Optional local LLM")
  }
  Rel(user, frontend, "Opens", "HTTPS")
  Rel(frontend, api, "Calls", "REST/JSON")
  Rel(api, sqlite, "Stores and searches")
  Rel(api, tools, "Detects and orchestrates")
  Rel(api, ollama, "Drafts replies")
```

## Containers

```mermaid
C4Container
  title localhuman-mail containers
  Container(frontend, "GitHub Pages app", "React/Vite", "Static UI published from docs/")
  Container(nginx, "nginx", "nginx:alpine", "TLS, CORS, rate limits, port 25342")
  Container(api, "API", "Go/distroless", "Mailbox import, search, AI assist")
  ContainerDb(store, "Runtime volume", "SQLite", "Private local data")
  Container(prom, "Prometheus", "optional profile", "Scrapes /metrics internally")
  Rel(frontend, nginx, "REST/JSON")
  Rel(nginx, api, "Proxy")
  Rel(api, store, "Read/write")
  Rel(prom, api, "Scrape")
```

## Boundaries

- `frontend/` is public and static.
- `docs/` is the GitHub Pages publish directory.
- `cmd/server/` is the backend entry point.
- `internal/api/` exposes HTTP handlers.
- `internal/mailbox/` parses imported messages.
- `internal/search/` owns runtime storage and fallback search.
- `internal/ai/` owns local LLM integration.
- `deploy/` owns server deployment files.

