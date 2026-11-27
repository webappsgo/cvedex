# cvedex AI.md Alignment Tasks

## Critical (P0)

- [ ] Finish server page and API route coverage against IDEA.md and AI.md PARTS 13, 14, 16, and 17; ensure every public data endpoint has a working server-rendered page and page errors follow the shared negotiation and accessibility rules.
- [ ] Verify and finish the existing GraphQL API and GraphiQL page in `src/graphql/`; keep query fields synchronized with the read API contract and confirm versioned and unversioned routes.
- [ ] Finish admin authentication, settings, security, status, and audit-log coverage required by PART 17; expose every applicable setting in cvedex's admin UI with validation, help, defaults, and restart feedback.
- [ ] Implement required `/.well-known/security.txt` and `/.well-known/llms.txt` endpoints, the security contact/policy flows and per-request URL resolution required by PARTS 8, 11, 16, and 17; preserve operator override and public error/accessibility rules.
- [ ] Complete the contact and coordinated security-report pipeline from PARTS 11 and 16: configured CAPTCHA providers and no-JS fallbacks, translated field labels, canonical encryption-key use, researcher acknowledgment/status token, report status pages and authenticated admin workflow, PGP handling, retention, and contact configuration/admin controls.
- [ ] Complete the health contract in PART 10: report maintenance and self-healing details, use configurable project branding, expose the dynamic dependency checks, and verify every JSON/text status variant.
- [ ] Add the required browser viewport conformance checks for essential toggles at 320, 768, and 1280 pixels.

## High (P1)

- [ ] Implement the required shared theme, i18n, and accessibility systems from PARTS 16 and 31 across pages, admin, CLI, and health output.
- [ ] Complete the existing `/server/{admin_path}/metrics` settings and token-management page from PART 21; audit scheduler/database metric families as those subsystems are implemented.
- [ ] Complete the always-on HTTP sliding-window limiter: configurable endpoint classes, per-token identity for every valid token type, failed-auth lockout, strict-action owner limits, persisted hot counters, cache backend/fallback, and admin controls. Read/write/health/login/strict limits, operator-token grouping, and the canonical 429 response are implemented but remain unverified against full PART 12 coverage.
- [ ] Implement required scheduler tasks and admin controls from PART 19, preserving cvedex's daily refresh behavior.
- [ ] Add the mandatory cache subsystem from AI.md's Cache Configuration section: none/memory/Valkey/Redis modes, URL or host credentials, TLS, pools, prefix, TTL and cluster nodes; connect sessions, distributed rate limits, health, metrics, heartbeat, Pub/Sub and locks; expose settings in admin UI and docs.
- [ ] Resolve safe offline navigation caching after the AI.md cache-eligibility contradiction is corrected; current worker deliberately serves a synthesized offline page instead of a language/theme-specific cached document.
- [ ] Finish and verify the SQLite concurrency contract: audit all queries for bounded timeouts and all mutations for the write helper, then run the required 20-writer/20-reader concurrency test from PART 10. Separate handles, immediate transactions, retry/backoff, and shutdown checkpointing are implemented.
- [ ] Implement mandatory backup scheduling, retention, verification, and atomic restore behavior from PART 22.
- [ ] Implement the mandatory TLS/certificate lifecycle surfaces from PART 15 where supported by the server deployment model.
- [ ] Complete the required client experience from PART 33, including TUI, shared theme, all API operations, shell completions, config, and failover discovery.
- [ ] Complete required service management and privilege handling from PARTS 24 and 25 for Linux, macOS, Windows, and BSD.
- [ ] Complete the update command and its stable/beta/devel, rollback, and admin notification behavior from PART 23.
- [ ] Implement mandatory Tor support and the documented container/host I2P behavior from PART 32.
- [ ] Finish the full server configuration surface, validation, and operator UI for every active setting in PARTS 5 and 12.

## Medium (P2)

- [ ] Bring README and every `docs/` page, example, OpenAPI document, and Compose sample into agreement with the completed implementation and IDEA.md.
- [ ] Complete CI/CD, release, security, license, and documentation gates for all configured providers and run the required project validation suite.
