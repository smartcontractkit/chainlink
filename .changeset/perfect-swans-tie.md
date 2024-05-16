---
"ccip": patch
---

Get liquidity managers token in bridge constructor and save for reuse, upon calls to GetTransfer() use the address to compare with remoteToken and localToken.
