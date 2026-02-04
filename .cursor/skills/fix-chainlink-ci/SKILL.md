---
name: fix-chainlink-ci
description: Triage and fix Chainlink GitHub Actions CI failures for pull requests (smartcontractkit/chainlink), focusing on merge conflicts/rebases, gofmt/gomod issues, failing Go tests, and flaky jobs. Use when the user mentions Chainlink CI, GitHub Actions, failing checks, merge conflicts, rebasing, flaky tests, or rerunning workflows.
---

# Fix Chainlink CI (GitHub Actions)

## Goal

Get a Chainlink PR to green by iterating: **inspect failing checks → classify cause → apply minimal fix → push → watch → rerun flakes**.

This is optimized for Chainlink’s typical pain points (high-churn files + multiple Go modules), e.g. PR #20708 touches `core/config/toml/types.go`, `docs/CONFIG.md`, and multiple `go.mod`/`go.sum` files and is a good “conflict magnet” example.

## Quick start

1. Work in `~/repos/chainlink` (or use `-R smartcontractkit/chainlink` with `gh` commands).
2. Collect facts:
   - `gh pr view <pr> --json url,title,headRefName,baseRefName,mergeable --jq '{url,title,head:.headRefName,base:.baseRefName,mergeable}'`
   - `gh pr checks <pr> --required`
   - Optional (for local fixes): `gh pr checkout <pr>`
3. For the first failure, prefer: `gh pr checks <pr> --watch --fail-fast`
4. Fix by category (below), push, then watch again until green.

## Workflow (loop until green)

### 1) Snapshot: what’s failing?

- Use `gh pr checks <pr> --required` to list required checks + links.
- If you need logs, identify the workflow run for the current head SHA/branch, then:
  - `gh run view <run-id> --log-failed`

### 2) Classify and fix

#### A) Merge conflict / PR not mergeable

Signals:
- `gh pr view <pr> --json mergeable --jq .mergeable` returns `CONFLICTING`, or UI says “has conflicts”.
- A workflow that checks mergeability fails, or checkout of PR merge ref fails.

Default fix (rebase onto base, usually `develop`):
- `git fetch origin <base>`
- `git rebase origin/<base>`
- Resolve conflicts (repeat until rebase completes):
  - `git status` → identify unmerged paths
  - edit files, remove conflict markers, keep intended behavior
  - `git add <files>`
  - `git rebase --continue`
- If the rebase requires rewriting history, push with `--force-with-lease` (never plain `--force`). If session rules require it, propose the exact command and wait for explicit confirmation.

Chainlink-specific tips:
- Conflicts commonly cluster in “high churn” files like `docs/CONFIG.md`, `core/config/toml/types.go`, and workflow syncer code under `core/services/workflows/...` (PR #20708 touches all of these).
- If conflicts involve `go.mod`/`go.sum`, finish the merge, then run the module tidy step (see C).

#### B) Lint / formatting

Signals:
- `gofmt`/`goimports` diffs, lint failures, or “files not formatted”.

Fix:
- Run the exact formatter/lint command shown in CI logs (preferred).
- For Go formatting, ensure files are gofmt’d before committing.

#### C) Go module / go.sum mismatch

Signals:
- CI complains `go.mod`/`go.sum` are dirty, `go list` errors, missing sums, or “tidy required”.

Fix:
- Chainlink has multiple Go modules; use the repo’s tidy helper when available:
  - `make gomodtidy`
- Otherwise run `go mod tidy` in each relevant module directory and commit resulting `go.mod`/`go.sum` updates.

#### D) Deterministic test failure (unit/integration)

Signals:
- Same test fails consistently; stack trace/assertion points to code.

Fix:
- Use `gh run view <run-id> --log-failed` to get the failing package/test name.
- Reproduce locally if feasible (copy the command from CI), otherwise make a minimal fix and rely on CI.
- Run a targeted test locally when possible (e.g., `go test ./path/to/pkg -run TestName -count=1`).
- When editing Go code, follow the repo Go guidelines (`.cursor/rules/go.mdc`): minimal diffs, explicit errors, gofmt/goimports, and minimal tests.

#### E) Flake / CI infrastructure noise

Signals (common):
- timeouts, transient network errors, runner failures, Docker pull flakes, “connection reset”, “context deadline exceeded”.
- A test that passes on rerun without code changes.

Fix:
- Don’t push “retry-only” commits.
- Wait for the run to finish, then rerun failed jobs (state-changing; get explicit confirmation if required by session rules):
  - `gh run rerun <run-id> --failed`
- If only a single job is flaky, rerun just that job using its database ID:
  - `gh run view <run-id> --json jobs --jq '.jobs[] | {name, databaseId, conclusion}'`
  - `gh run rerun <run-id> --job <databaseId>`
- After triggering a rerun, watch it:
  - `gh run watch <run-id> --exit-status`

### 3) Push + watch

- After applying a real fix, push the branch (respect repo/session rules on committing and force-pushing).
- Watch required checks:
  - `gh pr checks <pr> --watch`
- If you’re in a “conflict chain” (fix → new conflict → fix), repeat from snapshot.

## Guardrails

- Never bypass signing/hooks/auth (`--no-verify`, disabling GPG, etc.).
- If git/SSH auth or signing prompts appear stuck, remind the user to touch the YubiKey.
- Prefer minimal, mechanical fixes; avoid refactors while in “make CI green” mode.

## Additional resources

- Command cookbook: [reference.md](reference.md)
- Worked example (PR #20708): [examples.md](examples.md)
