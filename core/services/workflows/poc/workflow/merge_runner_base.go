package workflow

import (
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

func (m mergeRunnerBase) capabilityType() commoncap.CapabilityType {
	return commoncap.CapabilityTypeAction
}

func mergeOutputs(cs ...capability) map[string]any {
	outputs := make(map[string]any)
	for _, c := range cs {
		outputs[c.Ref()] = c.Outputs()
	}
	return outputs
}
