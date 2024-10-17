package capabilities

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/webapi/webapicap"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
)

const (
	// NOTE: more methods will go here. HTTP trigger/action/target; etc.
	MethodWebAPITarget  = "web_api_target"
	MethodWebAPITrigger = "web_api_trigger"
	MethodComputeAction = "compute_action"
)

type handler struct {
	config          HandlerConfig
	don             handlers.DON
	donConfig       *config.DONConfig
	savedCallbacks  map[string]*savedCallback
	mu              sync.Mutex
	lggr            logger.Logger
	httpClient      network.HTTPClient
	nodeRateLimiter *common.RateLimiter
	wg              sync.WaitGroup
}

type HandlerConfig struct {
	NodeRateLimiter         common.RateLimiterConfig `json:"nodeRateLimiter"`
	MaxAllowedMessageAgeSec uint                     `json:"maxAllowedMessageAgeSec"`
}

type savedCallback struct {
	id         string
	callbackCh chan<- handlers.UserCallbackPayload
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
		config:          cfg,
		don:             don,
		donConfig:       donConfig,
		lggr:            lggr.Named("WebAPIHandler." + donConfig.DonId),
		httpClient:      httpClient,
		nodeRateLimiter: nodeRateLimiter,
		wg:              sync.WaitGroup{},
		savedCallbacks:  make(map[string]*savedCallback),
	}, nil
}

