#!/usr/bin/env bash

# Generates the golang wrapper of the LINK ERC20 token, which is represented by
# a non-standard compiler argument. Takes no arguments.

function cleanup() { # Release resources on script exit
    rm -rf "$TMP_DIR"
}
# trap cleanup EXIT

CDIR=$(dirname "$0")

TMP_DIR=$(mktemp -d /tmp/link_token.XXXXXXXXX)

LINK_COMPILER_ARTIFACT_PATH="$CDIR/../../../../evm-test-helpers/src/LinkToken.json"

ABI=$(cat "$LINK_COMPILER_ARTIFACT_PATH" | jq -c -r .abi)
ABI_PATH="${TMP_DIR}/abi.json"
echo "$ABI" > "$ABI_PATH"

BIN=$(cat "$LINK_COMPILER_ARTIFACT_PATH" | jq -r .bytecode)
BIN_PATH="${TMP_DIR}/bin"
echo "$BIN" > "$BIN_PATH"

CLASS_NAME="LinkToken"
PKG_NAME="link_token_interface"
OUT_PATH="$TMP_DIR/$PKG_NAME.go"

# shellcheck source=common.sh
source "$(dirname "$0")/common.sh"

ABIGEN_ARGS=( -bin "$BIN_PATH" -abi "$ABI_PATH" -out "$OUT_PATH"
              -type "$CLASS_NAME" -pkg "$PKG_NAME")

# Geth version from which native abigen was built, or v.
NATIVE_ABIGEN_VERSION=v"$(
    abigen --version 2> /dev/null | \
    grep -E -o '([0-9]+\.[0-9]+\.[0-9]+)'
)" || true

GETH_VERSION=$(go list -json -m github.com/ethereum/go-ethereum | jq -r .Version)

# Generate golang wrapper
if [ "$NATIVE_ABIGEN_VERSION" == "$GETH_VERSION" ]; then
    abigen "${ABIGEN_ARGS[@]}"
else
    echo "must install correct version of abigen"
    echo "(`make abigen` in the chainlink root dir)"
    exit 1
fi

TARGET_DIR="./generated/$PKG_NAME/"
mkdir -p "$TARGET_DIR"
cp "$OUT_PATH" "$TARGET_DIR"
