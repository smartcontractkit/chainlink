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

if [ "$#" -lt 5 ]; then
echo "Generates Markdown Slither reports and saves them to a target directory."
echo "Usage: $0 <https://github.com/ORG/REPO/blob/COMMIT/> <config-file> <root-directory-with–contracts> <comma-separated list of contracts> <where-to-save-reports> [slither extra params]"
exit 1
fi

REPO_URL=$1
CONFIG_FILE=$2
SOURCE_DIR=$3
FILES=${4// /}  # Remove any spaces from the list of files
TARGET_DIR=$5
SLITHER_EXTRA_PARAMS=$6

extract_product() {
    local path=$1

    echo "$path" | awk -F'src/[^/]*/' '{print $2}' | cut -d'/' -f1
}

run_slither() {
    local FILE=$1
    local TARGET_DIR=$2

    # needed, because the action we use returns all modified files, also deleted ones and we must skip those
    if [ ! -f "$FILE" ]; then
      echo "Warning: File not found: $FILE"
      echo "Skipping..."
      return
    fi

    source ./contracts/scripts/ci/select_solc_version.sh "$FILE"
    if [ $? -ne 0 ]; then
        echo "Error: Failed to select Solc version for $FILE"
        exit 1
    fi

    SLITHER_OUTPUT_FILE="$TARGET_DIR/$(basename "${FILE%.sol}")-slither-report.md"

    output=$(slither --config-file "$CONFIG_FILE" "$FILE" --checklist --markdown-root "$REPO_URL" --fail-none $SLITHER_EXTRA_PARAMS)
    if [ $? -ne 0 ]; then
        echo "Slither failed for $FILE"
        exit 1
    fi
    output=$(echo "$output" | sed '/\*\*THIS CHECKLIST IS NOT COMPLETE\*\*. Use `--show-ignored-findings` to show all the results./d'  | sed '/Summary/d')

    echo "# Summary for $FILE" > "$SLITHER_OUTPUT_FILE"
    echo "$output" >> "$SLITHER_OUTPUT_FILE"
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

process_files "$SOURCE_DIR" "$TARGET_DIR" "${FILES[@]}"

if [ $? -ne 0 ]; then
    echo "Error: Failed to generate Slither reports"
    exit 1
fi

echo "Slither reports saved in $TARGET_DIR folder"
