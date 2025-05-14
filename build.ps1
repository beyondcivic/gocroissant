#!/usr/bin/env pwsh

# Set environment variables for static build
$env:CGO_ENABLED = "1"
$env:GOOS = "linux"
$APP_VERSION = "25.05.dev"
$GIT_HASH = $(git log -1 --format='%H')
$BUILD_TIME = $(Get-Date -UFormat '%Y-%m-%dT%H:%M:%SZ')

# Create out directory if it doesn't exist
New-Item -ItemType Directory -Force -Path "bin"

# Build static binary with version information
$env:CGO_ENABLED = "1"
$env:GOOS = "linux"
$env:GOARCH = "amd64"

go build -v `
    -ldflags "-X 'version.Version=$APP_VERSION' -X 'version.GitHash=$GIT_HASH' -X 'version.BuildTime=$BUILD_TIME'" `
    -o ./bin/gocroissant ./cmd/gocroissant
