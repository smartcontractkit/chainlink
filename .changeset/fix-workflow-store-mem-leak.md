---
"chainlink": patch
---

#bugfix Rebuild the in-memory workflow execution store's map when pruning so old bucket storage becomes eligible for GC. Go maps never shrink after deletes, which stranded memory as the store churned through millions of executions.
