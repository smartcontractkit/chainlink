# DON Topology

- Config: `configs/examples/workflow-don-overrides.toml`
- Class: `single-don`
- Infra: `docker`

## Capability Matrix

This matrix is the source of truth for capability placement by DON.

| Capability | `workflow` |
|---|---|
| `evm` | `local (2337)` |
| `vault` | `local` |

## DONs

### `workflow`

- Types: `workflow`
- Nodes: `5`
- Roles: `plugin`
- EVM chains: `2337`
- Exposes remote capabilities: `false`

