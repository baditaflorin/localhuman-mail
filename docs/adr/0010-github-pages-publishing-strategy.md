# 0010 - GitHub Pages Publishing Strategy

## Status

Accepted

## Context

The live GitHub Pages URL is a first-class deliverable from commit one. No GitHub Actions are allowed, so publishing must happen from committed build output.

## Decision

Publish from `main` branch `/docs` folder at https://baditaflorin.github.io/localhuman-mail/.

The Vite frontend builds into `docs/`. The folder is committed and intentionally not ignored. Asset filenames are hashed by Vite. A `docs/404.html` SPA fallback is generated after each build.

## Consequences

- The repo contains built frontend artifacts.
- `make build` must refresh `docs/` before pushing frontend changes.
- Rollback is a normal git revert of the publishing commit.
- Custom domain support can be added later with `docs/CNAME`.

## Alternatives Considered

- `gh-pages` branch: rejected because no GitHub Actions means branch publishing adds extra local steps.
- `main /`: rejected to keep source files and published files separate.
- `dist/`: rejected because `dist/` is conventionally ignored and not directly selectable as a Pages source unless committed.

