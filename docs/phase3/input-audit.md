# Phase 3 Input Pathway Audit

Status legend: green = works end-to-end on real data; yellow = partial or backend-dependent; red = claimed/expected but not usable; gray = deliberately out of scope for this product mode.

| Pathway | Before | Evidence | Gap | Decision |
| --- | --- | --- | --- | --- |
| Single `.eml` file picker | yellow | Hidden file input posts one file to `/api/v1/import/eml`; Phase 2 fixtures parse 10 real messages. | No progress, only first selected file, backend required, error lacks per-file context in UI. | Finish. |
| Drag-and-drop `.eml` | red | No drop handler in `frontend/src/App.tsx`. | Common desktop mail-export workflow is missing. | Finish. |
| Multi-file batch `.eml` | red | File input does not use `multiple`; API imports one multipart file. | Real users export batches, not single messages. | Finish. |
| Paste raw `.eml` text | red | No paste box or text import endpoint. | Users cannot paste an exported message or copied source. | Finish with browser-side upload as `Blob`. |
| Paste rendered HTML | gray | No HTML-email raw import path. | Importing rendered HTML loses RFC 5322 headers and is not equivalent to email import. | Out of scope; guide users to paste raw `.eml` text instead. |
| Paste image | gray | No OCR or image mailbox import claim. | Email intelligence should not infer from screenshots in Phase 3. | Out of scope. |
| URL input | gray | Browser CORS blocks arbitrary mailbox or webmail URLs; no server fetch endpoint. | A URL cannot safely represent private mailbox data without auth/secrets. | Permanently out of scope for Mode C v1; document paste/upload instead. |
| Clipboard read button | red | No `navigator.clipboard.readText` pathway. | Browser permission path not exercised. | Finish with fallback to paste text area. |
| Mobile picker | yellow | Browser file input exists; no `multiple`, no mobile guidance. | It should accept `.eml` from Files on mobile. | Finish through standard file input and document limitation. |
| Folder import | gray | No File System Access API or server folder import. | Browser folder upload has inconsistent support and can expose too much. | Out of scope; use multi-file. |
| Demo sample | green | Demo data appears without backend; backend demo import works when connected. | Demo import label does not distinguish local sample vs backend seed. | Keep and clarify. |
| Deep links | yellow | Pages SPA fallback exists; no message selection in URL hash. | Users cannot share or restore selected message view. | Finish for small client-side share/state snapshots. |
| Imported state file | red | No state import format. | Users cannot leave and reload a curated session elsewhere. | Finish. |
| Restored autosave | yellow | Backend URL persists in `localStorage`; selected message, query, tone, draft do not. | Refresh loses work-in-progress UI context. | Finish with versioned UI state. |

Initial counts: green 1, yellow 4, red 6, gray 4.
