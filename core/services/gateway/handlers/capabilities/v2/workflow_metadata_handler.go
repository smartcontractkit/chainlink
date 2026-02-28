package v2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/gateway"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/common/aggregation"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities/v2/metrics"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const ecdsaPubKeyHexLen = 42 // 2 (0x prefix) + 40 (hex digits)

type workflowReference struct {
	workflowOwner string
	workflowName  string
	workflowTag   string
}

// jwtReplayCache manages used JWT IDs to prevent replay attacks
type jwtReplayCache struct {
	mu            sync.RWMutex
	cleanupPeriod time.Duration
	cache         map[string]time.Time // jti -> timestamp
}

type WorkflowMetadataHandler struct {
	services.StateMachine
	lggr              logger.Logger
	mu                sync.RWMutex
	authorizedKeys    map[string]map[gateway.AuthorizedKey]struct{} // workflow ID -> authorized keys
	workflowRefToID   map[workflowReference]string                  // workflow reference -> workflow ID
	workflowIDToRef   map[string]workflowReference                  // workflow ID -> workflow reference
	workflowIDToShard map[string]*ShardInfo                         // workflow ID -> shard that advertises it
	aggByShard        map[*ShardInfo]*aggregation.WorkflowMetadataAggregator
	shards            []*ShardInfo
	nodeAddrToShard   map[string]*ShardInfo
	config            ServiceConfig
	stopCh            services.StopChan
	metrics           *metrics.Metrics
	jwtCache          *jwtReplayCache
	wg                sync.WaitGroup
	startTime         time.Time
}

// NewWorkflowMetadataHandler creates a new WorkflowMetadataHandler.
func NewWorkflowMetadataHandler(lggr logger.Logger, cfg ServiceConfig, shards []*ShardInfo, nodeAddrToShard map[string]*ShardInfo, metrics *metrics.Metrics) *WorkflowMetadataHandler {
	aggByShard := make(map[*ShardInfo]*aggregation.WorkflowMetadataAggregator, len(shards))
	for _, shard := range shards {
		threshold := shard.DONConfig.F + 1
		aggByShard[shard] = aggregation.NewWorkflowMetadataAggregator(
			lggr, threshold, time.Duration(cfg.CleanUpPeriodMs)*time.Millisecond, metrics,
		)
	}
	return &WorkflowMetadataHandler{
		lggr:              logger.Named(lggr, "HTTPTriggerWorkflowMetadataHandler"),
		authorizedKeys:    make(map[string]map[gateway.AuthorizedKey]struct{}),
		workflowRefToID:   make(map[workflowReference]string),
		workflowIDToRef:   make(map[string]workflowReference),
		workflowIDToShard: make(map[string]*ShardInfo),
		aggByShard:        aggByShard,
		shards:            shards,
		nodeAddrToShard:   nodeAddrToShard,
		config:            cfg,
		stopCh:            make(services.StopChan),
		metrics:           metrics,
		jwtCache:          newJWTReplayCache(time.Duration(cfg.JWTReplayPeriodMs) * time.Millisecond),
	}
}

func (h *WorkflowMetadataHandler) Authorize(workflowID string, token string, req *jsonrpc.Request[json.RawMessage]) (*gateway.AuthorizedKey, error) {
	claims, signer, err := utils.VerifyRequestJWT(token, *req)
	if err != nil {
		h.lggr.Errorw("Failed to verify JWT", "error", err)
		return nil, err
	}

	if h.jwtCache.isReplay(claims.ID) {
		h.lggr.Warnw("JWT token has already been used", "workflowID", workflowID, "signer", signer.Hex(), "jti", claims.ID)
		return nil, errors.New("JWT token has already been used. Please generate a new one with new id (jti)")
	}

	keys, exists := h.authorizedKeys[workflowID]
	if !exists {
		h.lggr.Errorw("Workflow ID not found in authorized keys", "workflowID", workflowID)
		return nil, fmt.Errorf("workflow ID %s not found", workflowID)
	}
	key := gateway.AuthorizedKey{
		KeyType:   gateway.KeyTypeECDSAEVM,
		PublicKey: strings.ToLower(signer.Hex()),
	}
	if _, exists = keys[key]; !exists {
		h.lggr.Errorw("Signer not found in authorized keys", "signer", signer.Hex())
		return nil, fmt.Errorf("signer '%s' is not authorized for workflow '%s'. Ensure that the signer is registered in the workflow definition", signer.Hex(), workflowID)
	}
	h.jwtCache.recordUsage(claims.ID)

	return &key, nil
}

