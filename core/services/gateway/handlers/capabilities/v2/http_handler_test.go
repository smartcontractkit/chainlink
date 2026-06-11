package v2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
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

		handler, err := NewGatewayHandler(configBytes, donConfig, mockDon, mockHTTPClient, lggr, limits.Factory{Logger: lggr}, defaultTestHTTPClientFactory)
		require.NoError(t, err)
		require.NotNil(t, handler)
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

		handler, err := NewGatewayHandler(invalidConfig, donConfig, mockDon, mockHTTPClient, lggr, limits.Factory{Logger: lggr}, defaultTestHTTPClientFactory)
		require.Error(t, err)
		require.Nil(t, handler)
	})

	t.Run("applies default config values", func(t *testing.T) {
		cfg := ServiceConfig{
			// CleanUpPeriodMs not set - should get default
		}
		configBytes, err := json.Marshal(cfg)
		require.NoError(t, err)

		donConfig := &config.DONConfig{DonId: "test-don"}
		mockDon := handlermocks.NewDON(t)
		mockHTTPClient := httpmocks.NewHTTPClient(t)
		lggr := logger.Test(t)

		handler, err := NewGatewayHandler(configBytes, donConfig, mockDon, mockHTTPClient, lggr, limits.Factory{Logger: lggr}, defaultTestHTTPClientFactory)
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

	t.Run("successful node message handling with MultiHeaders", func(t *testing.T) {
		mockDon := handler.don.(*handlermocks.DON)
		mockHTTPClient := handler.httpClient.(*httpmocks.HTTPClient)

		// Prepare outbound request
		outboundReq := gateway_common.OutboundHTTPRequest{
			Method:        "GET",
			URL:           "https://example.com/api/multiheaders-test",
			TimeoutMs:     5000,
			MultiHeaders:  map[string][]string{"Content-Type": {"application/json"}},
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

		// Response with multiple Set-Cookie headers
		httpResp := &network.HTTPResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Set-Cookie": "sessionid=abc123; Path=/; HttpOnly",
			},
			MultiHeaders: map[string][]string{
				"Set-Cookie": {
					"sessionid=abc123; Path=/; HttpOnly",
					"csrf_token=xyz789; Path=/; Secure",
				},
			},
			Body: []byte(`{"result": "success"}`),
		}

		mockHTTPClient.EXPECT().Send(mock.Anything, mock.MatchedBy(func(req network.HTTPRequest) bool {
			return req.Method == "GET" && req.URL == "https://example.com/api/multiheaders-test"
		})).Return(httpResp, nil).Once()

		capturedResponse := &gateway_common.OutboundHTTPResponse{}
		mockDon.EXPECT().SendToNode(mock.Anything, "node1", mock.MatchedBy(func(req *jsonrpc.Request[json.RawMessage]) bool {
			if req.Params == nil {
				return false
			}
			paramsStr := string(*req.Params)
			if !json.Valid(*req.Params) {
				return false
			}
			err2 := json.Unmarshal(*req.Params, capturedResponse)
			if err2 != nil {
				t.Logf("Failed to unmarshal response: %v, params: %s", err2, paramsStr)
				return false
			}
			if capturedResponse.StatusCode != 200 {
				return false
			}
			return req.ID == id
		})).Return(nil)

		err = handler.HandleNodeMessage(testutils.Context(t), resp, "node1")
		require.NoError(t, err)
		handler.wg.Wait()

		// Verify the response was captured
		require.Equal(t, 200, capturedResponse.StatusCode, "Response should have status code 200")
		require.NotNil(t, capturedResponse.MultiHeaders, "MultiHeaders should not be nil")
		require.NotEmpty(t, capturedResponse.MultiHeaders, "MultiHeaders should not be empty")
		setCookieValues, ok := capturedResponse.MultiHeaders["Set-Cookie"]
		require.True(t, ok, "Set-Cookie header should be in MultiHeaders, got: %+v", capturedResponse.MultiHeaders)
		require.Len(t, setCookieValues, 2, "Should have 2 Set-Cookie headers, got: %v", setCookieValues)
		require.Contains(t, setCookieValues, "sessionid=abc123; Path=/; HttpOnly")
		require.Contains(t, setCookieValues, "csrf_token=xyz789; Path=/; Secure")

		// Verify backward compatibility: all keys in MultiHeaders should be in Headers
		verifyBackwardCompatibility(t, capturedResponse.Headers, capturedResponse.MultiHeaders) //nolint:staticcheck // SA1019: intentionally asserting deprecated Headers for backward compatibility
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

func (m *mockResponseCache) Set(req gateway_common.OutboundHTTPRequest, response gateway_common.OutboundHTTPResponse) {
	m.setCallCount++
}

func (m *mockResponseCache) Fetch(ctx context.Context, req gateway_common.OutboundHTTPRequest, fetchFn func() gateway_common.OutboundHTTPResponse, storeOnFetch bool) gateway_common.OutboundHTTPResponse {
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

	handler, err := NewGatewayHandler(configBytes, donConfig, mockDon, mockHTTPClient, lggr, limits.Factory{Logger: lggr}, defaultTestHTTPClientFactory)
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
	return WithDefaults(ServiceConfig{})
}

func createTestHandler(t *testing.T) *gatewayHandler {
	cfg := serviceCfg()
	return createTestHandlerWithConfig(t, cfg)
}

// verifyBackwardCompatibility checks that all keys in MultiHeaders are also present in Headers
// with non-empty values. Same logic as in gateway/network/httpclient_test.go (package boundary).
func verifyBackwardCompatibility(t *testing.T, headers map[string]string, multiHeaders map[string][]string) {
	for key := range maps.Keys(multiHeaders) {
		require.NotEmpty(t, headers[key], "Headers should contain %s for backward compatibility", key)
	}
}

func createTestHandlerWithConfig(t *testing.T, cfg ServiceConfig) *gatewayHandler {
	configBytes, err := json.Marshal(cfg)
	require.NoError(t, err)

	donConfig := &config.DONConfig{
		DonId: "test-don",
		Members: []config.NodeConfig{
			{Name: "node1", Address: "node1"},
			{Name: "node2", Address: "node2"},
		},
	}
	mockDon := handlermocks.NewDON(t)
	mockHTTPClient := httpmocks.NewHTTPClient(t)
	lggr := logger.Test(t)

	handler, err := NewGatewayHandler(configBytes, donConfig, mockDon, mockHTTPClient, lggr, limits.Factory{Logger: lggr}, defaultTestHTTPClientFactory)
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

	t.Run("response with MultiHeaders is passed through correctly", func(t *testing.T) {
		handler := createTestHandler(t)
		mockHTTPClient := handler.httpClient.(*httpmocks.HTTPClient)

		expectedResp := &network.HTTPResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Set-Cookie": "sessionid=abc123; Path=/; HttpOnly, csrf_token=xyz789; Path=/; Secure",
				"Via":        "1.0 proxy1,1.1 proxy2",
			},
			MultiHeaders: map[string][]string{
				"Set-Cookie": {
					"sessionid=abc123; Path=/; HttpOnly",
					"csrf_token=xyz789; Path=/; Secure",
				},
				"Via": {
					"1.0 proxy1",
					"1.1 proxy2",
				},
			},
			Body: []byte(`{"result": "success"}`),
		}

		mockHTTPClient.EXPECT().Send(mock.Anything, mock.Anything).Return(expectedResp, nil)

		callback := handler.createHTTPRequestCallback(ctx, requestID, httpReq, outboundReq)
		response := callback()

		require.Equal(t, expectedResp.StatusCode, response.StatusCode)
		require.Equal(t, expectedResp.Body, response.Body)
		require.Empty(t, response.ErrorMessage)

		// Verify MultiHeaders are passed through
		require.NotNil(t, response.MultiHeaders, "MultiHeaders should not be nil")
		require.Len(t, response.MultiHeaders["Set-Cookie"], 2, "Should have 2 Set-Cookie headers")
		require.Contains(t, response.MultiHeaders["Set-Cookie"], "sessionid=abc123; Path=/; HttpOnly")
		require.Contains(t, response.MultiHeaders["Set-Cookie"], "csrf_token=xyz789; Path=/; Secure")
		require.Len(t, response.MultiHeaders["Via"], 2, "Should have 2 Via headers")
		require.Contains(t, response.MultiHeaders["Via"], "1.0 proxy1")
		require.Contains(t, response.MultiHeaders["Via"], "1.1 proxy2")

		// Verify Headers field is also set (for backward compatibility)
		require.NotNil(t, response.Headers, "Headers should not be nil")                         //nolint:staticcheck // SA1019: assert deprecated Headers for backward compatibility
		require.NotEmpty(t, response.Headers["Set-Cookie"], "Headers should contain Set-Cookie") //nolint:staticcheck // SA1019: assert deprecated Headers for backward compatibility
		require.NotEmpty(t, response.Headers["Via"], "Headers should contain Via")               //nolint:staticcheck // SA1019: assert deprecated Headers for backward compatibility

		// Verify backward compatibility: all keys in MultiHeaders should be in Headers
		verifyBackwardCompatibility(t, response.Headers, response.MultiHeaders) //nolint:staticcheck // SA1019: intentionally asserting deprecated Headers for backward compatibility
	})

	t.Run("response with empty MultiHeaders still sets Headers", func(t *testing.T) {
		handler := createTestHandler(t)
		mockHTTPClient := handler.httpClient.(*httpmocks.HTTPClient)

		expectedResp := &network.HTTPResponse{
			StatusCode: 200,
			Headers:    map[string]string{"Content-Type": "application/json"},
			MultiHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			Body: []byte(`{"result": "success"}`),
		}

		mockHTTPClient.EXPECT().Send(mock.Anything, mock.Anything).Return(expectedResp, nil)

		callback := handler.createHTTPRequestCallback(ctx, requestID, httpReq, outboundReq)
		response := callback()

		require.Equal(t, expectedResp.StatusCode, response.StatusCode)
		require.NotNil(t, response.MultiHeaders)
		require.Equal(t, []string{"application/json"}, response.MultiHeaders["Content-Type"])
		require.Equal(t, "application/json", response.Headers["Content-Type"]) //nolint:staticcheck // SA1019: assert deprecated Headers for backward compatibility

		// Verify backward compatibility: all keys in MultiHeaders should be in Headers
		verifyBackwardCompatibility(t, response.Headers, response.MultiHeaders) //nolint:staticcheck // SA1019: intentionally asserting deprecated Headers for backward compatibility
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

func TestMakeOutgoingRequest_SendResponseUsesIndependentContext(t *testing.T) {
	handler := createTestHandler(t)
	mockDon := handler.don.(*handlermocks.DON)
	mockHTTPClient := handler.httpClient.(*httpmocks.HTTPClient)

	outboundReq := gateway_common.OutboundHTTPRequest{
		Method:    "GET",
		URL:       "https://slow-endpoint.com/api",
		TimeoutMs: 50, // very short timeout so the HTTP context expires quickly
	}
	reqBytes, err := json.Marshal(outboundReq)
	require.NoError(t, err)

	id := fmt.Sprintf("%s/%s", gateway_common.MethodHTTPAction, uuid.New().String())
	rawRequest := json.RawMessage(reqBytes)
	resp := &jsonrpc.Response[json.RawMessage]{
		ID:     id,
		Result: &rawRequest,
	}

	// Simulate a slow HTTP endpoint that blocks until the context deadline expires.
	mockHTTPClient.EXPECT().Send(mock.Anything, mock.Anything).RunAndReturn(
		func(ctx context.Context, req network.HTTPRequest) (*network.HTTPResponse, error) {
			<-ctx.Done()
			return nil, ctx.Err()
		},
	)

	// The critical assertion: SendToNode must receive a non-expired context.
	// Before the fix, the same context was shared between the HTTP call and
	// sendResponseToNode, so an expired HTTP timeout would also prevent
	// delivering the response back to the node.
	mockDon.EXPECT().SendToNode(mock.MatchedBy(func(ctx context.Context) bool {
		return ctx.Err() == nil
	}), "node1", mock.Anything).Return(nil)

	err = handler.HandleNodeMessage(testutils.Context(t), resp, "node1")
	require.NoError(t, err)
	handler.wg.Wait()
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
	t.Run("per-node rate limiting", func(t *testing.T) {
		handler, resp, mockHTTPClient, mockDon := setupRateLimitingTest(t, ServiceConfig{})
		handler.globalNodeRateLimiter = limits.GlobalRateLimiter(100, 100)    // high global rate
		handler.perNodeRateLimiters["node1"] = limits.GlobalRateLimiter(1, 1) // very low per-node rate to trigger limits
		handler.perNodeRateLimiters["node2"] = limits.GlobalRateLimiter(1, 1) // very low per-node rate to trigger limits

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
		handler, resp, mockHTTPClient, mockDon := setupRateLimitingTest(t, ServiceConfig{})
		handler.globalNodeRateLimiter = limits.GlobalRateLimiter(1, 1)            // very low global rate to trigger limits
		handler.perNodeRateLimiters["node1"] = limits.GlobalRateLimiter(100, 100) // high per-node rate for node1
		handler.perNodeRateLimiters["node2"] = limits.GlobalRateLimiter(100, 100) // high per-node rate for node2

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

// defaultTestHTTPClientFactory returns a factory that always errors. Tests that
// exercise the mTLS path inject their own factory via handler.httpClientFactory.
var defaultTestHTTPClientFactory network.HTTPClientFactory = func(config network.HTTPClientConfig) (network.HTTPClient, error) {
	return nil, errors.New("httpClientFactory not configured for this test")
}

func TestGatewayHandler_Send_NoMtls_UsesDefaultClient(t *testing.T) {
	handler := createTestHandler(t)
	mockHTTPClient := handler.httpClient.(*httpmocks.HTTPClient)

	httpReq := network.HTTPRequest{Method: "GET", URL: "https://example.com/api"}
	outboundReq := gateway_common.OutboundHTTPRequest{Method: "GET", URL: "https://example.com/api"}

	expectedResp := &network.HTTPResponse{StatusCode: 200, Body: []byte("ok")}
	mockHTTPClient.EXPECT().Send(mock.Anything, httpReq).Return(expectedResp, nil).Once()

	resp, err := handler.send(testutils.Context(t), httpReq, outboundReq)
	require.NoError(t, err)
	require.Equal(t, expectedResp, resp)
}

func TestGatewayHandler_Send_MtlsBlockedByRateLimit(t *testing.T) {
	handler := createTestHandler(t)
	handler.mtlsRequestRateLimiter = limits.GlobalRateLimiter(0, 0)
	handler.httpClientFactory = func(config network.HTTPClientConfig) (network.HTTPClient, error) {
		return httpmocks.NewHTTPClient(t), nil
	}

	httpReq := network.HTTPRequest{Method: "GET", URL: "https://example.com/api"}
	outboundReq := gateway_common.OutboundHTTPRequest{
		Method: "GET",
		URL:    "https://example.com/api",
		Mtls: &gateway_common.MtlsAuth{
			PrivateKey:  []byte("private-key"),
			Certificate: []byte("certificate"),
		},
	}

	resp, err := handler.send(testutils.Context(t), httpReq, outboundReq)
	require.Error(t, err)
	require.Nil(t, resp)
	require.ErrorIs(t, err, network.ErrBlockedRequest)
	require.Contains(t, err.Error(), "global mtls request rate limit exceeded")
}

// TestGatewayHandler_Send_MtlsPassesConcurrencyLimiterToFactory verifies the
// handler hands its shared mtls concurrency limiter to the client factory. The
// limiter is enforced inside the HTTP client (on the request's capped-timeout
// context); that enforcement is covered by the network package tests.
func TestGatewayHandler_Send_MtlsPassesConcurrencyLimiterToFactory(t *testing.T) {
	handler := createTestHandler(t)
	handler.mtlsRequestRateLimiter = limits.GlobalRateLimiter(100, 100)

	httpReq := network.HTTPRequest{Method: "GET", URL: "https://example.com/api"}
	outboundReq := gateway_common.OutboundHTTPRequest{
		Method: "GET",
		URL:    "https://example.com/api",
		Mtls: &gateway_common.MtlsAuth{
			PrivateKey:  []byte("private-key"),
			Certificate: []byte("certificate"),
		},
	}

	expectedResp := &network.HTTPResponse{StatusCode: 200, Body: []byte("ok")}
	mtlsClient := httpmocks.NewHTTPClient(t)
	mtlsClient.EXPECT().Send(mock.Anything, httpReq).Return(expectedResp, nil).Once()

	var capturedConfig network.HTTPClientConfig
	handler.httpClientFactory = func(config network.HTTPClientConfig) (network.HTTPClient, error) {
		capturedConfig = config
		return mtlsClient, nil
	}

	resp, err := handler.send(testutils.Context(t), httpReq, outboundReq)
	require.NoError(t, err)
	require.Equal(t, expectedResp, resp)
	require.NotNil(t, capturedConfig.ConcurrencyLimiter, "handler must pass its mtls concurrency limiter to the client factory")
	require.Equal(t, handler.mtlsConcurrencyLimiter, capturedConfig.ConcurrencyLimiter)
}

func TestGatewayHandler_Send_MtlsUsesFactory(t *testing.T) {
	handler := createTestHandler(t)
	handler.mtlsRequestRateLimiter = limits.GlobalRateLimiter(100, 100)

	httpReq := network.HTTPRequest{Method: "GET", URL: "https://example.com/api"}
	outboundReq := gateway_common.OutboundHTTPRequest{
		Method: "GET",
		URL:    "https://example.com/api",
		Mtls: &gateway_common.MtlsAuth{
			PrivateKey:  []byte("private-key-bytes"),
			Certificate: []byte("certificate-bytes"),
		},
	}

	expectedResp := &network.HTTPResponse{StatusCode: 200, Body: []byte("mtls-ok")}
	mtlsClient := httpmocks.NewHTTPClient(t)
	mtlsClient.EXPECT().Send(mock.Anything, httpReq).Return(expectedResp, nil).Once()

	var capturedConfig network.HTTPClientConfig
	var factoryCalls int
	handler.httpClientFactory = func(config network.HTTPClientConfig) (network.HTTPClient, error) {
		factoryCalls++
		capturedConfig = config
		return mtlsClient, nil
	}

	resp, err := handler.send(testutils.Context(t), httpReq, outboundReq)
	require.NoError(t, err)
	require.Equal(t, expectedResp, resp)
	require.Equal(t, 1, factoryCalls, "factory should be called exactly once per mtls request")
	require.NotNil(t, capturedConfig.Mtls)
	require.Equal(t, []byte("private-key-bytes"), []byte(capturedConfig.Mtls.PrivateKey))
	require.Equal(t, []byte("certificate-bytes"), capturedConfig.Mtls.Certificate)

	// Default httpClient must not be used for mtls requests.
	mockDefault := handler.httpClient.(*httpmocks.HTTPClient)
	mockDefault.AssertNotCalled(t, "Send", mock.Anything, mock.Anything)
}

func TestGatewayHandler_Send_MtlsFactoryError(t *testing.T) {
	handler := createTestHandler(t)
	handler.mtlsRequestRateLimiter = limits.GlobalRateLimiter(100, 100)

	httpReq := network.HTTPRequest{Method: "GET", URL: "https://example.com/api"}
	outboundReq := gateway_common.OutboundHTTPRequest{
		Method: "GET",
		URL:    "https://example.com/api",
		Mtls:   &gateway_common.MtlsAuth{PrivateKey: []byte("k"), Certificate: []byte("c")},
	}

	factoryErr := errors.New("bad cert material")
	handler.httpClientFactory = func(config network.HTTPClientConfig) (network.HTTPClient, error) {
		return nil, factoryErr
	}

	resp, err := handler.send(testutils.Context(t), httpReq, outboundReq)
	require.Error(t, err)
	require.Nil(t, resp)
	require.ErrorIs(t, err, factoryErr, "factory error should be wrapped and discoverable via errors.Is")
}

// TestGatewayHandler_Send_InvalidMtlsCertDoesNotConsumeGlobalTokens verifies that a
// request carrying invalid mTLS credentials does not consume a global rate-limit token.
// Otherwise a malicious user could cheaply drain the shared mtls token bucket by spamming
// requests with bogus certificates. It uses the real HTTP client factory so that the
// production code path is what rejects the certificate as invalid.
func TestGatewayHandler_Send_InvalidMtlsCertDoesNotConsumeGlobalTokens(t *testing.T) {
	handler := createTestHandler(t)
	// Burst of exactly 1: only a single mtls request may pass the rate limiter.
	handler.mtlsRequestRateLimiter = limits.GlobalRateLimiter(1, 1)
	handler.httpClientFactory = network.NewHTTPClientFactory(network.HTTPClientConfig{}, logger.Test(t))

	httpReq := network.HTTPRequest{Method: "GET", URL: "https://example.com/api"}
	outboundReq := gateway_common.OutboundHTTPRequest{
		Method: "GET",
		URL:    "https://example.com/api",
		Mtls:   &gateway_common.MtlsAuth{PrivateKey: []byte("not-a-key"), Certificate: []byte("not-a-cert")},
	}

	ctx := testutils.Context(t)
	resp, err := handler.send(ctx, httpReq, outboundReq)
	require.Error(t, err)
	require.Nil(t, resp)
	require.Contains(t, err.Error(), "failed to parse MtlsAuth into KeyPair",
		"the real client factory should reject the invalid certificate material")

	// The single available token must still be present: the failed request above must not
	// have consumed it.
	require.True(t, handler.mtlsRequestRateLimiter.Allow(ctx),
		"global mtls rate-limit token must not be consumed by a request with an invalid certificate")
}

// TestGatewayHandler_Send_MtlsRoutesThroughCallbackOnly_DefaultClientUntouched
// verifies that an mTLS request flowing through the full callback path does not
// touch the default (shared) http client even when the factory returns a working
// client. This is the core property that prevents auth'd connections from
// leaking between users.
func TestGatewayHandler_Send_MtlsRoutesThroughCallbackOnly_DefaultClientUntouched(t *testing.T) {
	handler := createTestHandler(t)
	handler.mtlsRequestRateLimiter = limits.GlobalRateLimiter(100, 100)

	httpReq := network.HTTPRequest{Method: "GET", URL: "https://example.com/api", Timeout: 5 * time.Second}
	outboundReq := gateway_common.OutboundHTTPRequest{
		Method: "GET",
		URL:    "https://example.com/api",
		Mtls:   &gateway_common.MtlsAuth{PrivateKey: []byte("k"), Certificate: []byte("c")},
	}

	mtlsClient := httpmocks.NewHTTPClient(t)
	mtlsClient.EXPECT().Send(mock.Anything, mock.Anything).
		Return(&network.HTTPResponse{StatusCode: 200, Body: []byte("body")}, nil).Once()
	handler.httpClientFactory = func(config network.HTTPClientConfig) (network.HTTPClient, error) {
		return mtlsClient, nil
	}

	callback := handler.createHTTPRequestCallback(testutils.Context(t), "req-id", httpReq, outboundReq)
	resp := callback()
	require.Empty(t, resp.ErrorMessage)
	require.Equal(t, 200, resp.StatusCode)
	require.Equal(t, []byte("body"), resp.Body)

	mockDefault := handler.httpClient.(*httpmocks.HTTPClient)
	mockDefault.AssertNotCalled(t, "Send", mock.Anything, mock.Anything)
}

func TestGatewayHandler_Send_MtlsBlockedRequestIsValidationError(t *testing.T) {
	handler := createTestHandler(t)
	handler.mtlsRequestRateLimiter = limits.GlobalRateLimiter(0, 0)
	handler.httpClientFactory = func(config network.HTTPClientConfig) (network.HTTPClient, error) {
		return httpmocks.NewHTTPClient(t), nil
	}

	httpReq := network.HTTPRequest{Method: "GET", URL: "https://example.com/api", Timeout: 5 * time.Second}
	outboundReq := gateway_common.OutboundHTTPRequest{
		Method: "GET",
		URL:    "https://example.com/api",
		Mtls:   &gateway_common.MtlsAuth{PrivateKey: []byte("k"), Certificate: []byte("c")},
	}

	callback := handler.createHTTPRequestCallback(testutils.Context(t), "req-id", httpReq, outboundReq)
	resp := callback()
	require.NotEmpty(t, resp.ErrorMessage)
	require.True(t, resp.IsValidationError, "blocked mtls should be reported as a validation error")
	require.False(t, resp.IsExternalEndpointError, "rate-limited requests are not external endpoint errors")
	require.Equal(t, 0, resp.StatusCode)
}

// TestGatewayHandler_Send_MtlsRateLimitEnabledByDefault verifies that a handler
// built from the default settings has the mtls request rate limiter wired up and
// active, so mtls requests are gated without any per-test override. The default
// rate (cresettings.Default.GatewayHTTPActionMtlsRequestRate) has a zero burst,
// meaning mtls is blocked out of the box.
func TestGatewayHandler_Send_MtlsRateLimitEnabledByDefault(t *testing.T) {
	handler := createTestHandler(t)
	handler.httpClientFactory = func(config network.HTTPClientConfig) (network.HTTPClient, error) {
		return httpmocks.NewHTTPClient(t), nil
	}

	httpReq := network.HTTPRequest{Method: "GET", URL: "https://example.com/api"}
	outboundReq := gateway_common.OutboundHTTPRequest{
		Method: "GET",
		URL:    "https://example.com/api",
		Mtls: &gateway_common.MtlsAuth{
			PrivateKey:  []byte("private-key"),
			Certificate: []byte("certificate"),
		},
	}

	resp, err := handler.send(testutils.Context(t), httpReq, outboundReq)
	require.Error(t, err)
	require.Nil(t, resp)
	require.ErrorIs(t, err, network.ErrBlockedRequest)
	require.Contains(t, err.Error(), "global mtls request rate limit exceeded")
}
