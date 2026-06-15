# DON Topology

- Config: `configs/workflow-gateway-capabilities-multi-gateway-don.toml`
- Class: `multi-don`
- Infra: `docker`

## Capability Matrix

This matrix is the source of truth for capability placement by DON.

| Capability | `bootstrap-gateway-us` | `capabilities` | `gateway-eu` | `workflow` |
|---|---|---|---|---|
| `consensus` | `-` | `-` | `-` | `local` |
| `cron` | `-` | `-` | `-` | `local` |
| `don-time` | `-` | `-` | `-` | `local` |
| `evm` | `-` | `remote-exposed (2337)` | `-` | `local (1337)` |
| `http-action` | `-` | `-` | `-` | `local` |
| `http-trigger` | `-` | `-` | `-` | `local` |
| `vault` | `-` | `remote-exposed` | `-` | `-` |

## DONs

### `bootstrap-gateway-us`

- Types: `bootstrap`, `gateway`
- Nodes: `1`
- Roles: `bootstrap`, `gateway`
- EVM chains: `1337,2337`
- Gateway DON ID: `gateway_don_us`
- Exposes remote capabilities: `false`

### `capabilities`

- Types: `capabilities`
- Nodes: `4`
- Roles: `plugin`
- EVM chains: `1337,2337`
- Exposes remote capabilities: `true`

### `gateway-eu`

- Types: `gateway`
- Nodes: `1`
- Roles: `gateway`
- EVM chains: `1337,2337`
- Gateway DON ID: `gateway_don_eu`
- Exposes remote capabilities: `false`

### `workflow`

- Types: `workflow`
- Nodes: `4`
- Roles: `plugin`
- EVM chains: `1337,2337`
- Exposes remote capabilities: `false`

