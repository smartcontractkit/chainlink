package workflow

import (
	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/values"

	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/poc/capabilities"
)

type triggerRunner[O any] struct {
	capabilities.Trigger[O]
}

func (t triggerRunner[O]) Output() string {
	return "$(trigger.outputs)"
}

func (t triggerRunner[O]) Inputs() map[string]any {
	return map[string]any{}
}

func (t triggerRunner[O]) Run(value values.Value) (values.Value, bool, error) {
	vmap := value.(*values.Map).Underlying["action"]
	output, err := t.Transform(vmap)
	if err != nil {
		return nil, false, err
	}
	wrapped, err := values.Wrap(output)
	return wrapped, err == nil, err
}

func (t triggerRunner[O]) capabilityType() commoncap.CapabilityType {
	return commoncap.CapabilityTypeTrigger
}

func (t triggerRunner[O]) private() {}

var _ capability = &triggerRunner[any]{}
