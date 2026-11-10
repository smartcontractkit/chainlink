---
"chainlink": patch
---

Add `[Telemetry] LoopIDEnabled` (default `false`). When enabled, each LOOP plugin gets a telemetry-only identity via `CL_TELEMETRY_LOOP_ID`, which it exports as the `chainlink.loop.id` resource attribute, e.g. `Relayer.evm.1` or `StandardCapabilities.<jobID>.<jobName>`. This keeps series from different LOOPs on the same node from merging. #added
