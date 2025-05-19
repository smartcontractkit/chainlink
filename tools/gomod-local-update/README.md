# gomod-local-update

Updates any module that is `replace`'d with a local path to have its required module version in `go.mod` to match the latest git SHA from a remote branch.

Is meant to run within each directory where a `go.mod` file is present.

## Configuration

Command Line Flags:

```shell
Optional:
  -org-name       Organization name (default: smartcontractkit)
  -repo-name      Repository name (default: chainlink)
  -repo-remote    Git remote to use (default: origin)
  -branch-trunk   Branch to get SHA from (default: develop)
  -dry-run        Preview changes without applying them (default: false)
```

## Installation

The installed binary will be placed in your `$GOPATH/bin` directory. Make sure this directory is in your system's PATH to run the command from anywhere. From the root of this repository, run:

```shell
go install ./tools/gomod-local-update/cmd/gomod-local-update
```

## Usage Examples

Run from the root of a go module directory.

```shell
gomod-local-update
```

Was designed to be used with [gomods](https://github.com/jmank88/gomods) like:

```shell
gomods -w gomod-local-update
gomods tidy
```
