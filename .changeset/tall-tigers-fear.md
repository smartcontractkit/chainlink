---
"chainlink": minor
---

RPCs that sustain polling error rates above 50% will now eventually be marked as unreachable, in addition to previous behaviour of `PollFailureThreshold` failures in a row. #updated #nops
