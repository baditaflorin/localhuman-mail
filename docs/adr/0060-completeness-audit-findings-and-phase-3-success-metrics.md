# 0060 - Completeness Audit Findings And Phase 3 Success Metrics

## Status

Accepted

## Context

Phase 3 asks whether a stranger can use the live app end-to-end on their own data. The audits in `docs/phase3/` found that the core parser is now substantive, but the product path is incomplete around loading multiple messages, exporting work, restoring state, and matching docs to shipped behavior.

## Decision

Use `docs/phase3/findings.md` and `docs/phase3/plan.md` as the Phase 3 rubric. Success requires all non-gray input/output/control rows to be green or documented as permanently out of scope, zero unsafe TypeScript casts in authored frontend code, e2e coverage for upload/export/copy/state/share paths, and a passing full local hook chain.

## Consequences

- Work that does not move an audit row or success metric waits.
- Documentation drift is treated as a product bug.
- Pages and Docker deployment mode remain Mode C.

## Alternatives Considered

- Polish-first UI pass: rejected because output/input holes would remain.
- New engine work: rejected because Phase 2 already set the parser floor.
