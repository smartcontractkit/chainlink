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

# Modern solc objects have metadata suffixes which vary depending on
# incidental compilation context like absolute paths to source files. See
# https://solidity.readthedocs.io/en/v0.6.2/metadata.html#encoding-of-the-metadata-hash-in-the-bytecode
# Since this suffix varies so much, it can't be included in a reliable check
# that the golang wrapper is up-to-date, so remove it from the message hash.
BINLEN="${#BIN}"
TRUNCLEN="$((BINLEN - 106))" # 106/2=53=length of metadata hash in bytes
TRUNCATED="${BIN:0:$TRUNCLEN}"
SUFFIX="${BIN:$TRUNCLEN:106}" # The actual metadata hash, in hex.

# Verify that the suffix follows the pattern outlined in the above link, to
# ensure that we're actually truncating what we think we are.
SUFFIX_REGEXP='^a264697066735822[[:xdigit:]]{68}64736f6c6343[[:xdigit:]]{6}0033$'
if [[ ! $SUFFIX =~ $SUFFIX_REGEXP ]]; then
    echo "binary suffix has unexpected format; giving up"
    exit 1
fi

CONSTANT_SUFFIX="a264697066735822000000000000000000000000000000000000000000000"
CONSTANT_SUFFIX+="0000000000000000000000064736f6c63430000000033"

BIN_PATH="${TMP_DIR}/bin"
echo "${TRUNCATED}${CONSTANT_SUFFIX}" > "$BIN_PATH"

OUT_PATH="${TMP_DIR}/$PKG_NAME.go"

"$CDIR"/abigen.sh "$BIN_PATH" "$ABI_PATH" "$OUT_PATH" "$CLASS_NAME" "$PKG_NAME"

TARGET_DIR="./generated/$PKG_NAME/"
mkdir -p "$TARGET_DIR"
cp "$OUT_PATH" "$TARGET_DIR"

"$CDIR/record_versions.sh" "$@" "$ABI" "$BIN"
