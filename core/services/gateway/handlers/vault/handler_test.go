package vault

import (
	"context"
	"encoding/json"
	"sync"
	"testing"

	p2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/ratelimit"
	vaultcapmocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/mocks"
)

var NodeOne = config.NodeConfig{
	Name:    "node1",
	Address: "0x1234",
}

func setupHandler(t *testing.T) (handlers.Handler, chan handlers.UserCallbackPayload, *mocks.DON) {
	lggr := logger.Test(t)
	don := mocks.NewDON(t)
	donConfig := &config.DONConfig{
		DonId:   "test_don_id",
		Members: []config.NodeConfig{NodeOne},
	}
	handlerConfig := Config{
		RequestTimeoutSec: 30,
		NodeRateLimiter: ratelimit.RateLimiterConfig{
			GlobalRPS:      100,
			GlobalBurst:    100,
			PerSenderRPS:   10,
			PerSenderBurst: 10,
		},
	}
	methodConfig, err := json.Marshal(handlerConfig)
	require.NoError(t, err)

	requestAuthorizer := vaultcapmocks.NewRequestAuthorizer(t)
	requestAuthorizer.On("AuthorizeRequest", mock.Anything, mock.Anything).Return(true, owner, nil).Maybe()
	handler, err := NewHandler(methodConfig, donConfig, don, nil, requestAuthorizer, lggr)
	require.NoError(t, err)
	handler.aggregator = &mockAggregator{}
	return handler, make(chan handlers.UserCallbackPayload), don
}

type mockAggregator struct {
	err error
}

func (m *mockAggregator) Aggregate(_ context.Context, _ logger.Logger, _ *activeRequest, currResp *jsonrpc.Response[json.RawMessage]) (*jsonrpc.Response[json.RawMessage], error) {
	if m.err != nil {
		return nil, m.err
	}
	return currResp, nil
}

type mockCapabilitiesRegistry struct {
	F     uint8
	Nodes []capabilities.Node
}

var owner = "test_owner"

func (m *mockCapabilitiesRegistry) DONsForCapability(_ context.Context, _ string) ([]capabilities.DONWithNodes, error) {
	members := []p2ptypes.PeerID{}
	for _, n := range m.Nodes {
		members = append(members, *n.PeerID)
	}
	return []capabilities.DONWithNodes{
		{
			DON: capabilities.DON{
				F:       m.F,
				Members: members,
			},
			Nodes: m.Nodes,
		},
	}, nil
}

