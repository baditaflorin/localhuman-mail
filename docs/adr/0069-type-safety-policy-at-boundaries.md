# 0069 - Type-Safety Policy At Boundaries

## Status

Accepted

## Context

OpenAPI types provide API safety, but multipart and imported JSON were using or would invite unsafe casts.

## Decision

Authored frontend code must not contain `any`, `// @ts-ignore`, `as never`, or unsafe assertions. Boundary data is validated with zod or narrowed with explicit type guards. Generated files are excluded from this authored-code rule.

## Consequences

- Multipart upload goes through a typed wrapper.
- State/share imports are validated before use.
- `rg` checks are part of the final audit.

## Alternatives Considered

- Allow casts in UI handlers: rejected because these are exactly where user data enters.
