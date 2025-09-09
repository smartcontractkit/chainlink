package v2

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	gateway_common "github.com/smartcontractkit/chainlink-common/pkg/types/gateway"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities/v2/metrics"
)

func createCacheTestMetrics(t *testing.T) *metrics.Metrics {
	m, err := metrics.NewMetrics()
	require.NoError(t, err)
	return m
}

func createTestRequest(method, url string) gateway_common.OutboundHTTPRequest {
	return gateway_common.OutboundHTTPRequest{
		Method: method,
		URL:    url,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: []byte(`{"test": "data"}`),
		CacheSettings: gateway_common.CacheSettings{
			ReadFromCache: true,
			MaxAgeMs:      5000, // 5 seconds
		},
	}
}

func createTestResponse(statusCode int, body string) gateway_common.OutboundHTTPResponse {
	return gateway_common.OutboundHTTPResponse{
		StatusCode: statusCode,
		Body:       []byte(body),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
}

func TestIsCacheableStatusCode(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		expected   bool
	}{
		// 2xx status codes - cacheable
		{"200 OK", 200, true},
		{"201 Created", 201, true},
		{"204 No Content", 204, true},

		// 4xx status codes - cacheable
		{"400 Bad Request", 400, true},
		{"401 Unauthorized", 401, true},
		{"404 Not Found", 404, true},

		// 1xx status codes - not cacheable
		{"100 Continue", 100, false},
		{"102 Processing", 102, false},
		{"199 Custom 1xx", 199, false},

		// 3xx status codes - not cacheable
		{"300 Multiple Choices", 300, false},
		{"301 Moved Permanently", 301, false},
		{"399 Custom 3xx", 399, false},

		// 5xx status codes - not cacheable
		{"500 Internal Server Error", 500, false},
		{"502 Bad Gateway", 502, false},
		{"599 Custom 5xx", 599, false},

		// Edge cases
		{"0 Invalid", 0, false},
		{"600 Invalid", 600, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isCacheableStatusCode(tt.statusCode)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestCacheKey(t *testing.T) {
	req := createTestRequest("GET", "https://example.com")
	workflowID := "workflow-123"

	t.Run("generates consistent key", func(t *testing.T) {
		key1 := cacheKey(workflowID, req)
		key2 := cacheKey(workflowID, req)
		require.Equal(t, key1, key2)
	})

	t.Run("includes workflow ID", func(t *testing.T) {
		key := cacheKey(workflowID, req)
		require.Contains(t, key, workflowID)
	})

	t.Run("different workflow IDs generate different keys", func(t *testing.T) {
		key1 := cacheKey("workflow-1", req)
		key2 := cacheKey("workflow-2", req)
		require.NotEqual(t, key1, key2)
	})

	t.Run("different requests generate different keys", func(t *testing.T) {
		req1 := createTestRequest("GET", "https://example.com/path1")
		req2 := createTestRequest("GET", "https://example.com/path2")

		key1 := cacheKey(workflowID, req1)
		key2 := cacheKey(workflowID, req2)
		require.NotEqual(t, key1, key2)
	})
}

func TestIsExpiredOrNotCached(t *testing.T) {
	testMetrics := createCacheTestMetrics(t)
	cache := newResponseCache(logger.Test(t), 1000, testMetrics) // 1 second TTL
	workflowID := "workflow-123"
	req := createTestRequest("GET", "https://example.com")

	t.Run("returns true for non-existent entry", func(t *testing.T) {
		result := cache.isExpiredOrNotCached(workflowID, req)
		require.True(t, result)
	})

	t.Run("returns false for non-expired entry", func(t *testing.T) {
		key := cacheKey(workflowID, req)
		cache.cache[key] = &cachedResponse{
			response: createTestResponse(200, "test"),
			storedAt: time.Now(),
		}

		result := cache.isExpiredOrNotCached(workflowID, req)
		require.False(t, result)
	})

	t.Run("returns true for expired entry", func(t *testing.T) {
		key := cacheKey(workflowID, req)
		cache.cache[key] = &cachedResponse{
			response: createTestResponse(200, "test"),
			storedAt: time.Now().Add(-2 * time.Second),
		}

		result := cache.isExpiredOrNotCached(workflowID, req)
		require.True(t, result)
	})
}

func TestCachedFetch(t *testing.T) {
	testMetrics := createCacheTestMetrics(t)
	cache := newResponseCache(logger.Test(t), 10000, testMetrics) // 10 seconds TTL
	workflowID := "workflow-123"

	t.Run("calls fetchFn when cache miss", func(t *testing.T) {
		req := createTestRequest("GET", "https://example.com/miss")
		expectedResp := createTestResponse(200, "fresh data")

		var fetchCalled bool
		fetchFn := func() gateway_common.OutboundHTTPResponse {
			fetchCalled = true
			return expectedResp
		}

		result := cache.CachedFetch(t.Context(), workflowID, req, fetchFn)

		require.True(t, fetchCalled)
		require.Equal(t, expectedResp, result)
	})

	t.Run("returns cached response when cache hit", func(t *testing.T) {
		req := createTestRequest("GET", "https://example.com/hit")
		cachedResp := createTestResponse(200, "cached data")

		// Pre-populate cache
		key := cacheKey(workflowID, req)
		cache.cache[key] = &cachedResponse{
			response: cachedResp,
			storedAt: time.Now(),
		}

		var fetchCalled bool
		fetchFn := func() gateway_common.OutboundHTTPResponse {
			fetchCalled = true
			return createTestResponse(200, "should not be called")
		}

		result := cache.CachedFetch(t.Context(), workflowID, req, fetchFn)

		require.False(t, fetchCalled, "fetchFn should not be called on cache hit")
		require.Equal(t, cachedResp, result)
	})

	t.Run("calls fetchFn when cached entry is expired by MaxAgeMs", func(t *testing.T) {
		req := createTestRequest("GET", "https://example.com/expired")
		req.CacheSettings.MaxAgeMs = 100

		key := cacheKey(workflowID, req)
		cache.cache[key] = &cachedResponse{
			response: createTestResponse(200, "old data"),
			storedAt: time.Now().Add(-200 * time.Millisecond),
		}

		expectedResp := createTestResponse(200, "fresh data")
		var fetchCalled bool
		fetchFn := func() gateway_common.OutboundHTTPResponse {
			fetchCalled = true
			return expectedResp
		}

		result := cache.CachedFetch(t.Context(), workflowID, req, fetchFn)

		require.True(t, fetchCalled)
		require.Equal(t, expectedResp, result)
	})

	t.Run("caches cacheable responses", func(t *testing.T) {
		req := createTestRequest("GET", "https://example.com/cacheable")
		response := createTestResponse(200, "cacheable response")

		fetchFn := func() gateway_common.OutboundHTTPResponse {
			return response
		}

		cache.CachedFetch(t.Context(), workflowID, req, fetchFn)

		key := cacheKey(workflowID, req)
		cachedEntry, exists := cache.cache[key]
		require.True(t, exists)
		require.Equal(t, response, cachedEntry.response)
	})

	t.Run("does not cache non-cacheable responses", func(t *testing.T) {
		req := createTestRequest("GET", "https://example.com/noncacheable")
		response := createTestResponse(500, "server error")

		fetchFn := func() gateway_common.OutboundHTTPResponse {
			return response
		}

		result := cache.CachedFetch(t.Context(), workflowID, req, fetchFn)

		// Should return the response but not cache it
		require.Equal(t, response, result)

		key := cacheKey(workflowID, req)
		_, exists := cache.cache[key]
		require.False(t, exists, "5xx response should not be cached")
	})
}

func TestSet(t *testing.T) {
	testMetrics := createCacheTestMetrics(t)
	cache := newResponseCache(logger.Test(t), 10000, testMetrics)
	workflowID := "workflow-123"

	t.Run("sets cacheable response", func(t *testing.T) {
		req := createTestRequest("GET", "https://example.com/set")
		response := createTestResponse(200, "response to cache")

		cache.Set(workflowID, req, response)

		key := cacheKey(workflowID, req)
		cachedEntry, exists := cache.cache[key]
		require.True(t, exists)
		require.Equal(t, response, cachedEntry.response)
	})

	t.Run("does not set non-cacheable response", func(t *testing.T) {
		req := createTestRequest("GET", "https://example.com/nonset")
		response := createTestResponse(500, "server error")

		cache.Set(workflowID, req, response)

		key := cacheKey(workflowID, req)
		_, exists := cache.cache[key]
		require.False(t, exists, "5xx response should not be cached")
	})

	t.Run("does not overwrite non-expired entry", func(t *testing.T) {
		req := createTestRequest("GET", "https://example.com/nooverwrite")
		originalResponse := createTestResponse(200, "original")
		newResponse := createTestResponse(200, "new")

		cache.Set(workflowID, req, originalResponse)

		// Immediately try to set again
		cache.Set(workflowID, req, newResponse)

		key := cacheKey(workflowID, req)
		cachedEntry, exists := cache.cache[key]
		require.True(t, exists)
		require.Equal(t, originalResponse, cachedEntry.response)
	})

	t.Run("overwrites expired entry", func(t *testing.T) {
		req := createTestRequest("GET", "https://example.com/overwrite")

		key := cacheKey(workflowID, req)
		cache.cache[key] = &cachedResponse{
			response: createTestResponse(200, "expired"),
			storedAt: time.Now().Add(-20 * time.Second),
		}

		newResponse := createTestResponse(200, "fresh")
		cache.Set(workflowID, req, newResponse)

		cachedEntry, exists := cache.cache[key]
		require.True(t, exists)
		require.Equal(t, newResponse, cachedEntry.response)
	})
}

func TestDeleteExpired(t *testing.T) {
	testMetrics := createCacheTestMetrics(t)
	cache := newResponseCache(logger.Test(t), 1000, testMetrics)

	t.Run("deletes expired entries and returns count", func(t *testing.T) {
		workflowID := "workflow-123"

		expiredReq1 := createTestRequest("GET", "https://example.com/expired1")
		expiredReq2 := createTestRequest("GET", "https://example.com/expired2")
		validReq := createTestRequest("GET", "https://example.com/valid")

		expiredTime := time.Now().Add(-2 * time.Second)
		validTime := time.Now()

		cache.cache[cacheKey(workflowID, expiredReq1)] = &cachedResponse{
			response: createTestResponse(200, "expired1"),
			storedAt: expiredTime,
		}
		cache.cache[cacheKey(workflowID, expiredReq2)] = &cachedResponse{
			response: createTestResponse(200, "expired2"),
			storedAt: expiredTime,
		}
		cache.cache[cacheKey(workflowID, validReq)] = &cachedResponse{
			response: createTestResponse(200, "valid"),
			storedAt: validTime,
		}

		count := cache.DeleteExpired(t.Context())

		require.Equal(t, 2, count, "should delete 2 expired entries")
		require.Len(t, cache.cache, 1, "should have 1 entry remaining")

		// Valid entry should still exist
		_, exists := cache.cache[cacheKey(workflowID, validReq)]
		require.True(t, exists)
	})

	t.Run("returns zero when cache is empty", func(t *testing.T) {
		testMetrics := createCacheTestMetrics(t)
		emptyCache := newResponseCache(logger.Test(t), 1000, testMetrics)
		count := emptyCache.DeleteExpired(t.Context())
		require.Equal(t, 0, count)
	})
}

func TestEdgeCases(t *testing.T) {
	t.Run("zero TTL cache", func(t *testing.T) {
		testMetrics := createCacheTestMetrics(t)
		cache := newResponseCache(logger.Test(t), 0, testMetrics)
		workflowID := "workflow-123"
		req := createTestRequest("GET", "https://example.com/zero-ttl")

		require.True(t, cache.isExpiredOrNotCached(workflowID, req))

		cache.Set(workflowID, req, createTestResponse(200, "test"))
		count := cache.DeleteExpired(t.Context())
		require.Equal(t, 1, count, "entry should be immediately expired")
	})

	t.Run("handles nil response headers", func(t *testing.T) {
		testMetrics := createCacheTestMetrics(t)
		cache := newResponseCache(logger.Test(t), 5000, testMetrics)
		workflowID := "workflow-123"
		req := createTestRequest("GET", "https://example.com/nil-headers")

		resp := gateway_common.OutboundHTTPResponse{
			StatusCode: 200,
			Body:       []byte("test"),
			Headers:    nil,
		}

		cache.Set(workflowID, req, resp)

		result := cache.CachedFetch(t.Context(), workflowID, req, func() gateway_common.OutboundHTTPResponse {
			return resp
		})
		require.Equal(t, resp, result)
	})

	t.Run("handles empty request", func(t *testing.T) {
		testMetrics := createCacheTestMetrics(t)
		cache := newResponseCache(logger.Test(t), 5000, testMetrics)
		workflowID := "workflow-123"

		emptyReq := gateway_common.OutboundHTTPRequest{
			CacheSettings: gateway_common.CacheSettings{MaxAgeMs: 1000},
		}

		key := cacheKey(workflowID, emptyReq)
		require.NotEmpty(t, key)

		cache.Set(workflowID, emptyReq, createTestResponse(200, "test"))
	})
}
