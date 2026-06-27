# Binary Rules (PART 7, 8, 33)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- ❌ Enable CGO (CGO_ENABLED=0 always, no exceptions)
- ❌ Build on the host — use Docker (casjaysdev/go:latest)
- ❌ Add short flags other than -h and -v
- ❌ Omit any required CLI flag from the list below
- ❌ Make cvedex-cli optional — it is REQUIRED
- ❌ Use external assets — all assets must be embedded via go:embed
- ❌ Add musl suffix to binary names

## CRITICAL - ALWAYS DO
- ✅ CGO_ENABLED=0 on every build command
- ✅ Build via Docker with casjaysdev/go:latest
- ✅ Embed all assets (templates, static files, i18n) with go:embed
- ✅ Ship three binaries: cvedex, cvedex-cli (required), cvedex-agent (optional)
- ✅ Implement ALL required CLI flags on each binary
- ✅ Support -h (short for --help) and -v (short for --version) ONLY as short flags
- ✅ Use GOFLAGS=-buildvcs=false when building in Docker with volume mount

## SERVER BINARY (cvedex) — REQUIRED FLAGS
| Flag | Description |
|------|-------------|
| --help / -h | Show help |
| --version / -v | Show version |
| --mode | Run mode (production/development) |
| --config | Config file path |
| --data | Data directory path |
| --log | Log directory path |
| --pid | PID file path |
| --address | Listen address |
| --port | Listen port |
| --baseurl | Base URL for generated links |
| --debug | Enable debug logging |
| --status | Show service status |
| --service {start,restart,stop,reload,--install,--uninstall,--disable,--help} | Service management |
| --daemon | Run as background daemon |
| --maintenance | Enable maintenance mode |
| --update | Check for and apply updates |

## CLIENT BINARY (cvedex-cli) — REQUIRED
| Flag/Command | Description |
|-------------|-------------|
| lookup | Look up a CVE by ID |
| search | Search CVEs by keyword |
| kev | List/query CISA KEV entries |
| zone | Fetch DNS zone file |
| --server | Server URL |
| --output | Output format (json/text/table) |
| --help / -h | Show help |
| --version / -v | Show version |
| --debug | Enable debug output |
| --color | Enable/disable color output |
| --lang | Language for output |

## AGENT BINARY (cvedex-agent) — OPTIONAL
| Flag | Description |
|------|-------------|
| --server | Server URL to connect to |
| --token | Authentication token |
| --config | Config file path |
| --help / -h | Show help |
| --version / -v | Show version |
| --debug | Enable debug output |
| --color | Enable/disable color output |
| --lang | Language for output |

---
For complete details, see AI.md PART 7, 8, 33
