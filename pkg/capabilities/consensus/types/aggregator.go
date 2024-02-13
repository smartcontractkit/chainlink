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
