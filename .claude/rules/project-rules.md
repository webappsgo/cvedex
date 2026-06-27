# Project Rules (PART 2, 3, 4)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- ❌ Use any license other than MIT
- ❌ Build on the host — Docker only (casjaysdev/go:latest)
- ❌ Enable CGO (CGO_ENABLED=0 always)
- ❌ Use a musl suffix on binary names
- ❌ Put source outside ./src/
- ❌ Use plural directory names in Go packages (Go uses singular: handler/, model/, service/)
- ❌ Use camelCase filenames — snake_case only
- ❌ Put docker-compose.yml or Dockerfile at the project root
- ❌ Commit binaries/, releases/, or volumes/ — all are gitignored

## CRITICAL - ALWAYS DO
- ✅ MIT license in LICENSE.md
- ✅ Build all 8 platform targets
- ✅ CGO_ENABLED=0 on every build invocation
- ✅ Name binaries: `cvedex-{os}-{arch}` (windows adds .exe)
- ✅ Keep all Go source under ./src/
- ✅ Use singular directory names for Go packages
- ✅ Use snake_case for all filenames
- ✅ Keep docker/ directory for all Docker artifacts

## BUILD TARGETS (all 8 required)
| OS | Arch | Output Binary |
|----|------|---------------|
| linux | amd64 | cvedex-linux-amd64 |
| linux | arm64 | cvedex-linux-arm64 |
| darwin | amd64 | cvedex-darwin-amd64 |
| darwin | arm64 | cvedex-darwin-arm64 |
| windows | amd64 | cvedex-windows-amd64.exe |
| windows | arm64 | cvedex-windows-arm64.exe |
| freebsd | amd64 | cvedex-freebsd-amd64 |
| freebsd | arm64 | cvedex-freebsd-arm64 |

## DIRECTORY STRUCTURE
```
cvedex/
├── src/               # All Go source code
│   ├── cmd/           # Main entry points
│   ├── handler/       # HTTP handlers (singular — Go package)
│   ├── model/         # Data models (singular — Go package)
│   ├── service/       # Business logic (singular — Go package)
│   └── data/          # Embedded assets (i18n, templates)
├── docker/            # Dockerfile, docker-compose.yml, rootfs/
├── binaries/          # Build output (gitignored)
├── releases/          # Release artifacts (gitignored)
└── volumes/           # Docker volume mounts (gitignored)
```

## OS-SPECIFIC PATHS
| Mode | Config | Data | Logs |
|------|--------|------|------|
| Linux privileged | /etc/casapps/cvedex/ | /var/lib/casapps/cvedex/ | /var/log/casapps/cvedex/ |
| Linux user | ~/.config/casapps/cvedex/ | ~/.local/share/casapps/cvedex/ | ~/.local/log/casapps/cvedex/ |
| macOS | ~/Library/Application Support/casapps/cvedex/ | ~/Library/Application Support/casapps/cvedex/ | ~/Library/Logs/casapps/cvedex/ |
| Windows | %APPDATA%\casapps\cvedex\ | %APPDATA%\casapps\cvedex\ | %APPDATA%\casapps\cvedex\ |

---
For complete details, see AI.md PART 2, 3, 4
