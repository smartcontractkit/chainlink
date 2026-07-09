#!/usr/bin/env bash
#
# Builds the enginefleet Docker image.
#
# Compiles a Linux enginefleet binary inside a golang container (reusing the
# host Go module cache read-only, so no network or private-module auth is
# needed), then packages it into the slim runtime image defined by Dockerfile.
#
# Usage:
#   ./build.sh [-t image:tag]
#
set -euo pipefail

IMAGE="enginefleet:latest"
while getopts ":t:h" opt; do
  case "$opt" in
    t) IMAGE="$OPTARG" ;;
    h) echo "Usage: $0 [-t image:tag]"; exit 0 ;;
    *) echo "Usage: $0 [-t image:tag]" >&2; exit 1 ;;
  esac
done

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR" && git rev-parse --show-toplevel)"

# Go version must match the monorepo (go.mod) so the module cache is compatible.
GO_VERSION="$(awk '/^go / {print $2; exit}' "$REPO_ROOT/go.mod")"
GO_IMAGE="golang:${GO_VERSION}-bookworm"

GOMODCACHE="$(cd "$REPO_ROOT" && go env GOMODCACHE)"
BUILD_CACHE_VOL="enginefleet-go-build-cache"
PKG="./core/services/workflows/cmd/enginefleet"

echo ">> Compiling Linux binary with ${GO_IMAGE} (CGO, offline module cache)"
docker volume create "$BUILD_CACHE_VOL" >/dev/null
docker run --rm \
  -v "$REPO_ROOT":/src:ro \
  -v "$GOMODCACHE":/go/pkg/mod:ro \
  -v "$BUILD_CACHE_VOL":/root/.cache/go-build \
  -v "$SCRIPT_DIR":/out \
  -e CGO_ENABLED=1 \
  -e GOPROXY=off \
  -e GOFLAGS=-mod=mod \
  -w /src \
  "$GO_IMAGE" \
  go build -o /out/enginefleet "$PKG"

echo ">> Building image ${IMAGE}"
docker build -t "$IMAGE" -f "$SCRIPT_DIR/Dockerfile" "$SCRIPT_DIR"

echo ">> Done: ${IMAGE}"
