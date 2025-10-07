package v2

import (
	"context"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/gateway"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities/v2/metrics"
)

// responseCache is a thread-safe cache for storing HTTP responses.
// It uses a map to store responses keyed by a unique identifier generated from the request
// cache key is prefixed by workflowID to avoid collisions between different workflows.
type responseCache struct {
	cacheMu sync.Mutex
	cache   map[string]*cachedResponse
	lggr    logger.Logger
	ttl     time.Duration
	metrics *metrics.Metrics
}

type cachedResponse struct {
	response gateway.OutboundHTTPResponse
	storedAt time.Time
}

func newResponseCache(lggr logger.Logger, ttlMs int, metrics *metrics.Metrics) *responseCache {
	return &responseCache{
		cache:   make(map[string]*cachedResponse),
		lggr:    logger.Named(lggr, "ResponseCache"),
		ttl:     time.Duration(ttlMs) * time.Millisecond,
		metrics: metrics,
	}
}

// isCacheableStatusCode returns true if the HTTP status code indicates a cacheable response.
// This includes successful responses (2xx) and client errors (4xx)
func isCacheableStatusCode(statusCode int) bool {
	return (statusCode >= 200 && statusCode < 300) || (statusCode >= 400 && statusCode < 500)
}

// isExpiredOrNotCached returns true if the cached response is expired or not cached.
// IMPORTANT: this method does not lock the cache map. MUST be called with the cacheMu locked.
func (rc *responseCache) isExpiredOrNotCached(workflowID string, req gateway.OutboundHTTPRequest) bool {
	cachedResp, exists := rc.cache[req.Hash()]
	if !exists || time.Now().After(cachedResp.storedAt.Add(rc.ttl)) {
		return true
	}
	return false
}

// Fetch fetches a response from the cache if it exists and
// the age of cached response is less than the max age of the request.
// If the cached response is expired or not cached, it fetches a new response from the fetchFn.
// and caches the response if it is cacheable and storeOnFetch is true.
func (rc *responseCache) Fetch(ctx context.Context, workflowID string, req gateway.OutboundHTTPRequest, fetchFn func() gateway.OutboundHTTPResponse, storeOnFetch bool) gateway.OutboundHTTPResponse {
	rc.cacheMu.Lock()
	defer rc.cacheMu.Unlock()
	cacheMaxAge := time.Duration(req.CacheSettings.MaxAgeMs) * time.Millisecond
	cachedResp, exists := rc.cache[req.Hash()]
	if exists && cachedResp.storedAt.Add(cacheMaxAge).After(time.Now()) {
		rc.metrics.Action.IncrementCacheHitCount(ctx, rc.lggr)
		return cachedResp.response
	}
	response := fetchFn()
	if storeOnFetch && isCacheableStatusCode(response.StatusCode) && rc.isExpiredOrNotCached(workflowID, req) {
		rc.cache[req.Hash()] = &cachedResponse{
			response: response,
			storedAt: time.Now(),
		}
	}
	return response
}

// Set caches a response if it is cacheable (2xx or 4xx and cache is empty or expired for the given request)
func (rc *responseCache) Set(workflowID string, req gateway.OutboundHTTPRequest, response gateway.OutboundHTTPResponse) {
	rc.cacheMu.Lock()
	defer rc.cacheMu.Unlock()
	if isCacheableStatusCode(response.StatusCode) && rc.isExpiredOrNotCached(workflowID, req) {
		rc.cache[req.Hash()] = &cachedResponse{
			response: response,
			storedAt: time.Now(),
		}
	}
}

func (rc *responseCache) DeleteExpired(ctx context.Context) int {
	rc.cacheMu.Lock()
	defer rc.cacheMu.Unlock()
	now := time.Now()
	var expiredCount int
	for key, cachedResp := range rc.cache {
		if now.After(cachedResp.storedAt.Add(rc.ttl)) {
			delete(rc.cache, key)
			expiredCount++
		}
	}
	rc.lggr.Debugw("Removed expired cached HTTP responses", "count", expiredCount, "remaining", len(rc.cache))
	rc.metrics.Action.IncrementCacheCleanUpCount(ctx, int64(expiredCount), rc.lggr)
	rc.metrics.Action.RecordCacheSize(ctx, int64(len(rc.cache)), rc.lggr)
	return expiredCount
}
