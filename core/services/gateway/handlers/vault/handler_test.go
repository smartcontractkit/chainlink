package vault

import (
	"encoding/json"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/ratelimit"

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
		NodeRateLimiterConfig: ratelimit.RateLimiterConfig{
			GlobalRPS:      100,
			GlobalBurst:    100,
			PerSenderRPS:   10,
			PerSenderBurst: 10,
		},
	}
	methodConfig, err := json.Marshal(handlerConfig)
	require.NoError(t, err)

	handler, err := NewHandler(methodConfig, donConfig, don, lggr)
	require.NoError(t, err)

	return handler, make(chan handlers.UserCallbackPayload), don
}

func TestVaultHandler_HandleJSONRPCUserMessage(t *testing.T) {
	createSecretsRequest := SecretsCreateRequest{
		ID:    "test_id",
		Value: "test_value",
	}
	params, err2 := json.Marshal(createSecretsRequest)
	require.NoError(t, err2)

	t.Run("happy path", func(t *testing.T) {
		var wg sync.WaitGroup
		handler, callbackCh, don := setupHandler(t)
		don.On("SendToNode", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		validJSONRequest := jsonrpc.Request[json.RawMessage]{
			ID:     "1",
			Method: MethodSecretsCreate,
			Params: (*json.RawMessage)(&params),
		}

		responseData := SecretsCreateResponse{
			ResponseBase: ResponseBase{
				Success: true,
			},
			SecretID: createSecretsRequest.ID,
		}
		resultBytes, err := json.Marshal(responseData)
		require.NoError(t, err)
		response := jsonrpc.Response[json.RawMessage]{
			ID:     "1",
			Result: (*json.RawMessage)(&resultBytes),
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			callback := <-callbackCh
			var secretsResponse jsonrpc.Response[SecretsCreateResponse]
			err2 := json.Unmarshal(callback.RawResponse, &secretsResponse)
			assert.NoError(t, err2)
			assert.Equal(t, validJSONRequest.ID, secretsResponse.ID, "Request ID should match")
			assert.Equal(t, createSecretsRequest.ID, secretsResponse.Result.SecretID, "Secret ID should match")
			assert.True(t, secretsResponse.Result.Success, "Success should be true")
		}()

		err = handler.HandleJSONRPCUserMessage(t.Context(), validJSONRequest, callbackCh)
		require.NoError(t, err)

		err = handler.HandleNodeMessage(t.Context(), &response, NodeOne.Address)
		require.NoError(t, err)
		wg.Wait()
	})

	t.Run("unsupported method", func(t *testing.T) {
		var wg sync.WaitGroup
		handler, callbackCh, don := setupHandler(t)
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
			var secretsResponse jsonrpc.Response[SecretsCreateResponse]
			err := json.Unmarshal(callback.RawResponse, &secretsResponse)
			assert.NoError(t, err)
			assert.Equal(t, unsupportedMethodRequest.ID, secretsResponse.ID, "Request ID should match")
			assert.Equal(t, "unsupported method: "+unsupportedMethodRequest.Method, secretsResponse.Error.Message, "Error message should match")
			assert.Equal(t, api.ToJSONRPCErrorCode(api.UnsupportedMethodError), secretsResponse.Error.Code, "Error code should match")
		}()

		err := handler.HandleJSONRPCUserMessage(t.Context(), unsupportedMethodRequest, callbackCh)
		require.NoError(t, err)
		wg.Wait()
	})

	t.Run("empty params error", func(t *testing.T) {
		var wg sync.WaitGroup
		handler, callbackCh, don := setupHandler(t)
		// Don't expect SendToNode to be called for parse errors
		don.AssertNotCalled(t, "SendToNode")

		emptyParamsRequest := jsonrpc.Request[json.RawMessage]{
			ID:     "3",
			Method: MethodSecretsCreate,
			Params: &json.RawMessage{},
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			callback := <-callbackCh
			var secretsResponse jsonrpc.Response[SecretsCreateResponse]
			err := json.Unmarshal(callback.RawResponse, &secretsResponse)
			assert.NoError(t, err)
			assert.Equal(t, emptyParamsRequest.ID, secretsResponse.ID, "Request ID should match")
			assert.Equal(t, "user message parse error: unexpected end of JSON input", secretsResponse.Error.Message, "Error message should match")
			assert.Equal(t, api.ToJSONRPCErrorCode(api.UserMessageParseError), secretsResponse.Error.Code, "Error code should match")
		}()

		err := handler.HandleJSONRPCUserMessage(t.Context(), emptyParamsRequest, callbackCh)
		require.NoError(t, err)
		wg.Wait()
	})

	t.Run("invalid params error", func(t *testing.T) {
		var wg sync.WaitGroup
		handler, callbackCh, don := setupHandler(t)
		// Don't expect SendToNode to be called for invalid params
		don.AssertNotCalled(t, "SendToNode")

		invalidParams := json.RawMessage(`{"id": "empty_value_field"}`)
		invalidParamsRequest := jsonrpc.Request[json.RawMessage]{
			ID:     "4",
			Method: MethodSecretsCreate,
			Params: &invalidParams,
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			callback := <-callbackCh
			var secretsResponse jsonrpc.Response[SecretsCreateResponse]
			err := json.Unmarshal(callback.RawResponse, &secretsResponse)
			assert.NoError(t, err)
			assert.Equal(t, invalidParamsRequest.ID, secretsResponse.ID, "Request ID should match")
			assert.Equal(t, "invalid params error: secret id and value cannot be empty", secretsResponse.Error.Message, "Error message should match")
			assert.Equal(t, api.ToJSONRPCErrorCode(api.InvalidParamsError), secretsResponse.Error.Code, "Error code should match")
		}()

		err := handler.HandleJSONRPCUserMessage(t.Context(), invalidParamsRequest, callbackCh)
		require.NoError(t, err)
		wg.Wait()
	})

	t.Run("stale node response", func(t *testing.T) {
		handler, callbackCh, _ := setupHandler(t)

		// Create a response for a request that was never sent or has already been processed
		responseData := SecretsCreateResponse{
			ResponseBase: ResponseBase{
				Success: true,
			},
			SecretID: "stale_secret_id",
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
