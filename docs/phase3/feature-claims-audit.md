# Phase 3 Feature Claims Audit

| Claim Source | Claim | Before | Reality | Decision |
| --- | --- | --- | --- | --- |
| README | GitHub Pages frontend with repo, PayPal, version, and commit metadata visible. | green | UI shows all four. | Keep tested. |
| README | Go backend with health, readiness, metrics, version, mailbox import, search, reply assist. | green | Routes exist and smoke covers a happy path. | Keep. |
| README | EML import, demo import, SQLite storage, local search. | yellow | Single EML works; real users need batch import and clearer errors. | Finish. |
| README | Local LLM reply assist through Ollama, with fallback. | green | Implemented in `internal/ai`. | Keep. |
| README | Capability detection for native tools. | yellow | Detection is present; integrations are not first-class workflows. | Reword as detection-only. |
| README | Docker Compose production topology with nginx on port 25342. | green | Deploy files exist. | Keep. |
| ADR 0005 | Frontend remembers layout density and last selected message. | red | Only backend URL persists. | Fix docs or implement persistence. |
| ADR 0006 | Native mailbox tooling can run in backend. | yellow | Capability detection only, no execution UI. | Reword as future integration boundary. |
| docs/postmortem.md | PST import not exposed yet. | green | Correct limitation. | Keep. |
| In-app copy | Local boundary says only UI settings are stored. | yellow | Draft/query not persisted yet, mailbox content not stored in frontend. | Keep; clarify once UI state export exists. |

## After Phase 3

| Claim | After | Evidence |
| --- | --- | --- |
| Pages frontend with repo, PayPal, version, commit | green | Existing e2e and UI metadata. |
| Backend health/readiness/metrics/version/import/search/assist | green | Smoke script checks health, readiness, version, demo import, search. |
| EML import, demo import, SQLite storage, local search | green | Batch EML and real fixture e2e added; parser fixtures pass. |
| Local LLM assist with fallback | green | Existing backend route and frontend fallback. |
| Capability detection | green | README now phrases optional tools as detection/boundaries, not full workflows. |
| Docker Compose topology | green | Deploy files unchanged; final Docker build/push is attempted as release step. |
| ADR 0005 persistence claim | green | UI state persistence now covers backend URL, query, selected message, tone, and draft; layout density remains removed from claims. |
| Native mailbox tooling execution | green | Docs state capability detection only unless a workflow exists. |
| Local boundary copy | green | UI now says browser state export/import is explicit. |

Final counts: green 9, yellow 0, red 0.
