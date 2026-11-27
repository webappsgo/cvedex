# Administration

cvedex has a first-run administrator setup page, a protected server dashboard,
and separate operator APIs. The default panel root is
`/server/administration`; set `ADMIN_PATH` to choose another path segment.
Public pages do not link to the panel. Open the configured path directly.

On first start, cvedex prints a one-time setup token and a one-time operator
token. Use the setup token to create the first administrator. The setup token
is stored as a hash in the private configuration directory. The operator
token is stored in `secrets.yml` with mode `0600`; only its SHA-256 digest is
held in memory for API authentication. Create the first administrator before
exposing the panel. Passwords use Argon2id, sessions use HttpOnly,
SameSite=Strict cookies with absolute and idle expiry, and login failures are
rate-limited. Setup, login, settings, and other state-changing browser forms
use CSRF tokens.

The administrator dashboard shows active administrator names and links to
status, settings, metrics, security, backup, and audit-log pages. Active
session addresses, durations, activity timestamps, and browser details are
not shown.

## Settings and security

The settings page is `/server/{admin_path}/settings`. It configures the
upstream redirect URL template, AXFR/IXFR client ACL, NOTIFY targets, DNSSEC
signing, and DNS response rate limiting. Dataset refresh timing is configured
through the `cve_update` task at
`server.scheduler.tasks.cve_update.schedule`; see [Configuration](configuration.md).
Validated settings are persisted in `server.yml`; values that support live
application take effect without restarting. DNSSEC key rotation is available
from this page. HTTP request limiting applies by endpoint class, but its full
class settings are not yet editable in the panel; see [Security](security.md)
for the current limits and controls.

The security page at `/server/{admin_path}/security` displays the operator
token masked. Regeneration requires the current administrator password,
replaces the token in `secrets.yml`, immediately invalidates the old token,
and records `security.operator_token_regenerated` in the audit log. The new
token is displayed once and is not returned by an API.

The status page at `/server/{admin_path}/status` shows the public health
document, process state, restart requirements, activity counters, and
component checks. Its JSON counterpart is
`GET /api/v1/server/{admin_path}/status`.

