---
"chainlink": patch
---

#internal Add node-measured round-trip metrics to the confidential workflows ConfidentialModule: `platform_engine_confidential_execution_time_ms` (histogram) and `platform_engine_confidential_execution_failures` (counter), labeled by workflow. These are trusted (node-measured) and complement the non-attested `enclave.*` metrics forwarded from the enclave.
