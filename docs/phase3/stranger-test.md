# Phase 3 Stranger Test

Date: 2026-05-09

Method: substitute cold walkthrough in a fresh Playwright browser context against the built Pages app and a local backend. Input was the real fixture `test/fixtures/realdata/06-calendar-invite-rfc5545.eml`, not the canned demo.

## Walkthrough

1. Opened the Pages build.
2. Set backend URL to `http://127.0.0.1:18080`.
3. Uploaded the calendar `.eml` through the file picker.
4. Searched for `Design Review`.
5. Copied message body.
6. Generated and copied a draft.
7. Downloaded a state file.
8. Created a share URL.

## Findings

| Finding | Severity | Response |
| --- | --- | --- |
| Success toast blocked the Draft button and trapped clicks. | high | Made toast pointer-events disabled in `frontend/src/components/Toast.tsx`. |
| Share URL used the whole current source list and exceeded the limit during a filtered workflow. | high | Share/export toolbar now uses the current filtered list; share encoding uses compact JSON; URL limit raised to 20KB for small snapshots. |
| Playwright shape assertion was ambiguous because the shape appears in both list and detail. | low | Tightened e2e assertion to the detail confidence row. |

## Result

The cold path now passes: real `.eml` upload, search, body copy, draft generation/copy, state download, and share URL copy all complete without manual help.
