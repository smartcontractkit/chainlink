#!/usr/bin/env bash

set -euo pipefail
# Disable history expansion so that '!' characters are taken literally.
# To handle paths like: /go/pkg/mod/github.com/!cosm!wasm
set +H

LOOPINSTALL_OUTPUT_DIR="${1:-/tmp/loopinstall-output}"
LIBS_DIR="${2:-/tmp/lib}"

if [ ! -d "${LOOPINSTALL_OUTPUT_DIR}" ]; then
  echo "Error: Directory ${LOOPINSTALL_OUTPUT_DIR} does not exist."
  exit 1
fi

if [ ! -d "${LIBS_DIR}" ]; then
  echo "Error: Directory ${LIBS_DIR} does not exist."
  exit 1
fi

echo "Processing loopinstall installation records from: ${LOOPINSTALL_OUTPUT_DIR}"
# Debug: print the JSON files' content.
cat "${LOOPINSTALL_OUTPUT_DIR}"/*.json

# Enable nullglob so that non-matching glob patterns expand to an empty array.
shopt -s nullglob

# Extract library glob patterns from the JSON files and process each one.
find "${LOOPINSTALL_OUTPUT_DIR}" -name "*.json" -print0 | xargs -0 jq -r '
    if .sources? then
      .sources[][] |
      if .libs? then
        .libs[]? // empty
      else
        empty
      end
    else
      empty
    end
  ' 2>/dev/null | while IFS= read -r pattern; do
    # Use mapfile to safely create array from glob pattern
    mapfile -t files < <(compgen -G "$pattern" 2>/dev/null || echo "")
    
    if [ ${#files[@]} -eq 0 ]; then
      echo "Warning: Source file not found: $pattern"
    else
      for file in "${files[@]}"; do
        if [ -f "$file" ]; then
          chmod 0755 "$file"
          echo "Copying $file to ${LIBS_DIR}"
          cp "$file" "$LIBS_DIR"
        fi
      done
    fi
done