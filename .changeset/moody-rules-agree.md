---
"chainlink": patch
---

- register polling subscription logic in rpc client so when node unhealthy new susbcription will be used
- add a temporary special treatment for SubscribeNewHead before we replace it with SubscribeToHeads. Add a goroutine that forwards new head from poller to caller channel.
#fixed
