package monitoring

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/smartcontractkit/chainlink-relay/pkg/monitoring/config"
)

type Monitor struct {
	RootContext context.Context

	ChainConfig ChainConfig
	Config      config.Config

	Log            Logger
	Producer       Producer
	Metrics        Metrics
	ChainMetrics   ChainMetrics
	SchemaRegistry SchemaRegistry

	SourceFactories   []SourceFactory
	ExporterFactories []ExporterFactory

	RDDSource Source
	RDDPoller Poller

	Manager Manager

	HTTPServer HTTPServer
}

func NewMonitor(
	rootCtx context.Context,
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

	metrics := NewMetrics(log.With("component", "metrics"))
	chainMetrics := NewChainMetrics(chainConfig)

	sourceFactories := []SourceFactory{envelopeSourceFactory, txResultsSourceFactory}

	producer, err := NewProducer(rootCtx, log.With("component", "producer"), cfg.Kafka)
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
		log.With("component", "prometheus-exporter"),
		metrics,
	)
	kafkaExporterFactory, err := NewKafkaExporterFactory(
		log.With("component", "kafka-exporter"),
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
		log.With("component", "rdd-source"),
	)

	rddPoller := NewSourcePoller(
		rddSource,
		log.With("component", "rdd-poller"),
		cfg.Feeds.RDDPollInterval,
		cfg.Feeds.RDDReadTimeout,
		0, // no buffering!
	)

	manager := NewManager(
		log.With("component", "manager"),
		rddPoller,
	)

	// Configure HTTP server
	httpServer := NewHTTPServer(rootCtx, cfg.HTTP.Address, log.With("component", "http-server"))
	httpServer.Handle("/metrics", metrics.HTTPHandler())
	httpServer.Handle("/debug", manager.HTTPHandler())

	return &Monitor{
		rootCtx,

		chainConfig,
		cfg,

		log,
		producer,
		metrics,
		chainMetrics,
		schemaRegistry,

		sourceFactories,
		exporterFactories,

		rddSource,
		rddPoller,

		manager,

		httpServer,
	}, nil
}

func (m Monitor) Run() {
	rootCtx, cancel := context.WithCancel(m.RootContext)
	defer cancel()
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		m.RDDPoller.Run(rootCtx)
	}()

	// Instrument all source factories
	instrumentedSourceFactories := []SourceFactory{}
	for _, factory := range m.SourceFactories {
		instrumentedSourceFactories = append(instrumentedSourceFactories,
			NewInstrumentedSourceFactory(factory, m.ChainMetrics))
	}

	monitor := NewMultiFeedMonitor(
		m.ChainConfig,
		m.Log,
		instrumentedSourceFactories,
		m.ExporterFactories,
		100, // bufferCapacity for source pollers
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		m.Manager.Run(rootCtx, func(localCtx context.Context, data RDDData) {
			m.ChainMetrics.SetNewFeedConfigsDetected(float64(len(data.Feeds)))
			monitor.Run(localCtx, data)
		})
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		m.HTTPServer.Run(rootCtx)
	}()

	// Handle signals from the OS
	wg.Add(1)
	go func() {
		defer wg.Done()
		osSignalsCh := make(chan os.Signal, 1)
		signal.Notify(osSignalsCh, syscall.SIGINT, syscall.SIGTERM)
		var sig os.Signal
		select {
		case sig = <-osSignalsCh:
			m.Log.Infow("received signal. Stopping", "signal", sig)
			cancel()
		case <-rootCtx.Done():
		}
	}()

	wg.Wait()
}
