package v2

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
	"sync/atomic"
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
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/common/aggregation"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities/v2/metrics"
	hc "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
	handlermocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const (
	workflowID    = "0x1234567890abcdef1234567890abcdef12345678901234567890abcdef123456"
	workflowOwner = "0x1234567890abcdef1234567890abcdef12345678"
	requestID     = "test-request-id"
)

func createTestMetrics(t *testing.T, donConfig *config.DONConfig) *metrics.Metrics {
	m, err := metrics.NewMetrics(donConfig.Members)
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

		err = handler.HandleUserTriggerRequest(t.Context(), req, callback, time.Now())
		require.NoError(t, err)

		handler.callbacksMu.Lock()
		saved, exists := handler.callbacks[requestID]
		handler.callbacksMu.Unlock()

		require.True(t, exists)
		require.Equal(t, callback, saved.Callback)
		require.NotNil(t, saved.responseAggregators)
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

		err = handler.HandleUserTriggerRequest(t.Context(), req, callback, time.Now())
		require.NoError(t, err)

		handler.callbacksMu.Lock()
		saved, exists := handler.callbacks[requestID]
		handler.callbacksMu.Unlock()

		require.True(t, exists)
		require.Equal(t, callback, saved.Callback)
		require.NotNil(t, saved.responseAggregators)
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

		err = handler.HandleUserTriggerRequest(t.Context(), req, callback, time.Now())
		require.NoError(t, err)

		handler.callbacksMu.Lock()
		saved, exists := handler.callbacks[requestID]
		handler.callbacksMu.Unlock()

		require.True(t, exists)
		require.Equal(t, callback, saved.Callback)
		require.NotNil(t, saved.responseAggregators)
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

		err = handler.HandleUserTriggerRequest(t.Context(), req, callback, time.Now())
		require.NoError(t, err)

		handler.callbacksMu.Lock()
		saved, exists := handler.callbacks[requestID]
		handler.callbacksMu.Unlock()

		require.True(t, exists)
		require.Equal(t, callback, saved.Callback)
		require.NotNil(t, saved.responseAggregators)
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

		err := handler.HandleUserTriggerRequest(t.Context(), req, callback, time.Now())
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

		err := handler.HandleUserTriggerRequest(t.Context(), req, callback, time.Now())
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

		err = handler.HandleUserTriggerRequest(t.Context(), req, callback, time.Now())
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

		err = handler.HandleUserTriggerRequest(t.Context(), req, callback, time.Now())
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

		err = handler.HandleUserTriggerRequest(t.Context(), req, callback, time.Now())
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
		err = handler.HandleUserTriggerRequest(t.Context(), req, callback1, time.Now())
		require.NoError(t, err)

		// Second request with same ID should fail
		req.Auth = createTestJWTToken(t, req, privateKey)
		err = handler.HandleUserTriggerRequest(t.Context(), req, callback2, time.Now())
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
		err = handler.HandleUserTriggerRequest(t.Context(), req, callback1, time.Now())
		require.NoError(t, err)

		// Second request with same ID should fail
		err = handler.HandleUserTriggerRequest(t.Context(), req, callback2, time.Now())
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

		err := handler.HandleUserTriggerRequest(t.Context(), req, callback, time.Now())
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
		err = handler.HandleUserTriggerRequest(t.Context(), req, callback, time.Now())
		require.NoError(t, err)
		// Create node responses
		rawRes := json.RawMessage(`{"result":"success"}`)
		nodeResp := &jsonrpc.Response[json.RawMessage]{
			Version: "2.0",
			ID:      requestID,
			Result:  &rawRes,
		}

		// Send responses from multiple nodes (need (N+F)//2+1 = (3+1)//2+1 = 3 for N=3, F=1)
		err = handler.HandleNodeTriggerResponse(t.Context(), nodeResp, "node1")
		require.NoError(t, err)

		err = handler.HandleNodeTriggerResponse(t.Context(), nodeResp, "node2")
		require.NoError(t, err)

		// Third response should trigger aggregation
		err = handler.HandleNodeTriggerResponse(t.Context(), nodeResp, "node3")
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

		err := handler.HandleNodeTriggerResponse(t.Context(), nodeResp, "node1")
		require.Error(t, err)
		require.Contains(t, err.Error(), "callback not found")
	})
}

