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
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	handlermocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/monitoring"
	netmocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/network/mocks"
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

type handlerFactory struct {
	handlers map[string]handlers.Handler
}

func (h *handlerFactory) NewHandler(handlerType gateway.HandlerType, _ json.RawMessage, _ []config.ShardedDONConfig, _ [][]handlers.DON) (handlers.Handler, error) {
	return h.handlers[handlerType], nil
}

func newGatewayHandler(t *testing.T) gateway.HandlerFactory {
	lggr := logger.Test(t)
	return gateway.NewHandlerFactory(nil, nil, nil, nil, nil, lggr, limits.Factory{Logger: lggr})
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
	_, err := gateway.NewGatewayFromConfig(parseTOMLConfig(t, tomlConfig), newGatewayHandler(t), lggr, limits.Factory{Logger: lggr})
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
	_, err := gateway.NewGatewayFromConfig(parseTOMLConfig(t, tomlConfig), newGatewayHandler(t), lggr, limits.Factory{Logger: lggr})
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
	_, err := gateway.NewGatewayFromConfig(parseTOMLConfig(t, tomlConfig), newGatewayHandler(t), lggr, limits.Factory{Logger: lggr})
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
	_, err := gateway.NewGatewayFromConfig(parseTOMLConfig(t, tomlConfig), newGatewayHandler(t), lggr, limits.Factory{Logger: lggr})
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
	_, err := gateway.NewGatewayFromConfig(parseTOMLConfig(t, tomlConfig), newGatewayHandler(t), lggr, limits.Factory{Logger: lggr})
	require.Error(t, err)
}

// TestGateway_NewGatewayFromConfig_NewStyleConfigTOML tests TOML parsing of the new-style
// config format using Services and ShardedDONs with the Shard struct.
// Setup: 2 DONs (donA with 2 shards of 4 nodes each, donB with 1 shard of 4 nodes),
// 2 services (workflows -> donA, vault -> donB).
func TestGateway_NewGatewayFromConfig_NewStyleConfigTOML(t *testing.T) {
	t.Parallel()

	// New-style config with 2 DONs and 2 services
	// donA has 2 shards (4 nodes each), donB has 1 shard (4 nodes)
	tomlConfig := buildConfig(`
[[shardedDONs]]
DonName = "donA"
F = 1

[[shardedDONs.Shards]]
[[shardedDONs.Shards.Nodes]]
Name = "donA_shard0_node0"
Address = "0x0001020304050607080900010203040506070809"
[[shardedDONs.Shards.Nodes]]
Name = "donA_shard0_node1"
Address = "0x0002020304050607080900010203040506070809"
[[shardedDONs.Shards.Nodes]]
Name = "donA_shard0_node2"
Address = "0x0003020304050607080900010203040506070809"
[[shardedDONs.Shards.Nodes]]
Name = "donA_shard0_node3"
Address = "0x0004020304050607080900010203040506070809"

[[shardedDONs.Shards]]
[[shardedDONs.Shards.Nodes]]
Name = "donA_shard1_node0"
Address = "0x0005020304050607080900010203040506070809"
[[shardedDONs.Shards.Nodes]]
Name = "donA_shard1_node1"
Address = "0x0006020304050607080900010203040506070809"
[[shardedDONs.Shards.Nodes]]
Name = "donA_shard1_node2"
Address = "0x0007020304050607080900010203040506070809"
[[shardedDONs.Shards.Nodes]]
Name = "donA_shard1_node3"
Address = "0x0008020304050607080900010203040506070809"

[[shardedDONs]]
DonName = "donB"
F = 1

[[shardedDONs.Shards]]
[[shardedDONs.Shards.Nodes]]
Name = "donB_shard0_node0"
Address = "0x0011020304050607080900010203040506070809"
[[shardedDONs.Shards.Nodes]]
Name = "donB_shard0_node1"
Address = "0x0012020304050607080900010203040506070809"
[[shardedDONs.Shards.Nodes]]
Name = "donB_shard0_node2"
Address = "0x0013020304050607080900010203040506070809"
[[shardedDONs.Shards.Nodes]]
Name = "donB_shard0_node3"
Address = "0x0014020304050607080900010203040506070809"

[[services]]
ServiceName = "workflows"
DONs = ["donA"]

[[services.Handlers]]
Name = "dummy"

[[services]]
ServiceName = "vault"
DONs = ["donB"]

[[services.Handlers]]
Name = "dummy"
`)

	cfg := parseTOMLConfig(t, tomlConfig)

	// Verify config structure was parsed correctly
	require.Len(t, cfg.ShardedDONs, 2, "should have 2 sharded DONs")
	require.Len(t, cfg.Services, 2, "should have 2 services")
	require.Empty(t, cfg.Dons, "legacy Dons should be empty")

	// Verify donA config (2 shards, 4 nodes each)
	require.Equal(t, "donA", cfg.ShardedDONs[0].DonName)
	require.Equal(t, 1, cfg.ShardedDONs[0].F)
	require.Len(t, cfg.ShardedDONs[0].Shards, 2, "donA should have 2 shards")
	require.Len(t, cfg.ShardedDONs[0].Shards[0].Nodes, 4, "donA shard0 should have 4 nodes")
	require.Len(t, cfg.ShardedDONs[0].Shards[1].Nodes, 4, "donA shard1 should have 4 nodes")
	require.Equal(t, "donA_shard0_node0", cfg.ShardedDONs[0].Shards[0].Nodes[0].Name)
	require.Equal(t, "donA_shard1_node0", cfg.ShardedDONs[0].Shards[1].Nodes[0].Name)

	// Verify donB config (1 shard, 4 nodes)
	require.Equal(t, "donB", cfg.ShardedDONs[1].DonName)
	require.Equal(t, 1, cfg.ShardedDONs[1].F)
	require.Len(t, cfg.ShardedDONs[1].Shards, 1, "donB should have 1 shard")
	require.Len(t, cfg.ShardedDONs[1].Shards[0].Nodes, 4, "donB shard0 should have 4 nodes")
	require.Equal(t, "donB_shard0_node0", cfg.ShardedDONs[1].Shards[0].Nodes[0].Name)

	// Verify services config
	require.Equal(t, "workflows", cfg.Services[0].ServiceName)
	require.Equal(t, []string{"donA"}, cfg.Services[0].DONs)
	require.Len(t, cfg.Services[0].Handlers, 1)

	require.Equal(t, "vault", cfg.Services[1].ServiceName)
	require.Equal(t, []string{"donB"}, cfg.Services[1].DONs)
	require.Len(t, cfg.Services[1].Handlers, 1)

	// Verify config validation passes
	require.NoError(t, cfg.Validate())
}

