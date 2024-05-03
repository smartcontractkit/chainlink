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

func NewMercuryTrigger(ref string) *capabilities.RemoteTrigger[*MercuryTriggerResponse] {
	return &capabilities.RemoteTrigger[*MercuryTriggerResponse]{
		// TODO this would be what we can use to distinguish between different triggers
		// to allow data normalization
		RefName:  ref,
		TypeName: "mercury-trigger",
	}
}

type TriggerMetadata struct {
	TriggerRef string
}

type ChainReader interface {
	AddReadAction(ref string, wb *workflow.Builder[*MercuryTriggerResponse]) (workflow.Builder[*ChainReadResponse], error)
}

// This is where the concept of capability sets could come into play
// we would need to be able to generate the full name form the set, then the set would be mapped to the internal name and not each capability
// see the standalone POC.

func NewChainReader(typeName string) ChainReader {
	return &chainReader{typeName: typeName}
}

type chainReader struct {
	typeName string
}

func (c *chainReader) AddReadAction(ref string, wb *workflow.Builder[*MercuryTriggerResponse]) (*workflow.Builder[*ChainReadResponse], error) {
	action := &capabilities.RemoteAction[*MercuryTriggerResponse, *ChainReadResponse]{
		RefName:  ref,
		TypeName: c.typeName,
	}
	return workflow.AddStep[*MercuryTriggerResponse, *ChainReadResponse](wb, action)
}

type ChainReadResponse struct {
	Read string
}
