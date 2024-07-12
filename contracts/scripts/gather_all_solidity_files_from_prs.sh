#!/bin/bash

if [ "$#" -lt 4 ]; then
    echo "Copies all Solidity files from PRs in a comma-separated list of PR URLs to a target directory."
    echo "Usage: $0 <repository base URL> <comma-separated list of PR URLs> <source_dir> <target directory>"
    exit 1
fi

repo_url="$1"
pr_urls="$2"
source_dir="$3"
target_dir="$4"

IFS=',' read -r -a url_array <<< "$pr_urls"
pr_numbers=""
repo_name=$(echo "$repo_url" | cut -d '/' -f 4-)

for url in "${url_array[@]}"
do
  if [[ $url == *"${repo_name}/pull/"* ]]; then
    pr_number=$(echo "$url" | grep -o -E '/pull/[0-9]+' | grep -o -E '[0-9]+$')
    pr_numbers+="$pr_number,"
  else
    >&2 echo "Error: PR URL $url does not belong to the $repo_url repository."
    exit 1
  fi
done
pr_numbers=${pr_numbers%,}

>&2 echo "Will process PR numbers: $pr_numbers"
>&2 echo

IFS=',' read -r -a pr_array <<< "$pr_numbers"
pr_merge_info=""
for pr_number in "${pr_array[@]}"
do
  merge_commit_sha=$(gh pr view "$pr_number" --json mergeCommit -q '.mergeCommit.oid')
  merge_date=$(gh pr view "$pr_number" --json mergedAt -q '.mergedAt')
  pr_merge_info+="$pr_number,$merge_commit_sha,$merge_date"$'\n'
done
echo "$pr_merge_info" > pr_merge_info.txt

sorted_prs=$(sort -t, -k3 pr_merge_info.txt)
echo "$sorted_prs" > sorted_prs.txt

full_sorted_prs_list=$(cut -d',' -f1,3 sorted_prs.txt | tr '\n' ';' | sed 's/^;//')
IFS=';' read -r -a sorted_pr_array <<< "$full_sorted_prs_list"

>&2 echo "Sorted PRs from oldest to newest based on merge date:"
for pr_info in "${sorted_pr_array[@]}"
  do
    pr_number=$(echo "$pr_info" | cut -d',' -f1)
    date=$(echo "$pr_info" | cut -d',' -f2)
    echo "PR $pr_number merged at $date"
done

sorted_prs_list=$(cut -d',' -f1,2 sorted_prs.txt | tr '\n' ';' | sed 's/^;//')

IFS=';' read -r -a pr_array <<< "$sorted_prs_list"
mkdir -p "$target_dir"

# Calculate the index of the last element in the array
last_index=$(( ${#pr_array[@]} - 1 ))

# Process only the newest PR for files
latest_pr_info="${pr_array[$last_index]}"
latest_pr_number=$(echo "$latest_pr_info" | cut -d',' -f1)
latest_merge_commit_sha=$(echo "$latest_pr_info" | cut -d',' -f2)

>&2 echo
>&2 echo "Processing latest PR $latest_pr_number with merge commit SHA $latest_merge_commit_sha"

# Fetch all files from source directory
echo "Command: git ls-tree -r \"$latest_merge_commit_sha\" --name-only | grep \"^contracts/${source_dir}\""
src_files=$(git ls-tree -r "$latest_merge_commit_sha" --name-only | grep "^contracts/${source_dir}")
echo $src_files

for file in $src_files
do
  mkdir -p "$target_dir/$(dirname "$file")"
  git show "$latest_merge_commit_sha:$file" > "$target_dir/$file"
done

# Prepare the modified_contracts.txt file
modified_contracts_file="$target_dir/modified_contracts.txt"

# Get list of modified Solidity files and save to modified_contracts.txt
modified_files=$(git diff-tree --no-commit-id --name-only --diff-filter=AM -r "$latest_merge_commit_sha" | grep '\.sol$' | grep -v -E '/test/|/tests/')
for file in $modified_files
do
  echo "$file" >> "$modified_contracts_file"
done

# Gather changes from older PRs
for pr_info in "${pr_array[@]}"
do
  pr_number=$(echo "$pr_info" | cut -d',' -f1)
  merge_commit_sha=$(echo "$pr_info" | cut -d',' -f2)
  if [[ $pr_number != "$latest_pr_number" ]]; then
    >&2 echo "Gathering changes for older PR $pr_number"
    modified_files=$(git diff-tree --no-commit-id --name-only --diff-filter=AM -r "$merge_commit_sha" | grep '\.sol$' | grep -v -E '/test/|/tests/')
    for file in $modified_files
    do
      echo "$file" >> "$modified_contracts_file"
    done
  fi
done

rm pr_merge_info.txt
rm sorted_prs.txt

>&2 echo
>&2 echo "Done processing PRs"
