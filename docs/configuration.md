# Configuration

Configuration is loaded from an optional YAML file, then environment variables override it. CLI flags override both. Values are resolved at process startup.

| Variable | Default | Purpose |
|---|---:|---|
| `ADDRESS` | `0.0.0.0` | HTTP listener address |
| `PORT` | `64580` | HTTP listener port |
| `DNS_ADDRESS` | `0.0.0.0` | UDP and TCP DNS bind address |
| `DNS_PORT` | `53` | UDP and TCP DNS listener port |
| `DNS_ZONE` | `cve.` | Authoritative zone; a trailing dot is added when omitted |
| `DNS_IPV4` | unset | IPv4 address returned for CVE A queries |
| `DNS_IPV6` | unset | IPv6 address returned for CVE AAAA queries |
| `MODE` | `production` | Application mode |
| `DEBUG` | `false` | Enable debug behavior |
| `CONFIG_DIR` | OS-specific | Configuration directory |
| `DATA_DIR` | OS-specific | Dataset and runtime data directory |
| `DATABASE_DIR` | `{DATA_DIR}/db` (Docker: `/data/db/sqlite`) | Directory containing `server.db` and `users.db` |
| `LOG_DIR` | OS-specific | Log directory |
| `CACHE_DIR` | OS-specific | Cache directory |
| `DATA_URL` | `https://github.com/webappsgo/cvedex/releases/download/data-latest/cvedex-data.json.gz` | HTTPS URL for the generated dataset |
| `CHECKSUM_URL` | Dataset URL plus `.sha256` | HTTPS URL for the dataset SHA-256 file |
| `REFRESH_INTERVAL` | `24h` | Dataset refresh period; minimum `1m` |
| `DNS_ADDRESS` | `0.0.0.0` | DNS UDP and TCP bind address |
| `DNS_PORT` | `53` | DNS UDP and TCP listener port |
| `DNS_ZONE` | `cve.` | Authoritative zone |
| `DNS_IPV4` | unset | Address returned for CVE A records |
| `DNS_IPV6` | unset | Address returned for CVE AAAA records |

The YAML file uses the same namespaced keys, for example:

```yaml
server:
  mode: production
  address: 0.0.0.0
  port: 64580
  data_dir: /var/lib/casapps/cvedex
  database:
    sqlite:
      path: /var/lib/casapps/cvedex/db
  data_url: https://github.com/webappsgo/cvedex/releases/download/data-latest/cvedex-data.json.gz
  checksum_url: https://github.com/webappsgo/cvedex/releases/download/data-latest/cvedex-data.json.gz.sha256
  refresh_interval: 24h
  dns:
    address: 0.0.0.0
    port: 53
    zone: cve.
```

When running as root, the default directories are `/etc/casapps/cvedex`, `/var/lib/casapps/cvedex`, `/var/log/casapps/cvedex`, and `/var/cache/casapps/cvedex`. For a regular user, XDG directories are used when set; otherwise paths are under the user's home directory.

The service stores the verified snapshot as `cvedex-data.json.gz` with a matching `cvedex-data.json.gz.sha256` file. A refresh validates HTTPS, the SHA-256 checksum, and the schema before atomically replacing the compressed file and publishing the in-memory dataset. A failed refresh leaves the last verified dataset active and reports degraded health. On first start, the service attempts to download a dataset; until a valid snapshot is available, readiness is unhealthy.

The built-in scheduler refreshes the dataset daily by default. Set `DATA_URL` and `CHECKSUM_URL` to mirror a data release from another trusted distribution endpoint.
