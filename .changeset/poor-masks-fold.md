---
"chainlink": major
---

Add juels fee per coin cache freshness alert duration to config. This cfg option has a default of 24 hours which means that it doesn't have to be set unless we want to override the duration after which a stale cache should start throwing errors.
