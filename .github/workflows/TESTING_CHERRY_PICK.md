# Testing the Auto Cherry-Pick Workflow

This guide explains how to test the automated cherry-pick workflow without impacting production.

## Prerequisites

Before testing, ensure you have:
- [ ] A test repository or branch where you can safely test
- [ ] GitHub Actions enabled
- [ ] Permissions to create PRs and push branches
- [ ] A merged PR to `develop` branch to use as test case

## Testing Approaches

### 1. **Dry Run with Workflow Dispatch (Recommended)**

First, let's add a manual trigger to the workflow for testing:

```yaml
# Add to the 'on:' section in auto-cherry-pick.yml
on:
  workflow_dispatch:
    inputs:
      pr-number:
        description: 'PR number to cherry-pick'
        required: true
        type: number
      target-branches:
        description: 'Target branches (comma-separated, e.g., release-1.8,release-1.9)'
        required: true
        type: string
```

Then add a manual test job:

```yaml
# Add as a new job
  manual-test:
    if: github.event_name == 'workflow_dispatch'
    runs-on: ubuntu-latest
    steps:
      - name: Setup test inputs
        id: setup
        uses: actions/github-script@v8
        with:
          script: |
            const prNumber = '${{ github.event.inputs.pr-number }}';
            const { data: pr } = await github.rest.pulls.get({
              owner: context.repo.owner,
              repo: context.repo.repo,
              pull_number: prNumber
            });
            
            core.setOutput('commit-sha', pr.merge_commit_sha);
            core.setOutput('pr-title', pr.title);
            core.setOutput('target-branches', JSON.stringify(
              '${{ github.event.inputs.target-branches }}'.split(',').map(b => b.trim())
            ));
      
      # Then trigger the cherry-pick job with these outputs
```

### 2. **Test with Labels (Safe Production Test)**

The safest way to test in production:

1. **Create a test PR** to `develop` branch with a trivial change (e.g., update a comment in README)
2. **Add the label** `cherry-pick/test-release-branch` (create a test release branch first if needed)
3. **Merge the PR** to `develop`
4. **Watch the Actions tab** to see if the workflow triggers
5. **Verify** the cherry-pick PR is created

```bash
# Create a test release branch
git checkout develop
git pull origin develop
git checkout -b test-release-1.0
git push origin test-release-1.0

# Then create a test PR to develop, add label "cherry-pick/test-release-1.0", and merge
```

### 3. **Test with Comment Trigger**

After merging a test PR to `develop`:

1. **Comment on the merged PR**: `pick test-release-1.0`
2. **Watch the Actions tab** for the workflow run
3. **Verify** the cherry-pick PR is created

### 4. **Local Testing of the Helper Script**

Test the cherry-pick helper script locally before testing the full workflow:

```bash
# Navigate to the repo
cd /Users/kreherma/git/cll/deployment-workspace/chainlink

# Make sure you have a clean state
git checkout develop
git pull origin develop

# Test the script with a real commit and target branch
# Replace with actual values:
COMMIT_SHA="<some-commit-sha-from-develop>"
TARGET_BRANCH="test-release-1.0"

# Run the helper script
.github/scripts/cherry-pick-helper.sh "$COMMIT_SHA" "$TARGET_BRANCH"

# Check the results
git log --oneline -n 5
gh pr list --label cherry-pick
```

### 5. **Test the Custom Action Locally with Act**

