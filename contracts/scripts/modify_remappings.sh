#!/bin/bash

# Check if the correct number of arguments is provided
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <directory_prefix> <remappings_file>"
    exit 1
fi

# Get the directory prefix and remappings file location from the arguments
DIR_PREFIX=$1
REMAPPINGS_FILE=$2

# Check if the remappings file exists
if [ ! -f "$REMAPPINGS_FILE" ]; then
    echo "Error: Remappings file '$REMAPPINGS_FILE' not found."
    exit 1
fi

# Temporary file to store modified remappings
OUTPUT_FILE="remappings_modified.txt"

# Read the remappings file, prepend the directory prefix, and write to output file
while IFS= read -r line; do
    if [[ "$line" =~ ^[^=]+= ]]; then
        REMAPPED_PATH="${line#*=}"
        MODIFIED_LINE="${line%=*}=${DIR_PREFIX}/${REMAPPED_PATH}"
        echo "$MODIFIED_LINE" >> "$OUTPUT_FILE"
    else
        echo "$line" >> "$OUTPUT_FILE"
    fi
done < "$REMAPPINGS_FILE"

echo "Modified remappings have been saved to: $OUTPUT_FILE"
