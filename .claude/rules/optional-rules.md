# Optional Features Rules (PART 34, 35, 36)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## STATUS: NOT IMPLEMENTED

**PARTS 34, 35, and 36 are OPTIONAL and are NOT implemented for cvedex.**

cvedex has no Regular User concept per IDEA.md — all CVE/DNS data is anonymous-public read-only. There are no user accounts, organizations, or custom domains in the base implementation.

## CRITICAL - NEVER DO
- ❌ Implement any PART 34/35/36 feature without explicit declaration in SPEC.md
- ❌ Add user registration, login, or account management without activation
- ❌ Add organization or team features without activation
- ❌ Add custom domain routing without activation
- ❌ Partially implement these features — all or nothing

## CRITICAL - ALWAYS DO
- ✅ Check SPEC.md before assuming these features are active
- ✅ If activated, implement 100% completely — no stubs or partial work
- ✅ If activated, all rules below become NON-NEGOTIABLE

## ACTIVATION REQUIREMENT
To activate any of these optional parts, the user must declare them in SPEC.md. Without that declaration, treat these PARTs as non-existent.

---

## MULTI-USER (PART 34) — OPTIONAL, NOT ACTIVE

If ever activated, these rules apply:

| Item | Requirement |
|------|-------------|
| Admin DB | Separate tables for admins vs regular users |
| Auth methods | Local password (Argon2id) + OIDC + LDAP |
| Sessions | Secure session management, configurable TTL |
| Password reset | Email-based (requires SMTP) |
| Rate limiting | Login attempts rate-limited per IP |

## ORGANIZATIONS (PART 35) — OPTIONAL, NOT ACTIVE

If ever activated, these rules apply:

| Item | Requirement |
|------|-------------|
| Org model | org_id foreign key on all tenant-scoped resources |
| Roles | Owner, Admin, Member at minimum |
| Invites | Email-based org invitations |
| Isolation | Orgs cannot see each other's data |

## CUSTOM DOMAINS (PART 36) — OPTIONAL, NOT ACTIVE

If ever activated, these rules apply:

| Item | Requirement |
|------|-------------|
| Routing | Per-tenant domain routing based on Host header |
| TLS | Separate ACME cert per custom domain |
| Verification | Domain ownership verification before activation |
| Isolation | Custom domain maps to exactly one org |

---
For complete details, see AI.md PART 34, 35, 36