func TestHttpTriggerHandler_ServiceLifecycle(t *testing.T) {
	t.Run("start and stop", func(t *testing.T) {
		handler, _ := createTestTriggerHandler(t)

		ctx := t.Context()
		err := handler.Start(ctx)
		require.NoError(t, err)

		err = handler.Close()
		require.NoError(t, err)
	})

	t.Run("double start and close should errors", func(t *testing.T) {
		handler, _ := createTestTriggerHandler(t)

		ctx := t.Context()
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

func registerWorkflow(_ *testing.T, handler *httpTriggerHandler, workflowID string, privateKey *ecdsa.PrivateKey) {
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
	// Assign the workflow to every shard the metadata handler knows about so
	// that WorkflowShards(workflowID) returns them and the trigger handler's
	// fan-out (setupCallback/sendWithRetries) reaches all shards. The default
	// test harness uses a single shard, so this registers on that one shard;
	// registerWorkflowOnShards overrides this with a specific subset.
	assignWorkflowToAllShards(handler.workflowMetadataHandler, workflowID)
}

// assignWorkflowToAllShards writes workflowShards[workflowID] = copy(mh.shards)
// under the metadata handler's mutex. Used by registerWorkflow and by tests
// that populate the metadata maps directly.
func assignWorkflowToAllShards(mh *WorkflowMetadataHandler, workflowID string) {
	mh.mu.Lock()
	defer mh.mu.Unlock()
	mh.workflowShards[workflowID] = append([]*shardEndpoint(nil), mh.shards...)
}

// registerWorkflowOnShards is the sharding-aware variant of registerWorkflow.
// It registers a workflow's authorized key and selector, and assigns the
// workflow to the given shard endpoints so that WorkflowShards(workflowID)
// returns them and the trigger handler fans the request out only to those
// shards. The metadata handler's map is written directly under its mutex.
func registerWorkflowOnShards(t *testing.T, handler *httpTriggerHandler, workflowID string, privateKey *ecdsa.PrivateKey, assignedShards ...*shardEndpoint) {
	t.Helper()
	registerWorkflow(t, handler, workflowID, privateKey)

	handler.workflowMetadataHandler.mu.Lock()
	defer handler.workflowMetadataHandler.mu.Unlock()
	assigned := make([]*shardEndpoint, len(assignedShards))
	copy(assigned, assignedShards)
	handler.workflowMetadataHandler.workflowShards[workflowID] = assigned
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
		err = handler.HandleUserTriggerRequest(t.Context(), req, callback, time.Now())
		require.NoError(t, err)

		// Manually set the callback's createdAt to the past to simulate expiration
		handler.callbacksMu.Lock()
		if cb, exists := handler.callbacks[requestID]; exists {
			cb.createdAt = time.Now().Add(-time.Duration(cfg.CleanUpPeriodMs+1) * time.Millisecond)
			handler.callbacks[requestID] = cb
		}
		handler.callbacksMu.Unlock()

		// Manually trigger reaping
		handler.reapExpiredCallbacks(t.Context())

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
		err = handler.HandleUserTriggerRequest(t.Context(), req, callback, time.Now())
		require.NoError(t, err)

		// Optionally, set createdAt to now (should not be expired)
		handler.callbacksMu.Lock()
		if cb, exists := handler.callbacks[requestID]; exists {
			cb.createdAt = time.Now()
			handler.callbacks[requestID] = cb
		}
		handler.callbacksMu.Unlock()

		// Manually trigger reaping
		handler.reapExpiredCallbacks(t.Context())

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
	if testing.Short() {
		t.Skip("too slow for testing.Short")
	}

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
	testMetrics := createTestMetrics(t, donConfig)
	handler := newTestTriggerHandler(t, lggr, cfg, donConfig, mockDon, metadataHandler, userRateLimiter, testMetrics)
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

		err := handler.Start(t.Context())
		require.NoError(t, err)

		err = handler.HandleUserTriggerRequest(t.Context(), req, callback, time.Now())
		require.NoError(t, err)

		mockDon.AssertExpectations(t)
		err = handler.Close()
		require.NoError(t, err)
	})
}

func TestHttpTriggerHandler_HandleUserTriggerRequest_SendsToNodesInParallel(t *testing.T) {
	if testing.Short() {
		t.Skip("too slow for testing.Short")
	}

	t.Parallel()

	lggr := logger.Test(t)
	const nodeDelay = 200 * time.Millisecond
	cfg := WithDefaults(ServiceConfig{
		MaxTriggerRequestDurationMs: 5000,
		NodeSendTimeoutMs:           5000,
	})

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
	metadataHandler := createTestMetadataHandler(t)
	userRateLimiter := createTestUserRateLimiter()
	testMetrics := createTestMetrics(t, donConfig)
	handler := newTestTriggerHandler(t, lggr, cfg, donConfig, mockDon, metadataHandler, userRateLimiter, testMetrics)
	privateKey := createTestPrivateKey(t)
	registerWorkflow(t, handler, workflowID, privateKey)

	rawParams := json.RawMessage(`{"input":{},"workflow":{"workflowID":"0x1234567890abcdef1234567890abcdef12345678901234567890abcdef123456"}}`)
	req := &jsonrpc.Request[json.RawMessage]{
		ID:      "test-request-parallel",
		Method:  gateway_common.MethodWorkflowExecute,
		Params:  &rawParams,
		Version: "2.0",
	}
	req.Auth = createTestJWTToken(t, req, privateKey)
	callback := hc.NewCallback()

	// Every node takes nodeDelay to respond. If sends were sequential the whole
	// round would take ~3*nodeDelay; in parallel it should take ~nodeDelay.
	for _, addr := range []string{"node1", "node2", "node3"} {
		mockDon.On("SendToNode", mock.Anything, addr, mock.Anything).Run(func(args mock.Arguments) {
			time.Sleep(nodeDelay)
		}).Return(nil).Once()
	}

	err := handler.Start(t.Context())
	require.NoError(t, err)

	start := time.Now()
	err = handler.HandleUserTriggerRequest(t.Context(), req, callback, time.Now())
	require.NoError(t, err)
	elapsed := time.Since(start)

	require.Less(t, elapsed, 2*nodeDelay, "sends to DON members should happen in parallel, not sequentially")

	mockDon.AssertExpectations(t)
	err = handler.Close()
	require.NoError(t, err)
}

func TestHttpTriggerHandler_HandleUserTriggerRequest_SlowNodeDoesNotBlockOthers(t *testing.T) {
	t.Parallel()

	lggr := logger.Test(t)
	const nodeSendTimeoutMs = 100
	cfg := WithDefaults(ServiceConfig{
		MaxTriggerRequestDurationMs: 5000,
		NodeSendTimeoutMs:           nodeSendTimeoutMs,
		RetryConfig: RetryConfig{
			InitialIntervalMs: 10,
			MaxIntervalTimeMs: 50,
			Multiplier:        2,
		},
	})

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
	metadataHandler := createTestMetadataHandler(t)
	userRateLimiter := createTestUserRateLimiter()
	testMetrics := createTestMetrics(t, donConfig)
	handler := newTestTriggerHandler(t, lggr, cfg, donConfig, mockDon, metadataHandler, userRateLimiter, testMetrics)
	privateKey := createTestPrivateKey(t)
	registerWorkflow(t, handler, workflowID, privateKey)

	rawParams := json.RawMessage(`{"input":{},"workflow":{"workflowID":"0x1234567890abcdef1234567890abcdef12345678901234567890abcdef123456"}}`)
	req := &jsonrpc.Request[json.RawMessage]{
		ID:      "test-request-slow-node",
		Method:  gateway_common.MethodWorkflowExecute,
		Params:  &rawParams,
		Version: "2.0",
	}
	req.Auth = createTestJWTToken(t, req, privateKey)
	callback := hc.NewCallback()

	// node2 and node3 respond immediately and must not be re-sent even though
	// node1 hangs past its per-node timeout on the first attempt.
	mockDon.On("SendToNode", mock.Anything, "node2", mock.Anything).Return(nil).Once()
	mockDon.On("SendToNode", mock.Anything, "node3", mock.Anything).Return(nil).Once()

	var node1Attempts atomic.Int32
	mockDon.EXPECT().SendToNode(mock.Anything, "node1", mock.Anything).RunAndReturn(
		func(ctx context.Context, nodeAddress string, req *jsonrpc.Request[json.RawMessage]) error {
			if node1Attempts.Add(1) == 1 {
				deadline, ok := ctx.Deadline()
				require.True(t, ok, "per-node context should carry a deadline")
				require.WithinDuration(t, time.Now().Add(nodeSendTimeoutMs*time.Millisecond), deadline, 50*time.Millisecond,
					"per-node timeout should be applied, not the overall request duration")

				<-ctx.Done()
				return ctx.Err()
			}
			return nil
		}).Twice()

	err := handler.Start(t.Context())
	require.NoError(t, err)

	start := time.Now()
	err = handler.HandleUserTriggerRequest(t.Context(), req, callback, time.Now())
	require.NoError(t, err)
	elapsed := time.Since(start)

	require.Less(t, elapsed, 2*time.Second,
		"a hung node should be retried after its per-node timeout, not after the full request duration")

	mockDon.AssertExpectations(t)
	err = handler.Close()
	require.NoError(t, err)
}

// TestHttpTriggerHandler_HandleUserTriggerRequest_DeliversOnlyToRegisteredShards
// verifies the sharding fan-out behavior introduced by the sharding commit.
//
// Setup: a single DON with 3 shards (3 disjoint node sets, one connection
// manager per shard). Workflow X is registered on 2 of the 3 shards. A user
// trigger request for workflow X must be delivered to the members of those 2
// shards only; the 3rd shard's connection manager must never be contacted, and
// the per-request callback must carry exactly 2 response aggregators (one per
// assigned shard).
func TestHttpTriggerHandler_HandleUserTriggerRequest_DeliversOnlyToRegisteredShards(t *testing.T) {
	t.Parallel()

	lggr := logger.Test(t)
	cfg := WithDefaults(ServiceConfig{
		MaxTriggerRequestDurationMs: 5000,
		NodeSendTimeoutMs:           5000,
	})

	// Three shards, each with 3 disjoint nodes. donIDs: "don", "don_1", "don_2".
	shardedDONs := []config.ShardedDONConfig{
		{
			DonName: "don",
			F:       1, // (N+F)//2+1 = (3+1)//2+1 = 3 for quorum
			Shards: []config.Shard{
				{Nodes: []config.NodeConfig{{Address: "n1"}, {Address: "n2"}, {Address: "n3"}}},
				{Nodes: []config.NodeConfig{{Address: "n4"}, {Address: "n5"}, {Address: "n6"}}},
				{Nodes: []config.NodeConfig{{Address: "n7"}, {Address: "n8"}, {Address: "n9"}}},
			},
		},
	}
	shard0Don := handlermocks.NewDON(t)
	shard1Don := handlermocks.NewDON(t)
	shard2Don := handlermocks.NewDON(t)
	connMgrs := [][]handlers.DON{{shard0Don, shard1Don, shard2Don}}

	shards, nodeAddrToShard, err := buildShardEndpoints(shardedDONs, connMgrs)
	require.NoError(t, err)
	require.Len(t, shards, 3)
	// Sanity-check donID derivation used by the metadata/trigger handlers.
	require.Equal(t, "don", shards[0].donID)
	require.Equal(t, "don_1", shards[1].donID)
	require.Equal(t, "don_2", shards[2].donID)

	allMembersSlice := allMembers(shards)
	testMetrics, err := metrics.NewMetrics(allMembersSlice)
	require.NoError(t, err)
	metadataHandler := NewWorkflowMetadataHandler(lggr, cfg, shards, nodeAddrToShard, testMetrics)
	userRateLimiter := createTestUserRateLimiter()
	handler := NewHTTPTriggerHandler(lggr, cfg, shards, nodeAddrToShard, metadataHandler, userRateLimiter, testMetrics)

	privateKey := createTestPrivateKey(t)
	// Register workflow X on exactly 2 of the 3 shards (shard0 and shard1).
	registerWorkflowOnShards(t, handler, workflowID, privateKey, shards[0], shards[1])

	// Build and sign the user trigger request for workflow X.
	triggerReq := createTestTriggerRequest(workflowID)
	reqBytes, err := json.Marshal(triggerReq)
	require.NoError(t, err)
	rawParams := json.RawMessage(reqBytes)
	req := &jsonrpc.Request[json.RawMessage]{
		Version: "2.0",
		ID:      "test-request-sharded",
		Method:  gateway_common.MethodWorkflowExecute,
		Params:  &rawParams,
	}
	req.Auth = createTestJWTToken(t, req, privateKey)
	callback := hc.NewCallback()

	// The two registered shards' connection managers must each be hit once per
	// member (sends succeed on the first attempt, so no retries).
	for _, addr := range []string{"n1", "n2", "n3"} {
		shard0Don.EXPECT().SendToNode(mock.Anything, addr, mock.Anything).Return(nil).Once()
	}
	for _, addr := range []string{"n4", "n5", "n6"} {
		shard1Don.EXPECT().SendToNode(mock.Anything, addr, mock.Anything).Return(nil).Once()
	}
	// shard2Don gets NO expectations: any call to it would be a routing bug.

	require.NoError(t, handler.Start(t.Context()))
	t.Cleanup(func() { _ = handler.Close() })

	err = handler.HandleUserTriggerRequest(t.Context(), req, callback, time.Now())
	require.NoError(t, err)

	// Assert the callback was set up with one aggregator per *registered* shard.
	handler.callbacksMu.Lock()
	saved, exists := handler.callbacks[req.ID]
	handler.callbacksMu.Unlock()
	require.True(t, exists, "callback should be registered for the request ID")
	require.Len(t, saved.responseAggregators, 2, "exactly 2 shards are registered for this workflow")
	require.Contains(t, saved.responseAggregators, "don")
	require.Contains(t, saved.responseAggregators, "don_1")
	require.NotContains(t, saved.responseAggregators, "don_2")

	// Assert the mock connection managers saw exactly the expected sends. The
	// cleanup registered by NewDON will call AssertExpectations on test teardown,
	// which would fail if shard2Don had been contacted.
	shard0Don.AssertExpectations(t)
	shard1Don.AssertExpectations(t)
	shard2Don.AssertExpectations(t)
}

func TestHttpTriggerHandler_HandleUserTriggerRequest_JWTAuthorization(t *testing.T) {
	handler, mockDon := createTestTriggerHandler(t)
	ctx := t.Context()

	// Setup metadata handler with test data
	err := handler.workflowMetadataHandler.aggs[handler.workflowMetadataHandler.shards[0].donID].Start(ctx)
	require.NoError(t, err)
	defer handler.workflowMetadataHandler.aggs[handler.workflowMetadataHandler.shards[0].donID].Close()

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
	// Assign the workflow to all shards so setupCallback/sendWithRetries can
	// fan the request out (these tests populate the metadata maps directly
	// instead of calling registerWorkflow).
	assignWorkflowToAllShards(handler.workflowMetadataHandler, workflowID)

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
		require.Contains(t, err.Error(), "workflow not found")

		r, err2 := callback.Wait(t.Context())
		require.NoError(t, err2)
		requireUserErrorSent(t, r, jsonrpc.ErrInvalidRequest)
	})
}

