package v2

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	gateway_common "github.com/smartcontractkit/chainlink-common/pkg/types/gateway"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities/v2/metrics"
	hc "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
	handlermocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const workflowID = "0x1234567890abcdef1234567890abcdef12345678901234567890abcdef123456"
const workflowOwner = "0x1234567890abcdef1234567890abcdef12345678"
const requestID = "test-request-id"

func createTestMetrics(t *testing.T) *metrics.Metrics {
	m, err := metrics.NewMetrics()
	require.NoError(t, err)
	return m
}

func requireUserErrorSent(t *testing.T, payload handlers.UserCallbackPayload, errorCode int64) {
	require.NotEmpty(t, payload.RawResponse)
	require.Equal(t, api.FromJSONRPCErrorCode(errorCode), payload.ErrorCode)
}

func TestHttpTriggerHandler_HandleUserTriggerRequest(t *testing.T) {
	triggerReq := createTestTriggerRequest(workflowID)
	reqBytes, err := json.Marshal(triggerReq)
	require.NoError(t, err)

	rawParams := json.RawMessage(reqBytes)
	req := &jsonrpc.Request[json.RawMessage]{
		Version: "2.0",
		ID:      requestID,
		Method:  gateway_common.MethodWorkflowExecute,
		Params:  &rawParams,
	}
	privateKey := createTestPrivateKey(t)
	req.Auth = createTestJWTToken(t, req, privateKey)

	t.Run("successful trigger request", func(t *testing.T) {
		handler, mockDon := createTestTriggerHandler(t)
		registerWorkflow(t, handler, triggerReq.Workflow.WorkflowID, privateKey)
		callback := hc.NewCallback()

		// Mock DON to expect sends to all nodes
		mockDon.EXPECT().SendToNode(mock.Anything, "node1", mock.Anything).Return(nil)
		mockDon.EXPECT().SendToNode(mock.Anything, "node2", mock.Anything).Return(nil)
		mockDon.EXPECT().SendToNode(mock.Anything, "node3", mock.Anything).Return(nil)

		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callback, time.Now())
		require.NoError(t, err)

		handler.callbacksMu.Lock()
		saved, exists := handler.callbacks[requestID]
		handler.callbacksMu.Unlock()

		require.True(t, exists)
		require.Equal(t, callback, saved.Callback)
		require.NotNil(t, saved.responseAggregator)
	})

	t.Run("successful trigger request with missing 0x prefix", func(t *testing.T) {
		handler, mockDon := createTestTriggerHandler(t)
		registerWorkflow(t, handler, workflowID, privateKey)
		callback := hc.NewCallback()

		triggerReq := createTestTriggerRequest("1234567890abcdef1234567890abcdef12345678901234567890abcdef123456") // missing 0x prefix
		reqBytes, err := json.Marshal(triggerReq)
		require.NoError(t, err)

		rawParams := json.RawMessage(reqBytes)
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      requestID,
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &rawParams,
		}
		req.Auth = createTestJWTToken(t, req, privateKey)

		// Mock DON to expect sends to all nodes
		mockDon.EXPECT().SendToNode(mock.Anything, "node1", mock.Anything).Return(nil)
		mockDon.EXPECT().SendToNode(mock.Anything, "node2", mock.Anything).Return(nil)
		mockDon.EXPECT().SendToNode(mock.Anything, "node3", mock.Anything).Return(nil)

		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callback, time.Now())
		require.NoError(t, err)

		handler.callbacksMu.Lock()
		saved, exists := handler.callbacks[requestID]
		handler.callbacksMu.Unlock()

		require.True(t, exists)
		require.Equal(t, callback, saved.Callback)
		require.NotNil(t, saved.responseAggregator)
	})

	t.Run("successful trigger request with padded workflow ID", func(t *testing.T) {
		handler, mockDon := createTestTriggerHandler(t)
		workflowID := "0x00001234567890abcdef1234567890abcdef12345678901234567890abcdef12"
		registerWorkflow(t, handler, workflowID, privateKey)
		callback := hc.NewCallback()

		triggerReq := createTestTriggerRequest("0x1234567890abcdef1234567890abcdef12345678901234567890abcdef12") // missing 0s
		reqBytes, err := json.Marshal(triggerReq)
		require.NoError(t, err)

		rawParams := json.RawMessage(reqBytes)
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      requestID,
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &rawParams,
		}
		req.Auth = createTestJWTToken(t, req, privateKey)

		// Mock DON to expect sends to all nodes
		mockDon.EXPECT().SendToNode(mock.Anything, "node1", mock.Anything).Return(nil)
		mockDon.EXPECT().SendToNode(mock.Anything, "node2", mock.Anything).Return(nil)
		mockDon.EXPECT().SendToNode(mock.Anything, "node3", mock.Anything).Return(nil)

		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callback, time.Now())
		require.NoError(t, err)

		handler.callbacksMu.Lock()
		saved, exists := handler.callbacks[requestID]
		handler.callbacksMu.Unlock()

		require.True(t, exists)
		require.Equal(t, callback, saved.Callback)
		require.NotNil(t, saved.responseAggregator)
	})

	t.Run("successful trigger request with padded workflow ID and missing 0x prefix", func(t *testing.T) {
		handler, mockDon := createTestTriggerHandler(t)
		workflowID := "0x00001234567890abcdef1234567890abcdef12345678901234567890abcdef12"
		registerWorkflow(t, handler, workflowID, privateKey)
		callback := hc.NewCallback()

		triggerReq := createTestTriggerRequest("1234567890abcdef1234567890abcdef12345678901234567890abcdef12") // missing 0s
		reqBytes, err := json.Marshal(triggerReq)
		require.NoError(t, err)

		rawParams := json.RawMessage(reqBytes)
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      requestID,
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &rawParams,
		}
		req.Auth = createTestJWTToken(t, req, privateKey)

		// Mock DON to expect sends to all nodes
		mockDon.EXPECT().SendToNode(mock.Anything, "node1", mock.Anything).Return(nil)
		mockDon.EXPECT().SendToNode(mock.Anything, "node2", mock.Anything).Return(nil)
		mockDon.EXPECT().SendToNode(mock.Anything, "node3", mock.Anything).Return(nil)

		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callback, time.Now())
		require.NoError(t, err)

		handler.callbacksMu.Lock()
		saved, exists := handler.callbacks[requestID]
		handler.callbacksMu.Unlock()

		require.True(t, exists)
		require.Equal(t, callback, saved.Callback)
		require.NotNil(t, saved.responseAggregator)
	})

	t.Run("invalid JSON params", func(t *testing.T) {
		handler, _ := createTestTriggerHandler(t)
		callback := hc.NewCallback()

		rawParams := json.RawMessage(`{invalid json}`)
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      requestID,
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &rawParams,
		}

		err := handler.HandleUserTriggerRequest(testutils.Context(t), req, callback, time.Now())
		require.Error(t, err)
		r, err := callback.Wait(t.Context())
		require.NoError(t, err)
		requireUserErrorSent(t, r, jsonrpc.ErrParse)
	})

	t.Run("null JSON params", func(t *testing.T) {
		handler, _ := createTestTriggerHandler(t)
		callback := hc.NewCallback()

		rawParams := json.RawMessage(`null`)
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      requestID,
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &rawParams,
		}

		err := handler.HandleUserTriggerRequest(testutils.Context(t), req, callback, time.Now())
		require.Error(t, err)
		r, err := callback.Wait(t.Context())
		require.NoError(t, err)
		requireUserErrorSent(t, r, jsonrpc.ErrInvalidRequest)
	})

	t.Run("empty request ID", func(t *testing.T) {
		handler, _ := createTestTriggerHandler(t)
		callback := hc.NewCallback()

		triggerReq := createTestTriggerRequest(workflowID)
		reqBytes, err := json.Marshal(triggerReq)
		require.NoError(t, err)

		rawParams := json.RawMessage(reqBytes)
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      "", // Empty ID
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &rawParams,
		}

		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callback, time.Now())
		require.Error(t, err)
		require.Contains(t, err.Error(), "empty request ID")

		r, err := callback.Wait(t.Context())
		require.NoError(t, err)
		requireUserErrorSent(t, r, jsonrpc.ErrInvalidRequest)
	})

	t.Run("request ID contains slash", func(t *testing.T) {
		handler, _ := createTestTriggerHandler(t)
		callback := hc.NewCallback()

		triggerReq := gateway_common.HTTPTriggerRequest{
			Workflow: gateway_common.WorkflowSelector{
				WorkflowID: workflowID,
			},
			Input: []byte(`{"key": "value"}`),
		}
		reqBytes, err := json.Marshal(triggerReq)
		require.NoError(t, err)

		rawParams := json.RawMessage(reqBytes)
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      "test/request/id", // Contains slashes
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &rawParams,
		}

		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callback, time.Now())
		require.Error(t, err)
		require.Contains(t, err.Error(), "must not contain '/'")

		r, err := callback.Wait(t.Context())
		require.NoError(t, err)
		requireUserErrorSent(t, r, jsonrpc.ErrInvalidRequest)
	})

	t.Run("invalid method", func(t *testing.T) {
		handler, _ := createTestTriggerHandler(t)
		callback := hc.NewCallback()

		triggerReq := gateway_common.HTTPTriggerRequest{
			Workflow: gateway_common.WorkflowSelector{
				WorkflowID: workflowID,
			},
			Input: []byte(`{"key": "value"}`),
		}
		reqBytes, err := json.Marshal(triggerReq)
		require.NoError(t, err)

		rawParams := json.RawMessage(reqBytes)
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      requestID,
			Method:  "invalid-method",
			Params:  &rawParams,
		}

		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callback, time.Now())
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid method")

		r, err := callback.Wait(t.Context())
		require.NoError(t, err)
		requireUserErrorSent(t, r, jsonrpc.ErrMethodNotFound)
	})

	t.Run("duplicate request ID", func(t *testing.T) {
		handler, mockDon := createTestTriggerHandler(t)
		privateKey := createTestPrivateKey(t)
		registerWorkflow(t, handler, workflowID, privateKey)
		callback1 := hc.NewCallback()
		callback2 := hc.NewCallback()

		triggerReq := gateway_common.HTTPTriggerRequest{
			Workflow: gateway_common.WorkflowSelector{
				WorkflowID: workflowID,
			},
			Input: []byte(`{"key": "value"}`),
		}
		reqBytes, err := json.Marshal(triggerReq)
		require.NoError(t, err)

		rawParams := json.RawMessage(reqBytes)
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      requestID,
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &rawParams,
		}
		// First request should succeed
		req.Auth = createTestJWTToken(t, req, privateKey)
		mockDon.EXPECT().SendToNode(mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(3)
		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callback1, time.Now())
		require.NoError(t, err)

		// Second request with same ID should fail
		req.Auth = createTestJWTToken(t, req, privateKey)
		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callback2, time.Now())
		require.Error(t, err)
		require.Contains(t, err.Error(), "in-flight request")

		r, err := callback2.Wait(t.Context())
		require.NoError(t, err)
		requireUserErrorSent(t, r, jsonrpc.ErrConflict)
	})

	t.Run("duplicate JWT token and request ID", func(t *testing.T) {
		handler, mockDon := createTestTriggerHandler(t)
		privateKey := createTestPrivateKey(t)
		registerWorkflow(t, handler, workflowID, privateKey)
		callback1 := hc.NewCallback()
		callback2 := hc.NewCallback()

		triggerReq := gateway_common.HTTPTriggerRequest{
			Workflow: gateway_common.WorkflowSelector{
				WorkflowID: workflowID,
			},
			Input: []byte(`{"key": "value"}`),
		}
		reqBytes, err := json.Marshal(triggerReq)
		require.NoError(t, err)

		rawParams := json.RawMessage(reqBytes)
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      requestID,
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &rawParams,
		}
		// First request should succeed
		req.Auth = createTestJWTToken(t, req, privateKey)
		mockDon.EXPECT().SendToNode(mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(3)
		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callback1, time.Now())
		require.NoError(t, err)

		// Second request with same ID should fail
		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callback2, time.Now())
		require.Error(t, err)
		require.Contains(t, err.Error(), "token has already been used")

		r, err := callback2.Wait(t.Context())
		require.NoError(t, err)
		requireUserErrorSent(t, r, jsonrpc.ErrInvalidRequest)
	})

	t.Run("invalid input JSON", func(t *testing.T) {
		handler, _ := createTestTriggerHandler(t)
		callback := hc.NewCallback()

		rawParams := json.RawMessage([]byte(`{"workflow":{"workflowID":"0x1234567890abcdef1234567890abcdef12345678901234567890abcdef123456"},"input":{"invalid json"}`))
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      requestID,
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &rawParams,
		}

		err := handler.HandleUserTriggerRequest(testutils.Context(t), req, callback, time.Now())
		require.Error(t, err)
	})
}

