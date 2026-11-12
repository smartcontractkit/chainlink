---
"chainlink": patch
---

Add `[Telemetry] LoopIDEnabled` (default `false`). When enabled, each LOOP plugin gets a telemetry-only identity via `CL_TELEMETRY_LOOP_ID`, which it exports as the `chainlink.loop.id` resource attribute, e.g. `Relayer.evm.1` or `StandardCapabilities.<jobID>.<jobName>`. This keeps series from different LOOPs on the same node from merging. Add `[Telemetry] OTelResourceAttributes` (default `['process.executable.name', 'service.version']`, `process.pid` also supported): process attributes the node and every LOOP add to their telemetry resource, whether or not tracing is enabled. #added
