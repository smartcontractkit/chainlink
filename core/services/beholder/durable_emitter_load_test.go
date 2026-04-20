package beholder_test

// External Chip Ingress (integration):
//
//	Set CHIP_INGRESS_TEST_ADDR=host:port to dial a real Chip Ingress instead of the in-process mock.
//	Optional:
//	  CHIP_INGRESS_TEST_TLS=1|true          — use TLS (default: insecure plaintext gRPC)
//	  CHIP_INGRESS_TEST_BASIC_AUTH_USER     — basic auth user (e.g. admin)
//	  CHIP_INGRESS_TEST_BASIC_AUTH_PASS     — basic auth password
//
//	Tests that inject Chip failures or count in-process receives (outage, slow-Chip) are skipped
//	when CHIP_INGRESS_TEST_ADDR is set.
//
//	Running a real server: see atlas/chip-ingress/README.md. You need Kafka/Redpanda, the
//	`chip-demo` topic, and schema subject `chip-demo-pb.DemoClientPayload` (run
//	`make create-topic-and-schema` from atlas/chip-ingress, or equivalent rpk commands).
//	CRE local Beholder (`go run . env beholder start` / `env start --with-beholder`) creates
//	`chip-demo` and registers this schema automatically; see core/scripts/cre/environment/configs/chip-ingress.toml.
//	Tests call RegisterSchemas with the bundled proto; Chip still needs the topic to exist for Kafka.
//	External mode uses the Atlas demo shape: chip-demo / pb.DemoClientPayload + protobuf payload.
//	If unset, CHIP_INGRESS_TEST_BASIC_AUTH_USER/PASS default to chip-ingress-demo-client / password
//	(atlas docker-compose demo account). Set CHIP_INGRESS_TEST_SKIP_BASIC_AUTH=1 to omit auth.
//	Set CHIP_INGRESS_TEST_SKIP_SCHEMA_REGISTRATION=1 to skip RegisterSchemas (schema pre-created only).

