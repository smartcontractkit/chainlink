# Plugout: Plugins Out of Sync Checker

This tool ensures that plugin versions in Chainlink's plugin YAML files stay in sync with the versions defined in `go.mod`.

## Overview

The Chainlink core node uses plugin whose plugin definitions are stored in YAML files (`./plugins/plugins.*.yaml`). Some of these plugins, like chainlink-data-streams, chainlink-feeds, and chainlink-solana, are also defined as dependencies in the `go.mod` file as their packages are used by the core node software.

When updating a plugin's version in `go.mod` (if it exists there), the corresponding `gitRef` field in the plugin YAML files should also be updated to maintain consistency. This tool automates the checking and synchronization process.

## Usage

### Check Mode (Default)

To verify if all plugin versions in YAML files match their corresponding versions in `go.mod`:

```bash
go run ./tools/plugout/
```

The tool will report any mismatches but won't make changes. It will exit with code 1 if there are mismatches.

### Update Mode

To automatically update the `gitRef` values in the plugin YAML files to match `go.mod` versions:

```bash
go run ./tools/plugout/ --update
```

### Options

- `--go-mod <path>`: Path to the go.mod file (default: ./go.mod)
- `--module <module-uri>`: Go module URI to check (can be specified multiple times)
- `--plugin-file <file-path>`: Plugin YAML file to check (can be specified multiple times)
- `--update`: Write the gitRef using the go.mod version for matching plugins

### Examples

Check specific modules:

```bash
go run ./tools/plugout/ --module github.com/smartcontractkit/chainlink-data-streams --module github.com/smartcontractkit/chainlink-feeds
```

Check in a custom plugin file:

```bash
go run ./tools/plugout/ --plugin-file ./plugins/plugins.custom.yaml
```

Update plugin references in multiple files:

```bash
go run ./tools/plugout/ --update --plugin-file ./plugins/plugins.public.yaml --plugin-file ./plugins/plugins.private.yaml
```

## CI Integration

This tool is integrated into the CI pipeline to ensure that plugin versions remain in sync. If the job fails, it indicates that there are version mismatches between `go.mod` and the plugin YAML files.

To resolve a failed build:

1. Run the tool with the `--update` flag to update the plugin YAML files
2. Commit and push the changes

## Notes

- If a version in `go.mod` is a pseudo-version (e.g., `v0.1.1-0.20250325191518-036bb568a69d`), the tool will extract the commit hash part (`036bb568a69d`) for comparison with the gitRef in the plugins manifest. For a valid match, either:
  - The gitRef in the plugins manifest must exactly match the full commit hash that the pseudo-version in go.mod refers to, or
  - The commit hash part of the pseudo-version must match the gitRef in the plugins manifest (which might be a longer version of the same commit hash)
