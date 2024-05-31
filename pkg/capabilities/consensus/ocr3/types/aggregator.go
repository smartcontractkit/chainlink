package types

import (
	ocrcommon "github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

const MetadataFieldName = "INTERNAL_METADATA"

type Metadata struct {
	Version       uint32 // 1 byte
	ExecutionID   string // 32 hex bytes
	Timestamp     uint32 //  4 bytes
	DONID         string //  4 hex bytes
	WorkflowID    string // 32 hex bytes
	WorkflowName  string // 10 hex bytes
	WorkflowOwner string // 20 hex bytes
	ReportID      string //  2 hex bytes
}

type Aggregator interface {
	// Called by the Outcome() phase of OCR reporting.
	// The inner array of observations corresponds to elements listed in "inputs.observations" section.
	Aggregate(previousOutcome *AggregationOutcome, observations map[ocrcommon.OracleID][]values.Value, f int) (*AggregationOutcome, error)
}

func AppendMetadata(outcome *AggregationOutcome, meta *Metadata) (*AggregationOutcome, error) {
	metaWrapped, err := values.Wrap(meta)
	if err != nil {
		return nil, err
	}
	outcome.EncodableOutcome.Fields[MetadataFieldName] = values.Proto(metaWrapped)
	return outcome, nil
}

type AggregatorFactory func(name string, config values.Map, lggr logger.Logger) (Aggregator, error)
