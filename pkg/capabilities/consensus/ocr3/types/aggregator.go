package types

import (
	ocrcommon "github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

const (
	WorkflowIDFieldName  = "INTERNAL_workflow_id"
	ExecutionIDFieldName = "INTERNAL_execution_id"
)

type Aggregator interface {
	// Called by the Outcome() phase of OCR reporting.
	// The inner array of observations corresponds to elements listed in "inputs.observations" section.
	Aggregate(previousOutcome *AggregationOutcome, observations map[ocrcommon.OracleID][]values.Value) (*AggregationOutcome, error)
}

func AppendWorkflowIDs(outcome *AggregationOutcome, workflowID string, workflowExecutionID string) (*AggregationOutcome, error) {
	valueWID, err := values.Wrap(workflowID)
	if err != nil {
		return nil, err
	}
	protoWID, err := valueWID.Proto()
	if err != nil {
		return nil, err
	}
	outcome.EncodableOutcome.Fields[WorkflowIDFieldName] = protoWID
	valueWEID, err := values.Wrap(workflowExecutionID)
	if err != nil {
		return nil, err
	}
	protoWEID, err := valueWEID.Proto()
	if err != nil {
		return nil, err
	}
	outcome.EncodableOutcome.Fields[ExecutionIDFieldName] = protoWEID
	return outcome, nil
}
