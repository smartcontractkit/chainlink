package test_workflow

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/poc/workflow"
)

func CreateWorkflow() (*workflow.Spec, error) {
	root, trigger := workflow.NewWorkflowBuilder[*MercuryTriggerResponse](NewMercuryTrigger("mercury"))

	customLogic, err := workflow.AddTransform[*MercuryTriggerResponse, *CustomType](trigger, "custom_logic", func(_ *MercuryTriggerResponse) (*CustomType, error) {
		return &CustomType{Read: "output"}, nil
	})

	if err != nil {
		return nil, err
	}

	merge, err := workflow.Merge2("mergeReadAndLogic", trigger, customLogic, func(*MercuryTriggerResponse, *CustomType) (*CustomType, error) {
		return nil, nil
	})
	if err != nil {
		return nil, err
	}

	_ = merge

	return root.Build()
}

type CustomType struct {
	Read string
}
