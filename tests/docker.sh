#!/usr/bin/env bash
# Integration tests using Docker (alpine:latest).
set -eo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
BINARY="${PROJECT_DIR}/binaries/cvedex"

log() { echo "[docker.sh] $(date '+%Y-%m-%dT%H:%M:%S%z') $*"; }

if [ ! -f "${BINARY}" ]; then
    echo "ERROR: Binary not found at ${BINARY}"
    echo "Run 'make local' first to build the binary."
    exit 1
fi

CONTAINER_NAME="cvedex-test-$$"
log "Starting test container: ${CONTAINER_NAME}"

cleanup() {
    docker rm -f "${CONTAINER_NAME}" 2>/dev/null || true
}
trap cleanup EXIT

docker run --rm -d \
    --name "${CONTAINER_NAME}" \
    -v "${BINARY}:/app/cvedex:ro" \
    alpine:latest \
    sleep 60

log "Running --help test"
docker exec "${CONTAINER_NAME}" /app/cvedex --help

log "Running --version test"
docker exec "${CONTAINER_NAME}" /app/cvedex --version

log "All Docker integration tests passed"
