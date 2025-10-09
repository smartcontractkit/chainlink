package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jonboulle/clockwork"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	handlerscommon "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/common"
	gw_net "github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

var promRequest = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "gateway_request",
	Help: "Metric to track received requests and response codes",
}, []string{"response_code"})

type Gateway interface {
	job.ServiceCtx
	gw_net.HTTPRequestHandler

	GetUserPort() int
	GetNodePort() int
}

type HandlerType = string

type HandlerFactory interface {
	NewHandler(handlerType HandlerType, handlerConfig json.RawMessage, donConfig *config.DONConfig, don handlers.DON) (handlers.Handler, error)
}

type gateway struct {
	services.StateMachine

	codec              api.Codec
	httpServer         gw_net.HttpServer
	handlers           map[string]handlers.Handler
	serviceNameToDonID map[string]string
	connMgr            ConnectionManager
	lggr               logger.Logger
}

func NewGatewayFromConfig(cfg *config.GatewayConfig, handlerFactory HandlerFactory, lggr logger.Logger) (Gateway, error) {
	codec := &api.JsonRPCCodec{}
	httpServer := gw_net.NewHttpServer(&cfg.UserServerConfig, lggr)
	connMgr, err := NewConnectionManager(cfg, clockwork.NewRealClock(), lggr)
	if err != nil {
		return nil, err
	}

	handlerMap := make(map[string]handlers.Handler)
	serviceNameToDonID := make(map[string]string)

	for _, donConfig := range cfg.Dons {
		_, ok := handlerMap[donConfig.DonId]
		if ok {
			return nil, fmt.Errorf("duplicate DON ID %s", donConfig.DonId)
		}
		donConnMgr := connMgr.DONConnectionManager(donConfig.DonId)
		if donConnMgr == nil {
			return nil, fmt.Errorf("connection manager ID %s not found", donConfig.DonId)
		}
		for idx, nodeConfig := range donConfig.Members {
			donConfig.Members[idx].Address = strings.ToLower(nodeConfig.Address)
			if !common.IsHexAddress(nodeConfig.Address) {
				return nil, fmt.Errorf("invalid node address %s", nodeConfig.Address)
			}
		}

		// Convert old-style handler config to the new style.
		var handlers []config.Handler
		if donConfig.HandlerName != "" {
			handlers = append(handlers, config.Handler{
				Name:   donConfig.HandlerName,
				Config: donConfig.HandlerConfig,
			})
		}

		handlers = append(handlers, donConfig.Handlers...)
		handler, err := NewMultiHandler(handlerFactory, handlers, &donConfig, donConnMgr)
		if err != nil {
			return nil, fmt.Errorf("failed to create multi-handler for DON %s: %w", donConfig.DonId, err)
		}

		handlerMap[donConfig.DonId] = handler

		for _, h := range handlers {
			if h.ServiceName != "" {
				_, ok := serviceNameToDonID[h.ServiceName]
				if ok {
					return nil, fmt.Errorf("duplicate service name %s for DON ID %s", h.ServiceName, donConfig.DonId)
				}

				serviceNameToDonID[h.ServiceName] = donConfig.DonId
			}
		}

		donConnMgr.SetHandler(handler)
	}
	return NewGateway(codec, httpServer, handlerMap, serviceNameToDonID, connMgr, lggr), nil
}

func NewGateway(codec api.Codec, httpServer gw_net.HttpServer, handlers map[string]handlers.Handler, serviceNameToDonID map[string]string, connMgr ConnectionManager, lggr logger.Logger) Gateway {
	gw := &gateway{
		codec:              codec,
		httpServer:         httpServer,
		handlers:           handlers,
		serviceNameToDonID: serviceNameToDonID,
		connMgr:            connMgr,
		lggr:               logger.Named(lggr, "Gateway"),
	}
	httpServer.SetHTTPRequestHandler(gw)
	return gw
}

