package monitoring

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"
)

// MultiFeedMonitor manages the flow of data from multiple sources to
// multiple exporters for each feed in the configuration.
type MultiFeedMonitor interface {
	Run(ctx context.Context, data RDDData)
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
func (m *multiFeedMonitor) Run(ctx context.Context, data RDDData) {
	var subs utils.Subprocesses
	defer subs.Wait()

FEED_LOOP:
	for _, feedConfig := range data.Feeds {
		feedLogger := logger.With(m.log,
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
				logger.With(m.log, "component", "chain-poller", "source", sourceFactory.GetType()),
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
			exporter, err := exporterFactory.NewExporter(ExporterParams{
				m.chainConfig,
				feedConfig,
				data.Nodes,
			})
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
		for _, poller := range pollers {
			poller := poller
			subs.Go(func() {
				poller.Run(ctx)
			})
		}
		// Run feed monitor.
		feedMonitor := NewFeedMonitor(
			logger.With(m.log, "component", "feed-monitor"),
			pollers,
			exporters,
		)
		subs.Go(func() {
			feedMonitor.Run(ctx)
		})
	}
}
