package functions

import (
	"context"
	"crypto/ecdsa"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
)

type functionsConnectorHandler struct {
	connector connector.GatewayConnector
	signerKey *ecdsa.PrivateKey
	lggr      logger.Logger
}

var (
	_ connector.Signer                  = &functionsConnectorHandler{}
	_ connector.GatewayConnectorHandler = &functionsConnectorHandler{}
)

func NewFunctionsConnectorHandler(signerKey *ecdsa.PrivateKey, lggr logger.Logger) *functionsConnectorHandler {
	return &functionsConnectorHandler{
		signerKey: signerKey,
		lggr:      lggr,
	}
}

func (h *functionsConnectorHandler) SetConnector(connector connector.GatewayConnector) {
	h.connector = connector
}

func (h *functionsConnectorHandler) Sign(data ...[]byte) ([]byte, error) {
	return common.SignData(h.signerKey, data...)
}

func (h *functionsConnectorHandler) HandleGatewayMessage(gatewayId string, msg *api.Message) {
	h.lggr.Debugw("functionsConnectorHandler: received message from gateway", "id", gatewayId)
}

func (h *functionsConnectorHandler) Start(ctx context.Context) error {
	return nil
}

func (h *functionsConnectorHandler) Close() error {
	return nil
}
