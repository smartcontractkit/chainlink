# DON Topology

- Config: `configs/workflow-gateway-capabilities-don.toml`
- Class: `multi-don`
- Infra: `docker`

## Capability Matrix

This matrix is the source of truth for capability placement by DON.

| Capability | `bootstrap-gateway` | `capabilities` | `workflow` |
|---|---|---|---|
| `consensus` | `-` | `-` | `local` |
| `cron` | `-` | `-` | `local` |
| `custom-compute` | `-` | `-` | `local` |
| `don-time` | `-` | `-` | `local` |
| `evm` | `-` | `remote-exposed (2337)` | `local (1337)` |
| `http-action` | `-` | `-` | `local` |
| `http-trigger` | `-` | `-` | `local` |
| `ocr3` | `-` | `-` | `local` |
| `read-contract` | `-` | `remote-exposed (2337)` | `local (1337)` |
| `vault` | `-` | `remote-exposed` | `-` |
| `web-api-target` | `-` | `remote-exposed` | `-` |
| `web-api-trigger` | `-` | `-` | `local` |
| `write-evm` | `-` | `remote-exposed (2337)` | `local (1337)` |

## DONs

### `bootstrap-gateway`

- Types: `bootstrap`, `gateway`
- Nodes: `1`
- Roles: `bootstrap`, `gateway`
- EVM chains: `1337,2337`
- Exposes remote capabilities: `false`

### `capabilities`

- Types: `capabilities`
- Nodes: `4`
- Roles: `plugin`
- EVM chains: `1337,2337`
- Exposes remote capabilities: `true`

### `workflow`

- Types: `workflow`
- Nodes: `4`
- Roles: `plugin`
- EVM chains: `1337,2337`
- Exposes remote capabilities: `false`

