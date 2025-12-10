#!/bin/sh
set -e  # Exit on error
set -u  # Exit on unset variable

# This script configures git to use a GitHub token for authentication
# with private repositories during Docker build.
#
# Usage in Dockerfile (BuildKit):
#   # REQUIRED: Set GIT_CONFIG_GLOBAL to a temporary file path.
#   # This script will fail if GIT_CONFIG_GLOBAL is not set.
#   ENV GIT_CONFIG_GLOBAL=/tmp/gitconfig-github-token
#   RUN --mount=type=secret,id=GIT_AUTH_TOKEN \
#       --mount=type=cache,target=/go/pkg/mod \
#       --mount=type=cache,target=/root/.cache/go-build \
#       set -e && \
#       trap 'rm -f "$GIT_CONFIG_GLOBAL"' EXIT && \
#       ./plugins/scripts/setup_git_auth.sh && \
#       <your build commands>
#
# The RUN-level trap ensures the temporary git config is removed even if
# subsequent build commands fail.

if [ -z "${GIT_CONFIG_GLOBAL:-}" ]; then
  echo "ERROR: GIT_CONFIG_GLOBAL environment variable must be set to a temporary file path." >&2
  echo "Example: ENV GIT_CONFIG_GLOBAL=/tmp/gitconfig-github-token" >&2
  exit 1
fi

if [ -f "/run/secrets/GIT_AUTH_TOKEN" ]; then
  TOKEN=$(cat /run/secrets/GIT_AUTH_TOKEN)

  if [ -n "$TOKEN" ]; then
    git config --file "$GIT_CONFIG_GLOBAL" \
      url."https://oauth2:${TOKEN}@github.com/".insteadOf "https://github.com/"
    echo "Git configured to use authentication token for GitHub repositories"
  else
    echo "No GitHub token content found, continuing without authentication"
  fi
else
  echo "GIT_AUTH_TOKEN secret file not found, continuing without authentication"
fi
