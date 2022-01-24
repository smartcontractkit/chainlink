package monitoring

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/smartcontractkit/chainlink-relay/pkg/monitoring/config"
)

// Facade is the entrypoint to the monitoring service.
// All arguments are required!
// To terminate, cancel the context and wait for Facade to exit.
func Facade(
	ctx context.Context,
	log Logger,
	chainConfig ChainConfig,
	sourceFactory SourceFactory,
	feedParser FeedParser,
) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	bgCtx, cancelBgCtx := context.WithCancel(ctx)
	defer cancelBgCtx()

	cfg, err := config.Parse()
	if err != nil {
		log.Fatalw("failed to parse generic configuration", "error", err)
	}

	schemaRegistry := NewSchemaRegistry(cfg.SchemaRegistry, log)

	transmissionSchema, err := schemaRegistry.EnsureSchema(cfg.Kafka.TransmissionTopic+"-value", TransmissionAvroSchema)
	if err != nil {
		log.Fatalw("failed to prepare transmission schema", "error", err)
	}
	configSetSimplifiedSchema, err := schemaRegistry.EnsureSchema(cfg.Kafka.ConfigSetSimplifiedTopic+"-value", ConfigSetSimplifiedAvroSchema)
	if err != nil {
		log.Fatalw("failed to prepare config_set_simplified schema", "error", err)
	}

	producer, err := NewProducer(bgCtx, log.With("component", "producer"), cfg.Kafka)
	if err != nil {
		log.Fatalw("failed to create kafka producer", "error", err)
	}

	if cfg.Feature.TestOnlyFakeReaders {
		sourceFactory = NewRandomDataSourceFactory(bgCtx, wg, log.With("component", "rand-source"))
	}

	metrics := DefaultMetrics

	monitor := NewMultiFeedMonitor(
		chainConfig,

		log,
		sourceFactory,
		producer,
		metrics,

		cfg.Kafka.ConfigSetSimplifiedTopic,
		cfg.Kafka.TransmissionTopic,

		transmissionSchema,
		configSetSimplifiedSchema,
	)

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
	wg.Add(1)
	go func() {
		defer wg.Done()
		rddPoller.Run(bgCtx)
	}()

	manager := NewManager(
		log.With("component", "manager"),
		rddPoller,
	)
	wg.Add(1)
	go func() {
		defer wg.Done()
		manager.Run(bgCtx, monitor.Run)
	}()

	// Configure HTTP server
	http := NewHttpServer(bgCtx, cfg.Http.Address, log.With("component", "http-server"))
	http.Handle("/metrics", metrics.HTTPHandler())
	http.Handle("/debug", manager.HTTPHandler())
	wg.Add(1)
	go func() {
		defer wg.Done()
		http.Run(bgCtx)
	}()

	// Handle signals from the OS
	osSignalsCh := make(chan os.Signal, 1)
	signal.Notify(osSignalsCh, syscall.SIGINT, syscall.SIGTERM)
	sig := <-osSignalsCh
	log.Infow("received signal. Stopping", "signal", sig)
}
