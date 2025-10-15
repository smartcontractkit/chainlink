package v2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/ratelimit"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	gateway_common "github.com/smartcontractkit/chainlink-common/pkg/types/gateway"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	triggermocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities/v2/mocks"
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

		handler, err := NewGatewayHandler(configBytes, donConfig, mockDon, mockHTTPClient, lggr, limits.Factory{Logger: lggr})
		require.NoError(t, err)
		require.NotNil(t, handler)
		require.Equal(t, "test-don", handler.donConfig.DonId)
		require.NotNil(t, handler.responseCache)
		require.NotNil(t, handler.triggerHandler)
		require.NotNil(t, handler.metadataHandler)
	})

	t.Run("invalid config JSON", func(t *testing.T) {
		invalidConfig := []byte(`{invalid json}`)
		donConfig := &config.DONConfig{DonId: "test-don"}
		mockDon := handlermocks.NewDON(t)
		mockHTTPClient := httpmocks.NewHTTPClient(t)
		lggr := logger.Test(t)

		handler, err := NewGatewayHandler(invalidConfig, donConfig, mockDon, mockHTTPClient, lggr, limits.Factory{Logger: lggr})
		require.Error(t, err)
		require.Nil(t, handler)
	})

	t.Run("invalid rate limiter config", func(t *testing.T) {
		cfg := ServiceConfig{
			NodeRateLimiter: ratelimit.RateLimiterConfig{
				GlobalRPS:   -1, // Invalid negative rate
				GlobalBurst: 100,
			},
		}
		configBytes, err := json.Marshal(cfg)
		require.NoError(t, err)

		donConfig := &config.DONConfig{DonId: "test-don"}
		mockDon := handlermocks.NewDON(t)
		mockHTTPClient := httpmocks.NewHTTPClient(t)
		lggr := logger.Test(t)

		handler, err := NewGatewayHandler(configBytes, donConfig, mockDon, mockHTTPClient, lggr, limits.Factory{Logger: lggr})
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
			// CleanUpPeriodMs not set - should get default
		}
		configBytes, err := json.Marshal(cfg)
		require.NoError(t, err)

		donConfig := &config.DONConfig{DonId: "test-don"}
		mockDon := handlermocks.NewDON(t)
		mockHTTPClient := httpmocks.NewHTTPClient(t)
		lggr := logger.Test(t)

		handler, err := NewGatewayHandler(configBytes, donConfig, mockDon, mockHTTPClient, lggr, limits.Factory{Logger: lggr})
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
		outboundReq := gateway_common.OutboundHTTPRequest{
			Method:        "GET",
			URL:           "https://example.com/api",
			TimeoutMs:     5000,
			Headers:       map[string]string{"Content-Type": "application/json"},
			Body:          []byte(`{"test": "data"}`),
			CacheSettings: gateway_common.CacheSettings{},
		}
		reqBytes, err := json.Marshal(outboundReq)
		require.NoError(t, err)

		id := fmt.Sprintf("%s/%s", gateway_common.MethodHTTPAction, uuid.New().String())
		rawRequest := json.RawMessage(reqBytes)
		resp := &jsonrpc.Response[json.RawMessage]{
			ID:     id,
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
			return req.ID == id
		})).Return(nil)

		err = handler.HandleNodeMessage(testutils.Context(t), resp, "node1")
		require.NoError(t, err)
		handler.wg.Wait()
	})

	t.Run("returns cached response if available", func(t *testing.T) {
		outboundReq := gateway_common.OutboundHTTPRequest{
			Method:    "GET",
			URL:       "https://return-cached.com/api",
			TimeoutMs: 5000,
			CacheSettings: gateway_common.CacheSettings{
				Store:    true,
				MaxAgeMs: 600000, // Read from cache if cache entry is fresher than 10 minutes
			},
		}
		reqBytes, err := json.Marshal(outboundReq)
		require.NoError(t, err)
		id := fmt.Sprintf("%s/%s", gateway_common.MethodHTTPAction, uuid.New().String())
		rawRequest := json.RawMessage(reqBytes)
		resp := &jsonrpc.Response[json.RawMessage]{
			ID:     id,
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
			var cached gateway_common.OutboundHTTPResponse
			err2 := json.Unmarshal(*req.Params, &cached)
			return err2 == nil && string(cached.Body) == string(httpResp.Body)
		})).Return(nil)

		err = handler.HandleNodeMessage(testutils.Context(t), resp, "node1")
		require.NoError(t, err)
		handler.wg.Wait()
	})

	t.Run("status code 500 is not cached", func(t *testing.T) {
		outboundReq := gateway_common.OutboundHTTPRequest{
			Method:    "GET",
			URL:       "https://status-500.com/api",
			TimeoutMs: 5000,
			CacheSettings: gateway_common.CacheSettings{
				Store:    true,
				MaxAgeMs: 600000, // Read from cache if cache entry is fresher than 10 minutes
			},
		}
		reqBytes, err := json.Marshal(outboundReq)
		require.NoError(t, err)

		rawRequest := json.RawMessage(reqBytes)
		resp := &jsonrpc.Response[json.RawMessage]{
			ID:     fmt.Sprintf("%s/%s", gateway_common.MethodHTTPAction, uuid.New().String()),
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
			ID:     fmt.Sprintf("%s/%s", gateway_common.MethodHTTPAction, uuid.New().String()),
			Result: &rawRes,
		}

		err := handler.HandleNodeMessage(testutils.Context(t), resp, "node1")
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to unmarshal HTTP request")
		handler.wg.Wait()
	})
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
func TestHandleNodeMessage_RoutesToTriggerHandler(t *testing.T) {
	// This test covers the case where the response ID does not contain a "/"
	// and should be routed to the triggerHandler.HandleNodeTriggerResponse.
	mockTriggerHandler := triggermocks.NewHTTPTriggerHandler(t)
	handler := createTestHandler(t)
	handler.triggerHandler = mockTriggerHandler

	rawRes := json.RawMessage([]byte(`{}`))
	resp := &jsonrpc.Response[json.RawMessage]{
		ID:     "triggerResponseID", // No "/" in ID
		Result: &rawRes,
	}
	nodeAddr := "node1"

	mockTriggerHandler.
		On("HandleNodeTriggerResponse", mock.Anything, resp, nodeAddr).
		Return(nil).
		Once()

	err := handler.HandleNodeMessage(testutils.Context(t), resp, nodeAddr)
	require.NoError(t, err)
	mockTriggerHandler.AssertExpectations(t)
}

