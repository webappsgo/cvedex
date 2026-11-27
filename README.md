# cvedex

[![License](https://img.shields.io/github/license/webappsgo/cvedex)](LICENSE.md)
[![Release](https://img.shields.io/github/v/release/webappsgo/cvedex)](https://github.com/webappsgo/cvedex/releases)
[![CI](https://github.com/webappsgo/cvedex/actions/workflows/ci.yml/badge.svg)](https://github.com/webappsgo/cvedex/actions/workflows/ci.yml)
[![Docs](https://readthedocs.org/projects/casapps-cvedex/badge/?version=latest)](https://casapps-cvedex.readthedocs.io/)

## About

cvedex serves CVE data through an authoritative DNS zone and a read-only HTTP
gateway. The default zone is `cve.`; deployments can choose another zone.

## Features

- Authoritative DNS over UDP and TCP with CVE detail records, zone indexes, DNSSEC signing, ACL-gated transfers, and best-effort NOTIFY
- Verified, schema-versioned compressed datasets with SHA-256 checks and daily refreshes
- HTTP redirects for valid CVE hostnames and server-rendered pages for CVE details, search, catalogs, and service information
- JSON APIs for CVE details and search, vendor and product lookups, KEV data, and BIND-format zone downloads
- Read-only GraphQL and authenticated MCP JSON-RPC endpoints
- Authenticated administration for settings, service status, DNSSEC, audit logs, scheduled tasks, metrics credentials, backups, and network overlays
- Authenticated Prometheus metrics, Grafana dashboard export, and a bounded Loki-compatible log stream
- Managed Tor and I2P services with health reporting and operator controls
- HTTPS-updated IP and domain blocklists with per-source last-good data and health checks
- Guest theme, language, consent, and cross-device preference controls
- Installable progressive web app with versioned static assets and a translated offline fallback
- `cvedex-cli` for CVE lookup, search, local dataset generation, data synchronization, and zone downloads
- Six web interface languages: English, Spanish, French, German, Italian, and Arabic; Arabic uses right-to-left layout

DNS UDP and HTTP rate limiting are enabled by default. DNSSEC private keys are
stored under the configured data directory, and DNS transfers default to
loopback-only access. See [Security](docs/security.md), [DNSSEC](docs/dnssec.md),
and [Administration](docs/admin.md) for operator guidance.

## Production

### Docker Compose

The production Compose file publishes HTTP on `172.17.0.1:64580`, DNS on port
53 over UDP and TCP, and mounts persistent configuration and data directories.
From a source checkout:

```sh
docker compose -f docker/docker-compose.yml up -d
```

The image is also published as `ghcr.io/casapps/cvedex:latest`. Review the
Compose file and its machine port bindings before deployment. See
[Installation](docs/installation.md) and [Self-hosting](docs/self-hosting.md).

### Binary and service installation

Release binaries are published for Linux, macOS, Windows, and FreeBSD on
amd64 and arm64. Download the server binary for your operating system and
architecture from the [GitHub releases](https://github.com/webappsgo/cvedex/releases).
After placing the binary at `/usr/local/bin/cvedex`, install it on supported
service-manager platforms with:

```sh
sudo /usr/local/bin/cvedex --service --install
```

Installation details, including non-root and non-Linux service managers, are
in the [Installation guide](docs/installation.md).

## Client

The companion client is named `cvedex-cli`. Download the matching release
binary for your operating system and architecture from
[GitHub releases](https://github.com/webappsgo/cvedex/releases). For example,
to install the Linux amd64 binary in the current user's local bin directory:

```sh
mkdir -p "$HOME/.local/bin"
curl -q -LSsf -o "$HOME/.local/bin/cvedex-cli" \
  https://github.com/webappsgo/cvedex/releases/latest/download/cvedex-cli-linux-amd64
chmod 755 "$HOME/.local/bin/cvedex-cli"
```

The client can query a server, search CVEs, inspect KEV data, and download a
zone. Public data queries do not require a token:

```sh
cvedex-cli --server https://YOUR_SERVER get CVE-2024-3094
cvedex-cli --server https://YOUR_SERVER search "xz utils" --year 2024
```

See the [CLI reference](docs/cli.md) for configuration, output formats, and
data-maintenance commands. Translation contributors can follow the
[translation guide](docs/translating.md).

## Configuration

The server accepts `--config DIR` for YAML configuration, with environment
variables overriding YAML and command-line flags overriding both. Common
settings include `PORT` or `CVEDEX_PORT` for HTTP, `DNS_ADDRESS`, `DNS_PORT`,
and `DNS_ZONE` for DNS, and `DATA_URL` plus `CHECKSUM_URL` for dataset mirrors.
The server machine's HTTP port is selected once on first run and saved;
containers listen on port 80. See [Configuration](docs/configuration.md) for
defaults and all supported settings.

## API

The public HTTP API is read-only and rooted at `/api/v1`. The server also
provides rendered pages and browser documentation for GraphQL and OpenAPI.
Examples below use a local development server on port 64580:

```sh
curl -q -LSsf -H 'Accept: application/json' \
  http://127.0.0.1:64580/api/v1/cves/CVE-2024-3094
curl -q -LSsf -H 'Accept: application/json' \
  'http://127.0.0.1:64580/api/v1/search?q=linux&limit=10'
curl -q -LSsf -H 'Accept: application/json' \
  http://127.0.0.1:64580/api/v1/kev
```

The canonical health page is `/server/healthz`; it supports HTML, JSON, and
plain-text responses. See the [API reference](docs/api.md) for routes and
response formats.

## Other

### Troubleshooting

- If readiness is unhealthy on first start, wait for the initial dataset download and verification to finish.
- If a refresh fails, cvedex keeps serving the last verified snapshot; check the service logs and configured data URLs.
- DNS queries require UDP and TCP access to the configured DNS port. The default is port 53.
- Admin API requests use the operator Bearer token; browser pages use an authenticated administrator session.

See [Administration](docs/admin.md), [Integrations](docs/integrations.md), and
[Self-hosting](docs/self-hosting.md) for operational details.

## Development

Development instructions are for contributors. All Go builds and tests use
Docker with the `casjaysdev/go:latest` toolchain; do not install Go on the
development machine.

```sh
make build
make test
```

The debug server and CLI binaries are written under `binaries/`; `make clean`
removes generated binaries and release artifacts.
See [Development](docs/development.md) for the integration tiers and test
environment details. Do not commit generated binaries, datasets, or runtime
data.

## Disclaimer

This software is provided "as is" without warranty of any kind. Use it at your
own risk.

- **No Warranty:** The authors are not responsible for damages, data loss, or issues arising from use of this software.
- **Not Professional Advice:** This software does not provide legal, financial, medical, or other professional advice.
- **Third-Party Services:** External data sources and services have their own terms and conditions.
- **Security:** Following security practices cannot guarantee software is free of vulnerabilities.
- **Production Use:** Evaluate the software for your environment before deploying it in production.

By using this software, you acknowledge that you have read and understood this
disclaimer.

## License

MIT; see [LICENSE.md](LICENSE.md).
