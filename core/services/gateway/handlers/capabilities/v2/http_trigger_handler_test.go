package v2

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/ratelimit"
	gateway_common "github.com/smartcontractkit/chainlink-common/pkg/types/gateway"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities/v2/metrics"
	handlermocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const workflowID = "0x1234567890abcdef1234567890abcdef12345678901234567890abcdef123456"
const requestID = "test-request-id"

func createTestMetrics(t *testing.T) *metrics.Metrics {
	m, err := metrics.NewMetrics()
	require.NoError(t, err)
	return m
}

func requireUserErrorSent(t *testing.T, callbackCh chan handlers.UserCallbackPayload, errorCode int) {
	select {
	case payload := <-callbackCh:
		require.NotEmpty(t, payload.RawResponse)
		fmt.Printf("Received error payload: %+v\n", payload.RawResponse)
		require.Equal(t, api.ErrorCode(errorCode), payload.ErrorCode)
	case <-t.Context().Done():
		t.Fatal("Expected error callback")
	}
}

func TestHttpTriggerHandler_HandleUserTriggerRequest(t *testing.T) {

	triggerReq := createTestTriggerRequest()
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
		callbackCh := make(chan<- handlers.UserCallbackPayload, 1)

		// Mock DON to expect sends to all nodes
		mockDon.EXPECT().SendToNode(mock.Anything, "node1", mock.Anything).Return(nil)
		mockDon.EXPECT().SendToNode(mock.Anything, "node2", mock.Anything).Return(nil)
		mockDon.EXPECT().SendToNode(mock.Anything, "node3", mock.Anything).Return(nil)

		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callbackCh, time.Now())
		require.NoError(t, err)

		handler.callbacksMu.Lock()
		saved, exists := handler.callbacks[requestID]
		handler.callbacksMu.Unlock()

		require.True(t, exists)
		require.Equal(t, callbackCh, saved.callbackCh)
		require.NotNil(t, saved.responseAggregator)
	})

	t.Run("invalid JSON params", func(t *testing.T) {
		handler, _ := createTestTriggerHandler(t)
		callbackCh := make(chan handlers.UserCallbackPayload, 1)

		rawParams := json.RawMessage(`{invalid json}`)
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      requestID,
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &rawParams,
		}

		err := handler.HandleUserTriggerRequest(testutils.Context(t), req, callbackCh, time.Now())
		require.Error(t, err)

		requireUserErrorSent(t, callbackCh, int(jsonrpc.ErrParse))
	})

	t.Run("empty request ID", func(t *testing.T) {
		handler, _ := createTestTriggerHandler(t)
		callbackCh := make(chan handlers.UserCallbackPayload, 1)

		triggerReq := createTestTriggerRequest()
		reqBytes, err := json.Marshal(triggerReq)
		require.NoError(t, err)

		rawParams := json.RawMessage(reqBytes)
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      "", // Empty ID
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &rawParams,
		}

		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callbackCh, time.Now())
		require.Error(t, err)
		require.Contains(t, err.Error(), "empty request ID")

		requireUserErrorSent(t, callbackCh, int(jsonrpc.ErrInvalidRequest))
	})

	t.Run("request ID contains slash", func(t *testing.T) {
		handler, _ := createTestTriggerHandler(t)
		callbackCh := make(chan handlers.UserCallbackPayload, 1)

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

		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callbackCh, time.Now())
		require.Error(t, err)
		require.Contains(t, err.Error(), "must not contain '/'")

		requireUserErrorSent(t, callbackCh, int(jsonrpc.ErrInvalidRequest))
	})

	t.Run("invalid method", func(t *testing.T) {
		handler, _ := createTestTriggerHandler(t)
		callbackCh := make(chan handlers.UserCallbackPayload, 1)

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

		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callbackCh, time.Now())
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid method")

		requireUserErrorSent(t, callbackCh, int(jsonrpc.ErrMethodNotFound))
	})

	t.Run("duplicate request ID", func(t *testing.T) {
		handler, mockDon := createTestTriggerHandler(t)
		privateKey := createTestPrivateKey(t)
		registerWorkflow(t, handler, workflowID, privateKey)
		callbackCh1 := make(chan handlers.UserCallbackPayload, 1)
		callbackCh2 := make(chan handlers.UserCallbackPayload, 1)

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

		// First request should succeed
		mockDon.EXPECT().SendToNode(mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(3)
		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callbackCh1, time.Now())
		require.NoError(t, err)

		// Second request with same ID should fail
		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callbackCh2, time.Now())
		require.Error(t, err)
		require.Contains(t, err.Error(), "in-flight request")
		requireUserErrorSent(t, callbackCh2, int(jsonrpc.ErrConflict))
	})

	t.Run("invalid input JSON", func(t *testing.T) {
		handler, _ := createTestTriggerHandler(t)
		callbackCh := make(chan handlers.UserCallbackPayload, 1)

		rawParams := json.RawMessage([]byte(`{"workflow":{"workflowID":"0x1234567890abcdef1234567890abcdef12345678901234567890abcdef123456"},"input":{"invalid json"}`))
		req := &jsonrpc.Request[json.RawMessage]{
			Version: "2.0",
			ID:      requestID,
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &rawParams,
		}

		err := handler.HandleUserTriggerRequest(testutils.Context(t), req, callbackCh, time.Now())
		require.Error(t, err)
	})
}

