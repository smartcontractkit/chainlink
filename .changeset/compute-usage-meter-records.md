---
"chainlink": patch
---

Emit a `cre:workflow:compute` usage MeterRecord per workflow execution when `[Metering].MeterRecordsEnabled` is on. The compute duration is reported by the module that ran the workflow (the guest run for WASM, the enclave-measured run for confidential workflows). #internal
