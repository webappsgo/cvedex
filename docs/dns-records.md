# DNS records

The default authoritative zone is `cve.`. The CVE prefix is omitted from the owner name: `CVE-2024-3094` is `2024-3094.cve.`. Each query is answered authoritatively over UDP and TCP.

| Name | Type | Data |
|---|---|---|
| `2024-3094.cve.` | URI | Canonical CVE Program URL |
| `2024-3094.cve.` | A / AAAA | Redirector address when configured; no address is currently configured by default |
| `2024-3094.cve.` | TXT | Identifier, severity, score, description summary |
| `desc.2024-3094.cve.` | TXT | Description, split into DNS character strings no longer than 255 bytes |
| `cvss.2024-3094.cve.` | TXT | Highest available CVSS score |
| `severity.2024-3094.cve.` | TXT | `none`, `low`, `medium`, `high`, or `critical` |
| `vector.2024-3094.cve.` | TXT | CVSS vector |
| `cwe.2024-3094.cve.` | TXT | Comma-separated CWE identifiers |
| `published.2024-3094.cve.` / `modified.2024-3094.cve.` | TXT | Source timestamps |
| `state.2024-3094.cve.` | TXT | CVE record state |
| `vendors.2024-3094.cve.` / `products.2024-3094.cve.` | TXT | Associated vendor or product names |
| `refs.2024-3094.cve.` | TXT | One answer per reference URL |
| `kev.2024-3094.cve.` | TXT | `no` or KEV inclusion and metadata |
| `epss.2024-3094.cve.` | TXT | EPSS score and percentile, or `unavailable` |

Zone indexes include `count.cve.`, `count.kev.cve.`, `latest.cve.`, `built.cve.`, `source.cve.`, and `version.cve.`. Per-year and per-severity counters are available under `count.year-{year}.cve.` and `count.severity-{severity}.cve.`.

The URI record and HTTP redirector use the same `REDIRECT_URL_TEMPLATE`, with `{year}`, `{sequence}`, and `{cve_id}` substitutions. DNS UDP responses are rate-limited by client by default.

The apex serves SOA with a serial based on the last successful dataset build.
Unknown names return NXDOMAIN with SOA authority; existing names queried for
an absent type return NOERROR with SOA authority. Out-of-zone requests return
REFUSED. `ANY` queries receive a minimal RFC 8482 response. When DNSSEC is
enabled, the zone includes DNSKEY, RRSIG, NSEC3, and denial proof records. TCP
AXFR/IXFR is restricted by `DNS_TRANSFER_ACL`; successful refreshes send
best-effort NOTIFY to configured `DNS_NOTIFY_TARGETS`.
