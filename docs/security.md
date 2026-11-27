# Security

## Public surfaces

- DNS answers are authoritative for the configured zone and reject out-of-zone queries with `REFUSED`.
- DNS query logging is not enabled by the current server.
- HTTP security headers include a same-origin Content Security Policy, frame denial, content-type sniffing protection, and a restrictive referrer policy.
- CVE data API routes are read-only. The preferences endpoint accepts guest preference updates; administrator API routes require an operator Bearer token, while browser forms use an authenticated administrator session.
- `/server/healthz` and `/api/v1/server/healthz` expose service health without filesystem paths or secrets.

## Operational controls

DNS UDP equivalent-response rate limiting is enabled by default and can be
tuned with `DNS_RRL_ENABLED`, `DNS_RRL_RATE`, and `DNS_RRL_BURST`. HTTP request
limiting is always enabled. Read requests default to 120 per client per
configured minute, writes to 10, health checks to 120, login attempts to 5 per
15 minutes, and strict token/restore actions to 10 per 15 minutes. A 240 request
burst ceiling is shared across endpoint classes per client identity. Anonymous
traffic is keyed by the direct connection peer; valid operator-token requests
use a token digest. Forwarded headers are not trusted. Limits return JSON with
`429` and `Retry-After`. The current implementation keeps counters in memory,
so they reset on restart; persistent counters and cache-backed coordination are
not implemented yet. `HTTP_RATE_LIMIT_REQUESTS` and `HTTP_RATE_LIMIT_WINDOW`
control the general client read limit and window; other class defaults are
currently fixed.

Entries under `server.security.allowlist` bypass IP blocklists, HTTP rate
limits, and GeoIP country rules. They do not bypass administrator sessions,
API token validation, CSRF, or path security. Administrators manage entries at
`/server/{admin_path}/security/allowlist` or through the authenticated API;
changes persist in `server.yml` and take effect immediately.

GeoIP is a deny-only risk signal. It runs after HTTP rate limiting and before
route authentication, so blocked traffic consumes rate-limit capacity. Missing
or stale databases, failed lookups, unknown countries, and private/internal
addresses fail open to the rest of the security pipeline. GeoIP must never be
used as authentication or as a substitute for authorization. Operators can
configure country policies and named presets at
`/server/{admin_path}/network/geoip`.

IP and domain blocklists run before HTTP rate limiting and GeoIP. The built-in
FireHOL level-one IP list and enabled operator sources update over HTTPS, with
each valid source retained on disk and hot-loaded independently. Unsafe IP
prefixes and public-suffix wildcards are discarded during validation. Domain
matching uses Origin or Referer and submitted email or URL fields; it does not
use the request Host or reverse DNS. Blocked requests are logged without the
matched submitted domain, then rejected with `403` by default. See
[Configuration](configuration.md) for source settings and freshness behavior.

First-run administrator pages use session and CSRF protections. Admin API
routes require the operator token in an `Authorization: Bearer` header and do
not accept session cookies. The settings panel controls refresh, redirect,
transfer ACL, NOTIFY, DNSSEC, and DNS response limiting. HTTP limit settings are
not yet editable there. DNSSEC signs the published zone when a verified dataset
is available. Protect the private configuration directory, which contains the
operator token, and the persistent data directory, which contains private
DNSSEC keys.

## Reporting

Use GitHub's private vulnerability reporting for the `webappsgo/cvedex` repository when it is available. Avoid publishing exploit details in public issues.
