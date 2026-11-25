# Integrations

## DNS clients

The authoritative DNS listener accepts UDP and TCP queries. The default zone is
`cve.`; set `DNS_ZONE` and `DNS_PORT` to use another zone or port. In the
examples below, replace `127.0.0.1` with the server address and add
`-p PORT` when using a non-default DNS port.

### CVE record cookbook

The following queries use `CVE-2024-3094` as an example. The identifier's
`CVE-` prefix is omitted from DNS owner names.

```sh
dig URI 2024-3094.cve @127.0.0.1
dig A 2024-3094.cve @127.0.0.1
dig AAAA 2024-3094.cve @127.0.0.1
dig TXT 2024-3094.cve @127.0.0.1
dig TXT desc.2024-3094.cve @127.0.0.1
dig TXT cvss.2024-3094.cve @127.0.0.1
dig TXT severity.2024-3094.cve @127.0.0.1
dig TXT vector.2024-3094.cve @127.0.0.1
dig TXT cwe.2024-3094.cve @127.0.0.1
dig TXT published.2024-3094.cve @127.0.0.1
dig TXT modified.2024-3094.cve @127.0.0.1
dig TXT state.2024-3094.cve @127.0.0.1
dig TXT vendors.2024-3094.cve @127.0.0.1
dig TXT products.2024-3094.cve @127.0.0.1
dig TXT refs.2024-3094.cve @127.0.0.1
dig TXT kev.2024-3094.cve @127.0.0.1
dig TXT epss.2024-3094.cve @127.0.0.1
```

A and AAAA answers are available only when their corresponding DNS address is
configured. Long TXT values may be returned as multiple character strings.
Zone indexes include the following names:

```sh
dig TXT count.cve @127.0.0.1
dig TXT count.kev.cve @127.0.0.1
dig TXT latest.cve @127.0.0.1
dig TXT built.cve @127.0.0.1
dig TXT source.cve @127.0.0.1
dig TXT version.cve @127.0.0.1
dig TXT count.year-2024.cve @127.0.0.1
dig TXT count.severity-critical.cve @127.0.0.1
```

`dig SOA cve @127.0.0.1` returns the zone SOA. Unknown names return NXDOMAIN;
existing names queried for an absent type return NOERROR with SOA authority.
Out-of-zone queries return REFUSED. `ANY` receives a minimal RFC 8482 answer.
DNSSEC signing is enabled by default. TCP AXFR/IXFR is limited by the
transfer ACL, and NOTIFY is sent to configured targets after a successful
refresh; see [DNSSEC](dnssec.md) and [Self-hosting](self-hosting.md).

### Secondary DNS servers

cvedex serves AXFR and IXFR over TCP to addresses allowed by
`DNS_TRANSFER_ACL`. Successful dataset refreshes send best-effort NOTIFY to
`DNS_NOTIFY_TARGETS`. The generated signed zone can also be downloaded in
BIND master-file format from `/api/v1/zone`.

## HTTP JSON gateway

The public data API is rooted at `/api/v1` and is read-only. Admin settings
routes are separate and support authenticated updates. The routes and response
behavior are documented in [API reference](api.md). Examples:

The host HTTP port is persisted in `server.yml`; set `CVEDEX_PORT=64580` at
startup to use the example port below.

```sh
curl -H 'Accept: application/json' \
  'http://127.0.0.1:64580/api/v1/cves/CVE-2024-3094'
curl -H 'Accept: application/json' \
  'http://127.0.0.1:64580/api/v1/search?q=linux&severity=critical&limit=10'
curl -H 'Accept: application/json' \
  'http://127.0.0.1:64580/api/v1/kev'
curl -H 'Accept: application/json' \
  'http://127.0.0.1:64580/api/v1/zone' -o cvedex.zone
```

## Data sources and mirrors

The separate dataset-generation workflow consumes CVE Program v5 JSON
records, CISA KEV JSON, and FIRST EPSS CSV. A running server downloads the
generated, checksummed snapshot from `DATA_URL` and `CHECKSUM_URL`; it does not
fetch those upstream feeds directly. Operators can mirror the generated
snapshot by setting both URLs. See [Configuration](configuration.md) and
[Self-hosting](self-hosting.md) for details.

## Other protocols

MCP, GraphQL, OAuth/OIDC, webhooks, email delivery, and federation are not
enabled in the current build. The supported external interfaces are
authoritative DNS and the HTTP JSON gateway described above.
