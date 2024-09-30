---
"chainlink": patch
---

- register polling subscription to avoid subscription leaking when rpc client gets closed.
- add a temporary special treatment for SubscribeNewHead before we replace it with SubscribeToHeads. Add a goroutine that forwards new head from poller to caller channel.
- fix a deadlock in poller, by using a new lock for subs slice in rpc client.
#bugfix
