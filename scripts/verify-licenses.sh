#!/usr/bin/env bash
# Verify all Go dependencies use compatible (non-copyleft) licenses.
# Must be run inside the project build image (docker/Dockerfile.build).
set -eo pipefail

echo "Checking for incompatible licenses..."

# Require go-licenses — never install inline; pre-installed in docker/Dockerfile.build.
command -v go-licenses >/dev/null 2>&1 || {
    echo "ERROR: go-licenses not found — run inside the project build image (docker/Dockerfile.build)"
    exit 1
}

echo "Scanning dependencies..."
if go-licenses csv ./... | grep -iE 'GPL|AGPL|LGPL'; then
    echo "ERROR: Copyleft license detected!"
    echo "Remove the dependency or find an alternative."
    exit 1
fi

echo "All licenses are compatible"

echo "Generating license report..."
go-licenses csv ./... > licenses.csv
go-licenses save ./... --save_path=third_party_licenses

echo "License report saved to licenses.csv and third_party_licenses/"
echo ""
echo "Next steps:"
echo "1. Review licenses.csv"
echo "2. Update LICENSE.md with any new dependencies"
echo "3. Commit the changes"
