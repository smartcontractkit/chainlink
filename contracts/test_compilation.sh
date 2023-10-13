#!/bin/bash
#
# Description: Test compilation of hardhat-core, demonstrating the compiler downloader concurrency issue when multiple compilers are used across a large project
#
# Continuously run pnpm compile until it fails, storing logs in logs-iter-<count>.log
#
# Requirements:
# - pnpm v8.9.0
# - Linux, as this script assumes your cache is located at ~/.cache/hardhat-nodejs/compilers-v2
#
# Usage: ./test_compilation.sh [FIX] [LOCAL_PATH_TO_HH_REPO]
#
# If FIX is specified, it will use the local version of hardhat-core via pnpm link
# LOCAL_PATH_TO_HH_REPO is the path to the hardhat repo, e.g. ~/my/repos which will become ~/my/repos/hardhat/packages/hardhat-core
#
# Setup:
#
# 1. Clone hardhat repo: git@github.com:HenryNguyen5/hardhat.git
# 2. Checkout the test branch: git checkout fix/compiler_downloader_concurrency
# 3. Run this script, demonstrating failure case: ./test_compilation.sh
# 4. Run this script, demonstrating success case: ./test_compilation.sh FIX ~/my/repos
# 5. The failure case should fail within <10 iterations, while the success case should run indefinitely

count=0
if [ "$1" == "FIX" ]; then
  echo "Testing patched version via link to $2/hardhat/packages/hardhat-core"
else
  echo "Testing current version"
  pnpm unlink hardhat >/dev/null 2>&1
fi

while true; do
  ((count++))
  echo "Iteration: $count"
  rm -rf ~/.cache/hardhat-nodejs/compilers-v2
  pnpm i >/dev/null 2>&1
  # Run this if FIX is specified
  if [ "$1" == "FIX" ]; then
    pnpm link "$2/hardhat/packages/hardhat-core" >/dev/null 2>&1
  fi

  pnpm clean >/dev/null 2>&1
  pnpm compile >logs-iter-$count.log 2>&1
  exit_status=$?
  if [ $exit_status -ne 0 ]; then
    echo "Last command returned a non-zero exit status: $exit_status"
    break
  fi
done
