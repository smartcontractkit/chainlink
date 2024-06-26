---
"chainlink": minor
---

#changed Refactored the BlockHistoryEstimator check to prevent excessively bumping transactions. Check no longer waits for CheckInclusionBlocks to pass before assessing an attempt.
#bugfix Fixed a bug that would use the oldest blocks in the cached history instead of the latest to perform gas estimations.
