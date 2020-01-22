#!/bin/bash

GETH_VERSION=$(go list -json -m github.com/ethereum/go-ethereum | jq -r .Version)
export GETH_VERSION
