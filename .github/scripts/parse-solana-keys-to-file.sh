#!/bin/bash

set -euo pipefail

export RUSTUP_HOME="/root/.rustup"
export FORCE_COLOR=1

cd /solana/contracts

# Build and sync Anchor project
anchor keys sync
anchor build

# Extract program IDs and save to TOML file
anchor keys list | sed -E 's/ //g' | sed -E 's/([^:]*):*(.*)/\1 = "\2"/' > program_ids.toml

# Set permissions
chmod -R 755 ./target