// sendHTTPMessageToClient is an outgoing message from the gateway to external endpoints
// returns message to be sent back to the capability node
func (h *handler) sendHTTPMessageToClient(ctx context.Context, req network.HTTPRequest, msg *api.Message) (*api.Message, error) {
	var payload Response
	resp, err := h.httpClient.Send(ctx, req)
	if err != nil {
		return nil, err
	}
	payload = Response{
		ExecutionError: false,
		StatusCode:     resp.StatusCode,
		Headers:        resp.Headers,
		Body:           resp.Body,
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
	var targetPayload Request
	err := json.Unmarshal(msg.Body.Payload, &targetPayload)
	if err != nil {
		return err
	}
	// send message to target
	timeout := time.Duration(targetPayload.TimeoutMs) * time.Millisecond
	req := network.HTTPRequest{
		Method:  targetPayload.Method,
		URL:     targetPayload.URL,
		Headers: targetPayload.Headers,
		Body:    targetPayload.Body,
		Timeout: timeout,
	}
	// this handle method must be non-blocking
	// send response to node (target capability) async
	// if there is a non-HTTP error (e.g. malformed request), send payload with success set to false and error messages
	h.wg.Add(1)
	go func() {
		defer h.wg.Done()
		// not cancelled when parent is cancelled to ensure the goroutine can finish
		newCtx := context.WithoutCancel(ctx)
		newCtx, cancel := context.WithTimeout(newCtx, timeout)
		defer cancel()
		l := h.lggr.With("url", targetPayload.URL, "messageId", msg.Body.MessageId, "method", targetPayload.Method)
		respMsg, err := h.sendHTTPMessageToClient(newCtx, req, msg)
		if err != nil {
			l.Errorw("error while sending HTTP request to external endpoint", "err", err)
			payload := Response{
				ExecutionError: true,
				ErrorMessage:   err.Error(),
			}
			payloadBytes, err2 := json.Marshal(payload)
			if err2 != nil {
				// should not happen
				l.Errorw("error while marshalling payload", "err", err2)
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
		err = h.don.SendToNode(newCtx, nodeAddr, respMsg)
		if err != nil {
			l.Errorw("failed to send to node", "err", err, "to", nodeAddr)
			return
		}
	}()
	return nil
}

func (h *handler) handleWebAPITriggerMessage(ctx context.Context, msg *api.Message, nodeAddr string) error {
	h.mu.Lock()
	savedCb, found := h.savedCallbacks[msg.Body.MessageId]
	delete(h.savedCallbacks, msg.Body.MessageId)
	h.mu.Unlock()

	if found {
		// Send first response from a node back to the user, ignore any other ones.
		// TODO: in practice, we should wait for at least 2F+1 nodes to respond and then return an aggregated response
		// back to the user.
		savedCb.callbackCh <- handlers.UserCallbackPayload{Msg: msg, ErrCode: api.NoError, ErrMsg: ""}
		close(savedCb.callbackCh)
	}
	return nil
}

func (h *handler) handleWebAPIOutgoingMessage(ctx context.Context, msg *api.Message, nodeAddr string) error {
	h.lggr.Debugw("handling webAPI outgoing message", "messageId", msg.Body.MessageId, "nodeAddr", nodeAddr)
	if !h.nodeRateLimiter.Allow(nodeAddr) {
		return fmt.Errorf("rate limit exceeded for node %s", nodeAddr)
	}
	var payload Request
	err := json.Unmarshal(msg.Body.Payload, &payload)
	if err != nil {
		return err
	}

	timeout := time.Duration(payload.TimeoutMs) * time.Millisecond
	req := network.HTTPRequest{
		Method:  payload.Method,
		URL:     payload.URL,
		Headers: payload.Headers,
		Body:    payload.Body,
		Timeout: timeout,
	}

	// send response to node async
	h.wg.Add(1)
	go func() {
		defer h.wg.Done()
		// not cancelled when parent is cancelled to ensure the goroutine can finish
		newCtx := context.WithoutCancel(ctx)
		newCtx, cancel := context.WithTimeout(newCtx, timeout)
		defer cancel()
		l := h.lggr.With("url", payload.URL, "messageId", msg.Body.MessageId, "method", payload.Method)
		respMsg, err := h.sendHTTPMessageToClient(newCtx, req, msg)
		if err != nil {
			l.Errorw("error while sending HTTP request to external endpoint", "err", err)
			payload := Response{
				ExecutionError: true,
				ErrorMessage:   err.Error(),
			}
			payloadBytes, err2 := json.Marshal(payload)
			if err2 != nil {
				// should not happen
				l.Errorw("error while marshalling payload", "err", err2)
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
		err = h.don.SendToNode(newCtx, nodeAddr, respMsg)
		if err != nil {
			l.Errorw("failed to send to node", "err", err, "to", nodeAddr)
			return
		}
	}()
	return nil
}

func (h *handler) HandleNodeMessage(ctx context.Context, msg *api.Message, nodeAddr string) error {
	switch msg.Body.Method {
	case MethodWebAPITrigger:
		return h.handleWebAPITriggerMessage(ctx, msg, nodeAddr)
	case MethodWebAPITarget, MethodComputeAction:
		return h.handleWebAPIOutgoingMessage(ctx, msg, nodeAddr)
	default:
		return fmt.Errorf("unsupported method: %s", msg.Body.Method)
	}
}

func (h *handler) Start(context.Context) error {
	return nil
}

func (h *handler) Close() error {
	h.wg.Wait()
	return nil
}

func (h *handler) HandleUserMessage(ctx context.Context, msg *api.Message, callbackCh chan<- handlers.UserCallbackPayload) error {
	h.mu.Lock()
	h.savedCallbacks[msg.Body.MessageId] = &savedCallback{msg.Body.MessageId, callbackCh}
	don := h.don
	h.mu.Unlock()
	body := msg.Body
	var payload webapicap.TriggerRequestPayload
	err := json.Unmarshal(body.Payload, &payload)
	if err != nil {
		h.lggr.Errorw("error decoding payload", "err", err)
		callbackCh <- handlers.UserCallbackPayload{Msg: msg, ErrCode: api.UserMessageParseError, ErrMsg: fmt.Sprintf("error decoding payload %s", err.Error())}
		close(callbackCh)
		return nil
	}

	if payload.Timestamp == 0 {
		h.lggr.Errorw("error decoding payload")
		callbackCh <- handlers.UserCallbackPayload{Msg: msg, ErrCode: api.UserMessageParseError, ErrMsg: "error decoding payload"}
		close(callbackCh)
		return nil
	}

	if uint(time.Now().Unix())-h.config.MaxAllowedMessageAgeSec > uint(payload.Timestamp) {
		callbackCh <- handlers.UserCallbackPayload{Msg: msg, ErrCode: api.HandlerError, ErrMsg: "stale message"}
		close(callbackCh)
		return nil
	}
	// TODO: apply allowlist and rate-limiting here
	if msg.Body.Method != MethodWebAPITrigger {
		h.lggr.Errorw("unsupported method", "method", body.Method)
		callbackCh <- handlers.UserCallbackPayload{Msg: msg, ErrCode: api.HandlerError, ErrMsg: fmt.Sprintf("invalid method %s", msg.Body.Method)}
		close(callbackCh)
		return nil
	}

	// Send to all nodes.
	for _, member := range h.donConfig.Members {
		err = multierr.Combine(err, don.SendToNode(ctx, member.Address, msg))
	}
	return err
}
