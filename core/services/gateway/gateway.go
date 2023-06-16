package gateway

import (
	"context"
	"encoding/json"
	"fmt"

	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers"
	gw_net "github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type Gateway interface {
	job.ServiceCtx
	gw_net.HTTPRequestHandler
}

type HandlerType = string

type HandlerFactory interface {
	NewHandler(handlerType HandlerType, handlerConfig json.RawMessage, donConfig *config.DONConfig, don handlers.DON) (handlers.Handler, error)
}

type gateway struct {
	utils.StartStopOnce

	codec      api.Codec
	httpServer gw_net.HttpServer
	handlers   map[string]handlers.Handler
	connMgr    ConnectionManager
	lggr       logger.Logger
}

func NewGatewayFromConfig(config *config.GatewayConfig, handlerFactory HandlerFactory, lggr logger.Logger) (Gateway, error) {
	codec := &api.JsonRPCCodec{}
	httpServer := gw_net.NewHttpServer(&config.UserServerConfig, lggr)
	connMgr, err := NewConnectionManager(config, utils.NewRealClock(), lggr)
	if err != nil {
		return nil, err
	}

	handlerMap := make(map[string]handlers.Handler)
	for _, donConfig := range config.Dons {
		donConfig := donConfig
		_, ok := handlerMap[donConfig.DonId]
		if ok {
			return nil, fmt.Errorf("duplicate DON ID %s", donConfig.DonId)
		}
		donConnMgr := connMgr.DONConnectionManager(donConfig.DonId)
		if donConnMgr == nil {
			return nil, fmt.Errorf("connection manager ID %s not found", donConfig.DonId)
		}
		handler, err := handlerFactory.NewHandler(donConfig.HandlerName, donConfig.HandlerConfig, &donConfig, donConnMgr)
		if err != nil {
			return nil, err
		}
		handlerMap[donConfig.DonId] = handler
		donConnMgr.SetHandler(handler)
	}
	return NewGateway(codec, httpServer, handlerMap, connMgr, lggr), nil
}

func NewGateway(codec api.Codec, httpServer gw_net.HttpServer, handlers map[string]handlers.Handler, connMgr ConnectionManager, lggr logger.Logger) Gateway {
	gw := &gateway{
		codec:      codec,
		httpServer: httpServer,
		handlers:   handlers,
		connMgr:    connMgr,
		lggr:       lggr.Named("gateway"),
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
		err = multierr.Combine(err, g.httpServer.Close())
		err = multierr.Combine(err, g.connMgr.Close())
		for _, handler := range g.handlers {
			err = multierr.Combine(err, handler.Close())
		}
		return
	})
}

// Called by the server
func (g *gateway) ProcessRequest(ctx context.Context, rawRequest []byte) (rawResponse []byte, httpStatusCode int) {
	// decode
	msg, err := g.codec.DecodeRequest(rawRequest)
	if err != nil {
		return newError(g.codec, "", api.UserMessageParseError, err.Error())
	}
	// find correct handler
	handler, ok := g.handlers[msg.Body.DonId]
	if !ok {
		return newError(g.codec, msg.Body.MessageId, api.UnsupportedDONIdError, "unsupported DON ID")
	}
	// send to the handler
	responseCh := make(chan handlers.UserCallbackPayload, 1)
	err = handler.HandleUserMessage(ctx, msg, responseCh)
	if err != nil {
		return newError(g.codec, msg.Body.MessageId, api.InternalHandlerError, err.Error())
	}
	// await response
	var response handlers.UserCallbackPayload
	select {
	case <-ctx.Done():
		return newError(g.codec, msg.Body.MessageId, api.RequestTimeoutError, "handler timeout")
	case response = <-responseCh:
		break
	}
	if response.ErrCode != api.NoError {
		return newError(g.codec, msg.Body.MessageId, response.ErrCode, response.ErrMsg)
	}
	// encode
	rawResponse, err = g.codec.EncodeResponse(response.Msg)
	if err != nil {
		return newError(g.codec, msg.Body.MessageId, api.NodeReponseEncodingError, "")
	}
	return rawResponse, api.ToHttpErrorCode(api.NoError)
}

func newError(codec api.Codec, id string, errCode api.ErrorCode, errMsg string) ([]byte, int) {
	rawResponse, err := codec.EncodeNewErrorResponse(id, api.ToJsonRPCErrorCode(errCode), errMsg, nil)
	if err != nil {
		// we're not even able to encode a valid JSON response
		return []byte("fatal error"), api.ToHttpErrorCode(api.FatalError)
	}
	return rawResponse, api.ToHttpErrorCode(errCode)
}
