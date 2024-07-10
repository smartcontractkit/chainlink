#!/bin/bash

if [ -z "$1" ] || [ -z "$2" ]; then
    echo "Generates UML diagrams for all contracts in a directory after flattening them to avoid call stack overflows."
    echo "Usage: $0 <path to contracts> <path to target directory>"
    exit 1
fi

SOURCE_DIR=$1
TARGET_DIR=$2

flatten_contracts() {
    local SOURCE_DIR=$1
    local TARGET_DIR=$2

    mkdir -p "$TARGET_DIR"

    for ITEM in "$SOURCE_DIR"/*; do
        if [ -d "$ITEM" ]; then
            flatten_contracts "$ITEM" "$TARGET_DIR/$(basename "$ITEM")"
        elif [ -f "$ITEM" ] && [[ "$ITEM" == *.sol ]]; then
            FLATTENED_FILE="$TARGET_DIR/flattened_$(basename "$ITEM")"
            forge flatten "$ITEM" -o "$FLATTENED_FILE" &> /dev/null

            OUTPUT_FILE=${FLATTENED_FILE//"flattened_"/""}
            OUTPUT_FILE="${OUTPUT_FILE%.sol}.svg"
            sol2uml "$FLATTENED_FILE" -o "$OUTPUT_FILE"

            rm "$FLATTENED_FILE"
        fi
    done
}

flatten_contracts "$SOURCE_DIR" "$TARGET_DIR"

echo "UML diagrams saved in $TARGET_DIR"

#generate_uml "$MIRROR_DIR"


