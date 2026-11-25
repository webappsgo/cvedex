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
request limiting is enabled by default at 100 requests per client IP per
minute and can be tuned with `HTTP_RATE_LIMIT_ENABLED`,
`HTTP_RATE_LIMIT_REQUESTS`, and `HTTP_RATE_LIMIT_WINDOW`. Client addresses use
the direct connection peer; forwarded headers are not trusted. First-run
administrator authentication uses session and CSRF protections. The settings
panel controls refresh, redirect, transfer ACL, NOTIFY, DNSSEC, and DNS
response limiting. DNSSEC signs the published zone when a verified dataset is
available. Protect the persistent data directory because it contains private
DNSSEC keys.

## Reporting

Use GitHub's private vulnerability reporting for the `casapps/cvedex` repository. Avoid publishing exploit details in public issues.
