---
"chainlink": patch
---

#internal Add node-measured round-trip metrics to the confidential workflows ConfidentialModule: `enclave_execution_time_ms` (histogram) and `enclave_execution_failures` (counter), labeled by workflow. These are trusted (node-measured) and complement the non-attested `enclave.*` metrics forwarded from the enclave.
