# DON Topology

- Config: `configs/workflow-gateway-sharded-5-dons.toml`
- Class: `sharded`
- Infra: `docker`

## Capability Matrix

This matrix is the source of truth for capability placement by DON.

| Capability | `bootstrap-gateway` | `shard0` | `shard1` | `shard2` | `shard3` | `shard4` | `shard5` |
|---|---|---|---|---|---|---|---|
| `consensus` | `-` | `local` | `local` | `local` | `local` | `local` | `local` |
| `cron` | `-` | `local` | `local` | `local` | `local` | `local` | `local` |
| `custom-compute` | `-` | `local` | `local` | `local` | `local` | `local` | `local` |
| `don-time` | `-` | `local` | `local` | `local` | `local` | `local` | `local` |
| `evm` | `-` | `local (1337,2337)` | `local (1337,2337)` | `local (1337,2337)` | `local (1337,2337)` | `local (1337,2337)` | `local (1337,2337)` |
| `http-action` | `-` | `local` | `-` | `-` | `-` | `-` | `-` |
| `http-trigger` | `-` | `local` | `-` | `-` | `-` | `-` | `-` |
| `ocr3` | `-` | `local` | `local` | `local` | `local` | `local` | `local` |
| `read-contract` | `-` | `local (1337,2337)` | `local (1337,2337)` | `local (1337,2337)` | `local (1337,2337)` | `local (1337,2337)` | `local (1337,2337)` |
| `vault` | `-` | `local` | `-` | `-` | `-` | `-` | `-` |
| `web-api-target` | `-` | `local` | `local` | `local` | `local` | `local` | `local` |
| `web-api-trigger` | `-` | `local` | `local` | `local` | `local` | `local` | `local` |
| `write-evm` | `-` | `local (1337,2337)` | `local (1337,2337)` | `local (1337,2337)` | `local (1337,2337)` | `local (1337,2337)` | `local (1337,2337)` |

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

### `shard2`

- Types: `shard`, `workflow`
- Nodes: `4`
- Roles: `plugin`
- EVM chains: `1337,2337`
- Exposes remote capabilities: `false`

### `shard3`

- Types: `shard`, `workflow`
- Nodes: `4`
- Roles: `plugin`
- EVM chains: `1337,2337`
- Exposes remote capabilities: `false`

### `shard4`

- Types: `shard`, `workflow`
- Nodes: `4`
- Roles: `plugin`
- EVM chains: `1337,2337`
- Exposes remote capabilities: `false`

### `shard5`

- Types: `shard`, `workflow`
- Nodes: `4`
- Roles: `plugin`
- EVM chains: `1337,2337`
- Exposes remote capabilities: `false`

