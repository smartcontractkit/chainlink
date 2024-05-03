package test_workflow

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/poc/workflow"
)

func CreateWorkflow() (*workflow.Spec, error) {
	root, trigger := workflow.NewWorkflowBuilder[*MercuryTriggerResponse](NewMercuryTrigger("mercury"))
	read, err := NewChainReader("read_chain_action").AddReadAction("read_chain_action", trigger)
	if err != nil {
		return nil, err
	}


	compute := workflow.AddTransform(read, "compute", func(_ *MercuryTriggerResponse, _ *ChainReadResponse) (, error) {
		return , nil
	}

	return root.Build()
}
