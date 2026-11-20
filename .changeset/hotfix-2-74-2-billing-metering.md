---
"chainlink": patch
---

#changed Bump capability-evm/solana/aptos for the gas meter beholder logger: gas meter log lines now also emit as Beholder events, matching the workflow engine's dual logging
#changed Meter confidential workflow executions under cre:workflow:confidential-compute instead of cre:workflow:compute, so enclave compute is separately visible and billable downstream
