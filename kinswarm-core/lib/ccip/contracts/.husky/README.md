# Husky git hooks
The folder contains [husky](https://github.com/typicode/husky) git hooks that automate pre-commit and pre-push commands.

## Setup

Create an `.env` file in this folder to enable hooks:

```sh
# Can be left blank to compile everything
FOUNDRY_PROFILE=ccip
HUSKY_ENABLE_PUSH_HOOKS=true
HUSKY_ENABLE_COMMIT_HOOKS=true
UPSTREAM_BRANCH=origin/ccip-develop
```

```sh
# Automatically ran after pnpm install
pnpm prepare
```

### Script procedure

The setup is done via the `prepare.sh` script, which installs husky and enables the hooks.

The prepare step is skipped if we are in CI. This is checked via the `CI=true` flag, which is always set to true on GitHub actions.

## Hooks & Scripts

### Pre-commit
Runs [lint-staged](https://github.com/lint-staged/lint-staged), which automatically detects `.sol` file changes and applies `forge fmt` only on the changed files.

Configured in `package.json` in the `lint-staged` field.

### Pre-push
Runs forge build, test, solhint and optionally suggests to generate snapshots and wrappers.

Due to a [git workflow limitation](https://stackoverflow.com/questions/21334493/git-commit-in-pre-push-hook), generating wrappers & snapshots requires resubmitting the push (via `--no-verify` or by skiping the snapshot / wrappers).
