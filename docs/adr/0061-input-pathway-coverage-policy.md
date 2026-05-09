# 0061 - Input Pathway Coverage Policy

## Status

Accepted

## Context

The app accepted only one `.eml` through a hidden file input. Real users bring batches, drag-and-drop files, copied raw message text, mobile file pickers, and sometimes invalid files.

## Decision

Support `.eml` input through multi-file picker, drag-and-drop, raw text paste, and clipboard-read-with-fallback. Batch import reports per-file success/error and keeps partial success. URL, screenshot/OCR, rendered HTML paste, and folder import remain out of scope for v1.

## Consequences

- The frontend owns browser input ergonomics; the backend remains the parser authority.
- Invalid inputs are rejected before upload when possible and still receive server-side validation.
- Batch operations must never hide per-file failures.

## Alternatives Considered

- Add server-side URL fetch: rejected because private mailbox URLs require auth and secrets.
- Folder import: rejected because cross-browser support is uneven and multi-file covers the v1 need.
