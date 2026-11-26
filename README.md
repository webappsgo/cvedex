# cvedex

[![License](https://img.shields.io/github/license/webappsgo/cvedex)](LICENSE.md)
[![Release](https://img.shields.io/github/v/release/webappsgo/cvedex)](https://github.com/webappsgo/cvedex/releases)
[![CI](https://github.com/webappsgo/cvedex/actions/workflows/ci.yml/badge.svg)](https://github.com/webappsgo/cvedex/actions/workflows/ci.yml)
[![Documentation](https://readthedocs.org/projects/casapps-cvedex/badge/?version=latest)](https://casapps-cvedex.readthedocs.io/en/latest/)

cvedex serves CVE data through an authoritative DNS zone and a read-only HTTP gateway. The default zone is `cve.`; deployments can choose another zone.

## Current capabilities

- DNS over UDP and TCP, with CVE detail labels, index records, RFC 8482 minimal `ANY`, and `REFUSED` outside the configured zone
- HTTP redirects for valid CVE hostnames to the configured canonical CVE URL, plus a landing page at the zone apex
- Verified schema-versioned gzip dataset loading with SHA-256 validation
- First-run dataset download and daily refresh from a checksummed HTTPS release, preserving the last verified snapshot on refresh errors
- Offline DNSSEC signing with persisted ECDSA P-256 KSK/ZSK keys, NSEC3 denial proofs, and signed zone downloads
- ACL-gated TCP AXFR/IXFR and best-effort NOTIFY for configured secondaries
- JSON APIs for CVE detail and search, vendor and product lookup, KEV data, and BIND-format zone download
- Read-only GraphQL API with a browser query form at `/server/docs/graphql`
- Authenticated Prometheus metrics, a Grafana dashboard export, and a bounded Loki log stream
- Server-rendered public information, CVE detail, and catalog pages
- `cvedex-cli` for CVE lookup, search, dataset generation, data sync, and zone download
- YAML and environment-based server configuration

DNS UDP and HTTP rate limiting are enabled by default. DNSSEC keys are created under the configured data directory and restricted to the service account. Transfers default to loopback-only access. See [Security](docs/security.md), [DNSSEC](docs/dnssec.md), and [Administration](docs/admin.md) for operator guidance.

## Build

All Go builds use Docker with the `casjaysdev/go:latest` toolchain:

```sh
make local
make test
```

This creates `binaries/cvedex` and `binaries/cvedex-cli`.

## Generate a dataset

The scheduled GitHub workflow publishes the generated dataset to the `data-latest` prerelease. To build it yourself, use a local checkout of the CVE Program list; the command fetches CISA KEV and FIRST EPSS snapshots over HTTPS:

```sh
mkdir -p "${TMPDIR:-/tmp}/casapps"
ARTIFACT_DIR=$(mktemp -d "${TMPDIR:-/tmp}/casapps/cvedex-XXXXXX")
cvedex-cli data generate --source-dir ./cvelistV5 --output "$ARTIFACT_DIR/cvedex-data.json.gz"
```

This writes a `.sha256` file next to the dataset. Keep generated data in the temporary artifact directory rather than the source tree. The server downloads the published dataset on first start and refreshes daily. Set `DATA_URL` and `CHECKSUM_URL` to use a mirror.

## Run locally

```sh
mkdir -p "${TMPDIR:-/tmp}/casapps"
DATA_DIR=$(mktemp -d "${TMPDIR:-/tmp}/casapps/cvedex-XXXXXX")
DATA_DIR="$DATA_DIR" CVEDEX_PORT=64580 DNS_ADDRESS=127.0.0.1 DNS_PORT=1053 \
  ./binaries/cvedex --mode development
```

The native HTTP port is selected once from 64000–64999 and saved in `server.yml`; containers listen on port 80. DNS uses port 53 on UDP and TCP. Override `DNS_PORT` for local development as shown.

## DNS examples

```sh
dig TXT cvss.2024-3094.cve @127.0.0.1 -p 1053
dig TXT desc.2024-3094.cve @127.0.0.1 -p 1053
dig TXT kev.2024-3094.cve @127.0.0.1 -p 1053
dig URI 2024-3094.cve @127.0.0.1 -p 1053
```

The CVE prefix is omitted from DNS names. The full record and index formats are documented in [DNS Records](docs/dns-records.md).

## HTTP API

```sh
curl -q -LSsf -H 'Accept: application/json' http://127.0.0.1:64580/api/v1/cves/CVE-2024-3094
curl -q -LSsf -H 'Accept: application/json' 'http://127.0.0.1:64580/api/v1/search?q=linux&limit=10'
curl -q -LSsf -H 'Accept: application/json' http://127.0.0.1:64580/api/v1/kev
```

The status page is `/server/healthz`. See the [API reference](docs/api.md) for all routes and response formats.

## Configuration

The server accepts `--config DIR` for YAML configuration and environment overrides such as `LISTEN`, `PORT`, `DNS_ADDRESS`, `DNS_PORT`, `DNS_ZONE`, `DATA_DIR`, `DATA_URL`, and `CHECKSUM_URL`. See [Configuration](docs/configuration.md).

## Documentation

Browse the [documentation site](https://casapps-cvedex.readthedocs.io) for [installation](docs/installation.md), [self-hosting](docs/self-hosting.md), [security](docs/security.md), and [development](docs/development.md) guidance.
