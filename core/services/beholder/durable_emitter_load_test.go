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
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
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

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/chipingress"
	"github.com/smartcontractkit/chainlink-common/pkg/chipingress/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	beholdersvc "github.com/smartcontractkit/chainlink/v2/core/services/beholder"

	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
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

// loadTestServer is a controllable gRPC ChipIngress server for load tests.
type loadTestServer struct {
	pb.UnimplementedChipIngressServer

	mu           sync.Mutex
	publishErr   error
	batchErr     error
	publishDelay time.Duration

	publishCount atomic.Int64
	batchCount   atomic.Int64
	totalEvents  atomic.Int64
}

func (s *loadTestServer) Publish(_ context.Context, _ *cepb.CloudEvent) (*pb.PublishResponse, error) {
	if s.publishDelay > 0 {
		time.Sleep(s.publishDelay)
	}
	s.publishCount.Add(1)
	s.totalEvents.Add(1)
	s.mu.Lock()
	defer s.mu.Unlock()
	return &pb.PublishResponse{}, s.publishErr
}

func (s *loadTestServer) PublishBatch(_ context.Context, in *pb.CloudEventBatch) (*pb.PublishResponse, error) {
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
	envChipIngressTestAddr       = "CHIP_INGRESS_TEST_ADDR"
	envChipIngressTestTLS        = "CHIP_INGRESS_TEST_TLS"
	envChipIngressTestBasicUser  = "CHIP_INGRESS_TEST_BASIC_AUTH_USER"
	envChipIngressTestBasicPass  = "CHIP_INGRESS_TEST_BASIC_AUTH_PASS"
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
func TestFullStack_SustainedThroughput(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	srv, client := startChipIngressOrMock(t)
	store := beholdersvc.NewPgDurableEventStore(db)

	cfg := beholder.DefaultDurableEmitterConfig()
	cfg.RetransmitInterval = 500 * time.Millisecond
	cfg.RetransmitAfter = 2 * time.Second
	cfg.RetransmitBatchSize = 200
	cfg.PublishTimeout = 5 * time.Second

	em, err := beholder.NewDurableEmitter(store, client, cfg, logger.Test(t))
	require.NoError(t, err)

	ctx := testutils.Context(t)
	em.Start(ctx)
	defer em.Close()

	const (
		totalEvents = 1000
		concurrency = 10
	)

	payload := buildLoadTestPayload(256) // ~256 byte record (protobuf for external Chip)

	start := time.Now()

	var wg sync.WaitGroup
	var emitErrors atomic.Int64
	for w := 0; w < concurrency; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < totalEvents/concurrency; i++ {
				if err := em.Emit(ctx, payload, loadEmitAttrs()...); err != nil {
					emitErrors.Add(1)
				}
			}
		}()
	}
	wg.Wait()
	emitElapsed := time.Since(start)

	t.Logf("--- Emit Phase ---")
	t.Logf("Events emitted:   %d", totalEvents)
	t.Logf("Emit errors:      %d", emitErrors.Load())
	t.Logf("Elapsed:          %s", emitElapsed.Round(time.Millisecond))
	t.Logf("Emit rate:        %.0f events/sec", float64(totalEvents)/emitElapsed.Seconds())

	assert.Equal(t, int64(0), emitErrors.Load(), "all emits should succeed")

	// Wait for all events to be delivered and store to drain.
	require.Eventually(t, func() bool {
		pending, _ := store.ListPending(ctx, time.Now().Add(time.Hour), 1)
		return len(pending) == 0
	}, 30*time.Second, 100*time.Millisecond, "store should drain completely")

	totalElapsed := time.Since(start)

	t.Logf("--- Delivery Phase ---")
	t.Logf("Server received:  %s events (mock only; use external Chip metrics otherwise)", formatMockServerEvents(srv))
	if srv != nil {
		t.Logf("Publish calls:    %d", srv.publishCount.Load())
		t.Logf("Batch calls:      %d", srv.batchCount.Load())
	}
	t.Logf("Total elapsed:    %s", totalElapsed.Round(time.Millisecond))
	t.Logf("End-to-end rate:  %.0f events/sec", float64(totalEvents)/totalElapsed.Seconds())

	if srv != nil {
		assert.GreaterOrEqual(t, srv.totalEvents.Load(), int64(totalEvents),
			"server should have received all events (may have retransmit duplicates)")
	}
}

