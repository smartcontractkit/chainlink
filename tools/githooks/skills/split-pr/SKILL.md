---
name: split-pr
description: Context, strategies, and step-by-step workflows for an agent or developer to break a large PR/branch into smaller, reviewable chunks and GitHub stacked pull requests.
---

# Splitting Large PRs & GitHub Stacked Pull Requests

Use this skill when a Git branch or pull request is classified as **LARGE** (>500 effective lines) or when changes span multiple layers, modules, or concerns.

## Why PRs Must Be Small

1. **Review Speed & Quality**: Reviewers review diffs $\le 200$ lines in minutes with high scrutiny. Diff sizes $> 500$ lines experience $3\times$ higher review latency and higher defect escape rates.
2. **CI Isolation**: In the `chainlink` monorepo, smaller package-targeted diffs trigger focused test shards (`tools/githooks test` and CI shards) rather than entire module test cascades.
3. **Merge Safety**: Smaller PRs minimize merge conflict windows and make rollbacks trivial.

---

## PR Splitting Strategies

Choose the strategy that matches the nature of the branch:

### Strategy 1: Stacked Layers (Recommended for Features)
Split by architectural dependency order from foundational to high-level:
1. **Layer 1 (Bottom)**: Interfaces, proto definitions, database schemas, core types, shared models.
2. **Layer 2 (Middle)**: Core service logic, handlers, internal adapters.
3. **Layer 3 (Top)**: API routes, CLI commands, documentation, UI/integration tests.

```text
   ┌── feat/part3-api-and-cli  → PR #3 (base: feat/part2-service-logic)   ← Top
  ┌── feat/part2-service-logic → PR #2 (base: feat/part1-interfaces-types)
 ┌── feat/part1-interfaces-types → PR #1 (base: develop)                  ← Bottom
develop (default branch)
```

### Strategy 2: By Go Module / Package (Recommended for Refactors & Sweeps)
If changes touch multiple Go modules (e.g. root `.`, `deployment`, `integration-tests`, `tools/`):
- PR 1: Module A changes (`./deployment/...`)
- PR 2: Module B changes (`./core/...`)
- PR 3: Shared tooling / config updates

### Strategy 3: Refactor / Cleanup First, Feature Second
Separate mechanical refactoring (renaming, reformatting, signature changes) from behavioral additions:
- PR 1: Pure refactor (existing behavior preserved, tests green).
- PR 2: New feature logic built on top of refactored code.

### Strategy 4: Generated Files & Dependencies
Put massive generated files (mocks from mockery, protobuf generated `.pb.go`, config docs) into their own stacked PR layer after the generator input files.

---

## How to Create GitHub Stacked PRs Fast

GitHub supports stacked pull requests natively. Each PR targets the branch below it.

### Step 1: Identify the Commits or File Sets
List commits on the current large branch:
```bash
git log --oneline origin/develop..HEAD
```

### Step 2: Create Stacked Branches

#### Option A: Using Git Cherry-Pick (Commit-Based)
```bash
# 1. Branch 1: Base layer from default branch
git checkout -b <ticket>/layer-1-types origin/develop
git cherry-pick <commit-sha-1> <commit-sha-2>
git push -u origin <ticket>/layer-1-types

# 2. Branch 2: Next layer stacked ON TOP of layer 1
git checkout -b <ticket>/layer-2-logic <ticket>/layer-1-types
git cherry-pick <commit-sha-3> <commit-sha-4>
git push -u origin <ticket>/layer-2-logic

# 3. Branch 3: Top layer stacked ON TOP of layer 2
git checkout -b <ticket>/layer-3-cli <ticket>/layer-2-logic
git cherry-pick <commit-sha-5>
git push -u origin <ticket>/layer-3-cli
```

#### Option B: Using Mixed Reset (File-Set Based)
If commits are messy or everything is uncommitted:
```bash
# Save branch reference
git branch backup-large-branch

# Reset to merge-base
git reset $(git merge-base origin/develop HEAD)

# Layer 1: Stage foundational files
git checkout -b <ticket>/layer-1-types origin/develop
git add core/types/ schema/
git commit -m "feat(types): foundational interfaces and types"
git push -u origin <ticket>/layer-1-types

# Layer 2: Stage core logic
git checkout -b <ticket>/layer-2-logic <ticket>/layer-1-types
git add core/services/
git commit -m "feat(services): implement core logic"
git push -u origin <ticket>/layer-2-logic
```

---

### Step 3: Open PRs Targeting Previous Layers

Use GitHub CLI (`gh`) to open the pull requests with the correct base:

```bash
# PR 1 targets develop
gh pr create --base develop --head <ticket>/layer-1-types \
  --title "[<ticket>] Part 1/2: Interfaces & Types" \
  --body "Stacked PR 1 of 2. Foundational types."

# PR 2 targets layer-1-types
gh pr create --base <ticket>/layer-1-types --head <ticket>/layer-2-logic \
  --title "[<ticket>] Part 2/2: Service Logic" \
  --body "Stacked PR 2 of 2. Builds on top of #<PR-1-NUMBER>."
```

*(Note: When PR #1 merges into `develop`, GitHub automatically re-targets PR #2's base branch to `develop` and cascades rebases.)*

---

## Rebasing a Stacked PR Chain

When code in a lower layer changes (e.g. PR review feedback on Layer 1):
```bash
# 1. Update Layer 1
git checkout <ticket>/layer-1-types
# make edits...
git commit --amend # or git commit
git push origin <ticket>/layer-1-types --force-with-lease

# 2. Rebase Layer 2 onto updated Layer 1
git checkout <ticket>/layer-2-logic
git rebase <ticket>/layer-1-types
git push origin <ticket>/layer-2-logic --force-with-lease
```

---

## Agent Verification Checklist

When splitting a PR:
- [ ] Each layer compiles independently (`go build ./...`).
- [ ] Unit tests pass in each layer (`tools/githooks test`).
- [ ] PR size guard classifies each individual layer as **SMALL** or **MEDIUM** (`tools/githooks pr-size`).
- [ ] The PR description references the stack order (e.g. "Part 1 of 3: base for #...").
