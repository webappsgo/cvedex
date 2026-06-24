# cvedex

[![License](https://img.shields.io/github/license/casapps/cvedex)](LICENSE.md)
[![Release](https://img.shields.io/github/v/release/casapps/cvedex)](https://github.com/casapps/cvedex/releases)
[![CI](https://github.com/casapps/cvedex/actions/workflows/ci.yml/badge.svg)](https://github.com/casapps/cvedex/actions/workflows/ci.yml)

An authoritative DNS server that exposes the entire CVE corpus (~350k records) as queryable DNS names under the `.cve` TLD. Every CVE becomes `YYYY-NNNNN.cve` in DNS.

```
dig TXT cvss.2024-3094.cve
dig TXT severity.2024-3094.cve
dig TXT description.2024-3094.cve
```

## Features

- Authoritative DNS for the `cve.` zone (configurable to any zone)
- Every CVE queryable: CVSS, severity, vector, CWE, dates, KEV flag, EPSS score
- DNSSEC signed (offline by default, NSEC3, ECDSA P-256)
- AXFR/IXFR/NOTIFY for slaving into existing DNS infrastructure
- HTTP redirector: `https://2024-3094.cve` → canonical CVE page
- JSON gateway: per-CVE detail, vendor/product lookups, KEV catalog, zone-file download
- Daily refresh from cvelistV5, CISA KEV, and FIRST.org EPSS
- Single static binary, zero runtime dependencies
- Response Rate Limiting (RRL) enabled by default

## Quick Start

```bash
# Download the latest binary
curl -LSsf -o cvedex https://github.com/casapps/cvedex/releases/latest/download/cvedex-linux-amd64
chmod +x cvedex

# Start the server (first run auto-creates config and downloads data)
./cvedex

# Or with Docker
docker run -d \
  --name cvedex-app \
  -p 172.17.0.1:64580:80 \
  -v ./volumes/config:/config:z \
  -v ./volumes/data:/data:z \
  ghcr.io/casapps/cvedex:latest
```

## Operational Modes

| Mode | Description |
|------|-------------|
| `cvedex` (no args) | Start long-running DNS server |
| `cvedex --mode generate-zone` | One-shot zone file generation |
| `cvedex --mode generate-data` | One-shot upstream data processing |
| `cvedex --mode sync-only` | One-shot data refresh |

## DNS Query Examples

```bash
# Get CVSS score
dig TXT cvss.2024-3094.cve @127.0.0.1

# Get severity
dig TXT severity.2024-3094.cve @127.0.0.1

# Get full description
dig TXT description.2024-3094.cve @127.0.0.1

# Get EPSS score
dig TXT epss.2024-3094.cve @127.0.0.1

# Check if in CISA KEV
dig TXT kev.2024-3094.cve @127.0.0.1

# Get canonical URL (URI record)
dig URI 2024-3094.cve @127.0.0.1

# Slave the zone via AXFR (ACL-gated, localhost only by default)
dig AXFR cve. @127.0.0.1
```

## JSON Gateway

```bash
# Per-CVE detail
curl http://localhost:64580/api/v1/cve/CVE-2024-3094

# CISA KEV catalog
curl http://localhost:64580/api/v1/kev

# Zone file download
curl http://localhost:64580/api/v1/zone

# Health check
curl http://localhost:64580/health
```

## Configuration

Config file: `/etc/casapps/cvedex/server.yml` (root) or `~/.config/casapps/cvedex/server.yml` (user)

Key settings:

| Setting | Default | Description |
|---------|---------|-------------|
| `zone` | `cve.` | DNS zone to serve |
| `port` | `64580` | HTTP/DNS port |
| `dnssec` | `true` | Enable DNSSEC signing |
| `axfr_acl` | `127.0.0.1` | AXFR allowed IPs |
| `redirect_target` | `cve.org` | HTTP redirect upstream |

## Building from Source

```bash
# Development build (to temp dir)
make dev

# Production build (local platform)
make local

# Full cross-platform release build
make build

# Run tests
make test
```

Requires Docker (`casjaysdev/go:latest` — all builds run in Docker, never on host).

## Docker Compose

```bash
# Production
docker compose -f docker/docker-compose.yml up -d

# Development (with debug mode)
docker compose -f docker/docker-compose.dev.yml up
```

## Data Sources

| Source | Data | Update Frequency |
|--------|------|-----------------|
| [cvelistV5](https://github.com/CVEProject/cvelistV5) | CVE records | Daily |
| [CISA KEV](https://www.cisa.gov/known-exploited-vulnerabilities-catalog) | Known exploited vulnerabilities | Daily |
| [FIRST.org EPSS](https://www.first.org/epss/) | Exploit prediction scores | Daily |

## License

[MIT](LICENSE.md) — Copyright (c) 2026 casapps

## Links

- Official site: https://cvedex.casapps.us
- Repository: https://github.com/casapps/cvedex
- Issues: https://github.com/casapps/cvedex/issues
