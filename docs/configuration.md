# Configuration

Configuration is loaded from an optional YAML file, then environment variables override it. CLI flags override both. Values are resolved at process startup.

| Variable | Default | Purpose |
|---|---:|---|
| `LISTEN` | `::` | HTTP listener address |
| `PORT` | First-run random port from `64000`–`64999` on a host; `80` in a container | HTTP listener port; the first-run host value is saved to `server.yml` |
| `CVEDEX_PORT` | unset | Project-prefixed HTTP port override; takes precedence over `PORT` |
| `DNS_ADDRESS` | `0.0.0.0` | UDP and TCP DNS bind address |
| `DNS_PORT` | `53` | UDP and TCP DNS listener port |
| `DNS_ZONE` | `cve.` | Authoritative zone; a trailing dot is added when omitted |
| `DNS_IPV4` | unset | IPv4 address returned for CVE A queries |
| `DNS_IPV6` | unset | IPv6 address returned for CVE AAAA queries |
| `DNS_RRL_ENABLED` | `true` | Enable bounded DNS UDP response rate limiting |
| `DNS_RRL_RATE` | `20` | Equivalent responses per client per second |
| `DNS_RRL_BURST` | `40` | Initial equivalent-response burst per client |
| `DNSSEC_ENABLED` | `true` | Sign each built zone offline with the persisted KSK and ZSK |
| `DNS_TRANSFER_ACL` | Loopback only when unset | IP addresses or CIDRs allowed to request TCP AXFR/IXFR |
| `DNS_NOTIFY_TARGETS` | Empty | DNS secondary hostnames or `host:port` targets notified after a successful build |
| `REDIRECT_URL_TEMPLATE` | `https://www.cve.org/CVERecord?id={cve_id}` | Shared CVE redirect and DNS URI target; placeholders are `{year}`, `{sequence}`, and `{cve_id}` |
| `HTTP_RATE_LIMIT_REQUESTS` | `120` | Read requests allowed per anonymous client or valid operator token in the general HTTP window |
| `HTTP_RATE_LIMIT_WINDOW` | `1m` | Read, write, health, and global-burst rate-limit window; minimum `1s` |
| `METRICS_ENABLED` | `true` | Enable authenticated Prometheus, Grafana, and Loki metrics services |
| `METRICS_ROOT_ENABLED` | `true` | Mount root `/metrics[/{service}]` aliases |
| `METRICS_ALLOW_UNAUTHENTICATED` | `false` | Skip metrics bearer-token checks; use only on firewalled internal networks |
| `METRICS_INCLUDE_SYSTEM` | `true` | Collect CPU, memory, and data-directory disk metrics |
| `METRICS_INCLUDE_RUNTIME` | `true` | Include Go runtime metrics |
| `METRICS_LOKI_MAX_ENTRIES` | `1000` | Maximum recent log entries served by the Loki endpoint |
| `METRICS_LOKI_MAX_AGE` | `1h` | Maximum age of entries served by the Loki endpoint |
| `METRICS_DURATION_BUCKETS` | Built-in seconds buckets | Comma-separated HTTP duration histogram buckets |
| `METRICS_SIZE_BUCKETS` | Built-in byte buckets | Comma-separated HTTP request/response size histogram buckets |
| `BACKUP_DIR` | System backup directory when writable; otherwise mode-specific fallback | Backup destination directory |
| `CLI_BINARY_DIR` | `{DATA_DIR}/cli` | Directory scanned for signed-off CLI release files exposed by autodiscovery |
| `CLI_BINARY_DOWNLOAD_REQUIRE_AUTH` | `false` | Require authentication for CLI binary downloads |
| `CLI_MIN_VERSION` | Empty | Minimum supported CLI version published by autodiscovery |
| `PID_FILE` | OS-specific | Server PID file path |
| `BASEURL` | `/` | HTTP URL path prefix (`server.baseurl` in YAML) |
| `I2P_ENABLED` | `false` on hosts; `true` in containers | I2P runtime setting; process/listener integration is not currently available |
| `SECURE_COOKIES` | `false` | Require Secure on admin session cookies |
| `PRIVACY_DATA_SOLD` | `false` | Set `server.privacy.data.sold`; enables the CCPA opt-out controls and cookie behavior |
| `ADMIN_PATH` | `administration` | Admin panel path segment |
| `MODE` | `production` | Application mode |
| `DEBUG` | `false` | Enable debug behavior |
| `CONFIG_DIR` | OS-specific; `/config/cvedex` in Docker | Configuration directory |
| `DATA_DIR` | OS-specific; `/data/cvedex` in Docker | Dataset and runtime data directory |
| `DATABASE_DIR` | `{DATA_DIR}/db` (Docker: `/data/db/sqlite`) | Directory containing `server.db` and `users.db` |
| `LOG_DIR` | OS-specific | Log directory |
| `CACHE_DIR` | OS-specific | Cache directory |
| `DATA_URL` | `https://github.com/webappsgo/cvedex/releases/download/data-latest/cvedex-data.json.gz` | HTTPS URL for the generated dataset |
| `CHECKSUM_URL` | Dataset URL plus `.sha256` | HTTPS URL for the dataset SHA-256 file |
| `REFRESH_INTERVAL` | `24h` | Dataset refresh period; minimum `1m` |
| `TZ` | `America/New_York` | IANA timezone for scheduled tasks; overrides `server.scheduler.timezone` |

