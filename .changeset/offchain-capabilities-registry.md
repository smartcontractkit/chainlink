---
"chainlink": minor
---

#added Offchain capabilities registry: `cresettings` jobs accept `config_type = "capabilities_registry"` with a versioned `offchain_config` (proto JSON). With `[Capabilities.Local] UseOffchainRegistry = true` (default `false`), each capability's `spec_config` is resolved per (DON, capability) as TOML < on-chain < offchain, falling back to the legacy value for omitted keys; offchain-only updates reconfigure affected capabilities. At most one `capabilities_registry` job may exist, payload versions must increase (enforced durably, across job deletion and restarts), and the runtime config follows committed job state only, so rolled-back job changes have no effect. Deleting the job reverts to on-chain/TOML config. #db_update
