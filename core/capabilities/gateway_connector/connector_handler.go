package gateway_connector

import (
	"context"
	"crypto/ecdsa"

	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
)

type workflowConnectorHandler struct {
	services.StateMachine

	connector   connector.GatewayConnector
	signerKey   *ecdsa.PrivateKey
	nodeAddress string
	lggr        logger.Logger
}

type Response struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"error_message,omitempty"`
}

var (
	_ connector.Signer = &workflowConnectorHandler{}
	// _ connector.GatewayConnectorHandler = &workflowConnectorHandler{}
	// _ services.Service                  = &workflowConnectorHandler{}
)

// NOTE: the name is a little misleading - it's a generic handler to support all communication between Gateways and Workflow DONs
// We might want to come up with a cleaner split between capabilities.
func NewWorkflowConnectorHandler(config *config.GatewayConnector, signerKey *ecdsa.PrivateKey, lggr logger.Logger) (*workflowConnectorHandler, error) {
	return &workflowConnectorHandler{
		nodeAddress: (*config).NodeAddress(),
		signerKey:   signerKey,
		lggr:        lggr.Named("WorkflowConnectorHandler"),
	}, nil
}

func (h *workflowConnectorHandler) Sign(data ...[]byte) ([]byte, error) {
	return common.SignData(h.signerKey, data...)
}

func (h *workflowConnectorHandler) HandleGatewayMessage(ctx context.Context, gatewayId string, msg *api.Message) {
}
func (h *workflowConnectorHandler) Start(ctx context.Context) error {
	return h.StartOnce("WorkflowConnectorHandler", func() error {
		return nil
	})
}
func (h *workflowConnectorHandler) Close() error {
	return h.StopOnce("WorkflowConnectorHandler", func() (err error) {
		return nil
	})
}

func (h *workflowConnectorHandler) SetConnector(connector connector.GatewayConnector) {
	h.connector = connector
}
