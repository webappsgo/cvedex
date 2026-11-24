# Integrations

## DNS clients

The DNS listener supports UDP and TCP. Example queries for the default zone:

```sh
dig TXT severity.2024-3094.cve @127.0.0.1
dig TXT desc.2024-3094.cve @127.0.0.1
dig URI 2024-3094.cve @127.0.0.1
```

The DNS zone and port are configurable with `DNS_ZONE` and `DNS_PORT`.

## HTTP clients

Use the JSON gateway at `/api/v1`. CVE lookup, search, vendor and product reverse lookup, KEV catalog output, and BIND zone download are available when a valid dataset snapshot is loaded.

## Data sources

The dataset builder consumes CVE Program v5 JSON records, CISA KEV JSON, and FIRST EPSS CSV. Running server instances consume the generated compressed snapshot rather than fetching those upstream feeds directly.

## Other protocols

MCP, GraphQL, OAuth/OIDC, webhooks, email delivery, and federation are not enabled in the current build.
