#!/usr/bin/env bash

set -e

# Usage: abigen.sh <args>. See the following assignments for the argument list.
# $BIN_PATH, $ABI_PATH and $OUT_PATH must all be in the same directory

BIN_PATH="$1"   # Path to the contract binary, as 0x-hex
ABI_PATH="$2"   # Path to the contract ABI JSON
OUT_PATH="$3"   # Path at which to save the golang contract wrapper
CLASS_NAME="$4" # Name for the golang contract wrapper type
PKG_NAME="$5"   # Name for the package containing the wrapper

COMMON_PARENT_DIRECTORY=$(dirname "$BIN_PATH")
if [ "$(dirname "$ABI_PATH")" != "$COMMON_PARENT_DIRECTORY" ] || \
       [ "$(dirname "$OUT_PATH")" != "$COMMON_PARENT_DIRECTORY" ]; then
    # shellcheck disable=SC2016
    echo '$BIN_PATH, $ABI_PATH and $OUT_PATH must all be in the same directory'
    exit 1
fi

function cleanup() {
    rm -rf "$CONTAINER_NAME_PATH"
    if [ ! -z "$DOCKER_CONTAINER_NAME" ]; then
        docker rm -f "$DOCKER_CONTAINER_NAME" > /dev/null
    fi
}
trap cleanup EXIT

CONTAINER_NAME_PATH="$(mktemp)"
rm -f "$CONTAINER_NAME_PATH"

CDIR=$(dirname "$0")
# shellcheck source=common.sh
source "$CDIR/common.sh"

ABIGEN_ARGS=( -bin "$BIN_PATH" -abi "$ABI_PATH" -out "$OUT_PATH"
              -type "$CLASS_NAME" -pkg "$PKG_NAME")

# Geth version from which native abigen was built, or empty string.
NATIVE_ABIGEN_VERSION=v"$( (
    abigen --version 2> /dev/null | \
    grep -E -o '([0-9]+\.[0-9]+\.[0-9]+)'
  ) || true
)"

# Generate golang wrapper
if [ "$NATIVE_ABIGEN_VERSION" == "$GETH_VERSION" ]; then
    abigen "${ABIGEN_ARGS[@]}" # We can use native abigen, which is much faster
else # Must use dockerized abigen
    DOCKER_IMAGE="ethereum/client-go:alltools-$GETH_VERSION"
    echo -n "Native abigen unavailable, broken, or wrong version (need version "
    echo "$GETH_VERSION). Invoking abigen from $DOCKER_IMAGE docker image."
    docker run -v "${COMMON_PARENT_DIRECTORY}:${COMMON_PARENT_DIRECTORY}" \
           --cidfile="$CONTAINER_NAME_PATH" \
           "$DOCKER_IMAGE" \
           abigen "${ABIGEN_ARGS[@]}"
    DOCKER_CONTAINER_NAME=$(cat "$CONTAINER_NAME_PATH")
    if [ "$(docker wait "$DOCKER_CONTAINER_NAME")" != "0" ] ; then
        echo "Failed to build $CLASS_NAME golang wrapper"
        exit 1
    fi
    docker cp "$DOCKER_CONTAINER_NAME:${OUT_PATH}" "${OUT_PATH}"
fi