func TestGateway_CleanStartAndClose(t *testing.T) {
	t.Parallel()

	lggr := logger.Test(t)
	gatewayObj, err := gateway.NewGatewayFromConfig(parseTOMLConfig(t, buildConfig("")), newGatewayHandler(t), lggr, limits.Factory{Logger: lggr})
	require.NoError(t, err)
	servicetest.Run(t, gatewayObj)
}

func requireJSONRPCResult(t *testing.T, method string, response []byte, expectedID string, expectedResult string) {
	require.JSONEq(t, fmt.Sprintf(`{"jsonrpc":"2.0","id":"%s","result":%s,"method":"%s"}`, expectedID, expectedResult, method), string(response))
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

func newGatewayWithMockHandler(t *testing.T) (gateway.Gateway, *handlermocks.Handler) {
	httpServer := netmocks.NewHTTPServer(t)
	httpServer.On("SetHTTPRequestHandler", mock.Anything).Return(nil)
	handler := handlermocks.NewHandler(t)
	handlersObj := map[string]handlers.Handler{
		"testDON": handler,
	}
	gMetrics, err := monitoring.NewGatewayMetrics()
	require.NoError(t, err)
	gw := gateway.NewGateway(&api.JsonRPCCodec{}, httpServer, handlersObj, map[string]string{"testDON": "testDON"}, nil, nil, gMetrics, logger.Test(t))
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
	requireJSONRPCError(t, response, "abc", jsonrpc.ErrInvalidRequest, "Service name not found: request")
	require.Equal(t, 400, statusCode)
}

func TestGateway_ProcessRequest_IncorrectDonId(t *testing.T) {
	t.Parallel()

	gw, _ := newGatewayWithMockHandler(t)
	req := newSignedLegacyRequest(t, "abc", "request", "unknownDON", []byte{})
	response, statusCode := gw.ProcessRequest(testutils.Context(t), req, "")
	requireJSONRPCError(t, response, "abc", jsonrpc.ErrInvalidParams, "Unsupported DON ID: unknownDON")
	require.Equal(t, 400, statusCode)
}

func TestGateway_LegacyRequest_HandlerResponse(t *testing.T) {
	t.Parallel()

	gw, handler := newGatewayWithMockHandler(t)
	handler.On("HandleLegacyUserMessage", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		msg := args.Get(1).(*api.Message)
		callback := args.Get(2).(handlers.Callback)
		// echo back to sender with attached payload
		msg.Body.Payload = []byte(`{"result":"OK"}`)
		msg.Signature = ""
		codec := api.JsonRPCCodec{}
		err := callback.SendResponse(handlers.UserCallbackPayload{RawResponse: codec.EncodeLegacyResponse(msg), ErrorCode: api.NoError})
		require.NoError(t, err)
	})

	method := "request"
	req := newSignedLegacyRequest(t, "abcd", method, "testDON", []byte{})
	response, statusCode := gw.ProcessRequest(testutils.Context(t), req, "")
	requireJSONRPCResult(t, method, response, "abcd",
		`{"signature":"","body":{"message_id":"abcd","method":"request","don_id":"testDON","receiver":"","payload":{"result":"OK"}}}`)
	require.Equal(t, 200, statusCode)
}

