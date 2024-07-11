#!/bin/bash

if [ "$#" -lt 2 ]; then
    echo "Generates UML diagrams for all contracts in a directory after flattening them to avoid call stack overflows."
    echo "Usage: $0 <path to contracts> <path to target directory> [comma-separated list of files]"
    exit 1
fi

SOURCE_DIR=$1
TARGET_DIR=$2
FILES=${3// /}  # Remove any spaces from the list of files

flatten_and_generate_uml() {
    local FILE=$1
    local TARGET_DIR=$2

    FLATTENED_FILE="$TARGET_DIR/flattened_$(basename "$FILE")"
    forge flatten "$FILE" -o "$FLATTENED_FILE" &> /dev/null

    OUTPUT_FILE=${FLATTENED_FILE//"flattened_"/""}
    OUTPUT_FILE="${OUTPUT_FILE%.sol}.svg"
    sol2uml "$FLATTENED_FILE" -o "$OUTPUT_FILE"

    rm "$FLATTENED_FILE"
}

flatten_contracts_in_directory() {
    local SOURCE_DIR=$1
    local TARGET_DIR=$2

    mkdir -p "$TARGET_DIR"

    for ITEM in $(find "$SOURCE_DIR" -type f -name '*.sol'); do
        flatten_and_generate_uml "$ITEM" "$TARGET_DIR"
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
            flatten_and_generate_uml "$MATCH" "$TARGET_DIR"
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

echo "UML diagrams saved in $TARGET_DIR folder"
