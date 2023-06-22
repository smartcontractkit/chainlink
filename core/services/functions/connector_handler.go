package functions

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type functionsConnectorHandler struct {
	connector connector.GatewayConnector
	signerKey *ecdsa.PrivateKey
	signerID  string
	lggr      logger.Logger
}

var (
	_ connector.Signer                  = &functionsConnectorHandler{}
	_ connector.GatewayConnectorHandler = &functionsConnectorHandler{}
)

func NewFunctionsConnectorHandler(signerID string, signerKey *ecdsa.PrivateKey, lggr logger.Logger) *functionsConnectorHandler {
	return &functionsConnectorHandler{
		signerID:  signerID,
		signerKey: signerKey,
		lggr:      lggr.Named("functionsConnectorHandler"),
	}
}

func (h *functionsConnectorHandler) SetConnector(connector connector.GatewayConnector) {
	h.connector = connector
}

func (h *functionsConnectorHandler) Sign(data ...[]byte) ([]byte, error) {
	return common.SignData(h.signerKey, data...)
}

func (h *functionsConnectorHandler) HandleGatewayMessage(ctx context.Context, gatewayId string, msg *api.Message) {
	sender, err := msg.ValidateSignature()
	if err != nil {
		h.lggr.Errorw("failed to validate message signature", "id", gatewayId, "error", err)
		return
	}
	if utils.StringToHex(string(sender)) != msg.Body.Sender {
		h.lggr.Errorw("message signer does not match sender", "id", gatewayId, "sender", msg.Body.Sender)
		return
	}

	switch msg.Body.Method {
	case "ping":
		h.handlePingMessage(ctx, gatewayId, msg)
	default:
		h.lggr.Errorw("unsupported method", "id", gatewayId, "method", msg.Body.Method)
	}
}

func (h *functionsConnectorHandler) Start(ctx context.Context) error {
	return nil
}

func (h *functionsConnectorHandler) Close() error {
	return nil
}

func (h *functionsConnectorHandler) handlePingMessage(ctx context.Context, gatewayId string, msg *api.Message) {
	h.lggr.Debugw("ping message received from gateway", "id", gatewayId)

	type PongResponse struct {
		Success bool   `json:"success"`
		Error   string `json:"error,omitempty"`
	}

	response := PongResponse{Success: true}

	if err := h.signAndSendToGateway(ctx, gatewayId, msg.Body.MessageId, msg.Body.DonId, "pong", response); err != nil {
		h.lggr.Errorw("failed to send pong message to gateway", "id", gatewayId, "error", err)
	}
}

func (h *functionsConnectorHandler) signAndSendToGateway(ctx context.Context, gatewayId, messageId, donId, method string, payload any) error {
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	msg := &api.Message{
		Body: api.MessageBody{
			MessageId: messageId,
			DonId:     donId,
			Method:    method,
			Sender:    h.signerID,
			Payload:   payloadJson,
		},
	}
	if err = msg.Sign(h.signerKey); err != nil {
		return err
	}

	err = h.connector.SendToGateway(ctx, gatewayId, msg)
	if err == nil {
		h.lggr.Debugw("sent to gateway", "id", gatewayId, "messageId", messageId, "donId", donId, "method", method)
	}
	return err
}