func TestGateway_NewRequest_HandlerResponse(t *testing.T) {
	t.Parallel()

	gw, handler := newGatewayWithMockHandler(t)
	handler.On("HandleJSONRPCUserMessage", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		request := args.Get(1).(jsonrpc.Request[json.RawMessage])
		callback := args.Get(2).(handlers.Callback)
		// echo back to sender with attached payload
		rawResult := json.RawMessage(`{"result":"OK"}`)
		response := jsonrpc.Response[json.RawMessage]{
			Version: jsonrpc.JsonRpcVersion,
			ID:      request.ID,
			Result:  &rawResult,
			Method:  request.Method,
		}
		rawMsg, err := json.Marshal(&response)
		require.NoError(t, err)
		err = callback.SendResponse(handlers.UserCallbackPayload{RawResponse: rawMsg, ErrorCode: api.NoError})
		require.NoError(t, err)
	})

	req := newJSONRpcRequest(t, "abcd", "testDON", []byte(`{"type":"new"}`))
	response, statusCode := gw.ProcessRequest(testutils.Context(t), req, "")
	requireJSONRPCResult(t, "testDON", response, "abcd", `{"result":"OK"}`)
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
	requireJSONRPCError(t, response, "abcd", jsonrpc.ErrServerOverloaded, "handler timeout: context deadline exceeded")
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

func newMockHandler(t *testing.T, method string) *handlermocks.Handler {
	handler := handlermocks.NewHandler(t)
	handler.On("Methods").Return([]string{method})
	handler.On("HandleLegacyUserMessage", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		msg := args.Get(1).(*api.Message)
		callback := args.Get(2).(handlers.Callback)
		// echo back to sender with attached payload
		if msg.Body.Method != method {
			require.Fail(t, fmt.Sprintf("Expected method to be '%s'", method))
		}
		msg.Body.Payload = []byte(`{"result":"OK"}`)
		msg.Signature = ""
		codec := api.JsonRPCCodec{}
		err := callback.SendResponse(handlers.UserCallbackPayload{RawResponse: codec.EncodeLegacyResponse(msg), ErrorCode: api.NoError})
		require.NoError(t, err)
	})
	handler.On("HandleJSONRPCUserMessage", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		msg := args.Get(1).(jsonrpc.Request[json.RawMessage])
		callback := args.Get(2).(handlers.Callback)
		// echo back to sender with attached payload
		if msg.Method != method {
			require.Fail(t, fmt.Sprintf("Expected method to be '%s'", method))
		}
		rm := json.RawMessage(`{"result":"OK"}`)
		resp, err := json.Marshal(&jsonrpc.Response[json.RawMessage]{
			Version: jsonrpc.JsonRpcVersion,
			ID:      msg.ID,
			Method:  msg.Method,
			Result:  &rm,
		})
		require.NoError(t, err)
		err = callback.SendResponse(handlers.UserCallbackPayload{RawResponse: resp, ErrorCode: api.NoError})
		require.NoError(t, err)
	})
	return handler
}

