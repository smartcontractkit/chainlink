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
  merge_commit_sha=$(gh pr view $pr_number --json mergeCommit -q '.mergeCommit.oid')
  merge_date=$(gh pr view $pr_number --json mergedAt -q '.mergedAt')
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
for pr_info in "${pr_array[@]}"
do
  pr_number=$(echo $pr_info | cut -d',' -f1)
  merge_commit_sha=$(echo $pr_info | cut -d',' -f2)
  >&2 echo
  >&2 echo "Processing PR $pr_number with merge commit SHA $merge_commit_sha"

  # Create directory for this PR
  pr_dir="$target_dir/$pr_number"
  mkdir -p "$pr_dir"

  # Fetch all files from src directory
  echo "Command: git ls-tree -r \"$merge_commit_sha\" --name-only | grep \"^contracts/${source_dir}\""
  src_files=$(git ls-tree -r "$merge_commit_sha" --name-only | grep "^contracts/${source_dir}")
  echo $src_files

  for file in $src_files
  do
    mkdir -p "$pr_dir/$(dirname "$file")"
    git show "$merge_commit_sha:$file" > "$pr_dir/$file"
  done

  # Get list of modified Solidity files and save to PR_number.txt
  modified_files=$(git diff-tree --no-commit-id --name-only --diff-filter=AM -r "$merge_commit_sha" | grep '\.sol$' | grep -v -E '/test/|/tests/')
  for file in $modified_files
  do
    echo "$file" >> "$target_dir/$pr_number.txt"
  done
done

rm pr_merge_info.txt
rm sorted_prs.txt

>&2 echo
>&2 echo "Done processing PRs"
