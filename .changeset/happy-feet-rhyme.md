---
"chainlink": minor
---

This PR introduce few changes:
- Add a new flag support called `NewHeadsPollInterval` (0 by default indicate disabled), which is an interval for polling new block periodically using http client rather than subscribe to ws feed. 
- Updated new head handler for polling new head over http, and register the subscription in node lifecycle logic.
- If the polling new heads is enabled, WS new heads subscription will be replaced with the new http based polling.

Note: There will be another PR for making WS URL optional with some extra condition.
#added
