---
"chainlink": patch
---

Telemetry: with `TelemetryIngress.ChipIngressEnabled = true`, the node now sends legacy telemetry to chip-ingress through its own CSA-authenticated client and one shared batch client, instead of relying on Beholder's chip client (which silently dropped telemetry unless `Telemetry.ChipIngressEndpoint` was also set). Added `TelemetryIngress.ChipIngressEndpoint` (default `legacy-telemetry.prod.telemetry.chain.link:443`) and `TelemetryIngress.ChipIngressInsecureConnection`. `TelemetryIngress.Endpoints` entries still select which chains send telemetry. `telemetry_client_messages_sent{endpoint="chip-ingress"}` now counts messages instead of batches. #added
