#!/bin/bash

FILE="$1"

if [ "$#" -lt 1 ]; then
  echo "Detects the Solidity version of a file and selects the appropriate Solc version."
  echo "Usage: $0 <file>"
  exit 1
fi

if [ -z "$FILE" ]; then
  echo "Error: File not provided."
  exit 1
fi

echo "Detecting Solc version for $FILE"

if [[ -f "$FILE" ]]; then
SOLCVER="$(grep --no-filename '^pragma solidity' "$FILE" | cut -d' ' -f3)"
else
echo "Target is not a file"
exit 1
fi
SOLCVER="$(echo "$SOLCVER" | sed 's/[^0-9\.]//g')"

if [[ -z "$SOLCVER" ]]; then
# Fallback to latest version if the above fails.
SOLCVER="$(solc-select install | tail -1)"
fi

echo "Guessed $SOLCVER."

solc-select install "$SOLCVER"
solc-select use "$SOLCVER"
