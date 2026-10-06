---
"chainlink": minor
---

#added Offchain capabilities registry `method_configs` cutover: with `[Capabilities.Local] UseOffchainRegistry = true`, each capability's Don2Don `method_configs` are applied on top of the on-chain ones (offchain-wins per method; methods the payload omits keep their on-chain config, and capabilities/DONs absent offchain are unchanged). Payloads carrying `method_configs` are validated at ingestion, and the cross-validation `config_mismatch` divergence now also covers method_configs changes. Gate off (default) or no payload: behavior unchanged.
