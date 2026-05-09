FROM golang:1.25 AS builder # Build stage with Go toolchain.

WORKDIR /src # Working directory inside builder.

COPY go.mod go.sum ./ # Copy deps manifests first for cache.
RUN --mount=type=cache,target=/go/pkg/mod \
	go mod download # Cache Go modules.

COPY . . # Copy the rest of the source.

ARG APP_NAME=mmo_affiliate_app # App name for ldflags.
ARG APP_VERSION=0.0.0 # App version for ldflags.
ARG APP_ENV=prod # Build-time env injected into binary.

RUN --mount=type=cache,target=/go/pkg/mod \
	--mount=type=cache,target=/root/.cache/go-build \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -ldflags "-s -w -X main.Name=${APP_NAME} -X main.Version=${APP_VERSION} -X main.Env=${APP_ENV}" \
	-o /out/app ./cmd/app # Build a static binary.

FROM alpine:3.20 # Minimal runtime image.

RUN apk add --no-cache ca-certificates # TLS certs for HTTPS.

WORKDIR /app # Runtime working directory.

COPY --from=builder /out/app /app/app # Copy built binary.
COPY config /app/config # Copy config files.

EXPOSE 1424 # HTTP port.

ENTRYPOINT ["/app/app", "-env", "prod"] # Run with prod config.
