# Configuration

Configuration is currently read from environment variables. Values are applied at process startup.

| Variable | Default | Purpose |
|---|---:|---|
| `ADDRESS` | `0.0.0.0` | HTTP and DNS bind address |
| `PORT` | `64580` | HTTP listener port |
| `DNS_PORT` | `53` | UDP and TCP DNS listener port |
| `DNS_ZONE` | `cve.` | Authoritative zone; a trailing dot is added when omitted |
| `MODE` | `production` | Application mode |
| `DEBUG` | `false` | Enable debug behavior |
| `CONFIG_DIR` | OS-specific | Configuration directory |
| `DATA_DIR` | OS-specific | Dataset and runtime data directory |
| `LOG_DIR` | OS-specific | Log directory |
| `CACHE_DIR` | OS-specific | Cache directory |

When running as root, the default directories are `/etc/casapps/cvedex`, `/var/lib/casapps/cvedex`, `/var/log/casapps/cvedex`, and `/var/cache/casapps/cvedex`. For a regular user, XDG directories are used when set; otherwise paths are under the user's home directory.

The data snapshot is `cvedex-data.json.gz`, accompanied by `cvedex-data.json.gz.sha256`. The service verifies the checksum and schema before publishing records. Invalid snapshots fail startup and do not become live data.

Configuration-file support and live reload are not implemented yet; environment values are the supported interface in this version.
