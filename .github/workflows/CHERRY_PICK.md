# Automated Cherry-Pick Workflow

The auto-cherry-pick.yml workflow enables [IssueOps](https://issue-ops.github.io/docs/) for certain PR's that, when merged, should be be automatically cherry-picked into a target branch such as a release branch.

## Usage

Comment on any merged PR:
```
pick release-cre, release-2.29.1-ccip
```
This will create a signed PR against each of the specified branches

### Manual Cherry-Pick, Merge Conflicts

If automation fails, for instance because of merge conflict, you should manually create the PR using a helper script:
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

### Target branch doesn't exist
- Ensure the target branch exists in the repository
- Check branch name spelling in labels/comments

## Monitoring

Check the Actions tab for:
- `Auto Cherry-Pick` workflow runs
- Cherry-pick success/failure notifications
- Created PR links