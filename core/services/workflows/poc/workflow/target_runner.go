package workflow

import (
	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/values"

	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/poc/capabilities"
)

type targetRunner[O any] struct {
	nonTriggerCapability
	capabilities.Target[O]
}

func (t targetRunner[O]) Run(_ string, value values.Value) (values.Value, bool, error) {
	unwrapped, err := capabilities.UnwrapValue[O](value)
	if err != nil {
		return nil, false, err
	}
	return nil, false, t.Invoke(unwrapped)
}

func (t targetRunner[O]) CapabilityType() commoncap.CapabilityType {
	return commoncap.CapabilityTypeTarget
}
func (t targetRunner[O]) private() {}

var _ capability = &targetRunner[any]{}
