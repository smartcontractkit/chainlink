# `tools/ci` - Unified CI CLI Tool

`ci` is a Go CLI designed to replace brittle inline bash scripts and complex `jq`/`sed`/`awk` pipelines across CI workflows with typed, tested, and locally reproducible Go subcommands.

## Commands

| Command | Description | Replaces |
|---|---|---|
| `ci version` | Print version information (text or `--json`) | Custom version checks |

## Usage

```sh
# Build locally
make ci-cli

# Run
tools/ci/.bin/ci --help
tools/ci/.bin/ci version
tools/ci/.bin/ci version --json
```

## Testing

```sh
go -C tools/ci test ./...
```
