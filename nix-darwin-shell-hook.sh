#!/bin/bash

# This script is used to set up cross compilation to linux arm64 in a CRIB environment.
# It's loaded during the shell hook execution in shell.nix
main() {
    echo "Running in CRIB environment, setting up cross compilation to linux arm64..."

    if ! command -v brew >/dev/null 2>&1; then
        echo "Homebrew is not installed. Please install Homebrew first: https://brew.sh/"
        exit 1
    fi

    PACKAGE="aarch64-unknown-linux-gnu"
    if ! brew list --formula | grep $PACKAGE > /dev/null; then
        echo "The Homebrew package $PACKAGE is not installed."
        echo "Please install it by running: brew tap messense/macos-cross-toolchains && brew install ${PACKAGE}"
        exit 1
    fi

    export GOOS=linux
    export CC=/opt/homebrew/Cellar/aarch64-unknown-linux-gnu/13.3.0/bin/aarch64-linux-gnu-gcc
    export CXX=/opt/homebrew/Cellar/aarch64-unknown-linux-gnu/13.3.0/bin/aarch64-linux-gnu-g++
}

main "$@"
