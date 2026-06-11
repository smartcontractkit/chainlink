# mockchip — Local Chip Ingress Mock for DurableEmitter

`mockchip` is an in-process mock of the [Chip Ingress](https://github.com/smartcontractkit/chainlink-common/tree/main/pkg/chipingress)
gRPC service. It exists to verify `DurableEmitter` behaviour end-to-end on a
developer machine — and to deliberately simulate a Chip outage so the
durable queue + retransmit path can be observed.

It can be used two ways:

1. **As a library** (`mockchip.NewServer`) from Go tests — see
   `server_test.go` for the canonical patterns.
2. **As a standalone binary** (`cmd/mockchipendpoint`) bound to a local TCP
   port and optionally exposed to the public internet via
   [`ngrok tcp`](https://ngrok.com/docs/tcp/).

## What it does

- Implements `chippb.ChipIngressServer` (`Publish`, `PublishBatch`, `Ping`).
- Captures every received `CloudEvent` in memory together with the time it
  was received, for later inspection.
- Exposes an HTTP control plane to (a) list captured events, (b) flip an
  outage flag that causes every publish RPC to return
  `codes.Unavailable`, (c) reset captured state, and (d) read summary stats.

When the outage flag is active, `DurableEmitter` should **not** lose events
— they remain in its `DurableEventStore` (Postgres or in-memory) and are
re-delivered by the retransmit loop once the outage is cleared. Flipping
the flag back to off is the verification that the queue actually drains.

## Running locally

From the repo root:

```bash
go run ./core/services/durableemitter/mockchip/cmd/mockchipendpoint \
    -grpc :9095 -http :9096
```

Output:

```
mockchip: gRPC listening on [::]:9095
mockchip: HTTP control plane on http://[::]:9096
mockchip: outage_active=false
mockchip: expose to the internet with `ngrok tcp 9095`
```

### Exposing the gRPC endpoint with ngrok

The mock listens on plain TCP (no TLS) and Chip Ingress uses gRPC, so a
plain HTTP tunnel is not enough — you need a TCP tunnel:

```bash
ngrok tcp 9095
```

ngrok prints a forwarding address such as `tcp://4.tcp.ngrok.io:18342`.
Strip the `tcp://` scheme and pass `4.tcp.ngrok.io:18342` to whatever
configures `DurableEmitter`'s Chip Ingress endpoint (`SetupConfig.Endpoint`),
together with `InsecureConnection: true` (the mock does not terminate TLS).

> ngrok's free plan supports **one** simultaneous TCP tunnel. If the HTTP
> control plane also needs to be reachable externally, expose it through a
> second `ngrok http 9096` tunnel.

## HTTP control plane

| Method | Path             | Description                                                |
| ------ | ---------------- | ---------------------------------------------------------- |
| GET    | `/healthz`       | 200 OK liveness probe                                      |
| GET    | `/events`        | JSON list of every captured event (oldest first)           |
| GET    | `/events/count`  | Plaintext captured count (useful from shell loops)         |
| GET    | `/stats`         | JSON summary: captured count, RPC counters, outage flag    |
| POST   | `/outage/on`     | Begin failing every `Publish`/`PublishBatch` RPC           |
| POST   | `/outage/off`    | Restore normal behaviour; subsequent retransmits succeed   |
| POST   | `/reset`         | Clear captured events and counters (outage flag preserved) |

### Typical drain-verification flow

```bash
# 1. Start mock + ngrok in two terminals.
go run ./core/services/durableemitter/mockchip/cmd/mockchipendpoint -grpc :9095 -http :9096
ngrok tcp 9095

# 2. Point your DurableEmitter at the ngrok TCP address (insecure connection).

# 3. Confirm normal flow — events show up.
curl -s localhost:9096/events/count

# 4. Simulate outage. RPCs will start failing; the durable queue grows.
curl -s -X POST localhost:9096/outage/on
curl -s localhost:9096/stats        # outage_active=true, failed_calls climbs

# 5. End the outage. The retransmit loop should drain the queue.
curl -s -X POST localhost:9096/outage/off
watch -n1 'curl -s localhost:9096/events/count'

# 6. Compare captured count with the total events DurableEmitter emitted.
curl -s localhost:9096/events | jq 'length'
```

## Programmatic use

```go
srv := mockchip.NewServer()
addr, err := srv.Start("127.0.0.1:0")   // ":0" -> pick a free port
require.NoError(t, err)
defer srv.Stop()

// ... wire a DurableEmitter at `addr` with WithInsecureConnection ...

srv.SetOutage(true)
// emit a bunch of events; they queue durably
srv.SetOutage(false)

// Wait for DurableEmitter's retransmit loop to deliver them all.
require.NoError(t, srv.WaitFor(ctx, expected))
events := srv.Captured()
```

## Notes / caveats

- The mock is **for development and integration testing only**. It accepts
  every event unconditionally (no schema validation, no auth).
- All captured events live in process memory — there is no persistence.
  Restart = clean slate.
- gRPC uses plain TCP. If you front the mock with TLS-terminating
  infrastructure you must point the client at that infrastructure, not at
  the mock directly.
- The CloudEvent payload is preserved as a `*cepb.CloudEvent`. The HTTP
  `/events` endpoint base64-encodes binary payloads (`data` field) per Go's
  default JSON encoding for `[]byte`.
