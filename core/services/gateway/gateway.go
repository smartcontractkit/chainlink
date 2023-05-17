package gateway

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	gw_net "github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type Gateway interface {
	job.ServiceCtx
	gw_net.HTTPRequestHandler
}

type gateway struct {
	utils.StartStopOnce

	codec      Codec
	httpServer gw_net.HttpServer
	handlers   map[string]Handler
	connMgr    ConnectionManager
	lggr       logger.Logger
}

func NewGatewayFromConfig(config *GatewayConfig, lggr logger.Logger) (Gateway, error) {
	codec := &JsonRPCCodec{}
	httpServer := gw_net.NewHttpServer(&config.UserServerConfig, lggr.Named("http_server"))

	handlers := make(map[string]Handler)
	donConnMgrs := make(map[string]DONConnectionManager)
	for _, donConfig := range config.Dons {
		donConfig := donConfig
		if donConfig.DonId == "" {
			return nil, errors.New("empty DON ID")
		}
		_, ok := donConnMgrs[donConfig.DonId]
		if ok {
			return nil, fmt.Errorf("duplicate DON ID %s", donConfig.DonId)
		}
		donConnMgr := NewDONConnectionManager(&donConfig, codec)
		donConnMgrs[donConfig.DonId] = donConnMgr
		handler, err := NewHandler(donConfig.HandlerName, &donConfig, donConnMgr)
		if err != nil {
			return nil, err
		}
		handlers[donConfig.DonId] = handler
		donConnMgr.SetHandler(handler)
	}
	connMgr := NewConnectionManager(donConnMgrs, lggr)
	return NewGateway(codec, httpServer, handlers, connMgr, lggr), nil
}

func NewGateway(codec Codec, httpServer gw_net.HttpServer, handlers map[string]Handler, connMgr ConnectionManager, lggr logger.Logger) Gateway {
	gw := &gateway{
		codec:      codec,
		httpServer: httpServer,
		handlers:   handlers,
		connMgr:    connMgr,
		lggr:       lggr.Named("gateway"),
	}
	// Connect servers to the Gateway
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
		return g.httpServer.Start(ctx)
	})
}

func (g *gateway) Close() error {
	return g.StopOnce("Gateway", func() (err error) {
		g.lggr.Info("closing gateway")
		err = multierr.Combine(err, g.httpServer.Close())
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
		return newError(g.codec, "", UserMessageParseError, err.Error())
	}
	// find correct handler
	handler, ok := g.handlers[msg.Body.DonId]
	if !ok {
		return newError(g.codec, msg.Body.MessageId, UnsupportedDONIdError, "unsupported DON ID")
	}
	// send to the handler
	responseCh := make(chan UserCallbackPayload, 1)
	err = handler.HandleUserMessage(ctx, msg, responseCh)
	if err != nil {
		return newError(g.codec, msg.Body.MessageId, InternalHandlerError, err.Error())
	}
	// await response
	var response UserCallbackPayload
	select {
	case <-ctx.Done():
		return newError(g.codec, msg.Body.MessageId, RequestTimeoutError, "handler timeout")
	case response = <-responseCh:
		break
	}
	if response.ErrCode != NoError {
		return newError(g.codec, msg.Body.MessageId, response.ErrCode, response.ErrMsg)
	}
	// encode
	rawResponse, err = g.codec.EncodeResponse(response.Msg)
	if err != nil {
		return newError(g.codec, msg.Body.MessageId, NodeReponseEncodingError, "")
	}
	return rawResponse, ToHttpErrorCode(NoError)
}

func newError(codec Codec, id string, errCode ErrorCode, errMsg string) ([]byte, int) {
	rawResponse, err := codec.EncodeNewErrorResponse(id, ToJsonRPCErrorCode(errCode), errMsg, nil)
	if err != nil {
		// we're not even able to encode a valid JSON response
		return []byte("fatal error"), ToHttpErrorCode(FatalError)
	}
	return rawResponse, ToHttpErrorCode(errCode)
}
