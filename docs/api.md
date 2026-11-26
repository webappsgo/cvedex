# JSON API

The public data API uses version `v1` and is read-only. Admin settings routes
are separate and support authenticated changes. Responses use JSON. A missing
dataset returns HTTP 503; an unknown CVE returns HTTP 404.

| Method | Path | Behavior |
|---|---|---|
| GET | `/api/v1` | API identity and health link |
| GET | `/api/autodiscover` | Public instance metadata and checksummed CLI release inventory |
| GET | `/api/v1/server/healthz` | Public service status |
| GET | `/api/v1/cves` | Dataset record count and build metadata |
| GET | `/api/v1/cves/{cve_id}` | One CVE record, including CVSS, affected products, references, KEV, and EPSS when present |
| GET | `/api/v1/search?q={text}` | Description, identifier, vendor, and product search; supports `year`, `severity`, `kev`, `limit` (1–100), and numeric offset `cursor` |
| GET | `/api/v1/vendors?limit=20&cursor=0` | Paginated vendors with CVE counts |
| GET | `/api/v1/vendors/{vendor}?limit=20&cursor=0` | Paginated CVE IDs, severity counts, and most recent publish date for a vendor |
| GET | `/api/v1/vendors/{vendor}/products?limit=20&cursor=0` | Paginated products and CVE counts for a vendor |
| GET | `/api/v1/vendors/{vendor}/products/{product}?limit=20&cursor=0` | Paginated CVE IDs for a vendor/product pair |
| GET | `/api/v1/products?limit=20&cursor=0` | Paginated products with CVE counts |
| GET | `/api/v1/products/{product}?limit=20&cursor=0` | Paginated CVE IDs associated with a product across vendors |
| GET | `/api/v1/cves/latest?limit=20` | Most recently published CVEs; limit is configurable from 1 to 100 |
| GET | `/api/v1/kev` | KEV catalog entries in the dataset |
| GET | `/api/v1/zone` | BIND-format zone snapshot download |
| GET | `/api/v1/ds` | Current DNSSEC DS records; returns 503 until a signed zone is available |
| GET | `/api/swagger` | Direct alias for the current version's OpenAPI document |
| GET, POST | `/api/graphql` | Direct aliases for the current version's read-only GraphQL endpoint |
| GET | `/api/v1/server/graphql` | Execute a read-only GraphQL query supplied in the `query` URL parameter |
| POST | `/api/v1/server/graphql` | Execute a read-only GraphQL query supplied in a JSON request body |
| POST | `/api/v1/server/mcp` | MCP JSON-RPC transport; aliases `/api/mcp` and `/mcp` |

Public server information is available at `/api/v1/server/about`, `/privacy`, `/contact`, `/help`, `/help/mcp`, `/terms`, `/dpo`, `/security`, `/security/policy`, and `/security/thanks`. The corresponding HTML pages are under `/server/`.

Guest theme and consent settings are editable at `/server/preferences` and
available as JSON at `/api/v1/server/preferences` (GET and CSRF-protected
PATCH). The JSON API does not grant cookie consent; the visitor must make that
choice on the preferences page.

The OpenAPI document is served at `/api/v1/server/swagger` and its direct alias `/api/swagger`; its browser page is `/server/docs/swagger`.
The GraphQL query endpoint is `/api/v1/server/graphql`, with `/api/graphql` as
the direct alias for the current version. GET requests require a `query` URL
parameter; POST requests accept a JSON object with a `query` field. The
server-rendered query form is at `/server/docs/graphql`. GraphQL supports the
read operations `cve`, `search`, `vendors`, `vendor`, `vendor_products`,
`products`, `product`, `latest_cves`, `kev`, `dataset`, `ds`, and `zone`. Search,
catalog, and reverse lookup queries accept bounded `limit` and `cursor`
arguments.

## Model Context Protocol

The HTTP MCP transport accepts JSON-RPC 2.0 POST requests at `/api/v1/server/mcp`,
`/api/mcp`, and `/mcp`. Each request requires the operator Bearer token,
`MCP-Protocol-Version: 2026-07-28`, a matching `Mcp-Method` header, and a
matching protocol version in the JSON-RPC `_meta` object. The endpoint currently
generates its tool catalogue from the OpenAPI route definitions and dispatches
calls through the existing HTTP handlers. This implementation currently accepts
the operator bearer token; scoped user and organization tokens are unavailable
while PART 34 remains optional. That template conflict is recorded in
`~/Projects/github/webappsgo/TODO.SERVER.md`.

## Internal metrics

Metrics are enabled by default and require a service-specific bearer token
from `CONFIG_DIR/secrets.yml`. Tokens are never accepted in query parameters.
Prometheus exposition is available at `/server/metrics` and
`/server/metrics/prometheus`; the Grafana dashboard JSON and Loki-compatible
recent log stream are at `/server/metrics/grafana` and `/server/metrics/loki`.
The same handlers are mounted at `/api/v1/server/metrics[/{service}]`,
`/api/metrics[/{service}]`, and `/metrics[/{service}]`. The root aliases are enabled
by default and can be disabled with `server.metrics.root.enabled`.

Each service uses its own credential: `prometheus`, `grafana`, or `loki`.
Send it only as an HTTP `Authorization: Bearer <token>` header. The
`server.metrics.auth.allow_unauthenticated` setting disables authentication
for all services and is intended only for firewalled internal networks.

The registry includes application build and uptime information, normalized
HTTP traffic and latency, SQLite query and connection data, administrator
authentication and active-session aggregates, scheduler metric families, and
HTTP/DNS rate-limit outcomes. Runtime and CPU/memory/data-directory disk
families follow `include_runtime` and `include_system`. The Loki service exposes
only bounded recent structured logs after credential redaction.

Admin settings and DNSSEC key rotation use the server operator token as an
HTTP Bearer credential. See [Administration](admin.md) for the routes and
first-run token handling. A browser session cookie does not authenticate API
requests.

The host HTTP port is selected on first run and stored in `server.yml`. Set
`CVEDEX_PORT=64580` when starting the service to use the port shown in the
examples below.

## Search example

```sh
curl -q -LSsf -H 'Accept: application/json' 'http://127.0.0.1:64580/api/v1/search?q=linux&severity=critical&limit=10'
```

List and search endpoints use a numeric offset cursor. Responses include `total`, `limit`, `cursor`, and `next_cursor`; the next offset is returned as a string, or an empty string when there are no more results. Vendor lookup responses include `severity_counts` and `most_recent` across the full match set, while `cve_ids` is the requested page. Latest results are sorted by published date descending, with the CVE ID breaking ties.

## Lookup examples

```sh
curl -q -LSsf 'http://127.0.0.1:64580/api/v1/cves/latest?limit=10'
curl -q -LSsf 'http://127.0.0.1:64580/api/v1/vendors/acme?limit=20&cursor=0'
curl -q -LSsf 'http://127.0.0.1:64580/api/v1/vendors/acme/products/widget?limit=20'
curl -q -LSsf 'http://127.0.0.1:64580/api/v1/vendors/acme/products?limit=50'
```

## GraphQL example

```sh
curl -q -LSsf -H 'Content-Type: application/json' \
  -d '{"query":"{ dataset { record_count built_at } cve(id: \"CVE-2024-3094\") { id state description } }"}' \
  http://127.0.0.1:64580/api/graphql
```