func TestHttpTriggerHandler_HandleNodeTriggerResponse(t *testing.T) {
	t.Run("successful aggregation", func(t *testing.T) {
		handler, mockDon := createTestTriggerHandler(t)
		privateKey := createTestPrivateKey(t)
		registerWorkflow(t, handler, workflowID, privateKey)
		callbackCh := make(chan handlers.UserCallbackPayload, 1)

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
		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callbackCh, time.Now())
		require.NoError(t, err)
		// Create node responses
		rawRes := json.RawMessage(`{"result":"success"}`)
		nodeResp := &jsonrpc.Response[json.RawMessage]{
			Version: "2.0",
			ID:      requestID,
			Result:  &rawRes,
		}

		// Send responses from multiple nodes (need 2f+1 = 3 for f=1)
		err = handler.HandleNodeTriggerResponse(testutils.Context(t), nodeResp, "node1")
		require.NoError(t, err)

		err = handler.HandleNodeTriggerResponse(testutils.Context(t), nodeResp, "node2")
		require.NoError(t, err)

		// Third response should trigger aggregation
		err = handler.HandleNodeTriggerResponse(testutils.Context(t), nodeResp, "node3")
		require.NoError(t, err)

		// Check that callback was called
		select {
		case payload := <-callbackCh:
			require.NotEmpty(t, payload.RawResponse)
			require.Equal(t, api.NoError, payload.ErrorCode)

			var resp jsonrpc.Response[json.RawMessage]
			err := json.Unmarshal(payload.RawResponse, &resp)
			require.NoError(t, err)
			require.Equal(t, nodeResp.Result, resp.Result)
		case <-t.Context().Done():
			t.Fatal("Expected callback")
		}
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
		workflowOwner: "0x1234567890abcdef1234567890abcdef12345678",
		workflowName:  "test-workflow",
		workflowTag:   "v1.0",
	}
}

