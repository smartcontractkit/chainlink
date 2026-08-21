---
title: OTel metric export batching
sidebar_position: 4
---

# OTel metric export batching

When a node's OTLP metric collection can exceed the OTel Collector gRPC
receive limit, configure SDK-side metric export batching at deployment time:

```text
OTEL_GO_X_METRIC_EXPORT_BATCH_SIZE=<positive integer>
```

The value is the maximum number of metric data points in each OTLP exporter
request. It is count-based, not a serialized-byte limit. Large attributes,
histograms, or one oversized data point can still exceed a receiver limit.

This is an experimental OpenTelemetry Go SDK setting. It is process-wide and
is read when each `PeriodicReader` is constructed; changing the environment
after node startup does not reconfigure an existing reader. An unset, invalid,
zero, or negative value preserves unbatched export behavior.

The setting must be supplied before the node constructs its telemetry client.
Do not change production values without measuring request sizes, export
latency, failure rate, and node/collector resource usage. SDK batching protects
the node-to-collector request; a collector `batch` processor operates after
receipt and is not a substitute. Collector receive-limit changes are owned by
the deployment/infra team.

Node test and CRE environment definitions pass arbitrary node environment
variables through, so no packaging allowlist change is required for this
setting.
