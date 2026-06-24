# Infer PROJECTNAME and PROJECTORG from git remote or directory path (NEVER hardcode)
PROJECTNAME := $(shell git remote get-url origin 2>/dev/null | sed -E 's|.*/([^/]+)(\.git)?$$|\1|' || basename "$$(pwd)")
PROJECTORG := $(shell git remote get-url origin 2>/dev/null | sed -E 's|.*/([^/]+)/[^/]+(\.git)?$$|\1|' || basename "$$(dirname "$$(pwd)")")

# Version precedence: release.txt > env/default fallback
VERSION ?= $(shell cat release.txt 2>/dev/null || echo "devel")

# Build info — ISO 8601 / RFC 3339 UTC
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
COMMIT_ID := $(shell git rev-parse --short HEAD 2>/dev/null || echo "N/A")

# Official site URL (OPTIONAL — sourced from site.txt if present, else env, else empty)
OFFICIALSITE := $(shell [ -f site.txt ] && cat site.txt || echo "${OFFICIALSITE:-}")

# Linker flags to embed build info
LDFLAGS := -s -w \
	-X 'main.Version=$(VERSION)' \
	-X 'main.CommitID=$(COMMIT_ID)' \
	-X 'main.BuildDate=$(BUILD_DATE)' \
	-X 'main.OfficialSite=$(OFFICIALSITE)'

# Directories
BINDIR := binaries
RELDIR := releases

# Go module cache dirs (bind-mounted from host for speed)
GO_CACHE  ?= $(HOME)/go/pkg/mod
GO_BUILD  ?= $(HOME)/.cache/go-build

# Build targets
PLATFORMS ?= linux/amd64,linux/arm64

# Docker — all builds happen inside casjaysdev/go:latest
REGISTRY ?= ghcr.io/$(PROJECTORG)/$(PROJECTNAME)
GO_DOCKER := docker run --rm -it \
	--name $(PROJECTNAME)-$$(tr -dc 'a-z0-9' </dev/urandom | head -c8) \
	--entrypoint "" \
	-v $(PWD):/app \
	-v $(GO_CACHE):/usr/local/share/go/pkg/mod \
	-v $(GO_BUILD):/usr/local/share/go/cache \
	-w /app \
	-e CGO_ENABLED=0 \
	casjaysdev/go:latest

.PHONY: build local release docker test dev clean

# =============================================================================
# BUILD - Build all platforms + local binary (via Docker with cached modules)
# =============================================================================
build: clean
	@mkdir -p $(BINDIR) $(GO_CACHE) $(GO_BUILD)
	@echo "Building version $(VERSION)..."

	# Tidy and download modules
	@echo "Tidying and downloading Go modules..."
	@$(GO_DOCKER) sh -c "go mod tidy && go mod download"

	# Build for local OS/ARCH
	@echo "Building local binary..."
	@$(GO_DOCKER) sh -c "GOOS=$$(go env GOOS) GOARCH=$$(go env GOARCH) \
		go build -ldflags \"$(LDFLAGS)\" -o $(BINDIR)/$(PROJECTNAME) ./src"

	# Build server for all platforms
	@for platform in $$(echo $(PLATFORMS) | tr ',' ' '); do \
		OS=$${platform%/*}; \
		ARCH=$${platform#*/}; \
		OUTPUT=$(BINDIR)/$(PROJECTNAME)-$$OS-$$ARCH; \
		[ "$$OS" = "windows" ] && OUTPUT=$$OUTPUT.exe; \
		echo "Building server $$OS/$$ARCH..."; \
		$(GO_DOCKER) sh -c "GOOS=$$OS GOARCH=$$ARCH \
			go build -ldflags \"$(LDFLAGS)\" \
			-o $$OUTPUT ./src" || exit 1; \
	done

	# Build CLI for all platforms (if exists)
	@if [ -d "src/client" ]; then \
		for platform in $$(echo $(PLATFORMS) | tr ',' ' '); do \
			OS=$${platform%/*}; \
			ARCH=$${platform#*/}; \
			OUTPUT=$(BINDIR)/$(PROJECTNAME)-cli-$$OS-$$ARCH; \
			[ "$$OS" = "windows" ] && OUTPUT=$$OUTPUT.exe; \
			echo "Building CLI $$OS/$$ARCH..."; \
			$(GO_DOCKER) sh -c "GOOS=$$OS GOARCH=$$ARCH \
				go build -ldflags \"$(LDFLAGS)\" \
				-o $$OUTPUT ./src/client" || exit 1; \
		done; \
	fi

	@echo "Build complete: $(BINDIR)/"

