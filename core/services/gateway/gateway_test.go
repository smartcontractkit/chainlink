package gateway_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	handler_mocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/mocks"
	net_mocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/network/mocks"
)

func parseTOMLConfig(t *testing.T, tomlConfig string) *config.GatewayConfig {
	var cfg config.GatewayConfig
	err := toml.Unmarshal([]byte(tomlConfig), &cfg)
	require.NoError(t, err)
	return &cfg
}

func buildConfig(toAppend string) string {
	return `
	[userServerConfig]
	Path = "/user"
	[nodeServerConfig]
	Path = "/node"
	` + toAppend
}

func TestGateway_NewGatewayFromConfig_ValidConfig(t *testing.T) {
	t.Parallel()

	tomlConfig := buildConfig(`
[[dons]]
DonId = "my_don_1"
HandlerName = "dummy"

[[dons]]
DonId = "my_don_2"
HandlerName = "dummy"

[[dons.Members]]
Name = "node one"
Address = "0x0001020304050607080900010203040506070809"
`)

	lggr := logger.Test(t)
	_, err := gateway.NewGatewayFromConfig(parseTOMLConfig(t, tomlConfig), gateway.NewHandlerFactory(nil, nil, nil, lggr), lggr)
	require.NoError(t, err)
}

func TestGateway_NewGatewayFromConfig_DuplicateID(t *testing.T) {
	t.Parallel()

	tomlConfig := buildConfig(`
[[dons]]
DonId = "my_don"
HandlerName = "dummy"

[[dons]]
DonId = "my_don"
HandlerName = "dummy"
`)

	lggr := logger.Test(t)
	_, err := gateway.NewGatewayFromConfig(parseTOMLConfig(t, tomlConfig), gateway.NewHandlerFactory(nil, nil, nil, lggr), lggr)
	require.Error(t, err)
}

func TestGateway_NewGatewayFromConfig_InvalidHandler(t *testing.T) {
	t.Parallel()

	tomlConfig := buildConfig(`
[[dons]]
DonId = "my_don"
HandlerName = "no_such_handler"
`)

	lggr := logger.Test(t)
	_, err := gateway.NewGatewayFromConfig(parseTOMLConfig(t, tomlConfig), gateway.NewHandlerFactory(nil, nil, nil, lggr), lggr)
	require.Error(t, err)
}

func TestGateway_NewGatewayFromConfig_MissingID(t *testing.T) {
	t.Parallel()

	tomlConfig := buildConfig(`
[[dons]]
HandlerName = "dummy"
SomeOtherField = "abcd"
`)

	lggr := logger.Test(t)
	_, err := gateway.NewGatewayFromConfig(parseTOMLConfig(t, tomlConfig), gateway.NewHandlerFactory(nil, nil, nil, lggr), lggr)
	require.Error(t, err)
}

func TestGateway_NewGatewayFromConfig_InvalidNodeAddress(t *testing.T) {
	t.Parallel()

	tomlConfig := buildConfig(`
[[dons]]
HandlerName = "dummy"
DonId = "my_don"

[[dons.Members]]
Name = "node one"
Address = "0xnot_an_address"
`)

	lggr := logger.Test(t)
	_, err := gateway.NewGatewayFromConfig(parseTOMLConfig(t, tomlConfig), gateway.NewHandlerFactory(nil, nil, nil, lggr), lggr)
	require.Error(t, err)
}

func TestGateway_CleanStartAndClose(t *testing.T) {
	t.Parallel()

	lggr := logger.Test(t)
	gateway, err := gateway.NewGatewayFromConfig(parseTOMLConfig(t, buildConfig("")), gateway.NewHandlerFactory(nil, nil, nil, lggr), lggr)
	require.NoError(t, err)
	servicetest.Run(t, gateway)
}

func requireJSONRPCResult(t *testing.T, response []byte, expectedID string, expectedResult string) {
	require.JSONEq(t, fmt.Sprintf(`{"jsonrpc":"2.0","id":"%s","result":%s}`, expectedID, expectedResult), string(response))
}

