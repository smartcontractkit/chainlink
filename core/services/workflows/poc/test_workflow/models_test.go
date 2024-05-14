package test_workflow_test

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/poc/test_workflow"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/poc/wasm/wrokflowtesting"
)

// These would be generated from protos provided by the capability author
// For now, they mimic what's in the test for the multi input workflow

func AddMercuryTriggerToRegistry(ref string, registry *wrokflowtesting.Registry, value *test_workflow.MercuryTriggerResponse) error {
	return registry.RegisterTrigger(ref, value)
}

type ChainWriterMock struct {
	wrokflowtesting.TargetMock[*test_workflow.MercuryTriggerResponse]
}
