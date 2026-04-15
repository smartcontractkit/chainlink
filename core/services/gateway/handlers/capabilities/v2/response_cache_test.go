package v2

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	gateway_common "github.com/smartcontractkit/chainlink-common/pkg/types/gateway"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities/v2/metrics"
)

func createCacheTestMetrics(t *testing.T) *metrics.Metrics {
	m, err := metrics.NewMetrics(nil)
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
			MaxAgeMs: 5000, // Read from cache if cache entry is fresher than 5 seconds
			Store:    true, // Store responses in cache by default for tests
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

func TestRequestHash(t *testing.T) {
	req := createTestRequest("GET", "https://example.com")

	t.Run("generates consistent hash", func(t *testing.T) {
		hash1 := req.Hash()
		hash2 := req.Hash()
		require.Equal(t, hash1, hash2)
	})

	t.Run("different requests generate different hashes", func(t *testing.T) {
		req1 := createTestRequest("GET", "https://example.com/path1")
		req2 := createTestRequest("GET", "https://example.com/path2")

		hash1 := req1.Hash()
		hash2 := req2.Hash()
		require.NotEqual(t, hash1, hash2)
	})

	t.Run("same request with different method generates different hash", func(t *testing.T) {
		req1 := createTestRequest("GET", "https://example.com")
		req2 := createTestRequest("POST", "https://example.com")

		hash1 := req1.Hash()
		hash2 := req2.Hash()
		require.NotEqual(t, hash1, hash2)
	})

	t.Run("having different cacheSettings results in the same Hash", func(t *testing.T) {
		req1 := createTestRequest("GET", "https://example.com")
		req1.CacheSettings = gateway_common.CacheSettings{
			MaxAgeMs: 5000,
			Store:    true,
		}

		req2 := createTestRequest("GET", "https://example.com")
		req2.CacheSettings = gateway_common.CacheSettings{
			MaxAgeMs: 10000,
			Store:    false,
		}

		hash1 := req1.Hash()
		hash2 := req2.Hash()
		require.Equal(t, hash1, hash2, "Hash should be the same regardless of CacheSettings")
	})

	t.Run("having different workflowID results in same Hash", func(t *testing.T) {
		req1 := createTestRequest("GET", "https://example.com")
		req1.WorkflowID = "workflow-123"

		req2 := createTestRequest("GET", "https://example.com")
		req2.WorkflowID = "workflow-456"

		hash1 := req1.Hash()
		hash2 := req2.Hash()
		require.Equal(t, hash1, hash2, "Hash should be the same regardless of WorkflowID")
	})

	t.Run("having same workflowOwner results in the same Hash", func(t *testing.T) {
		req1 := createTestRequest("GET", "https://example.com")
		req1.WorkflowOwner = "workflow-owner-123"

		req2 := createTestRequest("GET", "https://example.com")
		req2.WorkflowOwner = "workflow-owner-123"

		hash1 := req1.Hash()
		hash2 := req2.Hash()
		require.Equal(t, hash1, hash2, "Hash should be the same for identical requests")
	})

	t.Run("having different workflowOwner results in different Hash", func(t *testing.T) {
		req1 := createTestRequest("GET", "https://example.com")
		req1.WorkflowOwner = "workflow-owner-123"

		req2 := createTestRequest("GET", "https://example.com")
		req2.WorkflowOwner = "workflow-owner-456"

		hash1 := req1.Hash()
		hash2 := req2.Hash()
		require.NotEqual(t, hash1, hash2, "Hash should be different for different workflow owner")
		require.NotEmpty(t, hash1, "Hash should not be empty")
		require.NotEmpty(t, hash2, "Hash should not be empty")
	})
}

func TestIsExpiredOrNotCached(t *testing.T) {
	testMetrics := createCacheTestMetrics(t)
	cache := newResponseCache(logger.Test(t), 1000, testMetrics) // 1 second TTL

	req := createTestRequest("GET", "https://example.com")

	t.Run("returns true for non-existent entry", func(t *testing.T) {
		result := cache.isExpiredOrNotCached(req)
		require.True(t, result)
	})

	t.Run("returns false for non-expired entry", func(t *testing.T) {
		cache.cache[req.Hash()] = &cachedResponse{
			response: createTestResponse(200, "test"),
			storedAt: time.Now(),
		}

		result := cache.isExpiredOrNotCached(req)
		require.False(t, result)
	})

	t.Run("returns true for expired entry", func(t *testing.T) {
		cache.cache[req.Hash()] = &cachedResponse{
			response: createTestResponse(200, "test"),
			storedAt: time.Now().Add(-2 * time.Second),
		}

		result := cache.isExpiredOrNotCached(req)
		require.True(t, result)
	})
}

func TestFetch(t *testing.T) {
	testMetrics := createCacheTestMetrics(t)
	cache := newResponseCache(logger.Test(t), 10000, testMetrics) // 10 seconds TTL

	t.Run("calls fetchFn when cache miss", func(t *testing.T) {
		req := createTestRequest("GET", "https://example.com/miss")
		expectedResp := createTestResponse(200, "fresh data")

		var fetchCalled bool
		fetchFn := func() gateway_common.OutboundHTTPResponse {
			fetchCalled = true
			return expectedResp
		}

		result := cache.Fetch(t.Context(), req, fetchFn, true)

		require.True(t, fetchCalled)
		require.Equal(t, expectedResp, result)
	})

	t.Run("returns cached response when cache hit", func(t *testing.T) {
		req := createTestRequest("GET", "https://example.com/hit")
		cachedResp := createTestResponse(200, "cached data")

		// Pre-populate cache
		cache.cache[req.Hash()] = &cachedResponse{
			response: cachedResp,
			storedAt: time.Now(),
		}

		var fetchCalled bool
		fetchFn := func() gateway_common.OutboundHTTPResponse {
			fetchCalled = true
			return createTestResponse(200, "should not be called")
		}

		result := cache.Fetch(t.Context(), req, fetchFn, true)

		require.False(t, fetchCalled, "fetchFn should not be called on cache hit")
		require.Equal(t, cachedResp, result)
	})

	t.Run("calls fetchFn when cached entry is expired by MaxAgeMs", func(t *testing.T) {
		req := createTestRequest("GET", "https://example.com/expired")
		req.CacheSettings.MaxAgeMs = 100

		cache.cache[req.Hash()] = &cachedResponse{
			response: createTestResponse(200, "old data"),
			storedAt: time.Now().Add(-200 * time.Millisecond),
		}

		expectedResp := createTestResponse(200, "fresh data")
		var fetchCalled bool
		fetchFn := func() gateway_common.OutboundHTTPResponse {
			fetchCalled = true
			return expectedResp
		}

		result := cache.Fetch(t.Context(), req, fetchFn, true)

		require.True(t, fetchCalled)
		require.Equal(t, expectedResp, result)
	})

	t.Run("caches cacheable responses when storeOnFetch is true", func(t *testing.T) {
		req := createTestRequest("GET", "https://example.com/cacheable")
		response := createTestResponse(200, "cacheable response")

		fetchFn := func() gateway_common.OutboundHTTPResponse {
			return response
		}

		cache.Fetch(t.Context(), req, fetchFn, true)

		cachedEntry, exists := cache.cache[req.Hash()]
		require.True(t, exists)
		require.Equal(t, response, cachedEntry.response)
	})

	t.Run("does not cache when storeOnFetch is false", func(t *testing.T) {
		req := createTestRequest("GET", "https://example.com/nostore")
		response := createTestResponse(200, "should not be stored")

		fetchFn := func() gateway_common.OutboundHTTPResponse {
			return response
		}

		result := cache.Fetch(t.Context(), req, fetchFn, false)

		// Should return the response but not cache it
		require.Equal(t, response, result)

		_, exists := cache.cache[req.Hash()]
		require.False(t, exists, "response should not be cached when storeOnFetch is false")
	})

	t.Run("does not cache non-cacheable responses", func(t *testing.T) {
		req := createTestRequest("GET", "https://example.com/noncacheable")
		response := createTestResponse(500, "server error")

		fetchFn := func() gateway_common.OutboundHTTPResponse {
			return response
		}

		result := cache.Fetch(t.Context(), req, fetchFn, true)

		// Should return the response but not cache it
		require.Equal(t, response, result)

		_, exists := cache.cache[req.Hash()]
		require.False(t, exists, "5xx response should not be cached")
	})
}

func TestFetch_ConcurrentDifferentKeys_RunInParallel(t *testing.T) {
	testMetrics := createCacheTestMetrics(t)
	cache := newResponseCache(logger.Test(t), 10000, testMetrics)

	const n = 10

	synctest.Test(t, func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(n)

		start := time.Now()
		for i := 0; i < n; i++ {
			go func(idx int) {
				defer wg.Done()
				req := createTestRequest("GET", fmt.Sprintf("https://example.com/parallel/%d", idx))
				// Non-uniform delays — synctest only cares about relative order.
				delay := time.Duration(idx+1) * time.Second
				fetchFn := func() gateway_common.OutboundHTTPResponse {
					time.Sleep(delay)
					return createTestResponse(200, fmt.Sprintf("response-%d", idx))
				}
				resp := cache.Fetch(t.Context(), req, fetchFn, true)
				assert.Equal(t, 200, resp.StatusCode)
			}(i)
		}
		wg.Wait()
		elapsed := time.Since(start)

		// If requests run in parallel, total time should be the longest delay.
		assert.Equal(t, time.Duration(n)*time.Second, elapsed,
			"concurrent fetches to different keys should run in parallel, took %v (expected %v)", elapsed, time.Duration(n)*time.Second)
	})
}

