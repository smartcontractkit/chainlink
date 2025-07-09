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
	methodConfig := json.RawMessage(`{"request_timeout_sec": 30}`)

	return NewHandler(methodConfig, donConfig, don, lggr), make(chan handlers.UserCallbackPayload), don
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
}