func TestHttpTriggerHandler_HandleUserTriggerRequest_WorkflowLookup(t *testing.T) {
	handler, mockDon := createTestTriggerHandler(t)
	ctx := t.Context()

	err := handler.workflowMetadataHandler.aggs[handler.workflowMetadataHandler.shards[0].donID].Start(ctx)
	require.NoError(t, err)
	defer handler.workflowMetadataHandler.aggs[handler.workflowMetadataHandler.shards[0].donID].Close()

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
	// Assign the workflow to all shards so setupCallback/sendWithRetries can
	// fan the request out (this test populates the metadata maps directly
	// instead of calling registerWorkflow).
	assignWorkflowToAllShards(handler.workflowMetadataHandler, workflowID)

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

		err = handler.HandleUserTriggerRequest(t.Context(), req, callback, time.Now())
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

		err = handler.HandleUserTriggerRequest(t.Context(), req, callback, time.Now())
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

		err = handler.HandleUserTriggerRequest(t.Context(), req, callback, time.Now())
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

		err = handler.HandleUserTriggerRequest(t.Context(), req, callback, time.Now())
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

		err = handler.HandleUserTriggerRequest(t.Context(), req, callback, time.Now())
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

		err = handler.HandleUserTriggerRequest(t.Context(), req, callback, time.Now())
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

		err = handler.HandleUserTriggerRequest(t.Context(), req, callback, time.Now())
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

		err = handler.HandleUserTriggerRequest(t.Context(), req, callback, time.Now())
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

		err = handler.HandleUserTriggerRequest(t.Context(), req, callback, time.Now())
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

		err = handler.HandleUserTriggerRequest(t.Context(), req, callback, time.Now())
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

		err = handler.HandleUserTriggerRequest(t.Context(), req, callback, time.Now())
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

		err = handler.HandleUserTriggerRequest(t.Context(), req, callback, time.Now())
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

		err = handler.HandleUserTriggerRequest(t.Context(), req, callback, time.Now())
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

		err = handler.HandleUserTriggerRequest(t.Context(), req, callback, time.Now())
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

		err = handler.HandleUserTriggerRequest(t.Context(), req, callback, time.Now())
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
	testMetrics := createTestMetrics(t, donConfig)
	shardedDONs := []config.ShardedDONConfig{
		{DonName: donConfig.DonId, F: donConfig.F, Shards: []config.Shard{{Nodes: donConfig.Members}}},
	}
	shards, nodeAddrToShard, err := buildShardEndpoints(shardedDONs, [][]handlers.DON{{mockDon}})
	require.NoError(t, err)
	return NewWorkflowMetadataHandler(lggr, cfg, shards, nodeAddrToShard, testMetrics)
}

