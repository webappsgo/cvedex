# Security

## Public surfaces

- DNS answers are authoritative for the configured zone and reject out-of-zone queries with `REFUSED`.
- DNS query logging is not enabled by the current server.
- HTTP security headers include a same-origin Content Security Policy, frame denial, content-type sniffing protection, and a restrictive referrer policy.
- Public API routes are read-only. Dataset files are checked against SHA-256 and schema validation before use.
- `/server/healthz` and `/api/v1/server/healthz` expose service health without host paths or secrets.

## Current limitations

DNSSEC signing, zone transfer ACLs (AXFR/IXFR), response-rate limiting, HTTP rate limiting, admin authentication, and a configurable vulnerability disclosure contact are not implemented yet. Do not treat the current unsigned DNS responses as authenticated data.

## Reporting

Use GitHub's private vulnerability reporting for the `casapps/cvedex` repository. Avoid publishing exploit details in public issues.
