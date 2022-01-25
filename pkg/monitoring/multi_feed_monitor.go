package monitoring

import (
	"context"
	"fmt"
	"sync"
)

type MultiFeedMonitor interface {
	Run(ctx context.Context, feeds []FeedConfig)
}

func NewMultiFeedMonitor(
	chainConfig ChainConfig,
	log Logger,

	sourceFactories []SourceFactory,
	exporterFactories []ExporterFactory,
) MultiFeedMonitor {
	return &multiFeedMonitor{
		chainConfig,
		log,

		sourceFactories,
		exporterFactories,
	}
}

type multiFeedMonitor struct {
	chainConfig ChainConfig

	log               Logger
	sourceFactories   []SourceFactory
	exporterFactories []ExporterFactory
}

const bufferCapacity = 100

// Run should be executed as a goroutine.
func (m *multiFeedMonitor) Run(ctx context.Context, feeds []FeedConfig) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

FEED_LOOP:
	for _, feedConfig := range feeds {
		feedLogger := m.log.With(
			"feed", feedConfig.GetName(),
			"network", m.chainConfig.GetNetworkName(),
		)
		// Create data sources
		pollers := make([]Poller, len(m.sourceFactories))
		for i, sourceFactory := range m.sourceFactories {
			source, err := sourceFactory.NewSource(m.chainConfig, feedConfig)
			if err != nil {
				feedLogger.Errorw("failed to create new source", "error", err, "source-type", fmt.Sprintf("%T", sourceFactory))
				continue FEED_LOOP
			}
			poller := NewSourcePoller(
				source,
				feedLogger.With("component", "chain-poller"),
				m.chainConfig.GetPollInterval(),
				m.chainConfig.GetReadTimeout(),
				bufferCapacity,
			)
			pollers[i] = poller
		}
		// Create exporters
		exporters := make([]Exporter, len(m.exporterFactories))
		for i, exporterFactory := range m.exporterFactories {
			exporter, err := exporterFactory.NewExporter(m.chainConfig, feedConfig)
			if err != nil {
				feedLogger.Errorw("failed to create new exporter", "error", err, "exporter-type", fmt.Sprintf("%T", exporterFactory))
				continue FEED_LOOP
			}
			exporters[i] = exporter
		}
		// Run poller goroutines.
		wg.Add(len(pollers))
		for _, poller := range pollers {
			go func(poller Poller) {
				defer wg.Done()
				poller.Run(ctx)
			}(poller)
		}
		// Run feed monitor.
		feedMonitor := NewFeedMonitor(
			feedLogger.With("component", "feed-monitor"),
			pollers,
			exporters,
		)
		wg.Add(1)
		go func() {
			defer wg.Done()
			feedMonitor.Run(ctx)
		}()
	}
}
