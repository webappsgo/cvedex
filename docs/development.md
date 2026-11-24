# Development

The project is written in Go. Follow the repository's `AI.md` for required architecture and conventions; it is the source of truth for implementation.

## Toolchain

Use Docker with `casjaysdev/go:latest` for all Go work. Do not install Go on the host.

## Build and checks

```sh
make local
make test
```

Test and lint commands are run through the session gate wrapper described in the repository's global instructions. `tests/run_tests.sh` runs the containerized integration test tier.

## Data package

`src/cvedata` parses CVE records and enrichment snapshots, validates schema-versioned datasets, and reads/writes deterministic gzip snapshots. `src/dnsserver` serves DNS responses. `src/server` provides the HTTP gateway and rendered pages. `src/service` coordinates listeners and shutdown.

## Contributions

Keep `AI.md` read-only. Update `IDEA.md` only when the product definition changes, and keep public docs synchronized with behavior. Add tests for behavior changes and do not commit generated binaries or runtime data.