func TestFetch_ConcurrentSameKey_Deduplicated(t *testing.T) {
	testMetrics := createCacheTestMetrics(t)
	cache := newResponseCache(logger.Test(t), 10000, testMetrics)

	const n = 5
	var fetchCount atomic.Int32

	req := createTestRequest("GET", "https://example.com/dedup")
	expectedResp := createTestResponse(200, "deduplicated")

	synctest.Test(t, func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(n)

		for i := 0; i < n; i++ {
			go func() {
				defer wg.Done()
				fetchFn := func() gateway_common.OutboundHTTPResponse {
					fetchCount.Add(1)
					time.Sleep(100 * time.Millisecond)
					return expectedResp
				}
				resp := cache.Fetch(t.Context(), req, fetchFn, true)
				assert.Equal(t, expectedResp, resp)
			}()
		}
		wg.Wait()

		// Singleflight should deduplicate: only 1 fetchFn call for the same key
		assert.Equal(t, int32(1), fetchCount.Load(), "singleflight should deduplicate concurrent requests to the same key")
	})
}

// TestFetch_StoreOnFetch_FlightLeaderDecides documents that only the singleflight
// leader's storeOnFetch value is used. Waiters' storeOnFetch is silently ignored
// because only the leader's closure executes inside flight.Do.
// This is not a production bug (callers with storeOnFetch=false bypass Fetch entirely),
// but it makes the contract explicit.
func TestFetch_StoreOnFetch_FlightLeaderDecides(t *testing.T) {
	testMetrics := createCacheTestMetrics(t)
	cache := newResponseCache(logger.Test(t), 10000, testMetrics)

	req := createTestRequest("GET", "https://example.com/leader-store")
	expectedResp := createTestResponse(200, "response")
	var fetchCount atomic.Int32
	leaderRunning := make(chan struct{})

	synctest.Test(t, func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(2)

		// Goroutine A: flight leader with storeOnFetch=false.
		go func() {
			defer wg.Done()
			fetchFn := func() gateway_common.OutboundHTTPResponse {
				fetchCount.Add(1)
				close(leaderRunning)
				time.Sleep(100 * time.Millisecond)
				return expectedResp
			}
			resp := cache.Fetch(t.Context(), req, fetchFn, false)
			assert.Equal(t, expectedResp, resp)
		}()

		// Goroutine B: waiter with storeOnFetch=true.
		go func() {
			defer wg.Done()
			<-leaderRunning // ensure A is the flight leader
			fetchFn := func() gateway_common.OutboundHTTPResponse {
				fetchCount.Add(1)
				return expectedResp
			}
			resp := cache.Fetch(t.Context(), req, fetchFn, true)
			assert.Equal(t, expectedResp, resp)
		}()

		wg.Wait()
	})

	assert.Equal(t, int32(1), fetchCount.Load(),
		"singleflight should deduplicate: only leader's fetchFn runs")

	_, exists := cache.cache[req.Hash()]
	assert.False(t, exists,
		"result should NOT be cached: leader had storeOnFetch=false, waiter's storeOnFetch=true was ignored")
}

