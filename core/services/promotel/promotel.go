package promotel

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/timeutil"

	promotelcommon "github.com/smartcontractkit/chainlink-common/pkg/promotel"
)

const (
	period              = 15 * time.Second
	heartbeatMetricName = "promotel_heartbeat"
)

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
	go reportHertbeatMetric(prometheus.DefaultRegisterer, f.lggr, 1*time.Second)

	go f.run()
	return nil
}

func (f *Forwarder) run() {

	defer close(f.done)
	ctx, cancel := f.stopCh.NewCtx()
	defer cancel()
	ticker := timeutil.NewTicker(func() time.Duration { return period })
	defer ticker.Stop()

	go reportHertbeatMetric(prometheus.DefaultRegisterer, f.lggr, 1*time.Second)

	exporter := startExporter(ctx, f.lggr)

	// Fetches metrics from in memory prometheus.Gatherer and converts to OTel format
	receiver := f.startMetricReceiver(func(ctx context.Context, md pmetric.Metrics) error {
		// todo: remove or make configurable
		f.logOtelMetric(md, f.lggr)

		// Exports the converted OTel metric
		return exporter.Consumer().ConsumeMetrics(ctx, md)
	})
	// Close the receiver and exporter when the forwarder is closed
	defer func() {
		if err := receiver.Close(); err != nil {
			f.lggr.Error("Failed to close scraper", zap.Error(err))
		}
		if err := exporter.Close(); err != nil {
			f.lggr.Error("Failed to close exporter", zap.Error(err))
		}
	}()

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

func reportHertbeatMetric(reg prometheus.Registerer, lggr logger.Logger, interval time.Duration) {
	heartbeat := promauto.With(reg).NewCounter(prometheus.CounterOpts{
		Name: heartbeatMetricName,
		ConstLabels: prometheus.Labels{
			"source": "promotel",
		},
	})
	for {
		heartbeat.Inc()
		lggr.Debugw("Heartbeat promotel")
		time.Sleep(interval)
	}
}

func startExporter(ctx context.Context, logger logger.Logger) promotelcommon.MetricExporter {
	expConfig, err := promotelcommon.NewExporterConfig(map[string]any{
		"endpoint": "localhost:4317",
		"tls": map[string]any{
			"insecure": true,
		},
	})
	if err != nil {
		logger.Fatal("Failed to create exporter config", zap.Error(err))
	}
	// Sends metrics data in OTLP format to otel-collector endpoint
	exporter, err := promotelcommon.NewMetricExporter(expConfig, logger)
	if err != nil {
		logger.Fatal("Failed to create metric exporter", zap.Error(err))
	}
	err = exporter.Start(ctx)
	if err != nil {
		logger.Fatal("Failed to start exporter", zap.Error(err))
	}
	return exporter
}

func (f *Forwarder) startMetricReceiver(next consumer.ConsumeMetricsFunc) promotelcommon.Runnable {
	f.lggr.Info("Starting promotel metric receiver")
	config, err := promotelcommon.NewDefaultReceiverConfig()
	if err != nil {
		f.lggr.Fatal("Failed to create config", zap.Error(err))
	}

	// Gather metrics via promotel
	// MetricReceiver fetches metrics from prometheus.Gatherer, then converts it to OTel format and writes formatted metrics to stdout
	receiver, err := promotelcommon.NewMetricReceiver(config, f.gatherer, next, f.lggr)
	if err != nil {
		f.lggr.Fatal("Failed to create debug metric receiver", zap.Error(err))
	}
	// Starts the promotel
	if err := receiver.Start(context.Background()); err != nil {
		f.lggr.Fatal("Failed to start metric receiver", zap.Error(err))
	}
	return receiver
}

func (f *Forwarder) logOtelMetric(md pmetric.Metrics, logger logger.Logger) {
	rms := md.ResourceMetrics()
	for i := 0; i < rms.Len(); i++ {
		rm := rms.At(i)
		ilms := rm.ScopeMetrics()
		for j := 0; j < ilms.Len(); j++ {
			ilm := ilms.At(j)
			metrics := ilm.Metrics()
			for k := 0; k < metrics.Len(); k++ {
				metric := metrics.At(k)
				logger.Debug("Exporting OTel metric ", zap.Any("name", metric.Name()), zap.Any("value", metric.Sum().DataPoints().At(0).DoubleValue()))
			}
		}
	}
}
