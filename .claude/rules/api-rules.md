# API Rules (PART 13, 14, 15)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- ❌ Expose Tier-1 secrets (tokens, passwords, DB paths) in /server/healthz
- ❌ Use unversioned API routes (always /api/v1/...)
- ❌ Use verbs in route paths (use nouns only)
- ❌ Add trailing slashes to routes
- ❌ Use camelCase or underscores in route paths — lowercase hyphens only
- ❌ Return CORS errors — always set Access-Control-Allow-Origin: *
- ❌ Use self-signed certs in production without config opt-in
- ❌ Redirect HTTP to HTTPS without a valid cert in place first

## CRITICAL - ALWAYS DO
- ✅ All API routes versioned under /api/v1/
- ✅ Route nouns are plural and lowercase
- ✅ Health endpoint returns only Tier-2 safe data
- ✅ Set Access-Control-Allow-Origin: * on all API responses
- ✅ TLS via ACME (autocert) with HTTP→HTTPS redirect
- ✅ Store ACME certs in {data_dir}
- ✅ Support manual cert path config option
- ✅ Fall back to self-signed cert for localhost/dev

## REQUIRED API ROUTES
| Method | Route | Description |
|--------|-------|-------------|
| GET | /api/v1/cve/{id} | Single CVE lookup by ID |
| GET | /api/v1/cves | CVE list/search |
| GET | /api/v1/vendors | Vendor list |
| GET | /api/v1/products | Product list |
| GET | /api/v1/kev | CISA KEV list |
| GET | /api/v1/dns/ds | DNSSEC DS records |
| GET | /api/v1/zone | DNS zone file download |
| GET | /api/v1/server/healthz | Health check |
| GET | /api/v1/server/about | Server info |

## HEALTH CHECK (/api/v1/server/healthz)
Tier-2 data ONLY — safe to expose publicly:

| Field | Source |
|-------|--------|
| app_name | "cvedex" |
| version | build-time constant |
| commit_hash | build-time constant |
| build_date | build-time constant |
| go_version | runtime.Version() |
| uptime | time since start |
| mode | production/development |
| db_type | "sqlite" |

## ROUTE NAMING RULES
- Plural nouns: /cves not /cve-list, /vendors not /vendor
- Lowercase only: /api/v1/kev not /api/v1/KEV
- Hyphens for multi-word: /api/v1/dns/ds not /api/v1/dns_ds
- No trailing slash: /api/v1/cves not /api/v1/cves/
- No verbs: /api/v1/cves not /api/v1/getCves

## SSL/TLS
| Scenario | Behavior |
|----------|----------|
| Production with domain | ACME via golang.org/x/crypto/acme/autocert |
| Manual cert configured | Use provided cert/key paths |
| Localhost/development | Self-signed cert (auto-generated) |
| HTTP request in TLS mode | 301 redirect to HTTPS |

---
For complete details, see AI.md PART 13, 14, 15
