package test_workflow

import (
	"errors"

	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/poc/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/poc/workflow"
)

func CreateWorkflow() (*workflow.Spec, error) {
	root, trigger := workflow.NewWorkflowBuilder[*MercuryTriggerResponse](NewMercuryTrigger("mercury", "mercury-trigger"))

	customLogic, err := workflow.AddTransform(trigger, "custom_logic", func(_ *MercuryTriggerResponse) (*CustomType, error) {
		return &CustomType{Read: "output"}, nil
	})

	if err != nil {
		return nil, err
	}

	merge, err := workflow.Merge2("mergeReadAndLogic", trigger, customLogic, func(mr *MercuryTriggerResponse, ct *CustomType) (*MercuryTriggerResponse, error) {
		// The test ignores the second input.  May as well verify it here...
		if ct.Read != "output" {
			return nil, errors.New("unexpected value")
		}
		return mr, nil
	})
	if err != nil {
		return nil, err
	}

	consensus := capabilities.NewPureConsensus("consensus", func(inputs []*MercuryTriggerResponse) (*MercuryTriggerResponse, error) {
		return inputs[0], nil
	})

	consensusStep, err := workflow.AddConsensus(merge, consensus)
	if err != nil {
		return nil, err
	}

	err = NewChainWriter("write_chain").AddWriteTarget("write", consensusStep)
	if err != nil {
		return nil, err
	}

	return root.Build()
}

type CustomType struct {
	Read string
}
