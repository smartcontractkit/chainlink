#!/bin/bash
# .github/scripts/cherry-pick-helper.sh
# Utility script for manual cherry-pick operations

set -euo pipefail

COMMIT_SHA="$1"
TARGET_BRANCH="$2"
SHORT_SHA="${COMMIT_SHA:0:7}"
BRANCH_NAME="cherry-pick/${SHORT_SHA}-to-${TARGET_BRANCH}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() { echo -e "${GREEN}ℹ️  $1${NC}"; }
log_warn() { echo -e "${YELLOW}⚠️  $1${NC}"; }
log_error() { echo -e "${RED}❌ $1${NC}"; }

# Validate inputs
if [[ -z "$COMMIT_SHA" || -z "$TARGET_BRANCH" ]]; then
    log_error "Usage: $0 <commit-sha> <target-branch>"
    exit 1
fi

log_info "Cherry-picking ${SHORT_SHA} to ${TARGET_BRANCH}"

# Update branches
log_info "Fetching latest changes..."
git fetch origin

log_info "Checking out develop branch..."
git checkout develop
git pull origin develop

log_info "Checking out target branch: ${TARGET_BRANCH}"
if ! git checkout "${TARGET_BRANCH}"; then
    log_error "Failed to checkout target branch: ${TARGET_BRANCH}"
    exit 1
fi
git pull origin "${TARGET_BRANCH}"

# Create cherry-pick branch
log_info "Creating cherry-pick branch: ${BRANCH_NAME}"
git checkout -b "${BRANCH_NAME}"

# Get commit message for PR
COMMIT_MSG=$(git log --format=%s -n 1 "${COMMIT_SHA}")

# Attempt cherry-pick
log_info "Attempting cherry-pick..."
if git cherry-pick "${COMMIT_SHA}"; then
    log_info "Cherry-pick successful! Pushing branch..."
    git push origin "${BRANCH_NAME}"
    
    # Create PR using GitHub CLI if available
    if command -v gh &> /dev/null; then
        log_info "Creating pull request..."
        gh pr create \
            --title "🍒 Cherry-pick: ${COMMIT_MSG} (${SHORT_SHA}) to ${TARGET_BRANCH}" \
            --body "Cherry-pick of ${COMMIT_SHA} to ${TARGET_BRANCH}

Original commit message:
${COMMIT_MSG}

This cherry-pick was created using the automated helper script." \
            --base "${TARGET_BRANCH}" \
            --head "${BRANCH_NAME}" \
            --label "cherry-pick" \
            --label "automated"
        
        log_info "✅ Cherry-pick PR created successfully!"
    else
        log_warn "GitHub CLI not available. Create PR manually:"
        log_info "Branch: ${BRANCH_NAME}"
        log_info "Target: ${TARGET_BRANCH}"
    fi
else
    log_error "Cherry-pick failed - conflicts detected"
    log_info "Conflict resolution required:"
    log_info "1. Resolve conflicts in the files"
    log_info "2. git add <resolved-files>"
    log_info "3. git cherry-pick --continue"
    log_info "4. git push origin ${BRANCH_NAME}"
    log_info "5. Create PR manually"
    exit 1
fi