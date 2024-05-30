package types

import (
	ocrcommon "github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

const (
	WorkflowIDFieldName    = "INTERNAL_workflow_id"
	DonIDFieldName         = "INTERNAL_don_id"
	ExecutionIDFieldName   = "INTERNAL_execution_id"
	WorkflowOwnerFieldName = "INTERNAL_workflow_owner"
)

type Aggregator interface {
	// Called by the Outcome() phase of OCR reporting.
	// The inner array of observations corresponds to elements listed in "inputs.observations" section.
	Aggregate(previousOutcome *AggregationOutcome, observations map[ocrcommon.OracleID][]values.Value, f int) (*AggregationOutcome, error)
}

func AppendWorkflowIDs(outcome *AggregationOutcome, workflowID string, donID string, workflowExecutionID string, workflowOwner string) (*AggregationOutcome, error) {
	valueWID, err := values.Wrap(workflowID)
	if err != nil {
		return nil, err
	}
	outcome.EncodableOutcome.Fields[WorkflowIDFieldName] = values.Proto(valueWID)
	valueDID, err := values.Wrap(donID)
	if err != nil {
		return nil, err
	}
	outcome.EncodableOutcome.Fields[DonIDFieldName] = values.Proto(valueDID)
	valueWEID, err := values.Wrap(workflowExecutionID)
	if err != nil {
		return nil, err
	}
	outcome.EncodableOutcome.Fields[ExecutionIDFieldName] = values.Proto(valueWEID)
	valueWOwner, err := values.Wrap(workflowOwner)
	if err != nil {
		return nil, err
	}
	outcome.EncodableOutcome.Fields[WorkflowOwnerFieldName] = values.Proto(valueWOwner)
	return outcome, nil
}

type AggregatorFactory func(name string, config values.Map, lggr logger.Logger) (Aggregator, error)
