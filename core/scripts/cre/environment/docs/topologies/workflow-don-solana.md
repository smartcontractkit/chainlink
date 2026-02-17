# DON Topology

- Config: `configs/workflow-don-solana.toml`
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
| `http-action` | `-` | `-` | `local` |
| `http-trigger` | `-` | `-` | `local` |
| `ocr3` | `-` | `-` | `local` |
| `solana` | `-` | `remote-exposed` | `-` |
| `vault` | `-` | `remote-exposed` | `-` |
| `web-api-target` | `-` | `remote-exposed` | `-` |
| `web-api-trigger` | `-` | `-` | `local` |

## DONs

### `bootstrap-gateway`

- Types: `bootstrap`, `gateway`
- Nodes: `1`
- Roles: `bootstrap`, `gateway`
- EVM chains: `1337`
- Exposes remote capabilities: `false`

### `capabilities`

- Types: `capabilities`
- Nodes: `4`
- Roles: `plugin`
- EVM chains: `1337`
- Solana chains: `22222222222222222222222222222222222222222222`
- Exposes remote capabilities: `true`

### `workflow`

- Types: `workflow`
- Nodes: `4`
- Roles: `plugin`
- EVM chains: `1337`
- Solana chains: `22222222222222222222222222222222222222222222`
- Exposes remote capabilities: `false`

