# Postmortem

## What Was Built

- Public GitHub repository: https://github.com/baditaflorin/localhuman-mail
- Live GitHub Pages frontend: https://baditaflorin.github.io/localhuman-mail/
- Mode C Go backend with health, readiness, metrics, mailbox import, search, capabilities, and reply assist.
- Docker Compose deployment with nginx on host port `25342`.
- Local hooks, Make targets, ADRs, OpenAPI spec, smoke tests, and operational docs.

## Was Mode C Correct?

Yes for v1. Browser-only IMAP is not feasible with normal GitHub Pages constraints, and Mode B would require publishing private mailbox artifacts somewhere. The backend is justified because mailbox credentials, local indexes, and AI prompts need a private runtime boundary.

Could part of this have stayed Mode A? The demo UI, public repo links, and future browser-local search experiments can stay static. Real mailbox import cannot.

## What Worked

- GitHub Pages from `main` `/docs` was simple and visible from the first commit.
- OpenAPI-first frontend typing kept the API boundary explicit.
- SQLite gave a reliable local fallback while leaving room for Tantivy.
- Ollama fallback behavior makes the product usable even without a local model.

## What Did Not

- Full browser/WASM-side mailbox processing was not practical for v1 because IMAP access is the hard blocker.
- Native tool integrations such as `readpst`, Tantivy, sentence-transformers, and `age` are capability-detected but not all first-class execution paths yet.
- Local Go tooling inherited machine-level CGO flags for `onnxruntime`; Make targets force `CGO_ENABLED=0` and repo-local temp storage to keep checks stable.

## Surprises

- `go test ./...` can wander into `frontend/node_modules` because one package ships Go files. The Makefile avoids that with explicit Go package scopes.
- The temp volume on the development machine filled during `go vet`; using `GOTMPDIR=tmp/go-build` fixed it.

## Accepted Tech Debt

- Search uses SQLite `LIKE` fallback in v1, while Tantivy is exposed as a capability boundary for future replacement.
- PST import is not exposed as a UI workflow yet; `readpst` capability detection is present.
- The frontend uses browser demo fallback data when no backend is connected.

## Next Improvements

1. Add a real IMAP account setup flow with OAuth/device-code-friendly providers and local encrypted credential storage.
2. Replace fallback search with a Tantivy-backed index adapter and add semantic ranking through sentence-transformers.
3. Add age-encrypted export/import for portable local backups.

## Time Spent Vs Estimate

Estimate for a complete production-grade clone would be multiple weeks. This v1 scaffold and working local product surface was built in one focused implementation pass, with deliberate limitations documented above.

