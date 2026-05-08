# Runbook

Live frontend: https://baditaflorin.github.io/localhuman-mail/

Repository: https://github.com/baditaflorin/localhuman-mail

## Local Debugging

```bash
make dev
curl http://localhost:8080/healthz
curl -X POST http://localhost:8080/api/v1/import/demo
curl "http://localhost:8080/api/v1/search?q=launch"
```

## Logs

Local backend logs are JSON on stdout.

Docker logs:

```bash
cd deploy
docker compose logs -f app
docker compose logs -f nginx
```

## Metrics

```bash
cd deploy
docker compose --profile metrics up -d prometheus
```

Prometheus config: deploy/prometheus/prometheus.yml

Grafana starter dashboard: deploy/grafana-dashboard.json

## Common Failures

- Frontend shows demo mode: backend is not reachable from the browser, or CORS does not include the Pages origin.
- EML import returns 400: upload is not multipart form data, or the message has no readable `text/plain` body.
- Reply assist uses fallback: Ollama is not reachable or the configured model is not available.
- `/metrics` works locally but not publicly: expected in production; nginx blocks it.
- `go test ./...` enters `frontend/node_modules`: use `make test`, which targets `./cmd/... ./internal/...`.

## Resource Sizing

Minimum backend: 1 CPU, 1GB RAM, 5GB disk.

Recommended backend: 2 CPU, 2GB RAM, 10GB disk.

Local LLM sizing depends on the Ollama model and runs outside this API container.

## Backups

Runtime data is in the Docker named volume `localhuman-data`.

```bash
cd deploy
docker run --rm -v localhuman-mail_localhuman-data:/data -v "$PWD":/backup alpine tar czf /backup/localhuman-mail-data.tgz /data
```

## Escalation

Security policy: SECURITY.md

Public issue tracker: https://github.com/baditaflorin/localhuman-mail/issues

