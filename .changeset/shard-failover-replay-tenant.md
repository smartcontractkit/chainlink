---
"chainlink": patch
---

#bugfix Sharded workflow DON failover now replays cached trigger events with the workflow's CRE tenant on the context, so replays no longer fail with `missing tenant for scope: workflow`.
