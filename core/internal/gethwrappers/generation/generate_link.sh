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

ABI=$(jq -c -r .abi < "$LINK_COMPILER_ARTIFACT_PATH")
ABI_PATH="${TMP_DIR}/abi.json"
echo "$ABI" > "$ABI_PATH"

BIN=$(jq -r .bytecode < "$LINK_COMPILER_ARTIFACT_PATH")
BIN_PATH="${TMP_DIR}/bin"
echo "$BIN" > "$BIN_PATH"

CLASS_NAME="LinkToken"
PKG_NAME="link_token_interface"
OUT_PATH="$TMP_DIR/$PKG_NAME.go"

"$CDIR"/abigen.sh "$BIN_PATH" "$ABI_PATH" "$OUT_PATH" "$CLASS_NAME" "$PKG_NAME"

TARGET_DIR="./generated/$PKG_NAME/"
mkdir -p "$TARGET_DIR"
cp "$OUT_PATH" "$TARGET_DIR"
"$CDIR/record_versions.sh" "$LINK_COMPILER_ARTIFACT_PATH" link_token_interface \
                           "$ABI" "$BIN" dont_truncate_binary
