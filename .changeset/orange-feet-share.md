---
"chainlink": minor
---

Implemented new EVM Multinode design. The Multinode is now called by chain clients to retrieve the best healthy RPC rather than performing RPC calls directly.
Multinode performs verious health checks on RPCs, and in turn increases reliability.
This new EVM Multinode design will also be implemented for non-EVMs chains in the future.
#updated #changed #internal
