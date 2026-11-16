package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func rateLimitEngine(period time.Duration, limit int64) *gin.Engine {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(rateLimiter(period, limit))
	engine.GET("/", func(c *gin.Context) { c.String(http.StatusOK, "ok") })
	return engine
}

func serveRateLimit(t *testing.T, engine *gin.Engine, remoteAddr string) *httptest.ResponseRecorder {
	t.Helper()
	request, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com/", nil)
	require.NoError(t, err)
	request.RemoteAddr = remoteAddr
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, request)
	return recorder
}

func TestLocalRateLimiter_AllowsThenRejects(t *testing.T) {
	engine := rateLimitEngine(time.Minute, 2)

	first := serveRateLimit(t, engine, "192.0.2.1:1234")
	assert.Equal(t, http.StatusOK, first.Code)
	assert.Equal(t, "2", first.Header().Get("X-RateLimit-Limit"))
	assert.Equal(t, "1", first.Header().Get("X-RateLimit-Remaining"))
	assert.NotEmpty(t, first.Header().Get("X-RateLimit-Reset"))

	second := serveRateLimit(t, engine, "192.0.2.1:1234")
	assert.Equal(t, http.StatusOK, second.Code)
	assert.Equal(t, "0", second.Header().Get("X-RateLimit-Remaining"))

	third := serveRateLimit(t, engine, "192.0.2.1:1234")
	assert.Equal(t, http.StatusTooManyRequests, third.Code)
	assert.Equal(t, "Limit exceeded", third.Body.String())
	assert.Equal(t, "0", third.Header().Get("X-RateLimit-Remaining"))
}

func TestLocalRateLimiter_KeysByClientIP(t *testing.T) {
	engine := rateLimitEngine(time.Minute, 1)

	assert.Equal(t, http.StatusOK, serveRateLimit(t, engine, "192.0.2.1:1234").Code)
	assert.Equal(t, http.StatusTooManyRequests, serveRateLimit(t, engine, "192.0.2.1:5678").Code)
	assert.Equal(t, http.StatusOK, serveRateLimit(t, engine, "192.0.2.2:1234").Code)
}

func TestLocalRateLimiter_ResetsAfterPeriod(t *testing.T) {
	limiter := newLocalRateLimiter(20*time.Millisecond, 1)

	remaining, _, reached := limiter.increment("192.0.2.1")
	assert.False(t, reached)
	assert.Equal(t, int64(0), remaining)

	_, _, reached = limiter.increment("192.0.2.1")
	assert.True(t, reached)

	time.Sleep(25 * time.Millisecond)
	remaining, _, reached = limiter.increment("192.0.2.1")
	assert.False(t, reached)
	assert.Equal(t, int64(0), remaining)
}