// TestFullStack_ChipOutage simulates Chip going down during sustained load,
// then recovering. Measures: how events accumulate in Postgres, and how
// fast they drain once Chip comes back.
func TestFullStack_ChipOutage(t *testing.T) {
	skipIfExternalChip(t, "inject Unavailable errors on mock server")

	db := pgtest.NewSqlxDB(t)
	srv, client := startChipIngressOrMock(t)
	require.NotNil(t, srv)
	store := beholdersvc.NewPgDurableEventStore(db)

	cfg := beholder.DefaultDurableEmitterConfig()
	cfg.RetransmitInterval = 200 * time.Millisecond
	cfg.RetransmitAfter = 100 * time.Millisecond
	cfg.RetransmitBatchSize = 100
	cfg.PublishTimeout = 1 * time.Second

	em, err := beholder.NewDurableEmitter(store, client, cfg, logger.Test(t))
	require.NoError(t, err)

	ctx := testutils.Context(t)
	em.Start(ctx)
	defer em.Close()

	// Phase 1: Chip is available — emit 200 events.
	for i := 0; i < 200; i++ {
		require.NoError(t, em.Emit(ctx, []byte("pre-outage"), loadEmitAttrs()...))
	}
	require.Eventually(t, func() bool {
		pending, _ := store.ListPending(ctx, time.Now().Add(time.Hour), 1)
		return len(pending) == 0
	}, 10*time.Second, 50*time.Millisecond, "pre-outage events should all deliver")
	t.Logf("Phase 1: %d events delivered pre-outage", srv.totalEvents.Load())

	// Phase 2: Chip goes down — emit 500 more events.
	srv.setPublishErr(status.Error(codes.Unavailable, "chip down"))
	srv.setBatchErr(status.Error(codes.Unavailable, "chip down"))

	outageStart := time.Now()
	for i := 0; i < 500; i++ {
		require.NoError(t, em.Emit(ctx, []byte("during-outage"), loadEmitAttrs()...))
	}
	t.Logf("Phase 2: emitted 500 events during outage in %s", time.Since(outageStart).Round(time.Millisecond))

	// Verify events are accumulating in Postgres.
	time.Sleep(500 * time.Millisecond) // let some retransmits fail
	pending, err := store.ListPending(ctx, time.Now().Add(time.Hour), 1000)
	require.NoError(t, err)
	t.Logf("Phase 2: %d events pending in Postgres during outage", len(pending))
	assert.Greater(t, len(pending), 0, "events should accumulate during outage")

	// Phase 3: Chip recovers.
	srv.setPublishErr(nil)
	srv.setBatchErr(nil)
	recoveryStart := time.Now()

	require.Eventually(t, func() bool {
		pending, _ := store.ListPending(ctx, time.Now().Add(time.Hour), 1)
		return len(pending) == 0
	}, 30*time.Second, 100*time.Millisecond, "all events should drain after recovery")

	recoveryElapsed := time.Since(recoveryStart)
	t.Logf("Phase 3: drained in %s after recovery (%.0f events/sec drain rate)",
		recoveryElapsed.Round(time.Millisecond),
		float64(500)/recoveryElapsed.Seconds())
	t.Logf("Total server events: %d", srv.totalEvents.Load())
}

