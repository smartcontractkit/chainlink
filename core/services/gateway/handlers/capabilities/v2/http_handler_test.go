package v2

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/ratelimit"
	"github.com/smartcontractkit/chainlink-common/pkg/types/gateway"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	handlermocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
	httpmocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/network/mocks"
)

func TestNewGatewayHandler(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		cfg := serviceCfg()
		configBytes, err := json.Marshal(cfg)
		require.NoError(t, err)

		donConfig := &config.DONConfig{
			DonId: "test-don",
		}
		mockDon := handlermocks.NewDON(t)
		mockHTTPClient := httpmocks.NewHTTPClient(t)
		lggr := logger.Test(t)

		handler, err := NewGatewayHandler(configBytes, donConfig, mockDon, mockHTTPClient, lggr)
		require.NoError(t, err)
		require.NotNil(t, handler)
		require.Equal(t, "test-don", handler.donConfig.DonId)
		require.NotNil(t, handler.responseCache)
	})

	t.Run("invalid config JSON", func(t *testing.T) {
		invalidConfig := []byte(`{invalid json}`)
		donConfig := &config.DONConfig{DonId: "test-don"}
		mockDon := handlermocks.NewDON(t)
		mockHTTPClient := httpmocks.NewHTTPClient(t)
		lggr := logger.Test(t)

		handler, err := NewGatewayHandler(invalidConfig, donConfig, mockDon, mockHTTPClient, lggr)
		require.Error(t, err)
		require.Nil(t, handler)
	})

	t.Run("invalid rate limiter config", func(t *testing.T) {
		cfg := ServiceConfig{
			NodeRateLimiter: ratelimit.RateLimiterConfig{
				GlobalRPS:   -1, // Invalid negative rate
				GlobalBurst: 100,
			},
			UserRateLimiter: ratelimit.RateLimiterConfig{
				GlobalRPS:   50,
				GlobalBurst: 50,
			},
		}
		configBytes, err := json.Marshal(cfg)
		require.NoError(t, err)

		donConfig := &config.DONConfig{DonId: "test-don"}
		mockDon := handlermocks.NewDON(t)
		mockHTTPClient := httpmocks.NewHTTPClient(t)
		lggr := logger.Test(t)

		handler, err := NewGatewayHandler(configBytes, donConfig, mockDon, mockHTTPClient, lggr)
		require.Error(t, err)
		require.Nil(t, handler)
	})

	t.Run("applies default config values", func(t *testing.T) {
		cfg := ServiceConfig{
			NodeRateLimiter: ratelimit.RateLimiterConfig{
				GlobalRPS:      100,
				GlobalBurst:    100,
				PerSenderRPS:   10,
				PerSenderBurst: 10,
			},
			UserRateLimiter: ratelimit.RateLimiterConfig{
				GlobalRPS:      50,
				GlobalBurst:    50,
				PerSenderRPS:   5,
				PerSenderBurst: 5,
			},
			// CleanUpPeriodMs not set - should get default
		}
		configBytes, err := json.Marshal(cfg)
		require.NoError(t, err)

		donConfig := &config.DONConfig{DonId: "test-don"}
		mockDon := handlermocks.NewDON(t)
		mockHTTPClient := httpmocks.NewHTTPClient(t)
		lggr := logger.Test(t)

		handler, err := NewGatewayHandler(configBytes, donConfig, mockDon, mockHTTPClient, lggr)
		require.NoError(t, err)
		require.NotNil(t, handler)
		require.Equal(t, defaultCleanUpPeriodMs, handler.config.CleanUpPeriodMs) // Default value
	})
}

