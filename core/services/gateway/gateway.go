package gateway

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type Gateway interface {
	job.ServiceCtx
}

type gateway struct {
	utils.StartStopOnce

	codec    Codec
	handlers map[string]Handler
	connMgr  ConnectionManager
	lggr     logger.Logger
}

func NewGatewayFromConfig(config *GatewayConfig, lggr logger.Logger) (Gateway, error) {
	codec := &JsonRPCCodec{}

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
	return NewGateway(codec, handlers, connMgr, lggr), nil
}

func NewGateway(codec Codec, handlers map[string]Handler, connMgr ConnectionManager, lggr logger.Logger) Gateway {
	return &gateway{
		codec:    codec,
		handlers: handlers,
		connMgr:  connMgr,
		lggr:     lggr.Named("gateway"),
	}
}

func (g *gateway) Start(ctx context.Context) error {
	return g.StartOnce("Gateway", func() error {
		g.lggr.Info("starting gateway")
		for _, handler := range g.handlers {
			if err := handler.Start(ctx); err != nil {
				return err
			}
		}
		return nil
	})
}

func (g *gateway) Close() error {
	return g.StopOnce("Gateway", func() (err error) {
		g.lggr.Info("closing gateway")
		for _, handler := range g.handlers {
			err = multierr.Combine(err, handler.Close())
		}
		return
	})
}
