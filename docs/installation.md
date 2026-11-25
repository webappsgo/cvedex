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

## Binary and systemd

Run the server binary directly to start it in the foreground. To install it as
a systemd service on Linux, first place the binary at its permanent path, then
run the following as root:

```sh
sudo /usr/local/bin/cvedex --service --install
```

The installer creates the `cvedex` system account when needed, writes and
enables the service unit, and starts the service. The generated system unit
runs unprivileged and grants only the `CAP_NET_BIND_SERVICE` capability needed
for DNS port 53. A non-root invocation installs a per-user unit instead; use
`systemctl --user status cvedex` to check it. Service management supports
`--service start`, `stop`, `restart`, `reload`, `--disable`, and `--uninstall`.
On a host without systemd, run the binary under the host's process supervisor.

## Runtime data

The server downloads and verifies the dataset at first start, then refreshes daily by default. A verified snapshot and checksum are persisted under `<data_dir>`. If distribution is unavailable, a cached snapshot continues serving; without any valid snapshot, readiness stays unhealthy. Dataset generation and scheduled publication are documented in [Development](development.md).

## Ports and privileges

The native HTTP port is selected on first run from 64000–64999 and persisted in
the configuration file; containers listen on port 80. DNS uses port 53 on UDP
and TCP. Set `CVEDEX_PORT` (or `PORT` in a container) to choose an HTTP port,
and `DNS_PORT` to choose an unprivileged DNS port.
