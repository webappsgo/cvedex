# Service Rules (PART 24, 25)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- ❌ Run the binary as OS root for port 53 — use setcap instead
- ❌ Use SUID bit to gain port binding privilege
- ❌ Request sudo unless --service install explicitly needs it
- ❌ Bundle service files as separate downloads — embed in binary
- ❌ Create OS-level service files by hand — generated and deployed by `cvedex --service install`
- ❌ Install service files without user confirmation

## CRITICAL - ALWAYS DO
- ✅ Embed all service file templates in the binary (go:embed)
- ✅ Generate and deploy service files via `cvedex --service install`
- ✅ Use setcap cap_net_bind_service=+ep on Linux for port 53 (NOT suid)
- ✅ Request sudo only when --service install actually needs it (file writes to system dirs)
- ✅ Support all three service managers
- ✅ Implement `cvedex --service {start,restart,stop,reload,--install,--uninstall,--disable,--help}`

## SERVICE FILE LOCATIONS
| Platform | Service Manager | File Path |
|----------|----------------|-----------|
| Linux | systemd | /etc/systemd/system/cvedex.service |
| macOS | launchd | /Library/LaunchDaemons/io.github.casapps.cvedex.plist |
| Windows | Windows Service | registered via golang.org/x/sys/windows/svc |

## LINUX (systemd)
```ini
[Unit]
Description=cvedex CVE Intelligence Server
After=network.target

[Service]
ExecStart=/usr/local/bin/cvedex --daemon
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

## MACOS (launchd)
- Label: io.github.casapps.cvedex
- Plist location: /Library/LaunchDaemons/io.github.casapps.cvedex.plist
- RunAtLoad: true
- KeepAlive: true

## WINDOWS (golang.org/x/sys/windows/svc)
- Service name: cvedex
- Display name: cvedex CVE Intelligence Server
- Start type: automatic

## PORT 53 PRIVILEGE (Linux only)
```bash
# Run after install — NOT suid
setcap cap_net_bind_service=+ep /usr/local/bin/cvedex
```

---
For complete details, see AI.md PART 24, 25
