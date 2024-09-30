#!/bin/bash
set -e

# Determine the current branch
current_branch=$(git branch --show-current)
upstream_branch="origin/ccip-develop"

# Compare the directory against the upstream branch
changes_root=$(git diff --name-only $upstream_branch...$current_branch -- ".changeset")

if ! [ -n "$changes_root" ]; then
    printf "\033[1;33mRoot changeset changes not found, Consider running pnpm changeset in the root directory if there is significant off-chain impact.\e[0m\n"
fi

changes_contracts=$(git diff --name-only $upstream_branch...$current_branch -- "contracts/.changeset")

if ! [ -n "$changes_contracts" ]; then
    printf "\033[0;31mContracts changeset changes not found, Make sure to run & commit \033[1;33mpnpm changeset\033[0;31m in the contracts directory.\n"
    exit 1
fi
