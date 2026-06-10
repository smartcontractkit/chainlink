# DON Topology

- Config: `configs/workflow-gateway-don-cache-soak-test.toml`
- Class: `single-don`
- Infra: `docker`

## Capability Matrix

This matrix is the source of truth for capability placement by DON.

| Capability | `bootstrap-gateway` | `workflow` |
|---|---|---|
| `consensus` | `-` | `local` |
| `cron` | `-` | `local` |
| `don-time` | `-` | `local` |
| `evm` | `-` | `local (1337,2337)` |
| `http-action` | `-` | `local` |
| `http-trigger` | `-` | `local` |
| `vault` | `-` | `local` |

## DONs

### `bootstrap-gateway`

- Types: `bootstrap`, `gateway`
- Nodes: `1`
- Roles: `bootstrap`, `gateway`
- EVM chains: `1337,2337`
- Exposes remote capabilities: `false`

### `workflow`

- Types: `workflow`
- Nodes: `4`
- Roles: `plugin`
- EVM chains: `1337,2337`
- Exposes remote capabilities: `false`