func TestHttpTriggerHandler_HandleNodeTriggerResponse(t *testing.T) {
	t.Run("successful aggregation", func(t *testing.T) {
		handler, mockDon := createTestTriggerHandler(t)
		privateKey := createTestPrivateKey(t)
		registerWorkflow(t, handler, workflowID, privateKey)
		callback := hc.NewCallback()

		// First, create a trigger request to set up the callback
		triggerReq := gateway_common.HTTPTriggerRequest{
			Workflow: gateway_common.WorkflowSelector{
				WorkflowID: workflowID,
			},
			Input: []byte(`{"key": "value"}`),
		}
		reqBytes, err := json.Marshal(triggerReq)
		require.NoError(t, err)

		rawParams := json.RawMessage(reqBytes)
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      requestID,
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &rawParams,
		}
		req.Auth = createTestJWTToken(t, req, privateKey)

		mockDon.EXPECT().SendToNode(mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(3)
		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callback, time.Now())
		require.NoError(t, err)
		// Create node responses
		rawRes := json.RawMessage(`{"result":"success"}`)
		nodeResp := &jsonrpc.Response[json.RawMessage]{
			Version: "2.0",
			ID:      requestID,
			Result:  &rawRes,
		}

		// Send responses from multiple nodes (need (N+F)//2+1 = (3+1)//2+1 = 3 for N=3, F=1)
		err = handler.HandleNodeTriggerResponse(testutils.Context(t), nodeResp, "node1")
		require.NoError(t, err)

		err = handler.HandleNodeTriggerResponse(testutils.Context(t), nodeResp, "node2")
		require.NoError(t, err)

		// Third response should trigger aggregation
		err = handler.HandleNodeTriggerResponse(testutils.Context(t), nodeResp, "node3")
		require.NoError(t, err)

		// Check that callback was called
		payload, err := callback.Wait(t.Context())
		require.NoError(t, err)
		require.NotEmpty(t, payload.RawResponse)
		require.Equal(t, api.NoError, payload.ErrorCode)

		var resp jsonrpc.Response[json.RawMessage]
		err = json.Unmarshal(payload.RawResponse, &resp)
		require.NoError(t, err)
		require.Equal(t, nodeResp.Result, resp.Result)
	})

	t.Run("callback not found", func(t *testing.T) {
		handler, _ := createTestTriggerHandler(t)

		rawRes := json.RawMessage(`{"result": "success"}`)
		nodeResp := &jsonrpc.Response[json.RawMessage]{
			Version: "2.0",
			ID:      "nonexistent-execution-id",
			Result:  &rawRes,
		}

		err := handler.HandleNodeTriggerResponse(testutils.Context(t), nodeResp, "node1")
		require.Error(t, err)
		require.Contains(t, err.Error(), "callback not found")
	})
}

