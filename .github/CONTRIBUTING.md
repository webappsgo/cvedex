# Contributing to cvedex

## Local setup and checks

Install Docker or Podman and GNU Make. The project builds and tests inside the
`casjaysdev/go` toolchain container; Python documentation builds use the pinned
dependencies in `docs/requirements.txt`.

```sh
make build
make test
docker run --rm --name "cvedex-docs-$(tr -dc 'a-z0-9' </dev/urandom | head -c8)" \
  -v "$PWD:/docs" -w /docs squidfunk/mkdocs-material:latest \
  build --strict --site-dir /tmp/cvedex-site
```

Open a focused pull request against `main`. Include the reason for the change,
test results, documentation and configuration updates, compatibility impact,
and security or privacy impact. Update tests and documentation with behavior
changes. Do not include credentials or unredacted private data in issues,
logs, or screenshots.

For support, start with the project documentation or `/server/contact`. Report
security vulnerabilities privately using the process in
[SECURITY.md](SECURITY.md); do not open a public issue for a vulnerability.
