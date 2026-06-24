# cvedex — Outstanding Implementation Work

Bootstrap (PARTS 0–6) is complete. All items below are post-bootstrap implementation
tasks drawn from PARTS 7–33. PARTS 34–36 (Multi-User, Organizations, Custom Domains)
are explicitly NOT implemented for this project per the spec.

---

## [ ] Implement DNS server core (authoritative .cve TLD)
Read: AI.md PART 7

DNS server listening on UDP/TCP 53 (configurable). Serves every CVE as
`YYYY-NNNNN.cve` with A/AAAA/TXT/SOA/NS records. RRL on by default.
Pure Go DNS library (miekg/dns). No CGO.

---

## [ ] Implement DNS mode dispatch and runtime flags
Read: AI.md PART 6

Wire up the four modes — `server` (default, long-running), `generate-zone` (one-shot
zone file output), `generate-data` (upstream processing pipeline), `sync-only`
(one-shot data refresh) — with proper exit codes for one-shot modes.

---

## [ ] Implement server binary CLI flags
Read: AI.md PART 8

Full flag set: `--config`, `--data`, `--port`, `--address`, `--mode`, `--debug`,
`--version`, `--help`, `--shell`, `--color`, `--lang`, `--status`, `--daemon`,
`--pid-file`. Respect NO_COLOR. Show actual binary name (filepath.Base(os.Args[0]))
in help/version/errors. Use hardcoded project name for User-Agent and default paths.

---

## [ ] Implement signal handler (platform-specific build tags)
Read: AI.md PART 8

`src/signal/signal_unix.go` (build tag `!windows`): SIGTERM, SIGINT, SIGQUIT,
SIGUSR1 (reopen logs), SIGUSR2 (status dump), SIGRTMIN+3 (Docker stop), ignore
SIGHUP. `src/signal/signal_windows.go`: os.Interrupt only. 10-step graceful shutdown
sequence with per-phase timeouts (30s requests, 10s children, 5s DB, 2s logs).

---

## [ ] Implement DNSSEC offline signing
Read: AI.md PART 7

KSK + ZSK ECDSA P-256 keys stored in data_dir. NSEC3 authenticated denial.
Zone signed offline (not on-the-fly). Keys rotatable via admin command. Publish
DS record in console output for registrar submission.

---

## [ ] Implement AXFR / IXFR / NOTIFY
Read: AI.md PART 7

Zone transfer support (AXFR) for secondary nameservers. Incremental transfer (IXFR)
where supported. NOTIFY to secondaries on zone update. Access-controlled by ACL list
in config (default: deny all).

---

## [ ] Implement two-tier data pipeline
Read: AI.md PART 7

Tier 1 (CI job): download cvelistV5 JSON, CISA KEV catalog, FIRST.org EPSS scores,
normalize to internal format, write compressed data file, publish as GitHub release
artifact. Tier 2 (running instance): fetch Tier 1 data file on startup and on
scheduler tick, load into SQLite, serve from in-memory cache.

---

## [ ] Implement CVE data schema and SQLite layer (modernc.org/sqlite only)
Read: AI.md PART 10

`CREATE TABLE IF NOT EXISTS` for: cves (id, year, number, cve_id, description,
severity, cvss_score, epss_score, kev_listed, published_at, updated_at),
cve_references, cve_vendors, cve_products. EnsureSchema() runs on startup
(idempotent). Use modernc.org/sqlite — never mattn/go-sqlite3 (CGO).

---

## [ ] Implement error handling layer and APIResponse type
Read: AI.md PART 9

`APIResponse{OK bool, Data any, Error string, Message string}`. Standard error codes
(BAD_REQUEST 400, UNAUTHORIZED 401, FORBIDDEN 403, NOT_FOUND 404, SERVER_ERROR 500,
etc.). Never expose stack traces in production. All errors logged with context.

---

## [ ] Implement in-memory cache layer
Read: AI.md PART 9

TTL-based cache for DNS responses and JSON gateway lookups. Cache invalidated on
data sync. Per-entry TTL matching DNS record TTL values.

---

## [ ] Implement security middleware (rate limiting, RRL, blocklists)
Read: AI.md PART 11

HTTP middleware chain (chi): request ID, real-IP extraction, RRL for DNS (token
bucket per client IP), IP blocklist check, CSRF guard for admin panel. Tier-1/2/3
data classification enforced at response layer (never leak DB creds, internal IPs,
tokens).

---

## [ ] Implement structured logging
Read: AI.md PART 11

Use slog (Go 1.21+). JSON output in production, human-readable in development.
Log levels: DEBUG, INFO, WARN, ERROR. Request log includes method, path, status,
latency, request-ID. Never log credentials or tokens.

---

## [ ] Implement HTTP redirector
Read: AI.md PART 12

Host-header based 302 redirect: `CVE-YYYY-NNNNN.cve` HTTP requests → canonical
NVD/MITRE CVE page. Handles both direct and Host-header-based requests.

---

## [ ] Implement health and version endpoints
Read: AI.md PART 13

