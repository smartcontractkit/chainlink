package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
)

type multiHandler struct {
	handlers map[string]handlers.Handler
}

func NewMultiHandler(handlerFactory HandlerFactory, hdlrs []config.Handler, donConfig *config.DONConfig, connMgr *donConnectionManager) (handlers.Handler, error) {

	handlerMap := map[string]handlers.Handler{}
	for _, h := range hdlrs {
		hdlr, err := handlerFactory.NewHandler(h.Name, h.Config, donConfig, connMgr)
		if err != nil {
			return nil, fmt.Errorf("failed to create handler: %w", err)
		}

		handlerMap[h.Name] = hdlr
	}

	return &multiHandler{
		handlers: handlerMap,
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
	// Short-circuit if there is only one handler.
	// This preserves backwards compatibility.
	if len(m.handlers) == 1 {
		for _, h := range m.handlers {
			return h, nil
		}
	}

	// Check that the method is fully-qualified (e.g., "handler.method").
	// Note: this requires callers to have been updated with fully-qualified method names.
	// This change is nevertheless backwards compatible as this behaviour only applies
	// if the job spec has been configured with multiple handlers.
	if strings.Contains(method, ".") {
		parts := strings.Split(method, ".")
		handlerName := parts[0]
		h, ok := m.handlers[handlerName]
		if !ok {
			return nil, fmt.Errorf("handler %s not found for method %s", handlerName, method)
		}

		return h, nil
	}

	return nil, fmt.Errorf("no handler found for method %s", method)
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
