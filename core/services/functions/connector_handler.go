package functions

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"

	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	hc "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/functions"
	"github.com/smartcontractkit/chainlink/v2/core/services/s4"
	"github.com/smartcontractkit/chainlink/v2/core/utils"

	ethCommon "github.com/ethereum/go-ethereum/common"
)

type functionsConnectorHandler struct {
	utils.StartStopOnce

	connector      connector.GatewayConnector
	signerKey      *ecdsa.PrivateKey
	nodeAddress    string
	storage        s4.Storage
	allowlist      functions.OnchainAllowlist
	rateLimiter    *hc.RateLimiter
	subscriptions  functions.OnchainSubscriptions
	minimumBalance assets.Link
	lggr           logger.Logger
}

var (
	_ connector.Signer                  = &functionsConnectorHandler{}
	_ connector.GatewayConnectorHandler = &functionsConnectorHandler{}
)

func NewFunctionsConnectorHandler(nodeAddress string, signerKey *ecdsa.PrivateKey, storage s4.Storage, allowlist functions.OnchainAllowlist, rateLimiter *hc.RateLimiter, subscriptions functions.OnchainSubscriptions, minimumBalance assets.Link, lggr logger.Logger) (*functionsConnectorHandler, error) {
	if signerKey == nil || storage == nil || allowlist == nil || rateLimiter == nil || subscriptions == nil {
		return nil, fmt.Errorf("signerKey, storage, allowlist, rateLimiter and subscriptions must be non-nil")
	}
	return &functionsConnectorHandler{
		nodeAddress:    nodeAddress,
		signerKey:      signerKey,
		storage:        storage,
		allowlist:      allowlist,
		rateLimiter:    rateLimiter,
		subscriptions:  subscriptions,
		minimumBalance: minimumBalance,
		lggr:           lggr.Named("FunctionsConnectorHandler"),
	}, nil
}

func (h *functionsConnectorHandler) SetConnector(connector connector.GatewayConnector) {
	h.connector = connector
}

func (h *functionsConnectorHandler) Sign(data ...[]byte) ([]byte, error) {
	return common.SignData(h.signerKey, data...)
}

func (h *functionsConnectorHandler) HandleGatewayMessage(ctx context.Context, gatewayId string, msg *api.Message) {
	body := &msg.Body
	fromAddr := ethCommon.HexToAddress(body.Sender)
	if !h.allowlist.Allow(fromAddr) {
		h.lggr.Errorw("allowlist prevented the request from this address", "id", gatewayId, "address", fromAddr)
		return
	}
	if !h.rateLimiter.Allow(body.Sender) {
		h.lggr.Errorw("request rate-limited", "id", gatewayId, "address", fromAddr)
		return
	}
	if balance, err := h.subscriptions.GetMaxUserBalance(fromAddr); err != nil || balance.Cmp(h.minimumBalance.ToInt()) < 0 {
		h.lggr.Errorw("user subscription has insufficient balance", "id", gatewayId, "address", fromAddr, "balance", balance, "minBalance", h.minimumBalance)
		response := functions.SecretsResponseBase{
			Success:      false,
			ErrorMessage: "user subscription has insufficient balance",
		}
		if err := h.sendResponse(ctx, gatewayId, body, response); err != nil {
			h.lggr.Errorw("failed to send response to gateway", "id", gatewayId, "error", err)
		}
		return
	}

	h.lggr.Debugw("handling gateway request", "id", gatewayId, "method", body.Method)

	switch body.Method {
	case functions.MethodSecretsList:
		h.handleSecretsList(ctx, gatewayId, body, fromAddr)
	case functions.MethodSecretsSet:
		h.handleSecretsSet(ctx, gatewayId, body, fromAddr)
	default:
		h.lggr.Errorw("unsupported method", "id", gatewayId, "method", body.Method)
	}
}

func (h *functionsConnectorHandler) Start(ctx context.Context) error {
	return h.StartOnce("FunctionsConnectorHandler", func() error {
		if err := h.allowlist.Start(ctx); err != nil {
			return err
		}
		return h.subscriptions.Start(ctx)
	})
}

func (h *functionsConnectorHandler) Close() error {
	return h.StopOnce("FunctionsConnectorHandler", func() (err error) {
		err = multierr.Combine(err, h.allowlist.Close())
		err = multierr.Combine(err, h.subscriptions.Close())
		return
	})
}

func (h *functionsConnectorHandler) handleSecretsList(ctx context.Context, gatewayId string, body *api.MessageBody, fromAddr ethCommon.Address) {
	var response functions.SecretsListResponse
	snapshot, err := h.storage.List(ctx, fromAddr)
	if err == nil {
		response.Success = true
		response.Rows = make([]functions.SecretsListRow, len(snapshot))
		for i, row := range snapshot {
			response.Rows[i] = functions.SecretsListRow{
				SlotID:     row.SlotId,
				Version:    row.Version,
				Expiration: row.Expiration,
			}
		}
	} else {
		response.ErrorMessage = fmt.Sprintf("Failed to list secrets: %v", err)
	}

	if err := h.sendResponse(ctx, gatewayId, body, response); err != nil {
		h.lggr.Errorw("failed to send response to gateway", "id", gatewayId, "error", err)
	}
}

func (h *functionsConnectorHandler) handleSecretsSet(ctx context.Context, gatewayId string, body *api.MessageBody, fromAddr ethCommon.Address) {
	var request functions.SecretsSetRequest
	var response functions.SecretsSetResponse
	err := json.Unmarshal(body.Payload, &request)
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
		h.lggr.Debugw("handling a secrets_set request", "address", fromAddr, "slotId", request.SlotID, "payloadVersion", request.Version, "expiration", request.Expiration)
		err = h.storage.Put(ctx, &key, &record, request.Signature)
		if err == nil {
			response.Success = true
		} else {
			response.ErrorMessage = fmt.Sprintf("Failed to set secret: %v", err)
		}
	} else {
		response.ErrorMessage = fmt.Sprintf("Bad request to set secret: %v", err)
	}

	if err := h.sendResponse(ctx, gatewayId, body, response); err != nil {
		h.lggr.Errorw("failed to send response to gateway", "id", gatewayId, "error", err)
	}
}

func (h *functionsConnectorHandler) sendResponse(ctx context.Context, gatewayId string, requestBody *api.MessageBody, payload any) error {
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

	err = h.connector.SendToGateway(ctx, gatewayId, msg)
	if err == nil {
		h.lggr.Debugw("sent to gateway", "id", gatewayId, "messageId", requestBody.MessageId, "donId", requestBody.DonId, "method", requestBody.Method)
	}
	return err
}
