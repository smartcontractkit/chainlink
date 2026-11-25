---
"chainlink": patch
---

#bugfix VRF: fulfillment transactions now reserve gas for the coordinator's accounting, so fulfillments complete even when a consumer's callback uses its entire callback gas limit, and the subscription balance pre-check now reserves the worst-case request cost (full callback gas limit at the gas lane maximum price, including the premium). #db_update Adds optional `gasOverhead` and `premiumPercentage` VRF job spec fields (`vrf_specs.gas_overhead` and `vrf_specs.premium_percentage` columns, nullable) to override the reserved overhead and premium per job.