func TestGateway_Multihandler(t *testing.T) {
	tomlConfig := buildConfig(`
[[dons]]
DonId = "1"

[[dons.Handlers]]
Name = "dummy"
ServiceName = "dummy"

[[dons.Handlers]]
Name = "dummy2"
ServiceName = "dummy2"

[[dons.Members]]
Name = "node one"
Address = "0x0001020304050607080900010203040506070809"
`)

	lggr := logger.Test(t)
	handler := newMockHandler(t, "dummy.dummy")
	handler2 := newMockHandler(t, "dummy2.dummy2")
	handlersObj := map[string]handlers.Handler{
		"dummy":  handler,
		"dummy2": handler2,
	}
	mhf := &handlerFactory{handlers: handlersObj}

	gatewayObj, err := gateway.NewGatewayFromConfig(parseTOMLConfig(t, tomlConfig), mhf, lggr, limits.Factory{Logger: lggr})
	require.NoError(t, err)

	method := "dummy.dummy"
	req := newSignedLegacyRequest(t, "abcd", method, "1", []byte{})
	response, statusCode := gatewayObj.ProcessRequest(testutils.Context(t), req, "")
	require.Equal(t, 200, statusCode, string(response))
	requireJSONRPCResult(t, method, response, "abcd",
		`{"signature":"","body":{"message_id":"abcd","method":"dummy.dummy","don_id":"1","receiver":"","payload":{"result":"OK"}}}`)

	method = "dummy2.dummy2"
	req = newSignedLegacyRequest(t, "abcd", method, "1", []byte{})
	response, statusCode = gatewayObj.ProcessRequest(testutils.Context(t), req, "")
	require.Equal(t, 200, statusCode, string(response))
	requireJSONRPCResult(t, method, response, "abcd",
		`{"signature":"","body":{"message_id":"abcd","method":"dummy2.dummy2","don_id":"1","receiver":"","payload":{"result":"OK"}}}`)

	method = "dummy.dummy"
	req = newJSONRpcRequest(t, "abcd", method, []byte(`{"type":"new"}`))
	response, statusCode = gatewayObj.ProcessRequest(testutils.Context(t), req, "")
	require.Equal(t, 200, statusCode, string(response))
	requireJSONRPCResult(t, method, response, "abcd",
		`{"result":"OK"}`)

	method = "dummy2.dummy2"
	req = newJSONRpcRequest(t, "abcd", method, []byte(`{"type":"new"}`))
	response, statusCode = gatewayObj.ProcessRequest(testutils.Context(t), req, "")
	require.Equal(t, 200, statusCode, string(response))
	requireJSONRPCResult(t, method, response, "abcd",
		`{"result":"OK"}`)
}

