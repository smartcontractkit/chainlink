package testsinkminimal

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"

	chippb "github.com/smartcontractkit/chainlink-common/pkg/chipingress/pb"
)

type PublishFn = func(ctx context.Context, event *pb.CloudEvent) (*chippb.PublishResponse, error)

// Config defines how the test sink listens.
type Config struct {
	// gRPC listen address for ChipIngress, e.g. ":9090".
	GRPCListen string

	// Optional upstream Chip Ingress endpoint to forward to.
	// If empty, no pass-through is performed.
	UpstreamEndpoint string

	PublishFunc PublishFn

	// Started optionally receives a signal once the gRPC listener is bound.
	Started chan<- struct{}
}

// Server implements the ChipIngress gRPC service + a tiny HTTP API.
type Server struct {
	cfg Config

	grpcServer *grpc.Server

	chippb.UnimplementedChipIngressServer

	// Optional pass-through client.
	upstream chippb.ChipIngressClient
}

// NewServer constructs a new test sink.
func NewServer(cfg Config) (*Server, error) {
	s := &Server{
		cfg: cfg,
	}

	// gRPC server
	s.grpcServer = grpc.NewServer()

	// Register ChipIngress service implementation on this server.
	chippb.RegisterChipIngressServer(s.grpcServer, s)

	// Optional upstream (pass-through)
	if cfg.UpstreamEndpoint != "" {
		conn, err := grpc.NewClient(cfg.UpstreamEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, fmt.Errorf("dial upstream chip ingress: %w", err)
		}
		s.upstream = chippb.NewChipIngressClient(conn)
	}

	return s, nil
}

// Run starts the gRPC server and blocks until it exits.
func (s *Server) Run() error {
	lc := &net.ListenConfig{}
	lis, err := lc.Listen(context.Background(), "tcp", s.cfg.GRPCListen)
	if err != nil {
		return fmt.Errorf("gRPC listen: %w", err)
	}

	notifyStarted(s.cfg.Started)
	log.Printf("[chip-testsink] gRPC listening on %s", s.cfg.GRPCListen)

	errCh := make(chan error, 1)
	go func() {
		errCh <- s.grpcServer.Serve(lis)
	}()

	if s.cfg.UpstreamEndpoint != "" {
		log.Printf("[chip-testsink] Forwarding to upstream Chip Ingress endpoint: %s", s.cfg.UpstreamEndpoint)
	}

	// Wait for first error.
	return <-errCh
}

//
// ===== gRPC: ChipIngressServer implementation =====
//

// Publish implements chippb.ChipIngressServer.Publish.
//
// Adjust the signature if your generated interface differs.
func (s *Server) Publish(ctx context.Context, event *pb.CloudEvent) (*chippb.PublishResponse, error) {
	go func() {
		if s.cfg.UpstreamEndpoint != "" {
			context, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancelFn()
			_, err := s.upstream.Publish(context, event)
			if err != nil {
				log.Printf("failed to forward to upstream: %v", err)
			}
		}
	}()

	return s.cfg.PublishFunc(ctx, event)
}

func (s *Server) Shutdown(ctx context.Context) {
	s.grpcServer.GracefulStop()
	log.Println("[chip-testsink] Server shutdown")
}

func notifyStarted(ch chan<- struct{}) {
	if ch == nil {
		return
	}

	select {
	case ch <- struct{}{}:
	default:
	}
}
