# 0044 - Confidence Model

## Status

Accepted

## Context

v1 either succeeds or fails. Real inputs need confidence and warnings so low-trust output is not presented as fact.

## Decision

Represent confidence as `{score, label, reasons}` for message-level and field-level inference. Labels are:

- `high`: score >= 0.8
- `medium`: score >= 0.5
- `low`: score < 0.5

Warnings carry severity, field, message, and next step.

## Consequences

- UI can expose low-confidence fields.
- Fixture expectations can assert confidence properties.
- Exports/API can carry confidence downstream.

## Alternatives Considered

- Boolean valid/invalid: rejected because many email imports are partially recoverable.

