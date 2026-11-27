# cvedex

cvedex makes CVE records queryable through authoritative DNS and a read-only HTTP API. The default zone is `cve.`; the configured zone may be a delegated domain or an internally routed custom zone.

HTTP requests to the zone apex show the landing page. Requests whose hostname is a
valid CVE name redirect to the configured upstream CVE URL when that record is
loaded.

The service loads a verified compressed dataset generated from CVE Program v5 records, CISA KEV, and FIRST EPSS data. It fetches the initial snapshot from the configured HTTPS distribution and refreshes it daily, keeping the last verified snapshot active when a refresh fails.

## Try a query locally

These examples assume a local development server using the default DNS zone.

```sh
dig TXT severity.2024-3094.cve @127.0.0.1
curl -q -LSsf -H 'Accept: application/json' http://127.0.0.1:64580/api/v1/cves/CVE-2024-3094
```

## Interfaces

- DNS over UDP and TCP for per-CVE fields and zone indexes
- HTTP redirects from loaded CVE hostnames to their canonical upstream pages
- HTTP JSON APIs for CVE lookup, search, reverse vendor/product lookups, KEV data, and zone download
- Read-only GraphQL with a browser explorer and OpenAPI documentation with a browser reference
- Authenticated MCP JSON-RPC transport for operator and metrics-scoped integrations
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
- [Administration](admin.md)
- [Security](security.md)
- [Integrations](integrations.md)
- [Self-hosting](self-hosting.md)
- [Development](development.md)
- [Translating](translating.md)
