---
"chainlink": patch
---

#added `Telemetry.MetricExportBatchSize` config limiting metric data points per OTLP export request (0 disables batching); passed to the Beholder client and LOOP plugins
