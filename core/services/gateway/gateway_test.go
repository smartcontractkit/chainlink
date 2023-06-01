package gateway_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway"
	gw_mocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/mocks"
	net_mocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/network/mocks"
)

func parseTOMLConfig(t *testing.T, tomlConfig string) *gateway.GatewayConfig {
	var cfg gateway.GatewayConfig
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
`)

	_, err := gateway.NewGatewayFromConfig(parseTOMLConfig(t, tomlConfig), logger.TestLogger(t))
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

	_, err := gateway.NewGatewayFromConfig(parseTOMLConfig(t, tomlConfig), logger.TestLogger(t))
	require.Error(t, err)
}

func TestGateway_NewGatewayFromConfig_InvalidHandler(t *testing.T) {
	t.Parallel()

	tomlConfig := buildConfig(`
[[dons]]
DonId = "my_don"
HandlerName = "no_such_handler"
`)

	_, err := gateway.NewGatewayFromConfig(parseTOMLConfig(t, tomlConfig), logger.TestLogger(t))
	require.Error(t, err)
}

func TestGateway_NewGatewayFromConfig_MissingID(t *testing.T) {
	t.Parallel()

	tomlConfig := buildConfig(`
[[dons]]
HandlerName = "dummy"
SomeOtherField = "abcd"
`)

	_, err := gateway.NewGatewayFromConfig(parseTOMLConfig(t, tomlConfig), logger.TestLogger(t))
	require.Error(t, err)
}

func requireJsonRPCResult(t *testing.T, response []byte, expectedId string, expectedResult string) {
	require.Equal(t, fmt.Sprintf(`{"jsonrpc":"2.0","id":"%s","result":%s}`, expectedId, expectedResult), string(response))
}

func requireJsonRPCError(t *testing.T, response []byte, expectedId string, expectedCode int, expectedMsg string) {
	require.Equal(t, fmt.Sprintf(`{"jsonrpc":"2.0","id":"%s","error":{"code":%d,"message":"%s"}}`, expectedId, expectedCode, expectedMsg), string(response))
}

func newGatewayWithMockHandler(t *testing.T) (gateway.Gateway, *gw_mocks.Handler) {
	httpServer := net_mocks.NewHttpServer(t)
	httpServer.On("SetHTTPRequestHandler", mock.Anything).Return(nil)
	handler := gw_mocks.NewHandler(t)
	handlers := map[string]gateway.Handler{
		"testDON": handler,
	}
	gw := gateway.NewGateway(&gateway.JsonRPCCodec{}, httpServer, handlers, nil, logger.TestLogger(t))
	return gw, handler
}

func TestGateway_ProcessRequest_ParseError(t *testing.T) {
	t.Parallel()

	gw, _ := newGatewayWithMockHandler(t)
	response, statusCode := gw.ProcessRequest(testutils.Context(t), []byte("{{}"))
	requireJsonRPCError(t, response, "", -32700, "invalid character '{' looking for beginning of object key string")
	require.Equal(t, 400, statusCode)
}

func TestGateway_ProcessRequest_IncorrectDonId(t *testing.T) {
	t.Parallel()

	gw, _ := newGatewayWithMockHandler(t)
	response, statusCode := gw.ProcessRequest(testutils.Context(t), []byte(`{"jsonrpc":"2.0", "id": "abc", "method": "request", "params": {}}`))
	requireJsonRPCError(t, response, "abc", -32602, "unsupported DON ID")
	require.Equal(t, 400, statusCode)

	response, statusCode = gw.ProcessRequest(testutils.Context(t), []byte(`{"jsonrpc":"2.0", "id": "abc", "method": "request", "params": {"body": {"don_id": "bad"}}}`))
	requireJsonRPCError(t, response, "abc", -32602, "unsupported DON ID")
	require.Equal(t, 400, statusCode)
}

func TestGateway_ProcessRequest_HandlerResponse(t *testing.T) {
	t.Parallel()

	gw, handler := newGatewayWithMockHandler(t)
	handler.On("HandleUserMessage", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		msg := args.Get(1).(*gateway.Message)
		callbackCh := args.Get(2).(chan<- gateway.UserCallbackPayload)
		// echo back to sender with attached payload
		msg.Body.Payload = []byte(`{"result":"OK"}`)
		callbackCh <- gateway.UserCallbackPayload{Msg: msg, ErrCode: gateway.NoError, ErrMsg: ""}
	})

	response, statusCode := gw.ProcessRequest(testutils.Context(t),
		[]byte(`{"jsonrpc":"2.0", "method": "request", "id": "abcd", "params": {"body":{"don_id": "testDON"}}}`))
	requireJsonRPCResult(t, response, "abcd",
		`{"signature":"","body":{"message_id":"abcd","method":"request","don_id":"testDON","sender":"","payload":{"result":"OK"}}}`)
	require.Equal(t, 200, statusCode)
}

func TestGateway_ProcessRequest_HandlerTimeout(t *testing.T) {
	t.Parallel()

	gw, handler := newGatewayWithMockHandler(t)
	handler.On("HandleUserMessage", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	timeoutCtx, cancel := context.WithTimeout(testutils.Context(t), time.Duration(time.Millisecond*10))
	defer cancel()

	response, statusCode := gw.ProcessRequest(timeoutCtx, []byte(`{"jsonrpc":"2.0", "method": "request", "id": "abcd", "params": {"body":{"don_id": "testDON"}}}`))
	requireJsonRPCError(t, response, "abcd", -32000, "handler timeout")
	require.Equal(t, 504, statusCode)
}

func TestGateway_ProcessRequest_HandlerError(t *testing.T) {
	t.Parallel()

	gw, handler := newGatewayWithMockHandler(t)
	handler.On("HandleUserMessage", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("failure"))

	response, statusCode := gw.ProcessRequest(testutils.Context(t), []byte(`{"jsonrpc":"2.0", "method": "request", "id": "abcd", "params": {"body":{"don_id": "testDON"}}}`))
	requireJsonRPCError(t, response, "abcd", -32000, "failure")
	require.Equal(t, 500, statusCode)
}
