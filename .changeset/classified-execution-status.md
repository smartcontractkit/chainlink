---
"chainlink": minor
---

#added The v2 workflow engine now emits `ClassifiedExecutionStatus` on `WorkflowExecutionFinished` events, distinguishing failures caused by the user's workflow (`USER_ERROR`) from platform/infrastructure failures (`SYSTEM_ERROR`). The v1 engine is unaffected.
