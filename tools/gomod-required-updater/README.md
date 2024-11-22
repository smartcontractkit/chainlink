# gomod-required-updater

Updates required module versions in go.mod files to match the latest git SHA from a remote branch.

## Configuration

Command Line Flags:

```shell
Optional:
  -org-name         Organization name to update modules for (default: smartcontractkit)
  -repo-name        Repository name to update modules for (default: chainlink)
  -repo-remote      Git remote to use (default: origin)
  -branch-trunk     Branch to get SHA from (default: develop)
  -dry-run         Preview changes without applying them (default: false)
  -
```

## Installation

The installed binary will be placed in your `$GOPATH/bin` directory. Make sure this directory is in your system's PATH to run the command from anywhere. From the root of this repository, run:

```shell
go install ./tools/gomod-required-updater/cmd/gomod-required-updater
```

## Usage Examples

```shell
make gomodrequiredupdater
```
