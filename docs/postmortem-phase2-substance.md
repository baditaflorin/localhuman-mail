# Phase 2 Substance Postmortem

## Real-Data Pass Rate

Before: 2/10 fixtures produced a useful import/search result without manual recovery.

After: 10/10 fixtures pass `internal/mailbox/fixture_test.go` with expected shape, required text, minimum confidence, warnings, attachment counts, and deterministic output.

## Top 5 Logic Gaps

1. Plaintext-only body parsing: closed with MIME walking, HTML-to-text fallback, charset decoder registration, and partial MIME recovery.
2. No domain shape inference: closed with deterministic shape classification and per-field confidence.
3. No warning/provenance model: closed with warnings, parser/schema versions, source hash, and size metadata.
4. Attachment/calendar blindness: closed with attachment metadata and `text/calendar` extraction.
5. Parser output dropped at storage/API: closed by extending SQLite schema, scan paths, OpenAPI, and frontend demo data.

## Smart Behaviors Delivered

- Importing real `.eml` exports produces a useful first guess: shape, readable body, snippet, tags, warnings, and confidence.
- Broken or partial inputs recover when possible and fail with a what/why/now-what error when not.
- Same input produces deterministic message JSON across repeated parses.
- Low-trust facts show up as warnings or lower confidence instead of silent certainty.

## Determinism

All 10 fixtures pass repeated parse JSON equality after confidence scores were rounded to stable precision.

## Performance

`go test ./internal/mailbox` completes under one second on the fixture set. The importer caps complete `.eml` reads at 25MB and indexed part reads at 2MB.

## Surprises

The public SpamAssassin corpus contains enough mbox envelope lines, malformed MIME, and charset weirdness to break a strict happy-path parser quickly. The most useful improvement was not "more AI"; it was tolerant, explicit parsing with honest warnings.

## Still Open

1. PST/Maildir/IMAP import are still capability boundaries, not first-class workflows.
2. Tantivy and sentence-transformers are still not wired as production search/ranking paths.
3. Parser internals remain large and should be split after the behavior is stable.
4. HTML extraction is readable but not layout-aware.
5. Attachment contents are detected, not indexed.

## Honest Take

The app no longer feels like a toy at the parser layer. It can ingest messy real mail and tell the user what it thinks, how sure it is, and what may be missing. Before Phase 3 it still felt incomplete as a product because users could not batch import or export their work; that is why Phase 3 followed immediately.
