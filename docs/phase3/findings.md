# Phase 3 Findings Synthesis

## Top 5 Usability Gaps

1. Real users can import only one `.eml` at a time through a hidden picker; no drag/drop, batch, or paste path.
2. Users cannot export normalized messages, drafts, or session state from the UI.
3. Refresh loses query, selected message, tone, and draft, so the app does not feel like a real workbench yet.
4. Import failures are technically actionable in the API but still too compressed in the UI and not tied to a specific file in batch workflows.
5. README and ADRs imply more persistence and native-tool capability than the first screen actually exposes.

## Top 5 Half-Baked Features

| Feature | Decision | Rationale |
| --- | --- | --- |
| EML import | Finish | It is core to real-user use. |
| Demo button | Finish/clarify | Demo mode is useful offline; backend seed should not be confused with local sample data. |
| Draft textarea | Finish | It needs copy and persistence to be usable output. |
| Capability detection | Keep but reword | Detection is honest; execution workflows are not shipped. |
| Settings-like backend URL | Finish into minimal settings/state controls | It already persists and should gain reset/export/import. |

## Top 5 Codebase Pain Points

1. `frontend/src/App.tsx` owns too much behavior and slows every UI completeness change.
2. File import, export, state persistence, and share logic need single-purpose utilities.
3. Error rendering is repeated in mutation handlers instead of flowing through one user-facing helper.
4. Generated OpenAPI types are strong, but multipart upload currently needs an unsafe cast.
5. UI e2e tests cover presence, not real-user workflows.

## Top 5 Documentation/Reality Mismatches

1. ADR 0005 says last selected message/layout density persist; they do not.
2. README says mailbox import broadly; UI only imports single `.eml`.
3. Capability detection can sound like feature execution unless worded carefully.
4. API/curl readiness is linked, but no verified example is prominent.
5. Local boundary copy says UI settings only; after Phase 3, exported local session snapshots must be described precisely.

## Fully Usable Means

- A stranger can load one or many `.eml` files by picker, drag/drop, or paste and see per-file results.
- They can search/select, generate or edit a draft, copy it, and export messages/draft/state without asking where the data went.
- A reload restores non-sensitive UI work context without storing mailbox contents in frontend storage.
- Every visible control either completes its label's promise or is documented as out of scope.
- README claims map to a passing test or are removed.

## Success Metrics

- Input audit: all non-gray rows green or ADR-documented.
- Output audit: all non-gray rows green or ADR-documented.
- Controls audit: 100% green for production controls.
- TypeScript unsafe casts in `frontend/src` excluding generated code: 0.
- Real-user e2e paths covered: upload fixture, copy draft/body, export/import state, share hash restore.
- `make test`, `make lint`, `make smoke`, pre-push hook chain pass.

## Out Of Scope

- New mailbox engines, IMAP sync, PST execution UI, OCR, screenshots, embed code, hosted auth, cross-device sync, and visual polish.
