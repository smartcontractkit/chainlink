# release-changelog

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
go run ./tools/release/release-changelog/cmd/release-changelog \
    --old v2.55.0 --new release/2.56.0 \
    [--out report.md] [--slack-thread https://<ws>.slack.com/archives/<channel>/p<ts>]
```

Both `--old` and `--new` accept any ref: SHAs, tags (`v2.55.0`,
`v2.56.1-rc.3`), or branches (`release/2.57.2`). Resolution order: local
refs, then `origin/<ref>`, then an on-demand `git fetch origin <ref>` (into
`FETCH_HEAD` only — no local branches/tags are created). If a ref can't be
found anywhere, the error suggests similarly named remote refs (e.g. asking
for `v2.56.1` when only `release/2.56.1` and `v2.56.1-rc.N` exist).

You can also pass **CCIP image tags or image URIs** directly — the exact
version string used across the release process:

```
--old 2.56.1-ccip-rc.2
--old v2.56.1-ccip-rc.2   # hybrid form (v-prefix + -ccip-) also accepted
--old public.ecr.aws/chainlink/ccip:2.56.1-ccip-rc.2
```

`build-publish.yml` derives the image tag from the git tag that built it
(`v2.56.1-rc.2` → image `2.56.1-ccip-rc.2`), so the tool inverts that mapping
locally — **no image is ever pulled or inspected**. The report header notes
the mapping (`image tag → git tag`) for the audit trail.

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

The code is split into a product-agnostic engine and product packages:

- **[`internal/engine/`](./internal/engine/)** — everything generic: git ref
  resolution, go.mod/plugins.yaml parsing, compare-API changelogs, path
  filtering, risk flags, report rendering, Slack delivery. It knows nothing
  about CCIP; the `Product` it runs on is passed in by the caller.
- **[`internal/products/ccip/`](./internal/products/ccip/)** — the **CCIP
  product definition**: the tracked repo list (`Repos`) and the CCIP ref
  normalizer (`normalizeRef`, which maps CCIP image tags/URIs to git tags).
- **[`cmd/release-changelog/main.go`](./cmd/release-changelog/main.go)** —
  wires `ccip.Product` into `engine.Generate`. The product choice is a
  visible, reviewable line of code.

**Adding support for another product** (e.g. Core releases): create
`internal/products/<name>/` with its own `Product` value (copy
`ccip.go`), then wire it into a main package. Once a second product exists,
a `--product` flag / workflow input can select between them at runtime.

There are **no CLI flags or workflow inputs for tracking** — edit
`internal/products/ccip/ccip.go` and re-run. Each `Repos` entry is an
`engine.RepoConfig`:

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

The golden tests render from an inline engine fixture (not the real product
config), so they stay stable across config edits — but engine changes to
rendering or flag logic require regenerating them:

```
UPDATE_GOLDEN=1 go test ./tools/release/release-changelog/...
go test ./tools/release/release-changelog/...
golangci-lint run ./tools/release/release-changelog/...
```

Then sanity-check a real run, e.g. `--old v2.55.0 --new release/2.56.0`.

### Related knobs (not in the product config)

All in [`internal/engine/`](./internal/engine/):

- Keyword callout pattern (`breaking|revert|hotfix|security|config|fix!`):
  `keywordPattern` in `analyze.go`.
- Slack message/markdown layout: `report.go`; compare-API behavior:
  `github.go`; git ref resolution: `refs.go`.

## Development

```
go test ./tools/release/release-changelog/...          # unit + golden tests
UPDATE_GOLDEN=1 go test ./tools/release/release-changelog/...  # regen goldens
```
