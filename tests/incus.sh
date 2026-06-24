#!/usr/bin/env bash
# Integration tests using Incus (Debian, full systemd — PREFERRED).
set -eo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
BINARY="${PROJECT_DIR}/binaries/cvedex"

log() { echo "[incus.sh] $(date '+%Y-%m-%dT%H:%M:%S%z') $*"; }

if ! command -v incus >/dev/null 2>&1; then
    echo "ERROR: incus is not available. Install Incus or use tests/docker.sh instead."
    exit 1
fi

if [ ! -f "${BINARY}" ]; then
    echo "ERROR: Binary not found at ${BINARY}"
    echo "Run 'make local' first to build the binary."
    exit 1
fi

CONTAINER_NAME="cvedex-test-$$"
log "Launching Incus container: ${CONTAINER_NAME}"

cleanup() {
    incus delete --force "${CONTAINER_NAME}" 2>/dev/null || true
}
trap cleanup EXIT

incus launch images:debian/12 "${CONTAINER_NAME}"

log "Waiting for container to be ready..."
sleep 5

log "Copying binary to container"
incus file push "${BINARY}" "${CONTAINER_NAME}/usr/local/bin/cvedex"
incus exec "${CONTAINER_NAME}" -- chmod +x /usr/local/bin/cvedex

log "Running --help test"
incus exec "${CONTAINER_NAME}" -- /usr/local/bin/cvedex --help

log "Running --version test"
incus exec "${CONTAINER_NAME}" -- /usr/local/bin/cvedex --version

log "All Incus integration tests passed"
