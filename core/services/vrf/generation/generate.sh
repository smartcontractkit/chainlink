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
    if [ ! -z "$DOCKER_CONTAINER_NAME" ]; then
        docker rm -f "$DOCKER_CONTAINER_NAME" > /dev/null
    fi
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

CONTAINER_NAME_PATH="${TMP_DIR}/container_name"
OUT_PATH="${TMP_DIR}/$2.go"

ABIGEN_ARGS=(
    -bin "$BIN_PATH" \
    -abi "$ABI_PATH" \
    -out "$OUT_PATH" \
    -type "$CLASS_NAME" \
    -pkg "$2"
)

# Geth version from which native abigen was built, or empty string.
NATIVE_ABIGEN_VERSION=v"$((
    abigen --version 2> /dev/null | \
    egrep -o '([0-9]+\.[0-9]+\.[0-9]+)'
  ) || true
)"

# Generate golang wrapper
if [ "$NATIVE_ABIGEN_VERSION" == "$GETH_VERSION" ]; then
    abigen "${ABIGEN_ARGS[@]}"
else
    DOCKER_IMAGE="ethereum/client-go:alltools-$GETH_VERSION"
    echo -n "Native abigen unavailable, broken, or wrong version (need version "
    echo "$GETH_VERSION). Invoking abigen from $DOCKER_IMAGE docker image."
    docker run -v "${TMP_DIR}:${TMP_DIR}" --cidfile="$CONTAINER_NAME_PATH" \
           "$DOCKER_IMAGE" \
           abigen "${ABIGEN_ARGS[@]}"
    DOCKER_CONTAINER_NAME=$(cat "$CONTAINER_NAME_PATH")
    if [ "$(docker wait "$DOCKER_CONTAINER_NAME")" != "0" ] ; then
        echo "Failed to build $CLASS_NAME golang wrapper"
        exit 1
    fi
    docker cp "$DOCKER_CONTAINER_NAME:${OUT_PATH}" "${OUT_PATH}"
fi

TARGET_DIR="./generated/$2/"
mkdir -p "$TARGET_DIR"
cp "$OUT_PATH" "$TARGET_DIR"

"$CDIR/record_versions.sh" "$@"
