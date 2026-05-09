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

Initial counts: green 5, yellow 4, red 1.
