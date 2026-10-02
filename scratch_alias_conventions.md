# Alias & Import Standardization (PR #23883)

Alias convention applied repo-wide:

| Package | Import | Alias |
|---|---|---|
| `chainlink-common/pkg/config` | `github.com/smartcontractkit/chainlink-common/pkg/config` | `commonconfig` |
| `core/config` | `github.com/smartcontractkit/chainlink/v2/core/config` | `coreconfig` |
| `core/utils/config` | `github.com/smartcontractkit/chainlink/v2/core/utils/config` | `configutils` |

Rules applied:

- **`commonconfig`** replaces all ad-hoc aliases for the common pkg config (`config`, `commoncfg`, `pkgconfig`) — collision with `core/config` was the reason each file invented its own name before.
- **`coreconfig`** replaces bare `config` imports of `core/config` when the file also imports the common pkg — again avoids the name clash.
- **`configutils`** replaces bare `config` for `core/utils/config` (was shadowing/conflicting with the other two).
- Files importing only one of these may keep the bare `config` name (e.g. `core/config/env/env.go` imports common pkg as plain `config` since no clash exists there).
- No reordering beyond what goimports required — the sweep renamed aliases only, never moved symbols between packages.