// syncMetadata aggregates the authorized keys and workflow selectors from all
// per-shard WorkflowMetadataAggregators and updates the local cache.
func (h *WorkflowMetadataHandler) syncMetadata() {
	authorizedKeys := make(map[string]map[gateway.AuthorizedKey]struct{})
	workflowRefToID := make(map[workflowReference]string)
	workflowIDToRef := make(map[string]workflowReference)
	workflowIDToShard := make(map[string]*ShardInfo)

	for shard, agg := range h.aggByShard {
		metadata, err := agg.Aggregate()
		if err != nil {
			h.lggr.Errorw("Failed to aggregate auth data", "shard", shard.DONConfig.DonId, "error", err)
			continue
		}
		for _, data := range metadata {
			wfID := data.WorkflowSelector.WorkflowID
			workflowRef := workflowReference{
				workflowOwner: data.WorkflowSelector.WorkflowOwner,
				workflowName:  data.WorkflowSelector.WorkflowName,
				workflowTag:   data.WorkflowSelector.WorkflowTag,
			}
			if _, exists := workflowIDToRef[wfID]; exists {
				h.lggr.Debugw("Duplicate workflow ID across shards", "workflowID", wfID, "shard", shard.DONConfig.DonId)
				continue
			}
			if _, exists := workflowRefToID[workflowRef]; exists {
				h.lggr.Debugw("Duplicate workflow reference across shards", "workflowRef", workflowRef, "shard", shard.DONConfig.DonId)
				continue
			}
			workflowIDToRef[wfID] = workflowRef
			workflowRefToID[workflowRef] = wfID
			workflowIDToShard[wfID] = shard
			authorizedKeys[wfID] = make(map[gateway.AuthorizedKey]struct{})
			for _, key := range data.AuthorizedKeys {
				authorizedKeys[wfID][key] = struct{}{}
			}
		}
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if len(h.workflowIDToRef) == 0 && len(workflowIDToRef) > 0 {
		latencyMs := time.Since(h.startTime).Milliseconds()
		h.metrics.RecordMetadataSyncStartupLatency(context.Background(), latencyMs, h.lggr)
	}
	workflowIDs := make([]string, 0, len(workflowIDToRef))
	for wfID := range workflowIDToRef {
		workflowIDs = append(workflowIDs, wfID)
	}
	h.lggr.Debugw("Synced workflow metadata", "workflowIDs", workflowIDs, "count", len(workflowIDs))

	h.authorizedKeys = authorizedKeys
	h.workflowRefToID = workflowRefToID
	h.workflowIDToRef = workflowIDToRef
	h.workflowIDToShard = workflowIDToShard
	h.metrics.RecordLoadedMetadataSize(context.Background(), int64(len(h.workflowIDToRef)), h.lggr)
}

// sendMetadataPullRequest sends a request to all nodes across all shards to pull the latest metadata.
func (h *WorkflowMetadataHandler) sendMetadataPullRequest() error {
	timeout := time.Duration(h.config.MetadataPullRequestTimeoutMs) * time.Millisecond
	ctx, cancel := h.stopCh.CtxWithTimeout(timeout)
	defer cancel()

	req := &jsonrpc.Request[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      gateway.GetRequestID(gateway.MethodPullWorkflowMetadata),
		Method:  gateway.MethodPullWorkflowMetadata,
	}
	var combinedErr error
	for _, shard := range h.shards {
		for _, member := range shard.DONConfig.Members {
			h.metrics.IncrementTriggerCapabilityRequestCount(ctx, member.Address, gateway.MethodPullWorkflowMetadata, h.lggr)
			err := shard.DON.SendToNode(ctx, member.Address, req)
			if err != nil {
				h.metrics.IncrementTriggerCapabilityRequestFailures(ctx, member.Address, gateway.MethodPullWorkflowMetadata, h.lggr)
				combinedErr = errors.Join(combinedErr, fmt.Errorf("failed to send pull request to node %s (shard %s): %w", member.Address, shard.DONConfig.DonId, err))
			}
		}
	}
	return combinedErr
}

func (h *WorkflowMetadataHandler) aggForNode(nodeAddr string) (*aggregation.WorkflowMetadataAggregator, error) {
	shard, ok := h.nodeAddrToShard[strings.ToLower(nodeAddr)]
	if !ok {
		return nil, fmt.Errorf("unknown node address %s: no shard mapping", nodeAddr)
	}
	agg, ok := h.aggByShard[shard]
	if !ok {
		return nil, fmt.Errorf("no aggregator for shard %s", shard.DONConfig.DonId)
	}
	return agg, nil
}

// OnMetadataPush handles the push of metadata from a node when a new workflow is registered
func (h *WorkflowMetadataHandler) OnMetadataPush(ctx context.Context, resp *jsonrpc.Response[json.RawMessage], nodeAddr string) error {
	var metadata gateway.WorkflowMetadata
	if err := json.Unmarshal(*resp.Result, &metadata); err != nil {
		return fmt.Errorf("failed to unmarshal metadata: %w", err)
	}
	h.lggr.Debugw("Received metadata push", "workflowID", metadata.WorkflowSelector.WorkflowID, "nodeAddr", nodeAddr)
	if err := h.validateAuthMetadata(metadata); err != nil {
		return err
	}
	agg, err := h.aggForNode(nodeAddr)
	if err != nil {
		return err
	}
	if err := agg.Collect(&metadata, nodeAddr); err != nil {
		return fmt.Errorf("failed to collect observation: %w", err)
	}
	return nil
}

// OnMetadataPullResponse handles the response to the metadata pull request.
func (h *WorkflowMetadataHandler) OnMetadataPullResponse(ctx context.Context, resp *jsonrpc.Response[json.RawMessage], nodeAddr string) error {
	var metadata []gateway.WorkflowMetadata
	if err := json.Unmarshal(*resp.Result, &metadata); err != nil {
		return fmt.Errorf("failed to unmarshal metadata pull response: %w", err)
	}
	h.lggr.Debugw("Received metadata pull response", "nodeAddr", nodeAddr)
	for _, data := range metadata {
		if err := h.validateAuthMetadata(data); err != nil {
			return err
		}
	}
	agg, err := h.aggForNode(nodeAddr)
	if err != nil {
		return err
	}
	var combinedErr error
	for _, data := range metadata {
		err := agg.Collect(&data, nodeAddr)
		combinedErr = errors.Join(combinedErr, err)
	}
	return combinedErr
}

// Start begins the periodic pull loop.
func (h *WorkflowMetadataHandler) Start(ctx context.Context) error {
	return h.StartOnce("WorkflowMetadataHandler", func() error {
		h.lggr.Info("Starting HTTP Trigger Metadata Handler")
		h.startTime = time.Now()
		for shard, agg := range h.aggByShard {
			if err := agg.Start(ctx); err != nil {
				return fmt.Errorf("failed to start aggregator for shard %s: %w", shard.DONConfig.DonId, err)
			}
		}
		h.runTicker(time.Duration(h.config.MetadataPullIntervalMs)*time.Millisecond, func() {
			if err := h.sendMetadataPullRequest(); err != nil {
				h.lggr.Errorw("Failed to send pull request", "error", err)
			}
		})
		h.runTicker(time.Duration(h.config.MetadataAggregationIntervalMs)*time.Millisecond, h.syncMetadata)

		h.runTicker(h.jwtCache.cleanupPeriod, func() {
			now := time.Now()
			expiredCount := h.jwtCache.cleanupOldEntries(now.Add(-h.jwtCache.cleanupPeriod))
			h.metrics.IncrementJwtCacheCleanUpCount(context.Background(), int64(expiredCount), h.lggr)
			h.metrics.RecordJwtCacheSize(context.Background(), int64(len(h.jwtCache.cache)), h.lggr)
			h.lggr.Debugw("Workflow execution cache cleanup completed", "expired_entries", expiredCount, "remaining_entries", len(h.jwtCache.cache))
		})
		return nil
	})
}

func (h *WorkflowMetadataHandler) runTicker(period time.Duration, fn func()) {
	h.wg.Add(1)
	go func() {
		defer h.wg.Done()
		ticker := time.NewTicker(period)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				fn()
			case <-h.stopCh:
				return
			}
		}
	}()
}

