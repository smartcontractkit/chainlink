---
"chainlink": minor
---

#added Add countNilsAsFaults flag to MedianTask. When enabled, nil values are counted toward allowedFaults and filtered out before median calculation, preventing nils from crashing the task while preserving fault
