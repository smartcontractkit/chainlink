#!/usr/bin/env bash

set -e

# Record versions of each contract, to avoid use of stale wrappers
#
# Usage: record_versions.sh <json-artifact-path> <golang-wrapper-package-name> \
#                           <abi-string> <binary-object-string> \
#                           [<don't-truncate>]
#
# The db is a flat file, one record per line. There is one line of the form
# "GETH_VERSION: <current-go-ethereum-version". The other lines are of the form
#
#   <golang-wrapper-package-name>: <json-artifact-path> <hash-of-json-artifact>
#
# with each contract path followed by the current hash of that contract. These
# are checked in the golang test TestCheckContractHashesFromLastGoGenerate, to
# ensure that the golang wrappers are current.
#
# If something is passed in the <don't-truncate> slot, no attempt is made to
# remove the trailing metadata in the binary object.

SOL_PATH="$1"
PKG_NAME="$2"
ABI="$3"
BIN="$4"
DONT_TRUNCATE="$5"

if [[ $SOL_PATH =~ [[:space:]] ]]; then
    # The golang parser splits on whitespace, so don't allow it in the pathname
    echo "path to compiler artifact, '$SOL_PATH', cannot contain whitespace"
    exit 1
fi

CDIR=$(dirname "$0")
# shellcheck source=common.sh
source "$CDIR/common.sh"

VERSION_DB_PATH="$CDIR/generated-wrapper-dependency-versions-do-not-edit.txt"
touch "$VERSION_DB_PATH"

function blow_away_version_record() {
    TGT_RECORD="$1"
    (grep -v "$TGT_RECORD": "$VERSION_DB_PATH" > "$VERSION_DB_PATH.tmp") || true
    mv "$VERSION_DB_PATH.tmp" "$VERSION_DB_PATH"
}

blow_away_version_record GETH_VERSION

# go.mod geth version is of form v1.9.9. Strip leading v.
echo GETH_VERSION: "${GETH_VERSION//v/}" >> "$VERSION_DB_PATH"

blow_away_version_record "$PKG_NAME"

if [ -n "$DONT_TRUNCATE" ]; then # Caller has asked us not to trucate binary
    MSG_BIN="$BIN"
else
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
    MSG_BIN="$TRUNCATED"
fi

HASHMSG="$ABI$MSG_BIN"
echo "$PKG_NAME: $SOL_PATH $(sha256sum <<< "$HASHMSG" | cut -f 1 -d ' ')" >> \
     "$VERSION_DB_PATH"
sort -o "$VERSION_DB_PATH" "$VERSION_DB_PATH"
