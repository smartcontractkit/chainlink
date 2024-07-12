#!/bin/bash

if [ "$#" -lt 2 ]; then
    echo "Finds the newest PR based on merge date from a comma-separated list of PR URLs."
    echo "Usage: $0 <repository base URL> <comma-separated list of PR ULRs>"
    exit 1
fi

repo_url="$1"
pr_urls="$2"

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

full_sorted_prs_list=$(cut -d',' -f1,2,3 sorted_prs.txt | tr '\n' ';' | sed 's/^;//')
IFS=';' read -r -a sorted_pr_array <<< "$full_sorted_prs_list"

>&2 echo "Sorted PRs from oldest to newest based on merge date:"
for pr_info in "${sorted_pr_array[@]}"
do
  pr_number=$(echo "$pr_info" | cut -d',' -f1)
  sha=$(echo "$pr_info" | cut -d',' -f2)
  date=$(echo "$pr_info" | cut -d',' -f3)
  >&2 echo "PR $pr_number (SHA: $sha) merged at $date"
done

# Get the SHA of the newest PR based on merge date
newest_pr_info=$(tail -n 1 sorted_prs.txt)
newest_pr_sha=$(echo "$newest_pr_info" | cut -d',' -f2)

>&2 echo
>&2 echo "SHA of the newest PR based on merge date: $newest_pr_sha"

rm pr_merge_info.txt sorted_prs.txt

echo "$newest_pr_sha"
