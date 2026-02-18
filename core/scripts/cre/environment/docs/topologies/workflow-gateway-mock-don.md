# DON Topology

- Config: `configs/workflow-gateway-mock-don.toml`
- Class: `multi-don`
- Infra: `docker`

## Capability Matrix

This matrix is the source of truth for capability placement by DON.

| Capability | `bootstrap-gateway` | `capabilities` | `workflow` |
|---|---|---|---|
| `custom-compute` | `-` | `-` | `local` |
| `don-time` | `-` | `-` | `local` |
| `mock` | `-` | `local` | `-` |
| `ocr3` | `-` | `-` | `local` |
| `web-api-trigger` | `-` | `-` | `local` |
| `write-evm` | `-` | `-` | `local (1337)` |

## DONs

### `bootstrap-gateway`

- Types: `bootstrap`, `gateway`
- Nodes: `1`
- Roles: `gateway`
- Exposes remote capabilities: `false`

### `capabilities`

- Types: `capabilities`
- Nodes: `3`
- Roles: `plugin`
- Exposes remote capabilities: `false`

### `workflow`

- Types: `workflow`
- Nodes: `4`
- Roles: `plugin`
- EVM chains: `1337`
- Exposes remote capabilities: `false`

