# Reference: Chainlink CI fix cookbook

This file is intentionally “copy/paste friendly” for debugging Chainlink PR CI on GitHub Actions.

## Work in the right repo

```bash
cd ~/repos/chainlink
```

If you’re not in the repo, pass `-R smartcontractkit/chainlink` to `gh` commands.

## PR facts (mergeability + branches)

```bash
gh pr view <pr> --json url,title,headRefName,baseRefName,mergeable,headRefOid \
  --jq '{url,title,head:.headRefName,base:.baseRefName,mergeable,sha:.headRefOid}'
```

## Checkout / diff (optional but handy)

```bash
gh pr checkout <pr>
gh pr diff <pr>
gh pr view <pr> --json files --jq '.files[].path'
```

## List required checks

```bash
gh pr checks <pr> --required
```

Watch until completion (fast feedback):

```bash
gh pr checks <pr> --watch --fail-fast
```

## Find the workflow run for a PR head branch

```bash
HEAD_BRANCH="$(gh pr view <pr> --json headRefName --jq .headRefName)"
gh run list --branch "$HEAD_BRANCH" --event pull_request -L 10 \
  --json databaseId,workflowName,status,conclusion,createdAt,url \
  --jq '.[] | {id:.databaseId, workflow:.workflowName, status, conclusion, createdAt, url}'
```

## View failing logs

```bash
gh run view <run-id> --log-failed
```

If you need job-level drilldown:

```bash
gh run view <run-id> --json jobs --jq '.jobs[] | {name, databaseId, status, conclusion}'
gh run view <run-id> --log-failed --job <databaseId>
```

## Rerun (flakes)

Rerun only failed jobs:

```bash
gh run rerun <run-id> --failed
```

Rerun a specific job (note: must use `databaseId`, not the number from the URL):

```bash
gh run view <run-id> --json jobs --jq '.jobs[] | {name, databaseId, conclusion}'
gh run rerun <run-id> --job <databaseId>
```

Watch:

```bash
gh run watch <run-id> --exit-status
```

## Merge conflict / rebase loop

Fetch + rebase onto base (usually `develop`):

```bash
git fetch origin develop
git rebase origin/develop
```

During conflicts:

```bash
git status
git add <files>
git rebase --continue
```

Abort if you need to restart:

```bash
git rebase --abort
```

After a rebase, update the PR branch (use `--force-with-lease`, never `--force`):

```bash
git push --force-with-lease
```

## Go module hygiene (Chainlink-specific)

Chainlink has multiple Go modules. Prefer the repo helper:

```bash
make gomodtidy
```

## Local test quickies

Run all Go tests:

```bash
go test ./...
```

Run a targeted package/test:

```bash
go test ./path/to/pkg -run TestName -count=1
```