func TestHttpTriggerHandler_ServiceLifecycle(t *testing.T) {
	t.Run("start and stop", func(t *testing.T) {
		handler, _ := createTestTriggerHandler(t)

		ctx := testutils.Context(t)
		err := handler.Start(ctx)
		require.NoError(t, err)

		err = handler.Close()
		require.NoError(t, err)
	})

	t.Run("double start and close should errors", func(t *testing.T) {
		handler, _ := createTestTriggerHandler(t)

		ctx := testutils.Context(t)
		err := handler.Start(ctx)
		require.NoError(t, err)

		err = handler.Start(ctx)
		require.Error(t, err)

		err = handler.Close()
		require.NoError(t, err)

		err = handler.Close()
		require.Error(t, err)
	})
}

func registerWorkflow(t *testing.T, handler *httpTriggerHandler, workflowID string, privateKey *ecdsa.PrivateKey) {
	handler.workflowMetadataHandler.authorizedKeys[workflowID] = map[gateway_common.AuthorizedKey]struct{}{
		{
			KeyType:   gateway_common.KeyTypeECDSAEVM,
			PublicKey: strings.ToLower(crypto.PubkeyToAddress(privateKey.PublicKey).Hex()),
		}: {},
	}
	handler.workflowMetadataHandler.workflowIDToRef[workflowID] = workflowReference{
		workflowOwner: workflowOwner,
		workflowName:  "test-workflow",
		workflowTag:   "v1.0",
	}
}

func TestHttpTriggerHandler_ReapExpiredCallbacks(t *testing.T) {
	requestID := "test-request-id"
	triggerReq := gateway_common.HTTPTriggerRequest{
		Workflow: gateway_common.WorkflowSelector{
			WorkflowID: workflowID,
		},
		Input: []byte(`{"key": "value"}`),
	}
	reqBytes, err := json.Marshal(triggerReq)
	require.NoError(t, err)

	rawParams := json.RawMessage(reqBytes)
	req := &jsonrpc.Request[json.RawMessage]{
		Version: "2.0",
		ID:      requestID,
		Method:  gateway_common.MethodWorkflowExecute,
		Params:  &rawParams,
	}
	privateKey := createTestPrivateKey(t)
	cfg := ServiceConfig{
		CleanUpPeriodMs:             100,
		MaxTriggerRequestDurationMs: 50,
	}
	handler, mockDon := createTestTriggerHandlerWithConfig(t, cfg)
	registerWorkflow(t, handler, workflowID, privateKey)

	t.Run("reap expired callbacks", func(t *testing.T) {
		req.Auth = createTestJWTToken(t, req, privateKey)
		callback := hc.NewCallback()
		mockDon.EXPECT().SendToNode(mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(3)
		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callback, time.Now())
		require.NoError(t, err)

		// Manually set the callback's createdAt to the past to simulate expiration
		handler.callbacksMu.Lock()
		if cb, exists := handler.callbacks[requestID]; exists {
			cb.createdAt = time.Now().Add(-time.Duration(cfg.MaxTriggerRequestDurationMs+1) * time.Millisecond)
			handler.callbacks[requestID] = cb
		}
		handler.callbacksMu.Unlock()

		// Manually trigger reaping
		handler.reapExpiredCallbacks(testutils.Context(t))

		// Verify callback was removed
		handler.callbacksMu.Lock()
		_, exists := handler.callbacks[requestID]
		handler.callbacksMu.Unlock()
		require.False(t, exists)
	})

	t.Run("keep non-expired callbacks", func(t *testing.T) {
		req.Auth = createTestJWTToken(t, req, privateKey)
		callback := hc.NewCallback()

		mockDon.EXPECT().SendToNode(mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(3)
		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callback, time.Now())
		require.NoError(t, err)

		// Optionally, set createdAt to now (should not be expired)
		handler.callbacksMu.Lock()
		if cb, exists := handler.callbacks[requestID]; exists {
			cb.createdAt = time.Now()
			handler.callbacks[requestID] = cb
		}
		handler.callbacksMu.Unlock()

		// Manually trigger reaping
		handler.reapExpiredCallbacks(testutils.Context(t))

		// Verify callback still exists
		handler.callbacksMu.Lock()
		_, exists := handler.callbacks[requestID]
		handler.callbacksMu.Unlock()
		require.True(t, exists)
	})
}

func TestIsValidJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected bool
	}{
		{
			name:     "valid JSON object",
			input:    []byte(`{"key": "value"}`),
			expected: true,
		},
		{
			name:     "valid JSON array",
			input:    []byte(`[1, 2, 3]`),
			expected: true,
		},
		{
			name:     "invalid JSON",
			input:    []byte(`{invalid}`),
			expected: false,
		},
		{
			name:     "empty object",
			input:    []byte(`{}`),
			expected: true,
		},
		{
			name:     "null",
			input:    []byte(`null`),
			expected: false,
		},
		{
			name:     "empty string",
			input:    []byte(``),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidJSON(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestHttpTriggerHandler_HandleUserTriggerRequest_Retries(t *testing.T) {
	lggr := logger.Test(t)
	cfg := ServiceConfig{
		MaxTriggerRequestDurationMs: 2000, // 2 seconds for test
		CleanUpPeriodMs:             10000,
	}

	donConfig := &config.DONConfig{
		DonId: "test-don",
		F:     1, // 1 faulty node, so (N+F)//2+1=(3+1)//2+1=3 for threshold
		Members: []config.NodeConfig{
			{Address: "node1"},
			{Address: "node2"},
			{Address: "node3"},
		},
	}

	mockDon := handlermocks.NewDON(t)
	metadataHandler := createTestMetadataHandler(t)
	userRateLimiter := createTestUserRateLimiter()
	testMetrics := createTestMetrics(t)
	handler := NewHTTPTriggerHandler(lggr, cfg, donConfig, mockDon, metadataHandler, userRateLimiter, testMetrics)
	privateKey := createTestPrivateKey(t)
	registerWorkflow(t, handler, workflowID, privateKey)

	t.Run("retries failed nodes until success", func(t *testing.T) {
		rawParams := json.RawMessage(`{"input":{},"workflow":{"workflowID":"0x1234567890abcdef1234567890abcdef12345678901234567890abcdef123456"}}`)
		req := &jsonrpc.Request[json.RawMessage]{
			ID:      "test-request-id",
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &rawParams,
			Version: "2.0",
		}
		req.Auth = createTestJWTToken(t, req, privateKey)

		callback := hc.NewCallback()

		// First attempt: node1 succeeds, node2 and node3 fail
		mockDon.On("SendToNode", mock.Anything, "node1", mock.Anything).Return(nil).Once()
		mockDon.On("SendToNode", mock.Anything, "node2", mock.Anything).Return(errors.New("connection error")).Once()
		mockDon.On("SendToNode", mock.Anything, "node3", mock.Anything).Return(errors.New("connection error")).Once()

		// Retry: node2 succeeds, node3 still fails
		mockDon.On("SendToNode", mock.Anything, "node2", mock.Anything).Return(nil).Once()
		mockDon.On("SendToNode", mock.Anything, "node3", mock.Anything).Return(errors.New("still failing")).Once()

		// Final retry: node3 succeeds
		mockDon.On("SendToNode", mock.Anything, "node3", mock.Anything).Return(nil).Once()

		err := handler.Start(testutils.Context(t))
		require.NoError(t, err)

		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callback, time.Now())
		require.NoError(t, err)

		mockDon.AssertExpectations(t)
		err = handler.Close()
		require.NoError(t, err)
	})
}

func TestHttpTriggerHandler_HandleUserTriggerRequest_JWTAuthorization(t *testing.T) {
	handler, mockDon := createTestTriggerHandler(t)
	ctx := testutils.Context(t)

	// Setup metadata handler with test data
	err := handler.workflowMetadataHandler.agg.Start(ctx)
	require.NoError(t, err)
	defer handler.workflowMetadataHandler.agg.Close()

	// Create test keys
	privateKey := createTestPrivateKey(t)
	signerAddr := crypto.PubkeyToAddress(privateKey.PublicKey)

	// Add authorized key to metadata handler
	key := gateway_common.AuthorizedKey{
		KeyType:   gateway_common.KeyTypeECDSAEVM,
		PublicKey: strings.ToLower(signerAddr.Hex()),
	}
	handler.workflowMetadataHandler.authorizedKeys[workflowID] = map[gateway_common.AuthorizedKey]struct{}{key: {}}
	handler.workflowMetadataHandler.workflowIDToRef[workflowID] = workflowReference{
		workflowOwner: workflowOwner,
		workflowName:  "test-workflow",
		workflowTag:   "v1.0",
	}

	t.Run("successful JWT authorization", func(t *testing.T) {
		callback := hc.NewCallback()

		triggerReq := createTestTriggerRequest(workflowID)
		reqBytes, err2 := json.Marshal(triggerReq)
		require.NoError(t, err2)

		rawParams := json.RawMessage(reqBytes)
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      "test-request-id",
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &rawParams,
		}

		jwtToken := createTestJWTToken(t, req, privateKey)
		req.Auth = jwtToken

		mockDon.EXPECT().SendToNode(mock.Anything, "node1", mock.MatchedBy(func(r *jsonrpc.Request[json.RawMessage]) bool {
			var params gateway_common.HTTPTriggerRequest
			err = json.Unmarshal(*r.Params, &params)
			return err == nil && params.Key.PublicKey == key.PublicKey
		})).Return(nil)
		mockDon.EXPECT().SendToNode(mock.Anything, "node2", mock.Anything).Return(nil)
		mockDon.EXPECT().SendToNode(mock.Anything, "node3", mock.Anything).Return(nil)

		err = handler.HandleUserTriggerRequest(ctx, req, callback, time.Now())
		require.NoError(t, err)
		handler.callbacksMu.Lock()
		_, exists := handler.callbacks[req.ID]
		handler.callbacksMu.Unlock()
		require.True(t, exists)
	})

	t.Run("invalid JWT token", func(t *testing.T) {
		callback := hc.NewCallback()

		triggerReq := createTestTriggerRequest(workflowID)
		reqBytes, err2 := json.Marshal(triggerReq)
		require.NoError(t, err2)

		rawParams := json.RawMessage(reqBytes)
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      "test-request-id-2",
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &rawParams,
			Auth:    "invalid.jwt.token",
		}

		err = handler.HandleUserTriggerRequest(ctx, req, callback, time.Now())
		require.Error(t, err)
		require.Contains(t, err.Error(), "auth failure")

		r, err2 := callback.Wait(t.Context())
		require.NoError(t, err2)
		requireUserErrorSent(t, r, jsonrpc.ErrInvalidRequest)
	})

	t.Run("unauthorized signer", func(t *testing.T) {
		callback := hc.NewCallback()
		unauthorizedKey := createTestPrivateKey(t)

		triggerReq := createTestTriggerRequest(workflowID)
		reqBytes, err2 := json.Marshal(triggerReq)
		require.NoError(t, err2)

		rawParams := json.RawMessage(reqBytes)
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      "test-request-id-3",
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &rawParams,
		}

		jwtToken := createTestJWTToken(t, req, unauthorizedKey)
		req.Auth = jwtToken

		err = handler.HandleUserTriggerRequest(ctx, req, callback, time.Now())
		require.Error(t, err)
		require.Contains(t, err.Error(), "auth failure")

		r, err2 := callback.Wait(t.Context())
		require.NoError(t, err2)
		requireUserErrorSent(t, r, jsonrpc.ErrInvalidRequest)
	})

	t.Run("workflow not found", func(t *testing.T) {
		callback := hc.NewCallback()

		triggerReq := gateway_common.HTTPTriggerRequest{
			Workflow: gateway_common.WorkflowSelector{
				WorkflowID: "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
			},
			Input: []byte(`{"key": "value"}`),
		}
		reqBytes, err2 := json.Marshal(triggerReq)
		require.NoError(t, err2)

		rawParams := json.RawMessage(reqBytes)
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      "test-request-id-4",
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &rawParams,
		}

		jwtToken := createTestJWTToken(t, req, privateKey)
		req.Auth = jwtToken

		err = handler.HandleUserTriggerRequest(ctx, req, callback, time.Now())
		require.Error(t, err)
		require.Contains(t, err.Error(), "auth failure")

		r, err2 := callback.Wait(t.Context())
		require.NoError(t, err2)
		requireUserErrorSent(t, r, jsonrpc.ErrInvalidRequest)
	})
}