import (
	"context"
	"fmt"
	"net"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	cepb "github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/chipingress"
	"github.com/smartcontractkit/chainlink-common/pkg/chipingress/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	beholdersvc "github.com/smartcontractkit/chainlink/v2/core/services/beholder"

	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

// chipLoadTestDemoProto is the raw .proto registered with Chip for subject chip-demo-pb.DemoClientPayload
// (keep in sync with chip_load_test_demo.proto).
const chipLoadTestDemoProto = `syntax = "proto3";

option go_package = "github.com/smartcontractkit/chainlink/v2/core/services/beholder;beholder";

package pb;

message DemoClientPayload {
  string id = 1;
  string domain = 2;
  string entity = 3;
  int64 batch_num = 4;
  int64 message_num = 5;
  int64 batch_position = 6;
}
`

// sustainedThroughputMockPublishLatency is the in-process mock's server-side sleep per Publish
// RPC in TestFullStack_SustainedThroughput (only). External Chip ignores this.
const sustainedThroughputMockPublishLatency = 500 * time.Millisecond

// loadTestServer is a controllable gRPC ChipIngress server for load tests.
type loadTestServer struct {
	pb.UnimplementedChipIngressServer

	mu         sync.Mutex
	publishErr error
	batchErr   error
	// publishDelayNs is nanoseconds to sleep in Publish (0 = none). Atomic so handlers see
	// the value set before traffic without a data race on the hot path.
	publishDelayNs atomic.Int64

	publishCount atomic.Int64
	batchCount   atomic.Int64
	totalEvents  atomic.Int64
}

func (s *loadTestServer) Publish(_ context.Context, _ *cepb.CloudEvent) (*pb.PublishResponse, error) {
	if ns := s.publishDelayNs.Load(); ns > 0 {
		time.Sleep(time.Duration(ns))
	}
	s.publishCount.Add(1)
	s.totalEvents.Add(1)
	s.mu.Lock()
	defer s.mu.Unlock()
	return &pb.PublishResponse{}, s.publishErr
}

func (s *loadTestServer) PublishBatch(_ context.Context, in *pb.CloudEventBatch) (*pb.PublishResponse, error) {
	if ns := s.publishDelayNs.Load(); ns > 0 {
		time.Sleep(time.Duration(ns))
	}
	s.batchCount.Add(1)
	s.totalEvents.Add(int64(len(in.Events)))
	s.mu.Lock()
	defer s.mu.Unlock()
	return &pb.PublishResponse{}, s.batchErr
}

func (s *loadTestServer) Ping(context.Context, *pb.EmptyRequest) (*pb.PingResponse, error) {
	return &pb.PingResponse{Message: "pong"}, nil
}

func (s *loadTestServer) setPublishErr(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.publishErr = err
}

func (s *loadTestServer) setBatchErr(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.batchErr = err
}

func (s *loadTestServer) setPublishDelay(d time.Duration) {
	if d <= 0 {
		s.publishDelayNs.Store(0)
		return
	}
	s.publishDelayNs.Store(d.Nanoseconds())
}

func startLoadServer(t testing.TB) (*loadTestServer, string) {
	t.Helper()
	srv := &loadTestServer{}
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	gs := grpc.NewServer()
	pb.RegisterChipIngressServer(gs, srv)
	go func() { _ = gs.Serve(lis) }()
	t.Cleanup(func() { gs.GracefulStop() })

	return srv, lis.Addr().String()
}

func chipClient(t testing.TB, addr string) chipingress.Client {
	t.Helper()
	c, err := chipingress.NewClient(addr, chipingress.WithInsecureConnection())
	require.NoError(t, err)
	t.Cleanup(func() { _ = c.Close() })
	return c
}

const (
	envChipIngressTestAddr      = "CHIP_INGRESS_TEST_ADDR"
	envChipIngressTestTLS       = "CHIP_INGRESS_TEST_TLS"
	envChipIngressTestBasicUser = "CHIP_INGRESS_TEST_BASIC_AUTH_USER"
	envChipIngressTestBasicPass = "CHIP_INGRESS_TEST_BASIC_AUTH_PASS"
)

func externalChipConfigured() bool {
	return strings.TrimSpace(os.Getenv(envChipIngressTestAddr)) != ""
}

func newChipClientFromEnv(t testing.TB) chipingress.Client {
	t.Helper()
	addr := strings.TrimSpace(os.Getenv(envChipIngressTestAddr))
	require.NotEmpty(t, addr, envChipIngressTestAddr)

	var opts []chipingress.Opt
	if tlsEnv := os.Getenv(envChipIngressTestTLS); tlsEnv == "1" || strings.EqualFold(tlsEnv, "true") {
		opts = append(opts, chipingress.WithTLS())
	} else {
		opts = append(opts, chipingress.WithInsecureConnection())
	}
	user := os.Getenv(envChipIngressTestBasicUser)
	pass := os.Getenv(envChipIngressTestBasicPass)
	skipAuth := os.Getenv("CHIP_INGRESS_TEST_SKIP_BASIC_AUTH")
	if skipAuth != "1" && !strings.EqualFold(skipAuth, "true") {
		if user == "" && pass == "" {
			// Default matches atlas/chip-ingress docker-compose CE_SA_CHIP_INGRESS_DEMO_CLIENT.
			user = "chip-ingress-demo-client"
			pass = "password"
		}
	}
	if user != "" && pass != "" {
		opts = append(opts, chipingress.WithBasicAuth(user, pass))
	}

	c, err := chipingress.NewClient(addr, opts...)
	require.NoError(t, err)
	t.Cleanup(func() { _ = c.Close() })
	return c
}

// startChipIngressOrMock starts the in-process mock ChipIngress server unless
// CHIP_INGRESS_TEST_ADDR is set; then it returns mock=nil and a client to the external server.
func startChipIngressOrMock(t testing.TB) (mock *loadTestServer, client chipingress.Client) {
	t.Helper()
	if externalChipConfigured() {
		t.Logf("Using external Chip Ingress at %s (%s)", os.Getenv(envChipIngressTestAddr), envChipIngressTestAddr)
		c := newChipClientFromEnv(t)
		registerChipDemoSchema(t, c)
		return nil, c
	}
	mock, addr := startLoadServer(t)
	return mock, chipClient(t, addr)
}

// registerChipDemoSchema registers the demo protobuf with Chip (via chip-config) so GetSchema
// succeeds for subject chip-demo-pb.DemoClientPayload. Skip with CHIP_INGRESS_TEST_SKIP_SCHEMA_REGISTRATION=1.
func registerChipDemoSchema(t testing.TB, client chipingress.Client) {
	t.Helper()
	if os.Getenv("CHIP_INGRESS_TEST_SKIP_SCHEMA_REGISTRATION") == "1" ||
		strings.EqualFold(os.Getenv("CHIP_INGRESS_TEST_SKIP_SCHEMA_REGISTRATION"), "true") {
		t.Logf("skipping RegisterSchemas (%s)", "CHIP_INGRESS_TEST_SKIP_SCHEMA_REGISTRATION")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, err := client.RegisterSchemas(ctx, &pb.Schema{
		Subject: "chip-demo-pb.DemoClientPayload",
		Schema:  chipLoadTestDemoProto,
		Format:  pb.SchemaType_PROTOBUF,
	})
	if err != nil {
		// Common when schema was already registered (e.g. make create-schema).
		msg := strings.ToLower(err.Error())
		if strings.Contains(msg, "already") || strings.Contains(msg, "exists") || strings.Contains(msg, "duplicate") {
			t.Logf("RegisterSchemas: treating as OK (schema likely present): %v", err)
			return
		}
		require.NoError(t, err, "RegisterSchemas for chip-demo; try atlas/chip-ingress make create-topic-and-schema")
	}
}

func skipIfExternalChip(t *testing.T, reason string) {
	t.Helper()
	if externalChipConfigured() {
		t.Skipf("requires in-process mock Chip: %s (unset %s)", reason, envChipIngressTestAddr)
	}
}

func formatMockServerEvents(srv *loadTestServer) string {
	if srv == nil {
		return "N/A"
	}
	return strconv.FormatInt(srv.totalEvents.Load(), 10)
}

func loadEmitAttrs() []any {
	if externalChipConfigured() {
		// Wire-compatible with atlas chip-ingress demo (see chip_load_test_demo.proto).
		return []any{
			"source", "chip-demo",
			"type", "pb.DemoClientPayload",
			"datacontenttype", "application/protobuf",
			"dataschema", "https://example.com/demo-client-schema",
			"time", time.Now(),
		}
	}
	return []any{"source", "cre.billing", "type", "workflow_execution_finished"}
}

// buildLoadTestPayload returns raw bytes for Emit(). For the in-process mock, arbitrary bytes are
// fine. For real Chip Ingress, payload must protobuf-decode as pb.DemoClientPayload (subject
// chip-demo-pb.DemoClientPayload in schema registry).
func buildLoadTestPayload(targetSize int) []byte {
	if !externalChipConfigured() {
		if targetSize < 0 {
			targetSize = 0
		}
		b := make([]byte, targetSize)
		return b
	}
	if targetSize <= 0 {
		targetSize = 1
	}
	p := &beholdersvc.DemoClientPayload{
		Domain:        "chip-demo",
		Entity:        "pb.DemoClientPayload",
		BatchNum:      1,
		MessageNum:    1,
		BatchPosition: 0,
	}
	id := ""
	for range targetSize*4 + 512 {
		p.Id = id
		b, err := proto.Marshal(p)
		if err != nil {
			return []byte{0x0a, 0x00}
		}
		if len(b) >= targetSize {
			for len(id) > 0 && len(b) > targetSize {
				id = id[:len(id)-1]
				p.Id = id
				b, _ = proto.Marshal(p)
			}
			return b
		}
		id += "x"
	}
	b, _ := proto.Marshal(p)
	return b
}

// TestChipIngressExternalPing is a smoke test: verifies gRPC connectivity when CHIP_INGRESS_TEST_ADDR is set.
func TestChipIngressExternalPing(t *testing.T) {
	if !externalChipConfigured() {
		t.Skipf("set %s to dial a real Chip Ingress (e.g. 127.0.0.1:50051)", envChipIngressTestAddr)
	}
	client := newChipClientFromEnv(t)
	ctx := testutils.Context(t)
	_, err := client.Ping(ctx, &pb.EmptyRequest{})
	require.NoError(t, err)
	t.Logf("Ping OK to %s", os.Getenv(envChipIngressTestAddr))
}

// ---------- Full-stack load tests: DurableEmitter + Postgres + gRPC ----------

// TestFullStack_SustainedThroughput measures steady-state throughput with
// real Postgres persistence and gRPC delivery. This answers: "how many
// events/sec can we sustain end-to-end?"
//
// Wall time is dominated by totalEvents / sustained Emit (insert) rate — often
// many minutes at 100k events. Run with -short for a 10k-event run (~tens of s).
// Spurious retransmits happen if RetransmitAfter is shorter than tail
// MarkDelivered latency under load; we use a generous RetransmitAfter here.
// With the in-process mock, each Publish RPC sleeps sustainedThroughputMockPublishLatency
// (const); pipeline logs should show ~that much in immediate Publish p50/p99/mean.
func TestFullStack_SustainedThroughput(t *testing.T) {
	// Must use non-txdb Postgres: txdb is a single transaction; any SQL error
	// aborts it and all follow-up queries fail with SQLSTATE 25P02 under concurrent
	// purge/retransmit/mark-delivered (DurableEmitter background loops).
	db := directDB(t)
	srv, client := startChipIngressOrMock(t)
	if srv != nil {
		srv.setPublishDelay(sustainedThroughputMockPublishLatency)
		t.Logf("Sustained throughput: mock Chip Publish server delay = %s (const sustainedThroughputMockPublishLatency)",
			sustainedThroughputMockPublishLatency)
	}
	store := beholdersvc.NewPgDurableEventStore(db)

	ctx := testutils.Context(t)

	pipe := &pipelineDeliveryStats{}
	cfg := beholder.DefaultDurableEmitterConfig()
	cfg.QuietMode = true
	cfg.RetransmitInterval = 30 * time.Second
	cfg.RetransmitAfter = 5 * time.Minute
	cfg.RetransmitBatchSize = 50000 // Shouldn't matter since all go thorugh first time
	cfg.PublishTimeout = 5 * time.Second
	cfg.PublishBatchSize = 25_000
	cfg.PublishBatchWorkers = 3
	cfg.PublishBatchFlushInterval = 2 * time.Millisecond
	cfg.PublishBatchChannelSize = 2_000_000
	cfg.InsertBatchSize = 2000
	cfg.InsertBatchFlushInterval = 250 * time.Microsecond
	cfg.InsertBatchWorkers = 20
	cfg.DisablePruning = true
	cfg.Hooks = newPipelineHooks(pipe)

	totalEvents := 5_000_000
	//if testing.Short() {
	//totalEvents = 10_000
	//}
	const concurrency = 200

	// Target produce rate in msg/s. 0 = unlimited (fire-hose / max throughput).
	targetRate := 0 //5000 //0

	// Wire OTel metrics to the local obs stack when OTEL_EXPORTER_OTLP_ENDPOINT is set.
	// Start the obs stack first: ./bin/ctf obs up
	// Then run: OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317 CHIP_INGRESS_TEST_ADDR=... go test ...
	if otlpEndpoint := strings.TrimSpace(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")); otlpEndpoint != "" {
		otlpEndpoint = strings.TrimPrefix(otlpEndpoint, "http://")
		exp, otelErr := otlpmetricgrpc.New(ctx,
			otlpmetricgrpc.WithEndpoint(otlpEndpoint),
			otlpmetricgrpc.WithInsecure(),
		)
		require.NoError(t, otelErr, "otlp metric exporter")
		res := sdkresource.NewWithAttributes("",
			attribute.String("service.name", "durable-emitter-loadtest"),
		)
		mp := sdkmetric.NewMeterProvider(
			sdkmetric.WithResource(res),
			sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exp,
				sdkmetric.WithInterval(1*time.Second),
			)),
		)
		t.Cleanup(func() { _ = mp.Shutdown(context.Background()) })
		bc := beholder.NewNoopClient()
		bc.MeterProvider = mp
		bc.Meter = mp.Meter("beholder")
		beholder.SetClient(bc)
		t.Cleanup(func() { beholder.SetClient(beholder.NewNoopClient()) })
		cfg.Metrics = &beholder.DurableEmitterMetricsConfig{
			PollInterval:       1 * time.Second,
			RecordProcessStats: true,
		}
		t.Logf("OTel metrics enabled → %s (1s push interval, Grafana: http://localhost:3000)", otlpEndpoint)
	}

	em, err := beholder.NewDurableEmitter(store, client, cfg, logger.Test(t))
	require.NoError(t, err)

	em.Start(ctx)
	defer em.Close()

	t.Logf("Full-stack sustained throughput: totalEvents=%d, concurrency=%d, targetRate=%d msg/s",
		totalEvents, concurrency, targetRate)

	payload := buildLoadTestPayload(256) // ~256 byte record (protobuf for external Chip)

	// CPU snapshot at start (getrusage; cumulative from process start, so we diff).
	var cpuStart syscall.Rusage
	_ = syscall.Getrusage(syscall.RUSAGE_SELF, &cpuStart)

	start := time.Now()

	var wg sync.WaitGroup
	var emitErrors atomic.Int64
	var producerWallNs atomic.Int64 // cumulative wall time inside Emit() calls

	// Per-worker rate: divide target evenly across goroutines.
	perWorkerRate := targetRate / concurrency
	if perWorkerRate < 1 && targetRate > 0 {
		perWorkerRate = 1
	}

	for w := 0; w < concurrency; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			eventsPerWorker := totalEvents / concurrency

			if targetRate <= 0 {
				for i := 0; i < eventsPerWorker; i++ {
					t0 := time.Now()
					if err := em.Emit(ctx, payload, loadEmitAttrs()...); err != nil {
						emitErrors.Add(1)
					}
					producerWallNs.Add(int64(time.Since(t0)))
				}
				return
			}

			interval := time.Duration(float64(time.Second) / float64(perWorkerRate))
			ticker := time.NewTicker(interval)
			defer ticker.Stop()
			for i := 0; i < eventsPerWorker; i++ {
				<-ticker.C
				t0 := time.Now()
				if err := em.Emit(ctx, payload, loadEmitAttrs()...); err != nil {
					emitErrors.Add(1)
				}
				producerWallNs.Add(int64(time.Since(t0)))
			}
		}()
	}
	wg.Wait()
	emitElapsed := time.Since(start)

	assert.Equal(t, int64(0), emitErrors.Load(), "all emits should succeed")

	// Wait for all events to be delivered and store to drain (pending list empty;
	// Postgres may still have tombstones until purge loop catches up).
	drainWait := 60 * time.Second
	if totalEvents >= 100_000 {
		drainWait = 300 * time.Second
	}
	require.Eventually(t, func() bool {
		pending, _ := store.ListPending(ctx, time.Now().Add(time.Hour), 1)
		return len(pending) == 0
	}, drainWait, 100*time.Millisecond, "store should drain completely (no pending delivery)")

	totalElapsed := time.Since(start)

	// CPU diff: user + system seconds consumed over the whole test.
	var cpuEnd syscall.Rusage
	_ = syscall.Getrusage(syscall.RUSAGE_SELF, &cpuEnd)
	cpuUserSec := (float64(cpuEnd.Utime.Sec) + float64(cpuEnd.Utime.Usec)/1e6) -
		(float64(cpuStart.Utime.Sec) + float64(cpuStart.Utime.Usec)/1e6)
	cpuSysSec := (float64(cpuEnd.Stime.Sec) + float64(cpuEnd.Stime.Usec)/1e6) -
		(float64(cpuStart.Stime.Sec) + float64(cpuStart.Stime.Usec)/1e6)
	cpuTotalSec := cpuUserSec + cpuSysSec
	// Utilization: fraction of available CPU time (wall × GOMAXPROCS).
	cpuUtilPct := 100.0 * cpuTotalSec / (totalElapsed.Seconds() * float64(runtime.GOMAXPROCS(0)))

	insN := pipe.emitIns.count()
	insP50 := durMs(pipe.emitIns.percentile(0.50))
	insP99 := durMs(pipe.emitIns.percentile(0.99))
	insMean := durMs(pipe.emitIns.mean())

	// Batch publish loop stats (primary delivery path when PublishBatchSize > 0).
	bpN := pipe.batchLoopPub.count()
	bpP50 := durMs(pipe.batchLoopPub.percentile(0.50))
	bpP99 := durMs(pipe.batchLoopPub.percentile(0.99))
	bpMean := durMs(pipe.batchLoopPub.mean())
	bmN := pipe.batchLoopDel.count()
	bmP50 := durMs(pipe.batchLoopDel.percentile(0.50))
	bmP99 := durMs(pipe.batchLoopDel.percentile(0.99))
	bmMean := durMs(pipe.batchLoopDel.mean())
	bpEvents := pipe.batchLoopPubEvents.Load()
	bmEvents := pipe.batchLoopMarkEvents.Load()

	// Legacy per-event stats (only populated when PublishBatchSize == 0).
	pubN := pipe.immPub.count()
	pubErrs := pipe.immPubErr.Load()
	pubP50 := durMs(pipe.immPub.percentile(0.50))
	pubP99 := durMs(pipe.immPub.percentile(0.99))
	pubMean := durMs(pipe.immPub.mean())
	delN := pipe.immDel.count()
	delP50 := durMs(pipe.immDel.percentile(0.50))
	delP99 := durMs(pipe.immDel.percentile(0.99))
	delMean := durMs(pipe.immDel.mean())

	// DB-based end-to-end latency (delivered_at - created_at).
	dbE2E, dbErr := queryDBE2ELatency(ctx, db)
	if dbErr != nil {
		t.Logf("WARNING: failed to query DB e2e latency: %v", dbErr)
	}

	var serverLine string
	if srv != nil {
		serverLine = fmt.Sprintf("%d (mock; batches: %d, individual: %d)",
			srv.totalEvents.Load(), srv.batchCount.Load(), srv.publishCount.Load())
	} else {
		serverLine = "N/A  (use external Chip metrics)"
	}

	target := chipIngressTargetDescription(srv)
	batchMode := cfg.PublishBatchSize > 0
	batchLabel := "disabled (per-event goroutines)"
	if batchMode {
		workers := cfg.PublishBatchWorkers
		if workers <= 0 {
			workers = 1
		}
		batchLabel = fmt.Sprintf("%d events/batch, %d workers", cfg.PublishBatchSize, workers)
	}

	rateLabel := "unlimited (fire-hose)"
	if targetRate > 0 {
		rateLabel = fmt.Sprintf("%d msg/s", targetRate)
	}

	t.Logf("╔══════════════════════════════════════════════════════════════════════╗")
	t.Logf("║              SUSTAINED THROUGHPUT TEST RESULTS                      ║")
	t.Logf("╠══════════════════════════════════════════════════════════════════════╣")
	t.Logf("║ Target:      %-55s ║", target)
	t.Logf("║ Batch mode:  %-55s ║", batchLabel)
	t.Logf("║ Target rate: %-55s ║", rateLabel)
	t.Logf("╠══════════════════════════════════════════════════════════════════════╣")
	t.Logf("║ EMIT (DB insert, batched gRPC)                                     ║")
	t.Logf("║   Events:         %-49d ║", totalEvents)
	t.Logf("║   Errors:         %-49d ║", emitErrors.Load())
	t.Logf("║   Elapsed:        %-49s ║", emitElapsed.Round(time.Millisecond))
	t.Logf("║   Actual rate:    %-49s ║", fmt.Sprintf("%.0f msg/s", float64(totalEvents)/emitElapsed.Seconds()))
	t.Logf("╠══════════════════════════════════════════════════════════════════════╣")
	t.Logf("║ DELIVERY (PublishBatch → MarkDeliveredBatch)                       ║")
	t.Logf("║   Server received:%-49s ║", serverLine)
	t.Logf("║   Total elapsed:  %-49s ║", totalElapsed.Round(time.Millisecond))
	t.Logf("║   End-to-end rate:%-49s ║", fmt.Sprintf("%.0f events/sec", float64(totalEvents)/totalElapsed.Seconds()))
	t.Logf("╠══════════════════════════════════════════════════════════════════════╣")
	t.Logf("║ PENDING QUEUE DEPTH (exact atomic counter)                         ║")
	t.Logf("║   Max:            %-49d ║", em.PendingMax())
	t.Logf("║   Current:        %-49d ║", em.PendingDepth())
	t.Logf("╠══════════════════════════════════════════════════════════════════════╣")
	t.Logf("║ PIPELINE LATENCY (hooks)                                           ║")
	t.Logf("║                        %-10s %-10s %-10s %-10s ║", "n", "p50 (ms)", "p99 (ms)", "mean (ms)")
	t.Logf("║   Emit (INSERT):       %-10d %-10.2f %-10.2f %-10.2f ║", insN, insP50, insP99, insMean)
	if batchMode {
		avgBatchSize := float64(0)
		if bpN > 0 {
			avgBatchSize = float64(bpEvents) / float64(bpN)
		}
		t.Logf("║   PublishBatch (gRPC): %-10d %-10.2f %-10.2f %-10.2f ║", bpN, bpP50, bpP99, bpMean)
		t.Logf("║     └ events published:%-49d ║", bpEvents)
		t.Logf("║     └ avg batch size:  %-49s ║", fmt.Sprintf("%.1f events/batch (configured: %d)", avgBatchSize, cfg.PublishBatchSize))
		t.Logf("║   MarkDeliveredBatch:  %-10d %-10.2f %-10.2f %-10.2f ║", bmN, bmP50, bmP99, bmMean)
		t.Logf("║     └ events marked:   %-49d ║", bmEvents)
	} else {
		t.Logf("║   Publish (gRPC):      %-10d %-10.2f %-10.2f %-10.2f ║", pubN, pubP50, pubP99, pubMean)
		t.Logf("║   MarkDelivered:       %-10d %-10.2f %-10.2f %-10.2f ║", delN, delP50, delP99, delMean)
	}
	if pipe.batchPub.count() > 0 {
		t.Logf("║   Retransmit:          %-10d %-10.2f %-10s %-10.2f ║",
			pipe.batchPub.count(),
			durMs(pipe.batchPub.percentile(0.50)), "—",
			durMs(pipe.batchPub.mean()))
	}
	t.Logf("╠══════════════════════════════════════════════════════════════════════╣")
	t.Logf("║ END-TO-END LATENCY (DB: delivered_at − created_at)                 ║")
	if dbE2E.count > 0 {
		t.Logf("║   Events:              %-49d ║", dbE2E.count)
		t.Logf("║   p50:                 %-49s ║", fmt.Sprintf("%.2f ms", float64(dbE2E.p50.Microseconds())/1000.0))
		t.Logf("║   p99:                 %-49s ║", fmt.Sprintf("%.2f ms", float64(dbE2E.p99.Microseconds())/1000.0))
		t.Logf("║   mean:                %-49s ║", fmt.Sprintf("%.2f ms", float64(dbE2E.mean.Microseconds())/1000.0))
	} else {
		t.Logf("║   (no data — check DisablePruning=true)                          ║")
	}
	t.Logf("╠══════════════════════════════════════════════════════════════════════╣")
	t.Logf("║ PROCESS CPU (getrusage, GOMAXPROCS=%d)                             ║", runtime.GOMAXPROCS(0))
	t.Logf("║   User:           %-49s ║", fmt.Sprintf("%.2f s", cpuUserSec))
	t.Logf("║   System:         %-49s ║", fmt.Sprintf("%.2f s", cpuSysSec))
	t.Logf("║   Total:          %-49s ║", fmt.Sprintf("%.2f s", cpuTotalSec))
	t.Logf("║   Utilization:    %-49s ║", fmt.Sprintf("%.1f%% of %d cores × %.1fs wall", cpuUtilPct, runtime.GOMAXPROCS(0), totalElapsed.Seconds()))
	t.Logf("╠══════════════════════════════════════════════════════════════════════╣")
	availCPUSec := totalElapsed.Seconds() * float64(runtime.GOMAXPROCS(0))
	producerWallSec := float64(producerWallNs.Load()) / 1e9
	consumerPubWallSec := pipe.batchLoopPub.sum().Seconds()
	consumerMarkWallSec := pipe.batchLoopDel.sum().Seconds()
	consumerWallSec := consumerPubWallSec + consumerMarkWallSec
	producerPct := 100.0 * producerWallSec / availCPUSec
	consumerPct := 100.0 * consumerWallSec / availCPUSec
	overheadPct := cpuUtilPct - producerPct - consumerPct
	if overheadPct < 0 {
		overheadPct = 0
	}
	t.Logf("║ CPU BREAKDOWN (wall-time attribution)                              ║")
	t.Logf("║   Producer (Emit):     %-44s ║", fmt.Sprintf("%.2f s  (%.1f%%)", producerWallSec, producerPct))
	t.Logf("║     └ DB Insert:       %-44s ║", fmt.Sprintf("%.2f s", pipe.emitIns.sum().Seconds()))
	t.Logf("║     └ Event build:     %-44s ║", fmt.Sprintf("%.2f s", producerWallSec-pipe.emitIns.sum().Seconds()))
	t.Logf("║   Consumer (Emitter):  %-44s ║", fmt.Sprintf("%.2f s  (%.1f%%)", consumerWallSec, consumerPct))
	t.Logf("║     └ PublishBatch:    %-44s ║", fmt.Sprintf("%.2f s", consumerPubWallSec))
	t.Logf("║     └ MarkDelivered:   %-44s ║", fmt.Sprintf("%.2f s", consumerMarkWallSec))
	t.Logf("║   Overhead (GC/sched): %-44s ║", fmt.Sprintf("%.1f%%", overheadPct))
	t.Logf("╠══════════════════════════════════════════════════════════════════════╣")
	t.Logf("║   Publish errors (need retransmit): %-31d ║", pubErrs+pipe.batchPubEventErrs.Load())
	t.Logf("╚══════════════════════════════════════════════════════════════════════╝")

	if srv != nil {
		assert.GreaterOrEqual(t, srv.totalEvents.Load(), int64(totalEvents),
			"server should have received all events (may have retransmit duplicates)")
	}
}

