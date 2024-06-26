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
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
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
	userRateLimiter, err := hc.NewRateLimiter(hc.RateLimiterConfig{GlobalRPS: 100.0, GlobalBurst: 100, PerSenderRPS: 100.0, PerSenderBurst: 100})
	require.NoError(t, err)
	nodeRateLimiter, err := hc.NewRateLimiter(hc.RateLimiterConfig{GlobalRPS: 100.0, GlobalBurst: 100, PerSenderRPS: 100.0, PerSenderBurst: 100})
	require.NoError(t, err)
	pendingRequestsCache := hc.NewRequestCache[functions.PendingRequest](requestTimeout, 1000)
	allowedHeartbeatInititors := map[string]struct{}{heartbeatSender: {}}
	handler := functions.NewFunctionsHandler(cfg, donConfig, don, pendingRequestsCache, allowlist, subscriptions, minBalance, userRateLimiter, nodeRateLimiter, allowedHeartbeatInititors, logger.TestLogger(t))
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

func sendNodeReponses(t *testing.T, handler handlers.Handler, userRequestMsg api.Message, nodes []gc.TestNode, responses []bool) {
	for id, resp := range responses {
		nodeResponseMsg := userRequestMsg
		nodeResponseMsg.Body.Receiver = userRequestMsg.Body.Sender
		if resp {
			nodeResponseMsg.Body.Payload = []byte(`{"success":true}`)
		} else {
			nodeResponseMsg.Body.Payload = []byte(`{"success":false}`)
		}
		require.NoError(t, nodeResponseMsg.Sign(nodes[id].PrivateKey))
		_ = handler.HandleNodeMessage(testutils.Context(t), &nodeResponseMsg, nodes[id].Address)
	}
}

func TestFunctionsHandler_Minimal(t *testing.T) {
	t.Parallel()

	handler, err := functions.NewFunctionsHandlerFromConfig(json.RawMessage("{}"), &config.DONConfig{}, nil, nil, nil, logger.TestLogger(t))
	require.NoError(t, err)

	// empty message should always error out
	msg := &api.Message{}
	err = handler.HandleUserMessage(testutils.Context(t), msg, nil)
	require.Error(t, err)
}

func TestFunctionsHandler_CleanStartAndClose(t *testing.T) {
	t.Parallel()

	handler, err := functions.NewFunctionsHandlerFromConfig(json.RawMessage("{}"), &config.DONConfig{}, nil, nil, nil, logger.TestLogger(t))
	require.NoError(t, err)

	servicetest.Run(t, handler)
}

func TestFunctionsHandler_HandleUserMessage_SecretsSet(t *testing.T) {
	t.Parallel()

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

			callbachCh := make(chan handlers.UserCallbackPayload)
			done := make(chan struct{})
			go func() {
				defer close(done)
				// wait on a response from Gateway to the user
				response := <-callbachCh
				require.Equal(t, api.NoError, response.ErrCode)
				require.Equal(t, userRequestMsg.Body.MessageId, response.Msg.Body.MessageId)
				var payload functions.CombinedResponse
				require.NoError(t, json.Unmarshal(response.Msg.Body.Payload, &payload))
				require.Equal(t, test.expectedGatewayResult, payload.Success)
				require.Equal(t, test.expectedNodeMessageCount, len(payload.NodeResponses))
			}()

			allowlist.On("Allow", common.HexToAddress(user.Address)).Return(true, nil)
			subscriptions.On("GetMaxUserBalance", common.HexToAddress(user.Address)).Return(big.NewInt(1000), nil)
			don.On("SendToNode", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			require.NoError(t, handler.HandleUserMessage(testutils.Context(t), &userRequestMsg, callbachCh))
			sendNodeReponses(t, handler, userRequestMsg, nodes, test.nodeResults)
			<-done
		})
	}
}

func TestFunctionsHandler_HandleUserMessage_Heartbeat(t *testing.T) {
	t.Parallel()

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

			callbachCh := make(chan handlers.UserCallbackPayload)
			done := make(chan struct{})
			go func() {
				defer close(done)
				// wait on a response from Gateway to the user
				response := <-callbachCh
				require.Equal(t, api.NoError, response.ErrCode)
				require.Equal(t, userRequestMsg.Body.MessageId, response.Msg.Body.MessageId)
				var payload functions.CombinedResponse
				require.NoError(t, json.Unmarshal(response.Msg.Body.Payload, &payload))
				require.Equal(t, test.expectedGatewayResult, payload.Success)
				require.Equal(t, test.expectedNodeMessageCount, len(payload.NodeResponses))
			}()

			allowlist.On("Allow", common.HexToAddress(user.Address)).Return(true, nil)
			don.On("SendToNode", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			require.NoError(t, handler.HandleUserMessage(testutils.Context(t), &userRequestMsg, callbachCh))
			sendNodeReponses(t, handler, userRequestMsg, nodes, test.nodeResults)
			<-done
		})
	}
}

func TestFunctionsHandler_HandleUserMessage_InvalidMethod(t *testing.T) {
	t.Parallel()

	nodes, user := gc.NewTestNodes(t, 4), gc.NewTestNodes(t, 1)[0]
	handler, _, allowlist, _ := newFunctionsHandlerForATestDON(t, nodes, time.Hour*24, user.Address)
	userRequestMsg := newSignedMessage(t, "1234", "secrets_reveal_all_please", "don_id", user.PrivateKey)

	allowlist.On("Allow", common.HexToAddress(user.Address)).Return(true, nil)
	err := handler.HandleUserMessage(testutils.Context(t), &userRequestMsg, make(chan handlers.UserCallbackPayload))
	require.Error(t, err)
}

func TestFunctionsHandler_HandleUserMessage_Timeout(t *testing.T) {
	t.Parallel()

	nodes, user := gc.NewTestNodes(t, 4), gc.NewTestNodes(t, 1)[0]
	handler, don, allowlist, subscriptions := newFunctionsHandlerForATestDON(t, nodes, time.Millisecond*10, user.Address)
	userRequestMsg := newSignedMessage(t, "1234", "secrets_set", "don_id", user.PrivateKey)

	callbachCh := make(chan handlers.UserCallbackPayload)
	done := make(chan struct{})
	go func() {
		defer close(done)
		// wait on a response from Gateway to the user
		response := <-callbachCh
		require.Equal(t, api.RequestTimeoutError, response.ErrCode)
		require.Equal(t, userRequestMsg.Body.MessageId, response.Msg.Body.MessageId)
	}()

	allowlist.On("Allow", common.HexToAddress(user.Address)).Return(true, nil)
	subscriptions.On("GetMaxUserBalance", common.HexToAddress(user.Address)).Return(big.NewInt(1000), nil)
	don.On("SendToNode", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	require.NoError(t, handler.HandleUserMessage(testutils.Context(t), &userRequestMsg, callbachCh))
	<-done
}
