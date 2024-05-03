package workflow

import (
	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/values"

	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/poc/capabilities"
)

var _ capability = &actionRunner[any, any]{}

type actionRunner[I, O any] struct {
	nonTriggerCapability
	capabilities.Action[I, O]
}

func (a *actionRunner[I, O]) capabilityType() commoncap.CapabilityType {
	return commoncap.CapabilityTypeAction
}

func (a *actionRunner[I, O]) Run(value values.Value) (values.Value, bool, error) {
	i, err := capabilities.UnwrapValue[I](value)
	if err != nil {
		return nil, false, err
	}
	o, cont, err := a.Invoke(i)
	if err != nil || !cont {
		return nil, false, err
	}

	wrapped, err := values.Wrap(o)
	return wrapped, true, err
}
