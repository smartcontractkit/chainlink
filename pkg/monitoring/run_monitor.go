package monitoring

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"
)

// runMonitor is a reusable method for creating pollers, exporters, and monitors
// sources are created outside because they may require different input parameters
func runMonitor(
	ctx context.Context,
	subs *utils.Subprocesses,
	lgr Logger,
	chainConfig ChainConfig,
	bufferCapacity uint32,
	sources []Source,
	sourceTypes []string,
	exporterFactories []ExporterFactory,
	exporterParams ExporterParams,
) {
	// Create data sources
	pollers := []Poller{}
	for i, source := range sources {
		poller := NewSourcePoller(
			source,
			logger.With(lgr, "component", "chain-poller", "source", sourceTypes[i]),
			chainConfig.GetPollInterval(),
			chainConfig.GetReadTimeout(),
			bufferCapacity,
		)
		pollers = append(pollers, poller)
	}
	if len(pollers) == 0 {
		lgr.Errorw("not tracking feed because all sources failed to initialize")
		return
	}
	// Create exporters
	exporters := []Exporter{}
	for _, exporterFactory := range exporterFactories {
		exporter, err := exporterFactory.NewExporter(exporterParams)
		if err != nil {
			lgr.Errorw("failed to create new exporter", "error", err, "exporter-type", fmt.Sprintf("%T", exporterFactory))
			continue
		}
		exporters = append(exporters, exporter)
	}
	if len(exporters) == 0 {
		lgr.Errorw("not tracking feed because all exporters failed to initialize")
		return
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
		logger.With(lgr, "component", "feed-monitor"),
		pollers,
		exporters,
	)
	subs.Go(func() {
		feedMonitor.Run(ctx)
	})
}
