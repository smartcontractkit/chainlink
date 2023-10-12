#!/bin/bash

count=0
while true; do
  ((count++))
  echo "Iteration: $count"
  rm -rf ~/.cache/hardhat-nodejs/compilers-v2
  pnpm i >/dev/null 2>&1
  # Run this if FIX is specified
  if [ "$1" == "FIX" ]; then
    pnpm link ~/src/cl/hardhat/packages/hardhat-core >/dev/null 2>&1
  fi

  pnpm clean >/dev/null 2>&1
  pnpm compile >logs-iter-$count.log 2>&1
  exit_status=$?
  if [ $exit_status -ne 0 ]; then
    echo "Last command returned a non-zero exit status: $exit_status"
    break
  fi
done