// TestFullStack_ChipOutage runs sustained emit load while injecting periodic
// Chip outages at fixed intervals. Each cycle: Chip is up for outagePeriod,
// then down for outageDuration, then recovers. The test measures how the DB
// queue accumulates during each outage and drains after each recovery, giving
// a real view of back-pressure, retransmit drain rate, and DB load over time.
//
// OTel metrics are exported when OTEL_EXPORTER_OTLP_ENDPOINT is set (same as
// TestFullStack_SustainedThroughput). Start the obs stack first:
//
//	./bin/ctf obs up
//	OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317 go test ./core/services/beholder/ -run TestFullStack_ChipOutage -v -count=1 -timeout 20m
func TestFullStack_ChipOutage(t *testing.T) {
	skipIfExternalChip(t, "inject Unavailable errors on mock server")

	db := directDB(t)
	srv, client := startChipIngressOrMock(t)
	require.NotNil(t, srv)
	store := beholdersvc.NewPgDurableEventStore(db)

	ctx := testutils.Context(t)

	pipe := &pipelineDeliveryStats{}
	cfg := beholder.DefaultDurableEmitterConfig()
	cfg.QuietMode = true
	cfg.RetransmitInterval = 200 * time.Millisecond
	cfg.RetransmitAfter = 500 * time.Millisecond
	cfg.RetransmitBatchSize = 200
	cfg.PublishTimeout = 2 * time.Second
	cfg.PublishBatchSize = 100
	cfg.DisablePruning = true
	cfg.Hooks = newPipelineHooks(pipe)

	// OTel metrics wiring (same as SustainedThroughput).
	if otlpEndpoint := strings.TrimSpace(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")); otlpEndpoint != "" {
		otlpEndpoint = strings.TrimPrefix(otlpEndpoint, "http://")
		exp, otelErr := otlpmetricgrpc.New(ctx,
			otlpmetricgrpc.WithEndpoint(otlpEndpoint),
			otlpmetricgrpc.WithInsecure(),
		)
		require.NoError(t, otelErr, "otlp metric exporter")
		res := sdkresource.NewWithAttributes("",
			attribute.String("service.name", "durable-emitter-loadtest"),
		)
		mp := sdkmetric.NewMeterProvider(
			sdkmetric.WithResource(res),
			sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exp, sdkmetric.WithInterval(1*time.Second))),
		)
		t.Cleanup(func() { _ = mp.Shutdown(context.Background()) })
		bc := beholder.NewNoopClient()
		bc.MeterProvider = mp
		bc.Meter = mp.Meter("beholder")
		beholder.SetClient(bc)
		t.Cleanup(func() { beholder.SetClient(beholder.NewNoopClient()) })
		cfg.Metrics = &beholder.DurableEmitterMetricsConfig{
			PollInterval:       1 * time.Second,
			RecordProcessStats: true,
		}
		t.Logf("OTel metrics enabled → %s (1s push interval, Grafana: http://localhost:3000)", otlpEndpoint)
	}

	em, err := beholder.NewDurableEmitter(store, client, cfg, logger.Test(t))
	require.NoError(t, err)
	em.Start(ctx)
	defer em.Close()

	// Outage schedule.
	const (
		outageCycles      = 3
		upDuration        = 20 * time.Second // Chip healthy between outages
		outageDuration    = 10 * time.Second // Chip unavailable per cycle
		emitConcurrency   = 5
		emitRatePerWorker = 200 // events/s target per worker (throttled)
	)
	totalEmitted := outageCycles * int(upDuration.Seconds()+outageDuration.Seconds()) *
		emitConcurrency * emitRatePerWorker
	// Cap to a sane ceiling so the test completes quickly.
	if totalEmitted > 50_000 {
		totalEmitted = 50_000
	}

	t.Logf("ChipOutage: %d cycles  up=%s  down=%s  workers=%d  target=%d events",
		outageCycles, upDuration, outageDuration, emitConcurrency, totalEmitted)

	// CPU snapshot.
	var cpuStart syscall.Rusage
	_ = syscall.Getrusage(syscall.RUSAGE_SELF, &cpuStart)

	// Outage injector: runs cycles of up/down in a background goroutine.
	type cycleResult struct {
		cycle         int
		outageStart   time.Time
		recoveryStart time.Time
		drainElapsed  time.Duration
		peakQueue     int64
		drainRate     float64 // events/sec
	}
	cycleResults := make([]cycleResult, 0, outageCycles)
	var cyclesMu sync.Mutex

	outageCtx, outageCancel := context.WithCancel(ctx)
	defer outageCancel()

	go func() {
		for cycle := 1; cycle <= outageCycles; cycle++ {
			// Wait for "up" phase.
			select {
			case <-outageCtx.Done():
				return
			case <-time.After(upDuration):
			}

			// Take down Chip.
			srv.setPublishErr(status.Error(codes.Unavailable, "chip down"))
			srv.setBatchErr(status.Error(codes.Unavailable, "chip down"))
			outStart := time.Now()
			t.Logf("↓ Cycle %d/%d: Chip DOWN at %s", cycle, outageCycles, outStart.Format("15:04:05"))

			// Wait for outage duration; the exact peak is tracked by em.PendingMax().
			select {
			case <-outageCtx.Done():
				return
			case <-time.After(outageDuration):
			}
			cyclePeak := em.PendingDepth()

			// Restore Chip.
			srv.setPublishErr(nil)
			srv.setBatchErr(nil)
			recovStart := time.Now()
			t.Logf("↑ Cycle %d/%d: Chip UP at %s (was down %s, peak queue %d rows)",
				cycle, outageCycles, recovStart.Format("15:04:05"),
				recovStart.Sub(outStart).Round(time.Millisecond), cyclePeak)

			// Wait for drain.
			drainDeadline := time.Now().Add(60 * time.Second)
			var drainElapsed time.Duration
			for time.Now().Before(drainDeadline) {
				pending, _ := store.ListPending(outageCtx, time.Now().Add(time.Hour), 1)
				if len(pending) == 0 {
					drainElapsed = time.Since(recovStart)
					break
				}
				time.Sleep(100 * time.Millisecond)
			}

			cyclesMu.Lock()
			cycleResults = append(cycleResults, cycleResult{
				cycle:         cycle,
				outageStart:   outStart,
				recoveryStart: recovStart,
				drainElapsed:  drainElapsed,
				peakQueue:     cyclePeak,
				drainRate:     float64(cyclePeak) / max(drainElapsed.Seconds(), 0.001),
			})
			cyclesMu.Unlock()
		}
	}()

	// Emit loop: concurrent workers emit at a throttled rate until totalEmitted.
	testStart := time.Now()
	payload := buildLoadTestPayload(256)
	var emitErrors atomic.Int64
	var emitCount atomic.Int64
	var emitWg sync.WaitGroup
	eventsPerWorker := totalEmitted / emitConcurrency

	for w := 0; w < emitConcurrency; w++ {
		emitWg.Add(1)
		go func() {
			defer emitWg.Done()
			interval := time.Duration(float64(time.Second) / float64(emitRatePerWorker))
			ticker := time.NewTicker(interval)
			defer ticker.Stop()
			for i := 0; i < eventsPerWorker; i++ {
				<-ticker.C
				if err := em.Emit(ctx, payload, loadEmitAttrs()...); err != nil {
					emitErrors.Add(1)
				}
				emitCount.Add(1)
			}
		}()
	}
	emitWg.Wait()
	outageCancel() // stop outage injector
	emitElapsed := time.Since(testStart)

	// Wait for final drain.
	require.Eventually(t, func() bool {
		pending, _ := store.ListPending(ctx, time.Now().Add(time.Hour), 1)
		return len(pending) == 0
	}, 60*time.Second, 100*time.Millisecond, "all events should drain after final recovery")

	totalElapsed := time.Since(testStart)

	// CPU diff.
	var cpuEnd syscall.Rusage
	_ = syscall.Getrusage(syscall.RUSAGE_SELF, &cpuEnd)
	cpuUserSec := (float64(cpuEnd.Utime.Sec) + float64(cpuEnd.Utime.Usec)/1e6) -
		(float64(cpuStart.Utime.Sec) + float64(cpuStart.Utime.Usec)/1e6)
	cpuSysSec := (float64(cpuEnd.Stime.Sec) + float64(cpuEnd.Stime.Usec)/1e6) -
		(float64(cpuStart.Stime.Sec) + float64(cpuStart.Stime.Usec)/1e6)
	cpuTotalSec := cpuUserSec + cpuSysSec
	cpuUtilPct := 100.0 * cpuTotalSec / (totalElapsed.Seconds() * float64(runtime.GOMAXPROCS(0)))

	pubMean := durMs(pipe.immPub.mean())
	delMean := durMs(pipe.immDel.mean())

	outageTargetRate := emitConcurrency * emitRatePerWorker
	t.Logf("╔══════════════════════════════════════════════════════════════╗")
	t.Logf("║            CHIP OUTAGE TEST RESULTS                          ║")
	t.Logf("╠══════════════════════════════════════════════════════════════╣")
	t.Logf("║ Target rate: %-47s ║", fmt.Sprintf("%d msg/s", outageTargetRate))
	t.Logf("║ EMIT                                                         ║")
	t.Logf("║   Events emitted: %-42d ║", emitCount.Load())
	t.Logf("║   Errors:         %-42d ║", emitErrors.Load())
	t.Logf("║   Elapsed:        %-42s ║", emitElapsed.Round(time.Millisecond))
	t.Logf("║   Actual rate:    %-42s ║", fmt.Sprintf("%.0f msg/s", float64(emitCount.Load())/emitElapsed.Seconds()))
	t.Logf("╠══════════════════════════════════════════════════════════════╣")
	t.Logf("║ DELIVERY                                                     ║")
	t.Logf("║   Server received:%-42d ║", srv.totalEvents.Load())
	t.Logf("║   Total elapsed:  %-42s ║", totalElapsed.Round(time.Millisecond))
	t.Logf("║   E2E rate:       %-42s ║", fmt.Sprintf("%.0f events/sec", float64(emitCount.Load())/totalElapsed.Seconds()))
	t.Logf("╠══════════════════════════════════════════════════════════════╣")
	t.Logf("║ OUTAGE CYCLES (up=%s / down=%s per cycle)         ║", upDuration, outageDuration)
	t.Logf("║   %-4s %-12s %-12s %-12s %-12s ║", "Cyc", "Peak queue", "Drain time", "Drain rate", "")
	cyclesMu.Lock()
	for _, r := range cycleResults {
		drainStr := r.drainElapsed.Round(time.Millisecond).String()
		if r.drainElapsed == 0 {
			drainStr = "timeout"
		}
		t.Logf("║   %-4d %-12d %-12s %-12s           ║",
			r.cycle, r.peakQueue, drainStr,
			fmt.Sprintf("%.0f/s", r.drainRate))
	}
	cyclesMu.Unlock()
	t.Logf("╠══════════════════════════════════════════════════════════════╣")
	t.Logf("║ PENDING QUEUE DEPTH (exact atomic counter)                   ║")
	t.Logf("║   Max:            %-42d ║", em.PendingMax())
	t.Logf("║   Current:        %-42d ║", em.PendingDepth())
	t.Logf("╠══════════════════════════════════════════════════════════════╣")
	t.Logf("║ PIPELINE LATENCY                                              ║")
	t.Logf("║                   %-10s %-10s %-10s %-10s ║", "n", "p50 (ms)", "p99 (ms)", "mean (ms)")
	t.Logf("║   Emit (INSERT):  %-10d %-10.2f %-10.2f %-10.2f ║",
		pipe.emitIns.count(), durMs(pipe.emitIns.percentile(0.50)),
		durMs(pipe.emitIns.percentile(0.99)), durMs(pipe.emitIns.mean()))
	t.Logf("║   Publish (gRPC): %-10d %-10.2f %-10.2f %-10.2f ║",
		pipe.immPub.count(), durMs(pipe.immPub.percentile(0.50)),
		durMs(pipe.immPub.percentile(0.99)), pubMean)
	t.Logf("║   MarkDelivered:  %-10d %-10.2f %-10.2f %-10.2f ║",
		pipe.immDel.count(), durMs(pipe.immDel.percentile(0.50)),
		durMs(pipe.immDel.percentile(0.99)), delMean)
	if pipe.batchPub.count() > 0 {
		t.Logf("║   Retransmit:     %-10d %-10.2f %-10s %-10.2f ║",
			pipe.batchPub.count(), durMs(pipe.batchPub.percentile(0.50)), "—",
			durMs(pipe.batchPub.mean()))
	}
	t.Logf("╠══════════════════════════════════════════════════════════════╣")
	t.Logf("║ PROCESS CPU (getrusage, GOMAXPROCS=%d)                       ║", runtime.GOMAXPROCS(0))
	t.Logf("║   User:           %-42s ║", fmt.Sprintf("%.2f s", cpuUserSec))
	t.Logf("║   System:         %-42s ║", fmt.Sprintf("%.2f s", cpuSysSec))
	t.Logf("║   Total:          %-42s ║", fmt.Sprintf("%.2f s", cpuTotalSec))
	t.Logf("║   Utilization:    %-42s ║", fmt.Sprintf("%.1f%% of %d cores × %.1fs wall", cpuUtilPct, runtime.GOMAXPROCS(0), totalElapsed.Seconds()))
	t.Logf("╚══════════════════════════════════════════════════════════════╝")

	assert.Equal(t, int64(0), emitErrors.Load(), "no emit errors expected")
	assert.GreaterOrEqual(t, srv.totalEvents.Load(), int64(totalEmitted),
		"server should have received all events (may include retransmit duplicates)")
}

