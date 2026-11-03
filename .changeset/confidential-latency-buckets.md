---
"chainlink": patch
---

Give the confidential relay request latency and confidential workflows enclave execution histograms explicit buckets up to 90s and 300s. The SDK default buckets stopped at 10s, which pinned their P99 at 10,000ms. #bugfix