func TestHandleNodeMessage_UnsupportedMethod(t *testing.T) {
	handler := createTestHandler(t)
	rawRes := json.RawMessage([]byte(`{}`))
	resp := &jsonrpc.Response[json.RawMessage]{
		ID:     "unsupportedMethod/123",
		Result: &rawRes,
	}
	nodeAddr := "node1"

	err := handler.HandleNodeMessage(testutils.Context(t), resp, nodeAddr)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported method unsupportedMethod")
}

func TestHandleNodeMessage_EmptyID(t *testing.T) {
	handler := createTestHandler(t)
	rawRes := json.RawMessage([]byte(`{}`))
	resp := &jsonrpc.Response[json.RawMessage]{
		ID:     "",
		Result: &rawRes,
	}
	nodeAddr := "node1"

	err := handler.HandleNodeMessage(testutils.Context(t), resp, nodeAddr)
	require.Error(t, err)
	require.Contains(t, err.Error(), "empty request ID")
}

type mockResponseCache struct {
	deleteExpiredCh chan struct{}
	setCallCount    int
	fetchCallCount  int
}

func newMockResponseCache() *mockResponseCache {
	return &mockResponseCache{
		deleteExpiredCh: make(chan struct{}),
		setCallCount:    0,
		fetchCallCount:  0,
	}
}

func (m *mockResponseCache) Set(workflowID string, req gateway_common.OutboundHTTPRequest, response gateway_common.OutboundHTTPResponse) {
	m.setCallCount++
}

func (m *mockResponseCache) Fetch(ctx context.Context, workflowID string, req gateway_common.OutboundHTTPRequest, fetchFn func() gateway_common.OutboundHTTPResponse, storeOnFetch bool) gateway_common.OutboundHTTPResponse {
	m.fetchCallCount++
	return fetchFn()
}

func (m *mockResponseCache) DeleteExpired(ctx context.Context) int {
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

	handler, err := NewGatewayHandler(configBytes, donConfig, mockDon, mockHTTPClient, lggr, limits.Factory{Logger: lggr})
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
	cfg := ServiceConfig{
		NodeRateLimiter: ratelimit.RateLimiterConfig{
			GlobalRPS:      100,
			GlobalBurst:    100,
			PerSenderRPS:   10,
			PerSenderBurst: 10,
		},
	}
	return WithDefaults(cfg)
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

	handler, err := NewGatewayHandler(configBytes, donConfig, mockDon, mockHTTPClient, lggr, limits.Factory{Logger: lggr})
	require.NoError(t, err)
	require.NotNil(t, handler)

	return handler
}

