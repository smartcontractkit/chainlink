---
"chainlink": patch
---

#internal Label the `enclave_execution_failures` counter with `error_type` ("user" or "system") so alerts on confidential-workflow enclave failures can exclude user-caused ones, such as a workflow that exceeds its enclave execution budget.
