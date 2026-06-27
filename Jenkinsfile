pipeline {
    agent {
        docker {
            image 'casjaysdev/go:latest'
            args '-v /var/run/docker.sock:/var/run/docker.sock'
        }
    }

    environment {
        CGO_ENABLED = '0'
        GOFLAGS = '-buildvcs=false'
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Lint') {
            steps {
                sh 'go vet ./...'
                sh 'staticcheck ./... || true'
            }
        }

        stage('Test') {
            steps {
                sh '''
                    mkdir -p /tmp/casapps
                    COVDIR=$(mktemp -d /tmp/casapps/cvedex-XXXXXX)
                    go test -v -cover -coverprofile="$COVDIR/coverage.out" ./...
                    COVERAGE=$(go tool cover -func="$COVDIR/coverage.out" | grep total | awk '{print $3}' | sed 's/%//')
                    echo "Coverage: ${COVERAGE}%"
                '''
            }
        }

        stage('Build') {
            steps {
                sh '''
                    VERSION=$(cat release.txt 2>/dev/null || echo "devel")
                    BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
                    COMMIT_ID=$(git rev-parse --short HEAD 2>/dev/null || echo "N/A")
                    LDFLAGS="-s -w -X main.Version=$VERSION -X main.CommitID=$COMMIT_ID -X main.BuildDate=$BUILD_DATE"
                    mkdir -p binaries
                    go build -ldflags "$LDFLAGS" -buildvcs=false -o binaries/cvedex ./src
                '''
            }
        }

        stage('Release') {
            when {
                tag 'v*'
            }
            steps {
                sh '''
                    VERSION=$(cat release.txt 2>/dev/null || echo "devel")
                    BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
                    COMMIT_ID=$(git rev-parse --short HEAD 2>/dev/null || echo "N/A")
                    LDFLAGS="-s -w -X main.Version=$VERSION -X main.CommitID=$COMMIT_ID -X main.BuildDate=$BUILD_DATE"
                    for PLATFORM in linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64 windows/arm64 freebsd/amd64 freebsd/arm64; do
                        OS=${PLATFORM%/*}
                        ARCH=${PLATFORM#*/}
                        OUTPUT="binaries/cvedex-${OS}-${ARCH}"
                        [ "$OS" = "windows" ] && OUTPUT="${OUTPUT}.exe"
                        GOOS=$OS GOARCH=$ARCH go build -ldflags "$LDFLAGS" -buildvcs=false -o "$OUTPUT" ./src
                    done
                '''
            }
        }
    }

    post {
        always {
            cleanWs()
        }
    }
}
