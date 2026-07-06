---
"chainlink": minor
---

Add the ConfidentialWorkflows.Enabled feature gate to the v2 workflow engine. Confidential workflow execution can now be toggled per workflow/owner/org/global via the settings registry; when disabled, ConfidentialModule.Execute rejects the request.