func (h *WorkflowMetadataHandler) validateAuthMetadata(metadata gateway.WorkflowMetadata) error {
	if len(metadata.WorkflowSelector.WorkflowID) != workflowIDLength {
		return fmt.Errorf("invalid workflow ID: expected %d characters, got %d", workflowIDLength, len(metadata.WorkflowSelector.WorkflowID))
	}
	if len(metadata.WorkflowSelector.WorkflowOwner) != workflowOwnerLength {
		return fmt.Errorf("invalid workflow owner: expected %d characters, got %d", workflowOwnerLength, len(metadata.WorkflowSelector.WorkflowOwner))
	}
	if len(metadata.WorkflowSelector.WorkflowName) != WorkflowNameHashLength {
		return fmt.Errorf("invalid workflow name: expected %d characters, got %d", WorkflowNameHashLength, len(metadata.WorkflowSelector.WorkflowName))
	}
	if len(metadata.WorkflowSelector.WorkflowTag) == 0 || len(metadata.WorkflowSelector.WorkflowTag) > maxWorkflowTagLength {
		return fmt.Errorf("invalid workflow tag: expected non-empty and at most %d characters, got %d", maxWorkflowTagLength, len(metadata.WorkflowSelector.WorkflowTag))
	}
	if len(metadata.AuthorizedKeys) == 0 {
		return errors.New("no authorized keys")
	}
	for _, key := range metadata.AuthorizedKeys {
		if key.KeyType != gateway.KeyTypeECDSAEVM {
			return errors.New("invalid key type")
		}
		if key.PublicKey == "" || !strings.HasPrefix(key.PublicKey, "0x") || len(key.PublicKey) != ecdsaPubKeyHexLen {
			return fmt.Errorf("invalid public key: %s", key.PublicKey)
		}
		if key.PublicKey != strings.ToLower(key.PublicKey) {
			return errors.New("invalid public key: must be all lowercase")
		}
	}
	return nil
}

