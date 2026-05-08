# GitHub Pages Deployment

Live app: https://baditaflorin.github.io/localhuman-mail/

Repository: https://github.com/baditaflorin/localhuman-mail

GitHub Pages is configured from `main` branch `/docs`.

## Publish

```bash
make build
git add docs frontend
git commit -m "feat: update pages app"
git push
```

## Rollback

Revert the commit that changed `docs/`, then push `main`.

```bash
git revert <commit>
git push
```

## Pages Gotchas

- GitHub Pages does not support `_headers` or `_redirects`.
- SPA fallback is handled by copying `docs/index.html` to `docs/404.html`.
- Vite `base` is `/localhuman-mail/`.
- The service worker scope is `/localhuman-mail/`.
- Do not gitignore `docs/`; it is the publish directory.

