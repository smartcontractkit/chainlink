---
name: split-pr
description: Break large git branches into small, reviewable pull requests.
---

# Split PR & Stacked PRs

Use when branch is **LARGE** (>500 lines), touches multiple modules, or mixes refactoring with feature code.

## 0. Rules

- Commits must be signed with a touch-only GPG key. You must ask the user to make the commit, or prompt the user to tap their key.
- Generated code doesn't count towards layer size; generated code must be fully up to date.
- If `gh` or `gh stack` is unavailable, STOP! Prompt the user to install them before continuing.

## 1. Splitting Heuristics

Order layers bottom (trunk-adjacent) to top:
- **Layered Architecture**: Interfaces/types/schemas (bottom) -> Core service logic -> APIs/CLI/routes (top).
- **Module / Package**: Slice by Go module (`deployment/`, `core/`, `tools/`).
- **Refactor vs Feature**: Pure non-breaking refactor first -> new behavior on top.

```text
develop (trunk) <- layer-1-types <- layer-2-logic <- layer-3-cli (top)
```

## 2. Break Existing Branch into Stacked Layers

### Option A: Reset & Stage by File Set (Messy/Uncommitted History)
```bash
# 0. Help for stack CLI
gh stack --help

# 1. Backup current branch
git branch backup-feature-branch

# 2. Reset merge-base without losing changes
git reset $(git merge-base origin/develop HEAD)

# 3. Layer 1 (Bottom): Types/Interfaces
git checkout -b <ticket>/1-types origin/develop
git add path/to/types/ path/to/schemas/
git commit -m "feat(types): foundational interfaces"

# 4. Layer 2: Core Logic
git checkout -b <ticket>/2-logic <ticket>/1-types
git add path/to/services/
git commit -m "feat(services): implement core logic"

# 5. Layer 3 (Top): APIs / CLI
git checkout -b <ticket>/3-cli <ticket>/2-logic
git add path/to/cli/
git commit -m "feat(cli): endpoints and wiring"
```

### Option B: Cherry-Pick Commits (Clean Commit History)
```bash
git checkout -b <ticket>/1-types origin/develop && git cherry-pick <sha1> <sha2>
git checkout -b <ticket>/2-logic <ticket>/1-types && git cherry-pick <sha3>
git checkout -b <ticket>/3-cli <ticket>/2-logic && git cherry-pick <sha4>
```

## 3. Submit & Manage Stack with `gh stack`

Always use non-interactive flags to prevent agent hangs:

```bash
# Initialize stack from existing branches (bottom to top)
gh stack init <ticket>/1-types <ticket>/2-logic <ticket>/3-cli

# Submit/update PRs non-interactively (--auto avoids interactive editor)
gh stack submit --auto

# View stack state (JSON prevents TUI freeze)
gh stack view --json

# Navigate stack (never use interactive `gh stack switch`)
gh stack bottom
gh stack up
gh stack down
gh stack top

# Edit lower layer & propagate changes upstack
gh stack checkout <ticket>/1-types
# ...make edits...
git commit --amend --no-edit
gh stack rebase --upstack
gh stack submit --auto
```

## 4. Verification Checklist

Before PR submission:
- [ ] Generated code up to date: `make rm-mocked & make generate`
- [ ] Each layer compiles and passes lints independently: `tools/githooks lint`
- [ ] Unit tests pass per layer: `tools/githooks test`
- [ ] Layer size is **SMALL** or **MEDIUM**: `tools/githooks pr-size`
- [ ] Check that all commits are signed with a GPG key: `git log --show-signature`. If not, prompt the user to sign them.
