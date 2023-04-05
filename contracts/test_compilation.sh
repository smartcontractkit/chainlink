#!/bin/bash

count=0
while true; do
  ((count++))
  echo "Iteration: $count"
  git clean -dffx >/dev/null 2>&1
  rm -rf ~/.cache/hardhat-nodejs/
  pnpm i >/dev/null 2>&1
  pnpm link ~/src/cl/hardhat/packages/hardhat-core >/dev/null 2>&1
  pnpm compile >logs.log 2>&1
  exit_status=$?
  if [ $exit_status -ne 0 ]; then
    echo "Last command returned a non-zero exit status: $exit_status"
    break
  fi
done