func requireJSONRPCError(t *testing.T, responseBytes []byte, expectedID string, expectedCode int64, expectedMsg string) {
	var response jsonrpc.Response[json.RawMessage]
	err := json.Unmarshal(responseBytes, &response)
	require.NoError(t, err)
	require.Equal(t, jsonrpc.JsonRpcVersion, response.Version)
	require.Equal(t, expectedID, response.ID)
	require.Equal(t, expectedCode, response.Error.Code)
	require.Equal(t, expectedMsg, response.Error.Message)
	require.Nil(t, response.Error.Data)
}

func newGatewayWithMockHandler(t *testing.T) (gateway.Gateway, *handler_mocks.Handler) {
	httpServer := net_mocks.NewHttpServer(t)
	httpServer.On("SetHTTPRequestHandler", mock.Anything).Return(nil)
	handler := handler_mocks.NewHandler(t)
	handlers := map[string]handlers.Handler{
		"testDON": handler,
	}
	gw := gateway.NewGateway(&api.JsonRPCCodec{}, httpServer, handlers, nil, logger.Test(t))
	return gw, handler
}

// newSignedLegacyRequest creates a signed legacy request message for testing purposes.
// Legacy requests embed
func newSignedLegacyRequest(t *testing.T, messageID string, method string, donID string, payload []byte) []byte {
	msg := &api.Message{
		Body: api.MessageBody{
			MessageId: messageID,
			Method:    method,
			DonId:     donID,
			Payload:   payload,
		},
	}
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	require.NoError(t, msg.Sign(privateKey))
	codec := api.JsonRPCCodec{}
	rawRequest, err := codec.EncodeLegacyRequest(msg)
	require.NoError(t, err)
	return rawRequest
}

// newJSONRpcRequest creates a json rpc based request message for testing purposes.
func newJSONRpcRequest(t *testing.T, requestID string, method string, payload []byte) []byte {
	rawPayload := json.RawMessage(payload)
	request := jsonrpc.Request[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      requestID,
		Method:  method,
		Params:  &rawPayload,
	}
	rawRequest, err := json.Marshal(&request)
	require.NoError(t, err)
	return rawRequest
}

func TestGateway_ProcessRequest_ParseError(t *testing.T) {
	t.Parallel()

	gw, _ := newGatewayWithMockHandler(t)
	response, statusCode := gw.ProcessRequest(testutils.Context(t), []byte("{{}"), "")
	requireJSONRPCError(t, response, "", jsonrpc.ErrParse, "invalid character '{' looking for beginning of object key string")
	require.Equal(t, 400, statusCode)
}

func TestGateway_ProcessRequest_MessageValidationError(t *testing.T) {
	t.Parallel()

	gw, _ := newGatewayWithMockHandler(t)
	req := newSignedLegacyRequest(t, "abc", "request", api.NullChar, []byte{})
	response, statusCode := gw.ProcessRequest(testutils.Context(t), req, "")
	requireJSONRPCError(t, response, "abc", jsonrpc.ErrParse, "DON ID ending with null bytes")
	require.Equal(t, 400, statusCode)
}

func TestGateway_ProcessRequest_MissingDonId(t *testing.T) {
	t.Parallel()

	gw, _ := newGatewayWithMockHandler(t)
	req := newSignedLegacyRequest(t, "abc", "request", "", []byte{})
	response, statusCode := gw.ProcessRequest(testutils.Context(t), req, "")
	requireJSONRPCError(t, response, "abc", jsonrpc.ErrInvalidParams, "Unsupported DON ID or Handler: request")
	require.Equal(t, 400, statusCode)
}

