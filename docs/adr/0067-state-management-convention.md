# 0067 - State Management Convention

## Status

Accepted

## Context

Server data uses React Query, while local UI state used a single `localStorage` helper and transient `useState`.

## Decision

Server message data stays in React Query and the backend store. Non-sensitive UI state (backend URL, query, selected message ID, tone, draft, imported snapshot messages) is stored in a versioned localStorage record. Mailbox credentials and backend-imported mailbox contents are not silently mirrored into browser storage.

## Consequences

- Reload restores work context without violating the local boundary claim.
- Snapshot import/export is explicit user action, not hidden sync.
- State migration is required for future localStorage shape changes.

## Alternatives Considered

- Store every fetched backend message in localStorage: rejected because it quietly duplicates mailbox data into the browser.
