package gateway_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway"
)

type TestHandler struct {
	t     *testing.T
	block chan bool
	resp  *gateway.Message
}

var _ gateway.Handler = &TestHandler{}

func (h *TestHandler) Init(connMgr gateway.ConnectionManager, donConfig *gateway.GatewayDONConfig) {
	//no-op
}

func (h *TestHandler) HandleUserMessage(msg *gateway.Message, cb gateway.Callback) {
	h.resp = msg
	h.block <- true
}

func (h *TestHandler) HandleNodeMessage(msg *gateway.Message, nodeAddr string) {
	//no-op
}

func TestConnector_connect(t *testing.T) {
	t.Parallel()

	gatewayConnectorConfig := &gateway.GatewayConnectorConfig{
		DONID:            "functions_local",
		GatewayAddresses: []string{"localhost:8040"},
	}
	closeChan := make(chan bool)
	handler := &TestHandler{t: t, block: closeChan}
	lggr := logger.TestLogger(t)
	gatewayConnector := gateway.NewGatewayConnector(gatewayConnectorConfig, handler, lggr)
	gatewayConnector.Start(context.TODO())
	<-closeChan

	require.Equal(t, nil, handler.resp)
}
