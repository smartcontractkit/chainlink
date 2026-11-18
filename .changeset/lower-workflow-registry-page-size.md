---
"chainlink": patch
---

#bugfix Lower the Workflow Registry syncer page size from 1000 to 250 so a single `getWorkflowListByDON` eth_call stays under the common 16,777,216 gas cap and registry reconciliation no longer fails once a DON family indexes enough workflows.
