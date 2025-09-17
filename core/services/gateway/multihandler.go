package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"slices"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
)

type multiHandler struct {
	typeToHandler   map[string]handlers.Handler
	methodToHandler map[string]handlers.Handler
}

func NewMultiHandler(handlerFactory HandlerFactory, hdlrs []config.Handler, donConfig *config.DONConfig, connMgr *donConnectionManager) (handlers.Handler, error) {
	methodToHandler := map[string]handlers.Handler{}
	typeToHandler := map[string]handlers.Handler{}
	for _, h := range hdlrs {
		hdlr, err := handlerFactory.NewHandler(h.Name, h.Config, donConfig, connMgr)
		if err != nil {
			return nil, fmt.Errorf("failed to create handler %s: %w", h.Name, err)
		}

		typeToHandler[h.Name] = hdlr

		for _, method := range hdlr.Methods() {
			if _, exists := methodToHandler[method]; exists {
				return nil, fmt.Errorf("duplicate handler for method %s: methods must be globally unique across handlers", method)
			}

			methodToHandler[method] = hdlr
		}
	}

	return &multiHandler{
		methodToHandler: methodToHandler,
		typeToHandler:   typeToHandler,
	}, nil
}

func (m *multiHandler) Methods() []string {
	return slices.Collect(maps.Keys(m.methodToHandler))
}

func (m *multiHandler) HandleLegacyUserMessage(ctx context.Context, msg *api.Message, callback handlers.Callback) error {
	h, err := m.getHandler(msg.Body.Method)
	if err != nil {
		return fmt.Errorf("failed to get handler for method %s: %w", msg.Body.Method, err)
	}

	return h.HandleLegacyUserMessage(ctx, msg, callback)
}

func (m *multiHandler) HandleJSONRPCUserMessage(ctx context.Context, jsonRequest jsonrpc.Request[json.RawMessage], callback handlers.Callback) error {
	h, err := m.getHandler(jsonRequest.Method)
	if err != nil {
		return fmt.Errorf("failed to get handler for method %s: %w", jsonRequest.Method, err)
	}

	return h.HandleJSONRPCUserMessage(ctx, jsonRequest, callback)
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
	if len(m.typeToHandler) == 1 {
		for _, handler := range m.typeToHandler {
			return handler, nil
		}
	}

	handler, ok := m.methodToHandler[method]
	if !ok {
		return nil, errors.New("no handler found for method " + method)
	}

	return handler, nil
}

func (m *multiHandler) Start(ctx context.Context) error {
	for name, h := range m.typeToHandler {
		if err := h.Start(ctx); err != nil {
			return fmt.Errorf("failed to start handler %s: %w", name, err)
		}
	}
	return nil
}

func (m *multiHandler) Close() error {
	for name, h := range m.typeToHandler {
		if e := h.Close(); e != nil {
			return fmt.Errorf("failed to close handler %s: %w", name, e)
		}
	}
	return nil
}