func createTestUserRateLimiter() limits.RateLimiter {
	return limits.UnlimitedRateLimiter()
}

func newTestTriggerHandler(t *testing.T, lggr logger.Logger, cfg ServiceConfig, donConfig *config.DONConfig, mockDon *handlermocks.DON, metadataHandler *WorkflowMetadataHandler, userRateLimiter limits.RateLimiter, testMetrics *metrics.Metrics) *httpTriggerHandler {
	shardedDONs := []config.ShardedDONConfig{
		{DonName: donConfig.DonId, F: donConfig.F, Shards: []config.Shard{{Nodes: donConfig.Members}}},
	}
	shards, nodeAddrToShard, err := buildShardEndpoints(shardedDONs, [][]handlers.DON{{mockDon}})
	require.NoError(t, err)
	// Repoint the metadata handler at THIS handler's shard set so that the
	// workflowShards populated by registerWorkflow route sends to mockDon. The
	// metadata handler may have been built with its own shard set (and donIDs) by
	// createTestMetadataHandler, so rebuild its per-shard aggregators keyed by the
	// shared shards' donIDs. (donIDs differ when createTestMetadataHandler's
	// donConfig has an empty DonId.)
	metadataHandler.shards = shards
	metadataHandler.nodeAddrToShard = nodeAddrToShard
	metadataHandler.aggs = make(map[string]*aggregation.WorkflowMetadataAggregator, len(shards))
	for _, shard := range shards {
		threshold := shard.f + 1
		metadataHandler.aggs[shard.donID] = aggregation.NewWorkflowMetadataAggregator(metadataHandler.lggr, threshold, time.Duration(cfg.CleanUpPeriodMs)*time.Millisecond, testMetrics)
	}
	return NewHTTPTriggerHandler(lggr, cfg, shards, nodeAddrToShard, metadataHandler, userRateLimiter, testMetrics)
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
	testMetrics := createTestMetrics(t, donConfig)

	// Build ONE shared shard set so the metadata handler and the trigger handler
	// reference the SAME shardEndpoint instances (and the same mockDon). This is
	// required because registerWorkflow writes workflowShards to the metadata
	// handler's .shards, and sendWithRetries reads them back and sends via
	// shard.connMgr — both must hit mockDon, which the test sets expectations on.
	shardedDONs := []config.ShardedDONConfig{
		{DonName: donConfig.DonId, F: donConfig.F, Shards: []config.Shard{{Nodes: donConfig.Members}}},
	}
	shards, nodeAddrToShard, err := buildShardEndpoints(shardedDONs, [][]handlers.DON{{mockDon}})
	require.NoError(t, err)

	metadataHandler := NewWorkflowMetadataHandler(lggr, WithDefaults(cfg), shards, nodeAddrToShard, testMetrics)
	userRateLimiter := createTestUserRateLimiter()
	handler := NewHTTPTriggerHandler(lggr, cfg, shards, nodeAddrToShard, metadataHandler, userRateLimiter, testMetrics)
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
	testMetrics := createTestMetrics(t, donConfig)

	t.Run("successful rate limit check with CRE context", func(t *testing.T) {
		userRateLimiter := createTestUserRateLimiter() // Unlimited
		handler := newTestTriggerHandler(t, lggr, cfg, donConfig, mockDon, metadataHandler, userRateLimiter, testMetrics)

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

		err = handler.HandleUserTriggerRequest(t.Context(), req, callback, time.Now())
		require.NoError(t, err)
	})

	t.Run("rate limit exceeded returns proper error", func(t *testing.T) {
		// Create a rate limiter with very restrictive limits
		restrictiveRateLimiter := limits.WorkflowRateLimiter(1, 0)
		handler := newTestTriggerHandler(t, lggr, cfg, donConfig, mockDon, metadataHandler, restrictiveRateLimiter, testMetrics)

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
		err = handler.HandleUserTriggerRequest(t.Context(), req, callback, time.Now())
		require.Error(t, err)
		r, err := callback.Wait(t.Context())
		require.NoError(t, err)
		requireUserErrorSent(t, r, jsonrpc.ErrLimitExceeded)
	})
}

