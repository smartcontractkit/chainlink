---
"chainlink": patch
---

#bugfix VRF: fulfillment transactions now reserve gas for the coordinator's accounting, so fulfillments complete even when a consumer's callback uses its entire callback gas limit. #db_update Adds an optional `gasOverhead` VRF job spec field (`vrf_specs.gas_overhead` column, nullable) to override the reserved overhead per job.
