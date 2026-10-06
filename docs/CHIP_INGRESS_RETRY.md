---
title: Chip Ingress Retry Policy
sidebar_position: 4
---

# Chip Ingress retry policy

The node's CHIP ingress publishes (custom events, workflow telemetry, heartbeats) go over gRPC.
Without retries, every transient publish failure — an `UNAVAILABLE` disconnect, a
`RESOURCE_EXHAUSTED` server under pressure — drops the batch permanently. gRPC-level retries
can be enabled per node in Core TOML, scoped to the `ChipIngress` service:

```text
[Telemetry]
ChipIngressRetryEnabled = true
```

Off by default. When enabled, unset knobs fall back to the recommended policy
(3 attempts, 100ms initial backoff doubling to a 1s cap, retrying `Unavailable` and
`ResourceExhausted`), paired with gRPC retry throttling so a struggling server is not
amplified:

```text
[Telemetry]
ChipIngressRetryEnabled = true
ChipIngressRetryMaxAttempts = 5
ChipIngressRetryInitialBackoff = '100ms'
ChipIngressRetryMaxBackoff = '1s'
ChipIngressRetryBackoffMultiplier = 2.0
ChipIngressRetryableStatusCodes = ['Unavailable', 'ResourceExhausted']
```

- `ChipIngressRetryMaxAttempts` counts the original attempt and must be at least 2; gRPC
  silently installs no retry policy below that.
- `ChipIngressRetryableStatusCodes` uses `codes.Code.String()` spelling
  (e.g. `Unavailable`, `ResourceExhausted`, `DeadlineExceeded`); unknown names fail node
  startup instead of being ignored.
- Retries share the caller's RPC deadline; they cannot rescue calls that already exhausted it.

Every value can also be set via environment variable, e.g.
`CL_TELEMETRY_CHIP_INGRESS_RETRY_ENABLED=true`, and the same settings are passed to LOOP
plugin subprocesses (their CHIP publishes use the same policy).

Retried publishes can duplicate events: an `UNAVAILABLE` may be returned after the server
already did the work. Event producers that set the `idempotencykey` extension are protected;
otherwise weigh enabling retries against duplicate tolerance for your event types. The
durable emitter path is unaffected — it already guarantees delivery via its own retransmit
loop and is intentionally left without gRPC retries.
