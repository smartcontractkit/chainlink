#!/bin/bash

prev_output_file=$(mktemp)
next_output_file=$(mktemp)

count=0
while true; do
  ((count++))
  echo "Iteration: $count"
  git clean -dffx >/dev/null 2>&1
  rm -rf ~/.cache/hardhat-nodejs/
  pnpm i >/dev/null 2>&1
  pnpm link ~/src/cl/hardhat/packages/hardhat-core >/dev/null 2>&1
  pnpm compile >logs-iter-$count.log 2>&1
  exit_status=$?
  if [ $exit_status -ne 0 ]; then
    echo "Last command returned a non-zero exit status: $exit_status"
    break
  fi
done

function cleanup {
  rm -f "$prev_output_file" "$next_output_file"
}

trap cleanup EXIT

function collect_and_diff {
  find ~/.cache/hardhat-nodejs/ -type f -exec sh -c 'file="{}"; size=$(stat -c "%s" "$file"); perm=$(stat -c "%A" "$file"); name=$(basename "$file"); timestamp=$(stat -c "%y" "$file"); printf "%-12s\t%-12s\t%-40s\t%-32s\n" "$size" "$perm" "$name" "$timestamp"' \; >"$next_output_file"
  diff --color=always --unified=0 --suppress-common-lines "$prev_output_file" "$next_output_file" | tail -n +3
  mv "$next_output_file" "$prev_output_file"
}

# Run once initially to show the current state
collect_and_diff

# Watch for events in the directory
while inotifywait -q -r -e modify,create,delete ~/.cache/hardhat-nodejs/; do
  collect_and_diff
done
