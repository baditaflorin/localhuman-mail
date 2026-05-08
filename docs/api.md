# API

OpenAPI spec: api/openapi.yaml

Default local backend: http://localhost:8080

## Health

```bash
curl http://localhost:8080/healthz
curl http://localhost:8080/readyz
```

## Version

```bash
curl http://localhost:8080/api/v1/version
```

## Import Demo Messages

```bash
curl -X POST http://localhost:8080/api/v1/import/demo
```

## Import EML

```bash
curl -F "file=@message.eml" http://localhost:8080/api/v1/import/eml
```

## Search

```bash
curl "http://localhost:8080/api/v1/search?q=launch&limit=10"
```

## Draft Reply

```bash
curl -X POST http://localhost:8080/api/v1/assist/reply \
  -H "Content-Type: application/json" \
  -d '{"messageId":"demo-backend-1","tone":"concise"}'
```

## Metrics

```bash
curl http://localhost:8080/metrics
```

Production nginx blocks public `/metrics`; Prometheus scrapes it over the internal Docker network.