func TestHttpTriggerHandler_HandleUserTriggerRequest_WorkflowLookup(t *testing.T) {
	handler, mockDon := createTestTriggerHandler(t)
	ctx := testutils.Context(t)

	err := handler.workflowMetadataHandler.agg.Start(ctx)
	require.NoError(t, err)
	defer handler.workflowMetadataHandler.agg.Close()

	privateKey := createTestPrivateKey(t)
	signerAddr := crypto.PubkeyToAddress(privateKey.PublicKey)

	workflowName := "0x" + hex.EncodeToString([]byte(workflows.HashTruncateName("test-workflow")))
	workflowOwner := "0x00001234567890abcdef1234567890abcdef1234"
	workflowTag := "v1.0"

	key := gateway_common.AuthorizedKey{
		KeyType:   gateway_common.KeyTypeECDSAEVM,
		PublicKey: strings.ToLower(signerAddr.Hex()),
	}
	handler.workflowMetadataHandler.authorizedKeys[workflowID] = map[gateway_common.AuthorizedKey]struct{}{key: {}}
	workflowRef := workflowReference{
		workflowOwner: workflowOwner,
		workflowName:  workflowName,
		workflowTag:   workflowTag,
	}
	handler.workflowMetadataHandler.workflowIDToRef[workflowID] = workflowRef
	handler.workflowMetadataHandler.workflowRefToID[workflowRef] = workflowID

	t.Run("successful workflow lookup by name", func(t *testing.T) {
		callback := hc.NewCallback()

		triggerReq := gateway_common.HTTPTriggerRequest{
			Workflow: gateway_common.WorkflowSelector{
				WorkflowOwner: workflowOwner,
				WorkflowName:  "test-workflow", // Use original name, not hashed
				WorkflowTag:   workflowTag,
			},
			Input: []byte(`{"key": "value"}`),
		}
		reqBytes, err := json.Marshal(triggerReq)
		require.NoError(t, err)

		rawParams := json.RawMessage(reqBytes)
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      "test-request-id",
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &rawParams,
		}

		// Create JWT token
		jwtToken := createTestJWTToken(t, req, privateKey)
		req.Auth = jwtToken

		// Mock DON to expect sends to all nodes
		mockDon.EXPECT().SendToNode(mock.Anything, "node1", mock.Anything).Return(nil)
		mockDon.EXPECT().SendToNode(mock.Anything, "node2", mock.Anything).Return(nil)
		mockDon.EXPECT().SendToNode(mock.Anything, "node3", mock.Anything).Return(nil)

		err = handler.HandleUserTriggerRequest(ctx, req, callback, time.Now())
		require.NoError(t, err)
	})

	t.Run("successful workflow lookup by name with missing 0x prefix", func(t *testing.T) {
		callback := hc.NewCallback()

		triggerReq := gateway_common.HTTPTriggerRequest{
			Workflow: gateway_common.WorkflowSelector{
				WorkflowOwner: "00001234567890abcdef1234567890abcdef1234", // missing 0x prefix
				WorkflowName:  "test-workflow",                            // Use original name, not hashed
				WorkflowTag:   workflowTag,
			},
			Input: []byte(`{"key": "value"}`),
		}
		reqBytes, err := json.Marshal(triggerReq)
		require.NoError(t, err)

		rawParams := json.RawMessage(reqBytes)
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      "test-request-id2",
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &rawParams,
		}

		// Create JWT token
		jwtToken := createTestJWTToken(t, req, privateKey)
		req.Auth = jwtToken

		// Mock DON to expect sends to all nodes
		mockDon.EXPECT().SendToNode(mock.Anything, "node1", mock.Anything).Return(nil)
		mockDon.EXPECT().SendToNode(mock.Anything, "node2", mock.Anything).Return(nil)
		mockDon.EXPECT().SendToNode(mock.Anything, "node3", mock.Anything).Return(nil)

		err = handler.HandleUserTriggerRequest(ctx, req, callback, time.Now())
		require.NoError(t, err)
	})

	t.Run("successful workflow lookup by name with padded workflow owner", func(t *testing.T) {
		callback := hc.NewCallback()

		triggerReq := gateway_common.HTTPTriggerRequest{
			Workflow: gateway_common.WorkflowSelector{
				WorkflowOwner: "0x1234567890abcdef1234567890abcdef1234", // missing 0s
				WorkflowName:  "test-workflow",                          // Use original name, not hashed
				WorkflowTag:   workflowTag,
			},
			Input: []byte(`{"key": "value"}`),
		}
		reqBytes, err := json.Marshal(triggerReq)
		require.NoError(t, err)

		rawParams := json.RawMessage(reqBytes)
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      "test-request-id4",
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &rawParams,
		}

		// Create JWT token
		jwtToken := createTestJWTToken(t, req, privateKey)
		req.Auth = jwtToken

		// Mock DON to expect sends to all nodes
		mockDon.EXPECT().SendToNode(mock.Anything, "node1", mock.Anything).Return(nil)
		mockDon.EXPECT().SendToNode(mock.Anything, "node2", mock.Anything).Return(nil)
		mockDon.EXPECT().SendToNode(mock.Anything, "node3", mock.Anything).Return(nil)

		err = handler.HandleUserTriggerRequest(ctx, req, callback, time.Now())
		require.NoError(t, err)
	})

	t.Run("successful workflow lookup by name with padded workflow owner and missing 0x prefix", func(t *testing.T) {
		callback := hc.NewCallback()

		triggerReq := gateway_common.HTTPTriggerRequest{
			Workflow: gateway_common.WorkflowSelector{
				WorkflowOwner: "1234567890abcdef1234567890abcdef1234", // missing 0x prefix
				WorkflowName:  "test-workflow",                        // Use original name, not hashed
				WorkflowTag:   workflowTag,
			},
			Input: []byte(`{"key": "value"}`),
		}
		reqBytes, err := json.Marshal(triggerReq)
		require.NoError(t, err)

		rawParams := json.RawMessage(reqBytes)
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      "test-request-id3",
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &rawParams,
		}

		// Create JWT token
		jwtToken := createTestJWTToken(t, req, privateKey)
		req.Auth = jwtToken

		// Mock DON to expect sends to all nodes
		mockDon.EXPECT().SendToNode(mock.Anything, "node1", mock.Anything).Return(nil)
		mockDon.EXPECT().SendToNode(mock.Anything, "node2", mock.Anything).Return(nil)
		mockDon.EXPECT().SendToNode(mock.Anything, "node3", mock.Anything).Return(nil)

		err = handler.HandleUserTriggerRequest(ctx, req, callback, time.Now())
		require.NoError(t, err)
	})

	t.Run("workflow not found by name", func(t *testing.T) {
		callback := hc.NewCallback()

		triggerReq := gateway_common.HTTPTriggerRequest{
			Workflow: gateway_common.WorkflowSelector{
				WorkflowOwner: workflowOwner,
				WorkflowName:  "nonexistent-workflow",
				WorkflowTag:   workflowTag,
			},
			Input: []byte(`{"key": "value"}`),
		}
		reqBytes, err := json.Marshal(triggerReq)
		require.NoError(t, err)

		rawParams := json.RawMessage(reqBytes)
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      "test-request-id-2",
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &rawParams,
		}

		jwtToken := createTestJWTToken(t, req, privateKey)
		req.Auth = jwtToken

		err = handler.HandleUserTriggerRequest(ctx, req, callback, time.Now())
		require.Error(t, err)
		require.Contains(t, err.Error(), "workflow not found")

		r, err := callback.Wait(t.Context())
		require.NoError(t, err)
		requireUserErrorSent(t, r, jsonrpc.ErrInvalidRequest)
	})
}
func TestHttpTriggerHandler_HandleUserTriggerRequest_Validation(t *testing.T) {
	handler, mockDon := createTestTriggerHandler(t)

	t.Run("workflowID uppercase", func(t *testing.T) {
		callback := hc.NewCallback()
		triggerReq := gateway_common.HTTPTriggerRequest{
			Workflow: gateway_common.WorkflowSelector{
				WorkflowID: "0x1234567890ABCDEF1234567890abcdef12345678901234567890abcdef123456", // Contains uppercase
			},
			Input: []byte(`{"key": "value"}`),
		}
		reqBytes, err := json.Marshal(triggerReq)
		require.NoError(t, err)

		rawParams := json.RawMessage(reqBytes)
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      "test-request-id-uppercase-wf",
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &rawParams,
		}

		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callback, time.Now())
		require.Error(t, err)
		require.Contains(t, err.Error(), "workflowID must be lowercase")

		r, err := callback.Wait(t.Context())
		require.NoError(t, err)
		requireUserErrorSent(t, r, jsonrpc.ErrInvalidRequest)
	})

	t.Run("workflowOwner uppercase", func(t *testing.T) {
		callback := hc.NewCallback()
		triggerReq := gateway_common.HTTPTriggerRequest{
			Workflow: gateway_common.WorkflowSelector{
				WorkflowOwner: "0x1234567890ABCDEF1234567890abcdef12345678", // Contains uppercase
				WorkflowName:  "test-workflow",
				WorkflowTag:   "v1.0",
			},
			Input: []byte(`{"key": "value"}`),
		}
		reqBytes, err := json.Marshal(triggerReq)
		require.NoError(t, err)

		rawParams := json.RawMessage(reqBytes)
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      "test-request-id-uppercase-owner",
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &rawParams,
		}

		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callback, time.Now())
		require.Error(t, err)
		require.Contains(t, err.Error(), "workflowOwner must be lowercase")

		r, err := callback.Wait(t.Context())
		require.NoError(t, err)
		requireUserErrorSent(t, r, jsonrpc.ErrInvalidRequest)
	})

	t.Run("nil input should fail", func(t *testing.T) {
		callback := hc.NewCallback()
		triggerReq := gateway_common.HTTPTriggerRequest{
			Workflow: gateway_common.WorkflowSelector{
				WorkflowID: workflowID,
			},
			Input: nil,
		}
		reqBytes, err := json.Marshal(triggerReq)
		require.NoError(t, err)

		rawParams := json.RawMessage(reqBytes)
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      "test-request-nil-input",
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &rawParams,
		}

		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callback, time.Now())
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid params JSON")

		r, err := callback.Wait(t.Context())
		require.NoError(t, err)
		requireUserErrorSent(t, r, jsonrpc.ErrInvalidRequest)
	})

	t.Run("empty input should fail", func(t *testing.T) {
		callback := hc.NewCallback()
		triggerReq := gateway_common.HTTPTriggerRequest{
			Workflow: gateway_common.WorkflowSelector{
				WorkflowID: workflowID,
			},
			Input: []byte{},
		}
		reqBytes, err := json.Marshal(triggerReq)
		require.NoError(t, err)

		rawParams := json.RawMessage(reqBytes)
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      "test-request-empty-input",
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &rawParams,
		}

		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callback, time.Now())
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid params JSON")

		r, err := callback.Wait(t.Context())
		require.NoError(t, err)
		requireUserErrorSent(t, r, jsonrpc.ErrInvalidRequest)
	})

	t.Run("empty JSON input should pass", func(t *testing.T) {
		handler, mockDon = createTestTriggerHandler(t)
		privateKey := createTestPrivateKey(t)
		registerWorkflow(t, handler, workflowID, privateKey)

		callback := hc.NewCallback()
		triggerReq := gateway_common.HTTPTriggerRequest{
			Workflow: gateway_common.WorkflowSelector{
				WorkflowID: workflowID,
			},
			Input: []byte(`{}`),
		}
		reqBytes, err := json.Marshal(triggerReq)
		require.NoError(t, err)

		rawParams := json.RawMessage(reqBytes)
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      "test-request-empty-json-input",
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &rawParams,
		}
		req.Auth = createTestJWTToken(t, req, privateKey)

		mockDon.EXPECT().SendToNode(mock.Anything, "node1", mock.Anything).Return(nil)
		mockDon.EXPECT().SendToNode(mock.Anything, "node2", mock.Anything).Return(nil)
		mockDon.EXPECT().SendToNode(mock.Anything, "node3", mock.Anything).Return(nil)

		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callback, time.Now())
		require.NoError(t, err)
	})

	t.Run("null JSON input should fail", func(t *testing.T) {
		callback := hc.NewCallback()
		triggerReq := gateway_common.HTTPTriggerRequest{
			Workflow: gateway_common.WorkflowSelector{
				WorkflowID: workflowID,
			},
			Input: []byte(`null`),
		}
		reqBytes, err := json.Marshal(triggerReq)
		require.NoError(t, err)

		rawParams := json.RawMessage(reqBytes)
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      "test-request-null-json-input",
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &rawParams,
		}

		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callback, time.Now())
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid params JSON")

		r, err := callback.Wait(t.Context())
		require.NoError(t, err)
		requireUserErrorSent(t, r, jsonrpc.ErrInvalidRequest)
	})

	t.Run("workflowID invalid hex odd length", func(t *testing.T) {
		callback := hc.NewCallback()
		triggerReq := gateway_common.HTTPTriggerRequest{
			Workflow: gateway_common.WorkflowSelector{
				WorkflowID: "0x12345",
			},
			Input: []byte(`{"key": "value"}`),
		}
		reqBytes, err := json.Marshal(triggerReq)
		require.NoError(t, err)

		rawParams := json.RawMessage(reqBytes)
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      "test-request-short-workflow-id",
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &rawParams,
		}

		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callback, time.Now())
		require.Error(t, err)
		require.Contains(t, err.Error(), "workflowID must be a valid hex string")

		r, err := callback.Wait(t.Context())
		require.NoError(t, err)
		requireUserErrorSent(t, r, jsonrpc.ErrInvalidRequest)
	})

	t.Run("workflowOwner invalid hex odd length", func(t *testing.T) {
		callback := hc.NewCallback()
		triggerReq := gateway_common.HTTPTriggerRequest{
			Workflow: gateway_common.WorkflowSelector{
				WorkflowOwner: "0x12345",
				WorkflowName:  "test-workflow",
				WorkflowTag:   "v1.0",
			},
			Input: []byte(`{"key": "value"}`),
		}
		reqBytes, err := json.Marshal(triggerReq)
		require.NoError(t, err)

		rawParams := json.RawMessage(reqBytes)
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      "test-request-short-workflow-owner",
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &rawParams,
		}

		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callback, time.Now())
		require.Error(t, err)
		require.Contains(t, err.Error(), "workflowOwner must be a valid hex string")

		r, err := callback.Wait(t.Context())
		require.NoError(t, err)
		requireUserErrorSent(t, r, jsonrpc.ErrInvalidRequest)
	})

	t.Run("workflowName too long", func(t *testing.T) {
		callback := hc.NewCallback()
		longName := strings.Repeat("a", 65)
		triggerReq := gateway_common.HTTPTriggerRequest{
			Workflow: gateway_common.WorkflowSelector{
				WorkflowOwner: workflowOwner,
				WorkflowName:  longName,
				WorkflowTag:   "v1.0",
			},
			Input: []byte(`{"key": "value"}`),
		}
		reqBytes, err := json.Marshal(triggerReq)
		require.NoError(t, err)

		rawParams := json.RawMessage(reqBytes)
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      "test-request-long-workflow-name",
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &rawParams,
		}

		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callback, time.Now())
		require.Error(t, err)
		require.Contains(t, err.Error(), "workflowName cannot exceed 64 characters")

		r, err := callback.Wait(t.Context())
		require.NoError(t, err)
		requireUserErrorSent(t, r, jsonrpc.ErrInvalidRequest)
	})

	t.Run("workflowTag too long", func(t *testing.T) {
		callback := hc.NewCallback()
		longTag := strings.Repeat("a", 33)
		triggerReq := gateway_common.HTTPTriggerRequest{
			Workflow: gateway_common.WorkflowSelector{
				WorkflowOwner: workflowOwner,
				WorkflowName:  "test-workflow",
				WorkflowTag:   longTag,
			},
			Input: []byte(`{"key": "value"}`),
		}
		reqBytes, err := json.Marshal(triggerReq)
		require.NoError(t, err)

		rawParams := json.RawMessage(reqBytes)
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      "test-request-long-workflow-tag",
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &rawParams,
		}

		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callback, time.Now())
		require.Error(t, err)
		require.Contains(t, err.Error(), "workflowTag cannot exceed 32 characters")

		r, err := callback.Wait(t.Context())
		require.NoError(t, err)
		requireUserErrorSent(t, r, jsonrpc.ErrInvalidRequest)
	})

	t.Run("missing workflowName when workflowID not provided", func(t *testing.T) {
		callback := hc.NewCallback()
		triggerReq := gateway_common.HTTPTriggerRequest{
			Workflow: gateway_common.WorkflowSelector{
				WorkflowOwner: workflowOwner,
			},
			Input: []byte(`{"key": "value"}`),
		}
		reqBytes, err := json.Marshal(triggerReq)
		require.NoError(t, err)

		rawParams := json.RawMessage(reqBytes)
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      "test-request-missing-workflow-name",
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &rawParams,
		}

		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callback, time.Now())
		require.Error(t, err)
		require.Contains(t, err.Error(), "workflowName is required when workflowID is not provided")

		r, err := callback.Wait(t.Context())
		require.NoError(t, err)
		requireUserErrorSent(t, r, jsonrpc.ErrInvalidRequest)
	})

	t.Run("missing workflowOwner when workflowID not provided", func(t *testing.T) {
		callback := hc.NewCallback()
		triggerReq := gateway_common.HTTPTriggerRequest{
			Workflow: gateway_common.WorkflowSelector{
				WorkflowName: "test-workflow",
				WorkflowTag:  "v1.0",
			},
			Input: []byte(`{"key": "value"}`),
		}
		reqBytes, err := json.Marshal(triggerReq)
		require.NoError(t, err)

		rawParams := json.RawMessage(reqBytes)
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      "test-request-missing-workflow-owner",
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &rawParams,
		}

		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callback, time.Now())
		require.Error(t, err)
		require.Contains(t, err.Error(), "workflowOwner is required when workflowID is not provided")

		r, err := callback.Wait(t.Context())
		require.NoError(t, err)
		requireUserErrorSent(t, r, jsonrpc.ErrInvalidRequest)
	})

	t.Run("missing workflowTag when workflowID not provided", func(t *testing.T) {
		callback := hc.NewCallback()
		triggerReq := gateway_common.HTTPTriggerRequest{
			Workflow: gateway_common.WorkflowSelector{
				WorkflowOwner: workflowOwner,
				WorkflowName:  "test-workflow",
			},
			Input: []byte(`{"key": "value"}`),
		}
		reqBytes, err := json.Marshal(triggerReq)
		require.NoError(t, err)

		rawParams := json.RawMessage(reqBytes)
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      "test-request-missing-workflow-tag",
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &rawParams,
		}

		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callback, time.Now())
		require.Error(t, err)
		require.Contains(t, err.Error(), "workflowTag is required when workflowID is not provided")

		r, err := callback.Wait(t.Context())
		require.NoError(t, err)
		requireUserErrorSent(t, r, jsonrpc.ErrInvalidRequest)
	})

	t.Run("invalid hex in workflowID", func(t *testing.T) {
		callback := hc.NewCallback()
		triggerReq := gateway_common.HTTPTriggerRequest{
			Workflow: gateway_common.WorkflowSelector{
				WorkflowID: "0x1234567890abcdef1234567890abcdef12345678901234567890abcdef12345g",
			},
			Input: []byte(`{"key": "value"}`),
		}
		reqBytes, err := json.Marshal(triggerReq)
		require.NoError(t, err)

		rawParams := json.RawMessage(reqBytes)
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      "test-request-invalid-hex-workflow-id",
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &rawParams,
		}

		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callback, time.Now())
		require.Error(t, err)
		require.Contains(t, err.Error(), "workflowID must be a valid hex string")

		r, err := callback.Wait(t.Context())
		require.NoError(t, err)
		requireUserErrorSent(t, r, jsonrpc.ErrInvalidRequest)
	})

	t.Run("invalid hex in workflowOwner", func(t *testing.T) {
		callback := hc.NewCallback()
		triggerReq := gateway_common.HTTPTriggerRequest{
			Workflow: gateway_common.WorkflowSelector{
				WorkflowOwner: "0x1234567890abcdef1234567890abcdef1234567g",
				WorkflowName:  "test-workflow",
				WorkflowTag:   "v1.0",
			},
			Input: []byte(`{"key": "value"}`),
		}
		reqBytes, err := json.Marshal(triggerReq)
		require.NoError(t, err)

		rawParams := json.RawMessage(reqBytes)
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      "test-request-invalid-hex-workflow-owner",
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &rawParams,
		}

		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callback, time.Now())
		require.Error(t, err)
		require.Contains(t, err.Error(), "workflowOwner must be a valid hex string")

		r, err := callback.Wait(t.Context())
		require.NoError(t, err)
		requireUserErrorSent(t, r, jsonrpc.ErrInvalidRequest)
	})
}

