package promotel

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"

	promotelcommon "github.com/smartcontractkit/chainlink-common/pkg/promotel"
)

const (
	name = "PromOTELForwarder"
)

type Options = promotelcommon.ForwarderOptions
type ForwarderService struct {
	services.StateMachine
	lggr      logger.Logger
	forwarder *promotelcommon.Forwarder
}

func NewForwarderService(g prometheus.Gatherer, r prometheus.Registerer, lggr logger.Logger, opts Options) (*ForwarderService, error) {
	l := logger.Named(lggr, name)
	forwarder, err := promotelcommon.NewForwarder(g, r, l, opts)
	if err != nil {
		return nil, err
	}
	return &ForwarderService{
		lggr:      l,
		forwarder: forwarder,
	}, nil
}

func (f *ForwarderService) HealthReport() map[string]error {
	return map[string]error{f.Name(): f.Healthy()}
}

func (f *ForwarderService) Name() string { return f.lggr.Name() }

func (f *ForwarderService) Start(ctx context.Context) error {
	return f.StartOnce(name, func() error {
		return f.forwarder.Start(ctx)
	})
}

func (f *ForwarderService) Close() error {
	return f.StopOnce(name, f.forwarder.Close)
}

func DefaultOptions() Options {
	return promotelcommon.DefaultForwarderOptions()
}
