# =============================================================================
# Build Stage - Compile Go binary
# =============================================================================
FROM casjaysdev/go:latest AS builder

RUN apk add --no-cache git bash

ARG TARGETARCH
ARG VERSION=dev
ARG BUILD_DATE
ARG COMMIT_ID
ARG OFFICIAL_SITE

WORKDIR /app

# Copy go.mod first for layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build \
    -ldflags "-s -w \
      -X 'main.Version=${VERSION}' \
      -X 'main.CommitID=${COMMIT_ID}' \
      -X 'main.BuildDate=${BUILD_DATE}' \
      -X 'main.OfficialSite=${OFFICIAL_SITE}'" \
    -o /app/binary/cvedex ./src

# =============================================================================
# Runtime Stage - Minimal Alpine image
# =============================================================================
FROM alpine:latest

ARG VERSION=dev
ARG BUILD_DATE
ARG COMMIT_ID
ARG LICENSE=MIT

# No LABEL blocks — all image metadata is applied as OCI annotations at build time.
# Pass --annotation flags to docker buildx build (or use docker/metadata-action
# annotations: output in GitHub Actions). Annotations attach to the manifest index
# and are visible on multiarch pulls. LABEL applies only to per-platform layers.

RUN apk add --no-cache \
    git \
    curl \
    bash \
    tini \
    tor

# Copy binary from builder stage (multi-stage build)
COPY --from=builder /app/binary/cvedex /usr/local/bin/cvedex

# Copy build-time overlay (entrypoint.sh) from docker/rootfs/ into image
COPY docker/rootfs/ /

# Copy Dockerfile to image for reference
COPY docker/Dockerfile /root/Dockerfile

# Make all scripts in /usr/local/bin executable
RUN chmod 755 /usr/local/bin/*

# MODE=development allows localhost, .local, .test hostnames in the container
# DATABASE_DIR: SQLite location (explicit is clearer than auto-detection)
ENV MODE=development \
    TZ=America/New_York \
    DATABASE_DIR=/data/db/sqlite

EXPOSE 80

STOPSIGNAL SIGRTMIN+3

HEALTHCHECK --start-period=10m --interval=5m --timeout=15s --retries=3 \
    CMD /usr/local/bin/cvedex --status || exit 1

ENTRYPOINT [ "tini", "-p", "SIGTERM", "--", "/usr/local/bin/entrypoint.sh" ]
