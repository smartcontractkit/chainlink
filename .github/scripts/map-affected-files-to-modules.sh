#!/bin/bash
set -e

# This script:
# 1. Finds all modules.
# 2. Maps changed files (passed as a param) to found modules.
# 3. Prints out the affected modules.
# 4. Output the result (as JSON) to a GitHub Actions environment variable.

# Get the list of changed files as parameter (from JSON array)
changed_files=$(echo "$1" | jq -r '.[]')
echo "Changed files: $changed_files"

# 1. Find all modules in the repository, 
# - Strip the leading './' from the path 
# (necessary for comparison, affected files do not have leading './')
modules=$(find . -name 'go.mod' -exec dirname {} \; | sed 's|^./||' | uniq)
echo "Found modules: $modules"

# Use a Bash associative array to track unique modules
declare -A unique_modules

for path_to_file in $changed_files; do
	echo "Resolving a module affected by a file: '$path_to_file'"
  for module in $modules; do
    echo "Validating against module: '$module'"

    # if no slash in the path, it is the root 
    # (i.e. `main.go`, `.gitignore` vs `core/main.go`)
    if [[ ! $path_to_file =~ \/ ]]; then
      echo "File '$path_to_file' mapped to the "root" module."
      unique_modules["."]="."
      break
    # if a module's name matches with a file path 
    # add it, to the affected modules array, skipping the root (`.`) 
    elif [[ $module != "." && $path_to_file =~ ^$module* ]]; then
      echo "File '$path_to_file' mapped the module '$module'"
      unique_modules["$module"]="$module"
      break
    fi
	done
done

# Convert keys (module names) of the associative array to an indexed array
affected_modules=("${!unique_modules[@]}")
echo "Affected modules: ${affected_modules[@]}"

# Convert bash array to a JSON array for GitHub Actions
json_array=$(printf '%s\n' "${affected_modules[@]}" | jq -R . | jq -s . | jq -c)
echo "module_names=$json_array" >> "$GITHUB_OUTPUT"
