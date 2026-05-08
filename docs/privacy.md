# Privacy

localhuman-mail has no client analytics in v1.

The GitHub Pages frontend stores only non-sensitive UI settings such as the backend URL in browser `localStorage`.

The frontend must not store:

- mailbox credentials
- mailbox contents
- private keys
- API tokens
- local LLM prompts from private mailboxes

Mailbox contents are imported into the user-controlled backend runtime store. In Docker deployment, that store lives in the `localhuman-data` named volume.

The backend exposes Prometheus metrics without message subjects, senders, recipients, body text, prompts, or search query labels.

Repository: https://github.com/baditaflorin/localhuman-mail

Live app: https://baditaflorin.github.io/localhuman-mail/

