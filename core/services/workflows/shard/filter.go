package shard

import (
	"context"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	shardorchpb "github.com/smartcontractkit/chainlink-common/pkg/workflows/shardorchestrator/pb"
)

// Filter determines whether a workflow should be executed on this shard node
// by checking with the ShardOrchestrator and maintaining a TTL cache.
type Filter struct {
	orchestratorClient shardorchpb.ShardOrchestratorServiceClient
	myShardIndex       uint32
	cacheTTL           time.Duration
	logger             logger.Logger

	mu    sync.RWMutex
	cache map[string]*cacheEntry
}

type cacheEntry struct {
	shardID   uint32
	expiresAt time.Time
}

// Config holds the configuration for creating a shard filter.
type Config struct {
	// ShardIndex is this node's shard index (0 for shard-zero, 1 for shard-1, etc.)
	ShardIndex uint32
	// OrchestratorAddress is the gRPC address of the ShardOrchestrator (e.g., "localhost:60051")
	OrchestratorAddress string
	// CacheTTL is the duration for which workflow mappings are cached
	CacheTTL time.Duration
	// Logger for debug/error messages
	Logger logger.Logger
}

// NewFilter creates a new shard filter with the given configuration.
func NewFilter(cfg Config) (*Filter, error) {
	if cfg.OrchestratorAddress == "" {
		return nil, fmt.Errorf("orchestrator address is required")
	}
	if cfg.CacheTTL <= 0 {
		cfg.CacheTTL = 5 * time.Second
	}
	if cfg.Logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	conn, err := grpc.NewClient(cfg.OrchestratorAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client for ShardOrchestrator: %w", err)
	}

	client := shardorchpb.NewShardOrchestratorServiceClient(conn)

	return &Filter{
		orchestratorClient: client,
		myShardIndex:       cfg.ShardIndex,
		cacheTTL:           cfg.CacheTTL,
		logger:             cfg.Logger,
		cache:              make(map[string]*cacheEntry),
	}, nil
}

// ShouldExecute returns true if the workflow should be executed on this shard.
// It first checks the cache, and if there's a miss, fetches from the ShardOrchestrator.
func (f *Filter) ShouldExecute(ctx context.Context, workflowID string) bool {
	// Check cache first
	f.mu.RLock()
	entry, found := f.cache[workflowID]
	f.mu.RUnlock()

	if found && time.Now().Before(entry.expiresAt) {
		// Cache hit - return cached result
		shouldExecute := entry.shardID == f.myShardIndex
		f.logger.Debugw("Shard filter cache hit",
			"workflowID", workflowID,
			"cachedShardID", entry.shardID,
			"myShardIndex", f.myShardIndex,
			"shouldExecute", shouldExecute,
		)
		return shouldExecute
	}

	// Cache miss or expired - fetch from orchestrator
	f.logger.Debugw("Shard filter cache miss, fetching from orchestrator",
		"workflowID", workflowID,
		"myShardIndex", f.myShardIndex,
	)

	resp, err := f.orchestratorClient.GetWorkflowShardMapping(ctx, &shardorchpb.GetWorkflowShardMappingRequest{
		WorkflowIds: []string{workflowID},
	})
	if err != nil {
		f.logger.Errorw("Failed to get workflow shard mapping from orchestrator",
			"workflowID", workflowID,
			"error", err,
		)
		// Fail-safe: skip execution if we can't determine shard assignment
		return false
	}

	// Update cache with the response
	f.mu.Lock()
	now := time.Now()
	for wfID, shardID := range resp.Mappings {
		f.cache[wfID] = &cacheEntry{
			shardID:   shardID,
			expiresAt: now.Add(f.cacheTTL),
		}
	}
	f.mu.Unlock()

	assignedShardID, ok := resp.Mappings[workflowID]
	if !ok {
		f.logger.Warnw("Workflow not found in orchestrator response",
			"workflowID", workflowID,
		)
		return false
	}

	shouldExecute := assignedShardID == f.myShardIndex
	f.logger.Infow("Shard filter decision from orchestrator",
		"workflowID", workflowID,
		"assignedShardID", assignedShardID,
		"myShardIndex", f.myShardIndex,
		"shouldExecute", shouldExecute,
	)

	return shouldExecute
}

// MyShardIndex returns this filter's shard index.
func (f *Filter) MyShardIndex() uint32 {
	return f.myShardIndex
}

// ClearCache removes all cached entries (useful for testing or forced refresh).
func (f *Filter) ClearCache() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.cache = make(map[string]*cacheEntry)
	f.logger.Info("Shard filter cache cleared")
}
