package webapicapabilities

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
)

const (
	// NOTE: more methods will go here. HTTP trigger/action/target; etc.
	MethodWebAPITarget = "web_api_target"
)

type handler struct {
	don             handlers.DON
	lggr            logger.Logger
	httpClient      network.HTTPClient
	nodeRateLimiter *common.RateLimiter
}

type HandlerConfig struct {
	NodeRateLimiter common.RateLimiterConfig `json:"nodeRateLimiter"`
}

var _ handlers.Handler = (*handler)(nil)

func NewHandler(handlerConfig json.RawMessage, donConfig *config.DONConfig, don handlers.DON, httpClient network.HTTPClient, lggr logger.Logger) (*handler, error) {
	var cfg HandlerConfig
	err := json.Unmarshal(handlerConfig, &cfg)
	if err != nil {
		return nil, err
	}
	nodeRateLimiter, err := common.NewRateLimiter(cfg.NodeRateLimiter)
	if err != nil {
		return nil, err
	}

	return &handler{
		don:             don,
		lggr:            lggr.Named("WebAPIHandler." + donConfig.DonId),
		httpClient:      httpClient,
		nodeRateLimiter: nodeRateLimiter,
	}, nil
}

func (h *handler) HandleUserMessage(ctx context.Context, msg *api.Message, callbackCh chan<- handlers.UserCallbackPayload) error {
	return nil
}

// sendHTTPMessageToClient is an outgoing message from the gateway to external endpoints
// returns message to be sent back to the capability node
func (h *handler) sendHTTPMessageToClient(ctx context.Context, req network.HTTPRequest, msg *api.Message) (*api.Message, error) {
	var payload TargetResponsePayload
	resp, err := h.httpClient.Send(ctx, req)
	if err != nil {
		return nil, err
	} else {
		payload = TargetResponsePayload{
			Success:    true,
			StatusCode: uint8(resp.StatusCode),
			Headers:    resp.Headers,
			Body:       resp.Body,
		}
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return &api.Message{
		Body: api.MessageBody{
			MessageId: msg.Body.MessageId,
			Method:    msg.Body.Method,
			DonId:     msg.Body.DonId,
			Payload:   payloadBytes,
		},
	}, nil
}

func (h *handler) handleWebAPITargetMessage(ctx context.Context, msg *api.Message, nodeAddr string) error {
	h.lggr.Debugw("handling web api target message", "messageId", msg.Body.MessageId, "nodeAddr", nodeAddr)
	if !h.nodeRateLimiter.Allow(nodeAddr) {
		return fmt.Errorf("rate limit exceeded for node %s", nodeAddr)
	}
	var targetPayload TargetRequestPayload
	err := json.Unmarshal(msg.Body.Payload, &targetPayload)
	if err != nil {
		return err
	}
	// send message to target
	req := network.HTTPRequest{
		Method:  targetPayload.Method,
		URL:     targetPayload.URL,
		Headers: targetPayload.Headers,
		Body:    targetPayload.Body,
		Timeout: time.Duration(targetPayload.TimeoutMs) * time.Millisecond,
	}
	// this handle method must be non-blocking
	// send response to node (target capability) async
	// if there is a non-HTTP error (e.g. malformed request), send payload with success set to false and error messages
	go func() {
		l := h.lggr.With("url", targetPayload.URL, "messageId", msg.Body.MessageId, "method", targetPayload.Method)
		respMsg, err := h.sendHTTPMessageToClient(ctx, req, msg)
		if err != nil {
			l.Errorw("error while sending HTTP request to external endpoint", "err", err)
			payload := TargetResponsePayload{
				Success:      false,
				ErrorMessage: err.Error(),
			}
			payloadBytes, err := json.Marshal(payload)
			if err != nil {
				// should not happen
				l.Errorw("error while marshalling payload", "err", err)
				return
			}
			respMsg = &api.Message{
				Body: api.MessageBody{
					MessageId: msg.Body.MessageId,
					Method:    msg.Body.Method,
					DonId:     msg.Body.DonId,
					Payload:   payloadBytes,
				},
			}
		}
		err = h.don.SendToNode(ctx, nodeAddr, respMsg)
		if err != nil {
			l.Errorw("failed to send to node", "err", err, "to", nodeAddr)
			return
		}
	}()
	return nil
}

func (h *handler) HandleNodeMessage(ctx context.Context, msg *api.Message, nodeAddr string) error {
	switch msg.Body.Method {
	case MethodWebAPITarget:
		return h.handleWebAPITargetMessage(ctx, msg, nodeAddr)
	default:
		return fmt.Errorf("unsupported method: %s", msg.Body.Method)
	}
}

func (h *handler) Start(context.Context) error {
	return nil
}

func (h *handler) Close() error {
	return nil
}
