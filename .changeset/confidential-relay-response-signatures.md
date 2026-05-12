---
"chainlink": patch
---

Confidential relay handlers now sign relay-DON secrets and capability responses before returning them toward the enclave path, and the gateway aggregator buckets per-node responses by canonical logical hash so F+1 unique signers form quorum over a shared logical payload.

#added
