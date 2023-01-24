#!/usr/bin/env bash
set -euo pipefail

##
# Execute tests via: ./ontriggerlint_test.sh
#
# For the GitHub contexts that should be passed in as args, see:
# https://docs.github.com/en/actions/learn-github-actions/contexts
#
# Trigger golangci-lint job steps when event is one of:
#   1. on a schedule (GITHUB_EVENT_NAME)
#   2. on PR's where the target branch (GITHUB_BASE_REF) is not prefixed with 'release/*'
#      and where the relevant source code (SRC_CHANGED) has changed
#   3. on pushes to these branches (GITHUB_REF): staging, trying, rollup and where the
#      relevant source code (SRC_CHANGED) has changed
#
# env vars:
# GITHUB_EVENT_NAME GitHub's event name, ex: schedule|pull_request|push (GitHub context: github.event_name)
# GITHUB_BASE_REF   GitHub's base ref - target branch of pull request (GitHub context: github.base_ref)
# GITHUB_REF        GitHub's ref - branch or tag that triggered run (GitHub context: github.ref)
# SRC_CHANGED       Specific source paths that we want to trigger on were modified. One of: true|false
##

if [[ -z "${GITHUB_REF:-}" ]]; then GITHUB_REF=""; fi

# Strip out /refs/heads/ from GITHUB_REF leaving just the abbreviated name
ABBREV_GITHUB_REF="${GITHUB_REF#refs\/heads\/}"

ON_TRIGGER=false
if [[ "${GITHUB_EVENT_NAME:-}" = "schedule" ]]; then
    # Trigger on scheduled runs
    ON_TRIGGER=true
elif [[ "${SRC_CHANGED:-}" = "true" && "${GITHUB_EVENT_NAME:-}" = "pull_request" && "${GITHUB_BASE_REF:-}" != release/* ]]; then
    # Trigger if it's from a pull request targetting any branch EXCEPT the release branch
    ON_TRIGGER=true
elif [[ "${SRC_CHANGED:-}" = "true" && "${GITHUB_EVENT_NAME:-}" = "push" && "${ABBREV_GITHUB_REF}" =~ ^(staging|trying|rollup)$ ]]; then
    # Trigger if it's a push to specific branches
    ON_TRIGGER=true
fi

echo "on_trigger=${ON_TRIGGER}"
