# 0014 - Error Handling Conventions

## Status

Accepted

## Context

Mailbox parsing, external commands, local AI calls, and HTTP handlers all need consistent error handling.

## Decision

Use wrapped Go errors with `%w`, typed HTTP error responses, and an `internal/utils.HandleErrorOrLogWithMessages(err, errMsg, successMsg)` helper for command-style flows.

Frontend errors surface through an error boundary, API error views, and concise toasts. Never expose secrets, raw credentials, or full mailbox bodies in user-facing errors.

## Consequences

- Backend errors remain traceable in logs.
- Users get actionable failures.
- Sensitive data is not echoed into logs or browser UI.

## Alternatives Considered

- Panics: rejected.
- Silent fallbacks: rejected because ingestion and search failures must be visible.

