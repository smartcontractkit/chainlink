package workflow

import (
	"errors"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/values"

	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/poc/capabilities"
)

type multiTriggerRunner[O any] struct {
	triggers map[string]capability
}

func (m *multiTriggerRunner[O]) Inputs() map[string]any {
	return map[string]any{}
}

func (m *multiTriggerRunner[O]) Outputs() []string {
	return []string{"$(trigger.outputs)"}
}

func (m *multiTriggerRunner[O]) Ref() string {
	return "trigger"
}

func (m *multiTriggerRunner[O]) Type() string {
	return capabilities.LocalCodeActionCapability
}

func (m *multiTriggerRunner[O]) Run(value values.Value) (values.Value, bool, error) {
	// in real life, we would probably nest the values so metadata isn't in the head, but I didn't want to impact the real tests.
	// also be safer here?
	vals := map[string]any{}
	if err := value.(*values.Map).UnwrapTo(vals); err != nil {
		return nil, false, err
	}
	ref := vals["TriggerRef"].(string)
	trigger, ok := m.triggers[ref]
	if !ok {
		return nil, false, errors.New("unknown trigger ref")
	}
	return trigger.Run(value)
}

func (m *multiTriggerRunner[O]) capabilityType() commoncap.CapabilityType {
	return commoncap.CapabilityTypeTrigger
}

func (m *multiTriggerRunner[O]) private() {}

var _ capability = &multiTriggerRunner[any]{}
