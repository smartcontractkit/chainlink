#!/bin/bash
set -e

# Detect if in CI to skip hooks
# https://docs.github.com/en/actions/learn-github-actions/variables#default-environment-variables
if [[ $CI == "true" ]]; then
    exit 0
fi

# Skip hooks creation if unconfigured
if ! [ -f .husky/.env ]; then
    printf "\033[1;33mNo .env file found in contracts/.husky, skipping hooks setup.\e[0m\n"
    exit 0
fi

cd ../
chmod +x ./contracts/.husky/*.sh
pnpm husky ./contracts/.husky
echo "Husky hooks prepared."
