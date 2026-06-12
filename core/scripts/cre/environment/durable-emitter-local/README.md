# DurableEmitter local validation stack

Self-contained local harness for sanity-checking the durable-emitter metrics
pipeline end-to-end: **driver → OTel collector → Prometheus → Grafana**. No
Chainlink Testing Framework or live Chip Ingress required.

## What's in the box

- `docker-compose.yml` — Postgres, OTel collector, Prometheus, Grafana
- `otel-collector-config.yaml` — receives OTLP on `:4317`/`:4318`, exports to
  Prometheus on `:8889`
- `prometheus.yml` — scrapes the collector's Prom exporter
- `postgres/init.sql` — creates the `cre.chip_durable_events` table (mirrors
  migration `0300_chip_durable_events.sql`)
- `grafana/provisioning/` — datasource + dashboard provider
- `grafana/dashboards/durable_emitter.json` — generated from
  [chainlink-observability/resources/durable_emitter](https://github.com/smartcontractkit/chainlink-observability/tree/main/resources/durable_emitter)
  (dashboard-as-code; not committed in chainlink)
- `driver/main.go` — small Go program that builds a real `DurableEmitter` with
  metrics enabled, uses `chipingress.NoopClient` as the transport, and drives
  events at a configurable rate

## Run procedure

### 1. Generate the Grafana dashboard

The dashboard lives in **chainlink-observability** as Go code. Render the local
load-test JSON before starting Grafana:

```bash
git clone https://github.com/smartcontractkit/chainlink-observability.git
cd chainlink-observability

mkdir -p /path/to/chainlink/core/scripts/cre/environment/durable-emitter-local/grafana/dashboards

go run ./cmd/generate-durable-emitter-dashboard/main.go \
  --local-load-test \
  --format grafana \
  --dashboard-uid durable-emitter-load-test \
  --output /path/to/chainlink/core/scripts/cre/environment/durable-emitter-local/grafana/dashboards/durable_emitter.json
```

### 2. Bring up the stack

```bash
cd core/scripts/cre/environment/durable-emitter-local
docker compose up -d
```

Wait until Postgres is healthy (`docker compose ps` shows `(healthy)`).
Grafana lands at <http://localhost:3000> with anonymous admin login enabled.

### 3. Build and run the driver

The driver lives in the `core/scripts` Go module (separate from the main
`chainlink/v2` module), so build from that directory:

```bash
cd core/scripts
go build -o /tmp/de-driver ./cre/environment/durable-emitter-local/driver
/tmp/de-driver --rate 200 --duration 5m
```

Useful flags:

| Flag | Default | Purpose |
|---|---|---|
| `--rate` | `200` | events/sec |
| `--duration` | `5m` | how long to drive; `0` = until SIGINT |
| `--payload-bytes` | `512` | size of each event body |
| `--db` | `postgres://chainlink:chainlink@localhost:5432/chainlink_test?sslmode=disable` | Postgres DSN |
| `--otlp` | `localhost:4317` | OTLP gRPC endpoint |
| `--service` | `durable-emitter-driver` | OTel `service.name`; the dashboard filters on `exported_job="durable-emitter-driver"` so change both or neither |
| `--export-interval` | `5s` | metric export cadence; must be ≤ the Prom scrape interval (5s) |

### 4. View metrics

- **Grafana:** <http://localhost:3000> → Dashboards → "Durable Emitter".
- **Prometheus:** <http://localhost:9090>. Try `durable_emitter_queue_depth`,
  `rate(durable_emitter_emit_success_total[1m])`, etc.
- **Collector raw output:** `docker compose logs -f otel-collector` (the
  `debug` exporter is enabled at `verbosity: basic`).

### 5. Tear down

```bash
docker compose down -v
```

`-v` wipes the Postgres volume so the next run starts from an empty queue.

## Validating the new instruments

The metrics added in this PR/branch:

- `durable_emitter_batch_enqueue_buffer_full_total{phase}`
- `durable_emitter_insert_coalescer_queue_fill_ratio`
- `durable_emitter_fallback_in_flight`

Process CPU / heap panels need `pollProcessGauges` in chainlink-common (and
the extra batch/coalescer instruments below). Until that version is pinned in
`core/scripts/go.mod`, add a temporary replace directive:

```
replace github.com/smartcontractkit/chainlink-common => /path/to/chainlink-common
```

Then `go mod tidy`, rebuild the driver, and regenerate the Grafana JSON from
chainlink-observability. Revert the replace before opening a PR.

Additional instruments on recent chainlink-common:


`fallback_in_flight` will stay at 0 against `NoopClient` (the noop never
returns a publish error, so the fallback goroutine is never spawned). To force
fallback activity, swap in a fake client that returns errors — easy follow-up
if needed.

## Notes / caveats

- The dashboard expects an `exported_job` label, produced by the collector's
  `resource_to_telemetry_conversion: true` setting on the Prometheus exporter.
  If you change the OTel exporter config, panels may go blank until that label
  is preserved.
- Production/staging dashboards are deployed from chainlink-observability via
  `go run ./chainlink-observability deploy durable-emitter ...`. This harness
  only generates a local JSON snapshot for docker-compose Grafana.
- The driver uses `chipingress.NoopClient` so all publish RPCs "succeed"
  instantly. Queue depth will hover near zero unless you crank the rate well
  past the batch settings. To see backpressure metrics light up, drop
  `--rate` to something tiny and inspect the polling gauges, or wire in a
  delayed/error-returning fake client.
- This harness intentionally bypasses `durableemitter.Setup()` (which dials a
  real Chip Ingress endpoint) and uses `NewDurableEmitter` directly. The
  metrics surface is identical to the production wiring.
