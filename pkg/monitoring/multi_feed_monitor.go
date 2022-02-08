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

	bufferCapacity uint32,
) MultiFeedMonitor {
	return &multiFeedMonitor{
		chainConfig,
		log,

		sourceFactories,
		exporterFactories,

		bufferCapacity,
	}
}

type multiFeedMonitor struct {
	chainConfig ChainConfig

	log               Logger
	sourceFactories   []SourceFactory
	exporterFactories []ExporterFactory

	bufferCapacity uint32
}

// Run should be executed as a goroutine.
func (m *multiFeedMonitor) Run(ctx context.Context, feeds []FeedConfig) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

FEED_LOOP:
	for _, feedConfig := range feeds {
		feedLogger := m.log.With(
			"feed_name", feedConfig.GetName(),
			"feed_id", feedConfig.GetID(),
			"network", m.chainConfig.GetNetworkName(),
		)
		// Create data sources
		pollers := []Poller{}
		for _, sourceFactory := range m.sourceFactories {
			source, err := sourceFactory.NewSource(m.chainConfig, feedConfig)
			if err != nil {
				feedLogger.Errorw("failed to create source", "error", err, "source-type", fmt.Sprintf("%T", sourceFactory))
				continue
			}
			poller := NewSourcePoller(
				source,
				feedLogger.With("component", "chain-poller"),
				m.chainConfig.GetPollInterval(),
				m.chainConfig.GetReadTimeout(),
				m.bufferCapacity,
			)
			pollers = append(pollers, poller)
		}
		if len(pollers) == 0 {
			feedLogger.Errorw("not tracking feed because all sources failed to initialize")
			continue FEED_LOOP
		}
		// Create exporters
		exporters := []Exporter{}
		for _, exporterFactory := range m.exporterFactories {
			exporter, err := exporterFactory.NewExporter(m.chainConfig, feedConfig)
			if err != nil {
				feedLogger.Errorw("failed to create new exporter", "error", err, "exporter-type", fmt.Sprintf("%T", exporterFactory))
				continue
			}
			exporters = append(exporters, exporter)
		}
		if len(exporters) == 0 {
			feedLogger.Errorw("not tracking feed because all exporters failed to initialize")
			continue FEED_LOOP
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
