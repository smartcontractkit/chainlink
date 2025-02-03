#!/usr/bin/env bash

# Helper to setup Solana contracts for testing
# Note: this may stop working if the structure of the chainlink-ccip repository changes
set -e

GIT_ROOT=$(git rev-parse --show-toplevel)
if [[ -z "$GIT_ROOT" ]]; then
  echo "Error: This script must be run within a Git repository."
  exit 1
fi

GO_MOD_FILE="$GIT_ROOT/deployment/go.mod"
TEMP_REPO_DIR=$(mktemp -d)
TARGET_DIR="$GIT_ROOT/deployment/ccip/changeset/internal/solana_contracts"

cleanup() {
  echo "Cleaning up..."
  rm -rf "$TEMP_REPO_DIR"
}
trap cleanup EXIT

if [[ ! -f "$GO_MOD_FILE" ]]; then
  echo "Error: go.mod file not found in the current directory."
  exit 1
fi

# Parse the go.mod file for the specific entry
MOD_ENTRY=$(grep -E 'github\.com/smartcontractkit/chainlink-ccip/chains/solana\s+v[0-9]+\.[0-9]+\.[0-9]+-[0-9]+-[a-f0-9]+' "$GO_MOD_FILE")
if [[ -z "$MOD_ENTRY" ]]; then
  echo "Error: Could not find the required entry in go.mod."
  exit 1
fi

# Extract repo URL and pseudo-version
PSEUDO_VERSION=$(echo "$MOD_ENTRY" | awk '{print $2}')

# Extract commit SHA from pseudo-version (last 12 characters)
COMMIT_SHA=$(echo "$PSEUDO_VERSION" | grep -oE '[a-f0-9]{12}$')
if [[ -z "$COMMIT_SHA" ]]; then
  echo "Error: Could not extract commit SHA from pseudo-version: $PSEUDO_VERSION"
  exit 1
fi

echo "Cloning chainlink-ccip..."
git clone https://github.com/smartcontractkit/chainlink-ccip.git "$TEMP_REPO_DIR"

echo "Checking out commit $COMMIT_SHA..."
cd "$TEMP_REPO_DIR"
git checkout "$COMMIT_SHA"

echo "Building Solana contracts..."
cd chains/solana
make docker-build-contracts

mkdir -p "$TARGET_DIR"

echo "Copying compiled artifacts to $TARGET_DIR..."
cp -r contracts/target/deploy/* "$TARGET_DIR/"

echo "Script completed successfully."