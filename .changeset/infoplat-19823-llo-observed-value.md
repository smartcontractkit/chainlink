---
"chainlink": minor
---

#added llo_datasource_stream_observed_value gauge: latest observed value per LLO stream (value_type benchmark|bid|ask, config_digest). To forward it to Beholder, set `[Telemetry.PrometheusBridge] Enabled = true` and add `llo_datasource_stream_observed_value` to `Prefixes`; setting `Prefixes` replaces the default `["go_"]`, so keep `"go_"` and any prefixes already forwarded.
