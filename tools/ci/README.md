# `tools/ci` - Unified CI CLI Tool

`ci` is a Go CLI designed to replace brittle inline bash scripts and complex `jq`/`sed`/`awk` pipelines across CI workflows with typed, tested, and locally reproducible Go subcommands.

## Commands

| Command | Description | Replaces |
|---|---|---|
| `ci version` | Print version information (text or `--json`) | Custom version checks |
| `ci image resolve` | Resolve Chainlink Docker image URI (ECR public or SDLC) | `.github/scripts/resolve-chainlink-image.sh` |
| `ci ccip resolve-baseline` | Resolve the CCIP release baseline image tag for mixed-version/rollout tests | `.github/scripts/resolve-ccip-release-baseline.sh` |
| `ci changeset check-tags` | Validate semver in changeset frontmatter and check release tags | `.github/scripts/check-changeset-tags.sh` |

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
