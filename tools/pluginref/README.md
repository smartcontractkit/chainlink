# Pluginref: Resolve Plugin gitRef

Looks up a plugin's `gitRef` in a plugins YAML file and prints a ref suitable for `git rev-parse`.

## Usage

```bash
go run ./tools/pluginref/ \
  --plugin-file ./plugins/plugins.public.yaml \
  --plugin solana
```

Example output for the current Solana entry: `b5a89c32fdc1`

Resolve to a full commit SHA in a checked-out repo:

```bash
commit_ref=$(go run ./tools/pluginref/ --plugin-file ./plugins/plugins.public.yaml --plugin solana)
git -C ./chainlink-solana rev-parse "${commit_ref}^{}"
```

## Options

- `--plugin-file <path>`: Path to a plugins YAML file (required)
- `--plugin <name>`: Plugin key under `plugins:` (required), e.g. `solana`, `starknet`

## Supported gitRef formats

| gitRef example | Output |
|---|---|
| `v1.3.1-0.20260605202330-b5a89c32fdc1` | `b5a89c32fdc1` (pseudo-version commit hash) |
| `v0.0.0-20260609211101-71d38bd6a0a9` | `71d38bd6a0a9` |
| `73cd3d46ad0ce2871160369cb6447aeb9b48513f` | full SHA unchanged |
| `v1.2.3` | `v1.2.3` (tag) |
| `sub/v1.2.3` | `sub/v1.2.3` (prefixed tag) |

## CI

Used by the Solana smoke test workflow in `.github/workflows/integration-tests.yml` to read the plugin version from `plugins.public.yaml` instead of `go.mod`.
