# Features Rules (PART 18-23)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- ❌ Use an external cron daemon or OS scheduler — built-in scheduler only
- ❌ Fail hard when SMTP is not configured — graceful degradation required
- ❌ Hardcode GeoIP database paths — use {data_dir}/security/geoip/
- ❌ Expose /metrics without authentication (auth-gated by default)
- ❌ Store backups outside {data_dir}/backups/
- ❌ Skip checksum verification on self-update downloads
- ❌ Use external task queue or job runner libraries

## CRITICAL - ALWAYS DO
- ✅ Scheduler runs as a background goroutine inside the server process
- ✅ SMTP: auto-detect on first run, embed templates, degrade gracefully if absent
- ✅ GeoIP: use ip-location-db (mmdb format), databases in {data_dir}/security/geoip/
- ✅ Metrics: Prometheus-compatible endpoint at /metrics, auth-gated by default
- ✅ Backup: `cvedex backup` and `cvedex restore` commands
- ✅ Update: `cvedex update` checks GitHub releases, verifies checksum, replaces self
- ✅ All scheduler jobs are configurable (enable/disable, interval)

## EMAIL (PART 18)
| Item | Requirement |
|------|-------------|
| Config | SMTP auto-detection on first run |
| Templates | Embedded in src/server/template/email/ via go:embed |
| Degradation | All email features silently disabled if no SMTP configured |
| No hard deps | Server must start and run without email configured |

## SCHEDULER JOBS (PART 19) — background goroutine, NEVER external cron
| Job | Trigger |
|-----|---------|
| CVE data sync | Daily (configurable interval) |
| GeoIP database refresh | Configurable interval |
| IP blocklist refresh | Configurable interval |
| DNSSEC re-signing | On CVE/zone data change |

## GEOIP (PART 20)
- Library: ip-location-db (mmdb format)
- Database storage: {data_dir}/security/geoip/
- Auto-downloaded on first run if absent
- Refreshed by scheduler job

## METRICS (PART 21)
- Endpoint: /metrics
- Format: Prometheus-compatible text format
- Auth: gated by default (configurable)
- Must include: request counts, error rates, CVE lookup counts, sync status

## BACKUP & RESTORE (PART 22)
| Command | Behavior |
|---------|----------|
| cvedex backup | Archive {data_dir} to {data_dir}/backups/{timestamp}.tar.gz |
| cvedex restore {file} | Extract backup, verify integrity, replace data |

## UPDATE COMMAND (PART 23)
1. Check latest release via GitHub API
2. Compare version strings
3. Download binary for current OS/arch
4. Verify SHA-256 checksum against published checksums file
5. Replace current binary (atomic rename)
6. Exit with message to restart

---
For complete details, see AI.md PART 18, 19, 20, 21, 22, 23
