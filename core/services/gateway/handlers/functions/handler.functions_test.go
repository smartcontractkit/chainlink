package functions_test

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/assets"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/ratelimit"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	gc "github.com/smartcontractkit/chainlink/v2/core/services/gateway/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	hc "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/functions"
	allowlist_mocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/functions/allowlist/mocks"
	subscriptions_mocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/functions/subscriptions/mocks"
	handlers_mocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/mocks"
)

func newFunctionsHandlerForATestDON(t *testing.T, nodes []gc.TestNode, requestTimeout time.Duration, heartbeatSender string) (handlers.Handler, *handlers_mocks.DON, *allowlist_mocks.OnchainAllowlist, *subscriptions_mocks.OnchainSubscriptions) {
	cfg := functions.FunctionsHandlerConfig{}
	donConfig := &config.DONConfig{
		Members: []config.NodeConfig{},
		F:       1,
	}

	for id, n := range nodes {
		donConfig.Members = append(donConfig.Members, config.NodeConfig{
			Name:    fmt.Sprintf("node_%d", id),
			Address: n.Address,
		})
	}

	don := handlers_mocks.NewDON(t)
	allowlist := allowlist_mocks.NewOnchainAllowlist(t)
	subscriptions := subscriptions_mocks.NewOnchainSubscriptions(t)
	minBalance := assets.NewLinkFromJuels(100)
	userRateLimiter, err := ratelimit.NewRateLimiter(ratelimit.RateLimiterConfig{GlobalRPS: 100.0, GlobalBurst: 100, PerSenderRPS: 100.0, PerSenderBurst: 100})
	require.NoError(t, err)
	nodeRateLimiter, err := ratelimit.NewRateLimiter(ratelimit.RateLimiterConfig{GlobalRPS: 100.0, GlobalBurst: 100, PerSenderRPS: 100.0, PerSenderBurst: 100})
	require.NoError(t, err)
	pendingRequestsCache := hc.NewRequestCache[functions.PendingRequest](requestTimeout, 1000)
	allowedHeartbeatInititors := map[string]struct{}{heartbeatSender: {}}
	handler := functions.NewFunctionsHandler(cfg, donConfig, don, pendingRequestsCache, allowlist, subscriptions, minBalance, userRateLimiter, nodeRateLimiter, allowedHeartbeatInititors, logger.Test(t))
	return handler, don, allowlist, subscriptions
}

func newSignedMessage(t *testing.T, id string, method string, donId string, privateKey *ecdsa.PrivateKey) api.Message {
	msg := api.Message{
		Body: api.MessageBody{
			MessageId: id,
			Method:    method,
			DonId:     donId,
		},
	}
	require.NoError(t, msg.Sign(privateKey))
	return msg
}

func sendNodeReponse(t *testing.T, handler handlers.Handler, userRequestMsg api.Message, nodes []gc.TestNode, responses []bool) {
	for id, resp := range responses {
		nodeResponseMsg := userRequestMsg
		nodeResponseMsg.Body.Receiver = userRequestMsg.Body.Sender
		if resp {
			nodeResponseMsg.Body.Payload = []byte(`{"success":true}`)
		} else {
			nodeResponseMsg.Body.Payload = []byte(`{"success":false}`)
		}
		err := nodeResponseMsg.Sign(nodes[id].PrivateKey)
		var jsonResp *jsonrpc.Response[json.RawMessage]
		if err == nil {
			jsonResp, err = hc.ValidatedResponseFromMessage(&nodeResponseMsg) // ensure the message is valid
		}
		if err != nil {
			jsonResp = &jsonrpc.Response[json.RawMessage]{
				ID:     userRequestMsg.Body.MessageId,
				Result: nil,
				Error: &jsonrpc.WireError{
					Code:    jsonrpc.ErrInternal,
					Message: fmt.Sprintf("failed to prepare node response: %v", err),
				},
			}
		}
		_ = handler.HandleNodeMessage(testutils.Context(t), jsonResp, nodes[id].Address)
	}
}

func TestFunctionsHandler_Minimal(t *testing.T) {
	handler, err := functions.NewFunctionsHandlerFromConfig(json.RawMessage("{}"), &config.DONConfig{}, nil, nil, nil, logger.Test(t))
	require.NoError(t, err)

	// empty message should always error out
	msg := &api.Message{}
	err = handler.HandleLegacyUserMessage(testutils.Context(t), msg, nil)
	require.Error(t, err)
}

func TestFunctionsHandler_CleanStartAndClose(t *testing.T) {
	handler, err := functions.NewFunctionsHandlerFromConfig(json.RawMessage("{}"), &config.DONConfig{}, nil, nil, nil, logger.Test(t))
	require.NoError(t, err)

	servicetest.Run(t, handler)
}

