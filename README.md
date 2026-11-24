# cvedex

[![License](https://img.shields.io/github/license/webappsgo/cvedex)](LICENSE.md)
[![Release](https://img.shields.io/github/v/release/webappsgo/cvedex)](https://github.com/webappsgo/cvedex/releases)
[![CI](https://github.com/webappsgo/cvedex/actions/workflows/ci.yml/badge.svg)](https://github.com/webappsgo/cvedex/actions/workflows/ci.yml)

cvedex serves CVE data through an authoritative DNS zone and a read-only HTTP gateway. The default zone is `cve.`; deployments can choose another zone.

## Current capabilities

- DNS over UDP and TCP, with CVE detail labels, index records, RFC 8482 minimal `ANY`, and `REFUSED` outside the configured zone
- Verified schema-versioned gzip dataset loading with SHA-256 validation
- First-run dataset download and daily refresh from a checksummed HTTPS release, preserving the last verified snapshot on refresh errors
- JSON APIs for CVE detail and search, vendor and product lookup, KEV data, and BIND-format zone download
- Server-rendered public information, CVE detail, and catalog pages
- `cvedex-cli` for CVE lookup, search, dataset generation, data sync, and zone download
- YAML and environment-based server configuration

DNSSEC signing, zone transfers, rate limiting, and the authenticated admin panel remain required implementation work. See [Security](docs/security.md) and [Administration](docs/admin.md) for current limitations.

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
cvedex-cli data generate --source-dir ./cvelistV5 --output ./cvedex-data.json.gz
```

This writes a `.sha256` file next to the dataset. The server downloads the published dataset on first start and refreshes daily. Set `DATA_URL` and `CHECKSUM_URL` to use a mirror.

## Run locally

```sh
mkdir -p "${TMPDIR:-/tmp}/casapps"
DATA_DIR=$(mktemp -d "${TMPDIR:-/tmp}/casapps/cvedex-XXXXXX")
cp ./cvedex-data.json.gz* "$DATA_DIR/"
DATA_DIR="$DATA_DIR" DNS_ADDRESS=127.0.0.1 DNS_PORT=1053 \
  ./binaries/cvedex --mode development
```

The default HTTP port is 64580. The default DNS port is 53 and usually requires elevated bind capability. Override it for local development as shown.

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
curl -H 'Accept: application/json' http://127.0.0.1:64580/api/v1/cves/CVE-2024-3094
curl -H 'Accept: application/json' 'http://127.0.0.1:64580/api/v1/search?q=linux&limit=10'
curl -H 'Accept: application/json' http://127.0.0.1:64580/api/v1/kev
```

The status page is `/server/healthz`. See the [API reference](docs/api.md) for all routes and response formats.

## Configuration

The server accepts `--config PATH` for YAML configuration and environment overrides such as `ADDRESS`, `PORT`, `DNS_ADDRESS`, `DNS_PORT`, `DNS_ZONE`, and `DATA_DIR`. See [Configuration](docs/configuration.md).

## Documentation

Browse the [documentation site](https://casapps-cvedex.readthedocs.io) for [installation](docs/installation.md), [self-hosting](docs/self-hosting.md), [security](docs/security.md), and [development](docs/development.md) guidance.
