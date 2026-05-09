# Phase 3 Controls Audit

| Control | Before | Does the label hold on real data? | Decision |
| --- | --- | --- | --- |
| Backend URL input | green | Persists to `localStorage` and drives the API client. | Keep; add reset through settings/state controls. |
| Connected/Demo status | green | Reflects version query success. | Keep. |
| Star link | green | Opens https://github.com/baditaflorin/localhuman-mail. | Keep. |
| PayPal link | green | Opens https://www.paypal.com/paypalme/florinbadita. | Keep. |
| Demo button | yellow | Backend seed works only when connected; offline demo already shown but button tries backend. | Clarify and keep. |
| EML button | yellow | Opens one-file picker and imports one selected `.eml`. | Extend to multi-file/batch. |
| Sync button | green | Invalidates React Query caches. | Keep. |
| Hidden file input | yellow | Accepts `.eml`, one file only. | Extend to `multiple`. |
| Search input | green | Filters current list client-side. | Keep; include new fields in filter. |
| Message row buttons | green | Select a message and clear draft. | Keep. |
| GitHub detail link | green | Opens repo. | Keep. |
| Tone buttons | green | Mutate tone state for reply draft. | Persist. |
| Draft button | green | Uses backend assist when online, browser fallback offline. | Keep. |
| Draft textarea | yellow | Editable but not saved, copied, or exported. | Persist and add copy. |

## After Phase 3

| Control | After | Evidence |
| --- | --- | --- |
| Backend URL input | green | Persists in versioned UI state and legacy key for compatibility. |
| Connected/Demo status | green | Still driven by `/api/v1/version`. |
| Star link | green | Unchanged and tested. |
| PayPal link | green | Unchanged and tested. |
| Demo button | green | Clearly imports backend demo messages when connected; offline demo remains visible. |
| EML files button | green | Multi-file picker, per-file progress, real fixture e2e. |
| Clipboard button | green | Reads clipboard and imports when content looks like raw `.eml`. |
| Sync button | green | Invalidates query cache. |
| Paste raw EML text box/import button | green | Validates raw headers and imports pasted content. |
| Import state | green | Validates JSON snapshot before applying. |
| Export state | green | Downloads versioned state file. |
| Clear local | green | Resets UI state to defaults. |
| Search input | green | Filters subject, sender, body, shape, confidence, warnings, attachments, and tags. |
| Message row buttons | green | Select message and clear draft. |
| Detail export/copy/share/print controls | green | Copy, JSON, CSV, state/share, print are wired. |
| Tone buttons | green | Persist in UI state. |
| Draft button | green | Backend or fallback assist works. |
| Draft textarea | green | Editable, persistent, copyable, exportable through state. |

Final counts: green 18, yellow 0, red 0.
