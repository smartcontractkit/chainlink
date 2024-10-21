package webapicapabilities

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
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
	MethodWebAPITarget                = "web_api_target"
	MethodWebAPITrigger               = "web_api_trigger"
	MethodWebAPITriggerUpdateMetadata = "web_api_trigger_update_metadata"
)

type TriggersConfig struct {
	triggersConfig map[string]webapicap.TriggerConfig
}
type NodeTriggerConfig struct {
	lastUpdatedAt  time.Time
	triggerConfigs TriggersConfig
}

type AllNodesTriggersConfig struct {
	capabilities.Validator[TriggersConfig, struct{}, capabilities.TriggerResponse]

	triggersConfigMap map[string]NodeTriggerConfig
}

type handler struct {
	capabilities.Validator[webapicap.TriggerConfig, struct{}, capabilities.TriggerResponse]

	config          HandlerConfig
	don             handlers.DON
	donConfig       *config.DONConfig
	savedCallbacks  map[string]*savedCallback
	mu              sync.Mutex
	lggr            logger.Logger
	httpClient      network.HTTPClient
	nodeRateLimiter *common.RateLimiter
	wg              sync.WaitGroup
	// each gateway node has a map of trigger IDs to trigger configs
	triggersConfig AllNodesTriggersConfig
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
	lggr.Debugf("new web api handler with config: %s", string(handlerConfig))
	var cfg HandlerConfig
	err := json.Unmarshal(handlerConfig, &cfg)
	if err != nil {
		lggr.Errorf("error unmarshalling config: %s, err: %s", string(handlerConfig), err.Error())
		return nil, err
	}
	lggr.Debugw("new web api handler", "parsedConfig", cfg)

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
		triggersConfig:  AllNodesTriggersConfig{},
	}, nil
}

// sendHTTPMessageToClient is an outgoing message from the gateway to external endpoints
// returns message to be sent back to the capability node
func (h *handler) sendHTTPMessageToClient(ctx context.Context, req network.HTTPRequest, msg *api.Message) (*api.Message, error) {
	var payload TargetResponsePayload
	resp, err := h.httpClient.Send(ctx, req)
	if err != nil {
		return nil, err
	}
	payload = TargetResponsePayload{
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
	var targetPayload TargetRequestPayload
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
			payload := TargetResponsePayload{
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
	h.lggr.Debugw("handling web api trigger message", "messageId", msg.Body.MessageId, "nodeAddr", nodeAddr)

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

//	body := api.MessageBody{
//		MessageId: types.RandomID().String(),
//		DonId:     h.connector.DonID(),
//		Method:    webapicapabilities.MethodWebAPITriggerUpdateMetadata,
//		Receiver:  gatewayIDStr,
//		Payload:   payloadJSON,
//	}
func (h *handler) handleWebAPITriggerUpdateMetadata(ctx context.Context, msg *api.Message, nodeAddr string) error {
	body := msg.Body
	h.lggr.Debugw("handleWebAPITriggerUpdateMetadata", "body", body, "payload", string(body.Payload))

	var payload *values.Map
	err := json.Unmarshal(body.Payload, &payload)
	if err != nil {
		// errors here:
		// error decoding payload	{"version": "unset@unset",
		// "err": "json: cannot unmarshal object into Go struct field Map.Underlying of type values.Value",
		// "payload": "{\"Underlying\":{\"AllowedSenders\":{\"Underlying\":[{\"Underlying\":\"0x853d51d5d9935964267a5050aC53aa63ECA39bc5\"}]},\"AllowedTopics\":{\"Underlying\":[{\"Underlying\":\"daily_price_update\"},{\"Underlying\":\"ad_hoc_price_update\"}]},\"RateLimiter\":{\"Underlying\":{\"GlobalBurst\":{\"Underlying\":101},\"GlobalRPS\":{\"Underlying\":100},\"PerSenderBurst\":{\"Underlying\":103},\"PerSenderRPS\":{\"Underlying\":102}}},\"RequiredParams\":{\"Underlying\":[{\"Underlying\":\"bid\"},{\"Underlying\":\"ask\"}]}}}"}
		h.lggr.Errorw("error decoding payload", "err", err, "payload", string(body.Payload))
		// callbackCh <- handlers.UserCallbackPayload{Msg: msg, ErrCode: api.UserMessageParseError, ErrMsg: fmt.Sprintf("error decoding payload %s", err.Error())}
		// close(callbackCh)
		return err
	}
	reqConfig, err := h.triggersConfig.ValidateConfig(payload)
	if err != nil {
		h.lggr.Errorw("error validating config", "err", err)
		return err
	}
	h.triggersConfig.triggersConfigMap[body.DonId] = NodeTriggerConfig{lastUpdatedAt: time.Now(), triggerConfigs: *reqConfig}
	// h.updateTriggerConsensus()
	return nil
}

func (h *handler) updateTriggerConsensus() {

}

func (h *handler) HandleNodeMessage(ctx context.Context, msg *api.Message, nodeAddr string) error {
	switch msg.Body.Method {
	case MethodWebAPITrigger:
		return h.handleWebAPITriggerMessage(ctx, msg, nodeAddr)
	case MethodWebAPITarget:
		return h.handleWebAPITargetMessage(ctx, msg, nodeAddr)
	case MethodWebAPITriggerUpdateMetadata:
		return h.handleWebAPITriggerUpdateMetadata(ctx, msg, nodeAddr)
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
	h.lggr.Debugw("handling web api target user message", "messageId", msg.Body.MessageId)
	h.mu.Lock()
	h.savedCallbacks[msg.Body.MessageId] = &savedCallback{msg.Body.MessageId, callbackCh}
	don := h.don
	h.mu.Unlock()
	body := msg.Body
	var payload TriggerRequestPayload
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