func TestHttpTriggerHandler_HandleUserTriggerRequest_StopsRetriesOnQuorum(t *testing.T) {
	lggr := logger.Test(t)
	cfg := WithDefaults(ServiceConfig{})

	// 4 nodes, 1 faulty node, so (N+F)//2+1=(4+1)//2+1=3 for threshold
	// Quorum is reached when 3 nodes respond.
	donConfig := &config.DONConfig{
		DonId: "test-don",
		F:     1,
		Members: []config.NodeConfig{
			{Address: "node1"},
			{Address: "node2"},
			{Address: "node3"},
			{Address: "node4"},
		},
	}

	mockDon := handlermocks.NewDON(t)
	metadataHandler := createTestMetadataHandler(t)
	userRateLimiter := createTestUserRateLimiter()
	testMetrics := createTestMetrics(t, donConfig)
	handler := newTestTriggerHandler(t, lggr, cfg, donConfig, mockDon, metadataHandler, userRateLimiter, testMetrics)
	privateKey := createTestPrivateKey(t)
	registerWorkflow(t, handler, workflowID, privateKey)

	t.Run("stops retries when quorum is reached and callback is responded", func(t *testing.T) {
		rawParams := json.RawMessage(`{"input":{},"workflow":{"workflowID":"0x1234567890abcdef1234567890abcdef12345678901234567890abcdef123456"}}`)
		req := &jsonrpc.Request[json.RawMessage]{
			ID:      "test-request-quorum-stop",
			Method:  gateway_common.MethodWorkflowExecute,
			Params:  &rawParams,
			Version: "2.0",
		}
		req.Auth = createTestJWTToken(t, req, privateKey)
		callback := hc.NewCallback()

		// Use channel to signal when initial broadcast is complete
		broadcastComplete := make(chan struct{})
		var callCount atomic.Int64

		// Setup: node1, node2, node3 succeed, node4 fails indefinitely
		mockDon.On("SendToNode", mock.Anything, "node1", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			if callCount.Add(1) == 3 {
				close(broadcastComplete)
			}
		}).Once()
		mockDon.On("SendToNode", mock.Anything, "node2", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			if callCount.Add(1) == 3 {
				close(broadcastComplete)
			}
		}).Once()
		mockDon.On("SendToNode", mock.Anything, "node3", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			if callCount.Add(1) == 3 {
				close(broadcastComplete)
			}
		}).Once()
		mockDon.On("SendToNode", mock.Anything, "node4", mock.Anything).Return(errors.New("connection error"))

		err := handler.Start(t.Context())
		require.NoError(t, err)

		// Start the trigger request in a goroutine
		errCh := make(chan error, 1)
		go func() {
			errCh <- handler.HandleUserTriggerRequest(t.Context(), req, callback, time.Now())
		}()

		// Wait for initial broadcast
		select {
		case <-broadcastComplete:
		case <-t.Context().Done():
			t.Fatal("Context cancelled waiting for initial broadcast to complete")
		}

		rawRes := json.RawMessage(`{"result":"ACCEPTED"}`)
		nodeResp := &jsonrpc.Response[json.RawMessage]{
			Version: "2.0",
			ID:      req.ID,
			Result:  &rawRes,
		}

		// Send responses from node1 and node2 to reach threshold (need 3 for quorum with F=1).
		// Node4 is not included in the quorum since it failed to connect.
		err = handler.HandleNodeTriggerResponse(t.Context(), nodeResp, "node1")
		require.NoError(t, err)

		err = handler.HandleNodeTriggerResponse(t.Context(), nodeResp, "node2")
		require.NoError(t, err)

		err = handler.HandleNodeTriggerResponse(t.Context(), nodeResp, "node3")
		require.NoError(t, err)

		payload, err := callback.Wait(t.Context())
		require.NoError(t, err)
		require.NotEmpty(t, payload.RawResponse)
		require.Equal(t, api.NoError, payload.ErrorCode)

		select {
		case err = <-errCh:
			require.NoError(t, err)
		case <-t.Context().Done():
			t.Fatal("Context cancelled waiting for HandleUserTriggerRequest to complete")
		}

		// After the response is sent, the callback is kept and marked processed
		// (rather than removed) so that late responses are recognized.
		handler.callbacksMu.Lock()
		saved, exists := handler.callbacks[req.ID]
		handler.callbacksMu.Unlock()
		require.True(t, exists, "callback should be kept after response is sent so late responses are recognized")
		require.True(t, saved.processed, "callback should be marked as processed after response is sent")

		err = handler.Close()
		require.NoError(t, err)
		mockDon.AssertExpectations(t)
	})
}