func TestFunctionsHandler_HandleUserMessage_SecretsSet(t *testing.T) {
	tests := []struct {
		name                     string
		nodeResults              []bool
		expectedGatewayResult    bool
		expectedNodeMessageCount int
	}{
		{"three successful", []bool{true, true, true, false}, true, 2},
		{"two successful", []bool{false, true, false, true}, true, 2},
		{"one successful", []bool{false, true, false, false}, false, 3},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			nodes, user := gc.NewTestNodes(t, 4), gc.NewTestNodes(t, 1)[0]
			handler, don, allowlist, subscriptions := newFunctionsHandlerForATestDON(t, nodes, time.Hour*24, user.Address)
			userRequestMsg := newSignedMessage(t, "1234", "secrets_set", "don_id", user.PrivateKey)
			cb := hc.NewCallback()
			allowlist.On("Allow", common.HexToAddress(user.Address)).Return(true, nil)
			subscriptions.On("GetMaxUserBalance", common.HexToAddress(user.Address)).Return(big.NewInt(1000), nil)
			don.On("SendToNode", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			require.NoError(t, handler.HandleLegacyUserMessage(testutils.Context(t), &userRequestMsg, cb))

			done := make(chan struct{})
			go func() {
				defer close(done)
				// Ensure the response is sent on another thread to avoid deadlock
				sendNodeReponse(t, handler, userRequestMsg, nodes, test.nodeResults)
			}()

			// wait on a response from Gateway to the user
			response, err := cb.Wait(t.Context())
			require.NoError(t, err)
			// wait for goroutine to complete to avoid race condition
			<-done

			require.Equal(t, api.NoError, response.ErrorCode)
			codec := api.JsonRPCCodec{}
			msg, err := codec.DecodeLegacyResponse(response.RawResponse)
			require.NoError(t, err)
			require.Equal(t, userRequestMsg.Body.MessageId, msg.Body.MessageId)
			var payload functions.CombinedResponse
			require.NoError(t, json.Unmarshal(msg.Body.Payload, &payload))
			require.Equal(t, test.expectedGatewayResult, payload.Success)
			require.Len(t, payload.NodeResponses, test.expectedNodeMessageCount)
		})
	}
}

func TestFunctionsHandler_HandleUserMessage_Heartbeat(t *testing.T) {
	tests := []struct {
		name                     string
		nodeResults              []bool
		expectedGatewayResult    bool
		expectedNodeMessageCount int
	}{
		{"three successful", []bool{true, true, true, false}, true, 2},
		{"two successful", []bool{false, true, false, true}, true, 2},
		{"one successful", []bool{false, true, false, false}, true, 2},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			nodes, user := gc.NewTestNodes(t, 4), gc.NewTestNodes(t, 1)[0]
			handler, don, allowlist, _ := newFunctionsHandlerForATestDON(t, nodes, time.Hour*24, user.Address)
			userRequestMsg := newSignedMessage(t, "1234", "heartbeat", "don_id", user.PrivateKey)
			cb := hc.NewCallback()
			allowlist.On("Allow", common.HexToAddress(user.Address)).Return(true, nil)
			don.On("SendToNode", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			require.NoError(t, handler.HandleLegacyUserMessage(testutils.Context(t), &userRequestMsg, cb))

			done := make(chan struct{})
			go func() {
				defer close(done)
				// Ensure the response is sent on another thread to avoid deadlock
				sendNodeReponse(t, handler, userRequestMsg, nodes, test.nodeResults)
			}()

			// wait on a response from Gateway to the user
			response, err := cb.Wait(t.Context())
			require.NoError(t, err)
			// wait for goroutine to complete to avoid race condition
			<-done

			require.Equal(t, api.NoError, response.ErrorCode)
			codec := api.JsonRPCCodec{}
			msg, err := codec.DecodeLegacyResponse(response.RawResponse)
			require.NoError(t, err)
			require.Equal(t, userRequestMsg.Body.MessageId, msg.Body.MessageId)
			var payload functions.CombinedResponse
			require.NoError(t, json.Unmarshal(msg.Body.Payload, &payload))
			require.Equal(t, test.expectedGatewayResult, payload.Success)
			require.Len(t, payload.NodeResponses, test.expectedNodeMessageCount)
		})
	}
}

func TestFunctionsHandler_HandleUserMessage_InvalidMethod(t *testing.T) {
	nodes, user := gc.NewTestNodes(t, 4), gc.NewTestNodes(t, 1)[0]
	handler, _, allowlist, _ := newFunctionsHandlerForATestDON(t, nodes, time.Hour*24, user.Address)
	userRequestMsg := newSignedMessage(t, "1234", "secrets_reveal_all_please", "don_id", user.PrivateKey)

	allowlist.On("Allow", common.HexToAddress(user.Address)).Return(true, nil)
	cb := hc.NewCallback()
	err := handler.HandleLegacyUserMessage(testutils.Context(t), &userRequestMsg, cb)
	require.Error(t, err)
}

func TestFunctionsHandler_HandleUserMessage_Timeout(t *testing.T) {
	nodes, user := gc.NewTestNodes(t, 4), gc.NewTestNodes(t, 1)[0]
	handler, don, allowlist, subscriptions := newFunctionsHandlerForATestDON(t, nodes, time.Millisecond*10, user.Address)
	userRequestMsg := newSignedMessage(t, "1234", "secrets_set", "don_id", user.PrivateKey)
	cb := hc.NewCallback()
	allowlist.On("Allow", common.HexToAddress(user.Address)).Return(true, nil)
	subscriptions.On("GetMaxUserBalance", common.HexToAddress(user.Address)).Return(big.NewInt(1000), nil)
	don.On("SendToNode", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	require.NoError(t, handler.HandleLegacyUserMessage(testutils.Context(t), &userRequestMsg, cb))

	// wait on a response from Gateway to the user
	response, err := cb.Wait(t.Context())
	require.NoError(t, err)
	require.Equal(t, api.RequestTimeoutError, response.ErrorCode)
	codec := api.JsonRPCCodec{}
	msg, err := codec.DecodeLegacyResponse(response.RawResponse)
	require.NoError(t, err)
	require.Equal(t, userRequestMsg.Body.MessageId, msg.Body.MessageId)
}
