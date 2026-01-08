package arbiter

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/services"

	// TODO: Update this import path once proto is moved
	pb "github.com/smartcontractkit/chainlink/v2/core/services/arbiter/proto"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// Arbiter is the main service interface.
type Arbiter interface {
	services.Service
	HealthReport() map[string]error
}

type arbiter struct {
	services.StateMachine

	grpcServer  *grpc.Server
	grpcHandler *GRPCServer
	state       *State
	decision    DecisionEngine
	shardConfig ShardConfigReader
	lggr        logger.Logger

	grpcAddr string
	stopCh   services.StopChan
	wg       sync.WaitGroup
}

var _ Arbiter = (*arbiter)(nil)

// New creates a new Arbiter service.
// contractReaderFactory is used to create the contract reader for querying the ShardConfig contract.
// This follows the same pattern as the workflow registry syncer and capability registry syncer.
func New(
	lggr logger.Logger,
	contractReaderFactory ContractReaderFactory,
	shardConfigAddr string,
	port uint16,
	pollInterval time.Duration,
	retryInterval time.Duration,
) (Arbiter, error) {
	lggr = lggr.Named("Arbiter")

	// Create state
	state := NewState()

	// Create ShardConfig syncer (implements services.Service)
	shardConfig := NewShardConfigSyncer(contractReaderFactory, shardConfigAddr, pollInterval, retryInterval, lggr)

	// Create decision engine with sugared logger
	decision := NewDecisionEngine(shardConfig, logger.Sugared(lggr))

	// Create gRPC handler
	grpcHandler := NewGRPCServer(state, decision, lggr)

	// Create gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterArbiterServiceServer(grpcServer, grpcHandler)

	return &arbiter{
		grpcServer:  grpcServer,
		grpcHandler: grpcHandler,
		state:       state,
		decision:    decision,
		shardConfig: shardConfig,
		lggr:        lggr,
		grpcAddr:    fmt.Sprintf(":%d", port),
		stopCh:      make(services.StopChan),
	}, nil
}

// Start starts the Arbiter service.
func (a *arbiter) Start(ctx context.Context) error {
	return a.StartOnce("Arbiter", func() error {
		a.lggr.Info("Starting Arbiter service")

		// Start ShardConfig syncer first
		if err := a.shardConfig.Start(ctx); err != nil {
			return fmt.Errorf("failed to start shard config syncer: %w", err)
		}

		// Start gRPC server in a goroutine
		a.wg.Add(1)
		go func() {
			defer a.wg.Done()
			a.runGRPCServer(ctx)
		}()

		a.lggr.Infow("Arbiter service started",
			"grpcAddr", a.grpcAddr,
		)

		return nil
	})
}

// runGRPCServer starts the gRPC server and blocks until stopped.
func (a *arbiter) runGRPCServer(ctx context.Context) {
	var lc net.ListenConfig
	lis, err := lc.Listen(ctx, "tcp", a.grpcAddr)
	if err != nil {
		a.lggr.Errorw("Failed to listen for gRPC",
			"addr", a.grpcAddr,
			"error", err,
		)
		return
	}

	a.lggr.Infow("gRPC server listening",
		"addr", a.grpcAddr,
	)

	if err := a.grpcServer.Serve(lis); err != nil {
		// Check if this is a normal shutdown
		select {
		case <-a.stopCh:
			// Normal shutdown, don't log as error
			a.lggr.Debug("gRPC server stopped")
		default:
			a.lggr.Errorw("gRPC server error",
				"error", err,
			)
		}
	}
}

// Close stops the Arbiter service.
func (a *arbiter) Close() error {
	return a.StopOnce("Arbiter", func() (err error) {
		a.lggr.Info("Stopping Arbiter service")

		// Signal stop
		close(a.stopCh)

		// Graceful shutdown of gRPC server
		a.grpcServer.GracefulStop()
		a.lggr.Debug("gRPC server stopped gracefully")

		// Wait for gRPC goroutine
		a.wg.Wait()

		// Close ShardConfig syncer
		if err := a.shardConfig.Close(); err != nil {
			a.lggr.Errorw("Failed to close shard config syncer", "error", err)
		}

		a.lggr.Info("Arbiter service stopped")

		return nil
	})
}

// HealthReport returns the health status of the service.
func (a *arbiter) HealthReport() map[string]error {
	return map[string]error{
		a.Name(): a.Ready(),
	}
}

// Name returns the service name.
func (a *arbiter) Name() string {
	return a.lggr.Name()
}
