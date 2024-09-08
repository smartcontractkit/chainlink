#!/usr/bin/env bash

set -euo pipefail

function check_chainlink_dir() {
  local param_dir="chainlink"
  current_dir=$(pwd)

  current_base=$(basename "$current_dir")

  if [[ "$current_base" != "$param_dir" ]]; then
    >&2 echo "The script must be run from the root of $param_dir directory"
    exit 1
  fi
}

check_chainlink_dir

if [ "$#" -lt 5 ]; then
  >&2 echo "Generates Markdown Slither reports and saves them to a target directory."
  >&2 echo "Usage: $0 <https://github.com/ORG/REPO/blob/COMMIT/> <config-file> <root-directory-with–contracts> <comma-separated list of contracts> <where-to-save-reports> [slither extra params]"
  exit 1
fi

REPO_URL=$1
CONFIG_FILE=$2
SOURCE_DIR=$3
FILES=${4// /}  # Remove any spaces from the list of files
TARGET_DIR=$5
SLITHER_EXTRA_PARAMS=$6

run_slither() {
    local FILE=$1
    local TARGET_DIR=$2

    if [[ ! -f "$FILE" ]]; then
      >&2 echo "::error:File not found: $FILE"
      return 1
    fi

    set +e
    source ./contracts/scripts/ci/select_solc_version.sh "$FILE"
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
    local SOURCE_DIR=$1
    local TARGET_DIR=$2
    local FILES=(${3//,/ })  # Split the comma-separated list into an array

    mkdir -p "$TARGET_DIR"

    for FILE in "${FILES[@]}"; do
      FILE=${FILE//\"/}
      run_slither "$SOURCE_DIR/$FILE" "$TARGET_DIR"
    done
}

set +e
process_files "$SOURCE_DIR" "$TARGET_DIR" "${FILES[@]}"

if [[ $? -ne 0 ]]; then
    >&2 echo "::warning::Failed to generate some Slither reports"
    exit 0
fi

echo "Slither reports saved in $TARGET_DIR folder"
