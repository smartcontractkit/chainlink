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
  # the flag that indicates if the path matches any module
  is_path_in_modules=false
  for module in $modules; do
    echo "Validating against module: '$module'"
    # if a module's name matches with a file path 
    # add it, to the affected modules array, skipping the root (`.`) 
    if [[ $module != "." && $path_to_file =~ ^$module* ]]; then
      echo -e "File '$path_to_file' mapped to the module '$module'\n"
      unique_modules["$module"]="$module"
      is_path_in_modules=true
      break
    fi
  done
  # if no matched module default to root module
  if [[ $is_path_in_modules == false ]]; then
    echo "File '$path_to_file' did not match any module, defaulting to root '.'"
    unique_modules["."]="."
    is_path_in_modules=false
  fi
  is_path_in_modules=false
done 

# if the path is empty (for any reason), it will not get to the loop,
# so if the unique_modules array is empty, default to the root module
if [[ ${#unique_modules[@]} -eq 0 ]]; then
  echo "No files were changed, defaulting to the root module '.'"
  unique_modules["."]="."
fi

# Convert keys (module names) of the associative array to an indexed array
affected_modules=("${!unique_modules[@]}")
echo "Affected modules: ${affected_modules[@]}"

# Convert bash array to a JSON array for GitHub Actions
json_array=$(printf '%s\n' "${affected_modules[@]}" | jq -R . | jq -s . | jq -c)
echo "module_names=$json_array" >> "$GITHUB_OUTPUT"
