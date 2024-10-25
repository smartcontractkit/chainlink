---
'@chainlink/contracts': minor
---

#internal Add a missing condition for the Execution plugin in the \_afterOCR3ConfigSet function. Now, the function correctly reverts if signature verification is enabled for the Execution plugin


PR issue: CCIP-3799

Solidity Review issue: CCIP-3966