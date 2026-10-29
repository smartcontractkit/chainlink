package telemetry

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"

	cepb "github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/chipingress/batch"
	"github.com/smartcontractkit/chainlink-common/pkg/chipingress/pb"

	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
)

// receivedEvent is one event as chip ingress saw it.
type receivedEvent struct {
	method string
	domain string // the chainlink-domain header, "" if absent
	event  *cepb.CloudEvent
}

// capturingChipServer records PublishBatch / PublishBatchNoAck calls and answers each event with
// a positional success, as chip ingress does. With hold set, it answers nothing until released.
type capturingChipServer struct {
	pb.UnimplementedChipIngressServer
	hold chan struct{}

	mu     sync.Mutex
	events []receivedEvent
}

func (s *capturingChipServer) record(ctx context.Context, method string, in *pb.CloudEventBatch) (*pb.PublishResponse, error) {
	if s.hold != nil {
		select {
		case <-s.hold:
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	md, _ := metadata.FromIncomingContext(ctx)
	domain := ""
	if v := md.Get(beholder.DomainHeaderKey); len(v) > 0 {
		domain = v[0]
	}
	s.mu.Lock()
	for _, ev := range in.Events {
		s.events = append(s.events, receivedEvent{method: method, domain: domain, event: ev})
	}
	s.mu.Unlock()
	return &pb.PublishResponse{Results: make([]*pb.PublishResult, len(in.Events))}, nil
}

func (s *capturingChipServer) PublishBatch(ctx context.Context, in *pb.CloudEventBatch) (*pb.PublishResponse, error) {
	return s.record(ctx, "PublishBatch", in)
}

func (s *capturingChipServer) PublishBatchNoAck(ctx context.Context, in *pb.CloudEventBatch) (*pb.PublishResponse, error) {
	return s.record(ctx, "PublishBatchNoAck", in)
}

func (s *capturingChipServer) received() []receivedEvent {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]receivedEvent(nil), s.events...)
}

type nopLogExporter struct{}

func (nopLogExporter) Export(context.Context, []sdklog.Record) error { return nil }
func (nopLogExporter) Shutdown(context.Context) error                { return nil }
func (nopLogExporter) ForceFlush(context.Context) error              { return nil }

// startChipServer serves chip on a local listener and returns its address.
func startChipServer(t *testing.T, chip *capturingChipServer) string {
	t.Helper()
	lis, err := (&net.ListenConfig{}).Listen(t.Context(), "tcp", "127.0.0.1:0")
	require.NoError(t, err)
	srv := grpc.NewServer()
	pb.RegisterChipIngressServer(srv, chip)
	go func() { _ = srv.Serve(lis) }()
	t.Cleanup(srv.Stop)
	return lis.Addr().String()
}

// startBeholderClient starts a real Beholder client with domain emitters against addr. Its OTel
// exporters also point at addr, so shutdown gets a prompt answer instead of dialing an absent
// collector. closeClient closes it once and returns Close's error; cleanup closes it otherwise.
func startBeholderClient(t *testing.T, addr string, tune func(*beholder.Config)) (client *beholder.Client, closeClient func() error) {
	t.Helper()
	cfg := beholder.Config{
		OtelExporterGRPCEndpoint:         addr,
		InsecureConnection:               true,
		ChipIngressEmitterEnabled:        true,
		ChipIngressEmitterGRPCEndpoint:   addr,
		ChipIngressInsecureConnection:    true,
		ChipIngressSendInterval:          10 * time.Millisecond,
		ChipIngressDomainEmittersEnabled: true,
		AuthHeaders:                      map[string]string{"X-Beholder-Node-Auth-Token": "1:abc:2:def"},
	}
	if tune != nil {
		tune(&cfg)
	}
	client, err := beholder.NewGRPCClient(cfg, func(...otlploggrpc.Option) (sdklog.Exporter, error) { return nopLogExporter{}, nil })
	require.NoError(t, err)
	require.NoError(t, client.Start(t.Context()))

	var once sync.Once
	var closeErr error
	closeClient = func() error {
		once.Do(func() { closeErr = client.Close() })
		return closeErr
	}
	t.Cleanup(func() { _ = closeClient() })
	return client, closeClient
}

