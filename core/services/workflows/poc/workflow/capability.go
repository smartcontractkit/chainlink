package workflow

import (
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

type capability interface {
	Type() string
	Ref() string
	Inputs() map[string]any
	Output() string
	LocalCapability
	private()
}

type LocalCapability interface {
	Run(ref string, value values.Value) (values.Value, bool, error)
	CapabilityType() capabilities.CapabilityType
}
