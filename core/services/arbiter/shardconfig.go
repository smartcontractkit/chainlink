package arbiter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink-evm/contracts/cre/gobindings/shardconfig/generated/v1_0_0/shard_config"
	"github.com/smartcontractkit/chainlink-evm/pkg/config"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const (
	// ShardConfigContractName is the name used to identify the ShardConfig contract.
	ShardConfigContractName = "ShardConfig"

	// GetDesiredShardCountMethod is the method name for reading the desired shard count.
	GetDesiredShardCountMethod = "getDesiredShardCount"
)

// ContractReaderFactory creates a ContractReader from a config.
// This matches the pattern used in workflow registry syncer.
type ContractReaderFactory func(context.Context, []byte) (types.ContractReader, error)

// ShardConfigReader reads the desired shard count from the ShardConfig contract.
type ShardConfigReader interface {
	services.Service
	// GetDesiredShardCount retrieves the current desired shard count from cache.
	// Returns an error if the value hasn't been fetched yet.
	GetDesiredShardCount(ctx context.Context) (uint64, error)
}

// shardConfigSyncer implements ShardConfigReader with periodic polling and caching.
type shardConfigSyncer struct {
	services.StateMachine
	stopCh services.StopChan
	wg     sync.WaitGroup

	lggr                  logger.Logger
	shardConfigAddress    string
	contractReaderFactory ContractReaderFactory
	contractReader        types.ContractReader

	// Cached value with mutex protection
	cachedShardCount uint64
	cachedMu         sync.RWMutex

	// Polling configuration
	pollInterval time.Duration
	retryTimeout time.Duration
}

var _ ShardConfigReader = (*shardConfigSyncer)(nil)

// NewShardConfigSyncer creates a new ShardConfigReader that polls the contract periodically.
func NewShardConfigSyncer(
	contractReaderFactory ContractReaderFactory,
	shardConfigAddress string,
	pollInterval time.Duration,
	retryTimeout time.Duration,
	lggr logger.Logger,
) ShardConfigReader {
	return &shardConfigSyncer{
		lggr:                  lggr.Named("ShardConfigSyncer"),
		shardConfigAddress:    shardConfigAddress,
		contractReaderFactory: contractReaderFactory,
		pollInterval:          pollInterval,
		retryTimeout:          retryTimeout,
		stopCh:                make(services.StopChan),
	}
}

// Start starts the ShardConfig syncer service.
func (s *shardConfigSyncer) Start(ctx context.Context) error {
	return s.StartOnce("ShardConfigSyncer", func() error {
		s.lggr.Info("Starting ShardConfig syncer")

		// Start async initialization and polling
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.run(ctx)
		}()

		return nil
	})
}

// Close stops the ShardConfig syncer service.
func (s *shardConfigSyncer) Close() error {
	return s.StopOnce("ShardConfigSyncer", func() error {
		s.lggr.Info("Stopping ShardConfig syncer")

		// Signal stop
		close(s.stopCh)

		// Wait for goroutines
		s.wg.Wait()

		// Close contract reader if initialized
		if s.contractReader != nil {
			if err := s.contractReader.Close(); err != nil {
				s.lggr.Errorw("Failed to close contract reader", "error", err)
			}
		}

		s.lggr.Info("ShardConfig syncer stopped")
		return nil
	})
}

// Name returns the service name.
func (s *shardConfigSyncer) Name() string {
	return s.lggr.Name()
}

// HealthReport returns the health status of the service.
func (s *shardConfigSyncer) HealthReport() map[string]error {
	return map[string]error{
		s.Name(): s.Ready(),
	}
}

// GetDesiredShardCount retrieves the cached desired shard count.
func (s *shardConfigSyncer) GetDesiredShardCount(ctx context.Context) (uint64, error) {
	s.cachedMu.RLock()
	defer s.cachedMu.RUnlock()

	if s.cachedShardCount == 0 {
		return 0, errors.New("shard count not yet available")
	}

	return s.cachedShardCount, nil
}

