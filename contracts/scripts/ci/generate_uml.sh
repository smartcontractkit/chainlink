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

if [ "$#" -lt 2 ]; then
    echo "Generates UML diagrams for all contracts in a directory after flattening them to avoid call stack overflows."
    echo "Usage: $0 <path to contracts> <path to target directory> [comma-separated list of files]"
    exit 1
fi

SOURCE_DIR="$1"
TARGET_DIR="$2"
FILES=${3// /}  # Remove any spaces from the list of files
FAILED_FILES=()

flatten_and_generate_uml() {
    local FILE=$1
    local TARGET_DIR=$2

    FLATTENED_FILE="$TARGET_DIR/flattened_$(basename "$FILE")"
    echo "Flattening $FILE to $FLATTENED_FILE"
    forge flatten "$FILE" -o "$FLATTENED_FILE" --root contracts
    if [ $? -ne 0 ]; then
        echo "Error: Failed to flatten $FILE"
        FAILED_FILES+=("$FILE")
        return
    fi

    OUTPUT_FILE=${FLATTENED_FILE//"flattened_"/""}
    OUTPUT_FILE_SVG="${OUTPUT_FILE%.sol}.svg"
    echo "Generating SVG UML for $FLATTENED_FILE to $OUTPUT_FILE_SVG"
    sol2uml "$FLATTENED_FILE" -o "$OUTPUT_FILE_SVG"
    if [ $? -ne 0 ]; then
        echo "Error: Failed to generate UML diagram in SVG format for $FILE"
        FAILED_FILES+=("$FILE")
        rm "$FLATTENED_FILE"
        return
    fi
    OUTPUT_FILE_DOT="${OUTPUT_FILE%.sol}.dot"
    echo "Generating DOT UML for $FLATTENED_FILE to $OUTPUT_FILE_DOT"
    sol2uml "$FLATTENED_FILE" -o "$OUTPUT_FILE_DOT" -f dot
    if [ $? -ne 0 ]; then
        echo "Error: Failed to generate UML diagram in DOT format for $FILE"
        FAILED_FILES+=("$FILE")
        rm "$FLATTENED_FILE"
        return
    fi

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
            return 1
        fi
    done
}

if [ -z "$FILES" ]; then
    flatten_contracts_in_directory "$SOURCE_DIR" "$TARGET_DIR"
else
    process_files "$SOURCE_DIR" "$TARGET_DIR" "$FILES"
fi

if [ "${#FAILED_FILES[@]}" -gt 0 ]; then
    echo "Error: Failed to generate UML diagrams for ${#FAILED_FILES[@]} files:"
    echo "Failed to generate UML diagrams for:" > "$TARGET_DIR/failed_files.txt"
    for FILE in "${FAILED_FILES[@]}"; do
        echo "  $FILE"
        echo "$FILE" >> "$TARGET_DIR/failed_files.txt"
    done
fi

echo "UML diagrams saved in $TARGET_DIR folder"
