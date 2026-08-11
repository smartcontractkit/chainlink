---
"chainlink": patch
---

#internal Move `BuildWorkflowGetSecretsRequestID` to chainlink-common's vault capability package and use it in place of the `vaultutils` copy, so all consumers derive the VaultDON GetSecrets request ID from a single definition.
