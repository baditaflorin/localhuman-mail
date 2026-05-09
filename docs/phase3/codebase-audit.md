# Phase 3 Codebase Health Audit

Measured before implementation on 2026-05-09.

## Size And Responsibility

| Module | Lines | Finding |
| --- | ---: | --- |
| `frontend/src/App.tsx` | 438 | God component: API orchestration, upload handling, persistence, rendering, export candidates, and controls live together. |
| `internal/mailbox/parser.go` | 528 | Parser does MIME walking, normalization, HTML extraction, inference, warning generation, and provenance in one file. Phase 3 will not change engine behavior, but the boundary is documented debt. |
| `internal/search/store.go` | 291 | Store migration, insertion, search, and JSON serialization live together; below the 300-line guard but near the limit. |

## DRY Violations

| Area | Paths | Finding | Decision |
| --- | --- | --- | --- |
| Demo message metadata | `internal/mailbox/demo.go`, `frontend/src/features/mail/demo.ts` | Backend and frontend define parallel demo messages and confidence metadata. | Accept for now because frontend demo works offline and backend seed works online; document in ADR 0064. |
| Error formatting | `internal/api/json.go`, `frontend/src/api/client.ts`, UI mutation `onError` handlers | Actionable errors are shaped in backend but composed ad hoc in frontend. | Consolidate frontend error helper and add per-file batch results. |
| Import state | Upload mutation and hidden file input in `App.tsx` | Input flow is embedded in UI rendering. | Extract local utilities for file/state/export logic. |

## Dead Code

No unreferenced source files found by inspection. No commented-out code blocks in source. Generated `docs/assets/*` are build artifacts and excluded from source-health counts.

## TODO / FIXME / XXX / HACK

No source-code TODO/FIXME/XXX/HACK markers found. Existing "future" language is in docs/postmortem and ADRs, not executable code.

## Type Safety Holes

| Path | Finding | Decision |
| --- | --- | --- |
| `frontend/src/App.tsx` | `body as never` is used to satisfy `openapi-fetch` multipart typing. | Replace with a typed API helper at the boundary. |
| `frontend/src/api/client.ts` | Error response narrowing uses a cast. | Replace with schema-based validation or local type guard. |
| `internal/ai/service.go` | `map[string]any` builds Ollama JSON payload. | Accept as JSON-boundary Go code; document. |
| `internal/search/scan.go` | `Scan(dest ...any)` mirrors `database/sql`. | Accept as database-boundary Go code; document. |

## Inconsistent Patterns

- Frontend state mixes `useState`, `localStorage` hook, and server cache without a single UI-state policy.
- Search is client-side filtering over list data while backend search is available through `/api/v1/search`; current UI label does not clarify this.
- Backend errors are now domain-shaped, but frontend still shows one-line toast only.

## Test Coverage Holes

- No e2e coverage for uploading a real `.eml`.
- No tests for multi-file import, state export/import, share URL, copy controls, or persistence migration.
- Phase 2 fixture tests cover parser behavior, not UI completeness.

## After Phase 3

| Metric | Before | After | Evidence |
| --- | ---: | ---: | --- |
| DRY issues in core workflow modules | 3 | 0 blocking | Import/export/state/share logic moved into `frontend/src/features/mail/workbench.ts` and `frontend/src/lib/browser.ts`; frontend/backend demo duplication remains accepted by ADR 0064. |
| Source TODO/FIXME/XXX/HACK | 0 | 0 | `rg` audit. |
| Unsafe authored TypeScript casts / `any` / `@ts-ignore` | 3 | 0 | `rg " as never|@ts-ignore|\\bany\\b" frontend/src --glob '!api/schema.d.ts'` returns no matches. |
| Production UI stubs | 0 | 0 | Controls audit green. |
| Real-user path test gaps | 5 | 0 blocking | Unit tests cover state/export helpers; Playwright covers upload/copy/draft/state/share. |

Accepted residual debt: `frontend/src/App.tsx` remains a large composition component because splitting presentational layout further would be polish/refactor work after the Phase 3 functional gates. The core logic it used to own has been extracted and tested.
