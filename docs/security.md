# Security

## Public surfaces

- DNS answers are authoritative for the configured zone and reject out-of-zone queries with `REFUSED`.
- DNS query logging is not enabled by the current server.
- HTTP security headers include a same-origin Content Security Policy, frame denial, content-type sniffing protection, and a restrictive referrer policy.
- CVE data API routes are read-only. The preferences endpoint accepts guest preference updates; administrator API routes require an operator Bearer token, while browser forms use an authenticated administrator session.
- `/server/healthz` and `/api/v1/server/healthz` expose service health without filesystem paths or secrets.
- `/.well-known/security.txt` and `/.well-known/llms.txt` are public plain-text discovery resources. `/llms.txt` is an alias; the well-known directory itself and unknown well-known paths return 404.

## Operational controls

DNS UDP equivalent-response rate limiting is enabled by default and can be
tuned with `DNS_RRL_ENABLED`, `DNS_RRL_RATE`, and `DNS_RRL_BURST`. HTTP request
limiting is always enabled. Read requests default to 120 per client per
configured minute, writes to 10, health checks to 120, login attempts to 5 per
15 minutes, and strict token/restore actions to 10 per 15 minutes. A 240 request
burst ceiling is shared across endpoint classes per client identity. Anonymous
traffic is keyed by the normalized client IP; valid operator-token requests
use a token digest. Forwarded client IPs are accepted only from configured
trusted proxies. Limits return JSON with
`429` and `Retry-After`. In-memory counters are flushed to `server.db` every 30
seconds and on shutdown, then restored at startup. Redis or Valkey provides the
shared counter backend when configured; an unavailable backend falls back to
the persisted local counters. `HTTP_RATE_LIMIT_REQUESTS` and
`HTTP_RATE_LIMIT_WINDOW` control the general client read limit and window.

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
matched submitted domain, then rejected with `403` by default. IP rules use the
normalized client address resolved through explicitly trusted proxies. See
[Configuration](configuration.md) for source settings and freshness behavior.

First-run administrator pages use session and CSRF protections. Admin API
routes require the operator token in an `Authorization: Bearer` header and do
not accept session cookies. The settings panel controls refresh, redirect,
transfer ACL, NOTIFY, DNSSEC, DNS response limiting, and cache configuration.
HTTP request limits are enforced by endpoint class and configured through the
server configuration. DNSSEC signs the published zone when a verified dataset
is available. Protect the private configuration directory, which contains the
operator token, and the persistent data directory, which contains private
DNSSEC keys.

## Reporting

For a vulnerability in the cvedex source or release, use [GitHub private vulnerability reporting](https://github.com/webappsgo/cvedex/security/advisories/new). Do not publish exploit details in public issues.

The generated `/.well-known/security.txt` lists the private reporting channel, the instance form when enabled, the configured security contact email, and the request host's policy URL. Set `web.security.report_url` or `SECURITY_REPORT_URL` to an HTTPS reporting URL for a self-hosted distribution. The operator may replace the generated text by placing a regular file at `{DATA_DIR}/web/.well-known/security.txt`; symlinked, oversized, or unreadable overrides are ignored.

`/server/security` provides the instance reporting overview. Its rotating `security_id` link opens the security-report form at `/server/contact`; without a valid id, that form offers abuse, DMCA, and privacy reports only. It is not a support or general-contact form. No user-submitted report content is placed in `llms.txt`.

The authenticated `/server/{admin_path}/security/keypair` page manages the
project OpenPGP key used for encrypted security reports. Key rotations
re-encrypt retained report bodies before committing the new fingerprint;
private-key export, import, and deletion require password reauthentication
and are audited. Configure `web.security.keyservers` to publish the public key.

`llms.txt` is generated from the public API route metadata and request-specific base URL. It describes anonymous read-only data access and identifies operator-only MCP/admin authentication without listing protected administrator routes. An operator can replace its body at `{DATA_DIR}/web/.well-known/llms.txt`.
