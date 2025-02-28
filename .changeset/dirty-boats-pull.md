---
"chainlink": patch
---

Add support for custom deviation functions in median plugin. #added

Do this like so, by adding the following config to a "median" job spec:

```toml
[pluginConfig.deviationFunc]
expiresAt = 1739895051.0 # REQUIRED. Unix timestamp indicating the expiry date. Should be specified as a float64 so even for integer values, add a decimal point.
type = "pendle" # REQUIRED. Currently only "pendle" is supported.
multiplier = "1000000" # OPTIONAL. Must be supplied as string integer. Default is 1e18 if omitted.
```
