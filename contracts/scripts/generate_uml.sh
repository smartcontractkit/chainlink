#!/bin/bash

if [ "$#" -lt 2 ]; then
    echo "Generates UML diagrams for all contracts in a directory after flattening them to avoid call stack overflows."
    echo "Usage: $0 <path to contracts> <path to target directory> [comma-separated list of files]"
    exit 1
fi

SOURCE_DIR="$1"
TARGET_DIR="$2"
FILES=${3// /}  # Remove any spaces from the list of files

flatten_and_generate_uml() {
    local FILE=$1
    local TARGET_DIR=$2

    FLATTENED_FILE="$TARGET_DIR/flattened_$(basename "$FILE")"
    echo "Flattening $FILE to $FLATTENED_FILE"
    forge flatten "$FILE" -o "$FLATTENED_FILE" &> /dev/null

    OUTPUT_FILE=${FLATTENED_FILE//"flattened_"/""}
    OUTPUT_FILE="${OUTPUT_FILE%.sol}.svg"
    echo "Generating UML for $FLATTENED_FILE to $OUTPUT_FILE"
    sol2uml "$FLATTENED_FILE" -o "$OUTPUT_FILE"

    rm "$FLATTENED_FILE"
}

flatten_contracts_in_directory() {
    local SOURCE_DIR=$1
    local TARGET_DIR=$2

    mkdir -p "$TARGET_DIR"

    find "$SOURCE_DIR" -type f -name '*.sol' | while read -r ITEM; do
        flatten_and_generate_uml "$ITEM" "$TARGET_DIR"
    done
}

process_files() {
    local SOURCE_DIR=$1
    local TARGET_DIR=$2
    local FILES=(${3//,/ })  # Split the comma-separated list into an array

    mkdir -p "$TARGET_DIR"

    for FILE in "${FILES[@]}"; do
        FILE=${FILE//\"/}
        MATCHES=($(find "$SOURCE_DIR" -type f -path "*/$FILE"))

        if [ ${#MATCHES[@]} -gt 1 ]; then
            echo "Error: Multiple matches found for $FILE:"
            for MATCH in "${MATCHES[@]}"; do
                echo "  $MATCH"
            done
            exit 1
        elif [ ${#MATCHES[@]} -eq 1 ]; then
            echo "File found: ${MATCHES[0]}"
            flatten_and_generate_uml "${MATCHES[0]}" "$TARGET_DIR"
        else
            echo "File $FILE does not exist within the source directory $SOURCE_DIR."
        fi
    done
}

if [ -z "$FILES" ]; then
    flatten_contracts_in_directory "$SOURCE_DIR" "$TARGET_DIR"
else
    process_files "$SOURCE_DIR" "$TARGET_DIR" "$FILES"
fi

echo "UML diagrams saved in $TARGET_DIR folder"
