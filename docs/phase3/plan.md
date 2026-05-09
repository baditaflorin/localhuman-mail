# Phase 3 Implementation Plan

Ranked by real-user impact on the audit findings.

| Rank | Catalog Item | Work | Evidence Target |
| ---: | --- | --- | --- |
| 1 | #1 Input pathways | Add drag/drop, multi-file picker, paste text, clipboard read fallback. | Input audit non-gray rows green. |
| 2 | #4 Batch input | Batch import `.eml` files with progress and per-file results. | UI shows imported/skipped/errors per file. |
| 3 | #2 Format detection | Sniff extension, MIME, and raw header pattern before import. | Wrong file gets actionable guidance. |
| 4 | #8 Resume | Persist query, selected message, tone, and draft with versioned UI state. | Reload restores non-sensitive work context. |
| 5 | #9 Export formats | Add JSON and CSV exports for current message list. | Downloaded files contain real messages. |
| 6 | #10 Copy output | Add copy body, copy draft, copy JSON controls. | Clipboard tests pass. |
| 7 | #11 State file | Export/import a versioned state snapshot. | Round-trip restores UI state and messages snapshot. |
| 8 | #12 Share URL | Hash-encode small state snapshots with documented limit. | New page load restores shared state. |
| 9 | #13 Print/PDF | Add print button and print CSS that hides controls. | Print view includes selected message and draft. |
| 10 | #15 Half-baked triage | Finish/clarify EML, demo, draft, capability, settings/state features. | ADR 0063 and README match behavior. |
| 11 | #18 Settings completeness | Add minimal settings/state panel: backend URL reset, clear UI state, export/import. | Every setting does something. |
| 12 | #19 Help/docs alignment | Update README feature checklist and limitations. | Feature claims audit green. |
| 13 | #20 DRY extraction | Extract import/export/state helpers from `App.tsx`. | App component shrinks and duplicate logic drops. |
| 14 | #21 Consolidation | One frontend error formatter and one download/copy helper. | Reused by all controls. |
| 15 | #22 Canonical types | Define local UI state and snapshot types in one module. | No parallel snapshot shapes. |
| 16 | #23 Validation schemas | Zod-validate local snapshots, share payloads, and GitHub API. | Invalid imports are actionable. |
| 17 | #24 Split god module | Move file input/export/share logic into `features/mail/workbench.ts`. | Smaller UI component. |
| 18 | #26 Inject dependencies | Pass copy/download/import helpers through functions, no hidden globals except browser boundary helpers. | Unit-testable helpers. |
| 19 | #28 Dead code | Remove unused or misleading claims, no dormant controls. | Codebase audit updated. |
| 20 | #31 Error convention | Apply what/why/now-what wording to UI import errors. | Per-file errors are understandable. |
| 21 | #32 State convention | Keep server data in React Query, non-sensitive UI state in versioned localStorage, snapshots explicit. | ADR 0067. |
| 22 | #35 Type safety | Remove unsafe casts from app code. | `rg " as never|@ts-ignore|any"` clean in frontend source. |
| 23 | #36 Boundary validation | Validate imported state and share URL payloads with zod. | Bad state file does not crash. |
| 24 | #38 Persistence | Save draft/tone/query/selection; clear-state works. | Vitest coverage. |
| 25 | #39 Migration | Add versioned UI state migration from missing/old shapes. | Old storage keys do not break app. |
| 26 | #41 Round-trip | Exported state imports back equivalently for supported fields. | Unit test. |
| 27 | #42 README checklist | Verified feature checklist with test names. | Claims are traceable. |
| 28 | #43 Quickstart | Ensure quickstart still reaches canonical story. | `make smoke` and docs. |
| 29 | #44 Inline help | Add short labels/ARIA descriptions for non-obvious import/export controls. | Controls audit green. |
| 30 | #45 Limitations | Honest limitations section for no IMAP/PST execution UI yet. | Docs mismatch closed. |
| 31 | #46 Stranger test | Private-window walkthrough with real fixture. | `docs/phase3/stranger-test.md`. |
| 32 | #47 Fix top 3 | Fix confusing dead-ends from stranger test. | Postmortem evidence. |

Implementation order: input completeness, output completeness, half-baked triage, codebase health, consistency/type safety, persistence/docs, stranger test, release.