func TestGateway_ProcessRequest_IncorrectDonId(t *testing.T) {
	t.Parallel()

	gw, _ := newGatewayWithMockHandler(t)
	req := newSignedLegacyRequest(t, "abc", "request", "unknownDON", []byte{})
	response, statusCode := gw.ProcessRequest(testutils.Context(t), req, "")
	requireJSONRPCError(t, response, "abc", jsonrpc.ErrInvalidParams, "Unsupported DON ID or Handler: unknownDON")
	require.Equal(t, 400, statusCode)
}

func TestGateway_LegacyRequest_HandlerResponse(t *testing.T) {
	t.Parallel()

	gw, handler := newGatewayWithMockHandler(t)
	handler.On("HandleLegacyUserMessage", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		msg := args.Get(1).(*api.Message)
		callbackCh := args.Get(2).(chan<- handlers.UserCallbackPayload)
		// echo back to sender with attached payload
		msg.Body.Payload = []byte(`{"result":"OK"}`)
		msg.Signature = ""
		codec := api.JsonRPCCodec{}
		callbackCh <- handlers.UserCallbackPayload{RawResponse: codec.EncodeLegacyResponse(msg), ErrorCode: api.NoError}
	})

	req := newSignedLegacyRequest(t, "abcd", "request", "testDON", []byte{})
	response, statusCode := gw.ProcessRequest(testutils.Context(t), req, "")
	requireJSONRPCResult(t, response, "abcd",
		`{"signature":"","body":{"message_id":"abcd","method":"request","don_id":"testDON","receiver":"","payload":{"result":"OK"}}}`)
	require.Equal(t, 200, statusCode)
}

func TestGateway_NewRequest_HandlerResponse(t *testing.T) {
	t.Parallel()

	gw, handler := newGatewayWithMockHandler(t)
	handler.On("HandleJSONRPCUserMessage", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		request := args.Get(1).(jsonrpc.Request[json.RawMessage])
		callbackCh := args.Get(2).(chan<- handlers.UserCallbackPayload)
		// echo back to sender with attached payload
		rawResult := json.RawMessage([]byte(`{"result":"OK"}`))
		response := jsonrpc.Response[json.RawMessage]{
			Version: jsonrpc.JsonRpcVersion,
			ID:      request.ID,
			Result:  &rawResult,
		}
		rawMsg, err := json.Marshal(&response)
		require.NoError(t, err)
		callbackCh <- handlers.UserCallbackPayload{RawResponse: rawMsg, ErrorCode: api.NoError}
	})

	req := newJSONRpcRequest(t, "abcd", "testDON", []byte(`{"type":"new"}`))
	response, statusCode := gw.ProcessRequest(testutils.Context(t), req, "")
	requireJSONRPCResult(t, response, "abcd", `{"result":"OK"}`)
	require.Equal(t, 200, statusCode)
}

func TestGateway_ProcessRequest_HandlerTimeout(t *testing.T) {
	t.Parallel()

	gw, handler := newGatewayWithMockHandler(t)
	handler.On("HandleLegacyUserMessage", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	timeoutCtx, cancel := context.WithTimeout(testutils.Context(t), time.Millisecond*10)
	defer cancel()

	req := newSignedLegacyRequest(t, "abcd", "request", "testDON", []byte{})
	response, statusCode := gw.ProcessRequest(timeoutCtx, req, "")
	requireJSONRPCError(t, response, "abcd", jsonrpc.ErrServerOverloaded, "handler timeout")
	require.Equal(t, 504, statusCode)
}

func TestGateway_ProcessRequest_HandlerError(t *testing.T) {
	t.Parallel()

	gw, handler := newGatewayWithMockHandler(t)
	handler.On("HandleLegacyUserMessage", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("failure"))

	req := newSignedLegacyRequest(t, "abcd", "request", "testDON", []byte{})
	response, statusCode := gw.ProcessRequest(testutils.Context(t), req, "")
	requireJSONRPCError(t, response, "abcd", jsonrpc.ErrInvalidRequest, "failure")
	require.Equal(t, 400, statusCode)
}
