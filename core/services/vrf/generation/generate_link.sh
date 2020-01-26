#!/usr/bin/env bash

set -e

function cleanup() { # Release resources on script exit
    rm -rf "$TMP_DIR"
}
trap cleanup EXIT

CDIR=$(dirname "$0")

TMP_DIR=$(mktemp -d /tmp/link_token.XXXXXXXXX)

LINK_COMPILER_ARTIFACT_PATH="$CDIR/../../../../evm/src/LinkToken.json"

ABI_PATH="${TMP_DIR}/abi.json"
jq .compilerOutput.abi < "$LINK_COMPILER_ARTIFACT_PATH" > "$ABI_PATH"

BIN_PATH="${TMP_DIR}/bin"
jq .bytecode < "$LINK_COMPILER_ARTIFACT_PATH" | tr -d '"' > "$BIN_PATH"

CLASS_NAME="LinkToken"
PKG_NAME="link_token_interface"
OUT_PATH="$TMP_DIR/$PKG_NAME.go"

"$CDIR"/abigen.sh "$BIN_PATH" "$ABI_PATH" "$OUT_PATH" "$CLASS_NAME" "$PKG_NAME"

TARGET_DIR="./generated/$PKG_NAME/"
mkdir -p "$TARGET_DIR"
"$CDIR/record_versions.sh" "$LINK_COMPILER_ARTIFACT_PATH" link_token_interface
