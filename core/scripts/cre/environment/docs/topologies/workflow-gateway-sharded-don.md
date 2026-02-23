# DON Topology

- Config: `configs/workflow-gateway-sharded-don.toml`
- Class: `sharded`
- Infra: `docker`

## Capability Matrix

This matrix is the source of truth for capability placement by DON.

| Capability | `bootstrap-gateway` | `shard0` | `shard1` |
|---|---|---|---|
| `consensus` | `-` | `local` | `local` |
| `cron` | `-` | `local` | `local` |
| `custom-compute` | `-` | `local` | `local` |
| `don-time` | `-` | `local` | `local` |
| `evm` | `-` | `local (1337,2337)` | `local (1337,2337)` |
| `http-action` | `-` | `local` | `-` |
| `http-trigger` | `-` | `local` | `-` |
| `ocr3` | `-` | `local` | `local` |
| `read-contract` | `-` | `local (1337,2337)` | `local (1337,2337)` |
| `vault` | `-` | `local` | `-` |
| `web-api-target` | `-` | `local` | `local` |
| `web-api-trigger` | `-` | `local` | `local` |
| `write-evm` | `-` | `local (1337,2337)` | `local (1337,2337)` |

## DONs

### `bootstrap-gateway`

- Types: `bootstrap`, `gateway`
- Nodes: `1`
- Roles: `bootstrap`, `gateway`
- EVM chains: `1337,2337`
- Exposes remote capabilities: `false`

### `shard0`

- Types: `shard`, `workflow`
- Nodes: `4`
- Roles: `plugin`
- EVM chains: `1337,2337`
- Exposes remote capabilities: `false`

### `shard1`

- Types: `shard`, `workflow`
- Nodes: `4`
- Roles: `plugin`
- EVM chains: `1337,2337`
- Exposes remote capabilities: `false`