func TestCreateHTTPRequestCallback(t *testing.T) {
	ctx := testutils.Context(t)

	requestID := "test-request-id"
	httpReq := network.HTTPRequest{
		Method:  "POST",
		URL:     "https://example.com/api",
		Headers: map[string]string{"Content-Type": "application/json"},
		Body:    []byte(`{"test": "data"}`),
		Timeout: 5 * time.Second,
	}
	outboundReq := gateway_common.OutboundHTTPRequest{
		Method:    "POST",
		URL:       "https://example.com/api",
		Headers:   map[string]string{"Content-Type": "application/json"},
		Body:      []byte(`{"test": "data"}`),
		TimeoutMs: 5000,
	}

	t.Run("successful HTTP request with latency measurement", func(t *testing.T) {
		handler := createTestHandler(t)
		mockHTTPClient := handler.httpClient.(*httpmocks.HTTPClient)

		expectedResp := &network.HTTPResponse{
			StatusCode: 200,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       []byte(`{"result": "success"}`),
		}

		mockHTTPClient.EXPECT().Send(mock.Anything, mock.Anything).Return(expectedResp, nil)

		callback := handler.createHTTPRequestCallback(ctx, requestID, httpReq, outboundReq)
		response := callback()

		require.Equal(t, expectedResp.StatusCode, response.StatusCode)
		require.Equal(t, expectedResp.Headers, response.Headers)
		require.Equal(t, expectedResp.Body, response.Body)
		require.Empty(t, response.ErrorMessage)
		require.False(t, response.IsExternalEndpointError)
		require.Positive(t, response.ExternalEndpointLatency)
	})

	t.Run("HTTP send error sets IsExternalEndpointError to true", func(t *testing.T) {
		handler := createTestHandler(t)
		mockHTTPClient := handler.httpClient.(*httpmocks.HTTPClient)

		mockHTTPClient.EXPECT().Send(mock.Anything, mock.Anything).Return(nil, network.ErrHTTPSend)

		callback := handler.createHTTPRequestCallback(ctx, requestID, httpReq, outboundReq)

		response := callback()

		require.NotEmpty(t, response.ErrorMessage, "Error message should not be empty")
		require.Equal(t, network.ErrHTTPSend.Error(), response.ErrorMessage)
		require.True(t, response.IsExternalEndpointError)
		require.Positive(t, response.ExternalEndpointLatency)
		require.Equal(t, 0, response.StatusCode)
		require.Nil(t, response.Headers)
		require.Nil(t, response.Body)
	})

	t.Run("HTTP read error sets IsExternalEndpointError to true", func(t *testing.T) {
		handler := createTestHandler(t)
		mockHTTPClient := handler.httpClient.(*httpmocks.HTTPClient)

		mockHTTPClient.EXPECT().Send(mock.Anything, mock.Anything).Return(nil, network.ErrHTTPRead)

		callback := handler.createHTTPRequestCallback(ctx, requestID, httpReq, outboundReq)

		response := callback()

		require.NotEmpty(t, response.ErrorMessage, "Error message should not be empty")
		require.Equal(t, network.ErrHTTPRead.Error(), response.ErrorMessage)
		require.True(t, response.IsExternalEndpointError)
		require.Positive(t, response.ExternalEndpointLatency)
		require.Equal(t, 0, response.StatusCode)
		require.Nil(t, response.Headers)
		require.Nil(t, response.Body)
	})

	t.Run("other errors set IsExternalEndpointError to false", func(t *testing.T) {
		handler := createTestHandler(t)
		mockHTTPClient := handler.httpClient.(*httpmocks.HTTPClient)

		genericError := errors.New("some other network error")
		mockHTTPClient.EXPECT().Send(mock.Anything, mock.Anything).Return(nil, genericError)

		callback := handler.createHTTPRequestCallback(ctx, requestID, httpReq, outboundReq)

		response := callback()

		require.NotEmpty(t, response.ErrorMessage, "Error message should not be empty")
		require.Equal(t, genericError.Error(), response.ErrorMessage)
		require.False(t, response.IsExternalEndpointError)
		require.Positive(t, response.ExternalEndpointLatency)
		require.Equal(t, 0, response.StatusCode)
		require.Nil(t, response.Headers)
		require.Nil(t, response.Body)
	})
}

