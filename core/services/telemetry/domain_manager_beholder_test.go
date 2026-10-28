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
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/chipingress/pb"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
)

type capturedBatch struct {
	method string
	md     metadata.MD
	events []*cepb.CloudEvent
}

// capturingChipServer records PublishBatch / PublishBatchNoAck calls and answers each event
// with a positional success, as chip ingress does.
type capturingChipServer struct {
	pb.UnimplementedChipIngressServer

	mu      sync.Mutex
	batches []capturedBatch
}

func (s *capturingChipServer) record(ctx context.Context, method string, in *pb.CloudEventBatch) (*pb.PublishResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	s.mu.Lock()
	s.batches = append(s.batches, capturedBatch{method: method, md: md, events: in.Events})
	s.mu.Unlock()
	return &pb.PublishResponse{Results: make([]*pb.PublishResult, len(in.Events))}, nil
}

func (s *capturingChipServer) PublishBatch(ctx context.Context, in *pb.CloudEventBatch) (*pb.PublishResponse, error) {
	return s.record(ctx, "PublishBatch", in)
}

func (s *capturingChipServer) PublishBatchNoAck(ctx context.Context, in *pb.CloudEventBatch) (*pb.PublishResponse, error) {
	return s.record(ctx, "PublishBatchNoAck", in)
}

// received returns one entry per event: the RPC method, its chainlink-domain header and the event.
func (s *capturingChipServer) received() (methods, domains []string, events []*cepb.CloudEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, b := range s.batches {
		for _, ev := range b.events {
			methods = append(methods, b.method)
			domains = append(domains, b.md.Get(beholder.DomainHeaderKey)...)
			events = append(events, ev)
		}
	}
	return methods, domains, events
}

type nopLogExporter struct{}

func (nopLogExporter) Export(context.Context, []sdklog.Record) error { return nil }
func (nopLogExporter) Shutdown(context.Context) error                { return nil }
func (nopLogExporter) ForceFlush(context.Context) error              { return nil }

func TestDomainManager_ThroughRealBeholderClient(t *testing.T) {
	t.Parallel()
	lis, err := (&net.ListenConfig{}).Listen(t.Context(), "tcp", "127.0.0.1:0")
	require.NoError(t, err)
	srv := grpc.NewServer()
	chip := &capturingChipServer{}
	pb.RegisterChipIngressServer(srv, chip)
	go func() { _ = srv.Serve(lis) }()
	t.Cleanup(srv.Stop)

	client, err := beholder.NewGRPCClient(beholder.Config{
		OtelExporterGRPCEndpoint:         "localhost:4317", // never reached: logs go to nopLogExporter
		InsecureConnection:               true,
		ChipIngressEmitterEnabled:        true,
		ChipIngressEmitterGRPCEndpoint:   lis.Addr().String(),
		ChipIngressInsecureConnection:    true,
		ChipIngressSendInterval:          10 * time.Millisecond,
		ChipIngressDomainEmittersEnabled: true,
		AuthHeaders:                      map[string]string{"X-Beholder-Node-Auth-Token": "1:abc:2:def"},
	}, func(...otlploggrpc.Option) (sdklog.Exporter, error) { return nopLogExporter{}, nil })
	require.NoError(t, err)
	require.NoError(t, client.Start(t.Context()))

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

	methods, domains, events := chip.received()
	require.Len(t, events, total)
	for i, ev := range events {
		assert.Equal(t, "PublishBatchNoAck", methods[i])
		switch ev.Source {
		case "data-feeds.telemetry.ocr2-median":
			assert.Equal(t, string(beholder.DomainDataFeeds), domains[i])
		case "data-streams.telemetry.llo-outcome":
			assert.Equal(t, string(beholder.DomainDataStreams), domains[i])
		default:
			t.Fatalf("unexpected source %q", ev.Source)
		}
		assert.Equal(t, "true", ev.Attributes["legacytelemetry"].GetCeString())
		assert.Equal(t, testCSAPubKeyHex, ev.Attributes["csapublickey"].GetCeString())
	}

	// Closing the client drains and closes the domain emitters exactly once. OTel provider
	// shutdown may fail to flush to the absent collector; nothing else may fail.
	if err := client.Close(); err != nil {
		assert.NotContains(t, err.Error(), "client connection is closing")
	}
	assert.Zero(t, counterValue(t, reader, metricMessagesSendErrors))
	assert.Zero(t, counterValue(t, reader, metricMessagesDropped))
}
