# 0005 - Client-Side Storage Strategy

## Status

Accepted

## Context

The frontend needs to remember non-sensitive preferences such as backend URL, layout density, and last selected message. It must not store mailbox credentials or mailbox contents.

## Decision

Use `localStorage` for small non-sensitive UI settings in v1. Do not store mailbox contents, tokens, credentials, private keys, or AI prompts in the browser.

## Consequences

- Users can reconnect to their local backend without retyping the URL.
- Clearing browser storage removes all client-side state.
- Cross-device sync is explicitly out of scope.

## Alternatives Considered

- IndexedDB: deferred until offline mailbox cache becomes a requirement.
- OPFS: rejected for v1 because mailbox data belongs in the backend runtime store.
- Cookies: rejected because the frontend does not authenticate with cookie sessions in v1.

