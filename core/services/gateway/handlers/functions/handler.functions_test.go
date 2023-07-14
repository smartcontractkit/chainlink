package functions_test

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	hc "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/functions"
	functions_mocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/functions/mocks"
	handlers_mocks "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/mocks"
)

type node struct {
	address    string
	privateKey *ecdsa.PrivateKey
}

func newNodes(t *testing.T, n int) []node {
	nodes := make([]node, n)
	for i := 0; i < n; i++ {
		privateKey, err := crypto.GenerateKey()
		require.NoError(t, err)
		address := strings.ToLower(crypto.PubkeyToAddress(privateKey.PublicKey).Hex())
		nodes[i] = node{address: address, privateKey: privateKey}
	}
	return nodes
}

func newFunctionsHandlerForATestDON(t *testing.T, nodes []node, requestTimeout time.Duration) (handlers.Handler, *handlers_mocks.DON, *functions_mocks.OnchainAllowlist) {
	cfg := functions.FunctionsHandlerConfig{}
	donConfig := &config.DONConfig{
		Members: []config.NodeConfig{},
		F:       1,
	}

	for id, n := range nodes {
		donConfig.Members = append(donConfig.Members, config.NodeConfig{
			Name:    fmt.Sprintf("node_%d", id),
			Address: n.address,
		})
	}

	don := handlers_mocks.NewDON(t)
	allowlist := functions_mocks.NewOnchainAllowlist(t)
	userRateLimiter, err := hc.NewRateLimiter(hc.RateLimiterConfig{GlobalRPS: 100.0, GlobalBurst: 100, PerUserRPS: 100.0, PerUserBurst: 100})
	require.NoError(t, err)
	nodeRateLimiter, err := hc.NewRateLimiter(hc.RateLimiterConfig{GlobalRPS: 100.0, GlobalBurst: 100, PerUserRPS: 100.0, PerUserBurst: 100})
	require.NoError(t, err)
	pendingRequestsCache := hc.NewRequestCache[functions.PendingSecretsRequest](requestTimeout)
	handler := functions.NewFunctionsHandler(cfg, donConfig, don, pendingRequestsCache, allowlist, userRateLimiter, nodeRateLimiter, logger.TestLogger(t))
	return handler, don, allowlist
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

func sendNodeReponses(t *testing.T, handler handlers.Handler, userRequestMsg api.Message, nodes []node, responses []bool) {
	for id, resp := range responses {
		nodeResponseMsg := userRequestMsg
		nodeResponseMsg.Body.Receiver = userRequestMsg.Body.Sender
		if resp {
			nodeResponseMsg.Body.Payload = []byte(`{"success":true}`)
		} else {
			nodeResponseMsg.Body.Payload = []byte(`{"success":false}`)
		}
		require.NoError(t, nodeResponseMsg.Sign(nodes[id].privateKey))
		_ = handler.HandleNodeMessage(context.Background(), &nodeResponseMsg, nodes[id].address)
	}
}

func TestFunctionsHandler_Minimal(t *testing.T) {
	t.Parallel()

	handler, err := functions.NewFunctionsHandlerFromConfig(json.RawMessage("{}"), &config.DONConfig{}, nil, nil, logger.TestLogger(t))
	require.NoError(t, err)

	// empty message should always error out
	msg := &api.Message{}
	err = handler.HandleUserMessage(testutils.Context(t), msg, nil)
	require.Error(t, err)
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
			nodes, user := newNodes(t, 4), newNodes(t, 1)[0]
			handler, don, allowlist := newFunctionsHandlerForATestDON(t, nodes, time.Hour*24)
			userRequestMsg := newSignedMessage(t, "1234", "secrets_set", "don_id", user.privateKey)

			callbachCh := make(chan handlers.UserCallbackPayload)
			done := make(chan struct{})
			go func() {
				defer close(done)
				// wait on a response from Gateway to the user
				response := <-callbachCh
				require.Equal(t, api.NoError, response.ErrCode)
				require.Equal(t, userRequestMsg.Body.MessageId, response.Msg.Body.MessageId)
				var payload functions.CombinedSecretsResponse
				require.NoError(t, json.Unmarshal(response.Msg.Body.Payload, &payload))
				require.Equal(t, test.expectedGatewayResult, payload.Success)
				require.Equal(t, test.expectedNodeMessageCount, len(payload.NodeResponses))
			}()

			allowlist.On("Allow", common.HexToAddress(user.address)).Return(true, nil)
			don.On("SendToNode", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			require.NoError(t, handler.HandleUserMessage(context.Background(), &userRequestMsg, callbachCh))
			sendNodeReponses(t, handler, userRequestMsg, nodes, test.nodeResults)
			<-done
		})
	}
}

func TestFunctionsHandler_HandleUserMessage_InvalidMethod(t *testing.T) {
	t.Parallel()

	nodes, user := newNodes(t, 4), newNodes(t, 1)[0]
	handler, _, allowlist := newFunctionsHandlerForATestDON(t, nodes, time.Hour*24)
	userRequestMsg := newSignedMessage(t, "1234", "secrets_reveal_all_please", "don_id", user.privateKey)

	allowlist.On("Allow", common.HexToAddress(user.address)).Return(true, nil)
	err := handler.HandleUserMessage(context.Background(), &userRequestMsg, make(chan handlers.UserCallbackPayload))
	require.Error(t, err)
}

func TestFunctionsHandler_HandleUserMessage_Timeout(t *testing.T) {
	t.Parallel()

	nodes, user := newNodes(t, 4), newNodes(t, 1)[0]
	handler, don, allowlist := newFunctionsHandlerForATestDON(t, nodes, time.Millisecond*10)
	userRequestMsg := newSignedMessage(t, "1234", "secrets_set", "don_id", user.privateKey)

	callbachCh := make(chan handlers.UserCallbackPayload)
	done := make(chan struct{})
	go func() {
		defer close(done)
		// wait on a response from Gateway to the user
		response := <-callbachCh
		require.Equal(t, api.RequestTimeoutError, response.ErrCode)
		require.Equal(t, userRequestMsg.Body.MessageId, response.Msg.Body.MessageId)
	}()

	allowlist.On("Allow", common.HexToAddress(user.address)).Return(true, nil)
	don.On("SendToNode", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	require.NoError(t, handler.HandleUserMessage(context.Background(), &userRequestMsg, callbachCh))
	<-done
}
