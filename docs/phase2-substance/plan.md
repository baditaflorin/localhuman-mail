# Phase 2 Substance Plan

Ranking rule: user impact on the 10 real-data audit inputs, not implementation novelty.

## Selected Substance Items

1. #32 Actionable errors: every import/search failure says what, why, and now what.
2. #33 Validate at boundaries: EML upload gets size/type/body validation before storage.
3. #2 Encoding variants: normalize UTF-8 BOM, CRLF/LF, CP1252-compatible charsets, NBSP, smart whitespace.
4. #9 Format normalization: deterministic dates, normalized body text, stable source hashes.
5. #35 Deterministic outputs: fixture tests assert stable normalized output.
6. #38 Output provenance: every import carries source hash, parser version, warnings, and schema version.
7. #16 Confidence scores: message-level and field-level confidence.
8. #18 Surface anomalies: parser warnings for invalid dates, missing body, attachments, HTML fallback, malformed headers.
9. #4 Partial inputs: recover useful headers/body when possible and explain unrecoverable truncation.
10. #5 Adversarial input: malformed MIME, duplicate headers, Unicode lookalikes, oversized parts.
11. #3 Huge inputs: 25MB upload cap, indexed body cap, documented cliff.
12. #6 Auto-detect structure: classify common email shapes from MIME/header/body evidence.
13. #7 Auto-classify fields: sender/date/body/attachments/calendar/list fields get confidence.
14. #8 Useful first guess: import returns a message with primary body, shape, warnings, and confidence immediately.
15. #11 Domain vocabulary: API/UI uses mailbox words, not parser internals.
16. #12 Domain-aware validation: invoice attachments, calendar parts, suspicious dates, HTML-only mail.
17. #13 Recognize common shapes: personal reply, mailing list, notification, receipt/invoice, calendar invite, newsletter.
18. #14 Domain-aware export/API: message JSON includes provenance, shape, confidence, warnings, attachments.
19. #15 Domain conventions: MIME semantics, List-* headers, thread headers, text/html fallback, quoted reply trimming.
20. #17 Suggest fixes: warnings include next steps.
21. #19 Explain decisions: each inference records concise reasons.
22. #22 Stable IDs: source hash + Message-ID-aware deterministic IDs.
23. #24 State taxonomy: documented states for import/search/draft.
24. #25 No stuck states: every recoverable failure keeps prior mailbox data and suggests a next step.
25. #27 Concurrency safety: duplicate demo/import requests dedupe deterministically.
26. #28 Profile real-data inputs: fixture benchmark records median/worst parse duration.
27. #30 Stream where possible: backend caps and truncates indexed body instead of unbounded reads.
28. #31 Cache expensive things: normalized parser output is stored once and reused by search/draft.
29. #34 Recoverable vs fatal: error taxonomy distinguishes bad input from server failure.
30. #37 Debug surface: `?debug=1` shows parser confidence/warnings/provenance in the existing message panel.

## Implementation Order

1. Fixtures and expected properties.
2. Parser model: metadata, confidence, provenance, warnings, attachments, shape.
3. MIME/HTML/charset/body normalization.
4. Store/API/OpenAPI migration.
5. Search over primary normalized content and metadata.
6. UI surfacing for confidence/warnings/provenance/debug.
7. Performance/determinism tests and postmortem.

## Non-Goals

- No new mailbox sync setup flow.
- No visual redesign.
- No hosted backend.
- No semantic vector search replacement in this phase.
- No full attachment text extraction.

