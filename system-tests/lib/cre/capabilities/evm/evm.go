package evm

import (
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	evmregistry "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilityregistry/v1/evm"
	evmjobs "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/evm"
)

func New() (*capabilities.Capability, error) {
	return capabilities.New(
		cre.EVMCapability,
		capabilities.WithJobSpecFn(evmjobs.JobSpecFn),
		capabilities.WithCapabilityRegistryV1ConfigFn(evmregistry.CapabilityRegistryConfigFn),
	)
}
