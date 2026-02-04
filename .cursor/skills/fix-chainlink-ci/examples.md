# Examples

## Example: PR #20708 (“Supporting Private Workflow Registry”)

PR #20708 is a good stress-test example because it touches:
- “high churn” config/docs (`core/config/toml/types.go`, `docs/CONFIG.md`)
- workflow syncer code (`core/services/workflows/syncer/v2/...`)
- multiple Go module files (`core/scripts/go.mod`, `system-tests/lib/go.mod`, `system-tests/lib/go.sum`, etc.)

### 1) Triage failing checks

```bash
gh pr view 20708 --json url,title,headRefName,baseRefName,mergeable \
  --jq '{url,title,head:.headRefName,base:.baseRefName,mergeable}'
gh pr checks 20708 --required
```

If you want quick feedback:

```bash
gh pr checks 20708 --watch --fail-fast
```

### 2) Merge conflict loop (illustrative)

If `mergeable` is `CONFLICTING` (or a mergeability check fails), rebase onto `develop`:

```bash
git fetch origin develop
git checkout feat/multi-source-workflows
git rebase origin/develop
```

Resolve conflicts in the most common hotspots for this PR:
- `core/config/toml/types.go` (new `AdditionalWorkflowSource` TOML config)
- `docs/CONFIG.md` / `core/config/docs/core.toml` (docs move a lot on `develop`)
- `core/services/workflows/syncer/v2/workflow_registry.go` (large file, frequent churn)

After conflict resolution:

```bash
git add <files>
git rebase --continue
```

Then run module tidy (PR #20708 updates module files):

```bash
make gomodtidy
```

Push (rebase requires history rewrite):

```bash
git push --force-with-lease
```

### 3) Flake loop (illustrative)

If a job flakes (timeouts/runner issues/transient network), rerun failed jobs instead of pushing a “retry” commit.

Find the newest PR run for the head branch:

```bash
HEAD_BRANCH="$(gh pr view 20708 --json headRefName --jq .headRefName)"
gh run list --branch "$HEAD_BRANCH" --event pull_request -L 3 \
  --json databaseId,workflowName,status,conclusion,url \
  --jq '.[] | {id:.databaseId, workflow:.workflowName, status, conclusion, url}'
```

View the failing logs and rerun:

```bash
gh run view <run-id> --log-failed
gh run rerun <run-id> --failed
gh run watch <run-id> --exit-status
```
