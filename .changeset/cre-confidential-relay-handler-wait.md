---
"chainlink": patch
---

#internal Confidential relay: on the workflow node, briefly wait (after attestation and Workflow-DON authorization pass) for a not-yet-registered execution handler before failing the enclave's relay callback. This lets a node that has not yet started its copy of the DON-shared execution register and sign in time, instead of failing the callback outright and eroding the relay quorum. The wait is bounded so a callback for an execution the node never runs still fails promptly.
