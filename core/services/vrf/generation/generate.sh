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
}
trap cleanup EXIT

CDIR=$(dirname "$0")
# shellcheck source=common.sh
source "$CDIR/common.sh"

CLASS_NAME=$(basename "$1" .json)

TMP_DIR=$(mktemp -d "/tmp/${CLASS_NAME}.XXXXXXXXX")

ABI_PATH="${TMP_DIR}/abi.json"
jq .compilerOutput.abi < "$1" > "$ABI_PATH"

BIN_PATH="${TMP_DIR}/bin"
jq .compilerOutput.evm.bytecode.object < "$1" | tr -d '"' > "$BIN_PATH"

OUT_PATH="${TMP_DIR}/$2.go"

PKG_NAME="$2"
"$CDIR"/abigen.sh "$BIN_PATH" "$ABI_PATH" "$OUT_PATH" "$CLASS_NAME" "$PKG_NAME"

TARGET_DIR="./generated/$PKG_NAME/"
mkdir -p "$TARGET_DIR"
cp "$OUT_PATH" "$TARGET_DIR"

"$CDIR/record_versions.sh" "$@"
