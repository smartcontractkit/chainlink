package v2

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/uuid"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/ratelimit"
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

		handler, err := NewGatewayHandler(configBytes, donConfig, mockDon, mockHTTPClient, lggr)
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
				ReadFromCache: true,
				MaxAgeMs:      600000, // 10 minute TTL
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
				ReadFromCache: true,
				MaxAgeMs:      600000,
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
}

func newMockResponseCache() *mockResponseCache {
	return &mockResponseCache{
		deleteExpiredCh: make(chan struct{}),
	}
}

func (m *mockResponseCache) Set(workflowID string, req gateway_common.OutboundHTTPRequest, response gateway_common.OutboundHTTPResponse) {
}

func (m *mockResponseCache) CachedFetch(ctx context.Context, workflowID string, req gateway_common.OutboundHTTPRequest, fetchFn func() gateway_common.OutboundHTTPResponse) gateway_common.OutboundHTTPResponse {
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

	handler, err := NewGatewayHandler(configBytes, donConfig, mockDon, mockHTTPClient, lggr)
	require.NoError(t, err)
	require.NotNil(t, handler)

	return handler
}
