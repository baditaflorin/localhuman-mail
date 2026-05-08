# 0006 - WASM Modules

## Status

Accepted

## Context

The original concept mentions browser/WASM-side mailbox parsing and search. ADR 0001 chooses Mode C because private IMAP access cannot be completed from a static browser app.

## Decision

Do not ship WASM modules in v1. Keep the API and frontend lazy-loading structure compatible with future WASM modules, but place mailbox parsing and search execution in the local Docker backend for v1.

## Consequences

- GitHub Pages does not require COOP/COEP workarounds for v1.
- Initial JS payload stays below the budget more easily.
- Native mailbox tooling such as `readpst`, `file`, and optional Tantivy integration can run in the backend container.

## Alternatives Considered

- DuckDB-WASM/sql.js for browser-local indexes: rejected because mailbox data would still need an import path from IMAP.
- Tantivy compiled to WASM: deferred until a browser-only import story is realistic.
- Service worker header shims: unnecessary without WASM in v1.

