# cvedex

An authoritative DNS server that exposes the entire CVE corpus (~350k records) as queryable DNS names under the `.cve` TLD.

## Quick Start

```bash
dig TXT cvss.2024-3094.cve @127.0.0.1
dig TXT severity.2024-3094.cve @127.0.0.1
dig TXT kev.2024-3094.cve @127.0.0.1
```

## Features

- Authoritative DNS for the `cve.` zone (configurable)
- Every CVE queryable via standard DNS tools
- DNSSEC signed (NSEC3, ECDSA P-256)
- AXFR/IXFR/NOTIFY support
- HTTP redirector and JSON gateway
- Daily refresh from canonical sources
- Single static binary, zero runtime dependencies

## Quick Links

- [Getting Started](getting-started.md)
- [DNS Records Reference](dns-records.md)
- [JSON API Reference](api.md)
- [DNSSEC Setup](dnssec.md)
