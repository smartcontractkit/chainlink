---
"chainlink": patch
---

#bugfix Vault nodes serve each gateway request in its own goroutine instead of inline in the gateway connection's read loop, so a request waiting on authorization no longer delays every later request from that gateway
