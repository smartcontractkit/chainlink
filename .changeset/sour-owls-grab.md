---
"chainlink": patch
---

Added config option `HeadTracker.FinalityTagBypass` to force `HeadTracker` to track blocks up to `FinalityDepth` even if `FinalityTagsEnabled = true`. This option is a temporary measure to address high CPU usage on chains with extremely large actual finality depth (gap between the current head and the latest finalized block). #added

Added config option `HeadTracker.MaxAllowedFinalityDepth` maximum gap between current head to the latest finalized block that `HeadTracker` considers healthy. #added