// drainRunConfig defines the variable parameters for a single backlog drain run.
type drainRunConfig struct {
	retransmitBatchSize int
	publishBatchSize    int
}

func (c drainRunConfig) label() string {
	return fmt.Sprintf("rt%d_pb%d", c.retransmitBatchSize, c.publishBatchSize)
}

// drainRunResult holds the measured metrics from one backlog drain iteration.
type drainRunResult struct {
	cfg            drainRunConfig
	p1Elapsed      time.Duration
	p1InsertRate   float64
	p2DrainElapsed time.Duration
	p2DrainRate    float64
	p2InsP99       float64
	cpuUtilPct     float64
	peakHeapMB     float64
}

// TestFullStack_BacklogDrain measures maximum drain rate across multiple
// RetransmitBatchSize / PublishBatchSize configurations.
//
// Works with both the in-process mock AND a real Chip Ingress server.
// Phase 1 inserts events directly into the DB (no Chip needed), then
// Phases 2/3 drain via retransmit against whichever server is configured.
//
// Mock:
//
//	go test ./core/services/beholder/ -run TestFullStack_BacklogDrain -v -count=1 -timeout 30m
//
// Real Chip:
//
//	CHIP_INGRESS_TEST_ADDR=127.0.0.1:50051 go test ./core/services/beholder/ -run TestFullStack_BacklogDrain -v -count=1 -timeout 30m
func TestFullStack_BacklogDrain(t *testing.T) {
	db := directDB(t)

	ctx := testutils.Context(t)

	configs := []drainRunConfig{
		{retransmitBatchSize: 1000, publishBatchSize: 100},
		{retransmitBatchSize: 1000, publishBatchSize: 1000},
		//{retransmitBatchSize: 5000, publishBatchSize: 100},
		//{retransmitBatchSize: 5000, publishBatchSize: 1000},
		//{retransmitBatchSize: 10_000, publishBatchSize: 100},
		//{retransmitBatchSize: 10_000, publishBatchSize: 1000},
	}
	const backlogSize = 500_000
	const liveRate = 1000
	const producerConcurrency = 20

	results := make([]drainRunResult, 0, len(configs))

	// Pre-build the serialized proto payload once (used for all direct DB inserts).
	// This mirrors what Emit() does: NewEvent → EventToProto → proto.Marshal.
	payload := buildLoadTestPayload(256)
	attrs := loadEmitAttrs()
	sourceDomain, entityType, attrErr := beholder.ExtractSourceAndType(attrs...)
	require.NoError(t, attrErr, "extract source/type from attrs")
	attrMap := make(map[string]any)
	for i := 0; i+1 < len(attrs); i += 2 {
		if k, ok := attrs[i].(string); ok {
			attrMap[k] = attrs[i+1]
		}
	}
	sampleEvent, err := chipingress.NewEvent(sourceDomain, entityType, payload, attrMap)
	require.NoError(t, err, "build sample CloudEvent")
	samplePb, err := chipingress.EventToProto(sampleEvent)
	require.NoError(t, err, "convert CloudEvent to proto")
	protoPayload, err := proto.Marshal(samplePb)
	require.NoError(t, err, "marshal proto payload")

	usingRealChip := externalChipConfigured()
	if usingRealChip {
		t.Logf("Running against REAL Chip Ingress at %s", os.Getenv(envChipIngressTestAddr))
	} else {
		t.Logf("Running against in-process mock Chip (set %s for real Chip)", envChipIngressTestAddr)
	}

	for _, rc := range configs {
		rc := rc
		t.Run(rc.label(), func(t *testing.T) {
			// Clean the table between runs.
			_, truncErr := db.ExecContext(ctx, "DELETE FROM cre.chip_durable_events")
			require.NoError(t, truncErr, "truncate table between runs")

			srv, client := startChipIngressOrMock(t)
			if srv != nil {
				srv.setPublishDelay(2 * time.Millisecond)
			}
			store := beholdersvc.NewPgDurableEventStore(db)

			// ── PHASE 1: Build Backlog (direct DB inserts, no Chip) ─
			t.Logf("═══ PHASE 1: Inserting %d events directly into DB (retransmit=%d, publishBatch=%d) ═══",
				backlogSize, rc.retransmitBatchSize, rc.publishBatchSize)

			var cpuStart syscall.Rusage
			_ = syscall.Getrusage(syscall.RUSAGE_SELF, &cpuStart)

			p1Start := time.Now()
			var p1Errors atomic.Int64
			var p1Wg sync.WaitGroup
			eventsPerWorker := backlogSize / producerConcurrency
			for w := 0; w < producerConcurrency; w++ {
				p1Wg.Add(1)
				go func() {
					defer p1Wg.Done()
					for i := 0; i < eventsPerWorker; i++ {
						if _, e := store.Insert(ctx, protoPayload); e != nil {
							p1Errors.Add(1)
						}
					}
				}()
			}
			p1Wg.Wait()
			p1Elapsed := time.Since(p1Start)
			p1InsertRate := float64(backlogSize) / p1Elapsed.Seconds()

			t.Logf("Phase 1 done: %d events in %s (%.0f ev/s, errors=%d)",
				backlogSize, p1Elapsed.Round(time.Millisecond), p1InsertRate, p1Errors.Load())

			// ── PHASE 2: Start emitter and drain + live load ────────
			t.Logf("═══ PHASE 2: Chip UP — draining + %d msg/s live ═══", liveRate)

			phase2Stats := &pipelineDeliveryStats{}

			cfg := beholder.DefaultDurableEmitterConfig()
			cfg.QuietMode = true
			cfg.RetransmitInterval = 100 * time.Millisecond
			cfg.RetransmitAfter = 0 // all pending rows are eligible immediately
			cfg.RetransmitBatchSize = rc.retransmitBatchSize
			cfg.PublishTimeout = 5 * time.Second
			cfg.PublishBatchSize = rc.publishBatchSize
			cfg.PublishBatchWorkers = 4
			cfg.PublishBatchFlushInterval = 1000 * time.Millisecond
			cfg.PublishBatchChannelSize = 10000
			cfg.DisablePruning = true
			cfg.Hooks = newPipelineHooks(phase2Stats)

			em, emErr := beholder.NewDurableEmitter(store, client, cfg, logger.Test(t))
			require.NoError(t, emErr)

			em.Start(ctx)

			p2Start := time.Now()

			// Query actual pending count from DB (atomic counter is 0 since
			// we inserted directly, not through Emit).
			countPending := func() int64 {
				var n int64
				row := db.QueryRowContext(ctx, "SELECT count(*) FROM cre.chip_durable_events WHERE delivered_at IS NULL")
				if scanErr := row.Scan(&n); scanErr != nil {
					t.Logf("count pending: %v", scanErr)
				}
				return n
			}

			p2Ctx, p2Cancel := context.WithCancel(ctx)
			var p2LiveCount atomic.Int64
			var p2Wg sync.WaitGroup
			perWorkerRate := liveRate / producerConcurrency
			if perWorkerRate < 1 {
				perWorkerRate = 1
			}
			for w := 0; w < producerConcurrency; w++ {
				p2Wg.Add(1)
				go func() {
					defer p2Wg.Done()
					interval := time.Duration(float64(time.Second) / float64(perWorkerRate))
					ticker := time.NewTicker(interval)
					defer ticker.Stop()
					for {
						select {
						case <-p2Ctx.Done():
							return
						case <-ticker.C:
							_ = em.Emit(ctx, payload, loadEmitAttrs()...)
							p2LiveCount.Add(1)
						}
					}
				}()
			}

			steadyThreshold := int64(liveRate)
			drainTick := time.NewTicker(2 * time.Second)

			var p2DrainElapsed time.Duration
			var peakHeapBytes uint64
			for {
				<-drainTick.C

				var ms runtime.MemStats
				runtime.ReadMemStats(&ms)
				if ms.HeapInuse > peakHeapBytes {
					peakHeapBytes = ms.HeapInuse
				}

				depth := countPending()
				elapsed := time.Since(p2Start)
				drained := int64(backlogSize) - depth + p2LiveCount.Load()
				rate := float64(drained) / elapsed.Seconds()
				t.Logf("  [%s] pending=%d  drain_rate=%.0f ev/s  heap=%.1f MiB",
					elapsed.Round(100*time.Millisecond), depth, rate,
					float64(ms.HeapInuse)/(1024*1024))
				if depth <= steadyThreshold {
					p2DrainElapsed = elapsed
					break
				}
				if elapsed > 10*time.Minute {
					p2DrainElapsed = elapsed
					t.Logf("Phase 2 timed out with %d pending", depth)
					break
				}
			}
			drainTick.Stop()

			p2DrainRate := float64(backlogSize) / p2DrainElapsed.Seconds()
			p2InsP99 := durMs(phase2Stats.emitIns.percentile(0.99))

			// Stop live producers and shut down emitter.
			p2Cancel()
			p2Wg.Wait()
			em.Close()

			var cpuEnd syscall.Rusage
			_ = syscall.Getrusage(syscall.RUSAGE_SELF, &cpuEnd)
			cpuUser := (float64(cpuEnd.Utime.Sec) + float64(cpuEnd.Utime.Usec)/1e6) -
				(float64(cpuStart.Utime.Sec) + float64(cpuStart.Utime.Usec)/1e6)
			cpuSys := (float64(cpuEnd.Stime.Sec) + float64(cpuEnd.Stime.Usec)/1e6) -
				(float64(cpuStart.Stime.Sec) + float64(cpuStart.Stime.Usec)/1e6)
			totalWall := time.Since(p1Start).Seconds()
			cpuPct := 100.0 * (cpuUser + cpuSys) / (totalWall * float64(runtime.GOMAXPROCS(0)))

			peakHeapMB := float64(peakHeapBytes) / (1024 * 1024)

			t.Logf("Done: drain=%.0f ev/s, CPU=%.1f%%, peakHeap=%.1f MiB",
				p2DrainRate, cpuPct, peakHeapMB)

			results = append(results, drainRunResult{
				cfg:            rc,
				p1Elapsed:      p1Elapsed,
				p1InsertRate:   p1InsertRate,
				p2DrainElapsed: p2DrainElapsed,
				p2DrainRate:    p2DrainRate,
				p2InsP99:       p2InsP99,
				cpuUtilPct:     cpuPct,
				peakHeapMB:     peakHeapMB,
			})
		})
	}

	// ── Comparison Chart ───────────────────────────────────────────────
	t.Logf("")
	t.Logf("╔══════════════════════════════════════════════════════════════════════════════════════════════════════════╗")
	t.Logf("║   BACKLOG DRAIN CONFIG MATRIX  (backlog=%d, live=%d msg/s, retransmitInterval=100ms)               ║",
		backlogSize, liveRate)
	t.Logf("╠══════════╦══════════╦══════════════╦══════════════╦══════════╦════════╦════════╦══════════╦═══════════╣")
	t.Logf("║ Retrans  ║ Publish  ║ Drain+Live   ║ Drain Rate   ║ INS p99  ║ RPCs/  ║ CPU    ║ Peak     ║ Insert    ║")
	t.Logf("║ Batch    ║ Batch    ║ Elapsed      ║ (events/s)   ║ (contend)║ tick   ║ Util   ║ Heap     ║ Rate      ║")
	t.Logf("╠══════════╬══════════╬══════════════╬══════════════╬══════════╬════════╬════════╬══════════╬═══════════╣")
	for _, r := range results {
		rpcsPerTick := r.cfg.retransmitBatchSize / r.cfg.publishBatchSize
		if rpcsPerTick < 1 {
			rpcsPerTick = 1
		}
		t.Logf("║ %-8d ║ %-8d ║ %-12s ║ %-12s ║ %-8s ║ %-6d ║ %-6s ║ %-8s ║ %-9s ║",
			r.cfg.retransmitBatchSize,
			r.cfg.publishBatchSize,
			r.p2DrainElapsed.Round(time.Millisecond).String(),
			fmt.Sprintf("%.0f", r.p2DrainRate),
			fmt.Sprintf("%.1fms", r.p2InsP99),
			rpcsPerTick,
			fmt.Sprintf("%.1f%%", r.cpuUtilPct),
			fmt.Sprintf("%.1fMiB", r.peakHeapMB),
			fmt.Sprintf("%.0f/s", r.p1InsertRate),
		)
	}
	t.Logf("╚══════════╩══════════╩══════════════╩══════════════╩══════════╩════════╩════════╩══════════╩═══════════╝")

	if len(results) > 1 {
		best := results[0]
		for _, r := range results[1:] {
			if r.p2DrainRate > best.p2DrainRate {
				best = r
			}
		}
		t.Logf("")
		t.Logf("Winner: RetransmitBatch=%d + PublishBatch=%d → %.0f ev/s, INS p99=%.1fms, heap=%.1f MiB",
			best.cfg.retransmitBatchSize, best.cfg.publishBatchSize, best.p2DrainRate, best.p2InsP99, best.peakHeapMB)
	}
}

