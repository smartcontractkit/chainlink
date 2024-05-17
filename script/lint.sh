#!/bin/bash

# run_linter.sh
# This script runs golangci-lint either locally or in a Docker container based on version matching.
# Parameters
GOLANGCI_LINT_VERSION="$1"
COMMON_OPTS="$2"
DIRECTORY="$3"
EXTRA_OPTS="$4"
OUTPUT_FILE="$DIRECTORY/$(date +%Y-%m-%d_%H:%M:%S).txt"

# Prepare the lint directory
mkdir -p "$DIRECTORY"

DOCKER_CMD="docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint:v$GOLANGCI_LINT_VERSION golangci-lint run $COMMON_OPTS $EXTRA_OPTS"

if command -v golangci-lint >/dev/null 2>&1; then
    LOCAL_VERSION=$(golangci-lint version 2>&1 | grep -oE "version .[{0-9}.]+" |  sed "s|version .||")

    if [ "$LOCAL_VERSION" = "$GOLANGCI_LINT_VERSION" ]; then
        echo "Local golangci-lint version ($LOCAL_VERSION) matches desired version ($GOLANGCI_LINT_VERSION). Using local version."
        # shellcheck disable=SC2086
        golangci-lint run $COMMON_OPTS $EXTRA_OPTS  

    else
        echo "Local golangci-lint version ($LOCAL_VERSION) does not match desired version ($GOLANGCI_LINT_VERSION). Using Docker version."
        $DOCKER_CMD > "$OUTPUT_FILE"
    fi

else
    echo "Local golangci-lint not found. Using Docker version."
    $DOCKER_CMD > "$OUTPUT_FILE"
fi

echo "Linting complete. Results saved to $OUTPUT_FILE"
 