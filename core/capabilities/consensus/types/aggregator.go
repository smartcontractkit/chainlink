package types

import (
	ocrcommon "github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

type Aggregator interface {
	// Called by the Outcome() phase of OCR reporting.
	// The inner array of observations corresponds to elements listed in "inputs.observations" section.
	Aggregate(previousOutcome *AggregationOutcome, observations map[ocrcommon.OracleID][]values.Value) (*AggregationOutcome, error)
}

// TODO move to a factory object
//func NewAggregator(aggregationMethod string, aggregationConfig values.Map) (Aggregator, error) {
//	if aggregationMethod == "data_feeds_2_0" {
//		return datafeeds.NewDataFeedsAggregator(aggregationConfig)
//	} else {
//		return nil, fmt.Errorf("unknown aggregation method %s", aggregationMethod)
//	}
//}