`GET /server/healthz` — returns app_name, version, commit_hash, build_date,
go_version, uptime, mode, db_type, db_locality, request counts. Public-safe
Tier-2 data only. Never expose Tier-1 secrets.
`GET /server/version` — returns version string only.

---

## [ ] Implement JSON gateway API (chi router)
Read: AI.md PART 14

Versioned API at `/api/{api_version}/`. Routes (from IDEA.md):
- `GET /api/v1/cve/{id}` — per-CVE detail (description, CVSS, EPSS, KEV, references)
- `GET /api/v1/cves` — paginated list with filters (year, severity, kev, vendor, product)
- `GET /api/v1/vendors` — vendor list
- `GET /api/v1/products` — product list with optional vendor filter
- `GET /api/v1/kev` — CISA KEV catalog
- `GET /api/v1/dns/ds` — DNSSEC DS record for .cve TLD
- `GET /api/v1/zone` — full zone file download (gzip)
- `GET /api/v1/server/healthz` — alias for /server/healthz
- `GET /api/v1/server/about` — project info, links
All API routes: plural nouns, lowercase, hyphens, no trailing slash, no verbs.
CORS: Access-Control-Allow-Origin: *.

---

## [ ] Implement SSL/TLS and Let's Encrypt auto-cert
Read: AI.md PART 15

ACME via golang.org/x/crypto/acme/autocert. Cert storage in data_dir. HTTP→HTTPS
redirect. Manual cert path option (cert_file / key_file in config). Self-signed
fallback for localhost/development.

---

## [ ] Implement web frontend (HTML/CSS/JS, embedded)
Read: AI.md PART 16

Embedded in binary via go:embed `src/server/static/` and `src/server/template/`.
Pages: CVE search/browse, per-CVE detail page, KEV list, vendor/product browse,
DNS lookup demo, DNSSEC DS record display, zone file download link.
Dark/light/auto theme. Mobile responsive. WCAG 2.1 AA. PWA manifest + service worker.
URL normalization: strip trailing slash (301 redirect).

---

## [ ] Implement admin panel (web UI + API)
Read: AI.md PART 17

Routes under `/admin/*` (protected by session auth). Pages: dashboard (stats,
uptime, last sync), data sync trigger, DNSSEC key management (rotate KSK/ZSK,
view DS record), zone configuration (SOA params, NS records, TTLs), access logs,
config viewer, first-run setup wizard (shows setup token on first start).
Admin API under `/api/v1/admin/*`.

---

## [ ] Implement first-run experience and setup token
Read: AI.md PART 17

On first run (no config found): auto-create config with defaults, generate random
setup token, print banner with URLs and setup token to stdout. Setup token grants
one-time access to admin panel. Invalidated after first admin account created.

---

## [ ] Implement email / SMTP (notification templates)
Read: AI.md PART 18

SMTP auto-detection on first run (loopback → Docker bridge → gateway → FQDN).
Embedded default templates in `src/server/template/email/`. Custom templates in
config_dir/template/email/. Live reload (no restart needed). Admin configurable
via admin panel. Graceful degradation: if no SMTP found, continue without email.

---

## [ ] Implement scheduler (CVE data sync, GeoIP, blocklists)
Read: AI.md PART 19

Background goroutine with configurable cron-like schedule. Jobs:
- CVE data sync (daily from Tier-1 data file)
- CISA KEV catalog refresh (daily)
- EPSS score refresh (daily)
- GeoIP database refresh (daily)
- IP/domain blocklist refresh (daily)
- DNSSEC zone re-signing (on data change)
- PID file heartbeat
Scheduler state visible in admin dashboard. Manual trigger via admin API.

---

## [ ] Implement GeoIP integration
Read: AI.md PART 20

ip-location-db via mmdb (pure Go mmdb reader — maxmind/mmdbwriter or oschwald/geoip2-golang).
Databases: asn.mmdb, country.mmdb, city.mmdb, whois.mmdb — downloaded on first run,
kept in data_dir/security/geoip/. Used for: request logging enrichment, admin
dashboard geographic stats. NOT for blocking by default (blocklists handle that).

---

## [ ] Implement metrics
Read: AI.md PART 21

Prometheus-compatible metrics at `/metrics` (configurable path, disabled by default
or auth-gated). Counters: dns_queries_total (by type, rcode), http_requests_total
(by method, path, status), cve_lookups_total, kev_queries_total. Histograms:
dns_query_duration_seconds, http_request_duration_seconds. Gauges: cve_count,
kev_count, uptime_seconds, goroutines.

---

## [ ] Implement backup and restore commands
Read: AI.md PART 22

`cvedex backup [--output path]` — exports SQLite DB + config to tar.gz.
`cvedex restore [--input path]` — restores from backup, verifies integrity.
Admin API endpoints: `POST /api/v1/admin/backup`, `POST /api/v1/admin/restore`.
Backup stored in data_dir/backups/ by default.

---

## [ ] Implement update command
Read: AI.md PART 23

