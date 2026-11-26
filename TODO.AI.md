# cvedex AI.md Alignment Tasks

## Critical (P0)

- [ ] Finish server page and API route coverage against IDEA.md and AI.md PARTS 13, 14, 16, and 17; ensure every public data endpoint has a working server-rendered page and page errors follow the shared negotiation and accessibility rules.
- [ ] Implement the mandatory GraphQL API and GraphiQL page in `src/graphql/`, mount both versioned and unversioned routes, and keep query fields synchronized with the read API contract.
- [ ] Complete admin authentication, settings, security, status, and audit-log routes and pages required by PART 17; align settings with all active cvedex configuration.

## High (P1)

- [ ] Implement the required shared theme, i18n, and accessibility systems from PARTS 16 and 31 across pages, admin, CLI, and health output.
- [ ] Implement mandatory metrics endpoints and per-service bearer authentication from PART 21.
- [ ] Implement required scheduler tasks and admin controls from PART 19, preserving cvedex's daily refresh behavior.
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
