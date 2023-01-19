#!/usr/bin/env bash
# set -euo pipefail

##
# Trigger golangci-lint job steps when event is one of:
# 1. on a schedule (cron)
# 2. on PR's where the head branch is not prefixed with "release/*"
# 3. on pushes to these branches: staging, trying, rollup
##

# TODO: make this use getopts and pass in the github expressions as args or convert to a zx/typescript script
# then call it and have it echo output "false" or "true" and >> send to GITHUB_OUTPUT from the yaml file


on_trigger=false
if [[ "${{ github.event_name }}" = "schedule" ]]; then
    on_trigger=true
elif [[ "${{ github.event_name }}" = "pull_request" && "${{ github.head_ref }}" != release/* ]]; then
    on_trigger=true
elif [[ "${{ github.event_name }}" = "push" && "${{ github.head_ref }}" =~ ^(staging|trying|rollup)$ ]]; then
    on_trigger=true
fi
echo "on_trigger=${on_trigger}" >> $GITHUB_OUTPUT
