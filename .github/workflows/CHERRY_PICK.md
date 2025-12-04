# Automated Cherry-Pick Workflow

## Usage

### 1. Using Labels (Recommended)
Add labels to your PR before merging:
- `cherry-pick/release-2.29.1-ccip`
- `cherry-pick/release-cre`

The cherry-pick will automatically trigger when the PR is merged to `develop`.

### 2. Using Comments
Comment on any merged PR:
```
pick release-cre, release-2.29.1-ccip
```

### 3. Manual Cherry-Pick
```bash
# From repository root
./.github/scripts/cherry-pick-helper.sh <commit-sha> <target-branch>
```

## Branch Naming Convention

Cherry-pick branches follow the pattern:
```
cherry-pick/{short-sha}-to-{target-branch}
```

Example: `cherry-pick/abc1234-to-release-cre`

## Troubleshooting

### Cherry-pick fails with conflicts
1. The bot will comment on the original PR with manual instructions
2. Use the helper script: `.github/scripts/cherry-pick-helper.sh`
3. Resolve conflicts manually and create PR

### Target branch doesn't exist
- Ensure the target branch exists in the repository
- Check branch name spelling in labels/comments

## Monitoring

Check the Actions tab for:
- `Auto Cherry-Pick` workflow runs
- Cherry-pick success/failure notifications
- Created PR links