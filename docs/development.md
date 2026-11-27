# Development

The project is written in Go. Follow the repository's `AI.md` for required architecture and conventions; it is the source of truth for implementation.

## Toolchain

Use Docker with `casjaysdev/go:latest` for all Go work. Do not install Go on the host.

## Build and checks

```sh
make local
make test
CVEDEX_RUN_INCUS=1 make test
```

Test and lint commands are run through the session gate wrapper described in the repository's global instructions. `tests/run_tests.sh` runs the basic integration tier in Docker and does not fall back to Incus.

Set `CVEDEX_RUN_INCUS=1` to include the full-OS Incus tier in `make test`; set `CVEDEX_INCUS_MODE=systemd` or `openrc` to run one init system, or leave it unset for `all`. Otherwise run `tests/incus.sh [systemd|openrc|all]` explicitly. It exercises Debian Trixie with systemd and Alpine 3.24 with OpenRC, including service install, stop, start, restart, reload, and persistence after an Incus container restart. Each environment checks the HTTP health route and authoritative DNS behavior. The test uses the server binary already built by `make test`; for a standalone run, it uses `binaries/cvedex` unless `CVEDEX_TEST_BINARY` is set. Incus operations have time limits: defaults are 60 seconds per operation, 900 seconds for package installation, and 3600 seconds for the full run. `INCUS_OP_TIMEOUT`, `INCUS_INSTALL_TIMEOUT`, and `INCUS_RUN_TIMEOUT` override those limits; CPU and memory are capped at 2 CPUs and 4 GiB by default and can be adjusted with `INCUS_CPUS` and `INCUS_MEM`.

## Data package

`src/cvedata` parses CVE records and enrichment snapshots, validates schema-versioned datasets, and reads/writes deterministic gzip snapshots. `src/dnsserver` serves DNS responses. `src/server` provides the HTTP gateway and rendered pages. `src/service` coordinates listeners and shutdown.

## Contributions

Keep `AI.md` read-only. Update `IDEA.md` only when the product definition changes, and keep public docs synchronized with behavior. Add tests for behavior changes and do not commit generated binaries or runtime data.
