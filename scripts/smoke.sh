#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PORT="${LOCALHUMAN_SMOKE_PORT:-18080}"
DATA_DIR="$ROOT/tmp/smoke-data"
LOG_FILE="$ROOT/tmp/smoke-backend.log"

cleanup() {
  if [[ -n "${SERVER_PID:-}" ]]; then
    kill "$SERVER_PID" >/dev/null 2>&1 || true
  fi
}
trap cleanup EXIT

rm -rf "$DATA_DIR"
mkdir -p "$DATA_DIR" "$ROOT/tmp"

LOCALHUMAN_ADDR="127.0.0.1:$PORT" \
LOCALHUMAN_DATA_DIR="$DATA_DIR" \
LOCALHUMAN_ALLOWED_ORIGINS="http://127.0.0.1:4187,http://localhost:4187,https://baditaflorin.github.io" \
LOCALHUMAN_OLLAMA_URL="" \
"$ROOT/bin/localhuman-mail" >"$LOG_FILE" 2>&1 &
SERVER_PID=$!

for _ in {1..40}; do
  if curl -fsS "http://127.0.0.1:$PORT/healthz" >/dev/null 2>&1; then
    break
  fi
  sleep 0.25
done

curl -fsS "http://127.0.0.1:$PORT/readyz" | grep -q '"status":"ok"'
curl -fsS "http://127.0.0.1:$PORT/api/v1/version" | grep -q '"version"'
curl -fsS -X POST "http://127.0.0.1:$PORT/api/v1/import/demo" | grep -q '"imported"'
curl -fsS "http://127.0.0.1:$PORT/api/v1/search?q=launch" | grep -q "Thursday launch checklist"

npm --prefix "$ROOT/frontend" run e2e
