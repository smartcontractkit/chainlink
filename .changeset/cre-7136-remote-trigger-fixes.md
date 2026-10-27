---
"chainlink": patch
---

#bugfix Remote trigger fixes: the publisher forwards ACKs to the underlying trigger once per event, and the subscriber no longer panics or blocks when a trigger is unregistered while an event is being delivered.
