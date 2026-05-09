# 0040 - Real-Data Audit Findings and Substance Success Metrics

## Status

Accepted

## Context

The v1 happy path imports simple `text/plain` EML files, but real mailbox data contains HTML-only messages, attachments, calendar parts, malformed headers, quoted replies, batch exports, and huge nested MIME structures.

## Decision

Use the 10 inputs in `docs/phase2-substance/realdata-audit.md` as the Phase 2 grading rubric. Phase 2 succeeds when at least 7 of 10 fixtures produce a useful preview/search result without manual correction, deterministic normalized output, no silent invented critical fields, and actionable errors for failures.

## Consequences

- Fixture pass rate becomes the primary quality metric.
- Parser correctness is more important than UI expansion.
- Any regression on a real-data fixture blocks a push unless an ADR explains the tradeoff.

## Alternatives Considered

- Continue polishing the demo flow: rejected because it does not improve real mailbox handling.
- Add new mailbox setup features first: rejected because the parser would still mishandle real inputs.

