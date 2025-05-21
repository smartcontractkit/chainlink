package capabilities

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/aggregators"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/datafeeds"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/streams"
)

type Aggregator string

const (
	// DataFeedsAggregator is the name of the data feeds aggregator.
	DataFeedsAggregator Aggregator = "data_feeds"
	// IdenticalAggregator is the name of the identical aggregator.
	IdenticalAggregator Aggregator = "identical"
	// ReduceAggregator is the name of the reduce aggregator.
	ReduceAggregator Aggregator = "reduce"
	// StreamsAggregator is the name of the streams aggregator.
	LLOStreamsAggregator Aggregator = "llo_streams"
)

// NewAggregator creates a new aggregator based on the provided name and config.
func NewAggregator(name string, config values.Map, lggr logger.Logger) (types.Aggregator, error) {
	switch name {
	case string(DataFeedsAggregator):
		mc := streams.NewCodec(lggr)
		return datafeeds.NewDataFeedsAggregator(config, mc)
	case string(IdenticalAggregator):
		return aggregators.NewIdenticalAggregator(config)
	case string(ReduceAggregator):
		return aggregators.NewReduceAggregator(config)
	case string(LLOStreamsAggregator):
		return datafeeds.NewLLOAggregator(config)
	default:
		return nil, fmt.Errorf("aggregator %s not supported", name)
	}
}
