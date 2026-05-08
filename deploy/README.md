# Deployment

Frontend URL: https://baditaflorin.github.io/localhuman-mail/

Backend image: ghcr.io/baditaflorin/localhuman-mail:latest

Repository: https://github.com/baditaflorin/localhuman-mail

## Server Prerequisites

- Docker Engine with Compose plugin
- DNS record pointing your backend hostname to the server
- Let's Encrypt certificates mounted at `/etc/letsencrypt`
- Optional local Ollama endpoint for reply drafting

## First Run

```bash
cd deploy
cp .env.example .env
docker compose pull
docker compose up -d
```

The backend is available through nginx on host port `25342`.

## TLS

Before starting nginx, replace `your-domain.example` in `deploy/nginx/conf.d/localhuman-mail.conf` with the real certificate directory under `/etc/letsencrypt/live/`.

## Metrics

Prometheus is profile-gated:

```bash
docker compose --profile metrics up -d prometheus
```

Public `/metrics` access is blocked by nginx. Prometheus scrapes the app over the internal bridge network.

## Rollback

```bash
docker compose pull app
docker compose up -d app
```

For a specific version:

```bash
docker compose pull ghcr.io/baditaflorin/localhuman-mail:v0.1.0
```

Then update the image tag in `deploy/docker-compose.yml` and restart.

## Backups

Runtime state lives in the `localhuman-data` named volume. Back it up with:

```bash
docker run --rm -v localhuman-mail_localhuman-data:/data -v "$PWD":/backup alpine tar czf /backup/localhuman-mail-data.tgz /data
```

## Logs

```bash
docker compose logs -f app
docker compose logs -f nginx
```

## Resource Sizing

Minimum: 1 CPU, 1GB RAM, 5GB disk.

Recommended with local LLM nearby: 2 CPU, 2GB RAM for the API container, plus whatever Ollama requires for the selected model.

