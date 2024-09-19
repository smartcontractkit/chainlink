---
"chainlink": minor
---

Adding feature flag for `LogBroadcaster` called `LogBroadcasterEnabled`, which is `true` by default to support backwards compatibility.
Adding `LogBroadcasterEnabled` allows certain chains to completely disable the `LogBroadcaster` feature, which is an old feature (getting replaced by logPoller) that only few products are using it:
* OCR1 Median
* *OCR2 Median when ChainReader is disabled
* *pre-OCR2 Keeper 
* Flux Monitor 
* Direct RequestOCR1 Median

#added
