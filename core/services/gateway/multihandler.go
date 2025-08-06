package gateway

import (
	"context"
	"encoding/json"
	"fmt"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
)

type multiHandler struct {
	handlers             map[string]handlers.Handler
	handlerTypeForMethod func(string) (HandlerType, error)
}

func NewMultiHandler(handlerFactory HandlerFactory, handlerTypeForMethod func(string) (HandlerType, error), hdlrs []config.Handler, donConfig *config.DONConfig, connMgr *donConnectionManager) (handlers.Handler, error) {
	handlerMap := map[string]handlers.Handler{}
	for _, h := range hdlrs {
		hdlr, err := handlerFactory.NewHandler(h.Name, h.Config, donConfig, connMgr)
		if err != nil {
			return nil, fmt.Errorf("failed to create handler: %w", err)
		}

		handlerMap[h.Name] = hdlr
	}

	return &multiHandler{
		handlers:             handlerMap,
		handlerTypeForMethod: handlerTypeForMethod,
	}, nil
}

func (m *multiHandler) HandleLegacyUserMessage(ctx context.Context, msg *api.Message, callbackCh chan<- handlers.UserCallbackPayload) error {
	h, err := m.getHandler(msg.Body.Method)
	if err != nil {
		return fmt.Errorf("failed to get handler for method %s: %w", msg.Body.Method, err)
	}

	return h.HandleLegacyUserMessage(ctx, msg, callbackCh)
}
func (m *multiHandler) HandleJSONRPCUserMessage(ctx context.Context, jsonRequest jsonrpc.Request[json.RawMessage], callbackCh chan<- handlers.UserCallbackPayload) error {
	h, err := m.getHandler(jsonRequest.Method)
	if err != nil {
		return fmt.Errorf("failed to get handler for method %s: %w", jsonRequest.Method, err)
	}

	return h.HandleJSONRPCUserMessage(ctx, jsonRequest, callbackCh)
}

func (m *multiHandler) HandleNodeMessage(ctx context.Context, resp *jsonrpc.Response[json.RawMessage], nodeAddr string) error {
	h, err := m.getHandler(resp.Method)
	if err != nil {
		return fmt.Errorf("failed to get handler for method %s: %w", resp.Method, err)
	}

	return h.HandleNodeMessage(ctx, resp, nodeAddr)
}

func (m *multiHandler) getHandler(method string) (handlers.Handler, error) {
	// If there's only one handler, return it directly.
	// This preserves backwards compatibility for cases where the method
	// isn't specified on responses (and for cases where only one handler is registered more generally).
	if len(m.handlers) == 1 {
		for _, handler := range m.handlers {
			return handler, nil
		}
	}

	handlerType, err := m.handlerTypeForMethod(method)
	if err != nil {
		return nil, fmt.Errorf("no handler found for method: %w", err)
	}

	handler, ok := m.handlers[handlerType]
	if !ok {
		return nil, fmt.Errorf("no handler registered for method %s (type %s)", method, handlerType)
	}

	return handler, nil
}

func (m *multiHandler) Start(ctx context.Context) error {
	for name, h := range m.handlers {
		if err := h.Start(ctx); err != nil {
			return fmt.Errorf("failed to start handler %s: %w", name, err)
		}
	}
	return nil
}

func (m *multiHandler) Close() error {
	for name, h := range m.handlers {
		if e := h.Close(); e != nil {
			return fmt.Errorf("failed to close handler %s: %w", name, e)
		}
	}
	return nil
}
