---
"chainlink": minor
---

#added `TelemetryIngress.BeholderEnabled` (off by default): sends node telemetry through Beholder's per-domain emitters instead of the WSRPC `TelemetryIngress.Endpoints`. Requires `Telemetry.Enabled` and `Telemetry.ChipIngressEndpoint`; cannot be combined with `TelemetryIngress.ChipIngressEnabled`.
