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

Initial counts: green 0, yellow 3, red 6, gray 2.
