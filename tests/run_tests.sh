#!/usr/bin/env bash
# Auto-detect available container runtime and run integration tests.
set -eo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

log() { echo "[run_tests] $(date '+%Y-%m-%dT%H:%M:%S%z') $*"; }

log "cvedex integration test runner"

# Prefer Incus (full OS, systemd support) over Docker
if command -v incus >/dev/null 2>&1; then
    log "Using Incus for integration tests"
    exec "${SCRIPT_DIR}/incus.sh" "$@"
elif command -v docker >/dev/null 2>&1; then
    log "Using Docker for integration tests"
    exec "${SCRIPT_DIR}/docker.sh" "$@"
else
    echo "ERROR: Neither incus nor docker is available. Install one to run integration tests."
    exit 1
fi
