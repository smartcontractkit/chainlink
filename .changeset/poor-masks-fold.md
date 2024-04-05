---
"chainlink": minor
---

Move JuelsPerFeeCoinCacheDuration under JuelsPerFeeCoinCache struct in config. Rename JuelsPerFeeCoinCacheDuration to updateInterval. Add stalenessAlertThreshold to  JuelsPerFeeCoinCache config.
StalenessAlertThreshold cfg option has a default of 24 hours which means that it doesn't have to be set unless we want to override the duration after which a stale cache should start throwing errors.
