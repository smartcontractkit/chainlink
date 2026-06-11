---
"chainlink": patch
---

#internal Confidential workflows: stop setting the deprecated outside-envelope `ConfidentialWorkflowRequest.binary_url`. `binary_url` stays in the hashed `WorkflowExecution` (PublicData); the enclave reads it there.
