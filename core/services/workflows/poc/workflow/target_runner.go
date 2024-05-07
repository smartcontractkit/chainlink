package workflow

import (
	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/values"

	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/poc/capabilities"
)

type targetRunner[O any] struct {
	inputs map[string]any
	capabilities.Target[O]
}

func (t targetRunner[O]) Inputs() map[string]any {
	return t.inputs
}

func (t targetRunner[O]) Output() string {
	return "$(target.outputs)"
}

func (t targetRunner[O]) Run(value values.Value) (values.Value, bool, error) {
	unwrapped, err := capabilities.UnwrapValue[O](value)
	if err != nil {
		return nil, false, err
	}
	return nil, false, t.Invoke(unwrapped)
}

func (t targetRunner[O]) capabilityType() commoncap.CapabilityType {
	return commoncap.CapabilityTypeTarget
}
func (t targetRunner[O]) private() {}

var _ capability = &targetRunner[any]{}
