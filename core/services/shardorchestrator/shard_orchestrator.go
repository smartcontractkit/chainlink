package shardorchestrator

import (
	"context"
	"fmt"
	"net"
	"sync"

	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

// ShardOrchestrator is the main service interface for the orchestrator.
type ShardOrchestrator interface {
	services.Service
	HealthReport() map[string]error
	GetAddress() string
}

type orchestrator struct {
	services.StateMachine

	grpcServer  *grpc.Server
	grpcHandler *Server
	listener    net.Listener
	lggr        logger.Logger

	grpcAddr   string
	stopCh     services.StopChan
	wg         sync.WaitGroup
	listenerMu sync.RWMutex
}

var _ ShardOrchestrator = (*orchestrator)(nil)

// New creates a new ShardOrchestrator service.
func New(port int, store *Store, lggr logger.Logger) ShardOrchestrator {
	lggr = logger.Named(lggr, "ShardOrchestrator")

	grpcHandler := NewServer(store, lggr)

	grpcServer := grpc.NewServer()
	grpcHandler.RegisterWithGRPCServer(grpcServer)

	return &orchestrator{
		grpcServer:  grpcServer,
		grpcHandler: grpcHandler,
		lggr:        lggr,
		grpcAddr:    fmt.Sprintf(":%d", port),
		stopCh:      make(services.StopChan),
	}
}

// Start starts the ShardOrchestrator service.
func (o *orchestrator) Start(ctx context.Context) error {
	return o.StartOnce("ShardOrchestrator", func() error {
		o.lggr.Infow("Starting ShardOrchestrator service", "grpcAddr", o.grpcAddr)

		// Start gRPC server in a goroutine
		o.wg.Add(1)
		go func() {
			defer o.wg.Done()
			o.runGRPCServer(ctx)
		}()

		o.lggr.Infow("ShardOrchestrator service started", "grpcAddr", o.grpcAddr)
		return nil
	})
}

// runGRPCServer starts the gRPC server and blocks until stopped.
func (o *orchestrator) runGRPCServer(ctx context.Context) {
	var lc net.ListenConfig
	lis, err := lc.Listen(ctx, "tcp", o.grpcAddr)
	if err != nil {
		o.lggr.Errorw("Failed to listen for gRPC", "addr", o.grpcAddr, "error", err)
		return
	}

	o.listenerMu.Lock()
	o.listener = lis
	o.listenerMu.Unlock()
	o.lggr.Infow("gRPC server listening", "addr", lis.Addr().String())

	if err := o.grpcServer.Serve(lis); err != nil {
		// Check if this is a normal shutdown
		select {
		case <-o.stopCh:
			// Normal shutdown, don't log as error
			o.lggr.Debug("gRPC server stopped")
		default:
			o.lggr.Errorw("gRPC server error", "error", err)
		}
	}
}

// Close stops the ShardOrchestrator service gracefully.
func (o *orchestrator) Close() error {
	return o.StopOnce("ShardOrchestrator", func() error {
		o.lggr.Info("Stopping ShardOrchestrator service")

		// Signal stop
		close(o.stopCh)

		// Graceful shutdown of gRPC server
		o.grpcServer.GracefulStop()
		o.lggr.Debug("gRPC server stopped gracefully")

		// Wait for goroutines to finish
		o.wg.Wait()

		o.lggr.Info("ShardOrchestrator service stopped")
		return nil
	})
}

func (o *orchestrator) Ready() error                   { return nil }
func (o *orchestrator) HealthReport() map[string]error { return nil }
func (o *orchestrator) Name() string                   { return "ShardOrchestrator" }

// GetAddress returns the address the gRPC server is listening on.
// Returns empty string if the server hasn't started yet.
func (o *orchestrator) GetAddress() string {
	o.listenerMu.RLock()
	defer o.listenerMu.RUnlock()
	if o.listener == nil {
		return ""
	}
	return o.listener.Addr().String()
}
