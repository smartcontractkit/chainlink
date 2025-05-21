#!/usr/bin/env bash

# Helper to setup Solana contracts for testing
# Note: this may stop working if the structure of the chainlink-ccip repository changes
set -e

CDIR="$(dirname "$0")"

DEPLOYMENT_DIR=$(realpath "$CDIR/..") # deployment/
TARGET_DIR="$DEPLOYMENT_DIR/ccip/changeset/internal/solana_contracts"
mkdir -p "$TARGET_DIR"

SOLANA_DIR=$(go list -m -f "{{.Dir}}" github.com/smartcontractkit/chainlink-ccip/chains/solana)
cd $SOLANA_DIR
TARGET_DIR=${TARGET_DIR} make docker-build-contracts

# copy to top level
cp -r $TARGET_DIR/deploy/* $TARGET_DIR

echo "Script completed successfully."
echo "Artifacts are in $TARGET_DIR"
