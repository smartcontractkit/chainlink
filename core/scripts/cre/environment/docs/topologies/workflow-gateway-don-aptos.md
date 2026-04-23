# DON Topology

- Config: `configs/workflow-gateway-don-aptos.toml`
- Class: `single-don`
- Infra: `docker`

## Capability Matrix

This matrix is the source of truth for capability placement by DON.

| Capability | `bootstrap-gateway` | `workflow` |
|---|---|---|
| `aptos` | `-` | `local (4)` |
| `consensus` | `-` | `local` |
| `cron` | `-` | `local` |

## DONs

### `bootstrap-gateway`

- Types: `bootstrap`, `gateway`
- Nodes: `1`
- Roles: `bootstrap`, `gateway`
- EVM chains: `1337`
- Exposes remote capabilities: `false`

### `workflow`

- Types: `workflow`
- Nodes: `4`
- Roles: `plugin`
- EVM chains: `1337`
- Exposes remote capabilities: `false`

