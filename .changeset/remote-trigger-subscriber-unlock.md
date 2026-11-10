---
"chainlink": patch
---

#bugfix Remote trigger subscriber no longer leaves its mutex held when an event arrives for a trigger that was unregistered between the workflow and registration lookups, which blocked all later events, registrations and shutdown for that capability on the node.
