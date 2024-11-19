# gomod-required-updater

Updates required module versions in go.mod files to match the latest git SHA from a remote branch.

## Features

- Update modules to latest SHA from specified branch
- Auto-detect and update modules with local replace directives
- Update multiple modules at once
- Preview changes with dry run mode
- Configure via TOML or command line flags

## Configuration

Optional TOML config.

```toml
# List of modules to update
modules = [
    "github.com/smartcontractkit/chainlink/v2"
]
```

Command Line Flags:

```shell
Required (one of):
  -module             Module to update (can be specified multiple times)
  -config            Path to TOML config file
  -update-org-modules Auto-detect and update modules with local replaces

Optional:
  -repo-remote      Git remote to use (default: origin)
  -branch-trunk     Branch to get SHA from (default: develop)
  -root            Root path for searching go.mod files
  -dry-run         Preview changes without applying them
```

## Installation

The installed binary will be placed in your `$GOPATH/bin` directory. Make sure this directory is in your system's PATH to run the command from anywhere.

```shell
go install github.com/smartcontractkit/chainlink/tools/gomod-required-updater
```

## Usage Examples

Update Specific Modules:

```shell
# Update single module
gomod-required-updater -module github.com/org/repo`

# Update multiple modules
gomod-required-updater -module github.com/org/repo1 -module github.com/org/repo2

# Update using config file
gomod-required-updater -config modules.toml

# Using different remote/branch
gomod-required-updater -module github.com/org/repo -repo-remote upstream -branch-trunk main
```

Auto-detect and Update Local Modules:

```shell
# Update all local modules that have replace directives
gomod-required-updater -update-org-modules

# Preview changes first
gomod-required-updater -update-org-modules -dry-run
```

## Notes

- When using multiple module sources, precedence is:
  1. Command line `-module` flags
  2. Config file modules
  3. Auto-detected modules via `-update-org-modules` flag
- Use the `-dry-run` flag to safely preview changes
- Local replace directives are preserved during updates