func createTestTriggerRequest(workflowID string) gateway_common.HTTPTriggerRequest {
	return gateway_common.HTTPTriggerRequest{
		Workflow: gateway_common.WorkflowSelector{
			WorkflowID: workflowID,
		},
		Input: []byte(`{"key": "value"}`),
	}
}

func createTestPrivateKey(t *testing.T) *ecdsa.PrivateKey {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	return privateKey
}

func createTestJWTToken(t *testing.T, req *jsonrpc.Request[json.RawMessage], privateKey *ecdsa.PrivateKey) string {
	token, err := utils.CreateRequestJWT(*req)
	require.NoError(t, err)

	tokenString, err := token.SignedString(privateKey)
	require.NoError(t, err)
	return tokenString
}

func createTestMetadataHandler(t *testing.T) *WorkflowMetadataHandler {
	lggr := logger.Test(t)
	mockDon := handlermocks.NewDON(t)
	donConfig := &config.DONConfig{
		F: 1,
		Members: []config.NodeConfig{
			{Address: "node1"},
			{Address: "node2"},
			{Address: "node3"},
		},
	}
	cfg := WithDefaults(ServiceConfig{})
	testMetrics := createTestMetrics(t)
	return NewWorkflowMetadataHandler(lggr, cfg, mockDon, donConfig, testMetrics)
}

