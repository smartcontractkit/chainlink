---
"chainlink": minor
---

#added `[Capabilities.Local] UseOffchainRegistry` config flag (default `false`), gating the offchain capabilities registry cutover. #removed `[Capabilities.Local.Capabilities.<id>.Config]` node-local TOML capability config overrides; capability config now comes from the on-chain registry (and, with the gate, the offchain registry). `BinaryPathOverride` is unchanged.