func (h *WorkflowMetadataHandler) GetWorkflowID(workflowOwner, workflowName, workflowTag string) (string, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	workflowRef := workflowReference{
		workflowOwner: workflowOwner,
		workflowName:  workflowName,
		workflowTag:   workflowTag,
	}
	workflowID, exists := h.workflowRefToID[workflowRef]
	if !exists {
		return "", false
	}
	return workflowID, true
}

func (h *WorkflowMetadataHandler) GetWorkflowReference(workflowID string) (workflowReference, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	workflowRef, exists := h.workflowIDToRef[workflowID]
	return workflowRef, exists
}

func (h *WorkflowMetadataHandler) GetWorkflowShard(workflowID string) (*ShardInfo, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	shard, exists := h.workflowIDToShard[workflowID]
	return shard, exists
}

func (h *WorkflowMetadataHandler) Close() error {
	return h.StopOnce("WorkflowMetadataHandler", func() error {
		h.lggr.Info("Stopping HTTP Trigger Metadata Handler")
		for shard, agg := range h.aggByShard {
			if err := agg.Close(); err != nil {
				h.lggr.Errorw("Failed to close WorkflowMetadataAggregator", "shard", shard.DONConfig.DonId, "error", err)
			}
		}
		close(h.stopCh)
		h.wg.Wait()
		return nil
	})
}

func newJWTReplayCache(cleanupPeriod time.Duration) *jwtReplayCache {
	return &jwtReplayCache{
		cache:         make(map[string]time.Time),
		cleanupPeriod: cleanupPeriod,
	}
}

func (cache *jwtReplayCache) isReplay(jti string) bool {
	cache.mu.RLock()
	defer cache.mu.RUnlock()

	_, exists := cache.cache[jti]
	return exists
}

func (cache *jwtReplayCache) recordUsage(jti string) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	cache.cache[jti] = time.Now()
}

// cleanupOldEntries removes expired entries from the cache
func (cache *jwtReplayCache) cleanupOldEntries(cutoff time.Time) int {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	var expiredCount int
	for jti, createdAt := range cache.cache {
		if createdAt.Before(cutoff) {
			delete(cache.cache, jti)
			expiredCount++
		}
	}
	return expiredCount
}