Use [act](https://github.com/nektos/act) to test GitHub Actions locally:

```bash
# Install act if not already installed
brew install act

# Create a test event file
cat > test-event.json << 'EOF'
{
  "pull_request": {
    "merged": true,
    "number": 1234,
    "title": "Test PR",
    "base": { "ref": "develop" },
    "merge_commit_sha": "abc123def456",
    "labels": [
      { "name": "cherry-pick/test-release-1.0" }
    ]
  }
}
EOF

# Run the workflow locally
act pull_request --eventpath test-event.json -j validate-trigger

# Or test with a specific job
act pull_request --eventpath test-event.json -j cherry-pick
```

## Validation Checklist

After running any test, verify:

- [ ] **Workflow triggers correctly** from the expected event
- [ ] **Target branches are parsed** correctly from labels/comments
- [ ] **Cherry-pick branch is created** with correct naming: `cherry-pick/{sha}-to-{target}`
- [ ] **Cherry-pick includes `-x` flag** (check commit message in cherry-pick branch)
- [ ] **PR is created automatically** with proper title and description
- [ ] **PR has correct labels** (`cherry-pick`, `automated`)
- [ ] **Comments are posted** to original PR with status
- [ ] **Conflicts are handled gracefully** (test with a conflicting commit)
- [ ] **Multiple target branches work** in parallel

## Testing Edge Cases

### Test Conflict Handling

1. Create a commit that will conflict when cherry-picked
2. Trigger the cherry-pick workflow
3. Verify that:
   - The workflow fails gracefully
   - A helpful comment is posted with manual instructions
   - No orphaned branches are left

### Test Multiple Target Branches

1. Create a PR with multiple cherry-pick labels or comment with multiple targets
2. Verify that:
   - All target branches are processed in parallel
   - Each gets its own PR
   - Summary comment includes all targets

### Test Invalid Target Branch

1. Try to cherry-pick to a non-existent branch
2. Verify that:
   - The validation step catches this
   - A clear error message is shown
   - The workflow fails before attempting cherry-pick

## Quick Test Commands

Here's a bash script to set up a complete test scenario:

```bash
#!/bin/bash
# quick-test-cherry-pick.sh

set -e

REPO_DIR="/Users/kreherma/git/cll/deployment-workspace/chainlink"
cd "$REPO_DIR"

echo "🧪 Setting up cherry-pick workflow test..."

# 1. Create a test release branch
echo "📝 Creating test release branch..."
git checkout develop
git pull origin develop
git checkout -b test-release-auto-$(date +%s)
git push origin HEAD
TEST_BRANCH=$(git branch --show-current)

echo "✅ Created test branch: $TEST_BRANCH"

# 2. Create a test commit on develop
echo "📝 Creating test commit..."
git checkout develop
echo "# Test for cherry-pick automation $(date)" >> TEST_CHERRY_PICK.md
git add TEST_CHERRY_PICK.md
git commit -m "test: cherry-pick automation test"
git push origin develop
TEST_COMMIT=$(git rev-parse HEAD)

echo "✅ Created test commit: $TEST_COMMIT"

# 3. Test the helper script
echo "🍒 Testing cherry-pick helper script..."
.github/scripts/cherry-pick-helper.sh "$TEST_COMMIT" "$TEST_BRANCH"

echo "✅ Test complete!"
echo ""
echo "Next steps:"
echo "1. Check the created PR: gh pr list --label cherry-pick"
echo "2. Verify the cherry-pick branch: git log cherry-pick/${TEST_COMMIT:0:7}-to-$TEST_BRANCH"
echo "3. Clean up: git branch -D cherry-pick-* && git push origin --delete $TEST_BRANCH"
```

## Monitoring Test Results

Watch the workflow execution:

```bash
# Watch workflow runs
gh run list --workflow=auto-cherry-pick.yml

# View a specific run
gh run view <run-id>

# Watch logs in real-time
gh run watch <run-id>

# Check created PRs
gh pr list --label cherry-pick --state open
```

## Rollback Plan

If the workflow causes issues:

1. **Disable the workflow**: Go to Actions → Auto Cherry-Pick → Disable workflow
2. **Close any auto-generated PRs**: `gh pr close <pr-number>`
3. **Delete cherry-pick branches**: `git push origin --delete cherry-pick/*`
4. **Revert the workflow file**: `git revert <commit-sha>`

## Recommended Testing Sequence

1. ✅ **Start with local script testing** (lowest risk)
2. ✅ **Test with workflow_dispatch** (manual control)
3. ✅ **Test with a test branch** (isolated environment)
4. ✅ **Test comment trigger on test PR** (safe production test)
5. ✅ **Test label trigger on test PR** (safe production test)
6. ✅ **Test with multiple targets** (parallel processing)
7. ✅ **Test conflict scenarios** (error handling)
8. ✅ **Monitor for a week** before full rollout

---

**Need Help?**
- Check workflow logs in the Actions tab
- Review the script output for detailed error messages
- Test the helper script locally first to isolate issues
