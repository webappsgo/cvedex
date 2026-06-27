# Backend Rules (PART 9, 10, 11, 32)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- ❌ Use mattn/go-sqlite3 — CGO dependency, forbidden
- ❌ Use bcrypt — Argon2id only for passwords
- ❌ Store tokens in plaintext — SHA-256 hash always
- ❌ Expose stack traces in production responses
- ❌ Log credentials, tokens, or passwords at any log level
- ❌ Use an external caching service (Redis, Memcached) — in-memory TTL cache only
- ❌ Use fmt.Println or log package for structured logging — use slog
- ❌ Skip rate limiting on public endpoints
- ❌ Expose Tor .onion address unless tor.expose=true in config

## CRITICAL - ALWAYS DO
- ✅ Use modernc.org/sqlite for SQLite (pure Go, no CGO)
- ✅ Hash passwords with Argon2id
- ✅ Hash tokens with SHA-256 before storage
- ✅ Return APIResponse{OK, Data, Error, Message} for all API responses
- ✅ Use in-memory TTL cache, invalidated on data sync
- ✅ Use slog (Go 1.21+) — JSON in production, human-readable in development
- ✅ Apply rate limiting on all public endpoints
- ✅ Apply RRL (Response Rate Limiting) on DNS endpoints
- ✅ Maintain IP blocklist with configurable refresh
- ✅ Apply CSRF guard on all admin form endpoints
- ✅ Check for tor binary and config flag before enabling Tor

## API RESPONSE STRUCTURE
```go
type APIResponse struct {
    OK      bool   `json:"ok"`
    Data    any    `json:"data,omitempty"`
    Error   string `json:"error,omitempty"`
    Message string `json:"message,omitempty"`
}
```

## DATABASE
| File | Purpose |
|------|---------|
| {db_dir}/server.db | CVE data, DNS zones, KEV entries, system config |
| {db_dir}/users.db | Admin accounts, sessions, API tokens |

## LOGGING (slog)
| Level | When to use |
|-------|-------------|
| DEBUG | Detailed trace for development |
| INFO | Normal operational events |
| WARN | Recoverable issues, degraded state |
| ERROR | Failures that affect functionality |

## SECURITY CONTROLS
| Control | Implementation |
|---------|---------------|
| Passwords | Argon2id (NEVER bcrypt or MD5) |
| Tokens | SHA-256 hash before storage |
| Rate limiting | Per-IP on all public endpoints |
| CSRF | Token guard on all admin forms |
| IP blocklist | Configurable, auto-refreshed by scheduler |

## TOR HIDDEN SERVICE (PART 32)
- Auto-enabled if `tor` binary found AND `tor.enabled=true` in config
- Hidden service private keys stored in {data_dir}/tor/
- Expose .onion address in /api/v1/server/about only when `tor.expose=true`
- Graceful degradation if Tor is unavailable

---
For complete details, see AI.md PART 9, 10, 11, 32
