# DON Topology

- Config: `configs/workflow-don-tron.toml`
- Class: `single-don`
- Infra: `docker`

## Capability Matrix

This matrix is the source of truth for capability placement by DON.

| Capability | `bootstrap-gateway` | `workflow` |
|---|---|---|
| `consensus` | `-` | `local` |
| `cron` | `-` | `local` |
| `custom-compute` | `-` | `local` |
| `don-time` | `-` | `local` |
| `http-action` | `-` | `local` |
| `http-trigger` | `-` | `local` |
| `ocr3` | `-` | `local` |
| `read-contract` | `-` | `local (1337,3360022319)` |
| `vault` | `-` | `local` |
| `web-api-target` | `-` | `local` |
| `web-api-trigger` | `-` | `local` |
| `write-evm` | `-` | `local (1337,3360022319)` |

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
- EVM chains: `1337,3360022319`
- Exposes remote capabilities: `false`

