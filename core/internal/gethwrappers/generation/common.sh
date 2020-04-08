#!/bin/bash

# GETH_VERSION is the version of go-ethereum chainlink is using
GETH_VERSION=$(go list -json -m github.com/ethereum/go-ethereum | jq -r .Version)
export GETH_VERSION
