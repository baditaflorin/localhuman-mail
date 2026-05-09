# 0005 - Client-Side Storage Strategy

## Status

Accepted

## Context

The frontend needs to remember non-sensitive preferences such as backend URL, query, tone, draft, and last selected message. It must not store mailbox credentials or silently mirror backend mailbox contents.

## Decision

Use versioned `localStorage` for small non-sensitive UI settings in v1. Do not store mailbox credentials, tokens, private keys, or hidden backend mailbox mirrors in the browser. Browser state may contain message snapshots only after an explicit user action such as importing or sharing a state file.

## Consequences

- Users can reconnect to their local backend without retyping the URL and can restore draft/query/selection after reload.
- Clearing browser storage removes all client-side state.
- Cross-device sync is explicitly out of scope.

## Alternatives Considered

- IndexedDB: deferred until offline mailbox cache becomes a requirement.
- OPFS: rejected for v1 because mailbox data belongs in the backend runtime store.
- Cookies: rejected because the frontend does not authenticate with cookie sessions in v1.
