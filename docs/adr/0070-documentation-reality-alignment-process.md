# 0070 - Documentation-Reality Alignment Process

## Status

Accepted

## Context

The README and ADRs included a few broader claims than the UI currently fulfilled.

## Decision

README features become a verified checklist tied to tests or explicit limitations. Claims about optional native tools must say "capability detection" unless the product exposes an execution path. The live Pages URL, repo URL, PayPal URL, version, and commit remain first-class documented deliverables.

## Consequences

- Overpromising is fixed by shipping the path or editing the claim.
- Limitations are not buried in postmortems only.

## Alternatives Considered

- Leave aspirational text in docs: rejected because Phase 3 treats docs drift as a product bug.
