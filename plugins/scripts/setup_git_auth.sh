#!/bin/sh
set -e  # Exit on error
set -u  # Exit on unset variable

# This script configures git to use a GitHub token for authentication
# with private repositories during Docker build.
# Usage in Dockerfile: 
#   COPY ./scripts/setup_git_auth.sh /tmp/
#   RUN --mount=type=secret,id=GIT_AUTH_TOKEN /tmp/setup_git_auth.sh

# Check if the token is provided
if [ -f "/run/secrets/GIT_AUTH_TOKEN" ]; then
  TOKEN=$(cat /run/secrets/GIT_AUTH_TOKEN)
  
  if [ -n "$TOKEN" ]; then
    # Set git config to use the token for github.com
    git config --global url."https://oauth2:${TOKEN}@github.com/".insteadOf "https://github.com/"
    echo "Git configured to use authentication token for GitHub repositories"
  else
    echo "No GitHub token content found, continuing without authentication"
  fi
else
  echo "No GitHub token file found, continuing without authentication"
fi