// run handles initialization and polling loop.
func (s *shardConfigSyncer) run(ctx context.Context) {
	// Phase 1: Initialize contract reader with retries
	if err := s.initContractReader(ctx); err != nil {
		s.lggr.Errorw("Failed to initialize contract reader", "error", err)
		return
	}

	// Phase 2: Start polling loop
	s.pollLoop(ctx)
}

// initContractReader initializes the contract reader with retry logic.
// This follows the lazy initialization pattern used by workflow registry syncer.
func (s *shardConfigSyncer) initContractReader(ctx context.Context) error {
	ticker := time.NewTicker(s.retryTimeout)
	defer ticker.Stop()

	// Try immediately first
	reader, err := s.newContractReader(ctx)
	if err == nil {
		s.contractReader = reader
		s.lggr.Info("Contract reader initialized successfully")
		return nil
	}
	s.lggr.Infow("Contract reader unavailable, will retry", "error", err)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-s.stopCh:
			return nil
		case <-ticker.C:
			reader, err := s.newContractReader(ctx)
			if err != nil {
				s.lggr.Infow("Contract reader unavailable, retrying", "error", err)
				continue
			}
			s.contractReader = reader
			s.lggr.Info("Contract reader initialized successfully")
			return nil
		}
	}
}

// newContractReader creates and configures a new contract reader.
func (s *shardConfigSyncer) newContractReader(ctx context.Context) (types.ContractReader, error) {
	cfg := buildShardConfigReaderConfig()

	cfgBytes, err := json.Marshal(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %w", err)
	}

	reader, err := s.contractReaderFactory(ctx, cfgBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create contract reader: %w", err)
	}
	if reader == nil {
		return nil, errors.New("contract reader factory returned nil")
	}

	// Bind the contract address
	bc := types.BoundContract{
		Address: s.shardConfigAddress,
		Name:    ShardConfigContractName,
	}

	if err := reader.Bind(ctx, []types.BoundContract{bc}); err != nil {
		return nil, fmt.Errorf("failed to bind contract: %w", err)
	}

	if err := reader.Start(ctx); err != nil {
		return nil, fmt.Errorf("failed to start contract reader: %w", err)
	}

	return reader, nil
}

// buildShardConfigReaderConfig creates the ChainReaderConfig for the ShardConfig contract.
func buildShardConfigReaderConfig() config.ChainReaderConfig {
	return config.ChainReaderConfig{
		Contracts: map[string]config.ChainContractReader{
			ShardConfigContractName: {
				ContractABI: shard_config.ShardConfigABI,
				Configs: map[string]*config.ChainReaderDefinition{
					GetDesiredShardCountMethod: {
						ChainSpecificName: GetDesiredShardCountMethod,
						ReadType:          config.Method,
					},
				},
			},
		},
	}
}

// pollLoop periodically fetches the shard count from the contract.
func (s *shardConfigSyncer) pollLoop(ctx context.Context) {
	ticker := time.NewTicker(s.pollInterval)
	defer ticker.Stop()

	// Initial fetch
	s.fetchAndCache(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.stopCh:
			return
		case <-ticker.C:
			s.fetchAndCache(ctx)
		}
	}
}

// fetchAndCache fetches the shard count from the contract and updates the cache.
func (s *shardConfigSyncer) fetchAndCache(ctx context.Context) {
	if s.contractReader == nil {
		return
	}

	var result *big.Int
	bc := types.BoundContract{
		Address: s.shardConfigAddress,
		Name:    ShardConfigContractName,
	}

	err := s.contractReader.GetLatestValue(
		ctx,
		bc.ReadIdentifier(GetDesiredShardCountMethod),
		primitives.Finalized, // Use finalized for governance data
		nil,                  // No input params
		&result,
	)
	if err != nil {
		s.lggr.Errorw("Failed to fetch shard count from on-chain",
			"error", err,
			"contract", s.shardConfigAddress,
		)
		return
	}

	count := result.Uint64()

	s.cachedMu.Lock()
	s.cachedShardCount = count
	s.cachedMu.Unlock()

	// Update metrics
	SetOnChainShardNumber(count)

	s.lggr.Debugw("Fetched shard count from on-chain",
		"count", count,
		"contract", s.shardConfigAddress,
	)
}
