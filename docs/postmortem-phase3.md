# Phase 3 Postmortem

## Audit Grids

| Audit | Before | After |
| --- | --- | --- |
| Input pathways | green 1, yellow 4, red 6, gray 4 | green 10, yellow 0, red 0, gray 4 |
| Output pathways | green 0, yellow 3, red 6, gray 2 | green 9, yellow 0, red 0, gray 2 |
| Controls | green 9, yellow 5, red 0 | green 18, yellow 0, red 0 |
| Feature claims | green 5, yellow 4, red 1 | green 9, yellow 0, red 0 |

## Half-Baked Feature Triage

| Feature | Outcome | Rationale |
| --- | --- | --- |
| EML import | finished | Added picker multi-file, drag/drop, paste, clipboard, batch results. |
| Demo button | finished/clarified | Backend seed remains, offline demo remains. |
| Draft textarea | finished | Persists, copies, exports through state. |
| Backend URL setting | finished | Versioned state, reset, import/export state. |
| Capability detection | kept/reworded | Detection is true; execution workflows remain limitations. |
| URL/screenshot/folder import | out of scope | ADR 0061 documents why. |

## Codebase Health

| Metric | Before | After |
| --- | ---: | ---: |
| DRY issues in core workflows | 3 | 0 blocking |
| Source TODO/FIXME/XXX/HACK | 0 | 0 |
| Unsafe authored frontend casts / `any` / `@ts-ignore` | 3 | 0 |
| Production UI stubs | 0 | 0 |
| Real-user path test gaps | 5 | 0 blocking |

Residual accepted debt: `frontend/src/App.tsx` remains a large composition component. The core workflow logic it used to own now lives in tested helpers, but presentational splitting is a Phase 4 candidate.

## Stranger Test

Documented in `docs/phase3/stranger-test.md`.

Top 3 fixed:

1. Toast blocked the Draft button; toast now ignores pointer events.
2. Share URL was too large for a filtered single-message workflow; share/export now uses current filtered messages and compact encoding.
3. E2E assertion ambiguity hid duplicate shape labels; test now targets the detail confidence row.

## Documentation Reality

Fixed README feature claims, current limitations, ADR 0005 persistence language, and Phase 3 audits. Optional native tools are now described as capability detection unless a workflow exists.

## Surprises

The biggest usability bug was not missing AI. It was a toast sitting over a button. Also, share links balloon quickly when message provenance/confidence is included, so the app now draws a line: small filtered snapshots use URLs, larger sessions use state files.

## Still Open

1. Split `App.tsx` into presentational sections without changing behavior.
2. Build first-class PST/Maildir/IMAP workflows.
3. Wire Tantivy and sentence-transformers as actual search/ranking backends.
4. Add attachment content indexing/opening.
5. Test mobile file import on physical iOS/Android devices.

## Honest Take

Yes, a stranger can now use the app for a real `.eml` workflow end-to-end: load real mail, inspect inferred shape/confidence/warnings, search, draft, copy, export, restore, and share a small snapshot without asking where the controls are. It is still not a full Superhuman replacement because account sync, PST/Maildir, production-grade ranking, and attachment contents are not wired. But it no longer feels like a polished demo; it behaves like a usable local workbench for exported email.