// TestFullStack_SlowChip simulates a slow Chip server (high latency per
// publish). This tests whether the async design keeps Emit() fast even
// when gRPC is slow.
func TestFullStack_SlowChip(t *testing.T) {
	skipIfExternalChip(t, "inject publish latency on mock server")

	db := directDB(t)
	srv, client := startChipIngressOrMock(t)
	require.NotNil(t, srv)
	srv.setPublishDelay(100 * time.Millisecond) // 50ms per publish = ~20 RPS max
	store := beholdersvc.NewPgDurableEventStore(db)

	cfg := beholder.DefaultDurableEmitterConfig()
	cfg.RetransmitInterval = 500 * time.Millisecond
	cfg.RetransmitAfter = 2 * time.Second
	cfg.RetransmitBatchSize = 50

	em, err := beholder.NewDurableEmitter(store, client, cfg, logger.Test(t))
	require.NoError(t, err)

	ctx := testutils.Context(t)
	em.Start(ctx)
	defer em.Close()

	const totalEvents = 200

	// Emit should still be fast because it only does DB insert (async gRPC).
	start := time.Now()
	for i := 0; i < totalEvents; i++ {
		require.NoError(t, em.Emit(ctx, []byte("slow-chip-event"), loadEmitAttrs()...))
	}
	emitElapsed := time.Since(start)

	t.Logf("Emit %d events in %s (%.0f events/sec) despite 50ms server latency",
		totalEvents, emitElapsed.Round(time.Millisecond),
		float64(totalEvents)/emitElapsed.Seconds())

	// Emit rate should be much higher than the server can handle,
	// proving the async design works.
	assert.Less(t, emitElapsed, 5*time.Second,
		"Emit() should not be bottlenecked by slow gRPC server")

	// Wait for everything to eventually deliver.
	require.Eventually(t, func() bool {
		pending, _ := store.ListPending(ctx, time.Now().Add(time.Hour), 1)
		return len(pending) == 0
	}, 60*time.Second, 200*time.Millisecond, "all events should eventually deliver")

	t.Logf("All %d events delivered (server received %d, including retransmits)",
		totalEvents, srv.totalEvents.Load())
}

// Benchmark_FullStack_EmitThroughput benchmarks the Emit() path with real Postgres
// and a fast mock gRPC server. This gives the upper bound of events/sec.
func Benchmark_FullStack_EmitThroughput(b *testing.B) {
	db := directDB(b)
	_, client := startChipIngressOrMock(b)
	store := beholdersvc.NewPgDurableEventStore(db)

	cfg := beholder.DefaultDurableEmitterConfig()
	em, err := beholder.NewDurableEmitter(store, client, cfg, logger.Nop())
	require.NoError(b, err)

	ctx := testutils.Context(b)
	em.Start(ctx)
	defer em.Close()

	payload := buildLoadTestPayload(256)

	b.ResetTimer()
	for b.Loop() {
		err := em.Emit(ctx, payload, loadEmitAttrs()...)
		require.NoError(b, err)
	}
}

// Benchmark_FullStack_EmitPayloadSizes benchmarks Emit throughput at
// different payload sizes to understand the DB I/O impact.
func Benchmark_FullStack_EmitPayloadSizes(b *testing.B) {
	sizes := []int{64, 256, 1024, 4096}
	for _, size := range sizes {
		b.Run(fmt.Sprintf("%dB", size), func(b *testing.B) {
			db := directDB(b)
			_, client := startChipIngressOrMock(b)
			store := beholdersvc.NewPgDurableEventStore(db)

			cfg := beholder.DefaultDurableEmitterConfig()
			em, err := beholder.NewDurableEmitter(store, client, cfg, logger.Nop())
			require.NoError(b, err)

			ctx := testutils.Context(b)
			em.Start(ctx)
			defer em.Close()

			payload := buildLoadTestPayload(size)

			b.ResetTimer()
			for b.Loop() {
				err := em.Emit(ctx, payload, loadEmitAttrs()...)
				require.NoError(b, err)
			}
		})
	}
}

// ---------- 1k TPS Target Tests ----------

// tpsSummaryBlocks collects human-readable result blocks from each TPS test;
// TestMain prints them together after the full test run.
var (
	tpsSummaryMu     sync.Mutex
	tpsSummaryBlocks []string

	tpsRampMu   sync.Mutex
	tpsRampRows []string

	tpsPayloadMu   sync.Mutex
	tpsPayloadRows []string
)

func appendTPSummaryBlock(title string, lines ...string) {
	tpsSummaryMu.Lock()
	defer tpsSummaryMu.Unlock()
	var b strings.Builder
	b.WriteString("--- ")
	b.WriteString(title)
	b.WriteString(" ---\n")
	for _, ln := range lines {
		b.WriteString(ln)
		b.WriteByte('\n')
	}
	tpsSummaryBlocks = append(tpsSummaryBlocks, b.String())
}

func TestMain(m *testing.M) {
	code := m.Run()
	tpsSummaryMu.Lock()
	blocks := append([]string(nil), tpsSummaryBlocks...)
	tpsSummaryMu.Unlock()
	if len(blocks) > 0 {
		fmt.Println()
		fmt.Println(strings.Repeat("=", 72))
		fmt.Println("TPS LOAD TEST SUMMARY (full run)")
		fmt.Println(strings.Repeat("=", 72))
		for _, blk := range blocks {
			fmt.Print(blk)
			fmt.Println()
		}
		fmt.Println(strings.Repeat("=", 72))
	}
	os.Exit(code)
}

func progressBar(pct float64, width int) string {
	if pct < 0 {
		pct = 0
	}
	if pct > 1 {
		pct = 1
	}
	filled := int(pct * float64(width))
	if filled > width {
		filled = width
	}
	var b strings.Builder
	b.WriteByte('[')
	for i := 0; i < width; i++ {
		if i < filled {
			b.WriteRune('█')
		} else {
			b.WriteRune('░')
		}
	}
	b.WriteByte(']')
	b.WriteString(fmt.Sprintf(" %3.0f%%", pct*100))
	return b.String()
}

// directDB opens a real (non-txdb) Postgres connection for concurrent load tests.
// pgtest.NewSqlxDB uses txdb: one shared transaction per pool. Any SQL error
// aborts that transaction (SQLSTATE 25P02 on later queries). DurableEmitter’s
// concurrent purge/retransmit/mark-delivered + many goroutines requires
// autocommit statements and a real pool, not txdb.
func directDB(t testing.TB) *sqlx.DB {
	t.Helper()
	testutils.SkipShortDB(t)
	dbURL := string(env.DatabaseURL.Get())
	if dbURL == "" {
		t.Fatal("CL_DATABASE_URL is required for TPS tests")
	}
	// Append synchronous_commit=off to the DSN so every pooled connection
	// skips WAL fsync. Events are durable via Chip ACK; retransmit recovers
	// on crash. This typically 2-3x insert throughput.
	sep := "&"
	if !strings.Contains(dbURL, "?") {
		sep = "?"
	}
	dbURL += sep + "options=-c%20synchronous_commit%3Doff"

	db, err := sqlx.Open("postgres", dbURL)
	require.NoError(t, err)
	require.NoError(t, db.Ping())
	db.SetMaxOpenConns(60)
	db.SetMaxIdleConns(30)
	db.SetConnMaxIdleTime(30 * time.Second)

	// Kill stale connections from previous Ctrl+C'd runs, then TRUNCATE
	// for a clean table with no dead-tuple bloat.
	_, _ = db.Exec(`SELECT pg_terminate_backend(pid) FROM pg_stat_activity
		WHERE datname = current_database() AND pid <> pg_backend_pid()
		AND state = 'idle' AND state_change < now() - interval '10 seconds'`)
	_, _ = db.Exec("TRUNCATE cre.chip_durable_events RESTART IDENTITY")
	t.Cleanup(func() {
		_, _ = db.Exec("TRUNCATE cre.chip_durable_events RESTART IDENTITY")
		_ = db.Close()
	})
	return db
}

// emitLatencyStats tracks Emit() call latencies.
type emitLatencyStats struct {
	mu       sync.Mutex
	samples  []time.Duration
	failures atomic.Int64
}

func (s *emitLatencyStats) record(d time.Duration) {
	s.mu.Lock()
	s.samples = append(s.samples, d)
	s.mu.Unlock()
}

func (s *emitLatencyStats) percentile(p float64) time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.samples) == 0 {
		return 0
	}
	sorted := make([]time.Duration, len(s.samples))
	copy(sorted, s.samples)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
	idx := int(float64(len(sorted)-1) * p)
	return sorted[idx]
}

