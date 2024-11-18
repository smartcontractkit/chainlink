# gomod-required-updater

## Installation

The installed binary will be placed in your `$GOPATH/bin` directory. Make sure this directory is in your system's PATH to run the command from anywhere.

```shell
go install github.com/smartcontractkit/chainlink/tools/gomod-required-updater
```

## Config

```toml
# Add any modules in here that you want to update the required version for. Can optionally be provided via the command line flag `-module`.
modules = [
    "github.com/smartcontractkit/chainlink/v2"
]
```

## Usage

```shell
# Using config file
gomod-required-updater -config modules.toml
# Using command line modules
gomod-required-updater -module github.com/org/repo1 -module github.com/org/repo2
# Mixed (command line takes precedence)
gomod-required-updater -config modules.toml -module github.com/org/repo1
```

Even though both `-config` and `-module` are provided, only `github.com/org/repo1` will be used because the command line `-module` flag takes precedence, and the config file is completely ignored when any `-module` flags are present.
