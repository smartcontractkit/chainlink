# ccip-release-changelog

Generates a CCIP-focused release changelog between two refs of this repo
(SHAs, tags like `v2.55.0`, or release branches like `release/2.56.0`), for
use as a risk-audit artifact in the release process.

For each ref pair it produces:

1. **go.mod diff** of the CCIP-relevant modules (chainlink-ccip and its
   `chains/*` submodules, chainlink-evm and its submodules,
   chainlink-aptos/codec, chainlink-sui/codec, chainlink-ton).
2. **Commit changelog for [chainlink-ccip](https://github.com/smartcontractkit/chainlink-ccip)**
   between the pinned commits.
3. **Commit changelogs for the chain-specific repos** (chainlink-aptos,
   chainlink-sui, chainlink-solana, chainlink-ton, chainlink-evm) between the
   pinned commits, plus the core repo restricted to `core/capabilities/ccip/`.
4. **Risk flags**: plugin gitRef changes in `plugins/plugins.public.yaml`,
   plugin-vs-go.mod drift (TON, EVM), rollbacks/divergence, modules or plugins
   added/removed, and keyword callouts (`breaking`, `revert`, `hotfix`,
   `security`, `fix!`, `config`).

## Usage

From the repo root (requires a full git history checkout):

```
go run ./tools/ccip/ccip-release-changelog/cmd/ccip-release-changelog \
    --old v2.55.0 --new release/2.56.0 \
    [--out report.md] [--slack-thread https://<ws>.slack.com/archives/<channel>/p<ts>]
```

Refs can be SHAs, tags, or branch names. If a branch only exists remotely
(e.g. `release/2.56.0` has never been checked out locally), the tool
automatically falls back to `origin/<ref>`.

Environment:

- `GITHUB_TOKEN` / `GH_TOKEN` — GitHub compare API auth. Falls back to
  `gh auth token`. All tracked repos are public, so this only affects rate
  limits.
- `SLACK_BOT_TOKEN` — required with `--slack-thread`. The bot must be a member
  of the target channel. The summary and flags are posted as a message in the
  thread; the full markdown report is uploaded as a file in the same thread.

In CI, use the `ccip-release-changelog` workflow (workflow_dispatch) which
takes the same inputs and uses the `SLACK_BOT_TOKEN_RELENG` secret.

## Configuration

All tracking behavior lives in one place: the `TrackedRepos` variable in
[`internal/changelog/config.go`](./internal/changelog/config.go). There are
**no CLI flags or workflow inputs for tracking** — edit that file and
re-run. Each entry is a `RepoConfig`:

| Field | Meaning |
|---|---|
| `Name` / `Owner` | GitHub repo (`Owner/Name`). Used for the compare API call and for commit/PR links. |
| `GoModules` | Module paths in the **root `go.mod`** that come from this repo, most important first. All of them appear in the *go.mod changes* section and participate in divergence notes. |
| `PluginKeys` | Keys in `plugins/plugins.public.yaml` that install from this repo (e.g. `ton`, `evm`). |
| `IncludePaths` | If non-empty, only commits touching at least one of these path prefixes appear in the commit changelog. |
| `ExcludePaths` | Commits touching *only* these path prefixes are dropped (applied before `IncludePaths`). |
| `Local` | Read the commit log from the local git checkout instead of the compare API. Use for the core repo itself. |

### How the "primary pin" is chosen

The commit changelog compares exactly one old/new SHA pair per repo:

- If `PluginKeys` is set, the **plugin gitRef** is primary — that's what gets
  built into the release image (aptos, sui, solana, ton, evm).
- Otherwise the **first entry of `GoModules`** is primary (chainlink-ccip).

Consequences:

- **Dual-source drift flag** — when a repo has both a plugin entry and its
  `moduleURI` also listed in `GoModules` (today: ton, evm), the two SHAs must
  agree at each ref or a `DRIFT` flag is raised.
- **Divergence notes** — when a repo's various pins point at *different
  commits of the same repo* (e.g. `chainlink-ccip` vs
  `chainlink-ccip/chains/evm`, or `plugin:sui` vs `chainlink-sui/codec`), a
  note is rendered in that repo's section. These are informational, not
  flags, and are normal for repos that ship mixed pins.
- Only the primary pin's SHA range gets a commit changelog. Submodule bumps
  still show up in the go.mod diff section.

### Editing examples

**Add a new repo** (shows up in all sections; primary = plugin gitRef if it
has one, else first GoModule):

```go
{
    Name:      "chainlink-tron",
    Owner:     "smartcontractkit",
    GoModules: []string{"github.com/smartcontractkit/chainlink-tron/relayer"},
},
```

**Change which core-repo paths are tracked** — edit the `Local` entry's
filters, e.g. to also include CCIP deployment code:

```go
IncludePaths: []string{"core/capabilities/ccip/", "deployment/ccip/"},
```

**Stop tracking a module** — remove it from `GoModules` (see the
`contracts/cre/gobindings` comment in the chainlink-evm entry for precedent).

### After editing

The golden tests render from `TrackedRepos`, so they must be regenerated:

```
UPDATE_GOLDEN=1 go test ./tools/ccip/ccip-release-changelog/...
go test ./tools/ccip/ccip-release-changelog/...
golangci-lint run ./tools/ccip/ccip-release-changelog/...
```

Then sanity-check a real run, e.g. `--old v2.55.0 --new release/2.56.0`.

### Related knobs (not in config.go)

- Keyword callout pattern (`breaking|revert|hotfix|security|config|fix!`):
  `keywordPattern` in
  [`internal/changelog/analyze.go`](./internal/changelog/analyze.go).
- Slack message/markdown layout: `report.go`; compare-API behavior:
  `github.go`.

## Development

```
go test ./tools/ccip/ccip-release-changelog/...          # unit + golden tests
UPDATE_GOLDEN=1 go test ./tools/ccip/ccip-release-changelog/...  # regen goldens
```
