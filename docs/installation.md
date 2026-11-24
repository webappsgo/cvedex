# Installation

cvedex runs as a Go server with an HTTP listener and authoritative DNS listeners on both UDP and TCP. The default DNS zone is `cve.` and the default DNS port is 53.

## Docker Compose

The repository includes Compose files under `docker/`. The production image must publish DNS port 53 for both transports and the configured HTTP port. Mount persistent storage at `/data` and configuration at `/config`.

```sh
docker compose -f docker/docker-compose.yml up -d
```

The Compose manifests publish HTTP plus DNS over both UDP and TCP. Persistent `/data` storage holds the verified dataset and SQLite databases; `/config` holds operator configuration.

## Build from source

All Go builds use the project's `casjaysdev/go:latest` container toolchain. From the repository root:

```sh
make local
```

The generated server binary is placed under `binaries/`.

## Runtime data

The server downloads and verifies the dataset at first start, then refreshes daily by default. A verified snapshot and checksum are persisted under `<data_dir>`. If distribution is unavailable, a cached snapshot continues serving; without any valid snapshot, readiness stays unhealthy. Dataset generation and scheduled publication are documented in [Development](development.md).

## Ports and privileges

The default HTTP port is 64580. DNS uses port 53 on UDP and TCP, which usually requires elevated bind capability on Unix systems. Set `DNS_PORT` for an unprivileged development instance.