func createTestUserRateLimiter() limits.RateLimiter {
	return limits.UnlimitedRateLimiter()
}

func createTestTriggerHandler(t *testing.T) (*httpTriggerHandler, *handlermocks.DON) {
	cfg := ServiceConfig{
		CleanUpPeriodMs:             60000,
		MaxTriggerRequestDurationMs: 300000,
	}
	return createTestTriggerHandlerWithConfig(t, cfg)
}

func createTestTriggerHandlerWithConfig(t *testing.T, cfg ServiceConfig) (*httpTriggerHandler, *handlermocks.DON) {
	donConfig := &config.DONConfig{
		DonId: "test-don",
		F:     1, // This means we need (N+F)//2+1 = (3+1)//2+1 = 3 responses for consensus
		Members: []config.NodeConfig{
			{Address: "node1"},
			{Address: "node2"},
			{Address: "node3"},
		},
	}
	mockDon := handlermocks.NewDON(t)
	lggr := logger.Test(t)
	metadataHandler := createTestMetadataHandler(t)
	userRateLimiter := createTestUserRateLimiter()
	testMetrics := createTestMetrics(t)

	handler := NewHTTPTriggerHandler(lggr, cfg, donConfig, mockDon, metadataHandler, userRateLimiter, testMetrics)
	return handler, mockDon
}

