// Command durable-emitter-driver exercises the durable emitter against a real
// Postgres store and pushes OTel metrics to a local collector. Use it together
// with the docker-compose stack one directory up to validate the metrics
// pipeline end-to-end (driver → collector → Prometheus → Grafana dashboard).
//
// The chip-ingress transport is stubbed with chipingress.NoopClient so no real
// Chip Ingress endpoint is needed. Events are inserted into Postgres, picked
// up by the batch emitter (which immediately "succeeds" the publish via the
// noop), marked delivered, and purged. The retransmit and expiry loops also
// run so the polling-based gauges (queue depth, oldest age, capacity ratio)
// have data to report.
package main

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"

	"github.com/smartcontractkit/chainlink-common/pkg/chipingress"
	chipingressbatch "github.com/smartcontractkit/chainlink-common/pkg/chipingress/batch"
	"github.com/smartcontractkit/chainlink-common/pkg/durableemitter"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

func main() {
	var (
		dbURL          = flag.String("db", "postgres://chainlink:chainlink@localhost:5432/chainlink_test?sslmode=disable", "Postgres DSN")
		otlpEndpoint   = flag.String("otlp", "localhost:4317", "OTLP gRPC endpoint for metrics")
		serviceName    = flag.String("service", "durable-emitter-driver", "OTel resource service.name (matches the exported_job filter on the Grafana dashboard)")
		rate           = flag.Float64("rate", 200, "Events per second to emit")
		duration       = flag.Duration("duration", 5*time.Minute, "How long to drive events (0 = until SIGINT)")
		payloadSize    = flag.Int("payload-bytes", 512, "Size of each event body")
		exportInterval = flag.Duration("export-interval", 5*time.Second, "OTel metric export interval (must be >= the Grafana scrape interval to avoid gaps)")
	)
	flag.Parse()

	lggr, err := logger.New()
	if err != nil {
		fatalf("logger: %v", err)
	}
	slggr := logger.Sugared(lggr)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// --- OTel meter exporting to the local collector ----------------------
	exporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(*otlpEndpoint),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		fatalf("otlp exporter: %v", err)
	}
	res, err := resource.New(ctx,
		resource.WithAttributes(semconv.ServiceName(*serviceName)),
	)
	if err != nil {
		fatalf("otel resource: %v", err)
	}
	mp := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(exporter,
			metric.WithInterval(*exportInterval),
		)),
	)
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = mp.Shutdown(shutdownCtx)
	}()
	meter := mp.Meter("durable-emitter-driver")

	// --- Postgres store ---------------------------------------------------
	db, err := sqlx.ConnectContext(ctx, "postgres", *dbURL)
	if err != nil {
		fatalf("postgres connect: %v", err)
	}
	defer db.Close()
	pgStore := durableemitter.NewPgDurableEventStore(db)

	// --- Noop chip-ingress transport (no network needed) ------------------
	noop := chipingress.NoopClient{}
	batchClient, err := chipingressbatch.NewBatchClient(noop,
		chipingressbatch.WithBatchSize(50),
		chipingressbatch.WithBatchInterval(50*time.Millisecond),
		chipingressbatch.WithMaxConcurrentSends(4),
	)
	if err != nil {
		fatalf("batch client: %v", err)
	}

	// --- Emitter with metrics enabled (same config as chainlink wiring) ---
	cfg := durableemitter.DefaultConfig()
	cfg.Metrics = &durableemitter.DurableEmitterMetricsConfig{
		MaxQueuePayloadBytes: 512 << 20, // 512 MiB; matches application.go default
	}
	emitter, err := durableemitter.NewDurableEmitter(
		pgStore,
		batchClient,
		noop, // fallback client
		true, // retransmitEnabled
		cfg,
		lggr,
		meter,
	)
	if err != nil {
		fatalf("new emitter: %v", err)
	}
	if err := emitter.Start(ctx); err != nil {
		fatalf("emitter start: %v", err)
	}
	defer func() {
		if err := emitter.Close(); err != nil {
			slggr.Warnw("emitter close", "err", err)
		}
	}()

	// --- Event driver -----------------------------------------------------
	payload := make([]byte, *payloadSize)
	if _, err := rand.Read(payload); err != nil {
		fatalf("random payload: %v", err)
	}

	interval := time.Duration(float64(time.Second) / *rate)
	if interval <= 0 {
		fatalf("rate must be > 0")
	}
	deadline := time.Now().Add(*duration)
	if *duration == 0 {
		deadline = time.Time{} // run until SIGINT
	}

	slggr.Infow("driver starting",
		"rate", *rate,
		"duration", duration.String(),
		"otlp", *otlpEndpoint,
		"service", *serviceName,
		"payload_bytes", *payloadSize,
	)

	var emitted, errored atomic.Int64
	go reportProgress(ctx, slggr, &emitted, &errored, 5*time.Second)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			slggr.Infow("driver stopping (signal)", "emitted", emitted.Load(), "errored", errored.Load())
			return
		case <-ticker.C:
			if !deadline.IsZero() && time.Now().After(deadline) {
				slggr.Infow("driver stopping (duration reached)", "emitted", emitted.Load(), "errored", errored.Load())
				return
			}
			emitCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
			err := emitter.Emit(emitCtx, payload,
				"source", "cre.local",
				"type", "durable_emitter_driver_event",
			)
			cancel()
			if err != nil {
				errored.Add(1)
			} else {
				emitted.Add(1)
			}
		}
	}
}

func reportProgress(ctx context.Context, lggr logger.SugaredLogger, emitted, errored *atomic.Int64, every time.Duration) {
	t := time.NewTicker(every)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			lggr.Infow("driver progress",
				"emitted_total", emitted.Load(),
				"errored_total", errored.Load(),
			)
		}
	}
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "durable-emitter-driver: "+format+"\n", args...)
	os.Exit(1)
}
