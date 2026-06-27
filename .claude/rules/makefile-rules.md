# Makefile Rules (PART 26)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO
- ❌ Use Makefile in CI/CD — CI runs its own steps directly
- ❌ Build on the host — always use Docker (casjaysdev/go:latest)
- ❌ Use $(pwd) in docker -v flags — use $PWD (statically analyzable)
- ❌ Omit GOFLAGS=-buildvcs=false — required for Docker volume mounts
- ❌ Write coverage or test output to the project tree
- ❌ Use cd in Makefile rules without absolute paths

## CRITICAL - ALWAYS DO
- ✅ The Makefile is for LOCAL DEV ONLY — it is never invoked in CI/CD
- ✅ All builds go through Docker (casjaysdev/go:latest)
- ✅ Use $PWD in docker -v flags (not $(pwd))
- ✅ Set GOFLAGS=-buildvcs=false in GO_DOCKER
- ✅ Bind-mount GO_CACHE and GO_BUILD from host for build caching
- ✅ Write test/coverage output to /tmp/{project_org}/{project_name}-XXXXXX/
- ✅ Include all required targets

## REQUIRED MAKE TARGETS
| Target | Purpose |
|--------|---------|
| make dev | Quick build to temp dir for local testing |
| make local | Build single binary for current OS/arch |
| make build | Build all 8 platform binaries |
| make test | Run tests inside Docker |
| make release | Create GitHub release with all artifacts |
| make docker | Build Docker image |
| make clean | Remove build artifacts |

## DOCKER BUILD PATTERN
```makefile
GO_DOCKER = docker run --rm \
    -e GOFLAGS=-buildvcs=false \
    -e CGO_ENABLED=0 \
    -v $PWD:/app \
    -v $(GO_CACHE):/root/.cache/go-build \
    -v $(GO_BUILD):/go/pkg/mod \
    -w /app \
    casjaysdev/go:latest
```

## TEST OUTPUT PATH
```makefile
test:
    @COVDIR=$$(docker run --rm -v $PWD:/app casjaysdev/go:latest \
        sh -c 'mkdir -p /tmp/casapps && mktemp -d /tmp/casapps/cvedex-XXXXXX' ) && \
    $(GO_DOCKER) go test -cover -coverprofile=$$COVDIR/coverage.out ./...
```

Coverage and test output always go to /tmp/{project_org}/{project_name}-XXXXXX/ — NEVER to the project tree.

---
For complete details, see AI.md PART 26
