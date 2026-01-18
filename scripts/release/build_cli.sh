#!/usr/bin/env bash
set -euo pipefail

PLATFORM=${1:?"platform (e.g., darwin/arm64) required"}
EXT=${2:-}

GOOS=$(echo "$PLATFORM" | cut -d'/' -f1)
GOARCH=$(echo "$PLATFORM" | cut -d'/' -f2)

echo "Building CLI for $GOOS/$GOARCH"
GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=0 go build -o "build/bin/vsynx${EXT}" -ldflags="-s -w" .
