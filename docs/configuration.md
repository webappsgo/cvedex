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
| `HTTP_RATE_LIMIT_ENABLED` | `true` | Enable per-client HTTP request limiting |
| `HTTP_RATE_LIMIT_REQUESTS` | `100` | Requests allowed per client in each window |
| `HTTP_RATE_LIMIT_WINDOW` | `1m` | HTTP rate-limit window; minimum `1s` |
| `BACKUP_DIR` | System backup directory when writable; otherwise mode-specific fallback | Backup destination directory |
| `PID_FILE` | OS-specific | Server PID file path |
| `BASEURL` | `/` | HTTP URL path prefix (`server.baseurl` in YAML) |
| `I2P_ENABLED` | `false` on hosts; `true` in containers | Enable the optional I2P eepsite |
| `SECURE_COOKIES` | `false` | Require Secure on admin session cookies |
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
  database:
    sqlite:
      path: /var/lib/casapps/cvedex/db
  data_url: https://github.com/webappsgo/cvedex/releases/download/data-latest/cvedex-data.json.gz
  checksum_url: https://github.com/webappsgo/cvedex/releases/download/data-latest/cvedex-data.json.gz.sha256
  refresh_interval: 24h
  redirect:
    url_template: https://www.cve.org/CVERecord?id={cve_id}
  rate_limit:
    enabled: true
    requests: 100
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
```

The HTTP router accepts requests under the configured `baseurl` and prefixes
its page links, static assets, redirects, and admin routes. When deployed behind
a reverse proxy, it resolves `X-Forwarded-Prefix`, then `X-Forwarded-Path`,
then `X-Script-Name` before using the configured value. Set the proxy to replace
these headers with its trusted prefix.

When running as root, the default directories are `/etc/casapps/cvedex`, `/var/lib/casapps/cvedex`, `/var/log/casapps/cvedex`, and `/var/cache/casapps/cvedex`. For a regular user, XDG directories are used when set; otherwise paths are under the user's home directory. Backups use `/mnt/Backups/casapps/cvedex` when its parent is writable, then fall back to `{DATA_DIR}/backup` for system mode or the user's OS-specific backup directory. Containers default to `/data/backups/cvedex`.

The service stores the verified snapshot as `cvedex-data.json.gz` with a matching `cvedex-data.json.gz.sha256` file. A refresh validates HTTPS, the SHA-256 checksum, the schema, and the signed zone before replacing the compressed file and publishing the new snapshot. A failed refresh leaves the last verified dataset and signed zone active and reports degraded health. DNSSEC private keys are stored in `DATA_DIR/security/dnssec/` with owner-only permissions. With no valid local snapshot, readiness remains unhealthy and DNS returns a server failure for in-zone queries.

The built-in scheduler refreshes the dataset daily by default with ±10% jitter. `SIGHUP` requests an immediate refresh. Set `DATA_URL` and `CHECKSUM_URL` to mirror a data release from another trusted distribution endpoint.
