---
"chainlink": patch
---

Fixed a race condition bug around EVM nonce management, which could cause the Node to skip a nonce and get stuck.
