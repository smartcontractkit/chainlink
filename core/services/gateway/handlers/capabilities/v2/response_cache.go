package v2

import (
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/gateway"
)

// responseCache is a thread-safe cache for storing HTTP responses.
// It uses a map to store responses keyed by a unique identifier generated from the request
type responseCache struct {
	cacheMu sync.RWMutex
	cache   map[string]*cachedResponse
	lggr    logger.Logger
}

type cachedResponse struct {
	response gateway.OutboundHTTPResponse
	expiry   time.Time
}

func newResponseCache(lggr logger.Logger) *responseCache {
	return &responseCache{
		cache: make(map[string]*cachedResponse),
		lggr:  logger.Named(lggr, "ResponseCache"),
	}
}

func (rc *responseCache) Set(req gateway.OutboundHTTPRequest, response gateway.OutboundHTTPResponse, ttl time.Duration) {
	rc.cacheMu.Lock()
	defer rc.cacheMu.Unlock()
	rc.cache[req.Hash()] = &cachedResponse{
		response: response,
		expiry:   time.Now().Add(ttl),
	}
}

func (rc *responseCache) Get(req gateway.OutboundHTTPRequest) *gateway.OutboundHTTPResponse {
	rc.cacheMu.RLock()
	defer rc.cacheMu.RUnlock()
	cachedResp, exists := rc.cache[req.Hash()]
	if !exists || time.Now().After(cachedResp.expiry) {
		return nil
	}
	return &cachedResp.response
}

func (rc *responseCache) DeleteExpired() int {
	rc.cacheMu.Lock()
	defer rc.cacheMu.Unlock()
	now := time.Now()
	var expiredCount int
	for key, cachedResp := range rc.cache {
		if now.After(cachedResp.expiry) {
			delete(rc.cache, key)
			expiredCount++
		}
	}
	rc.lggr.Debugw("Removed expired cached HTTP responses", "count", expiredCount, "remaining", len(rc.cache))
	return expiredCount
}
