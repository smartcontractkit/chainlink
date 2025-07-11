package v2

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/gateway"
)

func TestResponseCache_SetAndGet(t *testing.T) {
	cache := newResponseCache(logger.Test(t))

	req := gateway.OutboundHTTPRequest{
		Method: "GET",
		URL:    "https://example.com/api",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: []byte(`{"test": "data"}`),
	}

	resp := gateway.OutboundHTTPResponse{
		StatusCode: 200,
		Body:       []byte(`{"result": "success"}`),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	t.Run("set and get within TTL", func(t *testing.T) {
		ttl := 1 * time.Hour
		cache.Set(req, resp, ttl)

		cached := cache.Get(req)
		require.NotNil(t, cached)
		assert.Equal(t, resp.StatusCode, cached.StatusCode)
		assert.Equal(t, resp.Body, cached.Body)
		assert.Equal(t, resp.Headers, cached.Headers)
	})

	t.Run("get returns nil for non-existent key", func(t *testing.T) {
		nonExistentReq := gateway.OutboundHTTPRequest{
			Method: "POST",
			URL:    "https://different.com/api",
		}

		cached := cache.Get(nonExistentReq)
		assert.Nil(t, cached)
	})

	t.Run("different requests have different cache keys", func(t *testing.T) {
		req1 := gateway.OutboundHTTPRequest{
			Method: "GET",
			URL:    "https://example.com/api/v1",
		}

		req2 := gateway.OutboundHTTPRequest{
			Method: "GET",
			URL:    "https://example.com/api/v2",
		}

		resp1 := gateway.OutboundHTTPResponse{StatusCode: 200, Body: []byte("response1")}
		resp2 := gateway.OutboundHTTPResponse{StatusCode: 201, Body: []byte("response2")}

		cache.Set(req1, resp1, time.Hour)
		cache.Set(req2, resp2, time.Hour)

		cached1 := cache.Get(req1)
		cached2 := cache.Get(req2)

		require.NotNil(t, cached1)
		require.NotNil(t, cached2)
		assert.Equal(t, []byte("response1"), cached1.Body)
		assert.Equal(t, []byte("response2"), cached2.Body)
	})
}

func TestResponseCache_Expiration(t *testing.T) {
	cache := newResponseCache(logger.Test(t))

	req := gateway.OutboundHTTPRequest{
		Method: "GET",
		URL:    "https://example.com/api",
	}

	resp := gateway.OutboundHTTPResponse{
		StatusCode: 200,
		Body:       []byte("test response"),
	}

	t.Run("expired entries return nil", func(t *testing.T) {
		// Set with very short TTL
		ttl := 10 * time.Millisecond
		cache.Set(req, resp, ttl)

		// Verify it's initially available
		cached := cache.Get(req)
		require.NotNil(t, cached)

		// Wait for expiration
		time.Sleep(1000 * time.Millisecond)

		// Should now return nil
		cached = cache.Get(req)
		assert.Nil(t, cached)
	})
}

func TestResponseCache_DeleteExpired(t *testing.T) {
	cache := newResponseCache(logger.Test(t))

	t.Run("deletes expired entries", func(t *testing.T) {
		// Add some entries with different TTLs
		req1 := gateway.OutboundHTTPRequest{Method: "GET", URL: "https://example.com/1"}
		req2 := gateway.OutboundHTTPRequest{Method: "GET", URL: "https://example.com/2"}
		req3 := gateway.OutboundHTTPRequest{Method: "GET", URL: "https://example.com/3"}

		resp := gateway.OutboundHTTPResponse{StatusCode: 200}

		// Set entries with different TTLs
		cache.Set(req1, resp, 10*time.Millisecond) // Will expire quickly
		cache.Set(req2, resp, 10*time.Millisecond) // Will expire quickly
		cache.Set(req3, resp, 1*time.Hour)         // Will not expire

		// Wait for some to expire
		time.Sleep(1000 * time.Millisecond)

		// Delete expired entries
		deletedCount := cache.DeleteExpired()

		// Should have deleted 2 expired entries
		assert.Equal(t, 2, deletedCount)

		// Non-expired entry should still be available
		cached := cache.Get(req3)
		assert.NotNil(t, cached)

		// Expired entries should not be available
		cached1 := cache.Get(req1)
		cached2 := cache.Get(req2)
		assert.Nil(t, cached1)
		assert.Nil(t, cached2)
	})

	t.Run("returns zero when no entries to delete", func(t *testing.T) {
		// Clear cache by creating new one
		cache = newResponseCache(logger.Test(t))

		deletedCount := cache.DeleteExpired()
		assert.Equal(t, 0, deletedCount)
	})
}

func TestResponseCache_KeyGeneration(t *testing.T) {
	cache := newResponseCache(logger.Test(t))

	resp := gateway.OutboundHTTPResponse{StatusCode: 200}

	t.Run("different methods generate different keys", func(t *testing.T) {
		req1 := gateway.OutboundHTTPRequest{Method: "GET", URL: "https://example.com"}
		req2 := gateway.OutboundHTTPRequest{Method: "POST", URL: "https://example.com"}

		cache.Set(req1, resp, time.Hour)
		cache.Set(req2, resp, time.Hour)

		// Both should be stored separately
		cached1 := cache.Get(req1)
		cached2 := cache.Get(req2)

		assert.NotNil(t, cached1)
		assert.NotNil(t, cached2)
	})

	t.Run("different URLs generate different keys", func(t *testing.T) {
		req1 := gateway.OutboundHTTPRequest{Method: "GET", URL: "https://example.com/path1"}
		req2 := gateway.OutboundHTTPRequest{Method: "GET", URL: "https://example.com/path2"}

		cache.Set(req1, resp, time.Hour)
		cache.Set(req2, resp, time.Hour)

		// Both should be stored separately
		cached1 := cache.Get(req1)
		cached2 := cache.Get(req2)

		assert.NotNil(t, cached1)
		assert.NotNil(t, cached2)
	})

	t.Run("different bodies generate different keys", func(t *testing.T) {
		req1 := gateway.OutboundHTTPRequest{
			Method: "POST",
			URL:    "https://example.com",
			Body:   []byte(`{"key": "value1"}`),
		}
		req2 := gateway.OutboundHTTPRequest{
			Method: "POST",
			URL:    "https://example.com",
			Body:   []byte(`{"key": "value2"}`),
		}

		cache.Set(req1, resp, time.Hour)
		cache.Set(req2, resp, time.Hour)

		// Both should be stored separately
		cached1 := cache.Get(req1)
		cached2 := cache.Get(req2)

		assert.NotNil(t, cached1)
		assert.NotNil(t, cached2)
	})

	t.Run("certain headers affect cache key", func(t *testing.T) {
		req1 := gateway.OutboundHTTPRequest{
			Method: "GET",
			URL:    "https://example.com",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
		req2 := gateway.OutboundHTTPRequest{
			Method: "GET",
			URL:    "https://example.com",
			Headers: map[string]string{
				"Content-Type": "application/xml",
			},
		}

		cache.Set(req1, resp, time.Hour)
		cache.Set(req2, resp, time.Hour)

		// Both should be stored separately if headers affect cache key
		cached1 := cache.Get(req1)
		cached2 := cache.Get(req2)

		assert.NotNil(t, cached1)
		assert.NotNil(t, cached2)
	})
}
