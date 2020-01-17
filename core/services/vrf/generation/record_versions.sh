#!/usr/bin/env bash

set -e

# Record versions of each contract, to avoid use of stale wrappers
#
# Usage: record_versions.sh <json-artifact-path> <golang-wrapper-package-name>
#
# The db is a flat file, one record per line. There is one line of the form
# "GETH_VERSION: <current-go-ethereum-version". The other lines are of the form
#
#   <golang-wrapper-package-name>: <json-artifact-path> <hash-of-json-artifact>
#
# with each contract path followed by the current hash of that contract. These
# are checked in the golang test TestCheckContractHashesFromLastGoGenerate, to
# ensure that the golang wrappers are current.
CDIR=$(dirname "$0")
VERSION_DB_PATH="$CDIR/generated-wrapper-dependency-versions-do-not-edit.txt"
touch "$VERSION_DB_PATH"

function blow_away_version_record() {
    (grep -v "$1": "$VERSION_DB_PATH" > "$VERSION_DB_PATH.tmp") || true
    mv "$VERSION_DB_PATH.tmp" "$VERSION_DB_PATH"
}

blow_away_version_record GETH_VERSION

GETH_VERSION=$(go list -m github.com/ethereum/go-ethereum | awk '{print $2}')
# go.mod geth version is of form v1.9.9. Strip leading v.
echo GETH_VERSION: "${GETH_VERSION//v/}" >> "$VERSION_DB_PATH"

blow_away_version_record "$2"
echo "$2: $1 $(md5sum "$1" | cut -f 1 -d ' ')" | sort >> "$VERSION_DB_PATH"
