---
"chainlink": patch
---

#bugfix VRF: fulfillment transactions now reserve gas for the coordinator's post-callback accounting, so fulfillments complete even when a consumer's callback uses its entire callback gas limit.