`cvedex update` — checks GitHub releases for newer version, downloads binary for
current GOOS/GOARCH, verifies checksum, replaces self, restarts (POSIX: exec;
Windows: batch wrapper). Shows changelog excerpt. `--check` flag for dry-run.

---

## [ ] Implement privilege escalation helpers
Read: AI.md PART 24

`cvedex --service install` — installs systemd unit (Linux), launchd plist (macOS),
Windows Service (Windows). Requests sudo only when needed (not at startup).
`cvedex --service uninstall`, `--service start`, `--service stop`, `--service status`.
Port 53 capability: uses `setcap cap_net_bind_service=+ep` on Linux (not suid).

---

## [ ] Implement service support files (systemd, launchd, Windows Service)
Read: AI.md PART 25

Embed systemd unit template in binary (served by `--service install`). Unit file:
`/etc/systemd/system/cvedex.service`. Launchd plist: `io.github.casapps.cvedex.plist`
at `/Library/LaunchDaemons/`. Windows: register as Windows Service via golang.org/x/sys/windows/svc.
polkit rule for privilege escalation if needed.

---

## [ ] Implement ReadTheDocs documentation structure
Read: AI.md PART 30

`docs/` directory with mkdocs.yml or Sphinx conf.py. Sections: Getting Started,
Configuration Reference, DNS Records, JSON API Reference, DNSSEC Setup, Admin Guide,
Self-Hosting, Contributing. Auto-generated from inline doc comments where possible.

---

## [ ] Implement i18n scaffolding
Read: AI.md PART 31

`--lang` flag. Language files in `src/data/i18n/` (JSON, keyed). Embedded in binary.
Minimum: en (English) baseline. Format: `{"key": "value"}`. Applied to: CLI output,
admin panel, web frontend error messages. Falls back to `en` if key missing.

---

## [ ] Implement Tor hidden service integration
Read: AI.md PART 32

Optional: start embedded Tor process if `tor.enabled: true` in config and `tor`
binary found in PATH. Generate or load hidden service keys from data_dir/tor/.
Expose .onion address in `/api/v1/server/about` and admin dashboard (if
`tor.expose: true`). Binary handles Tor process lifecycle (start/stop with main
process). Graceful degradation: if Tor not found, log info and continue.

---

## [ ] Implement agent binary (cvedex-agent)
Read: AI.md PART 33

Separate build target `cvedex-agent`. Reports system info to cvedex server.
Flags: `--server`, `--token`, `--config`, `--help`, `--version`, `--debug`,
`--color`, `--lang`. User-Agent: `cvedex-agent/{version}`.

---

## [ ] Implement client binary (cvedex-cli)
Read: AI.md PART 33

Separate build target `cvedex-cli`. DNS query helper + JSON gateway client.
Commands: `cvedex-cli lookup CVE-2021-44228`, `cvedex-cli search --vendor apache`,
`cvedex-cli kev`, `cvedex-cli zone --output zone.txt`. Flags: `--server`,
`--output` (table/json/csv), `--help`, `--version`, `--debug`, `--color`, `--lang`.
First-run/double-click setup wizard (CLI is the only binary with the setup wizard).

---

## [ ] Write unit tests (coverage ≥ 60%)
Read: AI.md PART 29

CI enforces 60% coverage threshold. Tests required for: config loading (bool.go,
config.go), mode parsing, DNS record construction, CVE ID parsing/validation,
JSON gateway handlers (httptest), error response format, ParseBool vocabulary,
cache TTL logic, DNSSEC key loading. Use `go test -cover ./...`.

---

## [ ] Implement PID file management
Read: AI.md PART 7

Write PID file to data_dir (or /run/cvedex/ when running as root) on startup.
Remove on clean shutdown. `--pid-file` flag to override path. Check for stale PID
on startup (process no longer running → delete and continue; process running → error).

---

## [ ] Implement embedded asset directories (go:embed)
Read: AI.md PART 7

`src/server/template/` — Go HTML templates (web frontend + email defaults).
`src/server/static/` — CSS, JS, fonts, icons.
`src/data/` — application JSON (i18n, default config snippets).
All embedded with `//go:embed` directives. Never embedded: security databases,
GeoIP files, zone data (those are external and updated at runtime).

---

## [ ] Add go.sum and verify all dependencies are pure Go
Read: AI.md PART 7

After adding all required dependencies (miekg/dns, modernc.org/sqlite, chi,
oschwald/geoip2-golang or pure-go mmdb reader, golang.org/x/crypto, etc.),
run `go mod tidy` and `scripts/verify-licenses.sh`. Confirm CGO_ENABLED=0 build
succeeds with all deps. No GPL/AGPL/LGPL deps allowed.

---

## [ ] Wire up Makefile targets to actual implementations
Read: AI.md PART 26

Once `cvedex-cli` and `cvedex-agent` build targets exist, update Makefile
`release` target to build all three binaries (server + cli + agent) for all 8
platform/arch combinations. Add `make generate-data` and `make generate-zone`
one-shot targets.
