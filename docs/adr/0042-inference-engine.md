# 0042 - Inference Engine

## Status

Accepted

## Context

The app needs to infer obvious mailbox concepts without asking the user: shape, primary body, attachments, calendar invite, notification, receipt, mailing list, and confidence.

## Decision

Use a deterministic rule-based inference engine for Phase 2. Evidence comes from MIME parts, headers, subject/body keywords, attachment metadata, and URLs.

Inferred shapes:

- `personal_reply`
- `mailing_list`
- `notification`
- `receipt_invoice`
- `calendar_invite`
- `newsletter`
- `attachment_only`
- `unknown`

## Consequences

- Behavior is explainable and fixture-testable.
- Confidence can cite concrete evidence.
- Local LLMs are not needed for import-time inference.

## Alternatives Considered

- LLM classification: rejected for determinism and privacy.
- No classification: rejected because users need useful first guesses.

