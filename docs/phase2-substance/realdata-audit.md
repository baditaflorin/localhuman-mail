# Phase 2 Substance Real-Data Audit

Date: 2026-05-08

Scope: current v1 primary flow, `EML upload -> parse -> store -> list/search -> draft reply`.

## 10 Real-World Inputs

| # | Input | What v1 does | What it should do | Why it fails or falls short | Failure mode | Manual work forced onto user |
|---|---|---|---|---|---|---|
| 1 | Clean RFC 5322 single-part `text/plain` email with From/To/Subject/Date | Imports, stores, lists, searches | Same, plus show parse confidence and provenance | Happy path only; no confidence or parse metadata | Mostly succeeds, but opaque | User trusts it blindly |
| 2 | Real mailing-list reply with quoted thread, `In-Reply-To`, `References`, list footer, unsubscribe headers | Imports as one flat body | Identify thread headers, quoted text, footer/list metadata, primary reply text | Parser treats all text as equal | Wrong-but-confident | User must mentally separate reply vs quoted history |
| 3 | GitHub/GitLab-style notification with `multipart/alternative`, HTML, plaintext, list headers, long machine-generated footer | Imports plaintext if present, ignores HTML/list semantics | Prefer useful plaintext, detect notification type, extract repo/thread/source URL | No domain-aware shape detection | Silent under-extraction | User must search full noisy body |
| 4 | HTML-only transactional receipt/newsletter with quoted-printable HTML and no `text/plain` part | Rejects upload as "could not parse EML" | Extract readable text from HTML, keep source HTML metadata, warn if confidence is low | Parser only accepts `text/plain` inline parts | Obvious but unhelpful | User has to convert HTML to text elsewhere |
| 5 | Outlook/Exchange forwarded email with CP1252 smart quotes, NBSP, CRLF, malformed folded headers | May import garbled text or fallback fields; invalid date becomes current time | Normalize encoding, preserve original, report repaired headers/date confidence | No explicit normalization policy; date fallback is nondeterministic | Wrong-but-confident | User must notice corrupted text/date |
| 6 | Google Calendar/Outlook invite where useful content is `text/calendar` plus ICS attachment | Usually rejects if no text body; otherwise ignores ICS | Recognize meeting invite, extract summary, organizer, start/end, location, RSVP URL | No calendar domain handling | Obvious but dumb | User must open another app |
| 7 | Attachment-heavy invoice email with short body and PDF attachment | Imports only the short body, drops attachment metadata | Surface attachment names/types/sizes and infer "invoice/receipt" from headers/body | Attachments are skipped completely | Silent wrongness | User cannot tell invoice attachment exists |
| 8 | Maildir folder with `cur/`, `new/`, duplicate messages, flags, mixed line endings | No UI/API happy path; only one EML upload exists | Accept directory/batch import through backend, dedupe by stable message id, preserve flags | v1 only handles one uploaded `.eml` file | Missing path | User must upload one message at a time |
| 9 | Truncated/corrupted EML from interrupted export, with headers and partial MIME boundary | Rejects with generic parse error or no-readable-body error | Detect truncation/partial MIME, import recoverable fields, explain what is missing | Error taxonomy is generic and not domain-specific | Obvious but unactionable | User does not know whether to retry export or fix file |
| 10 | Huge real newsletter/archive email, 8-15MB, deeply nested MIME, long base64/quoted-printable sections | Reads body into memory and blocks until done; no progress/cancel | Stream parse, size-budget, progress, cancellation, safe body limits | No performance budget or cancellation model | Slow/stuck state risk | User waits without knowing whether it is working |

At least these inputs are visibly mishandled today: #4 HTML-only email, #6 calendar invite, #7 attachment-heavy invoice, #8 Maildir batch, #9 truncated EML, #10 huge nested MIME.

## Top 5 Logic Gaps

1. The parser only understands `text/plain` bodies and treats every accepted message as a flat blob.
2. Invalid or missing structured fields are silently replaced with fallbacks, including nondeterministic `Date: now`.
3. Attachments, calendar parts, list headers, thread headers, unsubscribe headers, and source URLs are discarded.
4. Import is single-message only; real mailbox data arrives as batches, Maildir, PST exports, IMAP syncs, and duplicates.
5. Search and assist operate on noisy raw body text with no extracted primary text, quoted-text separation, or confidence metadata.

## Top 3 Intuition Failures

1. Uploading a real HTML-only email fails even though the email is perfectly readable in every mail client.
2. An invoice email imports successfully while silently hiding the actual invoice attachment.
3. A malformed date becomes "now", making the app look confident while corrupting timeline/search behavior.

## Top 3 Feels-Stupid Moments

1. The user has to know whether an email contains `text/plain`; the app should infer readable text from MIME/HTML.
2. The user has to upload messages one at a time even though mailbox data is naturally a folder/export/account.
3. The user has to inspect raw message noise manually because the app does not separate useful content from quoted replies, footers, and machine boilerplate.

## What Smart Means For localhuman-mail

- A messy EML should produce a useful first preview: subject, sender, recipients, date, primary readable body, detected shape, and confidence.
- The importer should recover what it can, explain what it repaired or could not trust, and never silently invent important fields.
- Common email shapes should be recognized: personal reply, mailing-list post, notification, receipt/invoice, calendar invite, newsletter, attachment-only message.
- Search should index normalized primary text plus domain metadata, not just raw MIME body text.
- Every imported message should carry deterministic provenance: source id, parser version, normalized fields, warnings, and confidence.

## Phase 2 Substance Success Metrics

- At least 7 of the 10 audit fixtures complete import -> useful preview -> searchable result with no manual correction.
- Re-importing the same fixture produces byte-identical normalized JSON for all deterministic fields.
- No fixture produces a silently invented date, sender, body, or attachment state.
- Every failed fixture returns a domain-specific what/why/now-what error.
- Fixtures up to 10MB do not freeze the UI; operations over 300ms expose progress, and operations over 5s are cancellable.
- Search finds the intended message on at least 8 of 10 fixtures using an obvious user query.

## Out Of Scope

- No new visual polish, dark mode, command palette, onboarding, or marketing work.
- No architecture escalation beyond Mode C.
- No hosted mailbox proxy or SaaS sync.
- No new user-facing product surface beyond making the existing import, list, search, and draft flows smarter.
- No Phase 3 polish work until the fixture pass rate and failure honesty improve.

## After Phase 2 Substance

The 10 real-data fixtures are committed under `test/fixtures/realdata/` and enforced by `internal/mailbox/fixture_test.go`.

| Metric | Before | After |
| --- | ---: | ---: |
| Useful import/classification pass rate | 2/10 | 10/10 |
| Deterministic fixture output | untested | 10/10 |
| Actionable oversize/corrupt input errors | partial | tested |
| Confidence/warnings surfaced by parser/API | no | yes |

Closed gaps:

1. Plaintext-only parsing now handles mbox envelopes, HTML fallback, non-UTF charsets, partial MIME recovery, attachments, calendars, and quoted reply trimming.
2. Message shape inference now identifies personal replies, mailing lists, notifications, newsletters, attachment-only messages, receipts, calendars, and unknowns with confidence.
3. Parser output now carries per-field confidence, warnings, provenance, source hashes, schema version, parser version, attachment metadata, and calendar metadata.
4. SQLite storage and OpenAPI now preserve the richer message model instead of dropping it after parse.
5. Fixture tests prevent silent regressions on the 10 real inputs.
