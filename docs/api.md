# JSON API

The API is read-only and currently uses version `v1`. Responses use JSON. A missing dataset returns HTTP 503; an unknown CVE returns HTTP 404.

| Method | Path | Behavior |
|---|---|---|
| GET | `/api/v1` | API identity and health link |
| GET | `/api/v1/server/healthz` | Public service status |
| GET | `/api/v1/cves` | Dataset record count and build metadata |
| GET | `/api/v1/cves/{cve_id}` | One CVE record, including CVSS, affected products, references, KEV, and EPSS when present |
| GET | `/api/v1/search?q={text}` | Description, identifier, vendor, and product search; supports `year`, `severity`, `kev`, `limit` (1–100), and numeric offset `cursor` |
| GET | `/api/v1/vendors` | Sorted vendor names |
| GET | `/api/v1/vendors/{vendor}` | CVEs associated with a vendor |
| GET | `/api/v1/products` | Sorted product names |
| GET | `/api/v1/products/{product}` | CVEs associated with a product |
| GET | `/api/v1/kev` | KEV catalog entries in the dataset |
| GET | `/api/v1/zone` | BIND-format zone snapshot download |
| GET | `/api/v1/ds` | DNSSEC DS status; currently returns 503 because signing is not configured |

Public server information is available at `/api/v1/server/about`, `/privacy`, `/contact`, `/help`, `/help/mcp`, `/terms`, `/dpo`, `/security`, `/security/policy`, and `/security/thanks`. The corresponding HTML pages are under `/server/`.

The generated OpenAPI document is served at `/api/v1/server/swagger`; its browser page is `/server/docs/swagger`.

## Search example

```sh
curl -H 'Accept: application/json' 'http://127.0.0.1:64580/api/v1/search?q=linux&severity=critical&limit=10'
```

Search cursors are numeric offsets into the stable CVE identifier ordering. The next offset is returned as `next_cursor`, or an empty string when there are no more results.
