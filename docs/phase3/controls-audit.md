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

Initial counts: green 9, yellow 5, red 0.
