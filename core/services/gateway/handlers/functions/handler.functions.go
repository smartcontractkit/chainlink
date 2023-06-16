package functions

import (
	"context"
	"encoding/json"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
)

type FunctionsHandlerConfig struct {
}

type functionsHandler struct {
	handlerConfig *FunctionsHandlerConfig
	donConfig     *config.DONConfig
	don           handlers.DON
	lggr          logger.Logger
}

var _ handlers.Handler = (*functionsHandler)(nil)

func NewFunctionsHandler(handlerConfig json.RawMessage, donConfig *config.DONConfig, don handlers.DON, lggr logger.Logger) (handlers.Handler, error) {
	var parsedConfig FunctionsHandlerConfig
	if err := json.Unmarshal(handlerConfig, &parsedConfig); err != nil {
		return nil, err
	}
	return &functionsHandler{
		handlerConfig: &parsedConfig,
		donConfig:     donConfig,
		don:           don,
		lggr:          lggr,
	}, nil
}

func (h *functionsHandler) HandleUserMessage(ctx context.Context, msg *api.Message, callbackCh chan<- handlers.UserCallbackPayload) error {
	if err := msg.Validate(); err != nil {
		h.lggr.Debug("received invalid message", "err", err)
		return err
	}
	h.lggr.Debug("received a valid message", "sender", msg.Body.Sender)
	return nil
}

func (h *functionsHandler) HandleNodeMessage(ctx context.Context, msg *api.Message, nodeAddr string) error {
	return nil
}

func (h *functionsHandler) Start(context.Context) error {
	return nil
}

func (h *functionsHandler) Close() error {
	return nil
}