// TestMakeOutgoingRequestCachingBehavior tests the specific caching logic in makeOutgoingRequest
func TestMakeOutgoingRequestCachingBehavior(t *testing.T) {
	t.Run("MaxAgeMs=0 and Store=true calls Set", func(t *testing.T) {
		handler := createTestHandler(t)
		mockCache := newMockResponseCache()
		handler.responseCache = mockCache

		outboundReq := gateway_common.OutboundHTTPRequest{
			Method:    "GET",
			URL:       "https://test-store-true.com/api",
			TimeoutMs: 5000,
			CacheSettings: gateway_common.CacheSettings{
				MaxAgeMs: 0,    // No cache read
				Store:    true, // But do store
			},
		}
		reqBytes, err := json.Marshal(outboundReq)
		require.NoError(t, err)
		id := fmt.Sprintf("%s/%s", gateway_common.MethodHTTPAction, uuid.New().String())
		rawRequest := json.RawMessage(reqBytes)
		resp := &jsonrpc.Response[json.RawMessage]{
			ID:     id,
			Result: &rawRequest,
		}

		mockDon := handler.don.(*handlermocks.DON)
		mockHTTPClient := handler.httpClient.(*httpmocks.HTTPClient)
		httpResp := &network.HTTPResponse{
			StatusCode: 200,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       []byte(`{"test": "data"}`),
		}
		mockHTTPClient.EXPECT().Send(mock.Anything, mock.Anything).Return(httpResp, nil).Once()
		mockDon.EXPECT().SendToNode(mock.Anything, "node1", mock.Anything).Return(nil)

		err = handler.HandleNodeMessage(testutils.Context(t), resp, "node1")
		require.NoError(t, err)
		handler.wg.Wait()

		// Verify Set was called once and CachedFetch was not called
		require.Equal(t, 1, mockCache.setCallCount, "Set should be called once when Store=true and MaxAgeMs=0")
		require.Equal(t, 0, mockCache.fetchCallCount, "CachedFetch should not be called when MaxAgeMs=0")
		mockHTTPClient.AssertExpectations(t)
		mockDon.AssertExpectations(t)
	})

	t.Run("MaxAgeMs=0 and Store=false does not call Set", func(t *testing.T) {
		handler := createTestHandler(t)
		mockCache := newMockResponseCache()
		handler.responseCache = mockCache

		outboundReq := gateway_common.OutboundHTTPRequest{
			Method:    "GET",
			URL:       "https://test-store-false.com/api",
			TimeoutMs: 5000,
			CacheSettings: gateway_common.CacheSettings{
				MaxAgeMs: 0,     // No cache read
				Store:    false, // Don't store
			},
		}
		reqBytes, err := json.Marshal(outboundReq)
		require.NoError(t, err)
		id := fmt.Sprintf("%s/%s", gateway_common.MethodHTTPAction, uuid.New().String())
		rawRequest := json.RawMessage(reqBytes)
		resp := &jsonrpc.Response[json.RawMessage]{
			ID:     id,
			Result: &rawRequest,
		}

		mockDon := handler.don.(*handlermocks.DON)
		mockHTTPClient := handler.httpClient.(*httpmocks.HTTPClient)
		httpResp := &network.HTTPResponse{
			StatusCode: 200,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       []byte(`{"test": "data"}`),
		}
		mockHTTPClient.EXPECT().Send(mock.Anything, mock.Anything).Return(httpResp, nil).Once()
		mockDon.EXPECT().SendToNode(mock.Anything, "node1", mock.Anything).Return(nil)

		err = handler.HandleNodeMessage(testutils.Context(t), resp, "node1")
		require.NoError(t, err)
		handler.wg.Wait()

		// Verify Set was not called and CachedFetch was not called
		require.Equal(t, 0, mockCache.setCallCount, "Set should not be called when Store=false and MaxAgeMs=0")
		require.Equal(t, 0, mockCache.fetchCallCount, "CachedFetch should not be called when MaxAgeMs=0")
		mockHTTPClient.AssertExpectations(t)
		mockDon.AssertExpectations(t)
	})

	t.Run("MaxAgeMs>0 calls CachedFetch", func(t *testing.T) {
		handler := createTestHandler(t)
		mockCache := newMockResponseCache()
		handler.responseCache = mockCache

		outboundReq := gateway_common.OutboundHTTPRequest{
			Method:    "GET",
			URL:       "https://test-cached-fetch.com/api",
			TimeoutMs: 5000,
			CacheSettings: gateway_common.CacheSettings{
				MaxAgeMs: 5000, // Cache read enabled
				Store:    true, // Store the response
			},
		}
		reqBytes, err := json.Marshal(outboundReq)
		require.NoError(t, err)
		id := fmt.Sprintf("%s/%s", gateway_common.MethodHTTPAction, uuid.New().String())
		rawRequest := json.RawMessage(reqBytes)
		resp := &jsonrpc.Response[json.RawMessage]{
			ID:     id,
			Result: &rawRequest,
		}

		mockDon := handler.don.(*handlermocks.DON)
		mockHTTPClient := handler.httpClient.(*httpmocks.HTTPClient)
		httpResp := &network.HTTPResponse{
			StatusCode: 200,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       []byte(`{"test": "cached"}`),
		}
		mockHTTPClient.EXPECT().Send(mock.Anything, mock.Anything).Return(httpResp, nil).Once()
		mockDon.EXPECT().SendToNode(mock.Anything, "node1", mock.Anything).Return(nil)

		err = handler.HandleNodeMessage(testutils.Context(t), resp, "node1")
		require.NoError(t, err)
		handler.wg.Wait()

		// Verify CachedFetch was called and Set was not called directly
		require.Equal(t, 1, mockCache.fetchCallCount, "CachedFetch should be called when MaxAgeMs>0")
		require.Equal(t, 0, mockCache.setCallCount, "Set should not be called directly when using CachedFetch")
		mockHTTPClient.AssertExpectations(t)
		mockDon.AssertExpectations(t)
	})
}

