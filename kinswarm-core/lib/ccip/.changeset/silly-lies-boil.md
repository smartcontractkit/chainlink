---
"chainlink": minor
---

Make websocket URL `WSURL` for `EVM.Nodes` optional, and apply logic so that:
* If WS URL was not provided, SubscribeFilterLogs should fail with an explicit error
* If WS URL was not provided LogBroadcaster should be disabled 
#nops
