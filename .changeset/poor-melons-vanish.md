---
"chainlink": minor
---

Add the `pool_rpc_node_highest_finalized_block` metric that tracks the highest finalized block seen per RPC. If `FinalityTagEnabled = true`, a positive `NodePool.FinalizedBlockPollInterval` is needed to collect the metric. If the finality tag is not enabled, the metric is populated with a calculated latest finalized block based on the latest head and finality depth.