// setupRateLimitingTest creates common test setup for rate limiting tests
func setupRateLimitingTest(t *testing.T, cfg ServiceConfig) (*gatewayHandler, *jsonrpc.Response[json.RawMessage], *httpmocks.HTTPClient, *handlermocks.DON) {
	handler := createTestHandlerWithConfig(t, cfg)

	outboundReq := gateway_common.OutboundHTTPRequest{
		Method:    "GET",
		URL:       "https://example.com/api",
		TimeoutMs: 5000,
	}
	reqBytes, err := json.Marshal(outboundReq)
	require.NoError(t, err)

	id := gateway_common.MethodHTTPAction + "/workflowId123/uuid456"
	rawRequest := json.RawMessage(reqBytes)
	resp := &jsonrpc.Response[json.RawMessage]{
		ID:     id,
		Result: &rawRequest,
	}

	mockHTTPClient := handler.httpClient.(*httpmocks.HTTPClient)
	mockDon := handler.don.(*handlermocks.DON)

	return handler, resp, mockHTTPClient, mockDon
}

// expectSuccessfulRequest sets up expectations for a successful HTTP request
func expectSuccessfulRequest(mockHTTPClient *httpmocks.HTTPClient, mockDon *handlermocks.DON, nodeAddr string) {
	httpResp := &network.HTTPResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       []byte(`{"result": "success"}`),
	}
	mockHTTPClient.EXPECT().Send(mock.Anything, mock.Anything).Return(httpResp, nil).Once()
	mockDon.EXPECT().SendToNode(mock.Anything, nodeAddr, mock.Anything).Return(nil).Once()
}

func TestGatewayHandler_MakeOutgoingRequest_NodeRateLimiting(t *testing.T) {
	t.Run("node rate limiting with AllowVerbose", func(t *testing.T) {
		cfg := ServiceConfig{
			NodeRateLimiter: ratelimit.RateLimiterConfig{
				GlobalRPS:      100,
				GlobalBurst:    100,
				PerSenderRPS:   1, // Very low per-sender rate to trigger limits
				PerSenderBurst: 1,
			},
		}
		handler, resp, mockHTTPClient, mockDon := setupRateLimitingTest(t, cfg)

		// First request should succeed
		expectSuccessfulRequest(mockHTTPClient, mockDon, "node1")

		err := handler.HandleNodeMessage(testutils.Context(t), resp, "node1")
		require.NoError(t, err)
		handler.wg.Wait()

		// Second request from same node should be rate limited
		err = handler.HandleNodeMessage(testutils.Context(t), resp, "node1")
		require.Error(t, err)
		require.Contains(t, err.Error(), "rate limit exceeded for node")
		handler.wg.Wait()
	})

	t.Run("global rate limiting", func(t *testing.T) {
		cfg := ServiceConfig{
			NodeRateLimiter: ratelimit.RateLimiterConfig{
				GlobalRPS:      1, // Very low global rate to trigger limits
				GlobalBurst:    1,
				PerSenderRPS:   100,
				PerSenderBurst: 100,
			},
		}
		handler, resp, mockHTTPClient, mockDon := setupRateLimitingTest(t, cfg)

		// First request should succeed
		expectSuccessfulRequest(mockHTTPClient, mockDon, "node1")

		err := handler.HandleNodeMessage(testutils.Context(t), resp, "node1")
		require.NoError(t, err)
		handler.wg.Wait()

		// Second request from different node should be globally rate limited
		err = handler.HandleNodeMessage(testutils.Context(t), resp, "node2")
		require.Error(t, err)
		require.Contains(t, err.Error(), "global rate limit exceeded")
		handler.wg.Wait()
	})
}