func TestFetch_PanicInFetchFn_PropagatedToCaller(t *testing.T) {
	testMetrics := createCacheTestMetrics(t)
	cache := newResponseCache(logger.Test(t), 10000, testMetrics)

	req := createTestRequest("GET", "https://example.com/panic")
	fetchFn := func() gateway_common.OutboundHTTPResponse {
		panic("unexpected error in HTTP callback")
	}

	require.Panics(t, func() {
		cache.Fetch(t.Context(), req, fetchFn, true)
	})
}

func TestFetch_PanicInFetchFn_PropagatedToAllWaiters(t *testing.T) {
	testMetrics := createCacheTestMetrics(t)
	cache := newResponseCache(logger.Test(t), 10000, testMetrics)

	const n = 5
	req := createTestRequest("GET", "https://example.com/panic-shared")

	synctest.Test(t, func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(n)

		for i := 0; i < n; i++ {
			go func() {
				defer wg.Done()
				defer func() {
					r := recover()
					assert.NotNil(t, r, "panic from fetchFn should propagate to all waiters")
				}()
				fetchFn := func() gateway_common.OutboundHTTPResponse {
					time.Sleep(50 * time.Millisecond)
					panic("shared panic")
				}
				cache.Fetch(t.Context(), req, fetchFn, true)
			}()
		}
		wg.Wait()
	})
}

