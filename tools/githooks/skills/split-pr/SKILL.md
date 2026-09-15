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
- **Lock Prevention**: Set `git config core.filesRefLockTimeout 5000` to avoid transient index lock failures from IDE file watchers.
- **Worktree Isolation**: ALWAYS perform branch splitting inside an isolated worktree (`.worktrees/split-<ticket>`) so the main repository's `.git/index.lock` and IDE file watchers are never disturbed.

## 1. Splitting Heuristics

Order layers bottom (trunk-adjacent) to top:
- **Layered Architecture**: Interfaces/types/schemas (bottom) -> Core service logic -> APIs/CLI/routes (top).
- **Module / Package**: Slice by Go module (`deployment/`, `core/`, `tools/`).
- **Refactor vs Feature**: Pure non-breaking refactor first -> new behavior on top.

```text
develop (trunk) <- layer-1-types <- layer-2-logic <- layer-3-cli (top)
```

## 2. Break Existing Branch into Stacked Layers

### Option A: Isolated Worktree Flow (Recommended)
Prevents `.git/index.lock` collisions and IDE buffer thrashing.

```bash
# 1. Ensure lock timeout is configured
git config core.filesRefLockTimeout 5000

# 2. Backup current branch
git branch backup-<ticket>-feature

# 3. Create and enter dedicated staging worktree
git worktree add .worktrees/split-<ticket> origin/develop
cd .worktrees/split-<ticket>

# 4. Layer 1 (Bottom): Types/Interfaces / Foundation
git checkout -b <ticket>/1-types origin/develop
git checkout backup-<ticket>-feature -- path/to/types/ path/to/schemas/
# Verify layer compiles & passes tests
go test ./...
git commit -m "feat(types): foundational interfaces"

# 5. Layer 2: Core Logic
git checkout -b <ticket>/2-logic <ticket>/1-types
git checkout backup-<ticket>-feature -- path/to/services/
go test ./...
git commit -m "feat(services): implement core logic"

# 6. Layer 3 (Top): APIs / CLI / Config
git checkout -b <ticket>/3-cli <ticket>/2-logic
git checkout backup-<ticket>-feature -- path/to/cli/
go test ./...
git commit -m "feat(cli): endpoints and wiring"

# 7. Clean up worktree and return to main workspace
cd -
git worktree remove --force .worktrees/split-<ticket>
```

### Option B: Cherry-Pick Commits (Clean Commit History)
Can also be executed inside an isolated worktree:
```bash
git worktree add .worktrees/split-<ticket> origin/develop && cd .worktrees/split-<ticket>
git checkout -b <ticket>/1-types origin/develop && git cherry-pick <sha1> <sha2>
git checkout -b <ticket>/2-logic <ticket>/1-types && git cherry-pick <sha3>
git checkout -b <ticket>/3-cli <ticket>/2-logic && git cherry-pick <sha4>
cd - && git worktree remove --force .worktrees/split-<ticket>
```

## 3. PR Description Format

For each layer in the stack, generate a structured PR description using this format:

```md
## Intent

Short summary of what the intent of changes are. How things ultimately affect the user.

## Changes

Bullet list of changes. Short description of what, longer explanation of why.
```

Apply descriptions after submission:
```bash
gh pr edit <ticket>/1-types --body "$LAYER_1_BODY"
gh pr edit <ticket>/2-logic --body "$LAYER_2_BODY"
gh pr edit <ticket>/3-cli --body "$LAYER_3_BODY"
```

## 4. Submit & Manage Stack with `gh stack`

Always use non-interactive flags to prevent agent hangs (`--open` marks PRs ready for review):

```bash
# Initialize stack from existing branches (bottom to top)
gh stack init <ticket>/1-types <ticket>/2-logic <ticket>/3-cli

# Submit/update PRs non-interactively and publish as ready for review
gh stack submit --auto --open

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
gh stack submit --auto --open
```

## 5. Verification Checklist

Before PR submission:
- [ ] Generated code up to date: `make rm-mocked & make generate`
- [ ] Each layer compiles and passes lints independently: `tools/githooks lint`
- [ ] Unit tests pass per layer: `tools/githooks test`
- [ ] Layer size is **SMALL** or **MEDIUM**: `tools/githooks pr-size`
- [ ] Check that all commits are signed with a GPG key: `git log --show-signature`. If not, prompt the user to sign them.
- [ ] PR descriptions drafted with `## Intent` and `## Changes` for each layer.

