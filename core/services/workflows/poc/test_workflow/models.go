package test_workflow

import (
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/poc/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/poc/workflow"
)

// These would be generated from protos provided by the capability author
// For now, they mimic what's in the test for the multi input workflow

type MercuryTriggerResponse struct {
	Values   map[string]decimal.Decimal
	Decimals map[string]int
	Metadata TriggerMetadata
}

func NewMercuryTrigger(ref, typeName string) *capabilities.RemoteTrigger[*MercuryTriggerResponse] {
	return &capabilities.RemoteTrigger[*MercuryTriggerResponse]{
		// TODO this would be what we can use to distinguish between different triggers
		// to allow data normalization
		RefName:  ref,
		TypeName: typeName,
	}
}

type TriggerMetadata struct {
	TriggerRef string
}

// Not realistic for a chain writer, just mimicking the test

type ChainWriter interface {
	AddWriteTarget(ref string, wb *workflow.Builder[capabilities.ConsensusResult[*MercuryTriggerResponse]]) error
}

type ChainWriteRequest struct{}

func NewChainWriter(typeName string) ChainWriter {
	return &chainWriter{typeName: typeName}
}

type chainWriter struct {
	typeName string
}

func (c *chainWriter) AddWriteTarget(ref string, wb *workflow.Builder[capabilities.ConsensusResult[*MercuryTriggerResponse]]) error {
	return workflow.AddTarget[*MercuryTriggerResponse](wb, &capabilities.RemoteTarget[*MercuryTriggerResponse]{
		RefName:  ref,
		TypeName: c.typeName,
	})
}

var _ ChainWriter = (*chainWriter)(nil)
