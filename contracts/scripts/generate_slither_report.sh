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

run_slither() {
    local FILE=$1
    local TARGET_DIR=$2

    SLITHER_OUTPUT_FILE="$TARGET_DIR/$(basename "${FILE%.sol}")-slither-report.md"
    slither --config-file "$CONFIG_FILE" "$FILE" --checklist --markdown-root "$REPO_URL"  > "$SLITHER_OUTPUT_FILE"
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
