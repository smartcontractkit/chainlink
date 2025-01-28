#!/bin/bash

# Array of directories to process
DIRS=(
    "."
    "core/scripts"
    "integration-tests"
    "integration-tests/load"
    "deployment"
)

# Store the original directory
ORIGINAL_DIR=$(pwd)

# Function to run go mod tidy and check for errors
run_tidy() {
    local dir=$1
    echo "Running go mod tidy in $dir..."
    cd "$dir" || exit 1
    if ! go mod tidy; then
        echo "Error: go mod tidy failed in $dir"
        cd "$ORIGINAL_DIR"
        exit 1
    fi
    cd "$ORIGINAL_DIR"
}

# Process each directory
for dir in "${DIRS[@]}"; do
    run_tidy "$dir"
done

echo "All go mod tidy operations completed successfully!"