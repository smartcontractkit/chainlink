---
"chainlink": minor
---

#added Offchain capabilities registry `method_configs` cutover: with `[Capabilities.Local] UseOffchainRegistry = true`, each capability's Don2Don `method_configs` are applied on top of the on-chain ones (offchain-wins per method; methods the payload omits keep their on-chain config, and capabilities/DONs absent offchain are unchanged). Payloads carrying `method_configs` are validated at ingestion, and the cross-validation `config_mismatch` divergence now also covers method_configs changes. Gate off (default) or no payload: behavior unchanged.
#fixed Fixed a lost-wakeup race in the workflow DON notifier (`donNotifier.Subscribe`): a `NotifyDonSet` arriving between a subscriber's initial value check and its channel registration was missed entirely, leaving the subscriber (workflow engines, shard failover manager, or `WaitForDon` callers) blocked indefinitely. The subscriber channel is now registered before the current value is read, so a concurrent notify is either delivered by the broadcast or observed by the read.
