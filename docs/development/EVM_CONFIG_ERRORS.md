# EVM configuration errors (operational behavior)

## Error classes

| Class | Examples | Node behavior |
|-------|-----------|----------------|
| **Fatal (abort startup)** | Duplicate EVM chain ID in TOML, database unreachable, invalid global secrets | `NewApplication` fails; process exits with error |
| **Per-chain skip** | Enabled chain fails `NewTOMLChain` (e.g. RPC client error), or `NewRelayer` fails (e.g. chain ID missing from embedded `chain-selectors`) | Other enabled chains start; skipped chains are logged with `error_class`; Prometheus `evm_chain_config_skipped_total{reason=...}` increments; `evm_chain_config_degraded` set to 1; Beholder `platform_node_evm_chain_config_skipped_total` increments |
| **Readiness** | Any per-chain skip occurred | `EVMChainConfigHealth` readiness check fails until configuration or binary image is fixed |

## Alerting

- **Prometheus**: alert on `evm_chain_config_degraded == 1` or rate on `evm_chain_config_skipped_total`.
- **Beholder / platform**: use `platform_node_evm_chain_config_skipped_total` for platform dashboards and paging.

## Runbook (when an alert fires)

1. Identify which chain IDs were skipped from node logs (`skipping EVM chain` / `skipping EVM relayer`).
2. Compare configured chain IDs with the `chain-selectors` version in the running image (upgrade image or roll back config).
3. Confirm `/health?full` shows `EVMChainConfigHealth` failing until all enabled chains load.

## Release hygiene

Pin and bump `github.com/smartcontractkit/chain-selectors` whenever NOP configs add new chains so the image always includes required selector mappings.
