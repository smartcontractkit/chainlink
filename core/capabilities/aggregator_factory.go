package capabilities

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/datafeeds"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/streams"
)

func NewAggregator(name string, config values.Map, lggr logger.Logger) (types.Aggregator, error) {
	switch name {
	case "data_feeds":
		mc := streams.NewCodec(lggr)
		return datafeeds.NewDataFeedsAggregator(config, mc, lggr)
	default:
		return nil, fmt.Errorf("aggregator %s not supported", name)
	}
}
