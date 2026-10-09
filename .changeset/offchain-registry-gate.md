---
"chainlink": minor
---

#added `[Capabilities.Local] UseOffchainRegistry` config flag (default `false`), gating the offchain capabilities registry cutover. #removed the `[Capabilities.Local.Capabilities.<id>]` per-capability TOML table entirely, including the `Config` key/value overrides and `BinaryPathOverride`; capability config now comes from the on-chain registry (and, with the gate, the offchain registry), and the capability binary is always resolved from the capability ID.
