package webapicapabilities

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
)

const (
	// NOTE: more methods will go here. HTTP trigger/action/target; etc.
	MethodWebAPITarget  = "web_api_target"
	MethodWebAPITrigger = "web_api_trigger"
)

type handler struct {
	config         HandlerConfig
	donConfig      *config.DONConfig
	don            handlers.DON
	savedCallbacks map[string]*savedCallback
	mu             sync.Mutex
	lggr           logger.Logger
}

type HandlerConfig struct {
	MaxAllowedMessageAgeSec uint
}
type savedCallback struct {
	id         string
	callbackCh chan<- handlers.UserCallbackPayload
}

var _ handlers.Handler = (*handler)(nil)

func NewWorkflowHandler(handlerConfig json.RawMessage, donConfig *config.DONConfig, don handlers.DON, lggr logger.Logger) (*handler, error) {
	var cfg HandlerConfig
	err := json.Unmarshal(handlerConfig, &cfg)
	if err != nil {
		return nil, err
	}

	return &handler{
		config:         cfg,
		donConfig:      donConfig,
		don:            don,
		savedCallbacks: make(map[string]*savedCallback),
		lggr:           lggr.Named("WorkflowHandler." + donConfig.DonId),
	}, nil
}

func (d *handler) HandleUserMessage(ctx context.Context, msg *api.Message, callbackCh chan<- handlers.UserCallbackPayload) error {
	d.mu.Lock()
	d.savedCallbacks[msg.Body.MessageId] = &savedCallback{msg.Body.MessageId, callbackCh}
	don := d.don
	d.mu.Unlock()
	body := msg.Body
	var payload TriggerRequestPayload
	err := json.Unmarshal(body.Payload, &payload)
	if err != nil {
		d.lggr.Errorw("error decoding payload", "err", err)
		callbackCh <- handlers.UserCallbackPayload{Msg: msg, ErrCode: api.UserMessageParseError, ErrMsg: fmt.Sprintf("error decoding payload %s", err.Error())}
		close(callbackCh)
		return nil
	}

	if payload.Timestamp == 0 {
		d.lggr.Errorw("error decoding payload")
		callbackCh <- handlers.UserCallbackPayload{Msg: msg, ErrCode: api.UserMessageParseError, ErrMsg: "error decoding payload"}
		close(callbackCh)
		return nil
	}

	if uint(time.Now().Unix())-d.config.MaxAllowedMessageAgeSec > uint(payload.Timestamp) {
		callbackCh <- handlers.UserCallbackPayload{Msg: msg, ErrCode: api.HandlerError, ErrMsg: "stale message"}
		close(callbackCh)
		return nil
	}
	// TODO: apply allowlist and rate-limiting here
	if msg.Body.Method != MethodWebAPITrigger {
		d.lggr.Errorw("unsupported method", "method", body.Method)
		callbackCh <- handlers.UserCallbackPayload{Msg: msg, ErrCode: api.HandlerError, ErrMsg: fmt.Sprintf("invalid method %s", msg.Body.Method)}
		close(callbackCh)
		return nil
	}

	// Send to all nodes.
	for _, member := range d.donConfig.Members {
		err = multierr.Combine(err, don.SendToNode(ctx, member.Address, msg))
	}
	return err
}

func (d *handler) HandleNodeMessage(ctx context.Context, msg *api.Message, _ string) error {
	d.mu.Lock()
	savedCb, found := d.savedCallbacks[msg.Body.MessageId]
	delete(d.savedCallbacks, msg.Body.MessageId)
	d.mu.Unlock()

	if found {
		// Send first response from a node back to the user, ignore any other ones.
		// TODO: in practice, we should wait for at least 2F+1 nodes to respond and then return an aggregated response
		// back to the user.
		savedCb.callbackCh <- handlers.UserCallbackPayload{Msg: msg, ErrCode: api.NoError, ErrMsg: ""}
		close(savedCb.callbackCh)
	}
	return nil
}

func (d *handler) Start(context.Context) error {
	return nil
}

func (d *handler) Close() error {
	return nil
}
