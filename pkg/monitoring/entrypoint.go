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

type Entrypoint struct {
	Context context.Context

	ChainConfig ChainConfig
	Config      config.Config

	Log            Logger
	Producer       Producer
	Metrics        Metrics
	SchemaRegistry SchemaRegistry

	SourceFactories   []SourceFactory
	ExporterFactories []ExporterFactory

	RDDSource Source
	RDDPoller Poller

	Manager Manager

	HTTPServer HTTPServer
}

func NewEntrypoint(
	ctx context.Context,
	log Logger,
	chainConfig ChainConfig,
	envelopeSourceFactory SourceFactory,
	txResultsSourceFactory SourceFactory,
	feedParser FeedParser,
) (*Entrypoint, error) {
	cfg, err := config.Parse()
	if err != nil {
		return nil, fmt.Errorf("failed to parse generic configuration: %w", err)
	}

	if cfg.Feature.TestOnlyFakeReaders {
		envelopeSourceFactory = &fakeRandomDataSourceFactory{make(chan interface{})}
		txResultsSourceFactory = &fakeRandomDataSourceFactory{make(chan interface{})}
	}
	sourceFactories := []SourceFactory{envelopeSourceFactory, txResultsSourceFactory}

	producer, err := NewProducer(ctx, log.With("component", "producer"), cfg.Kafka)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka producer: %w", err)
	}

	metrics := DefaultMetrics

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

	rddSource := NewRDDSource(cfg.Feeds.URL, feedParser)
	if cfg.Feature.TestOnlyFakeRdd {
		// Generate between 2 and 10 random feeds every RDDPollInterval.
		rddSource = NewFakeRDDSource(2, 10)
	}
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
	httpServer := NewHTTPServer(ctx, cfg.HTTP.Address, log.With("component", "http-server"))
	httpServer.Handle("/metrics", metrics.HTTPHandler())
	httpServer.Handle("/debug", manager.HTTPHandler())

	return &Entrypoint{
		ctx,

		chainConfig,
		cfg,

		log,
		producer,
		metrics,
		schemaRegistry,

		sourceFactories,
		exporterFactories,

		rddSource,
		rddPoller,

		manager,

		httpServer,
	}, nil
}

func (e Entrypoint) Run() {
	ctx, cancel := context.WithCancel(e.Context)
	defer cancel()
	wg := &sync.WaitGroup{}

	if e.Config.Feature.TestOnlyFakeReaders {
		envelopeFactory := e.SourceFactories[0].(*fakeRandomDataSourceFactory)
		txResultsFactory := e.SourceFactories[1].(*fakeRandomDataSourceFactory)
		wg.Add(2)
		go func(factory *fakeRandomDataSourceFactory) {
			defer wg.Done()
			factory.RunWithEnvelope(ctx, e.Log)
		}(envelopeFactory)
		go func(factory *fakeRandomDataSourceFactory) {
			defer wg.Done()
			factory.RunWithTxResults(ctx, e.Log)
		}(txResultsFactory)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		e.RDDPoller.Run(ctx)
	}()

	monitor := NewMultiFeedMonitor(
		e.ChainConfig,
		e.Log,
		e.SourceFactories,
		e.ExporterFactories,
		100, // bufferCapacity for source pollers
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		e.Manager.Run(ctx, monitor.Run)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		e.HTTPServer.Run(ctx)
	}()

	// Handle signals from the OS
	wg.Add(1)
	func() {
		defer wg.Done()
		osSignalsCh := make(chan os.Signal, 1)
		signal.Notify(osSignalsCh, syscall.SIGINT, syscall.SIGTERM)
		var sig os.Signal
		select {
		case sig = <-osSignalsCh:
			e.Log.Infow("received signal. Stopping", "signal", sig)
			cancel()
		case <-ctx.Done():
		}
	}()

	wg.Wait()
}
