# 0041 - Input Robustness and Normalization Policy

## Status

Accepted

## Context

Email inputs vary by MIME structure, transfer encoding, charset, line endings, malformed headers, and size. v1 has no explicit policy.

## Decision

Normalize at the import boundary:

- Cap EML upload at 25MB for Phase 2.
- Normalize CRLF/LF, UTF-8 BOM, NBSP, repeated whitespace, and common charset variants through `go-message` charset support.
- Prefer `text/plain`; fall back to extracted `text/html` text when plaintext is absent.
- Preserve raw source hash for provenance.
- Use deterministic fallback values for missing fields and mark them low-confidence.

## Consequences

- Same input yields stable output.
- HTML-only and charset-weird email become recoverable.
- Oversized inputs get actionable failure messages.

## Alternatives Considered

- Store raw MIME only and defer normalization: rejected because search/draft need useful first output.
- Unlimited reads: rejected because huge inputs can exhaust memory.