The metrics page at `/server/{admin_path}/metrics` configures the metrics
endpoints and rotates the independent Prometheus, Grafana, and Loki bearer
tokens. New tokens are shown once; the page never displays stored tokens.
System and Go runtime collector changes require a server restart. Metrics
endpoints are described in the [API reference](api.md#internal-metrics).

## Network overlays

The I2P page at `/server/{admin_path}/network/i2p` shows provider status,
address, and uptime, and lets administrators edit validated I2P settings,
restart the provider, or regenerate its destination. Changes are atomically
saved to configuration; changing whether I2P is enabled restarts the provider.
The Tor page at `/server/{admin_path}/network/tor` shows read-only status,
onion address, and uptime.

Both pages require an administrator session. Form submissions use POST and
CSRF protection. The operator API routes use a Bearer token and are listed in
the route inventory below. See [Integrations](integrations.md#tor-and-i2p-overlays)
for process behavior and recovery.

Browser requests authenticate with the administrator session. Admin API
requests use `Authorization: Bearer <operator-token>`; a browser session
cookie does not authenticate API requests. The implemented JSON routes are:

| Method | Path | Purpose |
|---|---|---|
| GET, PATCH | `/api/v1/server/{admin_path}/settings` | Read or update validated server settings |
| GET | `/api/v1/server/{admin_path}/status` | Read server status and restart information |
| GET | `/api/v1/server/{admin_path}/network/i2p` | Read I2P status and settings |
| PATCH | `/api/v1/server/{admin_path}/network/i2p` | Validate, persist, and apply partial I2P settings |
| POST | `/api/v1/server/{admin_path}/network/i2p/validate` | Validate candidate settings without applying them |
| POST | `/api/v1/server/{admin_path}/network/i2p/restart` | Request an asynchronous provider restart |
| POST | `/api/v1/server/{admin_path}/network/i2p/regenerate` | Regenerate the I2P destination and request startup |
| GET | `/api/v1/server/{admin_path}/network/tor` | Read Tor status |
| POST | `/api/v1/server/{admin_path}/security/dnssec/rotate` | Rotate DNSSEC keys |
| GET | `/api/v1/server/{admin_path}/logs/audit` | Search and page through audit entries, or export filtered results |
| GET, PATCH, POST | `/api/v1/server/{admin_path}/backup` | Read settings/history, update settings, or create a backup |
| POST | `/api/v1/server/{admin_path}/backup/restore` | Verify and restore a named backup |
| GET | `/api/v1/server/{admin_path}/scheduler` | List registered scheduled tasks and their persisted state |
| GET, PATCH | `/api/v1/server/{admin_path}/scheduler/{id}` | Read a task or update its enabled state and schedule |
| POST | `/api/v1/server/{admin_path}/scheduler/{id}/run` | Run an enabled task immediately |
| POST | `/api/v1/server/{admin_path}/scheduler/{id}/enable` | Enable a task when its configuration allows it |
| POST | `/api/v1/server/{admin_path}/scheduler/{id}/disable` | Disable a task when its configuration allows it |
| GET | `/api/v1/server/{admin_path}/scheduler/{id}/history` | Read recent task executions |

## Scheduler

The authenticated task page is `/server/{admin_path}/scheduler`. It displays
each registered task's schedule, enabled state, last and next run, recent
result, and execution history. Operators can enable or disable a task, change
its schedule, open the history page, or run it immediately. Browser schedule
changes POST to `/server/{admin_path}/scheduler/{id}`; task toggles and manual
runs also use POST forms. Critical tasks cannot be disabled. The API routes
are listed above and use the operator Bearer token; browser controls use the
authenticated session and CSRF tokens.

The scheduler registers `cve_update`, `session_cleanup`, `token_cleanup`,
`healthcheck_self`, and `update_check`. It also registers daily/hourly backups
when the backup runner is available, Tor health when Tor is installed, and I2P
health when I2P is enabled. `token_cleanup` removes revoked or expired
administrator tokens and used or expired password-reset and email-verification
tokens. TLS renewal, GeoIP and blocklist updates, security cleanup, keyserver
retries, log rotation, and cluster heartbeats do not have task runners yet.
Tor and I2P processes are managed separately from these scheduler callbacks;
see [Integrations](integrations.md#tor-and-i2p-overlays).

## Audit log

The audit viewer is `/server/{admin_path}/logs/audit`; its API is
`GET /api/v1/server/{admin_path}/logs/audit`. Both require administrator
authentication (a session in the browser, or the operator bearer token for the
API). The viewer supports pagination and filters for category, severity,
result, actor, IP address, event, and date range. Plain GET links export the
current page, all matching entries, or the full log as JSON or CSV.

Audit events are written as JSON Lines to `LOG_DIR/audit.log`. Entries are
sanitized for sensitive values and form a SHA-256 hash chain. The writer
verifies the active file at startup, keeps the chain anchor in
`LOG_DIR/audit.anchor`, and rotates daily by default. The current retention
policy is `keep: none`, so rotated files are removed after their last hash has
been anchored.

## Backup and restore

The browser interface is `/server/{admin_path}/backup`. It shows backup files
and history, configures the schedule, retention, and encryption password, and
can create or restore a backup.

Scheduled backups default to 02:00 in the scheduler timezone, which defaults
to `America/New_York`. Set `server.scheduler.timezone` in YAML or override it
with the `TZ` environment variable. Hourly incrementals are disabled by
default. Retention can keep daily, weekly, monthly, and yearly full backups and
apply a total-size limit. New backups are verified before old backups are
removed. The restore path validates file names and archive contents, applies
the configured archive and expansion limits, and coordinates service shutdown
and restart while replacing state. The standalone `--maintenance restore`
command requires the service to be stopped first.

The server binary also provides maintenance commands. By default, backups are
written to the resolved `BACKUP_DIR`. A terminal backup prompts for a password;
leaving it empty creates an unencrypted archive unless encryption is enabled
in the backup settings. The operator token is blanked in the backup copy.
Optional data and SSL contents can be added with `--include-data` and
`--include-ssl`.

```sh
cvedex --maintenance backup --password 'backup-password'
cvedex --maintenance backup --include-data --include-ssl
cvedex --maintenance backup test
cvedex --maintenance backup test cvedex_backup_2026-01-01.tar.gz.enc --password 'backup-password'
cvedex --maintenance restore --password 'backup-password' /path/to/cvedex_backup_2026-01-01.tar.gz.enc
```

Encrypted restore prompts for the password when it is not supplied through the
supported command option. Keep the backup password separately: it is required
to verify or restore encrypted archives.

## Updates

The server update command can inspect and install releases from the `stable`,
`beta`, or `devel` channel. The scheduled `update_check` honors
`server.update.defer_days`; the manual command ignores that delay. With
`server.update.auto_install: true`, a privileged Linux `--service-runtime`
stages a verified release and exits with status 75 so a supervising service
can restart it as root. The root startup phase fetches the release checksum
again before replacing the binary. Other runs notify administrators and leave
installation to an explicit elevated command.

```sh
cvedex --update check
sudo cvedex --update yes
sudo cvedex --update branch beta
```
