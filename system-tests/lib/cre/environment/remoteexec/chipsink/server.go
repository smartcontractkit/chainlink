package chipsink

// NOTE: This implementation intentionally mirrors the test helper sink from
// `system-tests/tests/test-helpers/chip-testsink`.
// We keep this copy under `system-tests/lib` so runtime code (agent/CLI) can
// depend on it without importing from test-only packages.
// If we later move the sink to a shared package, both callers should use that
// single canonical location.

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	chippb "github.com/smartcontractkit/chainlink-common/pkg/chipingress/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const listenerReadyTimeout = 5 * time.Second

type PublishFn func(ctx context.Context, event *pb.CloudEvent) (*chippb.PublishResponse, error)

type Config struct {
	GRPCListen       string
	UpstreamEndpoint string
	PublishFn        PublishFn
	Started          chan<- string
}

type Server struct {
	cfg Config

	grpcServer *grpc.Server
	upstream   chippb.ChipIngressClient
	onceStop   sync.Once

	chippb.UnimplementedChipIngressServer
}

func NewServer(cfg Config) (*Server, error) {
	s := &Server{cfg: cfg}
	s.grpcServer = grpc.NewServer()
	chippb.RegisterChipIngressServer(s.grpcServer, s)

	if cfg.UpstreamEndpoint != "" {
		conn, err := grpc.NewClient(cfg.UpstreamEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, fmt.Errorf("dial upstream chip ingress: %w", err)
		}
		s.upstream = chippb.NewChipIngressClient(conn)
	}

	return s, nil
}

func (s *Server) Run() error {
	lc := net.ListenConfig{}
	lis, err := lc.Listen(context.Background(), "tcp", s.cfg.GRPCListen)
	if err != nil {
		return fmt.Errorf("gRPC listen: %w", err)
	}
	addr := lis.Addr().String()

	errCh := make(chan error, 1)
	go func() {
		errCh <- s.grpcServer.Serve(lis)
	}()
	if err := waitForListenerReady(addr, listenerReadyTimeout); err != nil {
		s.grpcServer.Stop()
		return err
	}
	notifyStarted(s.cfg.Started, addr)

	return <-errCh
}

func (s *Server) Shutdown(context.Context) {
	s.onceStop.Do(func() {
		s.grpcServer.GracefulStop()
	})
}

func (s *Server) Publish(ctx context.Context, event *pb.CloudEvent) (*chippb.PublishResponse, error) {
	if s.cfg.UpstreamEndpoint != "" && s.upstream != nil {
		go func() {
			upstreamCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			_, _ = s.upstream.Publish(upstreamCtx, event)
		}()
	}

	if s.cfg.PublishFn != nil {
		return s.cfg.PublishFn(ctx, event)
	}
	return &chippb.PublishResponse{}, nil
}

func waitForListenerReady(addr string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	var lastErr error
	for time.Now().Before(deadline) {
		dialer := &net.Dialer{Timeout: 250 * time.Millisecond}
		conn, err := dialer.Dial("tcp", addr)
		if err == nil {
			_ = conn.Close()
			return nil
		}
		lastErr = err
		time.Sleep(50 * time.Millisecond)
	}
	if lastErr == nil {
		lastErr = fmt.Errorf("listener on %s not ready", addr)
	}
	return fmt.Errorf("timeout waiting for listener readiness: %w", lastErr)
}

func notifyStarted(ch chan<- string, addr string) {
	if ch == nil {
		return
	}
	select {
	case ch <- addr:
	default:
	}
}
