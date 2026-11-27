# JSON API

The public data API uses version `v1` and is read-only. Admin settings routes
are separate and support authenticated changes. Responses use JSON. A missing
dataset returns HTTP 503; an unknown CVE returns HTTP 404.

The authenticated network overlay controls are documented in
[Administration](admin.md#network-overlays). They require the operator Bearer
token and use these routes:

| Method | Path | Purpose |
|---|---|---|
| GET | `/api/v1/server/{admin_path}/network/i2p` | Read I2P status and settings |
| PATCH | `/api/v1/server/{admin_path}/network/i2p` | Validate and apply partial I2P settings |
| POST | `/api/v1/server/{admin_path}/network/i2p/validate` | Validate I2P settings without applying them |
| POST | `/api/v1/server/{admin_path}/network/i2p/restart` | Request an asynchronous I2P restart |
| POST | `/api/v1/server/{admin_path}/network/i2p/regenerate` | Regenerate I2P destination and request startup |
| GET | `/api/v1/server/{admin_path}/network/tor` | Read Tor status |

The email settings and editable template catalogue are also restricted to the
operator Bearer token. SMTP credentials are accepted only on writes and are
never returned by reads. The administrator settings API includes
`trusted_proxy_additional` (an array of IPs, CIDRs, or DNS names) and
`trusted_proxy_client_ip_header` (`auto` or a supported client-IP header).
Updating either trusted-proxy setting requires a server restart.

| Method | Path | Purpose |
|---|---|---|
| GET, PATCH | `/api/v1/server/{admin_path}/email` | Read sanitized SMTP settings or update relay settings and event preferences |
| POST | `/api/v1/server/{admin_path}/email/test` | Send a real test message using sample template data |
| GET | `/api/v1/server/{admin_path}/email/templates` | List available email templates |
| GET, PUT | `/api/v1/server/{admin_path}/email/templates/{name}` | Read or validate and save a template |
| POST | `/api/v1/server/{admin_path}/email/templates/{name}/reset` | Remove a custom override and restore the embedded template with `{"confirm":true}` |
| POST | `/api/v1/server/{admin_path}/email/templates/{name}/preview` | Render supplied template data without sending email |
| GET | `/api/v1/server/{admin_path}/security/allowlist` | List trusted IP and CIDR entries |
| POST | `/api/v1/server/{admin_path}/security/allowlist` | Validate and add a trusted IP or CIDR entry |
| DELETE | `/api/v1/server/{admin_path}/security/allowlist/{id}` | Remove one trusted entry |
| GET | `/api/v1/server/{admin_path}/security/allowlist/check/{ip}` | Check whether an IP address is allowlisted |

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

## Browser pages and response formats

The public pages and data pages are server-rendered. A browser receives HTML
by default; `Accept: text/html` requests HTML explicitly.
`Accept: application/json` returns the same JSON as the paired API route (for example,
`/search` pairs with `/api/v1/search`). `Accept: text/plain`, or an HTTP tool
User-Agent such as curl, receives a readable text conversion of the rendered
page. Errors use the requested JSON, text, or HTML format as well. Text
browsers that advertise HTML receive server-rendered HTML and do not need
JavaScript for navigation or forms.

The frontend health page at `/server/healthz` supports HTML, JSON, and plain
text; `/api/v1/server/healthz` returns JSON by default and
`/api/v1/server/healthz.txt` returns plain text. `/api/v1.txt` returns the API
identity in plain text.

Guest theme and consent settings are editable at `/server/preferences` and
available as JSON at `/api/v1/server/preferences` (GET and POST). API writes
accept partial JSON or form updates and do not grant cookie consent; writes are
ignored until the visitor accepts the Preferences category on the preferences
page. Web form posts use CSRF protection.

The OpenAPI document is served at `/api/v1/server/swagger` and its direct alias `/api/swagger`; its browser page is `/server/docs/swagger`.
The GraphQL query endpoint is `/api/v1/server/graphql`, with `/api/graphql` as
the direct alias for the current version. GET requests require a `query` URL
parameter; POST requests accept a JSON object with a `query` field. GraphQL
execution errors use HTTP 200 with the standard `data` and `errors` response;
requests rejected before execution use the canonical API error shape and status.
The browser query form is at `/server/docs/graphql`. GraphQL supports the
read operations `cve`, `search`, `vendors`, `vendor`, `vendor_products`,
`products`, `product`, `latest_cves`, `kev`, `dataset`, `ds`, and `zone`. Search,
catalog, and reverse lookup queries accept bounded `limit` and `cursor`
arguments.

## Model Context Protocol

The HTTP MCP transport is enabled and accepts JSON-RPC 2.0 POST requests at `/api/v1/server/mcp`,
`/api/mcp`, and `/mcp`. The versioned route is canonical; the other two are
aliases. Each request requires a supported Bearer token,
`MCP-Protocol-Version: 2026-07-28`, a matching `Mcp-Method` header, and a
matching protocol version in the JSON-RPC `_meta` object. The endpoint currently
generates its tool catalogue from the OpenAPI route definitions and dispatches
calls through the existing HTTP handlers. It accepts the operator token and
the service-specific Prometheus, Grafana, or Loki metrics tokens. The tool list
is restricted to the permissions attached to the supplied token. User and
organization token types are not part of cvedex's product scope.

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

Administrator API routes use the server operator token as an HTTP Bearer
credential. Metrics routes use their corresponding service token. See
[Administration](admin.md) for the routes and first-run token handling. A
browser session cookie does not authenticate API requests.

The server machine's HTTP port is selected on first run and stored in `server.yml`. Set
`CVEDEX_PORT=64580` or `PORT=64580` when starting the service to use the port
shown in the examples below; `CVEDEX_PORT` takes precedence.

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
