#!/bin/bash

# This checks for if at least one tag exists from a list of tags provided in a changeset file
#
# TAG LIST:
# #nops : For any feature that is NOP facing and needs to be in the official Release Notes for the release.
# #added : For any new functionality added.
# #changed : For any change to the existing functionality. 
# #removed : For any functionality/config that is removed.
# #updated : For any functionality that is updated.
# #deprecation_notice : For any upcoming deprecation functionality.
# #breaking_change : For any functionality that requires manual action for the node to boot.
# #db_update : For any feature that introduces updates to database schema.
# #wip : For any change that is not ready yet and external communication about it should be held off till it is feature complete.
# #bugfix - For bug fixes.
# #internal - For changesets that need to be excluded from the final changelog.

if [ $# -eq 0 ]; then
  echo "Error: No changeset file path provided."
  exit 1
fi

CHANGESET_FILE_PATH=$1
tags_list=( "#nops" "#added" "#changed" "#removed" "#updated" "#deprecation_notice" "#breaking_change" "#db_update" "#wip" "#bugfix" "#internal" )
has_tags=false
found_tags=()

if [[ ! -f "$CHANGESET_FILE_PATH" ]]; then
  echo "Error: File '$CHANGESET_FILE_PATH' does not exist."
  exit 1
fi

changeset_content=$(sed -n '/^---$/,/^---$/{ /^---$/!p; }' $CHANGESET_FILE_PATH)
semvar_value=$(echo "$changeset_content" | awk -F": " '/"chainlink"/ {print $2}')

if [[ "$semvar_value" != "major" && "$semvar_value" != "minor" && "$semvar_value" != "patch" ]]; then
  echo "Invalid changeset semvar value for 'chainlink'. Must be 'major', 'minor', or 'patch'."
  exit 1
fi

while IFS= read -r line; do
  for tag in "${tags_list[@]}"; do
    if [[ "$line" == *"$tag"* ]]; then
      found_tags+=("$tag")
      echo "Found tag: $tag in $CHANGESET_FILE_PATH"
      has_tags=true
    fi
  done
done < "$CHANGESET_FILE_PATH"

if [[ "$has_tags" == false ]]; then
  echo "Error: No tags found in $CHANGESET_FILE_PATH"
fi

echo "has_tags=$has_tags" >> $GITHUB_OUTPUT
echo "found_tags=$(jq -jR 'split(" ") | join(",")' <<< "${found_tags[*]}")" >> $GITHUB_OUTPUT
