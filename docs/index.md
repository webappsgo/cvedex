# cvedex

cvedex makes CVE records queryable through authoritative DNS and a read-only HTTP API. The default zone is `cve.`; the configured zone may be a delegated domain or an internally routed custom zone.

The service loads a verified compressed dataset generated from CVE Program v5 records, CISA KEV, and FIRST EPSS data. It fetches the initial snapshot from the configured HTTPS distribution and refreshes it daily, keeping the last verified snapshot active when a refresh fails.

## Try a query

```sh
dig TXT severity.2024-3094.cve @127.0.0.1
curl -H 'Accept: application/json' http://127.0.0.1:64580/api/v1/cves/CVE-2024-3094
```

## Interfaces

- DNS over UDP and TCP for per-CVE fields and zone indexes
- HTTP JSON APIs for CVE lookup, search, reverse vendor/product lookups, KEV data, and zone download
- Server-rendered pages for CVEs, search, catalogs, and service information
- `cvedex-cli` for data generation, lookup, search, sync, and zone download

DNS UDP and HTTP rate limiting are enabled by default. First-run admin
authentication and the settings panel are available. DNSSEC signing, ACL
restricted transfers, and NOTIFY are supported; see [Security](security.md)
and [Administration](admin.md) for operator guidance.

## Documentation

- [Install cvedex](installation.md)
- [Configure the server](configuration.md)
- [DNS record reference](dns-records.md)
- [HTTP API](api.md)
- [CLI reference](cli.md)
- [Self-hosting](self-hosting.md)
- [Development](development.md)
