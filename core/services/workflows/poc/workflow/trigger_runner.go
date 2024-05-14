package workflow

import (
	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/values"

	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/poc/capabilities"
)

type triggerRunner[O any] struct {
	capabilities.Trigger[O]
}

func (t triggerRunner[O]) StepDependencies() []string {
	return []string{}
}

func (t triggerRunner[O]) Output() string {
	return "trigger"
}

func (t triggerRunner[O]) Inputs() map[string]string {
	return map[string]string{}
}

func (t triggerRunner[O]) Run(_ string, value values.Value) (values.Value, bool, error) {
	vmap := value.(*values.Map).Underlying["action"]
	output, err := t.Transform(vmap)
	if err != nil {
		return nil, false, err
	}
	wrapped, err := values.Wrap(output)
	return wrapped, err == nil, err
}

func (t triggerRunner[O]) CapabilityType() commoncap.CapabilityType {
	return commoncap.CapabilityTypeTrigger
}

func (t triggerRunner[O]) private() {}

var _ capability = &triggerRunner[any]{}