// createMultiShardTriggerHandler builds a trigger handler over multiple shards
// of a single DON. Each shard gets its own mock DON connection manager. Returns
// the handler, the mock DONs (one per shard), and the shard endpoints.
func createMultiShardTriggerHandler(t *testing.T, donName string, shardNodeSets [][]string, f int) (*httpTriggerHandler, []*handlermocks.DON, []*shardEndpoint) {
	t.Helper()
	lggr := logger.Test(t)
	cfg := WithDefaults(ServiceConfig{
		MaxTriggerRequestDurationMs: 5000,
		NodeSendTimeoutMs:           5000,
	})

	shardsCfg := make([]config.Shard, len(shardNodeSets))
	connMgrs := make([]handlers.DON, len(shardNodeSets))
	mockDons := make([]*handlermocks.DON, len(shardNodeSets))
	for i, addrs := range shardNodeSets {
		members := make([]config.NodeConfig, len(addrs))
		for j, addr := range addrs {
			members[j] = config.NodeConfig{Address: addr}
		}
		shardsCfg[i] = config.Shard{Nodes: members}
		mockDons[i] = handlermocks.NewDON(t)
		connMgrs[i] = mockDons[i]
	}

	shardedDONs := []config.ShardedDONConfig{
		{DonName: donName, F: f, Shards: shardsCfg},
	}
	shards, nodeAddrToShard, err := buildShardEndpoints(shardedDONs, [][]handlers.DON{connMgrs})
	require.NoError(t, err)

	allMembersSlice := allMembers(shards)
	testMetrics, err := metrics.NewMetrics(allMembersSlice)
	require.NoError(t, err)
	metadataHandler := NewWorkflowMetadataHandler(lggr, cfg, shards, nodeAddrToShard, testMetrics)
	userRateLimiter := createTestUserRateLimiter()
	handler := NewHTTPTriggerHandler(lggr, cfg, shards, nodeAddrToShard, metadataHandler, userRateLimiter, testMetrics)
	return handler, mockDons, shards
}

// TestHttpTriggerHandler_MultiShardQuorumRace verifies the core sharding
// behavior: when a workflow is assigned to multiple shards, the FIRST shard to
// reach quorum produces the user response. Late responses from the other shard
// are recognized as already-processed (not errors) since the callback is kept
// until the periodic reaper removes it.
func TestHttpTriggerHandler_MultiShardQuorumRace(t *testing.T) {
	t.Parallel()

	// 2 shards, each with 3 nodes, F=1, threshold=(3+1)//2+1=3
	handler, mockDons, shards := createMultiShardTriggerHandler(t, "don",
		[][]string{{"n1", "n2", "n3"}, {"n4", "n5", "n6"}}, 1)

	privateKey := createTestPrivateKey(t)
	registerWorkflowOnShards(t, handler, workflowID, privateKey, shards[0], shards[1])

	triggerReq := createTestTriggerRequest(workflowID)
	reqBytes, err := json.Marshal(triggerReq)
	require.NoError(t, err)
	rawParams := json.RawMessage(reqBytes)
	req := &jsonrpc.Request[json.RawMessage]{
		Version: "2.0",
		ID:      "test-quorum-race",
		Method:  gateway_common.MethodWorkflowExecute,
		Params:  &rawParams,
	}
	req.Auth = createTestJWTToken(t, req, privateKey)
	callback := hc.NewCallback()

	// All 6 nodes receive the send successfully on the first attempt.
	for _, addr := range []string{"n1", "n2", "n3"} {
		mockDons[0].EXPECT().SendToNode(mock.Anything, addr, mock.Anything).Return(nil).Once()
	}
	for _, addr := range []string{"n4", "n5", "n6"} {
		mockDons[1].EXPECT().SendToNode(mock.Anything, addr, mock.Anything).Return(nil).Once()
	}

	require.NoError(t, handler.Start(t.Context()))
	t.Cleanup(func() { _ = handler.Close() })

	// HandleUserTriggerRequest blocks until all shard sends complete; since all
	// succeed on the first attempt, it returns immediately.
	err = handler.HandleUserTriggerRequest(t.Context(), req, callback, time.Now())
	require.NoError(t, err)

	// Shard 1 reaches quorum first (3 identical responses).
	rawRes := json.RawMessage(`{"result":"from-shard1"}`)
	nodeResp := &jsonrpc.Response[json.RawMessage]{
		Version: "2.0",
		ID:      req.ID,
		Result:  &rawRes,
	}
	for _, addr := range []string{"n4", "n5", "n6"} {
		require.NoError(t, handler.HandleNodeTriggerResponse(t.Context(), nodeResp, addr))
	}

	// User should have received the response from shard 1.
	payload, err := callback.Wait(t.Context())
	require.NoError(t, err)
	require.Equal(t, api.NoError, payload.ErrorCode)
	require.NotEmpty(t, payload.RawResponse)

	// Late response from shard 0 is not an error: the request was already
	// processed (shard 1 reached quorum first), so it is ignored with a debug log.
	err = handler.HandleNodeTriggerResponse(t.Context(), nodeResp, "n1")
	require.NoError(t, err)

	// Callback is kept and marked processed until the periodic reaper removes it.
	handler.callbacksMu.Lock()
	saved, exists := handler.callbacks[req.ID]
	handler.callbacksMu.Unlock()
	require.True(t, exists, "callback should be kept after processing so late responses are recognized")
	require.True(t, saved.processed, "callback should be marked as processed")
}