The YAML file uses the same namespaced keys, for example:

```yaml
server:
  mode: production
  address: 0.0.0.0
  port: 64000
  baseurl: /
  pidfile: /run/casapps/cvedex/cvedex.pid
  data_dir: /var/lib/casapps/cvedex
  backup_dir: /mnt/Backups/casapps/cvedex
  backup:
    enabled: true
    hourly_enabled: false
    schedule: "0 2 * * *"
    retention:
      max_backups: 1
      keep_weekly: 0
      keep_monthly: 0
      keep_yearly: 0
      max_total_size: 10%
    disk_threshold: 90
    encryption:
      enabled: false
  scheduler:
    timezone: America/New_York
    catch_up_window: 1h
    tasks:
      cve_update:
        schedule: "0 5 * * *"
        enabled: true
      backup_daily:
        schedule: "0 2 * * *"
        enabled: true
      backup_hourly:
        schedule: "@hourly"
        enabled: false
  database:
    sqlite:
      path: /var/lib/casapps/cvedex/db
  data_url: https://github.com/webappsgo/cvedex/releases/download/data-latest/cvedex-data.json.gz
  checksum_url: https://github.com/webappsgo/cvedex/releases/download/data-latest/cvedex-data.json.gz.sha256
  refresh_interval: 24h
  redirect:
    url_template: https://www.cve.org/CVERecord?id={cve_id}
  rate_limit:
    requests: 120
    window: 1m
  dns:
    address: 0.0.0.0
    port: 53
    zone: cve.
    transfer_acl:
      - 127.0.0.1/32
      - ::1/128
    rate_limit:
      enabled: true
      per_second: 20
      burst: 40
    notify_targets: []
  dnssec:
    enabled: true
  ssl:
    enabled: false
    cert: /etc/casapps/cvedex/ssl/server.crt
    key: /etc/casapps/cvedex/ssl/server.key
    min_version: TLS1.2
  privacy:
    data:
      sold: false
```

The TLS listener currently supports a manually supplied certificate and key
through YAML (`server.ssl.enabled`, `server.ssl.cert`, `server.ssl.key`, and
`server.ssl.min_version`). The certificate and key must both be configured;
the minimum version is `TLS1.2` or `TLS1.3`. ACME issuance and renewal are not
available yet. Keep the private key in a protected configuration path.

The HTTP router accepts requests under the configured `baseurl` and prefixes
its page links, static assets, redirects, and admin routes. When deployed behind
a reverse proxy, it resolves `X-Forwarded-Prefix`, then `X-Forwarded-Path`,
then `X-Script-Name` before using the configured value. Set the proxy to replace
these headers with its trusted prefix.

When running as root, the default directories are `/etc/casapps/cvedex`, `/var/lib/casapps/cvedex`, `/var/log/casapps/cvedex`, and `/var/cache/casapps/cvedex`. For a regular user, XDG directories are used when set; otherwise paths are under the user's home directory. Backups use `/mnt/Backups/casapps/cvedex` when its parent is writable, then fall back to `{DATA_DIR}/backup` for system mode or the user's OS-specific backup directory. Containers default to `/data/backups/cvedex`.

The service stores the verified snapshot as `cvedex-data.json.gz` with a matching `cvedex-data.json.gz.sha256` file. A refresh validates HTTPS, the SHA-256 checksum, the schema, and the signed zone before replacing the compressed file and publishing the new snapshot. A failed refresh leaves the last verified dataset and signed zone active and reports degraded health. DNSSEC private keys are stored in `DATA_DIR/security/dnssec/` with owner-only permissions. With no valid local snapshot, readiness remains unhealthy and DNS returns a server failure for in-zone queries.

The first service start creates `{CONFIG_DIR}/secrets.yml` with mode `0600`
and a generated `server.token` operator credential. The raw token is printed
once at startup. Keep the file private; do not include it in support bundles
or publish it in logs.

The same file stores independently generated `mtr_` credentials under
`server.metrics.auth.tokens.prometheus`, `.grafana`, and `.loki`. The server
does not print these tokens; read them from the owner-only secrets file and
send them only in the metrics `Authorization` header.

Scheduled backups default to 02:00 in the scheduler timezone. The current
timezone default is `America/New_York`; override it with `TZ` or
`server.scheduler.timezone`. Configure the daily and hourly task schedules with
`server.scheduler.tasks.backup_daily.schedule` and
`server.scheduler.tasks.backup_hourly.schedule`;
retention accepts a maximum daily count, weekly,
monthly, and yearly full-backup counts, a total size limit, and a disk-use
threshold. Set `server.backup.encryption.enabled` to require encrypted
scheduled archives. The encryption password is stored encrypted in the secrets
file and is not written as plain text to `server.yml`. Restore safety limits
are available under `server.restore` and apply to admin and maintenance
restores.

`/api/autodiscover` publishes instance metadata and checksummed client release
files from `CLI_BINARY_DIR` (or `{DATA_DIR}/cli`). Release binaries appear in
the response only when their names, release version, and SHA-256 checksums
match the files in that directory.

The built-in scheduler refreshes the dataset daily by default with ±10% jitter. `SIGHUP` requests an immediate refresh. Set `DATA_URL` and `CHECKSUM_URL` to mirror a data release from another trusted distribution endpoint.
