---
"chainlink": patch
---

Fixed local finality violation caused by an RPC lagging behind on latest finalized block.

Added `EVM.FinalizedBlockOffset` and `EVM.NodePool.EnforceRepeatableRead` config options.
With `EnforceRepeatableRead = true`, RPC is considered healthy only if its most recent finalized block is larger or equal to the highest finalized block observed by the Node minus `FinalizedBlockOffset`.
#bugfix