// TestHttpTriggerHandler_NodeResponseFromUnassignedShard verifies that a node
// response from a shard the workflow is NOT assigned to is rejected.
func TestHttpTriggerHandler_NodeResponseFromUnassignedShard(t *testing.T) {
	t.Parallel()

	// 3 shards; workflow registered on shard 0 and shard 1 only.
	handler, mockDons, shards := createMultiShardTriggerHandler(t, "don",
		[][]string{{"n1", "n2", "n3"}, {"n4", "n5", "n6"}, {"n7", "n8", "n9"}}, 1)

	privateKey := createTestPrivateKey(t)
	registerWorkflowOnShards(t, handler, workflowID, privateKey, shards[0], shards[1])

	triggerReq := createTestTriggerRequest(workflowID)
	reqBytes, err := json.Marshal(triggerReq)
	require.NoError(t, err)
	rawParams := json.RawMessage(reqBytes)
	req := &jsonrpc.Request[json.RawMessage]{
		Version: "2.0",
		ID:      "test-unassigned-shard",
		Method:  gateway_common.MethodWorkflowExecute,
		Params:  &rawParams,
	}
	req.Auth = createTestJWTToken(t, req, privateKey)
	callback := hc.NewCallback()

	for _, addr := range []string{"n1", "n2", "n3"} {
		mockDons[0].EXPECT().SendToNode(mock.Anything, addr, mock.Anything).Return(nil).Once()
	}
	for _, addr := range []string{"n4", "n5", "n6"} {
		mockDons[1].EXPECT().SendToNode(mock.Anything, addr, mock.Anything).Return(nil).Once()
	}
	// mockDons[2] (shard 2) gets NO expectations: any send to it would be a bug.

	require.NoError(t, handler.Start(t.Context()))
	t.Cleanup(func() { _ = handler.Close() })

	err = handler.HandleUserTriggerRequest(t.Context(), req, callback, time.Now())
	require.NoError(t, err)

	// Response from a node in the unassigned shard 2 must be rejected.
	rawRes := json.RawMessage(`{"result":"success"}`)
	nodeResp := &jsonrpc.Response[json.RawMessage]{
		Version: "2.0",
		ID:      req.ID,
		Result:  &rawRes,
	}
	err = handler.HandleNodeTriggerResponse(t.Context(), nodeResp, "n7")
	require.Error(t, err)
	require.Contains(t, err.Error(), "not assigned to workflow")
}

// TestHttpTriggerHandler_NodeResponseFromUnknownNode verifies that a response
// from a node address that does not appear in any shard is rejected.
func TestHttpTriggerHandler_NodeResponseFromUnknownNode(t *testing.T) {
	t.Parallel()

	handler, mockDons, shards := createMultiShardTriggerHandler(t, "don",
		[][]string{{"n1", "n2", "n3"}, {"n4", "n5", "n6"}}, 1)

	privateKey := createTestPrivateKey(t)
	registerWorkflowOnShards(t, handler, workflowID, privateKey, shards[0], shards[1])

	triggerReq := createTestTriggerRequest(workflowID)
	reqBytes, err := json.Marshal(triggerReq)
	require.NoError(t, err)
	rawParams := json.RawMessage(reqBytes)
	req := &jsonrpc.Request[json.RawMessage]{
		Version: "2.0",
		ID:      "test-unknown-node",
		Method:  gateway_common.MethodWorkflowExecute,
		Params:  &rawParams,
	}
	req.Auth = createTestJWTToken(t, req, privateKey)
	callback := hc.NewCallback()

	for _, addr := range []string{"n1", "n2", "n3"} {
		mockDons[0].EXPECT().SendToNode(mock.Anything, addr, mock.Anything).Return(nil).Once()
	}
	for _, addr := range []string{"n4", "n5", "n6"} {
		mockDons[1].EXPECT().SendToNode(mock.Anything, addr, mock.Anything).Return(nil).Once()
	}

	require.NoError(t, handler.Start(t.Context()))
	t.Cleanup(func() { _ = handler.Close() })

	err = handler.HandleUserTriggerRequest(t.Context(), req, callback, time.Now())
	require.NoError(t, err)

	rawRes := json.RawMessage(`{"result":"success"}`)
	nodeResp := &jsonrpc.Response[json.RawMessage]{
		Version: "2.0",
		ID:      req.ID,
		Result:  &rawRes,
	}
	err = handler.HandleNodeTriggerResponse(t.Context(), nodeResp, "totally-unknown-node")
	require.Error(t, err)
	require.Contains(t, err.Error(), "unknown node")
}

