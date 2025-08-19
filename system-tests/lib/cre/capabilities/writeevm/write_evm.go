package writeevm

import (
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	writeevmregistry "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilityregistry/v1/writeevm"
)

func New() (*capabilities.Capability, error) {
	return capabilities.New(
		cre.WriteEVMCapability,
		capabilities.WithCapabilityRegistryV1ConfigFn(writeevmregistry.CapabilityRegistryConfigFn),
	)
}
