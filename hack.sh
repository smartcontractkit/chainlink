#!/bin/bash

# Configuration
BASE_BRANCH=${1:-develop}
GOLANGCI_CONFIG=${2:-.golangci.yml}  # Optional config file
GOLANGCI_BIN=${3:-/opt/homebrew/bin/golangci-lint}     # Optional binary path/name

CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)

echo "Finding Go modules for packages changed between $CURRENT_BRANCH and $BASE_BRANCH..."

# Get the repository root directory
REPO_ROOT=$(git rev-parse --show-toplevel)

# Function to find the module for a package
find_module() {
    local pkg_path="$1"
    local dir="$REPO_ROOT/$pkg_path"
    
    # Walk up the directory tree until we find a go.mod file
    while [[ "$dir" != "$REPO_ROOT" && "$dir" != "/" ]]; do
        if [[ -f "$dir/go.mod" ]]; then
            # Found the module directory
            echo "${dir#$REPO_ROOT/}"
            return 0
        fi
        dir=$(dirname "$dir")
    done
    
    # Check the repo root itself
    if [[ -f "$REPO_ROOT/go.mod" ]]; then
        echo "."
        return 0
    fi
    
    # No go.mod found
    echo "No module found for $pkg_path"
    return 1
}


# Find all changed .go files and their packages
changed_packages=$(git diff --name-only "$BASE_BRANCH"..."$CURRENT_BRANCH" | grep '\.go$' | xargs -L1 dirname | sort -u)

if [ -z "$changed_packages" ]; then
    echo "No Go files changed."
    exit 0
fi

# Create a temporary file
TEMP_FILE=$(mktemp)
trap 'rm -f "$TEMP_FILE"' EXIT

# Find modules for each package
for pkg in $changed_packages; do
    dir="$REPO_ROOT/$pkg"
    module_path=""
    
    # Walk up the directory tree until we find a go.mod file
    while [[ "$dir" != "$REPO_ROOT" && "$dir" != "/" ]]; do
        if [[ -f "$dir/go.mod" ]]; then
            module_path="${dir#$REPO_ROOT/}"
            break
        fi
        dir=$(dirname "$dir")
    done
    
    # If no module found in parent directories, check repo root
    if [[ -z "$module_path" && -f "$REPO_ROOT/go.mod" ]]; then
        module_path="."
    fi
    
    # Store the mapping if a module was found
    if [[ -n "$module_path" ]]; then
        echo "$module_path $pkg" >> "$TEMP_FILE"
    fi
done

# Display modules and their packages
echo -e "\nModules with changed packages:"
sort "$TEMP_FILE" | awk '
    $1 != prev_module {
        if (NR > 1) print "";
        prev_module = $1;
        print "Module:", $1;
    }
    { print "  -", $2 }
'

# Run linter by module
echo -e "\nRunning linter on changed modules:"

# Get unique modules
modules=$(awk '{print $1}' "$TEMP_FILE" | sort -u)

for module in $modules; do
    module_path="$REPO_ROOT/$module"
    echo -e "\n=== Linting module $module ==="
    
    # Get packages for this module
    pkg_list=$(grep "^$module " "$TEMP_FILE" | cut -d' ' -f2)
    
    # Build arguments for golangci-lint
    lint_args=""
    for pkg in $pkg_list; do
        # Get package path relative to module
        rel_pkg="${pkg#$module/}"
        if [ "$rel_pkg" = "$pkg" ]; then
            rel_pkg="."
        fi
        lint_args="$lint_args ./$rel_pkg/..."
    done
    
    # Run linter, disable nolintlint https://github.com/golangci/golangci-lint/issues/3228
    cmd="$GOLANGCI_BIN run  --fix --new=false --new-from-rev= --new-from-merge-base=  --disable=nolintlint $lint_args"
    echo "running locally: $cmd"
    (cd "$module_path" && $cmd)
    #echo -e "\n=== Linting module $module with Docker ==="
    #echo "running in docker: docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint:v2.1.6 golangci-lint run --max-issues-per-linter 0 --max-same-issues 0 $lint_args"
    #(cd "$module_path" && docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint:v2.1.6 golangci-lint run --new=false --new-from-rev= --new-from-merge-base=  $lint_args)
done