// TestFullStack_SlowChip simulates a slow Chip server (high latency per
// publish). This tests whether the async design keeps Emit() fast even
// when gRPC is slow.
func TestFullStack_SlowChip(t *testing.T) {
	skipIfExternalChip(t, "inject publish latency on mock server")

	db := pgtest.NewSqlxDB(t)
	srv, client := startChipIngressOrMock(t)
	require.NotNil(t, srv)
	srv.publishDelay = 50 * time.Millisecond // 50ms per publish = ~20 RPS max
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
	db := pgtest.NewSqlxDB(b)
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
			db := pgtest.NewSqlxDB(b)
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
// txdb serializes all operations through a single transaction, which bottlenecks
// concurrent writes. For TPS testing we need real connection pooling.
func directDB(t testing.TB) *sqlx.DB {
	t.Helper()
	testutils.SkipShortDB(t)
	dbURL := string(env.DatabaseURL.Get())
	if dbURL == "" {
		t.Fatal("CL_DATABASE_URL is required for TPS tests")
	}
	db, err := sqlx.Open("postgres", dbURL)
	require.NoError(t, err)
	require.NoError(t, db.Ping())
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)

	// Clean the table before and after the test.
	_, _ = db.Exec("DELETE FROM cre.chip_durable_events")
	t.Cleanup(func() {
		_, _ = db.Exec("DELETE FROM cre.chip_durable_events")
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

// runRateLimitedEmit emits events at a target rate for the given duration,
// using the specified concurrency. Returns latency stats.
// If progressLabel is non-empty, prints a live progress bar and emit count to stdout every 500ms.
func runRateLimitedEmit(
	ctx context.Context,
	t testing.TB,
	em *beholder.DurableEmitter,
	targetTPS int,
	duration time.Duration,
	concurrency int,
	payloadSize int,
	progressLabel string,
) *emitLatencyStats {
	t.Helper()

	stats := &emitLatencyStats{}
	var emitCount atomic.Int64
	payload := buildLoadTestPayload(payloadSize)

	// Each worker gets an equal share of the target TPS.
	perWorkerTPS := targetTPS / concurrency
	if perWorkerTPS < 1 {
		perWorkerTPS = 1
	}
	interval := time.Duration(float64(time.Second) / float64(perWorkerTPS))

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
	return stats
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

	t.Logf("╔════════════════════════════════════════════════════════════════════════════════════════════╗")
	t.Logf("║                              TPS RAMP-UP TEST RESULTS                                      ║")
	t.Logf("╠═══════════╦══════════╦═════════════╦══════════╦══════════╦══════════╦══════════╦══════════╣")
	t.Logf("║ Target    ║ Achieved ║ Total emits ║ Emit p50 ║ Emit p99 ║ Failures ║ Server   ║ Queue    ║")
	t.Logf("║ TPS       ║ TPS      ║ (success)   ║ (ms)     ║ (ms)     ║          ║ recv*    ║ depth    ║")
	t.Logf("╠═══════════╬══════════╬═════════════╬══════════╬══════════╬══════════╬══════════╬══════════╣")

	for _, targetTPS := range levels {
		t.Run(fmt.Sprintf("%d_tps", targetTPS), func(t *testing.T) {
			levelStart := time.Now()
			t.Logf(">>> level %d TPS: provisioning direct DB + Chip endpoint...", targetTPS)

			db := directDB(t)
			srv, client := startChipIngressOrMock(t)
			store := beholdersvc.NewPgDurableEventStore(db)

			cfg := beholder.DefaultDurableEmitterConfig()
			cfg.RetransmitInterval = 1 * time.Second
			cfg.RetransmitAfter = 3 * time.Second
			cfg.RetransmitBatchSize = 500

			em, err := beholder.NewDurableEmitter(store, client, cfg, logger.Nop())
			require.NoError(t, err)
			ctx := testutils.Context(t)
			em.Start(ctx)
			defer em.Close()

			const duration = 10 * time.Second
			const concurrency = 20

			t.Logf(">>> level %d TPS: emitting for %s @ concurrency=%d (progress bar on stdout)", targetTPS, duration, concurrency)
			stats := runRateLimitedEmit(ctx, t, em, targetTPS, duration, concurrency, 256,
				fmt.Sprintf("ramp_up/%d_tps", targetTPS))
			emitPhase := time.Since(levelStart)
			t.Logf(">>> level %d TPS: emit phase wall time %s", targetTPS, emitPhase.Round(time.Millisecond))

			// Brief pause for async publishes to complete.
			t.Logf(">>> level %d TPS: sleeping 2s for async publishes...", targetTPS)
			time.Sleep(2 * time.Second)

			achieved := float64(stats.count()) / duration.Seconds()
			p50 := stats.percentile(0.50)
			p99 := stats.percentile(0.99)
			serverCol := formatMockServerEvents(srv)

			var queueDepth int64
			row := db.QueryRow("SELECT count(*) FROM cre.chip_durable_events")
			_ = row.Scan(&queueDepth)

			totalEmits := stats.count()
			rowLine := fmt.Sprintf("║ %-9d ║ %-8.0f ║ %-11d ║ %-8.2f ║ %-8.2f ║ %-8d ║ %-8s ║ %-8d ║",
				targetTPS, achieved, totalEmits,
				float64(p50.Microseconds())/1000.0,
				float64(p99.Microseconds())/1000.0,
				stats.failures.Load(),
				serverCol, queueDepth)
			t.Log(rowLine)

			tpsRampMu.Lock()
			tpsRampRows = append(tpsRampRows, rowLine)
			tpsRampMu.Unlock()
		})
	}

	t.Logf("╚═══════════╩══════════╩═════════════╩══════════╩══════════╩══════════╩══════════╩══════════╝")
	t.Logf("* Server recv: in-process mock gRPC publish/batch event count. With CHIP_INGRESS_TEST_ADDR (real Chip), "+
		"this is N/A — observe Kafka/Chip metrics instead. Total emits = successful Emit() completions in the window.")
	t.Logf("TestTPS_RampUp finished in %s", time.Since(testStart).Round(time.Millisecond))

	summaryLines := []string{
		fmt.Sprintf("total wall clock: %s", time.Since(testStart).Round(time.Millisecond)),
		"╔════════════════════════════════════════════════════════════════════════════════════════════╗",
		"║                              TPS RAMP-UP TEST RESULTS                                      ║",
		"╠═══════════╦══════════╦═════════════╦══════════╦══════════╦══════════╦══════════╦══════════╣",
		"║ Target    ║ Achieved ║ Total emits ║ Emit p50 ║ Emit p99 ║ Failures ║ Server   ║ Queue    ║",
		"║ TPS       ║ TPS      ║ (success)   ║ (ms)     ║ (ms)     ║          ║ recv*    ║ depth    ║",
		"╠═══════════╬══════════╬═════════════╬══════════╬══════════╬══════════╬══════════╬══════════╣",
	}
	tpsRampMu.Lock()
	summaryLines = append(summaryLines, tpsRampRows...)
	tpsRampMu.Unlock()
	summaryLines = append(summaryLines, "╚═══════════╩══════════╩═════════════╩══════════╩══════════╩══════════╩══════════╩══════════╝",
		"* Server recv: mock-only; N/A with real Chip. Total emits = successful Emit() calls per level.")
	appendTPSummaryBlock("TestTPS_RampUp", summaryLines...)
}

// TestTPS_Sustained1k runs at exactly 1000 TPS for 60 seconds and verifies
// the pipeline keeps up: deletes match inserts, queue stays bounded, and
// Emit() latency stays low.
func TestTPS_Sustained1k(t *testing.T) {
	testStart := time.Now()
	t.Logf("TestTPS_Sustained1k: provisioning DB + Chip server + emitter...")

	db := directDB(t)
	srv, client := startChipIngressOrMock(t)
	store := beholdersvc.NewPgDurableEventStore(db)

	cfg := beholder.DefaultDurableEmitterConfig()
	cfg.RetransmitInterval = 1 * time.Second
	cfg.RetransmitAfter = 3 * time.Second
	cfg.RetransmitBatchSize = 500

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

	stats := runRateLimitedEmit(ctx, t, em, targetTPS, duration, concurrency, 256, "sustained_1k")

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

	t.Logf("╔════════════════════════════════════════════════════╗")
	t.Logf("║       SUSTAINED 1k TPS TEST RESULTS               ║")
	t.Logf("╠════════════════════════════════════════════════════╣")
	t.Logf("║ Target TPS:       %-6d                          ║", targetTPS)
	t.Logf("║ Duration:         %-6s                          ║", duration)
	t.Logf("║ Total emitted:    %-6d                          ║", stats.count())
	t.Logf("║ Achieved TPS:     %-6.0f                          ║", achievedTPS)
	t.Logf("║ Emit failures:    %-6d                          ║", stats.failures.Load())
	t.Logf("║ Emit p50 latency: %-6.2f ms                      ║", float64(stats.percentile(0.50).Microseconds())/1000.0)
	t.Logf("║ Emit p99 latency: %-6.2f ms                      ║", float64(stats.percentile(0.99).Microseconds())/1000.0)
	t.Logf("║ Server received:  %-6s (mock event count)      ║", formatMockServerEvents(srv))
	t.Logf("║ Drain time:       %-6s                          ║", drainTime.Round(time.Millisecond))
	t.Logf("╚════════════════════════════════════════════════════╝")
	t.Logf("TestTPS_Sustained1k finished in %s", time.Since(testStart).Round(time.Millisecond))

	appendTPSummaryBlock("TestTPS_Sustained1k",
		fmt.Sprintf("total wall clock: %s", time.Since(testStart).Round(time.Millisecond)),
		fmt.Sprintf("emit phase: %s", time.Since(emitStart).Round(time.Millisecond)),
		fmt.Sprintf("target TPS: %d, achieved: %.0f, failures: %d", targetTPS, achievedTPS, stats.failures.Load()),
		fmt.Sprintf("emit p50/p99 ms: %.2f / %.2f", float64(stats.percentile(0.50).Microseconds())/1000.0, float64(stats.percentile(0.99).Microseconds())/1000.0),
		fmt.Sprintf("server events: %s, drain time: %s", formatMockServerEvents(srv), drainTime.Round(time.Millisecond)),
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
	phase1Stats := runRateLimitedEmit(ctx, t, em, targetTPS, 15*time.Second, concurrency, 256, "outage/phase1_healthy")
	t.Logf("Phase 1 emit finished in %s", time.Since(p1Start).Round(time.Millisecond))
	time.Sleep(3 * time.Second) // let pipeline drain
	t.Logf("Phase 1 done: %d events emitted (%.0f TPS)", phase1Stats.count(),
		float64(phase1Stats.count())/15.0)

	// Phase 2: Chip goes down. Continue emitting for 15s.
	t.Logf("Phase 2: Chip UNAVAILABLE — emitting at %d TPS for 15s...", targetTPS)
	srv.setPublishErr(status.Error(codes.Unavailable, "chip down"))
	srv.setBatchErr(status.Error(codes.Unavailable, "chip down"))

	p2Start := time.Now()
	phase2Stats := runRateLimitedEmit(ctx, t, em, targetTPS, 15*time.Second, concurrency, 256, "outage/phase2_chip_down")
	t.Logf("Phase 2 emit finished in %s", time.Since(p2Start).Round(time.Millisecond))

	// Check queue depth during outage.
	var queueDuringOutage int64
	_ = db.QueryRow("SELECT count(*) FROM cre.chip_durable_events").Scan(&queueDuringOutage)
	t.Logf("Phase 2 done: %d events emitted (%.0f TPS), queue depth: %d",
		phase2Stats.count(), float64(phase2Stats.count())/15.0, queueDuringOutage)

	assert.Equal(t, int64(0), phase2Stats.failures.Load(),
		"Emit must not fail during Chip outage — DB insert should still work")

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
	t.Logf("║ Phase 2 (Chip down):                              ║")
	t.Logf("║   Emitted:          %-6d events                  ║", phase2Stats.count())
	t.Logf("║   p99 latency:      %-6.2f ms                     ║", float64(phase2Stats.percentile(0.99).Microseconds())/1000.0)
	t.Logf("║   Emit failures:    %-6d                         ║", phase2Stats.failures.Load())
	t.Logf("║   Queue depth:      %-6d events                  ║", queueDuringOutage)
	t.Logf("║ Phase 3 (recovery):                               ║")
	t.Logf("║   Drain time:       %-6s                         ║", drainTime.Round(time.Millisecond))
	t.Logf("║   Drain rate:       %-6.0f events/sec              ║", drainRate)
	t.Logf("║   Server received:  %-6d total                   ║", srv.totalEvents.Load())
	t.Logf("╚════════════════════════════════════════════════════╝")
	t.Logf("TestTPS_1k_WithChipOutage finished in %s", time.Since(testStart).Round(time.Millisecond))

	appendTPSummaryBlock("TestTPS_1k_WithChipOutage",
		fmt.Sprintf("total wall clock: %s", time.Since(testStart).Round(time.Millisecond)),
		fmt.Sprintf("phase1 events: %d, phase2 events: %d, queue at outage: %d", phase1Stats.count(), phase2Stats.count(), queueDuringOutage),
		fmt.Sprintf("drain time: %s, drain rate: %.0f ev/s, server total: %d", drainTime.Round(time.Millisecond), drainRate, srv.totalEvents.Load()),
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

	t.Logf("╔══════════════════════════════════════════════════════════════════════════╗")
	t.Logf("║                 1k TPS × PAYLOAD SIZE SCALING                            ║")
	t.Logf("╠══════════╦══════════╦═════════════╦══════════╦══════════╦════════════════╣")
	t.Logf("║ Payload  ║ Achieved ║ Total emits ║ Emit p50 ║ Emit p99 ║ Failures      ║")
	t.Logf("║ Size     ║ TPS      ║ (success)   ║ (ms)     ║ (ms)     ║               ║")
	t.Logf("╠══════════╬══════════╬═════════════╬══════════╬══════════╬════════════════╣")

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

			em, err := beholder.NewDurableEmitter(store, client, cfg, logger.Nop())
			require.NoError(t, err)

			ctx := testutils.Context(t)
			em.Start(ctx)
			defer em.Close()

			const targetTPS = 1000
			const concurrency = 20

			t.Logf(">>> payload %s: emitting %d TPS for %s", s.name, targetTPS, payloadDuration)
			stats := runRateLimitedEmit(ctx, t, em, targetTPS, payloadDuration, concurrency, s.size,
				fmt.Sprintf("payload/%s", s.name))

			achieved := float64(stats.count()) / payloadDuration.Seconds()
			totalEmits := stats.count()

			rowLine := fmt.Sprintf("║ %-8s ║ %-8.0f ║ %-11d ║ %-8.2f ║ %-8.2f ║ %-14d ║",
				s.name, achieved, totalEmits,
				float64(stats.percentile(0.50).Microseconds())/1000.0,
				float64(stats.percentile(0.99).Microseconds())/1000.0,
				stats.failures.Load())
			t.Log(rowLine)

			tpsPayloadMu.Lock()
			tpsPayloadRows = append(tpsPayloadRows, rowLine)
			tpsPayloadMu.Unlock()
		})
	}

	t.Logf("╚══════════╩══════════╩═════════════╩══════════╩══════════╩════════════════╝")
	t.Logf("Total emits = successful Emit() calls in each %s window (per payload size).", payloadDuration)
	t.Logf("TestTPS_PayloadSizeScaling finished in %s", time.Since(testStart).Round(time.Millisecond))

	summaryLines := []string{
		fmt.Sprintf("total wall clock: %s", time.Since(testStart).Round(time.Millisecond)),
		"╔══════════════════════════════════════════════════════════════════════════╗",
		"║                 1k TPS × PAYLOAD SIZE SCALING                            ║",
		"╠══════════╦══════════╦═════════════╦══════════╦══════════╦════════════════╣",
		"║ Payload  ║ Achieved ║ Total emits ║ Emit p50 ║ Emit p99 ║ Failures      ║",
		"║ Size     ║ TPS      ║ (success)   ║ (ms)     ║ (ms)     ║               ║",
		"╠══════════╬══════════╬═════════════╬══════════╬══════════╬════════════════╣",
	}
	tpsPayloadMu.Lock()
	summaryLines = append(summaryLines, tpsPayloadRows...)
	tpsPayloadMu.Unlock()
	summaryLines = append(summaryLines, "╚══════════╩══════════╩═════════════╩══════════╩══════════╩════════════════╝")
	appendTPSummaryBlock("TestTPS_PayloadSizeScaling", summaryLines...)
}
