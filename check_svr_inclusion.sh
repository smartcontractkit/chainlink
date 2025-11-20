#!/bin/bash

# Script to check if Chainlink Core version 2.29.0 includes all commits from 2.28.0-svr-2
# Exits 0 if included, 1 if not included

set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo "🔍 Checking if 2.29.0 includes 2.28.0-svr-2..."

# Step 1: Fetch all tags from origin
echo "📡 Fetching tags from origin..."
git fetch --tags origin >/dev/null 2>&1 || true

# Step 2: Get commit SHAs for the tags
TAG_2290="v2.29.0"
TAG_SVR="v2.28.0-svr-2"

# Check if tags exist
if ! git rev-parse --verify "$TAG_2290" >/dev/null 2>&1; then
    echo -e "${RED}❌ Tag $TAG_2290 not found${NC}"
    exit 1
fi

if ! git rev-parse --verify "$TAG_SVR" >/dev/null 2>&1; then
    echo -e "${RED}❌ Tag $TAG_SVR not found${NC}"
    exit 1
fi

# Get the commit SHAs
COMMIT_2290=$(git rev-parse "$TAG_2290")
COMMIT_SVR=$(git rev-parse "$TAG_SVR")

echo "📍 $TAG_2290: $COMMIT_2290"
echo "📍 $TAG_SVR: $COMMIT_SVR"

# Step 3: Use git merge-base --is-ancestor to check if svr commit is an ancestor of 2.29.0
if git merge-base --is-ancestor "$COMMIT_SVR" "$COMMIT_2290"; then
    echo -e "${GREEN}✅ All svr fixes are included in 2.29.0${NC}"
    exit 0
else
    echo -e "${RED}❌ svr fixes not included. Commits missing:${NC}"
    echo ""
    echo -e "${BLUE}📋 Recent commits in 2.29.0 since 2.28.0-svr-2:${NC}"
    git log "$TAG_SVR".."$TAG_2290" --oneline --no-merges | head -n 10
    exit 1
fi
