#!/usr/bin/env bash
set -euo pipefail

echo "=== Fixing directory structure ==="

mkdir -p cmd/kinproof internal/attribution scripts

# Ensure go.mod is correct
cat <<'GO' > go.mod
module kinproof-go
go 1.23

require (
	github.com/ethereum/go-ethereum v1.14.12
	github.com/tyler-smith/go-bip32 v1.0.0
	github.com/tyler-smith/go-bip39 v1.1.0
)
GO

echo "=== Structure fixed. Trying to run now ==="

chmod +x scripts/run.sh 2>/dev/null || true
go mod tidy
go run ./cmd/kinproof
