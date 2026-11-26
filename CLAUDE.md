# Project SPEC

Project: CVEDEX
Role: Efficient loader for AI.md

⚠️ **THIS FILE IS AUTO-LOADED EVERY CONVERSATION. FOLLOW IT EXACTLY.** ⚠️

`AI.md` is the read-only source of truth for implementation rules. `IDEA.md`
defines cvedex's product behavior. Read the relevant parts before every task;
never guess or silently weaken a requirement.

## FIRST TURN - MANDATORY

On every new conversation or after context compaction:

1. Read the relevant `.claude/rules/*.md` files for the current task.
2. Verify implementation decisions against the corresponding `AI.md` parts.
3. Check `TODO.AI.md` and `TODO.md` if present; preserve human-owned tasks.

## Project-specific rules

- The server is authoritative for CVE data and serves the configured DNS zone.
- Read-only public data is exposed through DNS, HTTP pages, and the JSON API.
- Do not add user accounts or enable optional AI.md parts without `SPEC.md`.
- The required client binary is `cvedex-cli`; an agent is not in project scope.
- Use `casjaysdev/go:latest` through the Makefile for Go build and test commands.
- Keep README, operator docs, OpenAPI, configuration examples, and pages in sync.

## Key paths

- Project behavior: `IDEA.md`
- Full implementation requirements: `AI.md` (parts 0–36; part 37 is reference)
- Source: `src/`
- Container files: `docker/`
- Documentation: `docs/`

Never edit `AI.md`. Project-specific implementation changes belong in code,
documentation, or `IDEA.md` as appropriate.
