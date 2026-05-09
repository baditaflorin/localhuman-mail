# Phase 3 Output Pathway Audit

Status legend: green = works end-to-end; yellow = partial; red = not usable; gray = deliberately out of scope.

| Pathway | Before | Evidence | Gap | Decision |
| --- | --- | --- | --- | --- |
| Copy generated reply | yellow | Reply appears in a textarea. | User can select manually, but no explicit copy confirmation. | Finish. |
| Copy message body | red | No copy control. | Users often need to paste normalized body elsewhere. | Finish. |
| JSON export | red | API returns JSON, but UI has no export. | User cannot take the normalized message/intelligence out. | Finish. |
| CSV export | red | No CSV generation. | Spreadsheet workflow is missing. | Finish for current message list. |
| Downloadable state file | red | No client state export. | Cannot move/reload a session without backend DB. | Finish with versioned JSON state. |
| Import state file round-trip | red | No importer. | Output cannot become input. | Finish. |
| Shareable URL | red | No hash-encoded state. | Small sessions cannot be shared or reopened from a link. | Finish with documented size limit. |
| Print view | yellow | Browser can print, but no print button and chrome remains. | Print/PDF output is noisy. | Finish minimal print CSS + button. |
| Screenshot | gray | No screenshot feature. | Browser/OS already handles this; not a mailbox workflow. | Out of scope. |
| Embed code | gray | Not relevant for private mailbox content. | Would encourage sharing private data. | Out of scope. |
| API/curl-ready | yellow | OpenAPI exists; README links it. | UI does not surface curl-ready evidence. | Finish by documenting verified API path; no UI control needed. |

## After Phase 3

| Pathway | After | Evidence |
| --- | --- | --- |
| Copy generated reply | green | `Copy draft` button writes to clipboard; Playwright verifies the success toast. |
| Copy message body | green | `Body` copy button writes selected normalized body; Playwright verifies the success toast. |
| JSON export | green | Detail toolbar downloads current filtered message list as JSON. |
| CSV export | green | Detail toolbar downloads current filtered message list as CSV. |
| Downloadable state file | green | `Export state` downloads `localhuman-mail-state.json`; Playwright verifies filename. |
| Import state file round-trip | green | Versioned zod-validated state import is implemented and unit-tested. |
| Shareable URL | green | Hash-encoded snapshots are generated for small current lists and copied. |
| Print view | green | Print button calls `window.print()` and print CSS hides controls. |
| Screenshot | gray | Still out of scope. |
| Embed code | gray | Still out of scope for private mailbox data. |
| API/curl-ready | green | README and `docs/api.md` point at OpenAPI; smoke verifies REST endpoints. |

Final counts: green 9, yellow 0, red 0, gray 2.
