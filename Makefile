SHELL := /bin/bash

APP_NAME := localhuman-mail
REPO := ghcr.io/baditaflorin/$(APP_NAME)
VERSION := 0.1.0
COMMIT := $(shell git rev-parse --short=12 HEAD 2>/dev/null || echo dev)
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GO_PACKAGES := ./cmd/... ./internal/...
GOTMP := $(CURDIR)/tmp/go-build
GO_ENV := GOTMPDIR=$(GOTMP) CGO_ENABLED=0
DOCKER_BUILDER ?=
DOCKER_BUILDER_FLAG := $(if $(DOCKER_BUILDER),--builder $(DOCKER_BUILDER),)

.PHONY: help install-hooks generate dev build data test test-integration smoke lint fmt pages-preview docker-build docker-push release compose-up compose-down clean hooks-pre-commit hooks-commit-msg hooks-pre-push hooks-post-merge

help:
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z0-9_-]+:.*##/ {printf "%-22s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install-hooks: ## Wire .githooks into this clone
	git config core.hooksPath .githooks
	chmod +x .githooks/*

generate: ## Regenerate generated frontend API types
	npm --prefix frontend run generate:api

dev: ## Run backend and frontend locally
	mkdir -p tmp runtime
	(LOCALHUMAN_ADDR=:8080 LOCALHUMAN_DATA_DIR=./runtime $(GO_ENV) go run ./cmd/server) & \
	npm --prefix frontend run dev

build: generate ## Build backend binary and Pages-ready frontend
	mkdir -p bin $(GOTMP)
	npm --prefix frontend run build
	$(GO_ENV) go build -trimpath -ldflags="-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildTime=$(BUILD_TIME)" -o bin/$(APP_NAME) ./cmd/server
	test -f docs/index.html
	test -s docs/index.html

data: ## No-op in Mode C
	@echo "Mode C stores private mailbox data in the backend runtime store; no static data artifacts."

test: ## Run unit tests
	mkdir -p $(GOTMP)
	npm --prefix frontend run test
	$(GO_ENV) go test $(GO_PACKAGES)

test-integration: ## Run integration tests
	@echo "No integration fixtures are committed yet."

smoke: build ## Run Pages and backend smoke tests
	scripts/smoke.sh

lint: ## Run linters and type checks
	mkdir -p $(GOTMP)
	npm --prefix frontend run lint
	npm --prefix frontend run fmt:check
	npm --prefix frontend run typecheck
	$(GO_ENV) go vet $(GO_PACKAGES)
	@if command -v golangci-lint >/dev/null 2>&1; then golangci-lint run ./cmd/... ./internal/...; else echo "golangci-lint not installed; skipping optional lint"; fi

fmt: ## Autoformat source
	gofmt -w cmd internal
	npm --prefix frontend run fmt

pages-preview: build ## Serve the Pages build locally
	npm --prefix frontend run preview -- --port 4173

docker-build: ## Build linux/amd64 backend image
	docker buildx build $(DOCKER_BUILDER_FLAG) --platform linux/amd64 \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		-t $(REPO):latest \
		-t $(REPO):$(VERSION) \
		-t $(REPO):$(COMMIT) \
		--load .

docker-push: ## Push backend image to GHCR
	docker buildx build $(DOCKER_BUILDER_FLAG) --platform linux/amd64 \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		-t $(REPO):latest \
		-t $(REPO):$(VERSION) \
		-t $(REPO):$(COMMIT) \
		--push .

release: test build docker-push ## Tag and push a release
	git tag v$(VERSION)
	git push origin v$(VERSION)

compose-up: ## Run local Docker Compose stack
	docker compose -f deploy/docker-compose.yml -f deploy/docker-compose.dev.yml up -d --build

compose-down: ## Stop local Docker Compose stack
	docker compose -f deploy/docker-compose.yml -f deploy/docker-compose.dev.yml down

clean: ## Remove local build outputs
	rm -rf bin tmp runtime frontend/playwright-report frontend/test-results

hooks-pre-commit:
	.githooks/pre-commit

hooks-commit-msg:
	@test -n "$(MSG)" || (echo "usage: make hooks-commit-msg MSG=.git/COMMIT_EDITMSG" && exit 2)
	.githooks/commit-msg "$(MSG)"

hooks-pre-push:
	.githooks/pre-push

hooks-post-merge:
	.githooks/post-merge
