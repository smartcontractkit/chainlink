#!/usr/bin/env bash

set -euo pipefail

if [ "$#" -ne 2 ]; then
    >&2 echo "Usage: $0 <directory_prefix> <remappings_file>"
    exit 1
fi

DIR_PREFIX=$1
REMAPPINGS_FILE=$2

if [ ! -f "$REMAPPINGS_FILE" ]; then
    >&2 echo "::error:: Remappings file '$REMAPPINGS_FILE' not found."
    exit 1
fi

OUTPUT_FILE="remappings_modified.txt"

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