func TestHttpTriggerHandler_HandleUserTriggerRequest_RateLimiting(t *testing.T) {
	cfg := ServiceConfig{
		CleanUpPeriodMs:             60000,
		MaxTriggerRequestDurationMs: 300000,
	}

	donConfig := &config.DONConfig{
		DonId: "test-don",
		F:     1,
		Members: []config.NodeConfig{
			{Address: "node1"},
			{Address: "node2"},
			{Address: "node3"},
		},
	}

	mockDon := handlermocks.NewDON(t)
	lggr := logger.Test(t)
	metadataHandler := createTestMetadataHandler(t)
	testMetrics := createTestMetrics(t)

	t.Run("successful rate limit check with CRE context", func(t *testing.T) {
		userRateLimiter := createTestUserRateLimiter() // Unlimited
		handler := NewHTTPTriggerHandler(lggr, cfg, donConfig, mockDon, metadataHandler, userRateLimiter, testMetrics)

		privateKey := createTestPrivateKey(t)
		workflowID := "0x1234567890abcdef1234567890abcdef12345678901234567890abcdef123456"
		workflowOwner := "0x1234567890abcdef1234567890abcdef12345678"

		// Register workflow with reference
		registerWorkflow(t, handler, workflowID, privateKey)
		handler.workflowMetadataHandler.workflowIDToRef[workflowID] = workflowReference{
			workflowOwner: workflowOwner,
			workflowName:  "test-workflow",
			workflowTag:   "v1.0",
		}

		triggerReq := gateway_common.HTTPTriggerRequest{
			Workflow: gateway_common.WorkflowSelector{
				WorkflowID: workflowID,
			},
			Input: []byte(`{"key": "value"}`),
		}
		reqBytes, err := json.Marshal(triggerReq)
		require.NoError(t, err)

		rawParams := json.RawMessage(reqBytes)
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      "test-request-id",
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &rawParams,
		}
		req.Auth = createTestJWTToken(t, req, privateKey)

		callback := hc.NewCallback()

		// Mock DON to expect sends to all nodes
		mockDon.EXPECT().SendToNode(mock.Anything, "node1", mock.Anything).Return(nil)
		mockDon.EXPECT().SendToNode(mock.Anything, "node2", mock.Anything).Return(nil)
		mockDon.EXPECT().SendToNode(mock.Anything, "node3", mock.Anything).Return(nil)

		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callback, time.Now())
		require.NoError(t, err)
	})

	t.Run("rate limit exceeded returns proper error", func(t *testing.T) {
		// Create a rate limiter with very restrictive limits
		restrictiveRateLimiter := limits.WorkflowRateLimiter(1, 0)
		handler := NewHTTPTriggerHandler(lggr, cfg, donConfig, mockDon, metadataHandler, restrictiveRateLimiter, testMetrics)

		privateKey := createTestPrivateKey(t)
		workflowID := "0x1234567890abcdef1234567890abcdef12345678901234567890abcdef123456"
		workflowOwner := "0x1234567890abcdef1234567890abcdef12345678"

		// Register workflow with reference
		registerWorkflow(t, handler, workflowID, privateKey)
		handler.workflowMetadataHandler.workflowIDToRef[workflowID] = workflowReference{
			workflowOwner: workflowOwner,
			workflowName:  "test-workflow",
			workflowTag:   "v1.0",
		}

		triggerReq := gateway_common.HTTPTriggerRequest{
			Workflow: gateway_common.WorkflowSelector{
				WorkflowID: workflowID,
			},
			Input: []byte(`{"key": "value"}`),
		}
		reqBytes, err := json.Marshal(triggerReq)
		require.NoError(t, err)

		rawParams := json.RawMessage(reqBytes)
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      "test-request-id-rate-limit",
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &rawParams,
		}
		req.Auth = createTestJWTToken(t, req, privateKey)

		callback := hc.NewCallback()

		// First request should consume the burst capacity and exceed the rate limit
		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callback, time.Now())
		require.Error(t, err)
		r, err := callback.Wait(t.Context())
		require.NoError(t, err)
		requireUserErrorSent(t, r, jsonrpc.ErrLimitExceeded)
	})
}
