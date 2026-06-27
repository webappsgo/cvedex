# cvedex

Read `AI.md` and `IDEA.md` before acting on this project.

## FIRST TURN — MANDATORY

On every new conversation or after "context compacted":
1. Read the relevant `.claude/rules/*.md` for your current task
2. Never assume or guess — verify against AI.md before implementing

## Before ANY Code Change

1. Have I read the relevant PART in AI.md? (If no → read it)
2. Does this follow the spec EXACTLY? (If unsure → check spec)
3. Am I guessing or do I KNOW from the spec? (If guessing → read spec)
4. Would this pass the compliance checklist? (AI.md FINAL section)

**WHEN IN DOUBT: READ THE SPEC. DO NOT GUESS.**

## Key Facts

- Binary: `cvedex` (DNS+HTTP server), `cvedex-cli` (CLI client), `cvedex-agent` (optional)
- Module: `github.com/casapps/cvedex`
- CGO_ENABLED=0 always — pure Go, no CGO ever
- Password hashing: Argon2id only (never bcrypt)
- SQLite driver: modernc.org/sqlite only (never mattn/go-sqlite3)
- Dockerfile: `docker/Dockerfile` (never project root)
- Source: `src/` — singular Go dir names (handler/, model/, service/)

## Placeholders

- `{project_name}` = cvedex
- `{project_org}` = casapps
- `{internal_name}` = cvedex

## File Locations

| Purpose | Privileged | User |
|---------|-----------|------|
| Config | /etc/casapps/cvedex/ | ~/.config/casapps/cvedex/ |
| Data | /var/lib/casapps/cvedex/ | ~/.local/share/casapps/cvedex/ |
| Logs | /var/log/casapps/cvedex/ | ~/.local/log/casapps/cvedex/ |

## Where to Find Details

- Full spec: `AI.md` (SOURCE OF TRUTH)
- Rules by topic: `.claude/rules/*.md`
- Project variables: `IDEA.md`
- Outstanding tasks: `TODO.AI.md`