// TestHttpTriggerHandler_MultiShardSendFailureResilience verifies that when the
// request fan-out reaches multiple shards, a total failure in one shard does not
// prevent the user from receiving a response. The healthy shard reaches quorum,
// which closes doneCh and stops the failing shard's retry loop.
func TestHttpTriggerHandler_MultiShardSendFailureResilience(t *testing.T) {
	t.Parallel()

	lggr := logger.Test(t)
	cfg := WithDefaults(ServiceConfig{
		MaxTriggerRequestDurationMs: 5000,
		NodeSendTimeoutMs:           5000,
		RetryConfig: RetryConfig{
			InitialIntervalMs: 10,
			MaxIntervalTimeMs: 50,
			Multiplier:        2,
		},
	})

	// 2 shards, each with 3 nodes, F=1, threshold=3.
	shardedDONs := []config.ShardedDONConfig{
		{DonName: "don", F: 1, Shards: []config.Shard{
			{Nodes: []config.NodeConfig{{Address: "n1"}, {Address: "n2"}, {Address: "n3"}}},
			{Nodes: []config.NodeConfig{{Address: "n4"}, {Address: "n5"}, {Address: "n6"}}},
		}},
	}
	shard0Don := handlermocks.NewDON(t)
	shard1Don := handlermocks.NewDON(t)
	connMgrs := [][]handlers.DON{{shard0Don, shard1Don}}
	shards, nodeAddrToShard, err := buildShardEndpoints(shardedDONs, connMgrs)
	require.NoError(t, err)

	allMembersSlice := allMembers(shards)
	testMetrics, err := metrics.NewMetrics(allMembersSlice)
	require.NoError(t, err)
	metadataHandler := NewWorkflowMetadataHandler(lggr, cfg, shards, nodeAddrToShard, testMetrics)
	userRateLimiter := createTestUserRateLimiter()
	handler := NewHTTPTriggerHandler(lggr, cfg, shards, nodeAddrToShard, metadataHandler, userRateLimiter, testMetrics)

	privateKey := createTestPrivateKey(t)
	registerWorkflowOnShards(t, handler, workflowID, privateKey, shards[0], shards[1])

	triggerReq := createTestTriggerRequest(workflowID)
	reqBytes, err := json.Marshal(triggerReq)
	require.NoError(t, err)
	rawParams := json.RawMessage(reqBytes)
	req := &jsonrpc.Request[json.RawMessage]{
		Version: "2.0",
		ID:      "test-shard-failure-resilience",
		Method:  gateway_common.MethodWorkflowExecute,
		Params:  &rawParams,
	}
	req.Auth = createTestJWTToken(t, req, privateKey)
	callback := hc.NewCallback()

	// Shard 0's sends always fail (will keep retrying until doneCh closes).
	shard0Don.EXPECT().SendToNode(mock.Anything, mock.Anything, mock.Anything).Return(errors.New("connection error")).Maybe()

	// Shard 1's sends succeed; signal when all 3 have completed.
	shard1Sent := make(chan struct{})
	var shard1SendCount atomic.Int32
	for _, addr := range []string{"n4", "n5", "n6"} {
		shard1Don.EXPECT().SendToNode(mock.Anything, addr, mock.Anything).RunAndReturn(
			func(ctx context.Context, nodeAddress string, r *jsonrpc.Request[json.RawMessage]) error {
				if shard1SendCount.Add(1) == 3 {
					close(shard1Sent)
				}
				return nil
			}).Once()
	}

	require.NoError(t, handler.Start(t.Context()))
	t.Cleanup(func() { _ = handler.Close() })

	// HandleUserTriggerRequest blocks while shard 0 retries, so run it async.
	errCh := make(chan error, 1)
	go func() {
		errCh <- handler.HandleUserTriggerRequest(t.Context(), req, callback, time.Now())
	}()

	// Wait for shard 1's sends to complete.
	select {
	case <-shard1Sent:
	case <-t.Context().Done():
		t.Fatal("timed out waiting for shard 1 sends")
	}

	// Send responses from shard 1 to reach quorum.
	rawRes := json.RawMessage(`{"result":"from-shard1"}`)
	nodeResp := &jsonrpc.Response[json.RawMessage]{
		Version: "2.0",
		ID:      req.ID,
		Result:  &rawRes,
	}
	for _, addr := range []string{"n4", "n5", "n6"} {
		require.NoError(t, handler.HandleNodeTriggerResponse(t.Context(), nodeResp, addr))
	}

	// User should receive the response from the healthy shard.
	payload, err := callback.Wait(t.Context())
	require.NoError(t, err)
	require.Equal(t, api.NoError, payload.ErrorCode)
	require.NotEmpty(t, payload.RawResponse)

	// HandleUserTriggerRequest should return once doneCh closes (stopping shard 0).
	select {
	case err = <-errCh:
		require.NoError(t, err, "sendToShard returns nil when doneCh closes")
	case <-t.Context().Done():
		t.Fatal("timed out waiting for HandleUserTriggerRequest to return")
	}
}

// TestHttpTriggerHandler_WorkflowNotAssignedToAnyShard verifies that
// setupCallback returns an error when the workflow has no shard assignments.
func TestHttpTriggerHandler_WorkflowNotAssignedToAnyShard(t *testing.T) {
	t.Parallel()

	handler, _, _ := createMultiShardTriggerHandler(t, "don",
		[][]string{{"n1", "n2", "n3"}}, 1)

	privateKey := createTestPrivateKey(t)
	// Register the workflow's key and reference but assign it to NO shards.
	registerWorkflowOnShards(t, handler, workflowID, privateKey)

	triggerReq := createTestTriggerRequest(workflowID)
	reqBytes, err := json.Marshal(triggerReq)
	require.NoError(t, err)
	rawParams := json.RawMessage(reqBytes)
	req := &jsonrpc.Request[json.RawMessage]{
		Version: "2.0",
		ID:      "test-no-shards-assigned",
		Method:  gateway_common.MethodWorkflowExecute,
		Params:  &rawParams,
	}
	req.Auth = createTestJWTToken(t, req, privateKey)
	callback := hc.NewCallback()

	require.NoError(t, handler.Start(t.Context()))
	t.Cleanup(func() { _ = handler.Close() })

	err = handler.HandleUserTriggerRequest(t.Context(), req, callback, time.Now())
	require.Error(t, err)
	require.Contains(t, err.Error(), "not assigned to any shards")

	r, err := callback.Wait(t.Context())
	require.NoError(t, err)
	requireUserErrorSent(t, r, jsonrpc.ErrInternal)
}