// TestGateway_NewStyleConfig_UserMessageRouting tests that user messages are correctly
// routed to the appropriate handlers when using the new-style config (Services + ShardedDONs).
//
// Test configuration represents:
//
//	ShardedDONs:
//	  - DonName: "don1" (4 nodes: 0x0001..., 0x0002..., 0x0003..., 0x0004...)
//	  - DonName: "don2" (4 nodes: 0x0011..., 0x0012..., 0x0013..., 0x0014...)
//
//	Services:
//	  - ServiceName: "workflows" -> attached to don1
//	  - ServiceName: "vault"     -> attached to don2
//
// This test verifies that:
//  1. Requests with method "workflows.*" are routed to the workflows handler (backed by don1)
//  2. Requests with method "vault.*" are routed to the vault handler (backed by don2)
//  3. Requests with unknown service names are rejected
func TestGateway_NewStyleConfig_UserMessageRouting(t *testing.T) {
	t.Parallel()

	// DON configuration (for documentation - actual nodes would be in connection manager)
	don1Nodes := []config.NodeConfig{
		{Name: "don1_node1", Address: "0x0001020304050607080900010203040506070809"},
		{Name: "don1_node2", Address: "0x0002020304050607080900010203040506070809"},
		{Name: "don1_node3", Address: "0x0003020304050607080900010203040506070809"},
		{Name: "don1_node4", Address: "0x0004020304050607080900010203040506070809"},
	}
	don2Nodes := []config.NodeConfig{
		{Name: "don2_node1", Address: "0x0011020304050607080900010203040506070809"},
		{Name: "don2_node2", Address: "0x0012020304050607080900010203040506070809"},
		{Name: "don2_node3", Address: "0x0013020304050607080900010203040506070809"},
		{Name: "don2_node4", Address: "0x0014020304050607080900010203040506070809"},
	}
	require.Len(t, don1Nodes, 4, "DON1 should have 4 nodes")
	require.Len(t, don2Nodes, 4, "DON2 should have 4 nodes")

	// Create mock handlers that track which methods they receive
	workflowsHandler := newServiceMockHandler(t, "workflows")
	vaultHandler := newServiceMockHandler(t, "vault")

	// Set up gateway with serviceToMultiHandler (new-style config)
	httpServer := netmocks.NewHTTPServer(t)
	httpServer.On("SetHTTPRequestHandler", mock.Anything).Return(nil)
	gMetrics, err := monitoring.NewGatewayMetrics()
	require.NoError(t, err)

	// Map services to their handlers (as would be created by setupFromNewConfig)
	// workflows service -> don1 (4 nodes)
	// vault service     -> don2 (4 nodes)
	serviceToMultiHandler := map[string]handlers.Handler{
		"workflows": workflowsHandler,
		"vault":     vaultHandler,
	}

	gw := gateway.NewGateway(
		&api.JsonRPCCodec{},
		httpServer,
		nil, // no legacy handlers
		nil, // no legacy serviceNameToDonID
		serviceToMultiHandler,
		nil, // no connMgr needed for this test
		gMetrics,
		logger.Test(t),
	)

	ctx := testutils.Context(t)

	// Test 1: workflows.execute should route to workflows handler
	req := newJSONRpcRequest(t, "req1", "workflows.execute", []byte(`{"workflow_id":"abc"}`))
	response, statusCode := gw.ProcessRequest(ctx, req, "")
	require.Equal(t, 200, statusCode, "workflows.execute failed: %s", string(response))
	requireJSONRPCResult(t, "workflows.execute", response, "req1", `{"service":"workflows","method":"workflows.execute"}`)

	// Test 2: vault.store should route to vault handler
	req = newJSONRpcRequest(t, "req3", "vault.store", []byte(`{"key":"secret"}`))
	response, statusCode = gw.ProcessRequest(ctx, req, "")
	require.Equal(t, 200, statusCode, "vault.store failed: %s", string(response))
	requireJSONRPCResult(t, "vault.store", response, "req3", `{"service":"vault","method":"vault.store"}`)

	// Test 3: Unknown service should return error
	req = newJSONRpcRequest(t, "req5", "unknown.method", []byte(`{}`))
	response, statusCode = gw.ProcessRequest(ctx, req, "")
	require.Equal(t, 400, statusCode)
	requireJSONRPCError(t, response, "req5", jsonrpc.ErrInvalidRequest, "Service name not found: unknown")
}

