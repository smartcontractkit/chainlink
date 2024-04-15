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

	// iterate over each feed
	for _, feedConfig := range data.Feeds {
		feedLogger := logger.With(m.log,
			"feed_name", feedConfig.GetName(),
			"feed_id", feedConfig.GetID(),
			"network", m.chainConfig.GetNetworkName(),
		)

		// create sources outside of createMonitor
		sources := []Source{}
		sourceTypes := []string{}
		for _, sourceFactory := range m.sourceFactories {
			source, err := sourceFactory.NewSource(m.chainConfig, feedConfig)
			if err != nil {
				feedLogger.Errorw("failed to create source", "error", err, "source-type", fmt.Sprintf("%T", sourceFactory))
				continue
			}
			sources = append(sources, source)
			sourceTypes = append(sourceTypes, sourceFactory.GetType())
		}

		runMonitor(
			ctx,
			&subs,
			feedLogger,
			m.chainConfig,
			m.bufferCapacity,
			sources,
			sourceTypes,
			m.exporterFactories,
			ExporterParams{
				m.chainConfig,
				feedConfig,
				data.Nodes,
			},
		)
	}
}
