---
"chainlink": patch
---

Set `NodePool.EnforceRepeatableRead = true` by default for all chains. This forces Core to stop using RPCs behind on the latest finalized block. #changed #nops
