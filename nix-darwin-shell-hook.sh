#!/bin/bash

# This script is used to set up cross compilation to linux arm64 in a CRIB environment.
# It's loaded during the shell hook execution in shell.nix
main() {
    echo "Running in CRIB environment, setting up cross compilation to linux arm64..."
    PACKAGE="aarch64-unknown-linux-gnu"

    if ! command -v brew >/dev/null 2>&1; then
        echo -e "\e[31mHomebrew is not installed. Please install Homebrew first: https://brew.sh/\e[0m"
        exit 1
    fi

    if ! brew list --formula | grep $PACKAGE > /dev/null; then
        echo -e "\e[31mThe Homebrew package $PACKAGE is not installed.\e[0m"
        echo -e "\e[31mPlease install it by running: brew tap messense/macos-cross-toolchains && brew install ${PACKAGE}\e[0m"
        exit 1
    fi

    export GOOS=linux

    installed_version=$(brew list --versions $PACKAGE | awk '{print $2}')
    path_prefix=$(brew --prefix)
    bin_path=$path_prefix/Cellar/$PACKAGE/$installed_version/bin

    export CC=$bin_path/aarch64-linux-gnu-gcc
    export CXX=$bin_path/aarch64-linux-gnu-g++
}

main "$@"