func TestHttpTriggerHandler_ReapExpiredCallbacks(t *testing.T) {
	workflowID := "0x1234567890abcdef1234567890abcdef12345678901234567890abcdef123456"
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
	req.Auth = createTestJWTToken(t, req, privateKey)
	cfg := ServiceConfig{
		CleanUpPeriodMs:             100,
		MaxTriggerRequestDurationMs: 50,
	}
	handler, mockDon := createTestTriggerHandlerWithConfig(t, cfg)
	registerWorkflow(t, handler, workflowID, privateKey)

	t.Run("reap expired callbacks", func(t *testing.T) {
		callbackCh := make(chan handlers.UserCallbackPayload, 1)
		mockDon.EXPECT().SendToNode(mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(3)
		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callbackCh, time.Now())
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
		callbackCh := make(chan handlers.UserCallbackPayload, 1)

		mockDon.EXPECT().SendToNode(mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(3)
		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callbackCh, time.Now())
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
		F:     1, // 1 faulty node, so 2*1+1=3 for threshold
		Members: []config.NodeConfig{
			{Address: "node1"},
			{Address: "node2"},
			{Address: "node3"},
		},
	}

	mockDon := handlermocks.NewDON(t)
	metadataHandler := createTestMetadataHandler(t)
	userRateLimiter := createTestUserRateLimiter(t)
	testMetrics := createTestMetrics(t)
	handler := NewHTTPTriggerHandler(lggr, cfg, donConfig, mockDon, metadataHandler, userRateLimiter, testMetrics)
	workflowID := "0x1234567890abcdef1234567890abcdef12345678901234567890abcdef123456"
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

		callbackCh := make(chan handlers.UserCallbackPayload, 1)

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

		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callbackCh, time.Now())
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
	workflowID := "0x1234567890abcdef1234567890abcdef12345678901234567890abcdef123456"
	key := gateway_common.AuthorizedKey{
		KeyType:   gateway_common.KeyTypeECDSAEVM,
		PublicKey: strings.ToLower(signerAddr.Hex()),
	}
	handler.workflowMetadataHandler.authorizedKeys[workflowID] = map[gateway_common.AuthorizedKey]struct{}{key: {}}
	handler.workflowMetadataHandler.workflowIDToRef[workflowID] = workflowReference{
		workflowOwner: "0x1234567890abcdef1234567890abcdef12345678",
		workflowName:  "test-workflow",
		workflowTag:   "v1.0",
	}

	t.Run("successful JWT authorization", func(t *testing.T) {
		callbackCh := make(chan handlers.UserCallbackPayload, 1)

		triggerReq := createTestTriggerRequest()
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

		err = handler.HandleUserTriggerRequest(ctx, req, callbackCh, time.Now())
		require.NoError(t, err)
		handler.callbacksMu.Lock()
		_, exists := handler.callbacks[req.ID]
		handler.callbacksMu.Unlock()
		require.True(t, exists)
	})

	t.Run("invalid JWT token", func(t *testing.T) {
		callbackCh := make(chan handlers.UserCallbackPayload, 1)

		triggerReq := createTestTriggerRequest()
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

		err = handler.HandleUserTriggerRequest(ctx, req, callbackCh, time.Now())
		require.Error(t, err)
		require.Contains(t, err.Error(), "auth failure")

		requireUserErrorSent(t, callbackCh, int(jsonrpc.ErrInvalidRequest))
	})

	t.Run("unauthorized signer", func(t *testing.T) {
		callbackCh := make(chan handlers.UserCallbackPayload, 1)
		unauthorizedKey := createTestPrivateKey(t)

		triggerReq := createTestTriggerRequest()
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

		err = handler.HandleUserTriggerRequest(ctx, req, callbackCh, time.Now())
		require.Error(t, err)
		require.Contains(t, err.Error(), "auth failure")

		requireUserErrorSent(t, callbackCh, int(jsonrpc.ErrInvalidRequest))
	})

	t.Run("workflow not found", func(t *testing.T) {
		callbackCh := make(chan handlers.UserCallbackPayload, 1)

		triggerReq := gateway_common.HTTPTriggerRequest{
			Workflow: gateway_common.WorkflowSelector{
				WorkflowID: "0xabcdef12345678901234567890ab",
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

		err = handler.HandleUserTriggerRequest(ctx, req, callbackCh, time.Now())
		require.Error(t, err)
		require.Contains(t, err.Error(), "auth failure")

		requireUserErrorSent(t, callbackCh, int(jsonrpc.ErrInvalidRequest))
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

	workflowID := "0x1234567890abcdef1234567890abcdef12345678901234567890abcdef123456"
	workflowOwner := "0x1234567890abcdef1234567890abcdef12345678"
	workflowName := "0x" + hex.EncodeToString([]byte(workflows.HashTruncateName("test-workflow")))
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
		callbackCh := make(chan handlers.UserCallbackPayload, 1)

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

		err = handler.HandleUserTriggerRequest(ctx, req, callbackCh, time.Now())
		require.NoError(t, err)
	})

	t.Run("workflow not found by name", func(t *testing.T) {
		callbackCh := make(chan handlers.UserCallbackPayload, 1)

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

		err = handler.HandleUserTriggerRequest(ctx, req, callbackCh, time.Now())
		require.Error(t, err)
		require.Contains(t, err.Error(), "workflow not found")

		requireUserErrorSent(t, callbackCh, int(jsonrpc.ErrInvalidRequest))
	})
}
func TestHttpTriggerHandler_HandleUserTriggerRequest_Validation(t *testing.T) {
	handler, _ := createTestTriggerHandler(t)
	callbackCh := make(chan handlers.UserCallbackPayload, 1)

	t.Run("workflowID without 0x prefix", func(t *testing.T) {
		triggerReq := gateway_common.HTTPTriggerRequest{
			Workflow: gateway_common.WorkflowSelector{
				WorkflowID: "1234567890abcdef1234567890abcdef12345678901234567890abcdef123456", // Missing 0x
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

		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callbackCh, time.Now())
		require.Error(t, err)
		require.Contains(t, err.Error(), "workflowID must be prefixed with '0x'")

		requireUserErrorSent(t, callbackCh, int(jsonrpc.ErrInvalidRequest))
	})

	t.Run("workflowOwner without 0x prefix", func(t *testing.T) {
		triggerReq := gateway_common.HTTPTriggerRequest{
			Workflow: gateway_common.WorkflowSelector{
				WorkflowOwner: "1234567890abcdef1234567890abcdef12345678", // Missing 0x
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
			ID:      "test-request-id-2",
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &rawParams,
		}

		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callbackCh, time.Now())
		require.Error(t, err)
		require.Contains(t, err.Error(), "workflowOwner must be prefixed with '0x'")

		requireUserErrorSent(t, callbackCh, int(jsonrpc.ErrInvalidRequest))
	})

	t.Run("workflowID uppercase", func(t *testing.T) {
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

		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callbackCh, time.Now())
		require.Error(t, err)
		require.Contains(t, err.Error(), "workflowID must be lowercase")

		requireUserErrorSent(t, callbackCh, int(jsonrpc.ErrInvalidRequest))
	})

	t.Run("workflowOwner uppercase", func(t *testing.T) {
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

		err = handler.HandleUserTriggerRequest(testutils.Context(t), req, callbackCh, time.Now())
		require.Error(t, err)
		require.Contains(t, err.Error(), "workflowOwner must be lowercase")

		requireUserErrorSent(t, callbackCh, int(jsonrpc.ErrInvalidRequest))
	})
}

func createTestTriggerRequest() gateway_common.HTTPTriggerRequest {
	return gateway_common.HTTPTriggerRequest{
		Workflow: gateway_common.WorkflowSelector{
			WorkflowID: "0x1234567890abcdef1234567890abcdef12345678901234567890abcdef123456",
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

func createTestUserRateLimiter(t *testing.T) *ratelimit.RateLimiter {
	cfg := ratelimit.RateLimiterConfig{
		GlobalRPS:      50,
		GlobalBurst:    50,
		PerSenderRPS:   5,
		PerSenderBurst: 5,
	}
	limiter, err := ratelimit.NewRateLimiter(cfg)
	require.NoError(t, err)
	return limiter
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
		F:     1, // This means we need 2f+1 = 3 responses for consensus
		Members: []config.NodeConfig{
			{Address: "node1"},
			{Address: "node2"},
			{Address: "node3"},
		},
	}
	mockDon := handlermocks.NewDON(t)
	lggr := logger.Test(t)
	metadataHandler := createTestMetadataHandler(t)
	userRateLimiter := createTestUserRateLimiter(t)
	testMetrics := createTestMetrics(t)

	handler := NewHTTPTriggerHandler(lggr, cfg, donConfig, mockDon, metadataHandler, userRateLimiter, testMetrics)
	return handler, mockDon
}