func TestSet(t *testing.T) {
	testMetrics := createCacheTestMetrics(t)
	cache := newResponseCache(logger.Test(t), 10000, testMetrics)

	t.Run("sets cacheable response", func(t *testing.T) {
		req := createTestRequest("GET", "https://example.com/set")
		response := createTestResponse(200, "response to cache")

		cache.Set(req, response)

		cachedEntry, exists := cache.cache[req.Hash()]
		require.True(t, exists)
		require.Equal(t, response, cachedEntry.response)
	})

	t.Run("does not set non-cacheable response", func(t *testing.T) {
		req := createTestRequest("GET", "https://example.com/nonset")
		response := createTestResponse(500, "server error")

		cache.Set(req, response)

		_, exists := cache.cache[req.Hash()]
		require.False(t, exists, "5xx response should not be cached")
	})

	t.Run("does not overwrite non-expired entry", func(t *testing.T) {
		req := createTestRequest("GET", "https://example.com/nooverwrite")
		originalResponse := createTestResponse(200, "original")
		newResponse := createTestResponse(200, "new")

		cache.Set(req, originalResponse)

		// Immediately try to set again
		cache.Set(req, newResponse)

		cachedEntry, exists := cache.cache[req.Hash()]
		require.True(t, exists)
		require.Equal(t, originalResponse, cachedEntry.response)
	})

	t.Run("overwrites expired entry", func(t *testing.T) {
		req := createTestRequest("GET", "https://example.com/overwrite")

		cache.cache[req.Hash()] = &cachedResponse{
			response: createTestResponse(200, "expired"),
			storedAt: time.Now().Add(-20 * time.Second),
		}

		newResponse := createTestResponse(200, "fresh")
		cache.Set(req, newResponse)

		cachedEntry, exists := cache.cache[req.Hash()]
		require.True(t, exists)
		require.Equal(t, newResponse, cachedEntry.response)
	})
}

func TestDeleteExpired(t *testing.T) {
	testMetrics := createCacheTestMetrics(t)
	cache := newResponseCache(logger.Test(t), 1000, testMetrics)

	t.Run("deletes expired entries and returns count", func(t *testing.T) {
		expiredReq1 := createTestRequest("GET", "https://example.com/expired1")
		expiredReq2 := createTestRequest("GET", "https://example.com/expired2")
		validReq := createTestRequest("GET", "https://example.com/valid")

		expiredTime := time.Now().Add(-2 * time.Second)
		validTime := time.Now()

		cache.cache[expiredReq1.Hash()] = &cachedResponse{
			response: createTestResponse(200, "expired1"),
			storedAt: expiredTime,
		}
		cache.cache[expiredReq2.Hash()] = &cachedResponse{
			response: createTestResponse(200, "expired2"),
			storedAt: expiredTime,
		}
		cache.cache[validReq.Hash()] = &cachedResponse{
			response: createTestResponse(200, "valid"),
			storedAt: validTime,
		}

		count := cache.DeleteExpired(t.Context())

		require.Equal(t, 2, count, "should delete 2 expired entries")
		require.Len(t, cache.cache, 1, "should have 1 entry remaining")

		// Valid entry should still exist
		_, exists := cache.cache[validReq.Hash()]
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

		req := createTestRequest("GET", "https://example.com/zero-ttl")

		require.True(t, cache.isExpiredOrNotCached(req))

		cache.Set(req, createTestResponse(200, "test"))
		count := cache.DeleteExpired(t.Context())
		require.Equal(t, 1, count, "entry should be immediately expired")
	})

	t.Run("handles nil response headers", func(t *testing.T) {
		testMetrics := createCacheTestMetrics(t)
		cache := newResponseCache(logger.Test(t), 5000, testMetrics)

		req := createTestRequest("GET", "https://example.com/nil-headers")

		resp := gateway_common.OutboundHTTPResponse{
			StatusCode: 200,
			Body:       []byte("test"),
			Headers:    nil,
		}

		cache.Set(req, resp)

		result := cache.Fetch(t.Context(), req, func() gateway_common.OutboundHTTPResponse {
			return resp
		}, true)
		require.Equal(t, resp, result)
	})

	t.Run("handles empty request", func(t *testing.T) {
		testMetrics := createCacheTestMetrics(t)
		cache := newResponseCache(logger.Test(t), 5000, testMetrics)

		emptyReq := gateway_common.OutboundHTTPRequest{
			CacheSettings: gateway_common.CacheSettings{MaxAgeMs: 1000},
		}

		hash := emptyReq.Hash()
		require.NotEmpty(t, hash)

		cache.Set(emptyReq, createTestResponse(200, "test"))
	})
}
