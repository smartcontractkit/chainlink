package gateway_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway"
)

type testConnManager struct {
	handler     gateway.Handler
	sendCounter int
}

func (m *testConnManager) SetHandler(handler gateway.Handler) {
	m.handler = handler
}

func (m *testConnManager) SendToNode(ctx context.Context, nodeAddress string, msg *gateway.Message) error {
	m.sendCounter++
	return nil
}

func TestDummyHandler_BasicFlow(t *testing.T) {
	t.Parallel()

	config := gateway.DONConfig{
		Members: []gateway.NodeConfig{
			{Name: "node one", Address: "addr_1"},
			{Name: "node two", Address: "addr_2"},
		},
	}

	connMgr := testConnManager{}
	handler, err := gateway.NewDummyHandler(&config, &connMgr)
	require.NoError(t, err)
	connMgr.SetHandler(handler)

	// User request
	msg := gateway.Message{Body: gateway.MessageBody{MessageId: "1234"}}
	callbackChan := make(chan gateway.UserCallbackPayload, 1)
	require.NoError(t, handler.HandleUserMessage(context.Background(), &msg, callbackChan))
	require.Equal(t, 2, connMgr.sendCounter)

	// Responses from both nodes
	require.NoError(t, handler.HandleNodeMessage(context.Background(), &msg, "addr_1"))
	require.NoError(t, handler.HandleNodeMessage(context.Background(), &msg, "addr_1"))
	response := <-callbackChan
	require.Equal(t, "1234", response.Msg.Body.MessageId)
}
