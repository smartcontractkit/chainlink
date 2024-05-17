---
"chainlink": patch
---

Fixed CPU usage issues caused by inefficiencies in HeadTracker.

HeadTracker's support of finality tags caused a drastic increase in the number of tracked blocks on the Arbitrum chain (from 50 to 12,000), which has led to a 30% increase in CPU usage. 

The fix improves the data structure for tracking blocks and makes lookup more efficient. BenchmarkHeadTracker_Backfill shows 40x time reduction.
#bugfix
