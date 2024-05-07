package workflows

import (
	pocWorkflow "github.com/smartcontractkit/chainlink/v2/core/services/workflows/poc/workflow"
)

type codeSpecBuilder struct {
	CodeConfig *codeConfig
	Workflow   *pocWorkflow.Spec
}

func (c *codeSpecBuilder) Build() (workflowSpec, error) {
	return workflowSpec{
		Triggers:  c.convertDefs(c.Workflow.Triggers),
		Actions:   c.convertDefs(c.Workflow.Actions),
		Consensus: c.convertDefs(c.Workflow.Consensus),
		Targets:   c.convertDefs(c.Workflow.Targets),
	}, nil
}

func (c *codeSpecBuilder) convertDefs(pocDefs []pocWorkflow.StepDefinition) []stepDefinition {
	defs := make([]stepDefinition, len(pocDefs))
	for i, pocDef := range pocDefs {
		stepType := pocDef.TypeRef
		rawStepType, ok := c.CodeConfig.TypeMap[pocDef.TypeRef]
		if ok {
			stepType = rawStepType
		}

		defs[i] = stepDefinition{
			Type:   stepType,
			Ref:    pocDef.Ref,
			Inputs: pocDef.Inputs,
			Config: c.CodeConfig.Config[pocDef.Ref],
		}
	}
	return defs
}

var _ specBuilder = &codeSpecBuilder{}