// requireCleanClose asserts that closing the client did not fail because of the domain emitters.
// OTel provider shutdown may still report export errors (the test server is not a collector).
func requireCleanClose(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		assert.NotContains(t, err.Error(), "client connection is closing", "domain emitter connections close once")
	}
}

func TestDomainManager_ThroughRealBeholderClient(t *testing.T) {
	t.Parallel()
	chip := &capturingChipServer{}
	addr := startChipServer(t, chip)
	client, closeClient := startBeholderClient(t, addr, nil)
	m, reader := newTestDomainManager(t, client)

	// Concurrent senders across two domains, through single-type and multitype endpoints.
	const workers, perWorker = 4, 25
	var wg sync.WaitGroup
	for range workers {
		wg.Go(func() {
			feeds := m.GenMonitoringEndpoint("EVM", ethMainnetChainID, "0xabc", synchronization.OCR2Median)
			streams := m.GenMultitypeMonitoringEndpoint("evm", ethMainnetChainID, "0xdef/1")
			for range perWorker {
				feeds.SendLog([]byte("p"))
				streams.SendTypedLog(synchronization.LLOOutcome, []byte("p"))
			}
		})
	}
	wg.Wait()

	const total = 2 * workers * perWorker
	require.Eventually(t, func() bool { return counterValue(t, reader, metricMessagesSent) == total },
		10*time.Second, 20*time.Millisecond, "every event is acknowledged by the per-event result")

	received := chip.received()
	require.Len(t, received, total)
	for _, r := range received {
		assert.Equal(t, "PublishBatchNoAck", r.method)
		switch r.event.Source {
		case "data-feeds.telemetry.ocr2-median":
			assert.Equal(t, string(beholder.DomainDataFeeds), r.domain)
		case "data-streams.telemetry.llo-outcome":
			assert.Equal(t, string(beholder.DomainDataStreams), r.domain)
		default:
			t.Errorf("unexpected source %q", r.event.Source)
		}
		assert.Equal(t, "true", r.event.Attributes["legacytelemetry"].GetCeString())
		assert.Equal(t, testCSAPubKeyHex, r.event.Attributes["csapublickey"].GetCeString())
	}

	requireCleanClose(t, closeClient())
	// Nothing is counted twice, and nothing else is counted, once the emitters have drained.
	assert.Equal(t, int64(total), counterValue(t, reader, metricMessagesSent))
	assert.Zero(t, counterValue(t, reader, metricMessagesSendErrors))
	assert.Zero(t, counterValue(t, reader, metricMessagesDropped))
}

func TestDomainManager_ThroughRealBeholderClient_BufferFull(t *testing.T) {
	t.Parallel()
	chip := &capturingChipServer{hold: make(chan struct{})}
	addr := startChipServer(t, chip)
	client, closeClient := startBeholderClient(t, addr, func(c *beholder.Config) {
		c.ChipIngressBufferSize = 3
		c.ChipIngressMaxBatchSize = 1
		c.ChipIngressMaxConcurrentSends = 1
	})
	m, reader := newTestDomainManager(t, client)

	// The server holds every request, so the tiny buffer fills and the rest are dropped.
	const total = 50
	ep := m.GenMonitoringEndpoint("EVM", ethMainnetChainID, "0xabc", synchronization.OCR2Median)
	for range total {
		ep.SendLog([]byte("p"))
	}
	bufferFull := attribute.String("error_code", batch.ErrorCodeFor(batch.ErrMessageBufferFull))
	require.Positive(t, counterValue(t, reader, metricMessagesDropped, bufferFull), "drops are counted as dropped, not send errors")

	close(chip.hold) // let the held and queued events through, then drain on close
	requireCleanClose(t, closeClient())

	sent := counterValue(t, reader, metricMessagesSent)
	dropped := counterValue(t, reader, metricMessagesDropped)
	assert.Equal(t, int64(total), sent+dropped, "every event is counted exactly once")
	assert.Equal(t, sent, int64(len(chip.received())), "sent matches what the server received")
	assert.Zero(t, counterValue(t, reader, metricMessagesSendErrors))
}
