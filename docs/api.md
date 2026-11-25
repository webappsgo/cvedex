# JSON API

The public data API uses version `v1` and is read-only. Admin settings routes
are separate and support authenticated changes. Responses use JSON. A missing
dataset returns HTTP 503; an unknown CVE returns HTTP 404.

| Method | Path | Behavior |
|---|---|---|
| GET | `/api/v1` | API identity and health link |
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

Public server information is available at `/api/v1/server/about`, `/privacy`, `/contact`, `/help`, `/help/mcp`, `/terms`, `/dpo`, `/security`, `/security/policy`, and `/security/thanks`. The corresponding HTML pages are under `/server/`.

Guest theme and consent settings are editable at `/server/preferences` and
available as JSON at `/api/v1/server/preferences` (GET and CSRF-protected
PATCH). The JSON API does not grant cookie consent; the visitor must make that
choice on the preferences page.

The generated OpenAPI document is served at `/api/v1/server/swagger`; its browser page is `/server/docs/swagger`.

Admin settings and DNSSEC key rotation use the server operator token as an
HTTP Bearer credential. See [Administration](admin.md) for the routes and
first-run token handling. A browser session cookie does not authenticate API
requests.

The host HTTP port is selected on first run and stored in `server.yml`. Set
`CVEDEX_PORT=64580` when starting the service to use the port shown in the
examples below.

## Search example

```sh
curl -H 'Accept: application/json' 'http://127.0.0.1:64580/api/v1/search?q=linux&severity=critical&limit=10'
```

List and search endpoints use a numeric offset cursor. Responses include `total`, `limit`, `cursor`, and `next_cursor`; the next offset is returned as a string, or an empty string when there are no more results. Vendor lookup responses include `severity_counts` and `most_recent` across the full match set, while `cve_ids` is the requested page. Latest results are sorted by published date descending, with the CVE ID breaking ties.

## Lookup examples

```sh
curl 'http://127.0.0.1:64580/api/v1/cves/latest?limit=10'
curl 'http://127.0.0.1:64580/api/v1/vendors/acme?limit=20&cursor=0'
curl 'http://127.0.0.1:64580/api/v1/vendors/acme/products/widget?limit=20'
curl 'http://127.0.0.1:64580/api/v1/vendors/acme/products?limit=50'
```