func (s *emitLatencyStats) count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.samples)
}

func (s *emitLatencyStats) mean() time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.samples) == 0 {
		return 0
	}
	var sum time.Duration
	for _, v := range s.samples {
		sum += v
	}
	return sum / time.Duration(len(s.samples))
}

func (s *emitLatencyStats) sum() time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()
	var t time.Duration
	for _, v := range s.samples {
		t += v
	}
	return t
}

// pipelineDeliveryStats aggregates DurableEmitterHooks samples to compare
// Emit (DB Insert) vs Chip Publish vs store MarkDelivered cost.
type pipelineDeliveryStats struct {
	emitIns                            emitLatencyStats // store.Insert latency (blocks the Emit caller)
	immPub, immDel, batchPub, batchDel emitLatencyStats
	// batchLoopPub tracks latency of PublishBatch RPCs in the batch publish loop.
	batchLoopPub emitLatencyStats
	// batchLoopDel tracks latency of MarkDeliveredBatch in the batch publish loop.
	batchLoopDel           emitLatencyStats
	batchLoopPubEvents     atomic.Int64 // total events across successful batch publishes
	batchLoopMarkEvents    atomic.Int64 // total events marked delivered via batch
	immPubErr, batchPubErr atomic.Int64
	batchPubEventErrs      atomic.Int64
}

func newPipelineHooks(p *pipelineDeliveryStats) *beholder.DurableEmitterHooks {
	return &beholder.DurableEmitterHooks{
		OnEmitInsert: func(d time.Duration, _ error) {
			p.emitIns.record(d)
		},
		OnImmediatePublish: func(d time.Duration, err error) {
			if err != nil {
				p.immPubErr.Add(1)
			}
			p.immPub.record(d)
		},
		OnImmediateDelete: func(d time.Duration, _ error) {
			p.immDel.record(d)
		},
		OnBatchPublish: func(d time.Duration, batchSize int, err error) {
			if err != nil {
				p.batchPubErr.Add(1)
				p.batchPubEventErrs.Add(int64(batchSize))
			} else {
				p.batchLoopPubEvents.Add(int64(batchSize))
			}
			p.batchLoopPub.record(d)
		},
		OnBatchMarkDelivered: func(d time.Duration, count int) {
			p.batchLoopMarkEvents.Add(int64(count))
			p.batchLoopDel.record(d)
		},
		OnRetransmitBatchPublish: func(d time.Duration, eventCount int, err error) {
			if err != nil {
				p.batchPubErr.Add(1)
				p.batchPubEventErrs.Add(int64(eventCount))
			}
			p.batchPub.record(d)
		},
		OnRetransmitBatchDeletes: func(d time.Duration, _ int) {
			p.batchDel.record(d)
		},
	}
}

func durMs(d time.Duration) float64 {
	return float64(d.Microseconds()) / 1000.0
}

// dbE2ELatencyStats queries the chip_durable_events table for events with both
// created_at and delivered_at set and returns p50, p99, mean of (delivered_at - created_at).
// Requires DisablePruning=true so rows aren't deleted before we can read them.
type dbLatencyResult struct {
	count int
	p50   time.Duration
	p99   time.Duration
	mean  time.Duration
}

func queryDBE2ELatency(ctx context.Context, db *sqlx.DB) (dbLatencyResult, error) {
	const q = `
SELECT extract(epoch FROM (delivered_at - created_at)) AS latency_sec
FROM cre.chip_durable_events
WHERE delivered_at IS NOT NULL
ORDER BY latency_sec ASC`

	var latencies []float64
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return dbLatencyResult{}, fmt.Errorf("query e2e latency: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var sec float64
		if err := rows.Scan(&sec); err != nil {
			return dbLatencyResult{}, fmt.Errorf("scan e2e latency: %w", err)
		}
		latencies = append(latencies, sec)
	}
	if err := rows.Err(); err != nil {
		return dbLatencyResult{}, fmt.Errorf("rows e2e latency: %w", err)
	}
	if len(latencies) == 0 {
		return dbLatencyResult{}, nil
	}

	var sum float64
	for _, v := range latencies {
		sum += v
	}
	mean := sum / float64(len(latencies))
	p50 := latencies[int(float64(len(latencies)-1)*0.50)]
	p99 := latencies[int(float64(len(latencies)-1)*0.99)]

	return dbLatencyResult{
		count: len(latencies),
		p50:   time.Duration(p50 * float64(time.Second)),
		p99:   time.Duration(p99 * float64(time.Second)),
		mean:  time.Duration(mean * float64(time.Second)),
	}, nil
}

// chipIngressTargetDescription labels latency logs: mock gRPC server vs external Chip Ingress.
func chipIngressTargetDescription(srv *loadTestServer) string {
	if srv != nil {
		return "in-process mock ChipIngress (loadTestServer)"
	}
	addr := strings.TrimSpace(os.Getenv(envChipIngressTestAddr))
	if addr == "" {
		return "external Chip Ingress"
	}
	return fmt.Sprintf("external Chip Ingress (%s)", addr)
}

func logPipelineDeliverySummary(t *testing.T, pipe *pipelineDeliveryStats) {
	t.Helper()
	ipN := pipe.immPub.count()
	idN := pipe.immDel.count()
	t.Logf("Pipeline — immediate Publish: n=%d errs=%d p50=%.3f ms p99=%.3f ms mean=%.3f ms Σ=%.1f ms",
		ipN, pipe.immPubErr.Load(),
		durMs(pipe.immPub.percentile(0.50)), durMs(pipe.immPub.percentile(0.99)),
		durMs(pipe.immPub.mean()), durMs(pipe.immPub.sum()))
	t.Logf("Pipeline — immediate MarkDelivered: n=%d p50=%.3f ms p99=%.3f ms mean=%.3f ms Σ=%.1f ms",
		idN,
		durMs(pipe.immDel.percentile(0.50)), durMs(pipe.immDel.percentile(0.99)),
		durMs(pipe.immDel.mean()), durMs(pipe.immDel.sum()))

	bpN := pipe.batchPub.count()
	if bpN > 0 {
		t.Logf("Pipeline — retransmit Publish (serial): rpcs=%d rpc_errs=%d evt_errs=%d p50=%.3f ms mean=%.3f ms | retransmit_mark_delivered_hooks=%d mean_loop=%.3f ms",
			bpN, pipe.batchPubErr.Load(), pipe.batchPubEventErrs.Load(),
			durMs(pipe.batchPub.percentile(0.50)), durMs(pipe.batchPub.mean()),
			pipe.batchDel.count(), durMs(pipe.batchDel.mean()))
	}

	if ipN >= 50 && idN >= 50 {
		pm, dm := durMs(pipe.immPub.mean()), durMs(pipe.immDel.mean())
		switch {
		case pm > 3*dm && pm > 0.5:
			t.Logf("Bottleneck hint: Publish mean %.3f ms ≫ MarkDelivered mean %.3f ms — likely Chip / gRPC bound", pm, dm)
		case dm > 3*pm && dm > 0.5:
			t.Logf("Bottleneck hint: MarkDelivered mean %.3f ms ≫ Publish mean %.3f ms — likely Postgres UPDATE bound", dm, pm)
		default:
			t.Logf("Bottleneck hint: Publish %.3f ms vs MarkDelivered %.3f ms comparable (per successful immediate delivery)", pm, dm)
		}
	} else {
		t.Logf("Bottleneck hint: few completed immediate deliveries in window (pub=%d del=%d); extend duration or check async backlog", ipN, idN)
	}
}

// rateLimitEmitResult is the outcome of runRateLimitedEmit.
type rateLimitEmitResult struct {
	stats *emitLatencyStats
	// maxQueueDepth is the maximum observed pending row count in cre.chip_durable_events
	// (delivered_at IS NULL).
	// during the emit window (polled periodically; nil DB disables sampling).
	maxQueueDepth int64
	// maxQueuePayloadBytes is the maximum observed sum(octet_length(payload)) for
	// rows still in the queue (serialized CloudEvent bytes stored in BYTEA).
	maxQueuePayloadBytes int64
	// ImmPublishFails is the count of failed immediate Publish RPCs in this window (one event each; needs retransmit).
	ImmPublishFails int64
	// BatchPublishFailEvents is the sum of batch sizes for failed PublishBatch calls in this window.
	BatchPublishFailEvents int64
}

// formatPubFailColumn formats publish-failure counts for result tables (8-char column).
// If there were failed batches, shows "imm+batchEv" when it fits, else "imm+…".
func formatPubFailColumn(imm, batchEv int64) string {
	if batchEv == 0 {
		return fmt.Sprintf("%-8d", imm)
	}
	combo := fmt.Sprintf("%d+%d", imm, batchEv)
	if len(combo) <= 8 {
		return fmt.Sprintf("%-8s", combo)
	}
	return fmt.Sprintf("%-8s", fmt.Sprintf("%d+..", imm))
}

func bumpMaxQueueDepth(maxQ *atomic.Int64, c int64) {
	for {
		old := maxQ.Load()
		if c <= old {
			return
		}
		if maxQ.CompareAndSwap(old, c) {
			return
		}
	}
}

func bumpMaxQueuePayloadBytes(maxB *atomic.Int64, b int64) {
	for {
		old := maxB.Load()
		if b <= old {
			return
		}
		if maxB.CompareAndSwap(old, b) {
			return
		}
	}
}

// queuePayloadStats returns pending row count and payload bytes (delivered_at IS NULL).
func queuePayloadStats(db *sqlx.DB, ctx context.Context) (rows int64, payloadBytes int64, err error) {
	err = db.QueryRowContext(ctx,
		`SELECT count(*), coalesce(sum(octet_length(payload)), 0) FROM cre.chip_durable_events WHERE delivered_at IS NULL`,
	).Scan(&rows, &payloadBytes)
	return rows, payloadBytes, err
}

func formatQueueKB(payloadBytes int64) string {
	if payloadBytes == 0 {
		return "0.0"
	}
	return fmt.Sprintf("%.1f", float64(payloadBytes)/1024.0)
}

// runRateLimitedEmit emits events at a target rate for the given duration,
// using the specified concurrency. Returns latency stats and optional max queue depth.
// If maxQueueDB is non-nil, polls cre.chip_durable_events during the emit window to
// record peak backlog (async publish may lag inserts).
// If progressLabel is non-empty, prints a live progress bar and emit count to stdout every 500ms.
// If pipe is non-nil (same *pipelineDeliveryStats wired via cfg.Hooks), ImmPublishFails and
// BatchPublishFailEvents are deltas for this emit window only.
func runRateLimitedEmit(
	ctx context.Context,
	t testing.TB,
	em *beholder.DurableEmitter,
	targetTPS int,
	duration time.Duration,
	concurrency int,
	payloadSize int,
	progressLabel string,
	maxQueueDB *sqlx.DB,
	pipe *pipelineDeliveryStats,
) *rateLimitEmitResult {
	t.Helper()

	var imm0, batchEv0 int64
	if pipe != nil {
		imm0 = pipe.immPubErr.Load()
		batchEv0 = pipe.batchPubEventErrs.Load()
	}

	stats := &emitLatencyStats{}
	var maxQ, maxPayloadBytes atomic.Int64
	var emitCount atomic.Int64
	payload := buildLoadTestPayload(payloadSize)

	// Each worker gets an equal share of the target TPS.
	perWorkerTPS := targetTPS / concurrency
	if perWorkerTPS < 1 {
		perWorkerTPS = 1
	}
	interval := time.Duration(float64(time.Second) / float64(perWorkerTPS))

	pollCtx, pollCancel := context.WithCancel(ctx)
	defer pollCancel()
	if maxQueueDB != nil {
		go func() {
			ticker := time.NewTicker(50 * time.Millisecond)
			defer ticker.Stop()
			for {
				select {
				case <-pollCtx.Done():
					return
				case <-ticker.C:
					c, b, err := queuePayloadStats(maxQueueDB, pollCtx)
					if err != nil {
						continue
					}
					bumpMaxQueueDepth(&maxQ, c)
					bumpMaxQueuePayloadBytes(&maxPayloadBytes, b)
				}
			}
		}()
	}

	var wg sync.WaitGroup

	if progressLabel != "" {
		startAll := time.Now()
		deadline := time.After(duration)
		done := make(chan struct{})
		go func() {
			ticker := time.NewTicker(500 * time.Millisecond)
			defer ticker.Stop()
			for {
				select {
				case <-ctx.Done():
					fmt.Fprintf(os.Stdout, "\n")
					return
				case <-done:
					return
				case <-ticker.C:
					elapsed := time.Since(startAll)
					pct := float64(elapsed) / float64(duration)
					if pct >= 1 {
						fmt.Fprintf(os.Stdout, "\r%s %s | %s / %s | emits=%d\n",
							progressBar(1, 36), progressLabel,
							duration.Round(time.Millisecond), duration.Round(time.Millisecond), emitCount.Load())
						return
					}
					fmt.Fprintf(os.Stdout, "\r%s %s | %s / %s | emits=%d   ",
						progressBar(pct, 36), progressLabel,
						elapsed.Round(time.Millisecond), duration.Round(time.Millisecond), emitCount.Load())
				}
			}
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-deadline
			close(done)
		}()
	}

	for w := 0; w < concurrency; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ticker := time.NewTicker(interval)
			defer ticker.Stop()
			localDeadline := time.After(duration)

			for {
				select {
				case <-localDeadline:
					return
				case <-ctx.Done():
					return
				case <-ticker.C:
					start := time.Now()
					if err := em.Emit(ctx, payload, loadEmitAttrs()...); err != nil {
						stats.failures.Add(1)
					} else {
						emitCount.Add(1)
						stats.record(time.Since(start))
					}
				}
			}
		}()
	}

	wg.Wait()
	if maxQueueDB != nil {
		c, b, err := queuePayloadStats(maxQueueDB, ctx)
		if err == nil {
			bumpMaxQueueDepth(&maxQ, c)
			bumpMaxQueuePayloadBytes(&maxPayloadBytes, b)
		}
	}
	res := &rateLimitEmitResult{
		stats:                stats,
		maxQueueDepth:        maxQ.Load(),
		maxQueuePayloadBytes: maxPayloadBytes.Load(),
	}
	if pipe != nil {
		res.ImmPublishFails = pipe.immPubErr.Load() - imm0
		res.BatchPublishFailEvents = pipe.batchPubEventErrs.Load() - batchEv0
	}
	return res
}