func TestHandleNodeMessage(t *testing.T) {
	handler := createTestHandler(t)

	t.Run("successful node message handling", func(t *testing.T) {
		mockDon := handler.don.(*handlermocks.DON)
		mockHTTPClient := handler.httpClient.(*httpmocks.HTTPClient)

		// Prepare outbound request
		outboundReq := gateway.OutboundHTTPRequest{
			Method:        "GET",
			URL:           "https://example.com/api",
			TimeoutMs:     5000,
			Headers:       map[string]string{"Content-Type": "application/json"},
			Body:          []byte(`{"test": "data"}`),
			CacheSettings: gateway.CacheSettings{},
		}
		reqBytes, err := json.Marshal(outboundReq)
		require.NoError(t, err)
		rawRequest := json.RawMessage(reqBytes)
		resp := &jsonrpc.Response[json.RawMessage]{
			ID:     "test-request-id",
			Result: &rawRequest,
		}

		httpResp := &network.HTTPResponse{
			StatusCode: 200,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       []byte(`{"result": "success"}`),
		}
		mockHTTPClient.EXPECT().Send(mock.Anything, mock.MatchedBy(func(req network.HTTPRequest) bool {
			return req.Method == "GET" && req.URL == "https://example.com/api"
		})).Return(httpResp, nil)

		mockDon.EXPECT().SendToNode(mock.Anything, "node1", mock.MatchedBy(func(req *jsonrpc.Request[json.RawMessage]) bool {
			return req.ID == "test-request-id"
		})).Return(nil)

		err = handler.HandleNodeMessage(testutils.Context(t), resp, "node1")
		require.NoError(t, err)
		handler.wg.Wait()
	})

	t.Run("returns cached response if available", func(t *testing.T) {
		outboundReq := gateway.OutboundHTTPRequest{
			Method:    "GET",
			URL:       "https://return-cached.com/api",
			TimeoutMs: 5000,
			CacheSettings: gateway.CacheSettings{
				StoreInCache:  true,
				ReadFromCache: true,
				TTLMs:         600000, // 10 minute TTL
			},
		}
		reqBytes, err := json.Marshal(outboundReq)
		require.NoError(t, err)
		rawRequest := json.RawMessage(reqBytes)
		resp := &jsonrpc.Response[json.RawMessage]{
			ID:     "test-request-id",
			Result: &rawRequest,
		}

		mockDon := handler.don.(*handlermocks.DON)
		// First call: should fetch from HTTP client and cache the response
		httpResp := &network.HTTPResponse{
			StatusCode: 200,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       []byte(`{"cached": "response"}`),
		}
		mockHTTPClient := handler.httpClient.(*httpmocks.HTTPClient)
		mockHTTPClient.EXPECT().Send(mock.Anything, mock.Anything).Return(httpResp, nil).Once()
		mockDon.EXPECT().SendToNode(mock.Anything, "node1", mock.Anything).Return(nil)

		err = handler.HandleNodeMessage(testutils.Context(t), resp, "node1")
		require.NoError(t, err)
		handler.wg.Wait()

		// Second call: should return cached response (no HTTP client call)
		mockDon.EXPECT().SendToNode(mock.Anything, "node1", mock.MatchedBy(func(req *jsonrpc.Request[json.RawMessage]) bool {
			var cached gateway.OutboundHTTPResponse
			err2 := json.Unmarshal(*req.Params, &cached)
			return err2 == nil && string(cached.Body) == string(httpResp.Body)
		})).Return(nil)

		err = handler.HandleNodeMessage(testutils.Context(t), resp, "node1")
		require.NoError(t, err)
		handler.wg.Wait()
	})

	t.Run("status code 500 is not cached if StoreInCache is false", func(t *testing.T) {
		outboundReq := gateway.OutboundHTTPRequest{
			Method:    "GET",
			URL:       "https://status-500.com/api",
			TimeoutMs: 5000,
			CacheSettings: gateway.CacheSettings{
				StoreInCache:  true,
				ReadFromCache: true,
				TTLMs:         600000,
			},
		}
		reqBytes, err := json.Marshal(outboundReq)
		require.NoError(t, err)
		rawRequest := json.RawMessage(reqBytes)
		resp := &jsonrpc.Response[json.RawMessage]{
			ID:     "test-request-id-500",
			Result: &rawRequest,
		}

		mockDon := handler.don.(*handlermocks.DON)
		mockHTTPClient := handler.httpClient.(*httpmocks.HTTPClient)
		httpResp := &network.HTTPResponse{
			StatusCode: 500,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       []byte(`{"error": "bad request"}`),
		}
		mockHTTPClient.EXPECT().Send(mock.Anything, mock.Anything).Return(httpResp, nil).Once()
		mockDon.EXPECT().SendToNode(mock.Anything, "node1", mock.Anything).Return(nil)

		// First call: should fetch from HTTP client, but not cache the response
		err = handler.HandleNodeMessage(testutils.Context(t), resp, "node1")
		require.NoError(t, err)
		handler.wg.Wait()

		// Second call: should NOT return cached response, so HTTP client is called again
		mockHTTPClient.EXPECT().Send(mock.Anything, mock.Anything).Return(httpResp, nil).Once()
		mockDon.EXPECT().SendToNode(mock.Anything, "node1", mock.Anything).Return(nil)

		err = handler.HandleNodeMessage(testutils.Context(t), resp, "node1")
		require.NoError(t, err)
		handler.wg.Wait()
	})

	t.Run("empty request ID", func(t *testing.T) {
		rawRes := json.RawMessage([]byte(`{}`))
		resp := &jsonrpc.Response[json.RawMessage]{
			ID:     "",
			Result: &rawRes,
		}

		err := handler.HandleNodeMessage(testutils.Context(t), resp, "node1")
		require.Error(t, err)
		require.Contains(t, err.Error(), "empty request ID")
		handler.wg.Wait()
	})

	t.Run("invalid JSON in response result", func(t *testing.T) {
		rawRes := json.RawMessage([]byte(`{invalid json}`))
		resp := &jsonrpc.Response[json.RawMessage]{
			ID:     "test-request-id2",
			Result: &rawRes,
		}

		err := handler.HandleNodeMessage(testutils.Context(t), resp, "node1")
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to unmarshal HTTP request")
		handler.wg.Wait()
	})
}

