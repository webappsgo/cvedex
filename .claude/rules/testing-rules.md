# Testing & Documentation Rules (PART 29, 30, 31)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- ❌ Write test output or coverage files to the project tree
- ❌ Commit with test coverage below 60%
- ❌ Build or run tests on the host — use Docker
- ❌ Skip the test gate before committing
- ❌ Hardcode language strings in Go source — use i18n JSON files
- ❌ Omit the English (en) i18n baseline — it is the required fallback
- ❌ Create docs outside docs/ directory

## CRITICAL - ALWAYS DO
- ✅ Run tests inside Docker: go test -cover ./... (casjaysdev/go:latest)
- ✅ Maintain test coverage ≥ 60% (enforced in CI)
- ✅ Write coverage output to /tmp/{project_org}/{project_name}-XXXXXX/coverage.out
- ✅ Support --lang flag on cvedex-cli and cvedex-agent
- ✅ Embed i18n JSON files in binary via go:embed
- ✅ Fall back to English (en) when requested language is unavailable
- ✅ WCAG 2.1 AA on all HTML pages
- ✅ Maintain ReadTheDocs docs in docs/ with mkdocs.yml

## TESTING (PART 29)
| Item | Requirement |
|------|-------------|
| Test runner | go test -cover ./... inside Docker |
| Coverage minimum | 60% (CI enforced, PR blocked below threshold) |
| Output location | /tmp/casapps/cvedex-XXXXXX/ (never project tree) |
| Test types | Unit tests for handlers, services, models; integration tests for API routes |

## READTHEDOCS DOCUMENTATION (PART 30)
| Section | Content |
|---------|---------|
| Getting Started | Installation, first run, quick start |
| Configuration | server.yml reference, all options |
| DNS Records | Supported record types, zone file format |
| JSON API | All endpoints with curl examples |
| DNSSEC Setup | Key generation, signing, DS record publishing |
| Admin Guide | Dashboard, sync, user management |
| Self-Hosting | Docker, systemd, port 53 setup |
| Contributing | Dev setup, PR process, conventions |

Docs location: docs/ with mkdocs.yml at project root.

## I18N (PART 31)
| Item | Requirement |
|------|-------------|
| Flag | --lang on cvedex-cli and cvedex-agent |
| File format | JSON key-value files |
| Location | src/data/i18n/{lang}.json |
| Embedding | go:embed — all files in binary |
| Fallback | Always fall back to en (English) |
| Required baseline | en must always be complete |

```
src/data/i18n/
├── en.json    # Required — complete baseline
├── es.json
├── fr.json
└── de.json
```

## ACCESSIBILITY (PART 31)
- Standard: WCAG 2.1 AA
- Touch targets: minimum 44x44px
- Color contrast: minimum 4.5:1 ratio for normal text
- Keyboard navigation: all interactive elements reachable via keyboard
- Screen reader: semantic HTML, ARIA labels where needed

---
For complete details, see AI.md PART 29, 30, 31
