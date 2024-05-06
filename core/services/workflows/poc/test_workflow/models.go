package test_workflow

import (
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/poc/capabilities"
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

type ChainWrite struct {
}
