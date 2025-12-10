#!/bin/bash
# Quick local test of the cherry-pick helper script
# This tests the core functionality without triggering GitHub Actions

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

echo "🧪 Cherry-Pick Helper Test Script"
echo "=================================="
echo ""

# Change to repo root
cd "$REPO_ROOT"

# Check if we're in a git repo
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo "❌ Not in a git repository"
    exit 1
fi

# Prompt for inputs
echo "This script will test the cherry-pick helper without creating actual PRs."
echo ""
read -p "Enter commit SHA to cherry-pick (or 'HEAD~1' for previous commit): " COMMIT_SHA
read -p "Enter target branch name (must exist, e.g., 'release-1.8'): " TARGET_BRANCH

# Resolve commit SHA if needed
COMMIT_SHA=$(git rev-parse "$COMMIT_SHA")
SHORT_SHA="${COMMIT_SHA:0:7}"

echo ""
echo "📋 Test Configuration:"
echo "  Commit: $SHORT_SHA ($COMMIT_SHA)"
echo "  Target: $TARGET_BRANCH"
echo ""

# Check if target branch exists
if ! git rev-parse --verify "origin/$TARGET_BRANCH" > /dev/null 2>&1; then
    echo "❌ Target branch 'origin/$TARGET_BRANCH' does not exist"
    echo ""
    echo "Available release branches:"
    git branch -r | grep -E "origin/release-" | sed 's/origin\//  /'
    exit 1
fi

# Confirm
read -p "Proceed with test? This will create a local branch and attempt cherry-pick (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Cancelled."
    exit 0
fi

echo ""
echo "🚀 Starting test..."
echo ""

# Save current branch
ORIGINAL_BRANCH=$(git branch --show-current)
echo "📍 Current branch: $ORIGINAL_BRANCH"

# Create a cleanup function
cleanup() {
    echo ""
    echo "🧹 Cleaning up..."
    git checkout "$ORIGINAL_BRANCH" 2>/dev/null || true
    git branch -D "test-cherry-pick-${SHORT_SHA}-to-${TARGET_BRANCH}" 2>/dev/null || true
    echo "✅ Cleanup complete"
}
trap cleanup EXIT

# Test the cherry-pick process (without pushing)
echo "📥 Fetching latest changes..."
git fetch origin

echo "🌿 Checking out target branch..."
git checkout "$TARGET_BRANCH"
git pull origin "$TARGET_BRANCH"

BRANCH_NAME="test-cherry-pick-${SHORT_SHA}-to-${TARGET_BRANCH}"
echo "🌱 Creating test branch: $BRANCH_NAME"
git checkout -b "$BRANCH_NAME"

echo "🍒 Attempting cherry-pick..."
if git cherry-pick -x "$COMMIT_SHA"; then
    echo ""
    echo "✅ Cherry-pick successful!"
    echo ""
    echo "📝 Commit details:"
    git log -1 --oneline
    echo ""
    echo "🔍 Files changed:"
    git show --stat HEAD
    echo ""
    echo "✅ TEST PASSED - Cherry-pick would work in production"
    echo ""
    echo "The branch was created locally but NOT pushed to remote."
    echo "Branch will be cleaned up automatically."
else
    echo ""
    echo "❌ Cherry-pick failed (conflicts detected)"
    echo ""
    echo "This would require manual resolution in production."
    echo "The workflow would post a comment with instructions."
    echo ""
    
    # Show conflict details
    echo "📋 Conflicting files:"
    git status --short
    echo ""
    
    # Abort the cherry-pick
    git cherry-pick --abort
    
    echo "⚠️ TEST RESULT: Cherry-pick would fail in production workflow"
fi

echo ""
echo "🏁 Test complete!"