func TestIsCacheableStatusCode(t *testing.T) {
	tests := []struct {
		statusCode int
		expected   bool
	}{
		{200, true},  // Success
		{201, true},  // Created
		{299, true},  // Last 2xx
		{300, false}, // Redirect (not cacheable)
		{400, true},  // Bad Request (cacheable)
		{404, true},  // Not Found (cacheable)
		{499, true},  // Last 4xx
		{500, false}, // Server Error (not cacheable)
		{503, false}, // Service Unavailable (not cacheable)
		{100, false}, // Informational (not cacheable)
		{600, false}, // Invalid status code
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("status_%d", tt.statusCode), func(t *testing.T) {
			result := isCacheableStatusCode(tt.statusCode)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestServiceLifecycle(t *testing.T) {
	handler := createTestHandler(t)

	t.Run("start and stop", func(t *testing.T) {
		ctx := testutils.Context(t)

		err := handler.Start(ctx)
		require.NoError(t, err)

		healthReport := handler.HealthReport()
		require.NoError(t, healthReport[handlerName])

		require.Equal(t, handlerName, handler.Name())

		err = handler.Close()
		require.NoError(t, err)
	})
}

type mockResponseCache struct {
	deleteExpiredCh chan struct{}
}

func newMockResponseCache() *mockResponseCache {
	return &mockResponseCache{
		deleteExpiredCh: make(chan struct{}),
	}
}

func (m *mockResponseCache) Set(gateway.OutboundHTTPRequest, gateway.OutboundHTTPResponse, time.Duration) {
}

func (m *mockResponseCache) Get(gateway.OutboundHTTPRequest) *gateway.OutboundHTTPResponse {
	return nil
}

func (m *mockResponseCache) DeleteExpired() int {
	select {
	case m.deleteExpiredCh <- struct{}{}:
	default:
	}
	return 0
}

func TestGatewayHandler_Start_CallsDeleteExpired(t *testing.T) {
	cfg := serviceCfg()
	cfg.CleanUpPeriodMs = 100 // fast cleanup for test

	configBytes, err := json.Marshal(cfg)
	require.NoError(t, err)

	donConfig := &config.DONConfig{DonId: "test-don"}
	mockDon := handlermocks.NewDON(t)
	mockHTTPClient := httpmocks.NewHTTPClient(t)
	lggr := logger.Test(t)

	handler, err := NewGatewayHandler(configBytes, donConfig, mockDon, mockHTTPClient, lggr)
	require.NoError(t, err)
	require.NotNil(t, handler)
	mockCache := newMockResponseCache()
	handler.responseCache = mockCache

	ctx := t.Context()
	err = handler.Start(ctx)
	require.NoError(t, err)

	// Wait for DeleteExpired to be called at least once
	select {
	case <-mockCache.deleteExpiredCh:
		// Success
	case <-ctx.Done():
		t.Fatal("DeleteExpired was not called within context deadline")
	}
	err = handler.Close()
	require.NoError(t, err)
}

func serviceCfg() ServiceConfig {
	return ServiceConfig{
		NodeRateLimiter: ratelimit.RateLimiterConfig{
			GlobalRPS:      100,
			GlobalBurst:    100,
			PerSenderRPS:   10,
			PerSenderBurst: 10,
		},
		UserRateLimiter: ratelimit.RateLimiterConfig{
			GlobalRPS:      50,
			GlobalBurst:    50,
			PerSenderRPS:   5,
			PerSenderBurst: 5,
		},
		CleanUpPeriodMs: defaultCleanUpPeriodMs,
	}
}

func createTestHandler(t *testing.T) *gatewayHandler {
	cfg := serviceCfg()
	return createTestHandlerWithConfig(t, cfg)
}

func createTestHandlerWithConfig(t *testing.T, cfg ServiceConfig) *gatewayHandler {
	configBytes, err := json.Marshal(cfg)
	require.NoError(t, err)

	donConfig := &config.DONConfig{
		DonId: "test-don",
	}
	mockDon := handlermocks.NewDON(t)
	mockHTTPClient := httpmocks.NewHTTPClient(t)
	lggr := logger.Test(t)

	handler, err := NewGatewayHandler(configBytes, donConfig, mockDon, mockHTTPClient, lggr)
	require.NoError(t, err)
	require.NotNil(t, handler)

	return handler
}
