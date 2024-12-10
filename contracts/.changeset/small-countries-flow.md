---
'@chainlink/contracts': patch
---

Modify TokenPool.sol function setChainRateLimiterConfig to now accept an array of configs and set sequentially. Requested by front-end. PR issue CCIP-4329 #bugfix
