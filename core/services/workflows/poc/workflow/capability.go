package workflow

import (
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

type capability interface {
	Type() string
	Ref() string
	Inputs() map[string]any
	Outputs() []string
	capabilityType() capabilities.CapabilityType
	LocalCapability
	private()
}

type LocalCapability interface {
	Run(value values.Value) (values.Value, bool, error)
}
