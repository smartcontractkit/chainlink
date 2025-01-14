#!/bin/bash

# Get the root project directory
DEPLOYMENT_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../../../" && pwd)
INTERNAL_FOLDER="$DEPLOYMENT_ROOT/ccip/changeset/internal"
TEMP_FOLDER="$INTERNAL_FOLDER/.tmp-solana-ccip-repo"

cd $DEPLOYMENT_ROOT

# extract the chainlink-ccip revision from go.mod
CCIP_VERSION=$(grep "github.com/smartcontractkit/chainlink-ccip/chains/solana" $DEPLOYMENT_ROOT/go.mod | awk '{print $2}' | cut -d'-' -f3)

echo "Chainlink-CCIP version: $CCIP_VERSION"

if [ -z "$CCIP_VERSION" ]; then
    echo "Error: Could not find chainlink-ccip dependency in go.mod"
    exit 1
fi

# remove the existing chainlink-ccip repo if it exists
rm -rf $TEMP_FOLDER

# clone the chainlink-ccip and checkout the required revision
git clone git@github.com:smartcontractkit/chainlink-ccip.git $TEMP_FOLDER
cd $TEMP_FOLDER
git checkout $CCIP_VERSION

# build the solana contracts
cd $TEMP_FOLDER/chains/solana/contracts
anchor build

# copy the built programs to the tests directory
cp target/deploy/*.so $INTERNAL_FOLDER/solana_contracts

# cleanup
rm -rf $TEMP_FOLDER