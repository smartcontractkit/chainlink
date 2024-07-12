#!/bin/bash

if [ "$#" -lt 5 ]; then
echo "Generates Markdown Slither reports and saves them to a target directory."
echo "Usage: $0 <https://github.com/ORG/REPO/blob/COMMIT/> <config-file> <root-directory-with–contracts> <comma-separated list of contracts> <where-to-save-reports>"
exit 1
fi

REPO_URL=$1
CONFIG_FILE=$2
SOURCE_DIR=$3
FILES=${4// /}  # Remove any spaces from the list of files
TARGET_DIR=$5

extract_product() {
    local path=$1
    local product=$(echo "$path" | awk -F'src/[^/]*/' '{print $2}' | cut -d'/' -f1)
    echo "$product"
}

run_slither() {
    local FILE=$1
    local TARGET_DIR=$2

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

    SLITHER_OUTPUT_FILE="$TARGET_DIR/$(basename "${FILE%.sol}")-slither-report.md"
    PRODUCT=$(extract_product "$FILE")

    FOUNDRY_PROFILE=$PRODUCT slither --config-file "$CONFIG_FILE" "$FILE" --checklist --markdown-root "$REPO_URL"  > "$SLITHER_OUTPUT_FILE"
}

process_files() {
    local SOURCE_DIR=$1
    local TARGET_DIR=$2
    local FILES=(${3//,/ })  # Split the comma-separated list into an array

    mkdir -p "$TARGET_DIR"

    for FILE in "${FILES[@]}"; do
      run_slither "$SOURCE_DIR/$FILE" "$TARGET_DIR"
    done
}

process_files "$SOURCE_DIR" "$TARGET_DIR" "${FILES[@]}"

echo "UML diagrams and Slither reports saved in $TARGET_DIR folder"
