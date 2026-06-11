// Package mockchip provides an in-process mock of the Chip Ingress gRPC
// service. It is intended for local development and integration testing of
// DurableEmitter: it accepts CloudEvents over gRPC, captures them in memory
// for inspection, and can be toggled into an "outage" mode where every
// publish RPC fails. Once the outage is cleared, DurableEmitter's retransmit
// loop drains its persisted queue back into the mock and the captured
// events can be verified.
package mockchip

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	cepb "github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"

	chippb "github.com/smartcontractkit/chainlink-common/pkg/chipingress/pb"
)

const (
	// maxRecvMsgSize matches the chipingress client default (16 MiB) so that
	// batches sent by DurableEmitter's BatchEmitter are not rejected.
	maxRecvMsgSize = 16 * 1024 * 1024
)

// CapturedEvent is a snapshot of a single CloudEvent observed by the mock
// server, with the wall-clock receive time attached for drain-latency
// assertions in tests.
type CapturedEvent struct {
	ReceivedAt time.Time
	Event      *cepb.CloudEvent
}

// Stats summarizes server state for HTTP inspection and tests.
type Stats struct {
	Captured        int   `json:"captured"`
	PublishCalls    int64 `json:"publish_calls"`
	BatchCalls      int64 `json:"batch_calls"`
	FailedCalls     int64 `json:"failed_calls"`
	OutageActive    bool  `json:"outage_active"`
	OutageDurations int64 `json:"outage_count"`
}

// Server is a mock ChipIngress endpoint. It implements chippb.ChipIngressServer
// and is safe for concurrent use.
type Server struct {
	chippb.UnimplementedChipIngressServer

	grpcServer *grpc.Server
	lis        net.Listener

	mu       sync.Mutex
	captured []CapturedEvent

	outage           atomic.Bool
	publishCalls     atomic.Int64
	batchCalls       atomic.Int64
	failedCalls      atomic.Int64
	outageActivation atomic.Int64
}

// NewServer constructs a Server. The server is not started — call Start.
func NewServer() *Server {
	s := &Server{}
	s.grpcServer = grpc.NewServer(
		grpc.MaxRecvMsgSize(maxRecvMsgSize),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    20 * time.Second,
			Timeout: 5 * time.Second,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             5 * time.Second,
			PermitWithoutStream: true,
		}),
	)
	chippb.RegisterChipIngressServer(s.grpcServer, s)
	return s
}

// Start binds the gRPC listener on listenAddr (e.g. ":9095" or "127.0.0.1:0")
// and serves in a goroutine. It returns the resolved listen address, or an
// error if binding failed.
func (s *Server) Start(listenAddr string) (string, error) {
	lc := &net.ListenConfig{}
	lis, err := lc.Listen(context.Background(), "tcp", listenAddr)
	if err != nil {
		return "", fmt.Errorf("mockchip: listen %s: %w", listenAddr, err)
	}
	s.lis = lis
	go func() {
		if serveErr := s.grpcServer.Serve(lis); serveErr != nil && !errors.Is(serveErr, grpc.ErrServerStopped) {
			// gRPC's Serve returns nil on graceful stop; non-nil errors that
			// are not ErrServerStopped indicate a real failure. The mock is a
			// dev tool so we surface via stderr rather than a logger.
			fmt.Printf("mockchip: gRPC Serve returned: %v\n", serveErr)
		}
	}()
	return lis.Addr().String(), nil
}

// Stop gracefully stops the gRPC server.
func (s *Server) Stop() {
	s.grpcServer.GracefulStop()
}

// Addr returns the resolved listen address (host:port). Returns the empty
// string before Start.
func (s *Server) Addr() string {
	if s.lis == nil {
		return ""
	}
	return s.lis.Addr().String()
}

// SetOutage toggles failure injection. When true, every Publish/PublishBatch
// RPC returns codes.Unavailable. Captured events are preserved across
// transitions so that the drain after restore can be verified.
func (s *Server) SetOutage(on bool) {
	prev := s.outage.Swap(on)
	if on && !prev {
		s.outageActivation.Add(1)
	}
}

// OutageActive reports the current outage flag.
func (s *Server) OutageActive() bool { return s.outage.Load() }

// Captured returns a defensive copy of every event observed so far, oldest
// first.
func (s *Server) Captured() []CapturedEvent {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]CapturedEvent, len(s.captured))
	copy(out, s.captured)
	return out
}

// CapturedCount returns the number of events currently held.
func (s *Server) CapturedCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.captured)
}

// Reset clears captured events and counters but leaves outage state alone.
func (s *Server) Reset() {
	s.mu.Lock()
	s.captured = nil
	s.mu.Unlock()
	s.publishCalls.Store(0)
	s.batchCalls.Store(0)
	s.failedCalls.Store(0)
}

// Stats returns a snapshot of counters.
func (s *Server) Stats() Stats {
	return Stats{
		Captured:        s.CapturedCount(),
		PublishCalls:    s.publishCalls.Load(),
		BatchCalls:      s.batchCalls.Load(),
		FailedCalls:     s.failedCalls.Load(),
		OutageActive:    s.outage.Load(),
		OutageDurations: s.outageActivation.Load(),
	}
}

// WaitFor blocks until the captured count reaches n or ctx expires. It is a
// test helper for asserting that DurableEmitter has drained its queue.
func (s *Server) WaitFor(ctx context.Context, n int) error {
	for {
		if s.CapturedCount() >= n {
			return nil
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("mockchip: timed out waiting for %d events (have %d): %w",
				n, s.CapturedCount(), ctx.Err())
		case <-time.After(25 * time.Millisecond):
		}
	}
}

// Publish implements chippb.ChipIngressServer.Publish.
func (s *Server) Publish(_ context.Context, event *cepb.CloudEvent) (*chippb.PublishResponse, error) {
	s.publishCalls.Add(1)
	if s.outage.Load() {
		s.failedCalls.Add(1)
		return nil, status.Error(codes.Unavailable, "mockchip: outage active")
	}
	s.appendEvent(event)
	return &chippb.PublishResponse{
		Results: []*chippb.PublishResult{{EventId: event.GetId()}},
	}, nil
}

// PublishBatch implements chippb.ChipIngressServer.PublishBatch.
func (s *Server) PublishBatch(_ context.Context, batch *chippb.CloudEventBatch) (*chippb.PublishResponse, error) {
	s.batchCalls.Add(1)
	if s.outage.Load() {
		s.failedCalls.Add(1)
		return nil, status.Error(codes.Unavailable, "mockchip: outage active")
	}
	if batch == nil {
		return &chippb.PublishResponse{}, nil
	}
	results := make([]*chippb.PublishResult, 0, len(batch.Events))
	for _, ev := range batch.Events {
		s.appendEvent(ev)
		results = append(results, &chippb.PublishResult{EventId: ev.GetId()})
	}
	return &chippb.PublishResponse{Results: results}, nil
}

// Ping implements chippb.ChipIngressServer.Ping.
func (s *Server) Ping(_ context.Context, _ *chippb.EmptyRequest) (*chippb.PingResponse, error) {
	if s.outage.Load() {
		return nil, status.Error(codes.Unavailable, "mockchip: outage active")
	}
	return &chippb.PingResponse{Message: "pong"}, nil
}

func (s *Server) appendEvent(event *cepb.CloudEvent) {
	if event == nil {
		return
	}
	s.mu.Lock()
	s.captured = append(s.captured, CapturedEvent{
		ReceivedAt: time.Now().UTC(),
		Event:      event,
	})
	s.mu.Unlock()
}
