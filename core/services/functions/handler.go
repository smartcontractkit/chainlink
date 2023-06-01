package functions

import (
	"crypto/ecdsa"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway"
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
	return gateway.SignData(h.signerKey, data...)
}

func (h *functionsConnectorHandler) HandleGatewayMessage(gatewayId string, msg *gateway.Message) {
	h.lggr.Debugw("functionsConnectorHandler: received message from gateway", "id", gatewayId)
}