# =============================================================================
# LOCAL - Build local binaries only (fast development builds)
# =============================================================================
local: clean
	@mkdir -p $(BINDIR) $(GO_CACHE) $(GO_BUILD)
	@echo "Building local binaries version $(VERSION)..."

	# Tidy and download modules
	@echo "Tidying and downloading Go modules..."
	@$(GO_DOCKER) sh -c "go mod tidy && go mod download"

	# Build server binary
	@echo "Building $(PROJECTNAME)..."
	@$(GO_DOCKER) sh -c "GOOS=$$(go env GOOS) GOARCH=$$(go env GOARCH) \
		go build -ldflags \"$(LDFLAGS)\" -o $(BINDIR)/$(PROJECTNAME) ./src"

	@echo "Local build complete: $(BINDIR)/"

# =============================================================================
# RELEASE - Manual local release (stable only)
# =============================================================================
release: build
	@mkdir -p $(RELDIR)
	@echo "Preparing release $(VERSION)..."

	# Create version.txt
	@echo "$(VERSION)" > $(RELDIR)/version.txt

	# Copy binaries to releases (strip if needed)
	@for f in $(BINDIR)/$(PROJECTNAME)-*; do \
		[ -f "$$f" ] || continue; \
		strip "$$f" 2>/dev/null || true; \
		cp "$$f" $(RELDIR)/; \
	done

	# Create source archive (exclude VCS and build artifacts)
	@tar --exclude='.git' --exclude='.github' --exclude='.gitea' \
		--exclude='binaries' --exclude='releases' --exclude='*.tar.gz' \
		-czf $(RELDIR)/$(PROJECTNAME)-$(VERSION)-source.tar.gz .

	# Delete existing release/tag if exists
	@gh release delete $(VERSION) --yes 2>/dev/null || true
	@git tag -d $(VERSION) 2>/dev/null || true
	@git push origin :refs/tags/$(VERSION) 2>/dev/null || true

	# Create new release (stable)
	@gh release create $(VERSION) $(RELDIR)/* \
		--title "$(PROJECTNAME) $(VERSION)" \
		--notes "Release $(VERSION)" \
		--latest

	@echo "Release complete: $(VERSION)"

# =============================================================================
# DOCKER - Build container image (no push — CI/CD handles pushing)
# =============================================================================
docker:
	@echo "Building Docker image $(VERSION)..."

	# Ensure buildx is available
	@docker buildx version > /dev/null 2>&1 || (echo "docker buildx required" && exit 1)

	# Create/use builder
	@docker buildx create --name $(PROJECTNAME)-builder --use 2>/dev/null || \
		docker buildx use $(PROJECTNAME)-builder

	# Build multi-arch locally (no push — pushing is CI/CD's responsibility)
	@docker buildx build \
		-f docker/Dockerfile \
		--platform linux/amd64,linux/arm64 \
		--build-arg VERSION="$(VERSION)" \
		--build-arg BUILD_DATE="$(BUILD_DATE)" \
		--build-arg COMMIT_ID="$(COMMIT_ID)" \
		-t $(REGISTRY):$(VERSION) \
		-t $(REGISTRY):latest \
		.

	@echo "Docker build complete: $(REGISTRY):$(VERSION)"

# =============================================================================
# TEST - Run all tests with coverage enforcement (via Docker)
# =============================================================================
test:
	@echo "Running tests with coverage..."
	@$(GO_DOCKER) sh -c "go mod download && \
		go test -v -cover -coverprofile=coverage.out ./... && \
		COVERAGE=\$$(go tool cover -func=coverage.out | grep total | awk '{print \$$3}' | sed 's/%//'); \
		if [ \$$(echo \"\$$COVERAGE < 80\" | bc -l) -eq 1 ]; then \
			echo \"ERROR: Coverage is \$$COVERAGE%, must be >= 80%\"; \
			exit 1; \
		fi && \
		echo \"Tests complete - Coverage: \$$COVERAGE% (>= 80% required)\""

# =============================================================================
# DEV - Quick build for local development/testing (to random temp dir)
# =============================================================================
dev:
	@$(GO_DOCKER) sh -c "go mod tidy"
	@mkdir -p "$${TMPDIR:-/tmp}/$(PROJECTORG)" && \
		BUILD_DIR=$$(mktemp -d "$${TMPDIR:-/tmp}/$(PROJECTORG)/$(PROJECTNAME)-XXXXXX") && \
		echo "Quick dev build to $$BUILD_DIR..." && \
		$(GO_DOCKER) sh -c "go build -o $$BUILD_DIR/$(PROJECTNAME) ./src" && \
		echo "Built: $$BUILD_DIR/$(PROJECTNAME)" && \
		echo "Test: docker run --rm -it --name $(PROJECTNAME)-test -v $$BUILD_DIR:/app alpine:latest /app/$(PROJECTNAME) --help"

# =============================================================================
# CLEAN - Remove build artifacts
# =============================================================================
clean:
	@rm -rf $(BINDIR) $(RELDIR)
