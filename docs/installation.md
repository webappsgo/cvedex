# Installation

cvedex runs as a Go server with an HTTP listener and authoritative DNS listeners on both UDP and TCP. The default DNS zone is `cve.` and the default DNS port is 53.

## Docker Compose

The repository includes Compose files under `docker/`. The production image must publish DNS port 53 for both transports and the configured HTTP port. Mount persistent storage at `/data` and configuration at `/config`.

```sh
docker compose -f docker/docker-compose.yml up -d
```

The current Compose manifest is being brought into sync with the DNS service. Until that is complete, review its published ports before using it.

## Build from source

All Go builds use the project's `casjaysdev/go:latest` container toolchain. From the repository root:

```sh
make local
```

The generated server binary is placed under `binaries/`.

## Runtime data

The server loads a verified dataset snapshot from `<data_dir>/cvedex-data.json.gz` and its SHA-256 file at `<data_dir>/cvedex-data.json.gz.sha256`. A missing snapshot leaves the CVE data readiness check unhealthy. Dataset generation and scheduled publication are documented in [Development](development.md).

## Ports and privileges

The default HTTP port is 64580. DNS uses port 53 on UDP and TCP, which usually requires elevated bind capability on Unix systems. Set `DNS_PORT` for an unprivileged development instance.
