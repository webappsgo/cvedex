# Configuration Rules (PART 5, 6, 12)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- ❌ Use strconv.ParseBool() for boolean config values — use ParseBool() wrapper
- ❌ Hardcode dev machine hostname, IP, or CPU count
- ❌ Hardcode port numbers without config/env var override
- ❌ Use .env files for config — use server.yml
- ❌ Accept only true/false for boolean fields — must accept all variations
- ❌ Expose Tier-1 secrets (tokens, passwords, keys) in any config endpoint or log

## CRITICAL - ALWAYS DO
- ✅ Config file location: {config_dir}/server.yml
- ✅ Use custom ParseBool() that accepts all 40+ boolean variations
- ✅ Detect hostname via os.Hostname() at runtime
- ✅ Detect CPU count via runtime.NumCPU() at runtime
- ✅ Detect server IP from network interfaces at runtime
- ✅ Support PORT env var override for HTTP listen port
- ✅ Listen on 0.0.0.0:80 inside container by default
- ✅ Support production (default) and development modes
- ✅ First-run setup wizard when no config exists

## BOOLEAN PARSING (ParseBool wrapper — required)
Accept ALL of these as valid boolean inputs:

| Truthy | Falsy |
|--------|-------|
| true, yes, 1, on, enable, enabled | false, no, 0, off, disable, disabled |
| TRUE, YES, ON, ENABLE, ENABLED | FALSE, NO, OFF, DISABLE, DISABLED |
| True, Yes, On, Enable, Enabled | False, No, Off, Disable, Disabled |
| t, y | f, n |
| (and other common variations) | |

## APPLICATION MODES
| Mode | Behavior |
|------|----------|
| production (default) | JSON structured logs, minimal output, error-only panics |
| development | Human-readable logs, verbose output, stack traces shown |

## RUNTIME DETECTION (NEVER hardcode)
| Value | Detection Method |
|-------|-----------------|
| hostname | os.Hostname() |
| CPU count | runtime.NumCPU() |
| server IP | detect from non-loopback network interfaces |
| listen port | PORT env var, then config, then default 80 |

## SERVER CONFIGURATION
- HTTP server binds to 0.0.0.0:80 inside container
- Port configurable via PORT environment variable
- Config file: {config_dir}/server.yml
- Database files: {db_dir}/server.db and {db_dir}/users.db

---
For complete details, see AI.md PART 5, 6, 12
