#!/bin/bash

function check_chainlink_dir() {
  local param_dir="chainlink"
  current_dir=$(pwd)

  current_base=$(basename "$current_dir")

  if [ "$current_base" != "$param_dir" ]; then
    echo "The script must be run from the root of $param_dir directory"
    exit 1
  fi
}

check_chainlink_dir

FILE="$1"

if [ "$#" -lt 1 ]; then
  echo "Detects the Solidity version of a file and selects the appropriate Solc version."
  echo "If the version is not installed, it will be installed and selected."
  echo "Will prefer to use the version from Foundry profile if it satisfies the version in the file."
  echo "Usage: $0 <file>"
  exit 1
fi

if [ -z "$FILE" ]; then
  echo "Error: File not provided."
  exit 1
fi

extract_product() {
    local path=$1

    echo "$path" | awk -F'src/[^/]*/' '{print $2}' | cut -d'/' -f1
}

extract_pragma() {
  local FILE=$1

  if [[ -f "$FILE" ]]; then
    SOLCVER="$(grep --no-filename '^pragma solidity' "$FILE" | cut -d' ' -f3)"
  else
    echo "$FILE is not a file or it could not be found. Exiting."
    return 1
  fi
  SOLCVER="$(echo "$SOLCVER" | sed 's/[^0-9\.^]//g')"
  >&2 echo "::debug::Detected Solidity version in pragma: $SOLCVER"
  echo "$SOLCVER"
}

echo "Detecting Solc version for $FILE"

PRODUCT=$(extract_product "$FILE")
SOLC_IN_PROFILE=$(FOUNDRY_PROFILE=$PRODUCT forge config --json --root contracts | jq ".solc")
SOLC_IN_PROFILE=$(echo "$SOLC_IN_PROFILE" | tr -d "'\"")
echo "::debug::Detected Solidity version in profile: $SOLC_IN_PROFILE"

SOLCVER=$(extract_pragma "$FILE")

exit_code=$?
if [ $exit_code -ne 0 ]; then
  echo "Error: Failed to extract the Solidity version from $FILE."
  return 1
fi

SOLCVER=$(echo "$SOLCVER" | tr -d "'\"")

if [[ "$SOLC_IN_PROFILE" != "null" && -n "$SOLCVER" ]]; then
  set +e
  COMPAT_SOLC_VERSION=$(npx semver "$SOLC_IN_PROFILE" -r "$SOLCVER")
  exit_code=$?
  set -e
  if [[ $exit_code -eq 0 && -n "$COMPAT_SOLC_VERSION" ]]; then
    echo "::debug::Version $SOLC_IN_PROFILE satisfies the constraint $SOLCVER"
    SOLC_TO_USE="$SOLC_IN_PROFILE"
  else
    echo "::debug::Version $SOLC_IN_PROFILE does not satisfy the constraint $SOLCVER"
    SOLC_TO_USE="$SOLCVER"
  fi
 elif [[ "$SOLC_IN_PROFILE" != "null" && -z "$SOLCVER" ]]; then
    >&2 echo "No version found in the Solidity file. Exiting"
    return 1
  elif [[ "$SOLC_IN_PROFILE" == "null" && -n "$SOLCVER" ]]; then
    echo "::debug::Using the version from the file: $SOLCVER"
    SOLC_TO_USE="$SOLCVER"
  else
    >&2 echo "No version found in the profile or the Solidity file."
    return 1
fi

echo "Will use $SOLC_TO_USE"
SOLC_TO_USE=$(echo "$SOLC_TO_USE" | tr -d "'\"")
SOLC_TO_USE="$(echo "$SOLC_TO_USE" | sed 's/[^0-9\.]//g')"

INSTALLED_VERSIONS=$(solc-select versions)

if echo "$INSTALLED_VERSIONS" | grep -q "$SOLC_TO_USE"; then
  echo "::debug::Version $SOLCVER is already installed."
  if echo "$INSTALLED_VERSIONS" | grep "$SOLC_TO_USE" | grep -q "current"; then
    echo "::debug::Version $SOLCVER is already selected."
  else
    echo "::debug::Selecting $SOLC_TO_USE"
    solc-select use "$SOLC_TO_USE"
  fi
else
  echo "::debug::Version $SOLC_TO_USE is not installed."
  solc-select install "$SOLC_TO_USE"
  solc-select use "$SOLC_TO_USE"
fi
