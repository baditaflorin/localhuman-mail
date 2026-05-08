# Contributing

Thanks for helping improve localhuman-mail.

## Local Workflow

1. Install Go, Node.js, Docker, gitleaks, and lefthook-compatible shell tooling.
2. Run `make install-hooks`.
3. Create a focused branch.
4. Use Conventional Commits for commit messages.
5. Run `make test`, `make build`, and `make smoke` before pushing.

Do not commit secrets, mailbox contents, private keys, production `.env` files, or generated local runtime data.

