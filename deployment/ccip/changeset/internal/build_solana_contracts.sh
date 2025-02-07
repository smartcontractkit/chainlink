#!/bin/bash

set -euo pipefail

# Note: this version of the script expects anchor to be installed

# Get the root project directory
ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")" && while [[ $PWD != "/" && ${PWD##*/} != "chainlink" ]]; do cd ..; done && pwd)
DEPLOYMENT_FOLDER="$ROOT/deployment"
INTERNAL_FOLDER="$DEPLOYMENT_FOLDER/ccip/changeset/internal"
TEMP_FOLDER="$INTERNAL_FOLDER/.tmp-solana-ccip-repo"
CONTRACTS_FOLDER=$INTERNAL_FOLDER/solana_contracts

cd $DEPLOYMENT_FOLDER

# extract the chainlink-ccip revision from go.mod
CCIP_VERSION=$(grep "github.com/smartcontractkit/chainlink-ccip/chains/solana" $DEPLOYMENT_FOLDER/go.mod | awk '{print $2}' | cut -d'-' -f3)

echo "chainlink-CCIP version: $CCIP_VERSION"

if [ -z "$CCIP_VERSION" ]; then
    echo "Error: Could not find chainlink-ccip dependency in go.mod"
    exit 1
fi

# cleanup the existing chainlink-ccip repo if it exists
rm -rf $TEMP_FOLDER
rm -rf $CONTRACTS_FOLDER

mkdir -p $CONTRACTS_FOLDER

echo $CCIP_VERSION > $CONTRACTS_FOLDER/.solana_contracts_rev

# clone the chainlink-ccip and checkout the required revision
git clone https://github.com/smartcontractkit/chainlink-ccip.git $TEMP_FOLDER
cd $TEMP_FOLDER
git checkout $CCIP_VERSION

# build the solana contracts
cd $TEMP_FOLDER/chains/solana/contracts

anchor build

# copy the built programs to the tests directory
cp target/deploy/*.so $CONTRACTS_FOLDER

# cleanup
rm -rf $TEMP_FOLDER