func TestVaultHandler_HandleJSONRPCUserMessage(t *testing.T) {
	createSecretsRequest := &vaultcommon.CreateSecretsRequest{
		RequestId: "test_request_id",
		EncryptedSecrets: []*vaultcommon.EncryptedSecret{
			{
				Id: &vaultcommon.SecretIdentifier{
					Key:   "test_id",
					Owner: owner,
				},
				EncryptedValue: "test_value",
			},
		},
	}
	params, err2 := json.Marshal(createSecretsRequest)
	require.NoError(t, err2)

	t.Run("happy path", func(t *testing.T) {
		var wg sync.WaitGroup
		h, callbackCh, don := setupHandler(t)
		don.On("SendToNode", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		requestID := "1"
		validJSONRequest := jsonrpc.Request[json.RawMessage]{
			ID:     requestID,
			Method: vaulttypes.MethodSecretsCreate,
			Params: (*json.RawMessage)(&params),
		}

		responseData := &vaultcommon.CreateSecretsResponse{
			Responses: []*vaultcommon.CreateSecretResponse{
				{
					Id:      createSecretsRequest.EncryptedSecrets[0].Id,
					Success: true,
				},
			},
		}
		resultBytes, err := json.Marshal(responseData)
		require.NoError(t, err)
		expectedRequestID := owner + "::" + requestID
		response := jsonrpc.Response[json.RawMessage]{
			ID:     expectedRequestID,
			Result: (*json.RawMessage)(&resultBytes),
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			callback := <-callbackCh
			var secretsResponse jsonrpc.Response[vaultcommon.CreateSecretsResponse]
			err2 := json.Unmarshal(callback.RawResponse, &secretsResponse)
			assert.NoError(t, err2)
			assert.Equal(t, validJSONRequest.ID, secretsResponse.ID, "Request ID should match")
			assert.Len(t, secretsResponse.Result.Responses, 1, "Should have one encrypted secret in response")
			assert.Equal(t, createSecretsRequest.EncryptedSecrets[0].Id.Key, secretsResponse.Result.Responses[0].Id.Key, "Secret ID should match")
			assert.True(t, secretsResponse.Result.Responses[0].Success, "Success should be true")
		}()

		err = h.HandleJSONRPCUserMessage(t.Context(), validJSONRequest, callbackCh)
		require.NoError(t, err)

		err = h.HandleNodeMessage(t.Context(), &response, NodeOne.Address)
		require.NoError(t, err)
		wg.Wait()
	})

	t.Run("happy path - delete secrets", func(t *testing.T) {
		var wg sync.WaitGroup
		h, callbackCh, don := setupHandler(t)
		don.On("SendToNode", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		id := &vaultcommon.SecretIdentifier{
			Key:       "foo",
			Namespace: "default",
			Owner:     owner,
		}
		reqData := &vaultcommon.DeleteSecretsRequest{
			RequestId: "id",
			Ids: []*vaultcommon.SecretIdentifier{
				id,
			},
		}
		reqDataBytes, err := json.Marshal(reqData)
		require.NoError(t, err)
		requestID := "1"
		validJSONRequest := jsonrpc.Request[json.RawMessage]{
			ID:     requestID,
			Method: vaulttypes.MethodSecretsDelete,
			Params: (*json.RawMessage)(&reqDataBytes),
		}

		responseData := &vaultcommon.DeleteSecretsResponse{
			Responses: []*vaultcommon.DeleteSecretResponse{
				{
					Id:      id,
					Success: true,
				},
			},
		}
		resultBytes, err := json.Marshal(responseData)
		require.NoError(t, err)
		expectedRequestID := owner + "::" + requestID
		response := jsonrpc.Response[json.RawMessage]{
			ID:     expectedRequestID,
			Result: (*json.RawMessage)(&resultBytes),
			Method: vaulttypes.MethodSecretsDelete,
		}
		resultBytes, err = json.Marshal(responseData)
		require.NoError(t, err)

		wg.Add(1)
		go func() {
			defer wg.Done()
			callback := <-callbackCh
			var secretsResponse jsonrpc.Response[vaultcommon.DeleteSecretsResponse]
			err2 := json.Unmarshal(callback.RawResponse, &secretsResponse)
			assert.NoError(t, err2)
			assert.Equal(t, validJSONRequest.ID, secretsResponse.ID, "Request ID should match")
			assert.True(t, proto.Equal(secretsResponse.Result, responseData), "Response data should match")
		}()

		err = h.HandleJSONRPCUserMessage(t.Context(), validJSONRequest, callbackCh)
		require.NoError(t, err)

		err = h.HandleNodeMessage(t.Context(), &response, NodeOne.Address)
		require.NoError(t, err)
		wg.Wait()
	})

	t.Run("happy path - list secret identifiers", func(t *testing.T) {
		var wg sync.WaitGroup
		h, callbackCh, don := setupHandler(t)
		don.On("SendToNode", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		requestID := "1"
		reqData := &vaultcommon.ListSecretIdentifiersRequest{
			RequestId: requestID,
			Owner:     owner,
		}
		reqDataBytes, err := json.Marshal(reqData)
		require.NoError(t, err)

		validJSONRequest := jsonrpc.Request[json.RawMessage]{
			ID:     requestID,
			Method: vaulttypes.MethodSecretsList,
			Params: (*json.RawMessage)(&reqDataBytes),
		}

		responseData := &vaultcommon.ListSecretIdentifiersResponse{
			Identifiers: []*vaultcommon.SecretIdentifier{
				{
					Key:       "foo",
					Owner:     owner,
					Namespace: "default",
				},
			},
		}
		resultBytes, err := json.Marshal(responseData)
		require.NoError(t, err)
		expectedRequestID := owner + "::" + requestID
		response := jsonrpc.Response[json.RawMessage]{
			ID:     expectedRequestID,
			Result: (*json.RawMessage)(&resultBytes),
			Method: vaulttypes.MethodSecretsList,
		}
		resultBytes, err = json.Marshal(responseData)
		require.NoError(t, err)

		wg.Add(1)
		go func() {
			defer wg.Done()
			callback := <-callbackCh
			var secretsResponse jsonrpc.Response[vaultcommon.ListSecretIdentifiersResponse]
			err2 := json.Unmarshal(callback.RawResponse, &secretsResponse)
			assert.NoError(t, err2)
			assert.Equal(t, validJSONRequest.ID, secretsResponse.ID, "Request ID should match")
			assert.True(t, proto.Equal(secretsResponse.Result, responseData), "Response data should match")
		}()

		err = h.HandleJSONRPCUserMessage(t.Context(), validJSONRequest, callbackCh)
		require.NoError(t, err)

		err = h.HandleNodeMessage(t.Context(), &response, NodeOne.Address)
		require.NoError(t, err)
		wg.Wait()
	})

	t.Run("unhappy path - quorum unobtainable", func(t *testing.T) {
		var wg sync.WaitGroup
		h, callbackCh, don := setupHandler(t)
		h.(*handler).aggregator = &mockAggregator{err: errQuorumUnobtainable}

		don.On("SendToNode", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		requestID := "1"
		reqData := &vaultcommon.ListSecretIdentifiersRequest{
			RequestId: requestID,
			Owner:     owner,
		}
		reqDataBytes, err := json.Marshal(reqData)
		require.NoError(t, err)

		validJSONRequest := jsonrpc.Request[json.RawMessage]{
			ID:     requestID,
			Method: vaulttypes.MethodSecretsList,
			Params: (*json.RawMessage)(&reqDataBytes),
		}

		expectedRequestID := owner + "::" + requestID
		response := jsonrpc.Response[json.RawMessage]{
			ID:     expectedRequestID,
			Method: vaulttypes.MethodSecretsList,
			Error: &jsonrpc.WireError{
				Code:    -32603,
				Message: "quorum unobtainable",
			},
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			callback := <-callbackCh
			var secretsResponse jsonrpc.Response[vaultcommon.ListSecretIdentifiersResponse]
			err2 := json.Unmarshal(callback.RawResponse, &secretsResponse)
			assert.NoError(t, err2)
			assert.Equal(t, validJSONRequest.ID, secretsResponse.ID, "Request ID should match")
			assert.Equal(t, response.Error, secretsResponse.Error, "Response error should match")
		}()

		err = h.HandleJSONRPCUserMessage(t.Context(), validJSONRequest, callbackCh)
		require.NoError(t, err)

		err = h.HandleNodeMessage(t.Context(), &response, NodeOne.Address)
		require.NoError(t, err)
		wg.Wait()
	})

	t.Run("unsupported method", func(t *testing.T) {
		var wg sync.WaitGroup
		h, callbackCh, don := setupHandler(t)
		// Don't expect SendToNode to be called for unsupported methods
		don.AssertNotCalled(t, "SendToNode")

		unsupportedMethodRequest := jsonrpc.Request[json.RawMessage]{
			ID:     "2",
			Method: "vault.unsupported.method",
			Params: (*json.RawMessage)(&params),
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			callback := <-callbackCh
			var secretsResponse jsonrpc.Response[vaultcommon.CreateSecretsResponse]
			err := json.Unmarshal(callback.RawResponse, &secretsResponse)
			assert.NoError(t, err)
			assert.Equal(t, unsupportedMethodRequest.ID, secretsResponse.ID, "Request ID should match")
			assert.Equal(t, "unsupported method: "+unsupportedMethodRequest.Method, secretsResponse.Error.Message, "Error message should match")
			assert.Equal(t, api.ToJSONRPCErrorCode(api.UnsupportedMethodError), secretsResponse.Error.Code, "Error code should match")
		}()

		err := h.HandleJSONRPCUserMessage(t.Context(), unsupportedMethodRequest, callbackCh)
		require.NoError(t, err)
		wg.Wait()
	})

	t.Run("empty params error", func(t *testing.T) {
		var wg sync.WaitGroup
		h, callbackCh, don := setupHandler(t)
		// Don't expect SendToNode to be called for parse errors
		don.AssertNotCalled(t, "SendToNode")

		emptyParamsRequest := jsonrpc.Request[json.RawMessage]{
			ID:     "3",
			Method: vaulttypes.MethodSecretsCreate,
			Params: &json.RawMessage{},
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			callback := <-callbackCh
			var secretsResponse jsonrpc.Response[vaultcommon.CreateSecretsResponse]
			err := json.Unmarshal(callback.RawResponse, &secretsResponse)
			assert.NoError(t, err)
			assert.Equal(t, emptyParamsRequest.ID, secretsResponse.ID, "Request ID should match")
			assert.Equal(t, "user message parse error: unexpected end of JSON input", secretsResponse.Error.Message, "Error message should match")
			assert.Equal(t, api.ToJSONRPCErrorCode(api.UserMessageParseError), secretsResponse.Error.Code, "Error code should match")
		}()

		err := h.HandleJSONRPCUserMessage(t.Context(), emptyParamsRequest, callbackCh)
		require.NoError(t, err)
		wg.Wait()
	})

	t.Run("no request inside the batch request", func(t *testing.T) {
		var wg sync.WaitGroup
		h, callbackCh, don := setupHandler(t)
		// Don't expect SendToNode to be called for invalid params
		don.AssertNotCalled(t, "SendToNode")

		invalidParams := json.RawMessage(`{"request_id": "empty_value_field"}`)
		invalidParamsRequest := jsonrpc.Request[json.RawMessage]{
			ID:     "4",
			Method: vaulttypes.MethodSecretsCreate,
			Params: &invalidParams,
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			callback := <-callbackCh
			var secretsResponse jsonrpc.Response[vaultcommon.CreateSecretsResponse]
			err := json.Unmarshal(callback.RawResponse, &secretsResponse)
			assert.NoError(t, err)
			assert.Equal(t, invalidParamsRequest.ID, secretsResponse.ID, "Request ID should match")
			assert.Equal(t, "invalid params error: must have at least 1 request", secretsResponse.Error.Message, "Error message should match")
			assert.Equal(t, api.ToJSONRPCErrorCode(api.InvalidParamsError), secretsResponse.Error.Code, "Error code should match")
		}()

		err := h.HandleJSONRPCUserMessage(t.Context(), invalidParamsRequest, callbackCh)
		require.NoError(t, err)
		wg.Wait()
	})

	t.Run("invalid params error", func(t *testing.T) {
		var wg sync.WaitGroup
		h, callbackCh, don := setupHandler(t)
		// Don't expect SendToNode to be called for invalid params
		don.AssertNotCalled(t, "SendToNode")

		invalidParamsRequest := &vaultcommon.CreateSecretsRequest{
			RequestId: "test_request_id",
			EncryptedSecrets: []*vaultcommon.EncryptedSecret{
				{
					Id: &vaultcommon.SecretIdentifier{
						Key:   "",
						Owner: "test_owner",
					},
					EncryptedValue: "test_value",
				},
			},
		}
		params, err2 := json.Marshal(invalidParamsRequest) //nolint:govet // The lock field is not set on this proto
		require.NoError(t, err2)
		jsonRequest := jsonrpc.Request[json.RawMessage]{
			ID:     "4",
			Method: vaulttypes.MethodSecretsCreate,
			Params: (*json.RawMessage)(&params),
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			callback := <-callbackCh
			var secretsResponse jsonrpc.Response[vaultcommon.CreateSecretsResponse]
			err := json.Unmarshal(callback.RawResponse, &secretsResponse)
			assert.NoError(t, err)
			assert.Equal(t, jsonRequest.ID, secretsResponse.ID, "Request ID should match")
			assert.Equal(t, "invalid params error: secret id key, owner and EncryptedValue cannot be empty on index 0", secretsResponse.Error.Message, "Error message should match")
			assert.Equal(t, api.ToJSONRPCErrorCode(api.InvalidParamsError), secretsResponse.Error.Code, "Error code should match")
		}()

		err := h.HandleJSONRPCUserMessage(t.Context(), jsonRequest, callbackCh)
		require.NoError(t, err)
		wg.Wait()
	})

	t.Run("stale node response", func(t *testing.T) {
		handler, callbackCh, _ := setupHandler(t)

		// Create a response for a request that was never sent or has already been processed
		responseData := &vaultcommon.CreateSecretsResponse{
			Responses: []*vaultcommon.CreateSecretResponse{
				{
					Id:      createSecretsRequest.EncryptedSecrets[0].Id,
					Success: true,
				},
			},
		}
		resultBytes, err := json.Marshal(responseData)
		require.NoError(t, err)
		staleResponse := jsonrpc.Response[json.RawMessage]{
			ID:     "stale_request_id",
			Result: (*json.RawMessage)(&resultBytes),
		}

		// Handle the stale node response - this should not trigger any callback
		// since there's no matching pending request
		err = handler.HandleNodeMessage(t.Context(), &staleResponse, NodeOne.Address)
		require.NoError(t, err)

		// Verify that no callback was sent by checking that the channel is empty
		select {
		case <-callbackCh:
			t.Error("Expected no callback for stale node response, but received one")
		default:
			// Expected: no callback should be sent for stale responses
		}
	})
}
