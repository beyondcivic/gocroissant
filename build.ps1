#!/usr/bin/env pwsh

# Set environment variables for static build
$env:CGO_ENABLED = "1"
$env:GOOS = "linux"
$APP_VERSION = "$(git describe --tags --abbrev=0 2>/dev/null || echo '0.0.0')"

# Create out directory if it doesn't exist
New-Item -ItemType Directory -Force -Path "bin"

# Build static binary with version information
$env:CGO_ENABLED = "1"
$env:GOOS = "linux"
$env:GOARCH = "amd64"

go build -v `
    -ldflags "-X 'github.com/beyondcivic/gocroissant/pkg/version.Version=$APP_VERSION'" `
    -o ./bin/gocroissant .
