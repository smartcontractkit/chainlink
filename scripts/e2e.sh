#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
go mod tidy
go run ./cmd/kinproof
echo "KinProof-Go E2E execution completed successfully."
