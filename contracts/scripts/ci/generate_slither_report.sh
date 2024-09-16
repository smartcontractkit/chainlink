#!/usr/bin/env bash

set -euo pipefail

extract_product() {
    local path=$1

    echo "$path" | awk -F'src/[^/]*/' '{print $2}' | cut -d'/' -f1
}

extract_pragma() {
  local FILE=$1

  if [[ -f "$FILE" ]]; then
    SOLCVER="$(grep --no-filename '^pragma solidity' "$FILE" | cut -d' ' -f3)"
  else
    >&2 echo ":error::$FILE is not a file or it could not be found. Exiting."
    return 1
  fi
  SOLCVER="$(echo "$SOLCVER" | sed 's/[^0-9\.^]//g')"
  >&2 echo "::debug::Detected Solidity version in pragma: $SOLCVER"
  echo "$SOLCVER"
}

detect_solc_version() {
   local FOUNDRY_DIR=$1
   local FILE=$2

   echo "Detecting Solc version for $FILE"

   # Set FOUNDRY_PROFILE to the product name only if it is set; otherwise either already set value will be used or it will be empty
   PRODUCT=$(extract_product "$FILE")
   if [ -n "$PRODUCT" ]; then
     FOUNDRY_PROFILE="$PRODUCT"
   fi
   SOLC_IN_PROFILE=$(forge config --json --root "$FOUNDRY_DIR" | jq ".solc")
   SOLC_IN_PROFILE=$(echo "$SOLC_IN_PROFILE" | tr -d "'\"")
   echo "::debug::Detected Solidity version in profile: $SOLC_IN_PROFILE"

   set +e
   SOLCVER=$(extract_pragma "$FILE")

   if [[ $? -ne 0 ]]; then
     >&2 echo "::error:: Failed to extract the Solidity version from $FILE."
     return 1
   fi

   set -e

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
       >&2 echo "::error::No version found in the Solidity file. Exiting"
       return 1
     elif [[ "$SOLC_IN_PROFILE" == "null" && -n "$SOLCVER" ]]; then
       echo "::debug::Using the version from the file: $SOLCVER"
       SOLC_TO_USE="$SOLCVER"
     else
       >&2 echo "::error::No version found in the profile or the Solidity file."
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
}

if [ "$#" -lt 5 ]; then
  >&2 echo "Generates Markdown Slither reports and saves them to a target directory."
  >&2 echo "Usage: $0 <https://github.com/ORG/REPO/blob/COMMIT/> <config file> <directory with foundry.toml> <comma-separated list of contracts> <where-to-save-reports> [slither extra params]"
  exit 1
fi

REPO_URL=$1
CONFIG_FILE=$2
FOUNDRY_DIR=$3
FILES=${4// /}  # Remove any spaces from the list of files
TARGET_DIR=$5
SLITHER_EXTRA_PARAMS=${6-''}

run_slither() {
    local FILE=$1
    local TARGET_DIR=$2

    if [[ ! -f "$FILE" ]]; then
      >&2 echo "::error:File not found: $FILE"
      return 1
    fi

    set +e
    detect_solc_version "$FOUNDRY_DIR" "$FILE"
    if [[ $? -ne 0 ]]; then
        >&2 echo "::error::Failed to select Solc version for $FILE"
        return 1
    fi

    SLITHER_OUTPUT_FILE="$TARGET_DIR/$(basename "${FILE%.sol}")-slither-report.md"
    if ! output=$(slither --config-file "$CONFIG_FILE" "$FILE" --checklist --markdown-root "$REPO_URL" --fail-none $SLITHER_EXTRA_PARAMS); then
        >&2 echo "::warning::Slither failed for $FILE"
        return 0
    fi
    set -e
    output=$(echo "$output" | sed '/\*\*THIS CHECKLIST IS NOT COMPLETE\*\*. Use `--show-ignored-findings` to show all the results./d'  | sed '/Summary/d')

    echo "# Summary for $FILE" > "$SLITHER_OUTPUT_FILE"
    echo "$output" >> "$SLITHER_OUTPUT_FILE"

    if [[ -z "$output" ]]; then
        echo "No issues found." >> "$SLITHER_OUTPUT_FILE"
    fi
}

process_files() {
    local TARGET_DIR=$1
    local FILES=(${2//,/ })  # Split the comma-separated list into an array

    mkdir -p "$TARGET_DIR"

    for FILE in "${FILES[@]}"; do
      FILE=${FILE//\"/}
      run_slither "$FILE" "$TARGET_DIR"
    done
}

set +e
process_files "$TARGET_DIR" "${FILES[@]}"

if [[ $? -ne 0 ]]; then
    >&2 echo "::warning::Failed to generate some Slither reports"
    exit 0
fi

echo "Slither reports saved in $TARGET_DIR folder"
