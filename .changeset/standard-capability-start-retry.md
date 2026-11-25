---
"chainlink": patch
---

#bugfix A standard capability that fails to start is retried until it starts, instead of staying down until the node restarts. A retry after a failed Initialise runs on a new plugin process.
