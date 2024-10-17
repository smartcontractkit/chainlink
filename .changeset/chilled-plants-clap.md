---
"chainlink": patch
---

The findBroadcastedAttempts in detectStuckTransactionsHeuristic can returns uninitialized struct that potentially cause nil pointer error. Changed the return type of findBroadcastedAttempts to be pointers and added nil pointer check. #bugfix
