---
"chainlink": patch
---

#bugfix CRE: confidential-workflows and vault GetSecrets requests now carry the execution's DonTime-derived `ExecutionTimestamp` in `RequestMetadata`, so capability DONs evaluate atomic feature flags against the real execution time instead of the zero time.
