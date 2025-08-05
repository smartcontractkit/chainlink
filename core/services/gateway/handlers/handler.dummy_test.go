package handlers_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	hc "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
)

const (
	privateKey = "6c358b4f16344f03cfce12ebf7b768301bbe6a8977c98a2a2d76699f8bc56161"
)

type testConnManager struct {
	handler     handlers.Handler
	sendCounter int
}

func (m *testConnManager) SetHandler(handler handlers.Handler) {
	m.handler = handler
}

func (m *testConnManager) SendToNode(ctx context.Context, nodeAddress string, resp *jsonrpc.Request[json.RawMessage]) error {
	m.sendCounter++
	return nil
}

func TestDummyHandler_BasicFlow(t *testing.T) {
	t.Parallel()

	config := config.DONConfig{
		Members: []config.NodeConfig{
			{Name: "node one", Address: "addr_1"},
			{Name: "node two", Address: "addr_2"},
		},
	}

	connMgr := testConnManager{}
	handler, err := handlers.NewDummyHandler(&config, &connMgr, logger.Test(t))
	require.NoError(t, err)
	connMgr.SetHandler(handler)

	ctx := testutils.Context(t)

	// User request
	msg := api.Message{
		Body: api.MessageBody{
			MessageId: "1234",
			Method:    "testMethod",
			DonId:     "test_don",
		},
	}
	key, err := crypto.HexToECDSA(privateKey)
	require.NoError(t, err)
	err = msg.Sign(key)
	require.NoError(t, err)
	err = msg.Validate()
	require.NoError(t, err)
	callbackCh := make(chan handlers.UserCallbackPayload, 1)
	require.NoError(t, handler.HandleLegacyUserMessage(ctx, &msg, callbackCh))
	require.Equal(t, 2, connMgr.sendCounter)

	// Responses from both nodes
	resp, err := hc.ValidatedResponseFromMessage(&msg)
	require.NoError(t, err)
	require.NoError(t, handler.HandleNodeMessage(ctx, resp, msg.Body.Sender))
	require.NoError(t, handler.HandleNodeMessage(ctx, resp, msg.Body.Sender))
	response := <-callbackCh
	codec := api.JsonRPCCodec{}
	responseMsg, err := codec.DecodeLegacyResponse(response.RawResponse)
	require.NoError(t, err)
	require.Equal(t, "1234", responseMsg.Body.MessageId)
}
