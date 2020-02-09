#!/bin/bash

compile_command="yarn workspace chainlinkv0.6 compile"

# Only print compilation output on failure.
out="$($compile_command 2>&1)"

# shellcheck disable=SC2181
if [ "$?" != "0" ]; then
    echo
    echo "↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓"
    echo "Error while compiling solidity contracts. See below for output."
    echo "You can reproduce this error directly with the command"
    echo
    echo "   " "$compile_command"
    echo
    echo "This is probably a problem with a solidity contract, in the directory"
    echo "evm/v0.6/contracts."
    echo "↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑"
    echo
    echo "$out"
    exit 1
fi