func (g *gateway) Start(ctx context.Context) error {
	return g.StartOnce("Gateway", func() error {
		g.lggr.Info("starting gateway")
		for _, handler := range g.handlers {
			if err := handler.Start(ctx); err != nil {
				return err
			}
		}
		if err := g.connMgr.Start(ctx); err != nil {
			return err
		}
		return g.httpServer.Start(ctx)
	})
}

func (g *gateway) Close() error {
	return g.StopOnce("Gateway", func() (err error) {
		g.lggr.Info("closing gateway")
		err = errors.Join(err, g.httpServer.Close())
		err = errors.Join(err, g.connMgr.Close())
		for _, handler := range g.handlers {
			err = errors.Join(err, handler.Close())
		}
		return
	})
}

// Called by the server
func (g *gateway) ProcessRequest(ctx context.Context, rawRequest []byte, auth string) (rawResponse []byte, httpStatusCode int) {
	// decode
	jsonRequest, err := jsonrpc2.DecodeRequest[json.RawMessage](rawRequest, auth)
	if err != nil {
		return newError("", api.UserMessageParseError, err.Error())
	}
	msg, err := g.codec.DecodeJSONRequest(jsonRequest)
	if err != nil {
		return newError(jsonRequest.ID, api.UserMessageParseError, err.Error())
	}
	var isLegacyRequest = false
	var handlerKey string
	if msg == nil || msg.Body.DonId == "" {
		// if no DON ID is specified, it is a new JsonRPC request. Use the service name as handler key
		// Let's map the service name to the DON ID and find the handler
		hk, ok := g.serviceNameToDonID[jsonRequest.ServiceName()]
		if !ok {
			return newError(jsonRequest.ID, api.HandlerError, "Service name not found: "+jsonRequest.ServiceName())
		}
		handlerKey = hk
	} else {
		// Means legacy request. Proceed to validate it and fetch DonId
		isLegacyRequest = true
		if err = msg.Validate(); err != nil {
			return newError(jsonRequest.ID, api.UserMessageParseError, err.Error())
		}
		handlerKey = msg.Body.DonId
	}
	h, ok := g.handlers[handlerKey]
	if !ok {
		return newError(jsonRequest.ID, api.UnsupportedDONIdError, "Unsupported DON ID or Handler: "+handlerKey)
	}
	// send to the right handler
	callback := handlerscommon.NewCallback()
	if isLegacyRequest {
		err = h.HandleLegacyUserMessage(ctx, msg, callback)
	} else {
		err = h.HandleJSONRPCUserMessage(ctx, jsonRequest, callback)
	}
	if err != nil {
		return newError(jsonRequest.ID, api.HandlerError, err.Error())
	}
	// await response
	response, err := callback.Wait(ctx)
	if err != nil {
		return newError(jsonRequest.ID, api.RequestTimeoutError, "handler timeout: "+err.Error())
	}

	g.lggr.Debugw("received response from handler", "handler", handlerKey, "response", response, "requestID", jsonRequest.ID)
	promRequest.WithLabelValues(response.ErrorCode.String()).Inc()
	return response.RawResponse, api.ToHttpErrorCode(response.ErrorCode)
}

func newError(id string, errCode api.ErrorCode, errMsg string) ([]byte, int) {
	response := jsonrpc2.Response[json.RawMessage]{
		Version: jsonrpc2.JsonRpcVersion,
		ID:      id,
		Error: &jsonrpc2.WireError{
			Code:    api.ToJSONRPCErrorCode(errCode),
			Message: errMsg,
			Data:    nil,
		},
	}
	rawResponse, err := json.Marshal(response)
	if err != nil {
		rawResponse = []byte("fatal error" + err.Error())
	}
	promRequest.WithLabelValues(errCode.String()).Inc()
	return rawResponse, api.ToHttpErrorCode(errCode)
}

func (g *gateway) GetUserPort() int {
	return g.httpServer.GetPort()
}

func (g *gateway) GetNodePort() int {
	return g.connMgr.GetPort()
}
