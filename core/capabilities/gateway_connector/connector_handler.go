package gateway_connector

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"

	ethCommon "github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/workflow"
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
	_ connector.Signer                  = &workflowConnectorHandler{}
	_ connector.GatewayConnectorHandler = &workflowConnectorHandler{}
	_ services.Service                  = &workflowConnectorHandler{}
)

// NOTE: the name is a little misleading - it's a generic handler to support all communication between Gateways and Workflow DONs
// We might want to come up with a cleaner split between capabilities.
func NewWorkflowConnectorHandler(config *WorkflowConnectorConfig, signerKey *ecdsa.PrivateKey, lggr logger.Logger) (*workflowConnectorHandler, error) {
	return &workflowConnectorHandler{
		nodeAddress: config.GatewayConnectorConfig.NodeAddress,
		signerKey:   signerKey,
		lggr:        lggr.Named("WorkflowConnectorHandler"),
	}, nil
}

func (h *workflowConnectorHandler) Sign(data ...[]byte) ([]byte, error) {
	return common.SignData(h.signerKey, data...)
}

func (h *workflowConnectorHandler) HandleGatewayMessage(ctx context.Context, gatewayId string, msg *api.Message) {
	body := &msg.Body
	fromAddr := ethCommon.HexToAddress(body.Sender)
	// TODO: apply allowlist and rate-limiting
	h.lggr.Debugw("handling gateway request", "id", gatewayId, "method", body.Method, "address", fromAddr)

	switch body.Method {
	case workflow.MethodAddWorkflow:
		// TODO: add a new workflow spec and return success/failure
		// we need access to Job ORM or whatever CLO uses to fully launch a new spec
		h.lggr.Debugw("added workflow spec", "payload", string(body.Payload))
		response := Response{Success: true}
		h.sendResponse(ctx, gatewayId, body, response)
	default:
		h.lggr.Errorw("unsupported method", "id", gatewayId, "method", body.Method)
	}
}

func (h *workflowConnectorHandler) sendResponse(ctx context.Context, gatewayId string, requestBody *api.MessageBody, payload any) error {
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	msg := &api.Message{
		Body: api.MessageBody{
			MessageId: requestBody.MessageId,
			DonId:     requestBody.DonId,
			Method:    requestBody.Method,
			Receiver:  requestBody.Sender,
			Payload:   payloadJson,
		},
	}
	if err = msg.Sign(h.signerKey); err != nil {
		return err
	}
	return h.connector.SendToGateway(ctx, gatewayId, msg)
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

func (h *workflowConnectorHandler) Ready() error {
	return nil
}

func (h *workflowConnectorHandler) HealthReport() map[string]error {
	return nil
}

func (h *workflowConnectorHandler) Name() string {
	return "WorkflowConnectorHandler"
}
