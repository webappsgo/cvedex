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

The settings page is `/server/{admin_path}/settings`. It configures the data
refresh interval, upstream redirect URL template, AXFR/IXFR client ACL, NOTIFY
targets, DNSSEC signing, DNS response rate limiting, and HTTP rate limiting.
Validated settings are persisted in `server.yml`; values that support live
application take effect without restarting. DNSSEC key rotation is available
from this page.

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

Browser requests authenticate with the administrator session. Admin API
requests use `Authorization: Bearer <operator-token>`; a browser session
cookie does not authenticate API requests. The implemented JSON routes are:

| Method | Path | Purpose |
|---|---|---|
| GET, PATCH | `/api/v1/server/{admin_path}/settings` | Read or update validated server settings |
| GET | `/api/v1/server/{admin_path}/status` | Read server status and restart information |
| POST | `/api/v1/server/{admin_path}/security/dnssec/rotate` | Rotate DNSSEC keys |
| GET | `/api/v1/server/{admin_path}/logs/audit` | Search and page through audit entries, or export filtered results |
| GET, PATCH, POST | `/api/v1/server/{admin_path}/backup` | Read settings/history, update settings, or create a backup |
| POST | `/api/v1/server/{admin_path}/backup/restore` | Verify and restore a named backup |

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

Scheduled backups default to 02:00 UTC. Hourly incrementals are disabled by
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
`beta`, or `devel` channel. Installation requires the privileges needed to
replace the server binary; the service manager performs the restart.

```sh
cvedex --update check
sudo cvedex --update yes
sudo cvedex --update branch beta
```
