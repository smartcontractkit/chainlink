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
	lggr            logger.Logger
	mu              sync.RWMutex
	authorizedKeys  map[string]map[gateway.AuthorizedKey]struct{} // map of workflow ID to authorized keys
	workflowRefToID map[workflowReference]string                  // map of workflow reference to workflow ID
	workflowIDToRef map[string]workflowReference                  // map of workflow ID to workflow reference
	workflowShards  map[string][]*shardEndpoint                   // map of workflow ID to the shards it is assigned to (quorum reached)
	// aggs holds one WorkflowMetadataAggregator per shard, keyed by shard donID.
	aggs            map[string]*aggregation.WorkflowMetadataAggregator
	shards          []*shardEndpoint
	nodeAddrToShard map[string]*shardEndpoint
	config          ServiceConfig
	stopCh          services.StopChan
	metrics         *metrics.Metrics
	jwtCache        *jwtReplayCache // JWT replay protection cache
	wg              sync.WaitGroup
	startTime       time.Time // time when Start() was called
}

// NewWorkflowMetadataHandler creates a new WorkflowMetadataHandler spanning the
// full DON×shard matrix. Each shard gets its own aggregator with threshold F+1.
func NewWorkflowMetadataHandler(lggr logger.Logger, cfg ServiceConfig, shards []*shardEndpoint, nodeAddrToShard map[string]*shardEndpoint, metrics *metrics.Metrics) *WorkflowMetadataHandler {
	aggs := make(map[string]*aggregation.WorkflowMetadataAggregator, len(shards))
	for _, shard := range shards {
		threshold := shard.f + 1
		aggs[shard.donID] = aggregation.NewWorkflowMetadataAggregator(lggr, threshold, time.Duration(cfg.CleanUpPeriodMs)*time.Millisecond, metrics)
	}
	return &WorkflowMetadataHandler{
		lggr:            logger.Named(lggr, "HTTPTriggerWorkflowMetadataHandler"),
		authorizedKeys:  make(map[string]map[gateway.AuthorizedKey]struct{}),
		workflowRefToID: make(map[workflowReference]string),
		workflowIDToRef: make(map[string]workflowReference),
		workflowShards:  make(map[string][]*shardEndpoint),
		aggs:            aggs,
		shards:          shards,
		nodeAddrToShard: nodeAddrToShard,
		config:          cfg,
		stopCh:          make(services.StopChan),
		metrics:         metrics,
		jwtCache:        newJWTReplayCache(time.Duration(cfg.JWTReplayPeriodMs) * time.Millisecond),
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

// syncMetadata aggregates the authorized keys and workflow selectors from each
// shard's WorkflowMetadataAggregator and updates the local cache. A workflow is
// considered assigned to a shard once that shard's aggregator reports it (i.e.
// F+1 of the shard's nodes observed it).
func (h *WorkflowMetadataHandler) syncMetadata(ctx context.Context) {
	authorizedKeys := make(map[string]map[gateway.AuthorizedKey]struct{})
	workflowRefToID := make(map[workflowReference]string)
	workflowIDToRef := make(map[string]workflowReference)
	workflowShards := make(map[string][]*shardEndpoint)

	for _, shard := range h.shards {
		agg := h.aggs[shard.donID]
		metadata := agg.Aggregate()
		for _, data := range metadata {
			workflowID := data.WorkflowSelector.WorkflowID
			workflowRef := workflowReference{
				workflowOwner: data.WorkflowSelector.WorkflowOwner,
				workflowName:  data.WorkflowSelector.WorkflowName,
				workflowTag:   data.WorkflowSelector.WorkflowTag,
			}

			// Case 1: this workflow ID was already registered. If the reference
			// matches, this is the same workflow reported by another shard —
			// append the shard to its fan-out list. If the reference differs,
			// it's a conflicting observation; drop it.
			if existingRef, idExists := workflowIDToRef[workflowID]; idExists {
				if existingRef == workflowRef {
					workflowShards[workflowID] = append(workflowShards[workflowID], shard)
				} else {
					h.lggr.Debugw("Duplicate workflow ID with conflicting reference, dropping",
						"workflowID", workflowID, "existingRef", existingRef, "conflictingRef", workflowRef)
				}
				continue
			}

			// Case 2: this workflow reference was already registered under a
			// different workflow ID. First-wins by reference; drop the duplicate.
			if _, refExists := workflowRefToID[workflowRef]; refExists {
				h.lggr.Debugw("Duplicate workflow reference found, dropping",
					"workflowRef", workflowRef, "workflowID", workflowID)
				continue
			}

			// Case 3: new workflow ID and reference — register it.
			workflowIDToRef[workflowID] = workflowRef
			workflowRefToID[workflowRef] = workflowID
			authorizedKeys[workflowID] = make(map[gateway.AuthorizedKey]struct{})
			for _, key := range data.AuthorizedKeys {
				authorizedKeys[workflowID][key] = struct{}{}
			}
			workflowShards[workflowID] = append(workflowShards[workflowID], shard)
		}
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if len(h.workflowIDToRef) == 0 && len(workflowIDToRef) > 0 {
		latencyMs := time.Since(h.startTime).Milliseconds()
		h.metrics.RecordMetadataSyncStartupLatency(ctx, latencyMs, h.lggr)
	}
	// Log all registered workflow IDs
	workflowIDs := make([]string, 0, len(workflowIDToRef))
	for workflowID := range workflowIDToRef {
		workflowIDs = append(workflowIDs, workflowID)
	}
	h.lggr.Debugw("Synced workflow metadata", "workflowIDs", workflowIDs, "count", len(workflowIDs))

	h.authorizedKeys = authorizedKeys
	h.workflowRefToID = workflowRefToID
	h.workflowIDToRef = workflowIDToRef
	h.workflowShards = workflowShards
	h.metrics.RecordLoadedMetadataSize(ctx, int64(len(h.workflowIDToRef)), h.lggr)
}

// sendMetadataPullRequest sends a request to all nodes in every shard to pull
// the latest metadata. no retries are performed, as the caller is expected to
// poll periodically.
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
		for _, member := range shard.members {
			h.metrics.IncrementTriggerCapabilityRequestCount(ctx, member.Address, gateway.MethodPullWorkflowMetadata, h.lggr)
			err := shard.connMgr.SendToNode(ctx, member.Address, req)
			if err != nil {
				h.metrics.IncrementTriggerCapabilityRequestFailures(ctx, member.Address, gateway.MethodPullWorkflowMetadata, h.lggr)
				combinedErr = errors.Join(combinedErr, fmt.Errorf("failed to send pull request to node %s (shard %s): %w", member.Address, shard.donID, err))
			}
		}
	}
	return combinedErr
}

// OnMetadataPush handles the push of metadata from a node when a new workflow is registered
func (h *WorkflowMetadataHandler) OnMetadataPush(ctx context.Context, resp *jsonrpc.Response[json.RawMessage], nodeAddr string) error {
	var metadata gateway.WorkflowMetadata
	if err := json.Unmarshal(*resp.Result, &metadata); err != nil {
		return fmt.Errorf("failed to unmarshal metadata: %w", err)
	}
	h.lggr.Debugw("Received metadata push", "workflowID", metadata.WorkflowSelector.WorkflowID, "nodeAddr", nodeAddr)
	err := h.validateAuthMetadata(metadata)
	if err != nil {
		return err
	}
	agg, err := h.aggForNode(nodeAddr)
	if err != nil {
		return err
	}
	var combinedErr error
	err = agg.Collect(&metadata, nodeAddr)
	if err != nil {
		combinedErr = errors.Join(combinedErr, fmt.Errorf("failed to collect observation: %w", err))
	}
	return combinedErr
}

// OnMetadataPullResponse handles the response to the metadata pull request.
func (h *WorkflowMetadataHandler) OnMetadataPullResponse(ctx context.Context, resp *jsonrpc.Response[json.RawMessage], nodeAddr string) error {
	var metadata []gateway.WorkflowMetadata
	if err := json.Unmarshal(*resp.Result, &metadata); err != nil {
		return fmt.Errorf("failed to unmarshal metadata pull response: %w", err)
	}
	h.lggr.Debugw("Received metadata pull response", "nodeAddr", nodeAddr)
	for _, data := range metadata {
		err := h.validateAuthMetadata(data)
		if err != nil {
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

// aggForNode returns the per-shard aggregator that owns the given node address.
func (h *WorkflowMetadataHandler) aggForNode(nodeAddr string) (*aggregation.WorkflowMetadataAggregator, error) {
	shard, ok := h.nodeAddrToShard[nodeAddr]
	if !ok {
		return nil, fmt.Errorf("received metadata from unknown node %s (no owning shard)", nodeAddr)
	}
	return h.aggs[shard.donID], nil
}

// WorkflowShards returns the shards a workflow is currently assigned to (those
// whose quorum reported the workflow's metadata). The returned slice is a copy.
func (h *WorkflowMetadataHandler) WorkflowShards(workflowID string) []*shardEndpoint {
	h.mu.RLock()
	defer h.mu.RUnlock()
	shards := h.workflowShards[workflowID]
	out := make([]*shardEndpoint, len(shards))
	copy(out, shards)
	return out
}

// Start begins the periodic pull loop.
func (h *WorkflowMetadataHandler) Start(ctx context.Context) error {
	return h.StartOnce("WorkflowMetadataHandler", func() error {
		h.lggr.Info("Starting HTTP Trigger Metadata Handler")
		h.startTime = time.Now()
		for _, shard := range h.shards {
			if err := h.aggs[shard.donID].Start(ctx); err != nil {
				return fmt.Errorf("failed to start aggregator for shard %s: %w", shard.donID, err)
			}
		}
		h.runTicker(time.Duration(h.config.MetadataPullIntervalMs)*time.Millisecond, func(ctx context.Context) {
			err2 := h.sendMetadataPullRequest()
			if err2 != nil {
				h.lggr.Errorw("Failed to send pull request", "error", err2)
			}
		})
		h.runTicker(time.Duration(h.config.MetadataAggregationIntervalMs)*time.Millisecond, h.syncMetadata)

		h.runTicker(h.jwtCache.cleanupPeriod, func(ctx context.Context) {
			now := time.Now()
			expiredCount := h.jwtCache.cleanupOldEntries(now.Add(-h.jwtCache.cleanupPeriod))
			h.metrics.IncrementJwtCacheCleanUpCount(ctx, int64(expiredCount), h.lggr)
			h.metrics.RecordJwtCacheSize(ctx, int64(len(h.jwtCache.cache)), h.lggr)
			h.lggr.Debugw("Workflow execution cache cleanup completed", "expired_entries", expiredCount, "remaining_entries", len(h.jwtCache.cache))
		})
		return nil
	})
}

func (h *WorkflowMetadataHandler) runTicker(period time.Duration, fn func(ctx context.Context)) {
	h.wg.Go(func() {
		ctx, cancel := h.stopCh.NewCtx()
		defer cancel()
		ticker := time.NewTicker(period)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				fn(ctx)
			case <-ctx.Done():
				return
			}
		}
	})
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

func (h *WorkflowMetadataHandler) Close() error {
	return h.StopOnce("WorkflowMetadataHandler", func() error {
		h.lggr.Info("Stopping HTTP Trigger Metadata Handler")
		for _, shard := range h.shards {
			if err := h.aggs[shard.donID].Close(); err != nil {
				h.lggr.Errorw("Failed to close WorkflowMetadataAggregator", "shard", shard.donID, "error", err)
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
