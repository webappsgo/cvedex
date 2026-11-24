pipeline {
    agent none

    triggers {
        cron('0 3 * * *\n0 4 * * *')
    }

    options {
        disableConcurrentBuilds()
        timeout(time: 30, unit: 'MINUTES')
        timestamps()
    }

    environment {
        PROJECT_NAME = 'cvedex'
        PROJECT_ORG = 'casapps'
        INTERNAL_NAME = 'cvedex'
        CGO_ENABLED = '0'
        GOFLAGS = '-buildvcs=false'
        GO_CACHE = "${JENKINS_HOME}/cache/casapps/cvedex/go-mod"
        GO_BUILD = "${JENKINS_HOME}/cache/casapps/cvedex/go-build"
        GOMODCACHE = '/go/pkg/mod'
        GOCACHE = '/root/.cache/go-build'
    }

    stages {
        stage('Checkout') {
            agent { label 'amd64' }
            steps {
                checkout([$class: 'GitSCM', branches: scm.branches, extensions: [[$class: 'CloneOption', depth: 0, noTags: false, shallow: false]], userRemoteConfigs: scm.userRemoteConfigs])
            }
        }

        stage('Lint') {
            agent { docker { image 'casjaysdev/go:latest'; label 'amd64'; args "-v /var/run/docker.sock:/var/run/docker.sock -v ${env.JENKINS_HOME}/cache/casapps/cvedex/go-mod:/go/pkg/mod -v ${env.JENKINS_HOME}/cache/casapps/cvedex/go-build:/root/.cache/go-build" } }
            steps {
                checkout scm
                sh 'mkdir -p "$GO_CACHE" "$GO_BUILD"'
                sh 'go vet ./...'
                sh 'staticcheck ./...'
            }
        }

        stage('Test and coverage') {
            agent { docker { image 'casjaysdev/go:latest'; label 'amd64'; args "-v /var/run/docker.sock:/var/run/docker.sock -v ${env.JENKINS_HOME}/cache/casapps/cvedex/go-mod:/go/pkg/mod -v ${env.JENKINS_HOME}/cache/casapps/cvedex/go-build:/root/.cache/go-build" } }
            steps {
                checkout scm
                sh '''
                    set -eu
                    mkdir -p /tmp/casapps
                    COVDIR=$(mktemp -d /tmp/casapps/cvedex-XXXXXX)
                    go test -cover -coverprofile="$COVDIR/coverage.out" ./...
                    test -s "$COVDIR/coverage.out" || { echo "coverage profile missing"; exit 1; }
                    PCT=$(go tool cover -func="$COVDIR/coverage.out" | awk '/^total:/ {gsub("%", "", $3); print $3}')
                    test -n "$PCT" || { echo "coverage total missing"; exit 1; }
                    awk -v pct="$PCT" 'BEGIN { if (pct < 60) { printf "coverage %.1f%% < threshold 60%%\\n", pct; exit 1 } }'
                '''
            }
        }

        stage('Build') {
            agent { docker { image 'casjaysdev/go:latest'; label 'amd64'; args "-v /var/run/docker.sock:/var/run/docker.sock -v ${env.JENKINS_HOME}/cache/casapps/cvedex/go-mod:/go/pkg/mod -v ${env.JENKINS_HOME}/cache/casapps/cvedex/go-build:/root/.cache/go-build" } }
            steps {
                checkout scm
                sh 'go build -buildvcs=false -trimpath ./...'
            }
        }

        stage('Security') {
            parallel {
                stage('Secret scan') {
                    agent { label 'amd64' }
                    steps {
                        checkout([$class: 'GitSCM', branches: scm.branches, extensions: [[$class: 'CloneOption', depth: 0, noTags: false, shallow: false]], userRemoteConfigs: scm.userRemoteConfigs])
                        script {
                            env.IS_SCHEDULED = currentBuild.getBuildCauses().any { cause -> cause._class == 'hudson.triggers.TimerTrigger$TimerTriggerCause' } ? 'true' : 'false'
                        }
                        sh '''
                            set -eu
                            BASE=""
                            if [ -n "${CHANGE_TARGET:-}" ]; then
                                git fetch --no-tags origin "$CHANGE_TARGET"
                                BASE=$(git merge-base "origin/$CHANGE_TARGET" HEAD)
                            elif [ "${IS_SCHEDULED:-false}" != true ] && [ -n "${GIT_PREVIOUS_SUCCESSFUL_COMMIT:-}" ]; then
                                BASE="$GIT_PREVIOUS_SUCCESSFUL_COMMIT"
                            fi
                            if [ -n "$BASE" ]; then
                                docker run --rm -v "$WORKSPACE:/repo:ro" -w /repo trufflesecurity/trufflehog:3.100.1@sha256:9c7a8325d43a4af2960d3f21d61a08981d8a1bf0542e007efa2e47097e386452 git file:///repo --since-commit "$BASE" --results=verified,unknown --exclude-paths=.github/trufflehog-exclude.txt --fail
                            else
                                docker run --rm -v "$WORKSPACE:/repo:ro" -w /repo trufflesecurity/trufflehog:3.100.1@sha256:9c7a8325d43a4af2960d3f21d61a08981d8a1bf0542e007efa2e47097e386452 git file:///repo --results=verified --exclude-paths=.github/trufflehog-exclude.txt --fail
                            fi
                        '''
                    }
                }
                stage('Go vulnerability scan') {
                    agent { docker { image 'casjaysdev/go:latest'; label 'amd64'; args "-v ${env.JENKINS_HOME}/cache/casapps/cvedex/go-mod:/go/pkg/mod -v ${env.JENKINS_HOME}/cache/casapps/cvedex/go-build:/root/.cache/go-build" } }
                    steps {
                        checkout scm
                        sh 'govulncheck ./...'
                    }
                }
                stage('Workflow policy') {
                    agent { docker { image 'casjaysdev/go:latest'; label 'amd64'; args "-v ${env.JENKINS_HOME}/cache/casapps/cvedex/go-mod:/go/pkg/mod -v ${env.JENKINS_HOME}/cache/casapps/cvedex/go-build:/root/.cache/go-build" } }
                    steps {
                        checkout scm
                        sh '''
                            set -eu
                            if grep -R -n 'pull_request_target' .github/workflows .gitea/workflows .forgejo/workflows; then exit 1; fi
                            if grep -R -nE '^ +uses: +[^@[:space:]]+@' .github/workflows .gitea/workflows .forgejo/workflows | grep -vE '@[0-9a-f]{40} +# v[0-9]'; then exit 1; fi
                            if grep -R -nE '^ +uses: +[^@]+@[0-9a-f]{40}($| +#)' .github/workflows .gitea/workflows .forgejo/workflows | grep -v '# v[0-9]'; then exit 1; fi
                        '''
                    }
                }
                stage('Image scan') {
                    agent { label 'amd64' }
                    steps {
                        checkout scm
                        sh 'docker build --tag cvedex:ci --file docker/Dockerfile .'
                        sh 'docker run --rm -v /var/run/docker.sock:/var/run/docker.sock aquasec/trivy:0.70.0@sha256:be1190afcb28352bfddc4ddeb71470835d16462af68d310f9f4bca710961a41e image --exit-code 1 --severity CRITICAL,HIGH cvedex:ci'
                    }
                }
            }
        }

        stage('Release binaries') {
            when { buildingTag() }
            parallel {
                stage('amd64 artifacts') {
                    agent { docker { image 'casjaysdev/go:latest'; label 'amd64'; args "-v ${env.JENKINS_HOME}/cache/casapps/cvedex/go-mod:/go/pkg/mod -v ${env.JENKINS_HOME}/cache/casapps/cvedex/go-build:/root/.cache/go-build" } }
                    steps {
                        checkout scm
                        sh '''
                            set -eu
                            mkdir -p binaries
                            VERSION="${TAG_NAME#v}"
                            COMMIT_ID=$(git rev-parse --short=7 HEAD)
                            BUILD_EPOCH=$(date -u +%s)
                            OFFICIAL_SITE=$(cat site.txt 2>/dev/null || true)
                            LDFLAGS="-s -w -X main.Version=$VERSION -X main.CommitID=$COMMIT_ID -X main.BuildEpoch=$BUILD_EPOCH -X main.OfficialSite=$OFFICIAL_SITE"
                            for OS in linux darwin windows freebsd; do
                                OUTPUT="binaries/cvedex-${OS}-amd64"
                                [ "$OS" != windows ] || OUTPUT="$OUTPUT.exe"
                                GOOS="$OS" GOARCH=amd64 go build -trimpath -buildvcs=false -ldflags "$LDFLAGS" -o "$OUTPUT" ./src
                                sha256sum "$OUTPUT" > "$OUTPUT.sha256"
                                if [ -d src/client ]; then
                                    CLI_OUTPUT="binaries/cvedex-${OS}-amd64-cli"
                                    [ "$OS" != windows ] || CLI_OUTPUT="$CLI_OUTPUT.exe"
                                    GOOS="$OS" GOARCH=amd64 go build -trimpath -buildvcs=false -ldflags "$LDFLAGS" -o "$CLI_OUTPUT" ./src/client
                                    sha256sum "$CLI_OUTPUT" > "$CLI_OUTPUT.sha256"
                                fi
                            done
                        '''
                        stash name: 'release-amd64', includes: 'binaries/**'
                    }
                }
                stage('arm64 artifacts') {
                    agent { docker { image 'casjaysdev/go:latest'; label 'arm64'; args "-v ${env.JENKINS_HOME}/cache/casapps/cvedex/go-mod:/go/pkg/mod -v ${env.JENKINS_HOME}/cache/casapps/cvedex/go-build:/root/.cache/go-build" } }
                    steps {
                        checkout scm
                        sh '''
                            set -eu
                            mkdir -p binaries
                            VERSION="${TAG_NAME#v}"
                            COMMIT_ID=$(git rev-parse --short=7 HEAD)
                            BUILD_EPOCH=$(date -u +%s)
                            OFFICIAL_SITE=$(cat site.txt 2>/dev/null || true)
                            LDFLAGS="-s -w -X main.Version=$VERSION -X main.CommitID=$COMMIT_ID -X main.BuildEpoch=$BUILD_EPOCH -X main.OfficialSite=$OFFICIAL_SITE"
                            for OS in linux darwin windows freebsd; do
                                OUTPUT="binaries/cvedex-${OS}-arm64"
                                [ "$OS" != windows ] || OUTPUT="$OUTPUT.exe"
                                GOOS="$OS" GOARCH=arm64 go build -trimpath -buildvcs=false -ldflags "$LDFLAGS" -o "$OUTPUT" ./src
                                sha256sum "$OUTPUT" > "$OUTPUT.sha256"
                                if [ -d src/client ]; then
                                    CLI_OUTPUT="binaries/cvedex-${OS}-arm64-cli"
                                    [ "$OS" != windows ] || CLI_OUTPUT="$CLI_OUTPUT.exe"
                                    GOOS="$OS" GOARCH=arm64 go build -trimpath -buildvcs=false -ldflags "$LDFLAGS" -o "$CLI_OUTPUT" ./src/client
                                    sha256sum "$CLI_OUTPUT" > "$CLI_OUTPUT.sha256"
                                fi
                            done
                        '''
                        stash name: 'release-arm64', includes: 'binaries/**'
                    }
                }
            }
        }

        stage('Package stable release') {
            when { buildingTag() }
            agent { docker { image 'casjaysdev/go:latest'; label 'amd64'; args "-v ${env.JENKINS_HOME}/cache/casapps/cvedex/go-mod:/go/pkg/mod -v ${env.JENKINS_HOME}/cache/casapps/cvedex/go-build:/root/.cache/go-build" } }
            steps {
                checkout scm
                unstash 'release-amd64'
                unstash 'release-arm64'
                sh '''
                    set -eu
                    (cd binaries && sha256sum cvedex-* > SHA256SUMS)
                    cyclonedx-gomod app -json -output binaries/cvedex-sbom.cdx.json ./src
                    git archive --format=tar.gz --prefix="cvedex-${TAG_NAME}/" HEAD > "binaries/cvedex-${TAG_NAME}-source.tar.gz"
                '''
                stash name: 'stable-release', includes: 'binaries/**'
                archiveArtifacts artifacts: 'binaries/**', fingerprint: true, onlyIfSuccessful: true
            }
        }

        stage('Publish stable release') {
            when { buildingTag() }
            agent { label 'amd64' }
            steps {
                checkout scm
                unstash 'stable-release'
                withCredentials([string(credentialsId: 'github-token', variable: 'GH_TOKEN')]) {
                    sh '''
                        set -eu
                        command -v gh >/dev/null || { echo "GitHub CLI is required on the Jenkins amd64 agent"; exit 1; }
                        TAG="$TAG_NAME"
                        AUTH_HEADER="AUTHORIZATION: basic $(printf 'x-access-token:%s' "$GH_TOKEN" | base64 | tr -d '\\n')"
                        REMOTE_TAG=$(git -c "http.extraheader=$AUTH_HEADER" ls-remote origin "refs/tags/$TAG" | awk 'NR == 1 {print $1}')
                        if [ -n "$REMOTE_TAG" ]; then
                            if gh release view "$TAG" >/dev/null 2>&1; then gh release delete "$TAG" --yes; fi
                            git -c "http.extraheader=$AUTH_HEADER" push origin ":refs/tags/$TAG"
                        fi
                        git tag -f "$TAG" HEAD
                        git -c "http.extraheader=$AUTH_HEADER" push origin "refs/tags/$TAG"
                        GH_TOKEN="$GH_TOKEN" gh release create "$TAG" binaries/* --title "cvedex $TAG" --generate-notes --verify-tag
                    '''
                }
            }
        }

        stage('Publish beta and devel releases') {
            when {
                anyOf {
                    branch 'beta'
                    branch 'main'
                    branch 'master'
                    triggeredBy 'TimerTrigger'
                }
            }
            agent { label 'amd64' }
            steps {
                checkout scm
                withCredentials([string(credentialsId: 'github-token', variable: 'GH_TOKEN')]) {
                    sh '''
                        set -eu
                        command -v gh >/dev/null || { echo "GitHub CLI is required on the Jenkins amd64 agent"; exit 1; }
                        mkdir -p binaries
                        COMMIT_ID=$(git rev-parse --short=7 HEAD)
                        BUILD_EPOCH=$(date -u +%s)
                        OFFICIAL_SITE=$(cat site.txt 2>/dev/null || true)
                        if [ "${BRANCH_NAME:-}" = beta ]; then
                            VERSION="$(date -u +%Y%m%d%H%M%S)-beta"
                            RELEASE_TAG="$VERSION"
                            PRERELEASE=--prerelease
                        else
                            VERSION="$COMMIT_ID"
                            RELEASE_TAG=devel
                            PRERELEASE=--prerelease
                            if gh release view devel >/dev/null 2>&1; then gh release delete devel --yes --cleanup-tag; fi
                            REMOTE_TAG=$(git ls-remote origin refs/tags/devel | awk 'NR == 1 {print $1}')
                            if [ -n "$REMOTE_TAG" ]; then
                                AUTH_HEADER="AUTHORIZATION: basic $(printf 'x-access-token:%s' "$GH_TOKEN" | base64 | tr -d '\\n')"
                                git -c "http.extraheader=$AUTH_HEADER" push origin :refs/tags/devel
                            fi
                        fi
                        LDFLAGS="-s -w -X main.Version=$VERSION -X main.CommitID=$COMMIT_ID -X main.BuildEpoch=$BUILD_EPOCH -X main.OfficialSite=$OFFICIAL_SITE"
                        for PLATFORM in linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64 windows/arm64 freebsd/amd64 freebsd/arm64; do
                            OS=${PLATFORM%/*}; ARCH=${PLATFORM#*/}; OUTPUT="binaries/cvedex-$OS-$ARCH"
                            [ "$OS" != windows ] || OUTPUT="$OUTPUT.exe"
                            GOOS=$OS GOARCH=$ARCH CGO_ENABLED=0 docker run --rm -e GOOS -e GOARCH -e CGO_ENABLED -v "$WORKSPACE:/src" -w /src casjaysdev/go:latest sh -lc "go build -trimpath -buildvcs=false -ldflags '$LDFLAGS' -o '$OUTPUT' ./src"
                            sha256sum "$OUTPUT" > "$OUTPUT.sha256"
                            if [ -d src/client ]; then
                                CLI_OUTPUT="binaries/cvedex-$OS-$ARCH-cli"
                                [ "$OS" != windows ] || CLI_OUTPUT="$CLI_OUTPUT.exe"
                                GOOS=$OS GOARCH=$ARCH CGO_ENABLED=0 docker run --rm -e GOOS -e GOARCH -e CGO_ENABLED -v "$WORKSPACE:/src" -w /src casjaysdev/go:latest sh -lc "go build -trimpath -buildvcs=false -ldflags '$LDFLAGS' -o '$CLI_OUTPUT' ./src/client"
                                sha256sum "$CLI_OUTPUT" > "$CLI_OUTPUT.sha256"
                            fi
                        done
                        (cd binaries && sha256sum cvedex-* > SHA256SUMS)
                        docker run --rm -v "$WORKSPACE:/src" -w /src casjaysdev/go:latest sh -lc 'cyclonedx-gomod app -json -output binaries/cvedex-sbom.cdx.json ./src'
                        git archive --format=tar.gz --prefix="cvedex-${RELEASE_TAG}/" HEAD > "binaries/cvedex-${RELEASE_TAG}-source.tar.gz"
                        if [ "$RELEASE_TAG" = devel ]; then
                            AUTH_HEADER="AUTHORIZATION: basic $(printf 'x-access-token:%s' "$GH_TOKEN" | base64 | tr -d '\\n')"
                            git tag -f devel HEAD
                            git -c "http.extraheader=$AUTH_HEADER" push origin refs/tags/devel
                        fi
                        GH_TOKEN="$GH_TOKEN" gh release create "$RELEASE_TAG" binaries/* --title "cvedex $RELEASE_TAG" --generate-notes $PRERELEASE --verify-tag
                    '''
                    archiveArtifacts artifacts: 'binaries/**', fingerprint: true, onlyIfSuccessful: true
                }
            }
        }

        stage('Build Docker images') {
            when { not { changeRequest() } }
            agent { label 'amd64' }
            steps {
                checkout scm
                withCredentials([string(credentialsId: 'github-token', variable: 'REGISTRY_TOKEN')]) {
                    sh '''
                        set -eu
                        BUILD_EPOCH=$(date -u +%s)
                        BUILD_DATE=$(date -u +%Y-%m-%dT%H:%M:%SZ)
                        COMMIT_ID=$(git rev-parse --short=7 HEAD)
                        IMAGE=ghcr.io/casapps/cvedex
                        echo "$REGISTRY_TOKEN" | docker login ghcr.io -u webappsgo --password-stdin
                        docker buildx create --use || docker buildx use default
                        docker buildx build --platform linux/amd64,linux/arm64 --file docker/Dockerfile --build-arg VERSION="${TAG_NAME:-devel}" --build-arg BUILD_DATE="$BUILD_DATE" --build-arg BUILD_EPOCH="$BUILD_EPOCH" --build-arg COMMIT_ID="$COMMIT_ID" --tag "$IMAGE:$COMMIT_ID" --provenance=false --push .
                        docker run --rm aquasec/trivy:0.70.0@sha256:be1190afcb28352bfddc4ddeb71470835d16462af68d310f9f4bca710961a41e image --exit-code 1 --severity CRITICAL,HIGH "$IMAGE:$COMMIT_ID"
                        if [ "${BRANCH_NAME:-}" = beta ]; then
                            docker buildx build --platform linux/amd64,linux/arm64 --file docker/Dockerfile --build-arg VERSION=beta --build-arg BUILD_DATE="$BUILD_DATE" --build-arg BUILD_EPOCH="$BUILD_EPOCH" --build-arg COMMIT_ID="$COMMIT_ID" --tag "$IMAGE:beta" --provenance=false --push .
                        fi
                        if [ -n "${TAG_NAME:-}" ]; then
                            VERSION="${TAG_NAME#v}"
                            docker buildx build --platform linux/amd64,linux/arm64 --file docker/Dockerfile --build-arg VERSION="$VERSION" --build-arg BUILD_DATE="$BUILD_DATE" --build-arg BUILD_EPOCH="$BUILD_EPOCH" --build-arg COMMIT_ID="$COMMIT_ID" --tag "$IMAGE:$VERSION" --tag "$IMAGE:latest" --tag "$IMAGE:$(date -u +%y%m)" --provenance=false --push .
                        fi
                        if [ "${BRANCH_NAME:-}" = main ] || [ "${BRANCH_NAME:-}" = master ]; then
                            docker buildx build --platform linux/amd64,linux/arm64 --file docker/Dockerfile.dev --build-arg VERSION=devel --build-arg BUILD_DATE="$BUILD_DATE" --build-arg BUILD_EPOCH="$BUILD_EPOCH" --build-arg COMMIT_ID="$COMMIT_ID" --tag "$IMAGE:devel" --provenance=false --push .
                        fi
                    '''
                }
            }
        }
    }

    post {
        always { cleanWs() }
    }
}
