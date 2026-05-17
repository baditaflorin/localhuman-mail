FROM golang:1.26.3-alpine AS builder

RUN apk add --no-cache ca-certificates git
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION=0.3.1
ARG COMMIT=dev
ARG BUILD_TIME=unknown

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -trimpath \
    -ldflags="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.buildTime=${BUILD_TIME}" \
    -o /out/localhuman-mail ./cmd/server

# Pre-create /data so the named volume initializes nonroot-owned
# (ADR-0023 gap 3 — Docker copies image content + ownership to the
# named volume on first mount; without this, /data is root:root and
# the nonroot user cannot open the SQLite DB).
RUN mkdir -p /out/data

FROM gcr.io/distroless/static-debian12:nonroot

ARG VERSION=0.3.1
ARG COMMIT=dev
ARG BUILD_TIME=unknown

LABEL org.opencontainers.image.title="localhuman-mail" \
      org.opencontainers.image.description="Privacy-first AI email client backend." \
      org.opencontainers.image.source="https://github.com/baditaflorin/localhuman-mail" \
      org.opencontainers.image.revision="${COMMIT}" \
      org.opencontainers.image.version="${VERSION}" \
      org.opencontainers.image.created="${BUILD_TIME}" \
      org.opencontainers.image.licenses="MIT"

COPY --from=builder /out/localhuman-mail /localhuman-mail
COPY --from=builder --chown=nonroot:nonroot /out/data /data

USER nonroot:nonroot
EXPOSE 8080
VOLUME ["/data"]

HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 CMD ["/localhuman-mail", "-healthcheck"]
ENTRYPOINT ["/localhuman-mail"]
