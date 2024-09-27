package promotel

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/timeutil"
)

const period = 15 * time.Second

type Forwarder struct {
	services.StateMachine
	lggr          logger.Logger
	gatherer      prometheus.Gatherer
	meterProvider metric.MeterProvider
	stopCh        services.StopChan
	done          chan struct{}
}

func NewForwarder(lggr logger.Logger, gatherer prometheus.Gatherer, meterProvider metric.MeterProvider) *Forwarder {
	return &Forwarder{
		lggr:          logger.Named(lggr, "PromOTELForwarder"),
		gatherer:      gatherer,
		meterProvider: meterProvider,
		stopCh:        make(chan struct{}),
		done:          make(chan struct{}),
	}
}

func (f *Forwarder) HealthReport() map[string]error { return map[string]error{f.Name(): f.Healthy()} }

func (f *Forwarder) Name() string { return f.lggr.Name() }

func (f *Forwarder) Start(context.Context) error {
	go f.run()
	return nil
}

func (f *Forwarder) run() {
	defer close(f.done)
	ctx, cancel := f.stopCh.NewCtx()
	defer cancel()
	ticker := timeutil.NewTicker(func() time.Duration { return period })
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			f.forward(ctx)
		}
	}
}

func (f *Forwarder) forward(ctx context.Context) {
	mfs, err := f.gatherer.Gather()
	if err != nil {
		f.lggr.Errorw("Failed to gather prometheus metrics", "err", err)
	}
	for _, mf := range mfs {
		for range mf.Metric {
			if ctx.Err() != nil {
				return
			}

			//TODO f.meterProvider.Meter()
		}
	}
}

func (f *Forwarder) Close() error {
	close(f.stopCh)
	<-f.done
	return nil
}