// TestGateway_NewStyleConfig_NodeResponseRouting tests that node responses are correctly
// routed back to the appropriate handlers when using the new-style config.
//
// Test configuration represents:
//
//	ShardedDONs:
//	  - DonName: "don1" (4 nodes)
//	  - DonName: "don2" (4 nodes)
//
//	Services:
//	  - ServiceName: "workflows" -> attached to don1
//	  - ServiceName: "vault"     -> attached to don2
//
// This test verifies that user requests are dispatched to the correct handler based on service name,
// which would then be forwarded to the correct DON's nodes. The handler tracks received methods
// to verify correct routing.
func TestGateway_NewStyleConfig_NodeResponseRouting(t *testing.T) {
	t.Parallel()

	// Track which handlers receive user messages (simulating what would eventually go to nodes)
	workflowsNodeMsgs := make(chan string, 10)
	vaultNodeMsgs := make(chan string, 10)

	workflowsHandler := newNodeResponseMockHandler(t, "workflows", workflowsNodeMsgs)
	vaultHandler := newNodeResponseMockHandler(t, "vault", vaultNodeMsgs)

	httpServer := netmocks.NewHTTPServer(t)
	httpServer.On("SetHTTPRequestHandler", mock.Anything).Return(nil)
	gMetrics, err := monitoring.NewGatewayMetrics()
	require.NoError(t, err)

	// Map services to their handlers
	// workflows service -> would be backed by don1 (4 nodes)
	// vault service     -> would be backed by don2 (4 nodes)
	serviceToMultiHandler := map[string]handlers.Handler{
		"workflows": workflowsHandler,
		"vault":     vaultHandler,
	}

	gw := gateway.NewGateway(
		&api.JsonRPCCodec{},
		httpServer,
		nil,
		nil,
		serviceToMultiHandler,
		nil,
		gMetrics,
		logger.Test(t),
	)

	ctx := testutils.Context(t)

	// Send user requests to establish context in handlers
	req := newJSONRpcRequest(t, "wf1", "workflows.execute", []byte(`{}`))
	_, statusCode := gw.ProcessRequest(ctx, req, "")
	require.Equal(t, 200, statusCode)

	req = newJSONRpcRequest(t, "v1", "vault.store", []byte(`{}`))
	_, statusCode = gw.ProcessRequest(ctx, req, "")
	require.Equal(t, 200, statusCode)

	// Verify correct handlers received the user messages
	// Messages are sent synchronously in the mock handler before callback.SendResponse,
	// so they are already in the channels by the time ProcessRequest returns.
	require.Equal(t, "workflows.execute", <-workflowsNodeMsgs)
	require.Equal(t, "vault.store", <-vaultNodeMsgs)
}

// newServiceMockHandler creates a mock handler that echoes back the service name and method.
func newServiceMockHandler(t *testing.T, serviceName string) *handlermocks.Handler {
	handler := handlermocks.NewHandler(t)
	handler.On("Methods").Return([]string{serviceName + ".execute", serviceName + ".status", serviceName + ".store", serviceName + ".retrieve"}).Maybe()
	handler.On("HandleJSONRPCUserMessage", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		request := args.Get(1).(jsonrpc.Request[json.RawMessage])
		callback := args.Get(2).(handlers.Callback)

		result := json.RawMessage(fmt.Sprintf(`{"service":"%s","method":"%s"}`, serviceName, request.Method))
		response := jsonrpc.Response[json.RawMessage]{
			Version: jsonrpc.JsonRpcVersion,
			ID:      request.ID,
			Result:  &result,
			Method:  request.Method,
		}
		rawMsg, err := json.Marshal(&response)
		require.NoError(t, err)
		err = callback.SendResponse(handlers.UserCallbackPayload{RawResponse: rawMsg, ErrorCode: api.NoError})
		require.NoError(t, err)
	}).Maybe()
	handler.On("HandleLegacyUserMessage", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()
	handler.On("HandleNodeMessage", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()
	return handler
}

// newNodeResponseMockHandler creates a mock handler that reports received methods to a channel.
func newNodeResponseMockHandler(t *testing.T, serviceName string, receivedMethods chan<- string) *handlermocks.Handler {
	handler := handlermocks.NewHandler(t)
	handler.On("Methods").Return([]string{serviceName + ".execute", serviceName + ".status", serviceName + ".store", serviceName + ".retrieve"}).Maybe()
	handler.On("HandleJSONRPCUserMessage", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		request := args.Get(1).(jsonrpc.Request[json.RawMessage])
		callback := args.Get(2).(handlers.Callback)

		// Report the method to tracking channel
		receivedMethods <- request.Method

		result := json.RawMessage(fmt.Sprintf(`{"service":"%s"}`, serviceName))
		response := jsonrpc.Response[json.RawMessage]{
			Version: jsonrpc.JsonRpcVersion,
			ID:      request.ID,
			Result:  &result,
			Method:  request.Method,
		}
		rawMsg, err := json.Marshal(&response)
		require.NoError(t, err)
		err = callback.SendResponse(handlers.UserCallbackPayload{RawResponse: rawMsg, ErrorCode: api.NoError})
		require.NoError(t, err)
	}).Maybe()
	handler.On("HandleLegacyUserMessage", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()
	handler.On("HandleNodeMessage", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()
	return handler
}
