#!/usr/bin/env bash

# Usage:
#
#   ./generate.sh <sol-compiler-output-path> <package-name>
#
# This will output the generated file to ./<package-name>/<package-name>.go,
# where ./<package-name> is in the same directory as this script

set -e

function cleanup() { # Release resources on script exit
    rm -rf "$TMP_DIR"
    docker rm -f "$DOCKER_CONTAINER_NAME" > /dev/null
}
trap cleanup EXIT

CDIR=$(dirname "$0")

CLASS_NAME=$(basename "$1" .json)

TMP_DIR=$(mktemp -d "/tmp/${CLASS_NAME}.XXXXXXXXX")

ABI_PATH="${TMP_DIR}/abi.json"
jq .compilerOutput.abi < "$1" > "$ABI_PATH"

BIN_PATH=${TMP_DIR}/bin
jq .compilerOutput.evm.bytecode.object < "$1" | tr -d '"' > "$BIN_PATH"

GETH_VERSION=$(go list -json -m github.com/ethereum/go-ethereum | jq -r .Version)
CONTAINER_NAME_PATH="${TMP_DIR}/container_name"
DOCKER_BIN_PATH=/data/$(basename "$BIN_PATH")
DOCKER_ABI_PATH=/data/$(basename "$ABI_PATH")
DOCKER_OUT_PATH="/data/$2.go"

# Invoke abigen from the ethereum/go-client image to generate the golang wrapper
docker run -v "${TMP_DIR}:/data" --cidfile="$CONTAINER_NAME_PATH" \
       \
       "ethereum/client-go:alltools-$GETH_VERSION" \
       \
       abigen -bin "$DOCKER_BIN_PATH" \
              -abi "$DOCKER_ABI_PATH" \
              -out "$DOCKER_OUT_PATH" \
              -type "$CLASS_NAME" \
              -pkg "$2"

DOCKER_CONTAINER_NAME=$(cat "$CONTAINER_NAME_PATH")
if [ "$(docker wait "$DOCKER_CONTAINER_NAME")" != "0" ] ; then
    echo "Failed to build $CLASS_NAME golang wrapper"
    exit 1
fi

TARGET_DIR="./generated/$2/"
mkdir -p "$TARGET_DIR"
docker cp "$DOCKER_CONTAINER_NAME:$DOCKER_OUT_PATH" "$TARGET_DIR"

"$CDIR/record_versions.sh" "$@"
