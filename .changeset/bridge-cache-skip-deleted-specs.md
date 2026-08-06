---
"chainlink": patch
---

#bugfix Bridge cache: skip bridge responses whose pipeline spec has been deleted when flushing the cache to `bridge_last_value`. The cache holds responses in memory for the lifetime of the node, so after a job was deleted it kept offering rows for the removed spec. Since the flush is a single multi-row insert, one such row failed the whole batch on `bridge_last_value_spec_id_fkey`, logging a foreign key violation every upsert interval (5s by default) and preventing live bridge responses from being persisted.
