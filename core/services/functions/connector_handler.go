package functions

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	"github.com/smartcontractkit/chainlink/v2/core/services/s4"
	"github.com/smartcontractkit/chainlink/v2/core/utils"

	ethCommon "github.com/ethereum/go-ethereum/common"
)

type functionsConnectorHandler struct {
	connector connector.GatewayConnector
	signerKey *ecdsa.PrivateKey
	signerID  string
	storage   s4.Storage
	lggr      logger.Logger
}

const (
	methodSecretsSet  = "secrets_set"
	methodSecretsList = "secrets_list"
)

var (
	_ connector.Signer                  = &functionsConnectorHandler{}
	_ connector.GatewayConnectorHandler = &functionsConnectorHandler{}
)

func NewFunctionsConnectorHandler(signerID string, signerKey *ecdsa.PrivateKey, storage s4.Storage, lggr logger.Logger) *functionsConnectorHandler {
	return &functionsConnectorHandler{
		signerID:  signerID,
		signerKey: signerKey,
		storage:   storage,
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
	// Gateway should have signature verified, therefore if it fails validation here, it is likely a maliscious gateway.
	signer, err := msg.ValidateSignature()
	if err != nil {
		h.lggr.Errorw("failed to validate message signature", "id", gatewayId, "error", err)
		return
	}
	if utils.StringToHex(string(signer)) != msg.Body.Sender {
		h.lggr.Errorw("signer address does not match sender", "id", gatewayId, "sender", msg.Body.Sender, "signer", utils.StringToHex(string(signer)))
		return
	}
	fromAddr := ethCommon.BytesToAddress(signer)

	h.lggr.Debugw("handling gateway request", "id", gatewayId, "method", msg.Body.Method)

	switch msg.Body.Method {
	case methodSecretsList:
		h.handleSecretsList(ctx, gatewayId, msg, fromAddr)
	case methodSecretsSet:
		h.handleSecretsSet(ctx, gatewayId, msg, fromAddr)
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

func (h *functionsConnectorHandler) handleSecretsList(ctx context.Context, gatewayId string, msg *api.Message, fromAddr ethCommon.Address) {
	type ListRow struct {
		SlotID     uint   `json:"slot_id"`
		Version    uint64 `json:"version"`
		Expiration int64  `json:"expiration"`
	}

	type ListResponse struct {
		Success bool      `json:"success"`
		Error   string    `json:"error,omitempty"`
		Rows    []ListRow `json:"rows,omitempty"`
	}

	var response ListResponse
	snapshot, err := h.storage.List(ctx, fromAddr)
	if err == nil {
		response.Success = true
		response.Rows = make([]ListRow, len(snapshot))
		for i, row := range snapshot {
			response.Rows[i] = ListRow{
				SlotID:     row.SlotId,
				Version:    row.Version,
				Expiration: row.Expiration,
			}
		}
	} else {
		response.Error = fmt.Sprintf("Failed to list secrets: %v", err)
	}

	if err := h.sendResponse(ctx, gatewayId, msg, response); err != nil {
		h.lggr.Errorw("failed to send response to gateway", "id", gatewayId, "error", err)
	}
}

func (h *functionsConnectorHandler) handleSecretsSet(ctx context.Context, gatewayId string, msg *api.Message, fromAddr ethCommon.Address) {
	type SetRequest struct {
		SlotID     uint   `json:"slot_id"`
		Version    uint64 `json:"version"`
		Expiration int64  `json:"expiration"`
		Payload    []byte `json:"payload"`
		Signature  []byte `json:"signature"`
	}

	type SetResponse struct {
		Success bool   `json:"success"`
		Error   string `json:"error,omitempty"`
	}

	var request SetRequest
	var response SetResponse
	err := json.Unmarshal(msg.Body.Payload, &request)
	if err == nil {
		key := s4.Key{
			Address: fromAddr,
			SlotId:  request.SlotID,
			Version: request.Version,
		}
		record := s4.Record{
			Expiration: request.Expiration,
			Payload:    request.Payload,
		}
		err = h.storage.Put(ctx, &key, &record, request.Signature)
		if err == nil {
			response.Success = true
		} else {
			response.Error = fmt.Sprintf("Failed to set secret: %v", err)
		}
	} else {
		response.Error = fmt.Sprintf("Bad request to set secret: %v", err)
	}

	if err := h.sendResponse(ctx, gatewayId, msg, response); err != nil {
		h.lggr.Errorw("failed to send response to gateway", "id", gatewayId, "error", err)
	}
}

func (h *functionsConnectorHandler) sendResponse(ctx context.Context, gatewayId string, request *api.Message, payload any) error {
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	msg := &api.Message{
		Body: api.MessageBody{
			MessageId: request.Body.MessageId,
			DonId:     request.Body.DonId,
			Method:    request.Body.Method,
			Sender:    h.signerID,
			Payload:   payloadJson,
		},
	}
	if err = msg.Sign(h.signerKey); err != nil {
		return err
	}

	err = h.connector.SendToGateway(ctx, gatewayId, msg)
	if err == nil {
		h.lggr.Debugw("sent to gateway", "id", gatewayId, "messageId", request.Body.MessageId, "donId", request.Body.DonId, "method", request.Body.Method)
	}
	return err
}
