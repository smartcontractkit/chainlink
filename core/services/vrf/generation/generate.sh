#!/usr/bin/env bash

set -e

# Usage:
#
#   ./generate.sh <sol-compiler-output-path> <package-name>
#
# This will output the generated file to ./<package-name>/<package-name>.go,
# where ./<package-name> is in the same directory as this script

SOL_PATH="$1"
PKG_NAME="$2"

function cleanup() { # Release resources on script exit
    rm -rf "$TMP_DIR"
}
trap cleanup EXIT

CDIR=$(dirname "$0")
# shellcheck source=common.sh
source "$CDIR/common.sh"

CLASS_NAME=$(basename "$SOL_PATH" .json)

TMP_DIR=$(mktemp -d "/tmp/${CLASS_NAME}.XXXXXXXXX")

ABI=$(jq -c -r .compilerOutput.abi < "$SOL_PATH")
ABI_PATH="${TMP_DIR}/abi.json"
echo "$ABI" > "$ABI_PATH"

# We want the bytecode here, not the deployedByteCode. The latter does not
# include the initialization code.
# https://ethereum.stackexchange.com/questions/32234/difference-between-bytecode-and-runtime-bytecode
BIN=$(jq -r .compilerOutput.evm.bytecode.object < "$SOL_PATH")
BIN_PATH="${TMP_DIR}/bin"
echo "$BIN" > "$BIN_PATH"

OUT_PATH="${TMP_DIR}/$PKG_NAME.go"

"$CDIR"/abigen.sh "$BIN_PATH" "$ABI_PATH" "$OUT_PATH" "$CLASS_NAME" "$PKG_NAME"

TARGET_DIR="./generated/$PKG_NAME/"
mkdir -p "$TARGET_DIR"
cp "$OUT_PATH" "$TARGET_DIR"

"$CDIR/record_versions.sh" "$@" "$ABI" "$BIN"
