#!/bin/sh
set -e  # Exit on error
set -u  # Exit on unset variable

# This script configures git to use a GitHub token for authentication
# with private repositories during Docker build.
#
# Usage in Dockerfile (BuildKit):
#   ENV GIT_CONFIG_GLOBAL=/tmp/gitconfig-github-token
#   RUN --mount=type=secret,id=GIT_AUTH_TOKEN \
#       --mount=type=cache,target=/go/pkg/mod \
#       --mount=type=cache,target=/root/.cache/go-build \
#       set -e && \
#       trap 'rm -f "${GIT_CONFIG_GLOBAL:-$HOME/.gitconfig}"' EXIT && \
#       ./plugins/scripts/setup_git_auth.sh && \
#       <your build commands>
#
# The RUN-level trap ensures the temporary git config is removed even if
# subsequent build commands fail.

if [ -f "/run/secrets/GIT_AUTH_TOKEN" ]; then
  TOKEN=$(cat /run/secrets/GIT_AUTH_TOKEN)

  if [ -n "$TOKEN" ]; then
    CONFIG_FILE="${GIT_CONFIG_GLOBAL:-$HOME/.gitconfig}"
    git config --file "$CONFIG_FILE" \
      url."https://oauth2:${TOKEN}@github.com/".insteadOf "https://github.com/"
    echo "Git configured to use authentication token for GitHub repositories"
  else
    echo "No GitHub token content found, continuing without authentication"
  fi
else
  echo "GIT_AUTH_TOKEN secret file not found, continuing without authentication"
fi
