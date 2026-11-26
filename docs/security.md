# Security

## Public surfaces

- DNS answers are authoritative for the configured zone and reject out-of-zone queries with `REFUSED`.
- DNS query logging is not enabled by the current server.
- HTTP security headers include a same-origin Content Security Policy, frame denial, content-type sniffing protection, and a restrictive referrer policy.
- Public API routes are read-only. Dataset files are checked against SHA-256 and schema validation before use.
- `/server/healthz` and `/api/v1/server/healthz` expose service health without host paths or secrets.

## Operational controls

DNS UDP equivalent-response rate limiting is enabled by default and can be
tuned with `DNS_RRL_ENABLED`, `DNS_RRL_RATE`, and `DNS_RRL_BURST`. HTTP
request limiting is enabled by default at 120 requests per client IP per
minute and can be tuned with `HTTP_RATE_LIMIT_REQUESTS` and
`HTTP_RATE_LIMIT_WINDOW`. The limiter is always enabled. Client addresses use
the direct connection peer; forwarded headers are not trusted. First-run
administrator pages use session and CSRF protections. Admin API routes require
the operator token in an `Authorization: Bearer` header and do not accept
session cookies. The settings panel controls refresh, redirect, transfer ACL,
NOTIFY, DNSSEC, and DNS response limiting. DNSSEC signs the published zone
when a verified dataset is available. Protect the private configuration
directory, which contains the operator token, and the persistent data directory,
which contains private DNSSEC keys.

## Reporting

Use GitHub's private vulnerability reporting for the `webappsgo/cvedex` repository when it is available. Avoid publishing exploit details in public issues.
