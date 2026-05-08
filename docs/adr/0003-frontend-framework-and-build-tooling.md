# 0003 - Frontend Framework and Build Tooling

## Status

Accepted

## Context

The frontend needs a fast, typed, accessible product UI with simple static deployment to GitHub Pages.

## Decision

Use React, TypeScript strict mode, Vite, Tailwind CSS, TanStack Query, Zod, Lucide React, and Playwright.

## Consequences

- Vite can build directly to `docs/`.
- React supports the interaction density expected from an email client.
- TanStack Query handles API caching and stale-state behavior.
- Zod validates API payloads at the frontend boundary.
- Tailwind keeps styling local to the product without a large CSS framework.

## Alternatives Considered

- Vanilla TypeScript: rejected because complex client state would become harder to maintain.
- Next.js: rejected because server rendering is unnecessary and complicates Pages deployment.
- Svelte: viable, but React has broader library support for the chosen stack.

