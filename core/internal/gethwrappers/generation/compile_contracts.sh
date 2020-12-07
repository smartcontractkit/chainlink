#!/bin/bash

echo "compiling contracts"

CDIR=$(dirname "$0")
#COMPILE_COMMAND=$(<"$CDIR/compile_command.txt")
COMPILE_COMMAND=$(yarn workspace @chainlink/contracts belt compile solc)

# Only print compilation output on failure.
OUT="$(bash -c $COMPILE_COMMAND 2>&1)"
ERR="$?"

# shellcheck disable=SC2181
if [ "$ERR" != "0" ]; then
    echo
    echo "↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓"
    echo "Error while compiling solidity contracts. See below for output."
    echo "You can reproduce this error directly by running the command"
    echo
    echo "   " "$COMPILE_COMMAND"
    echo
    echo "in the directory $SOLIDITY_DIR"
    echo
    echo "This is probably a problem with a solidity contract, under the"
    echo "directory evm-contracts/src/."
    echo "↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑"
    echo
    echo "$OUT"
    exit 1
fi