// TestTPS_RampUp tests the durable emitter at increasing TPS levels to find
// the throughput ceiling. Each level gets its own DurableEmitter to avoid
// carry-over. Measures achieved rate, Emit() latency, and queue depth.
func TestTPS_RampUp(t *testing.T) {
	levels := []int{100, 500, 1000, 2000}
	testStart := time.Now()

	tpsRampMu.Lock()
	tpsRampRows = nil
	tpsRampMu.Unlock()

	t.Logf("TPS ramp-up: levels=%v (each level: fresh DB + server + emitter)", levels)

	t.Logf("╔══════════════════════════════════════════════════════════════════════════════════════════════╗")
	t.Logf("║                              TPS RAMP-UP TEST RESULTS                                        ║")
	t.Logf("╠═══════════╦══════════╦═════════════╦══════════╦══════════╦══════════╦══════════╦══════════╦══════════╣")
	t.Logf("║ Target    ║ Achieved ║ Total emits ║ Emit p50 ║ Emit p99 ║ Pub fail ║ Q max    ║ Q end    ║ Q max    ║")
	t.Logf("║ TPS       ║ TPS      ║ (success)   ║ (ms)     ║ (ms)     ║ (retry)* ║ (rows)   ║ (rows)   ║ (KB)*    ║")
	t.Logf("╠═══════════╬══════════╬═════════════╬══════════╬══════════╬══════════╬══════════╬══════════╬══════════╣")

	for _, targetTPS := range levels {
		t.Run(fmt.Sprintf("%d_tps", targetTPS), func(t *testing.T) {
			levelStart := time.Now()
			t.Logf(">>> level %d TPS: provisioning direct DB + Chip endpoint...", targetTPS)

			db := directDB(t)
			_, client := startChipIngressOrMock(t)
			store := beholdersvc.NewPgDurableEventStore(db)

			cfg := beholder.DefaultDurableEmitterConfig()
			cfg.RetransmitInterval = 1 * time.Second
			cfg.RetransmitAfter = 3 * time.Second
			cfg.RetransmitBatchSize = 500
			pipe := &pipelineDeliveryStats{}
			cfg.Hooks = newPipelineHooks(pipe)

			em, err := beholder.NewDurableEmitter(store, client, cfg, logger.Nop())
			require.NoError(t, err)
			ctx := testutils.Context(t)
			em.Start(ctx)
			defer em.Close()

			const duration = 1 * time.Minute
			const concurrency = 20

			t.Logf(">>> level %d TPS: emitting for %s @ concurrency=%d (progress bar on stdout)", targetTPS, duration, concurrency)
			emitRes := runRateLimitedEmit(ctx, t, em, targetTPS, duration, concurrency, 256,
				fmt.Sprintf("ramp_up/%d_tps", targetTPS), db, pipe)
			stats := emitRes.stats
			emitPhase := time.Since(levelStart)
			t.Logf(">>> level %d TPS: emit phase wall time %s", targetTPS, emitPhase.Round(time.Millisecond))

			// Brief pause for async publishes to complete.
			t.Logf(">>> level %d TPS: sleeping 2s for async publishes...", targetTPS)
			time.Sleep(2 * time.Second)

			t.Logf(">>> level %d TPS: pipeline delivery (Publish vs Delete)", targetTPS)
			logPipelineDeliverySummary(t, pipe)

			achieved := float64(stats.count()) / duration.Seconds()
			p50 := stats.percentile(0.50)
			p99 := stats.percentile(0.99)

			queueEnd, _, err := queuePayloadStats(db, ctx)
			require.NoError(t, err)

			totalEmits := stats.count()
			if stats.failures.Load() > 0 {
				t.Logf(">>> level %d TPS: Emit() (DB insert) failures: %d", targetTPS, stats.failures.Load())
			}
			rowLine := fmt.Sprintf("║ %-9d ║ %-8.0f ║ %-11d ║ %-8.2f ║ %-8.2f ║ %-8s ║ %-8d ║ %-8d ║ %-8s ║",
				targetTPS, achieved, totalEmits,
				float64(p50.Microseconds())/1000.0,
				float64(p99.Microseconds())/1000.0,
				formatPubFailColumn(emitRes.ImmPublishFails, emitRes.BatchPublishFailEvents),
				emitRes.maxQueueDepth, queueEnd,
				formatQueueKB(emitRes.maxQueuePayloadBytes))
			t.Log(rowLine)

			tpsRampMu.Lock()
			tpsRampRows = append(tpsRampRows, rowLine)
			tpsRampMu.Unlock()
		})
	}

	t.Logf("╚═══════════╩══════════╩═════════════╩══════════╩══════════╩══════════╩══════════╩══════════╩══════════╝")
	t.Logf("* Pub fail: immediate Publish RPC errors (events need retransmit). a+b = a immediate fails + b events in failed PublishBatch. " +
		"Emit() insert failures are logged per level if non-zero.")
	t.Logf("* Q max / Q end: peak & final row counts (polled ~50ms; Q end after settle). Q max KB* = sum(octet_length(payload))/1024 for queued rows.")
	t.Logf("TestTPS_RampUp finished in %s", time.Since(testStart).Round(time.Millisecond))

	summaryLines := []string{
		fmt.Sprintf("total wall clock: %s", time.Since(testStart).Round(time.Millisecond)),
		"╔══════════════════════════════════════════════════════════════════════════════════════════════╗",
		"║                              TPS RAMP-UP TEST RESULTS                                        ║",
		"╠═══════════╦══════════╦═════════════╦══════════╦══════════╦══════════╦══════════╦══════════╦══════════╣",
		"║ Target    ║ Achieved ║ Total emits ║ Emit p50 ║ Emit p99 ║ Pub fail ║ Q max    ║ Q end    ║ Q max    ║",
		"║ TPS       ║ TPS      ║ (success)   ║ (ms)     ║ (ms)     ║ (retry)* ║ (rows)   ║ (rows)   ║ (KB)*    ║",
		"╠═══════════╬══════════╬═════════════╬══════════╬══════════╬══════════╬══════════╬══════════╬══════════╣",
	}
	tpsRampMu.Lock()
	summaryLines = append(summaryLines, tpsRampRows...)
	tpsRampMu.Unlock()
	summaryLines = append(summaryLines, "╚═══════════╩══════════╩═════════════╩══════════╩══════════╩══════════╩══════════╩══════════╩══════════╝",
		"* Q max KB* = sum(octet_length(payload))/1024 for queued rows (see test log footnotes).")
	appendTPSummaryBlock("TestTPS_RampUp", summaryLines...)
}

// TestTPS_Sustained1k runs at exactly 1000 TPS for 60 seconds and verifies
// the pipeline keeps up: deletes match inserts, queue stays bounded, and
// Emit() latency stays low.
func TestTPS_Sustained1k(t *testing.T) {
	testStart := time.Now()
	t.Logf("TestTPS_Sustained1k: provisioning DB + Chip server + emitter...")

	db := directDB(t)
	_, client := startChipIngressOrMock(t)
	store := beholdersvc.NewPgDurableEventStore(db)

	cfg := beholder.DefaultDurableEmitterConfig()
	cfg.RetransmitInterval = 1 * time.Second
	cfg.RetransmitAfter = 3 * time.Second
	cfg.RetransmitBatchSize = 500
	pipe := &pipelineDeliveryStats{}
	cfg.Hooks = newPipelineHooks(pipe)

	em, err := beholder.NewDurableEmitter(store, client, cfg, logger.Nop())
	require.NoError(t, err)

	ctx := testutils.Context(t)
	em.Start(ctx)
	defer em.Close()

	const targetTPS = 1000
	const duration = 60 * time.Second
	const concurrency = 20

	t.Logf("Emit phase: target=%d TPS for %s @ concurrency=%d (progress bar on stdout)", targetTPS, duration, concurrency)
	emitStart := time.Now()

	emitRes := runRateLimitedEmit(ctx, t, em, targetTPS, duration, concurrency, 256, "sustained_1k", db, pipe)
	stats := emitRes.stats

	achievedTPS := float64(stats.count()) / duration.Seconds()
	t.Logf("Emit phase complete in %s: %d events (%.0f TPS)", time.Since(emitStart).Round(time.Millisecond), stats.count(), achievedTPS)

	// Wait for the pipeline to drain.
	t.Logf("Waiting for pipeline to drain...")
	drainStart := time.Now()
	require.Eventually(t, func() bool {
		var count int64
		_ = db.QueryRow("SELECT count(*) FROM cre.chip_durable_events").Scan(&count)
		return count == 0
	}, 30*time.Second, 500*time.Millisecond, "pipeline should drain after emit phase ends")
	drainTime := time.Since(drainStart)

	t.Logf("Pipeline delivery after drain (full async + retransmit settled):")
	logPipelineDeliverySummary(t, pipe)

	t.Logf("╔════════════════════════════════════════════════════╗")
	t.Logf("║       SUSTAINED 1k TPS TEST RESULTS               	║")
	t.Logf("╠════════════════════════════════════════════════════╣")
	t.Logf("║ Target TPS:       %-6d                          	║", targetTPS)
	t.Logf("║ Duration:         %-6s                          	║", duration)
	t.Logf("║ Total emitted:    %-6d                          	║", stats.count())
	t.Logf("║ Achieved TPS:     %-6.0f                          	║", achievedTPS)
	t.Logf("║ Pub fail (retry): %-8s (1st+batch ev)          	║", formatPubFailColumn(emitRes.ImmPublishFails, emitRes.BatchPublishFailEvents))
	t.Logf("║ Emit insert fail: %-6d (DB path)               	║", stats.failures.Load())
	t.Logf("║ Emit p50 latency: %-6.2f ms                      	║", float64(stats.percentile(0.50).Microseconds())/1000.0)
	t.Logf("║ Emit p99 latency: %-6.2f ms                      	║", float64(stats.percentile(0.99).Microseconds())/1000.0)
	t.Logf("║ Queue max (emit): %-6d rows                     	║", emitRes.maxQueueDepth)
	t.Logf("║ Queue max (emit): %-10s KB payload*             	║", formatQueueKB(emitRes.maxQueuePayloadBytes))
	t.Logf("║ Drain time:       %-6s                          	║", drainTime.Round(time.Millisecond))
	t.Logf("╚════════════════════════════════════════════════════╝")
	t.Logf("* Queue KB = sum(octet_length(payload))/1024 for queued rows (excludes index/heap overhead).")
	t.Logf("TestTPS_Sustained1k finished in %s", time.Since(testStart).Round(time.Millisecond))

	appendTPSummaryBlock("TestTPS_Sustained1k",
		fmt.Sprintf("total wall clock: %s", time.Since(testStart).Round(time.Millisecond)),
		fmt.Sprintf("emit phase: %s", time.Since(emitStart).Round(time.Millisecond)),
		fmt.Sprintf("target TPS: %d, achieved: %.0f, pub_fail imm/batch_ev: %d/%d, emit_insert_fail: %d",
			targetTPS, achievedTPS, emitRes.ImmPublishFails, emitRes.BatchPublishFailEvents, stats.failures.Load()),
		fmt.Sprintf("emit p50/p99 ms: %.2f / %.2f", float64(stats.percentile(0.50).Microseconds())/1000.0, float64(stats.percentile(0.99).Microseconds())/1000.0),
		fmt.Sprintf("queue max during emit: %d rows, %s KB payload (sum octet_length/1024)", emitRes.maxQueueDepth, formatQueueKB(emitRes.maxQueuePayloadBytes)),
		fmt.Sprintf("pipeline imm Publish/MarkDelivered means ms: %.3f / %.3f (n=%d/%d)", durMs(pipe.immPub.mean()), durMs(pipe.immDel.mean()), pipe.immPub.count(), pipe.immDel.count()),
		fmt.Sprintf("drain time: %s", drainTime.Round(time.Millisecond)),
	)

	assert.GreaterOrEqual(t, achievedTPS, float64(targetTPS)*0.9,
		"should achieve at least 90%% of target TPS")
	assert.Equal(t, int64(0), stats.failures.Load(),
		"no Emit() calls should fail")
	assert.Less(t, stats.percentile(0.99), 50*time.Millisecond,
		"p99 Emit() latency should be under 50ms")
}

