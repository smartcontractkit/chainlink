package monitoring

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/monitoring/config"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"
)

// Monitor is the entrypoint for an on-chain monitor integration.
// Monitors should only be created via NewMonitor()
type Monitor struct {
	StopCh services.StopRChan

	ChainConfig ChainConfig
	Config      config.Config

	Log            Logger
	Producer       Producer
	Metrics        Metrics
	ChainMetrics   ChainMetrics
	SchemaRegistry SchemaRegistry

	// per-feed
	SourceFactories   []SourceFactory
	ExporterFactories []ExporterFactory

	// single (network level, default empty)
	NetworkSourceFactories   []NetworkSourceFactory
	NetworkExporterFactories []ExporterFactory

	RDDSource Source
	RDDPoller Poller

	Manager Manager

	HTTPServer HTTPServer
}

// NewMonitor builds a new Monitor instance using dependency injection.
// If advanced configurations of the Monitor are required - for instance,
// adding a custom third party service to send data to - this method
// should provide a good starting template to do that.
func NewMonitor(
	stopCh services.StopRChan,
	log Logger,
	chainConfig ChainConfig,
	envelopeSourceFactory SourceFactory,
	txResultsSourceFactory SourceFactory,
	feedsParser FeedsParser,
	nodesParser NodesParser,
) (*Monitor, error) {
	cfg, err := config.Parse()
	if err != nil {
		return nil, fmt.Errorf("failed to parse generic configuration: %w", err)
	}

	metrics := NewMetrics(logger.With(log, "component", "metrics"))
	chainMetrics := NewChainMetrics(chainConfig)

	sourceFactories := []SourceFactory{envelopeSourceFactory, txResultsSourceFactory}

	producer, err := NewProducer(stopCh, logger.With(log, "component", "producer"), cfg.Kafka)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka producer: %w", err)
	}
	producer = NewInstrumentedProducer(producer, chainMetrics)

	schemaRegistry := NewSchemaRegistry(cfg.SchemaRegistry, log)

	transmissionSchema, err := schemaRegistry.EnsureSchema(
		SubjectFromTopic(cfg.Kafka.TransmissionTopic), TransmissionAvroSchema)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare transmission schema: %w", err)
	}
	configSetSimplifiedSchema, err := schemaRegistry.EnsureSchema(
		SubjectFromTopic(cfg.Kafka.ConfigSetSimplifiedTopic), ConfigSetSimplifiedAvroSchema)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare config_set_simplified schema: %w", err)
	}

	prometheusExporterFactory := NewPrometheusExporterFactory(
		logger.With(log, "component", "prometheus-exporter"),
		metrics,
	)
	kafkaExporterFactory, err := NewKafkaExporterFactory(
		logger.With(log, "component", "kafka-exporter"),
		producer,
		[]Pipeline{
			{cfg.Kafka.TransmissionTopic, MakeTransmissionMapping, transmissionSchema},
			{cfg.Kafka.ConfigSetSimplifiedTopic, MakeConfigSetSimplifiedMapping, configSetSimplifiedSchema},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka exporter: %w", err)
	}

	exporterFactories := []ExporterFactory{prometheusExporterFactory, kafkaExporterFactory}

	rddSource := NewRDDSource(
		cfg.Feeds.URL, feedsParser, cfg.Feeds.IgnoreIDs,
		cfg.Nodes.URL, nodesParser,
		logger.With(log, "component", "rdd-source"),
	)

	rddPoller := NewSourcePoller(
		rddSource,
		logger.With(log, "component", "rdd-poller"),
		cfg.Feeds.RDDPollInterval,
		cfg.Feeds.RDDReadTimeout,
		0, // no buffering!
	)

	manager := NewManager(
		logger.With(log, "component", "manager"),
		rddPoller,
	)

	// Configure HTTP server
	httpServer := NewHTTPServer(stopCh, cfg.HTTP.Address, logger.With(log, "component", "http-server"))
	httpServer.Handle("/metrics", metrics.HTTPHandler())
	httpServer.Handle("/debug", manager.HTTPHandler())
	// Required for k8s.
	httpServer.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	return &Monitor{
		StopCh:            stopCh,
		ChainConfig:       chainConfig,
		Config:            cfg,
		Log:               log,
		Producer:          producer,
		Metrics:           metrics,
		ChainMetrics:      chainMetrics,
		SchemaRegistry:    schemaRegistry,
		SourceFactories:   sourceFactories,
		ExporterFactories: exporterFactories,
		RDDSource:         rddSource,
		RDDPoller:         rddPoller,
		Manager:           manager,
		HTTPServer:        httpServer,
	}, nil
}

// Run() starts all the goroutines needed by a Monitor. The lifecycle of these routines
// is controlled by the context passed to the NewMonitor constructor.
func (m Monitor) Run() {
	rootCtx, cancel := m.StopCh.NewCtx()
	defer cancel()
	var subs utils.Subprocesses

	subs.Go(func() {
		m.RDDPoller.Run(rootCtx)
	})

	// Instrument all source factories
	instrumentedSourceFactories := []SourceFactory{}
	for _, factory := range m.SourceFactories {
		instrumentedSourceFactories = append(instrumentedSourceFactories,
			NewInstrumentedSourceFactory(factory, m.ChainMetrics))
	}

	// setup per-feed & network monitor
	monitor := NewMultiFeedMonitor(
		m.ChainConfig,
		m.Log,
		instrumentedSourceFactories,
		m.ExporterFactories,
		100, // bufferCapacity for source pollers
	)
	networkMonitor := NewNetworkMonitor(
		m.ChainConfig,
		m.Log,
		m.NetworkSourceFactories,
		m.NetworkExporterFactories,
		100, // bufferCapacity for source pollers
	)
	subs.Go(func() {
		m.Manager.Run(rootCtx,
			// run per-feed monitors
			func(localCtx context.Context, data RDDData) {
				m.ChainMetrics.SetNewFeedConfigsDetected(float64(len(data.Feeds)))
				m.Log.Infow("Starting Feed Monitor", "exporters", len(m.ExporterFactories), "sources", len(instrumentedSourceFactories), "feeds", len(data.Feeds))
				monitor.Run(localCtx, data) // blocking func controlled by ctx
			},
			// run network monitor if factories present
			func(localCtx context.Context, data RDDData) {
				if len(m.NetworkExporterFactories) != 0 || len(m.NetworkSourceFactories) != 0 {
					m.Log.Infow("Starting Network Monitor", "exporters", len(m.NetworkExporterFactories), "sources", len(m.NetworkSourceFactories))
					networkMonitor.Run(localCtx, data) // blocking func controlled by ctx
				}
			},
		)
	})

	subs.Go(func() {
		m.HTTPServer.Run(rootCtx)
	})

	// Handle signals from the OS
	subs.Go(func() {
		osSignalsCh := make(chan os.Signal, 1)
		signal.Notify(osSignalsCh, syscall.SIGINT, syscall.SIGTERM)
		var sig os.Signal
		select {
		case sig = <-osSignalsCh:
			m.Log.Infow("received signal. Stopping", "signal", sig)
			cancel()
		case <-rootCtx.Done():
		}
	})

	subs.Wait()
}
