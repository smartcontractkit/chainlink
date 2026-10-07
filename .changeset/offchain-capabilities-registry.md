---
"chainlink": minor
---

#added Offchain capabilities registry: `cresettings` jobs accept `config_type = "capabilities_registry"` in their `settings` TOML, with the versioned `OffchainCapabilitiesRegistry` proto-JSON payload under `offchain_config` (`domain`, `env`, `version`, and per-DON `capabilities` keyed by on-chain DON name). Unknown payload fields are rejected and versions must increase. With `[Capabilities.Local] UseOffchainRegistry = true` (default `false`), each capability's `spec_config` is resolved per (DON, capability) as on-chain < offchain, falling back to the on-chain value for omitted keys; offchain-only updates reconfigure affected capabilities. Like the other `cresettings` config types, at most one `capabilities_registry` job may be active; deleting it reverts to on-chain config. Emits `platform_cap_config_applied_version` and `platform_cap_config_apply_errors_total`.
