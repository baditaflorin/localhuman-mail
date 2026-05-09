# 0043 - Domain Vocabulary and UI Language Conventions

## Status

Accepted

## Context

Parser terms like "inline part" and "selector" are not useful to email users.

## Decision

Use mailbox language:

- "Readable body" instead of "inline text part".
- "Attachment" instead of "MIME attachment header".
- "Calendar invite" instead of "text/calendar".
- "Imported with warnings" instead of "partial parse".
- "Searchable text" instead of "indexed body".

## Consequences

- Errors become actionable.
- UI can surface inference without teaching MIME internals.

## Alternatives Considered

- Developer-centric messages: rejected because they make the app feel brittle.

