package workflow

import (
	"fmt"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"

	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/poc/capabilities"
)

type mergeRunnerBase struct {
	nonTriggerCapability
}

func (m mergeRunnerBase) Type() string {
	return capabilities.LocalCodeActionCapability
}

func (m mergeRunnerBase) Ref() string {
	return m.ref
}

func (m mergeRunnerBase) CapabilityType() commoncap.CapabilityType {
	return commoncap.CapabilityTypeAction
}

func mergeOutputs(cs ...capability) map[string]any {
	outputs := map[string]any{}
	for i, c := range cs {
		outputs[fmt.Sprintf("action%d", i+1)] = c.Output()
	}
	return outputs
}
