package monitoring

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"
)

// NetworkMOnitor manages the flow of data from sources to exporters
// for non-feed-specific metrics (blockheight, balances, etc)
type NetworkMonitor interface {
	Run(ctx context.Context, data RDDData)
}

func NewNetworkMonitor(
	chainConfig ChainConfig,
	log Logger,

	sourceFactories []NetworkSourceFactory,
	exporterFactories []ExporterFactory,

	bufferCapacity uint32,
) NetworkMonitor {
	return &networkMonitor{
		chainConfig,
		log,

		sourceFactories,
		exporterFactories,

		bufferCapacity,
	}
}

type networkMonitor struct {
	chainConfig ChainConfig

	log               Logger
	sourceFactories   []NetworkSourceFactory
	exporterFactories []ExporterFactory

	bufferCapacity uint32
}

// Run should be executed as a goroutine.
func (m *networkMonitor) Run(ctx context.Context, data RDDData) {
	var subs utils.Subprocesses
	defer subs.Wait()

	lgr := logger.With(m.log,
		"network", m.chainConfig.GetNetworkName(),
	)

	sources := []Source{}
	sourceTypes := []string{}
	for _, sourceFactory := range m.sourceFactories {
		source, err := sourceFactory.NewSource(m.chainConfig, data.Nodes)
		if err != nil {
			lgr.Errorw("failed to create source", "error", err, "source-type", fmt.Sprintf("%T", sourceFactory))
			continue
		}
		sources = append(sources, source)
		sourceTypes = append(sourceTypes, sourceFactory.GetType())
	}

	runMonitor(
		ctx,
		&subs,
		lgr,
		m.chainConfig,
		m.bufferCapacity,
		sources,
		sourceTypes,
		m.exporterFactories,
		ExporterParams{
			ChainConfig: m.chainConfig,
			Nodes:       data.Nodes,
		},
	)
}