// TestTPS_1k_WithChipOutage runs at 1000 TPS, takes Chip down mid-test,
// and verifies events accumulate safely then drain on recovery.
func TestTPS_1k_WithChipOutage(t *testing.T) {
	skipIfExternalChip(t, "inject Unavailable errors on mock server")

	testStart := time.Now()
	t.Logf("TestTPS_1k_WithChipOutage: provisioning...")

	db := directDB(t)
	srv, client := startChipIngressOrMock(t)
	require.NotNil(t, srv)
	store := beholdersvc.NewPgDurableEventStore(db)

	cfg := beholder.DefaultDurableEmitterConfig()
	cfg.RetransmitInterval = 1 * time.Second
	cfg.RetransmitAfter = 2 * time.Second
	cfg.RetransmitBatchSize = 500
	outagePipe := &pipelineDeliveryStats{}
	cfg.Hooks = newPipelineHooks(outagePipe)

	em, err := beholder.NewDurableEmitter(store, client, cfg, logger.Nop())
	require.NoError(t, err)

	ctx := testutils.Context(t)
	em.Start(ctx)
	defer em.Close()

	const targetTPS = 1000
	const concurrency = 20

	// Phase 1: 15s of healthy operation at 1k TPS.
	t.Logf("Phase 1: Healthy — emitting at %d TPS for 15s...", targetTPS)
	p1Start := time.Now()
	phase1Res := runRateLimitedEmit(ctx, t, em, targetTPS, 15*time.Second, concurrency, 256, "outage/phase1_healthy", db, outagePipe)
	phase1Stats := phase1Res.stats
	t.Logf("Phase 1 emit finished in %s", time.Since(p1Start).Round(time.Millisecond))
	time.Sleep(3 * time.Second) // let pipeline drain
	t.Logf("Phase 1 done: %d events emitted (%.0f TPS)", phase1Stats.count(),
		float64(phase1Stats.count())/15.0)

	// Phase 2: Chip goes down. Continue emitting for 15s.
	t.Logf("Phase 2: Chip UNAVAILABLE — emitting at %d TPS for 15s...", targetTPS)
	srv.setPublishErr(status.Error(codes.Unavailable, "chip down"))
	srv.setBatchErr(status.Error(codes.Unavailable, "chip down"))

	p2Start := time.Now()
	phase2Res := runRateLimitedEmit(ctx, t, em, targetTPS, 15*time.Second, concurrency, 256, "outage/phase2_chip_down", db, outagePipe)
	phase2Stats := phase2Res.stats
	t.Logf("Phase 2 emit finished in %s", time.Since(p2Start).Round(time.Millisecond))

	// Queue at end of outage phase (for drain math) + peak sampled during phase 2 emit window.
	queueDuringOutage, _, err := queuePayloadStats(db, ctx)
	require.NoError(t, err)
	t.Logf("Phase 2 done: %d events emitted (%.0f TPS), queue end: %d rows, queue max (emit): %d rows / %s KB*",
		phase2Stats.count(), float64(phase2Stats.count())/15.0,
		queueDuringOutage,
		phase2Res.maxQueueDepth, formatQueueKB(phase2Res.maxQueuePayloadBytes))

	assert.Equal(t, int64(0), phase2Stats.failures.Load(),
		"Emit() must not fail during Chip outage — DB insert should still work")

	// Phase 3: Chip recovers. Stop emitting. Measure drain.
	t.Logf("Phase 3: Chip RECOVERED — measuring drain...")
	srv.setPublishErr(nil)
	srv.setBatchErr(nil)

	drainStart := time.Now()
	require.Eventually(t, func() bool {
		var count int64
		_ = db.QueryRow("SELECT count(*) FROM cre.chip_durable_events").Scan(&count)
		return count == 0
	}, 60*time.Second, 500*time.Millisecond, "queue should drain after Chip recovery")
	drainTime := time.Since(drainStart)
	drainRate := float64(queueDuringOutage) / drainTime.Seconds()

	t.Logf("╔════════════════════════════════════════════════════╗")
	t.Logf("║    1k TPS WITH CHIP OUTAGE — RESULTS              ║")
	t.Logf("╠════════════════════════════════════════════════════╣")
	t.Logf("║ Phase 1 (healthy):                                ║")
	t.Logf("║   Emitted:          %-6d events                  ║", phase1Stats.count())
	t.Logf("║   p99 latency:      %-6.2f ms                     ║", float64(phase1Stats.percentile(0.99).Microseconds())/1000.0)
	t.Logf("║   Pub fail (retry): %-8s                        ║", formatPubFailColumn(phase1Res.ImmPublishFails, phase1Res.BatchPublishFailEvents))
	t.Logf("║   Queue max (emit): %-6d rows / %-8s KB*        ║", phase1Res.maxQueueDepth, formatQueueKB(phase1Res.maxQueuePayloadBytes))
	t.Logf("║ Phase 2 (Chip down):                              ║")
	t.Logf("║   Emitted:          %-6d events                  ║", phase2Stats.count())
	t.Logf("║   p99 latency:      %-6.2f ms                     ║", float64(phase2Stats.percentile(0.99).Microseconds())/1000.0)
	t.Logf("║   Pub fail (retry): %-8s (Publish RPC errors)   ║", formatPubFailColumn(phase2Res.ImmPublishFails, phase2Res.BatchPublishFailEvents))
	t.Logf("║   Emit insert fail: %-6d                         ║", phase2Stats.failures.Load())
	t.Logf("║   Queue max (emit): %-6d rows / %-8s KB*        ║", phase2Res.maxQueueDepth, formatQueueKB(phase2Res.maxQueuePayloadBytes))
	t.Logf("║   Queue end:        %-6d rows                   ║", queueDuringOutage)
	t.Logf("║ Phase 3 (recovery):                               ║")
	t.Logf("║   Drain time:       %-6s                         ║", drainTime.Round(time.Millisecond))
	t.Logf("║   Drain rate:       %-6.0f events/sec              ║", drainRate)
	t.Logf("╚════════════════════════════════════════════════════╝")
	t.Logf("* Queue max KB = sum(octet_length(payload))/1024 for queued rows (excludes index/heap overhead).")
	t.Logf("TestTPS_1k_WithChipOutage finished in %s", time.Since(testStart).Round(time.Millisecond))

	appendTPSummaryBlock("TestTPS_1k_WithChipOutage",
		fmt.Sprintf("total wall clock: %s", time.Since(testStart).Round(time.Millisecond)),
		fmt.Sprintf("phase1 events: %d, phase2 events: %d, queue end: %d rows, phase2 queue max: %d rows / %s KB",
			phase1Stats.count(), phase2Stats.count(), queueDuringOutage, phase2Res.maxQueueDepth, formatQueueKB(phase2Res.maxQueuePayloadBytes)),
		fmt.Sprintf("drain time: %s, drain rate: %.0f ev/s", drainTime.Round(time.Millisecond), drainRate),
	)
}

// TestTPS_PayloadSizeScaling tests 1k TPS at different payload sizes to
// understand how billing record size affects throughput.
func TestTPS_PayloadSizeScaling(t *testing.T) {
	testStart := time.Now()
	sizes := []struct {
		name string
		size int
	}{
		{"64B", 64},
		{"256B", 256},
		{"1KB", 1024},
		{"4KB", 4096},
	}

	tpsPayloadMu.Lock()
	tpsPayloadRows = nil
	tpsPayloadMu.Unlock()

	t.Logf("TestTPS_PayloadSizeScaling: 1k TPS × payload sizes %v", sizes)

	const payloadDuration = 15 * time.Second

	t.Logf("╔════════════════════════════════════════════════════════════════════════════════════════════╗")
	t.Logf("║                 1k TPS × PAYLOAD SIZE SCALING                                              ║")
	t.Logf("╠══════════╦══════════╦═════════════╦══════════╦══════════╦══════════╦══════════╦══════════╦══════════╣")
	t.Logf("║ Payload  ║ Achieved ║ Total emits ║ Emit p50 ║ Emit p99 ║ Pub fail ║ Q max    ║ Q end    ║ Q max    ║")
	t.Logf("║ Size     ║ TPS      ║ (success)   ║ (ms)     ║ (ms)     ║ (retry)* ║ (rows)   ║ (rows)   ║ (KB)*    ║")
	t.Logf("╠══════════╬══════════╬═════════════╬══════════╬══════════╬══════════╬══════════╬══════════╬══════════╣")

	for _, s := range sizes {
		t.Run(s.name, func(t *testing.T) {
			t.Logf(">>> payload %s: provisioning...", s.name)
			db := directDB(t)
			_, client := startChipIngressOrMock(t)
			store := beholdersvc.NewPgDurableEventStore(db)

			cfg := beholder.DefaultDurableEmitterConfig()
			cfg.RetransmitInterval = 1 * time.Second
			cfg.RetransmitAfter = 3 * time.Second
			cfg.RetransmitBatchSize = 500
			pipe := &pipelineDeliveryStats{}
			cfg.Hooks = newPipelineHooks(pipe)

			em, err := beholder.NewDurableEmitter(store, client, cfg, logger.Nop())
			require.NoError(t, err)

			ctx := testutils.Context(t)
			em.Start(ctx)
			defer em.Close()

			const targetTPS = 1000
			const concurrency = 20

			t.Logf(">>> payload %s: emitting %d TPS for %s", s.name, targetTPS, payloadDuration)
			emitRes := runRateLimitedEmit(ctx, t, em, targetTPS, payloadDuration, concurrency, s.size,
				fmt.Sprintf("payload/%s", s.name), db, pipe)
			stats := emitRes.stats

			queueEnd, _, err := queuePayloadStats(db, ctx)
			require.NoError(t, err)

			achieved := float64(stats.count()) / payloadDuration.Seconds()
			totalEmits := stats.count()

			if stats.failures.Load() > 0 {
				t.Logf(">>> payload %s: Emit() insert failures: %d", s.name, stats.failures.Load())
			}
			rowLine := fmt.Sprintf("║ %-8s ║ %-8.0f ║ %-11d ║ %-8.2f ║ %-8.2f ║ %-8s ║ %-8d ║ %-8d ║ %-8s ║",
				s.name, achieved, totalEmits,
				float64(stats.percentile(0.50).Microseconds())/1000.0,
				float64(stats.percentile(0.99).Microseconds())/1000.0,
				formatPubFailColumn(emitRes.ImmPublishFails, emitRes.BatchPublishFailEvents),
				emitRes.maxQueueDepth, queueEnd,
				formatQueueKB(emitRes.maxQueuePayloadBytes))
			t.Log(rowLine)

			tpsPayloadMu.Lock()
			tpsPayloadRows = append(tpsPayloadRows, rowLine)
			tpsPayloadMu.Unlock()
		})
	}

	t.Logf("╚══════════╩══════════╩═════════════╩══════════╩══════════╩══════════╩══════════╩══════════╩══════════╝")
	t.Logf("* Pub fail: failed Publish / PublishBatch (see ramp test footnote). Q max KB* = sum(octet_length(payload))/1024. Total emits = successful Emit() per %s.", payloadDuration)
	t.Logf("TestTPS_PayloadSizeScaling finished in %s", time.Since(testStart).Round(time.Millisecond))

	summaryLines := []string{
		fmt.Sprintf("total wall clock: %s", time.Since(testStart).Round(time.Millisecond)),
		"╔════════════════════════════════════════════════════════════════════════════════════════════╗",
		"║                 1k TPS × PAYLOAD SIZE SCALING                                              ║",
		"╠══════════╦══════════╦═════════════╦══════════╦══════════╦══════════╦══════════╦══════════╦══════════╣",
		"║ Payload  ║ Achieved ║ Total emits ║ Emit p50 ║ Emit p99 ║ Pub fail ║ Q max    ║ Q end    ║ Q max    ║",
		"║ Size     ║ TPS      ║ (success)   ║ (ms)     ║ (ms)     ║ (retry)* ║ (rows)   ║ (rows)   ║ (KB)*    ║",
		"╠══════════╬══════════╬═════════════╬══════════╬══════════╬══════════╬══════════╬══════════╬══════════╣",
	}
	tpsPayloadMu.Lock()
	summaryLines = append(summaryLines, tpsPayloadRows...)
	tpsPayloadMu.Unlock()
	summaryLines = append(summaryLines, "╚══════════╩══════════╩═════════════╩══════════╩══════════╩══════════╩══════════╩══════════╩══════════╝")
	appendTPSummaryBlock("TestTPS_PayloadSizeScaling", summaryLines...)
}
