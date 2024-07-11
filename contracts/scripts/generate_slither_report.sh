#!/bin/bash

if [ "$#" -lt 4 ]; then
echo "Generates Markdown Slither reports and saves them to a target directory."
echo "Usage: $0 <https://github.com/ORG/REPO/blob/COMMIT/> <config-file> <root-directory-with–contracts> <where-to-save-reports> [comma-separated list of contracts]"
exit 1
fi

REPO_URL=$1
CONFIG_FILE=$2
SOURCE_DIR=$3
TARGET_DIR=$4
FILES=${5// /}  # Remove any spaces from the list of files

run_slither() {
    local FILE=$1
    local TARGET_DIR=$2

    SLITHER_OUTPUT_FILE="$TARGET_DIR/$(basename "${FILE%.sol}")-slither-report.md"
    slither --config-file "$CONFIG_FILE" "$FILE" --checklist --markdown-root "$REPO_URL"  > "$SLITHER_OUTPUT_FILE"
}

flatten_contracts_in_directory() {
    local SOURCE_DIR=$1
    local TARGET_DIR=$2

    mkdir -p "$TARGET_DIR"

    for ITEM in $(find "$SOURCE_DIR" -type f -name '*.sol'); do
        run_slither "$ITEM" "$TARGET_DIR"
    done
}

process_files() {
    local SOURCE_DIR=$1
    local TARGET_DIR=$2
    local FILES=(${3//,/ })  # Split the comma-separated list into an array

    mkdir -p "$TARGET_DIR"

    for FILE in "${FILES[@]}"; do
        MATCH=$(find "$SOURCE_DIR" -type f -name "$(basename "$FILE")")
        if [ -n "$MATCH" ]; then
            run_slither "$MATCH" "$TARGET_DIR"
        else
            echo "File $FILE does not exist within the source directory."
        fi
    done
}

if [ -z "$FILES" ]; then
  flatten_contracts_in_directory "$SOURCE_DIR" "$TARGET_DIR"
else
  process_files "$SOURCE_DIR" "$TARGET_DIR" "$FILES"
fi

echo "UML diagrams and Slither reports saved in $TARGET_DIR folder"
