# Phase 2 State Taxonomy

## Import States

- `idle`: no import in progress.
- `validating`: upload received and size/type checks running.
- `parsing`: MIME structure is being parsed.
- `recovering`: parser found malformed input but is attempting partial recovery.
- `imported-clean`: message imported with high confidence and no warnings.
- `imported-with-warnings`: message imported, but warnings require review.
- `rejected-recoverable`: input cannot be imported, but mailbox state is intact and the user has a next step.
- `rejected-fatal`: server/runtime failure prevented import; mailbox state is intact.

## Message List States

- `loaded-empty`: backend connected, no messages.
- `loaded-some`: one or more messages.
- `loaded-many`: more than 100 messages, list is capped by API limit.
- `searching`: search request in flight.
- `search-empty`: query completed with no matches.
- `search-error-recoverable`: search failed but prior list remains visible.

## Draft States

- `draft-idle`: no draft requested.
- `drafting-local-llm`: local model request in progress.
- `drafting-fallback`: deterministic fallback draft used.
- `draft-ready`: draft text is available.
- `draft-error-recoverable`: draft failed, message remains visible.

Every state has an exit: retry, clear search, select another message, import another message, or keep existing mailbox state.

