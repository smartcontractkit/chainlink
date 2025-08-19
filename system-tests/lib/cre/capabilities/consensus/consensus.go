package consensus

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	consensusregistry "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilityregistry/v1/consensus"
	consensusjobs "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/consensus"
)

func NewV1(chainID uint64) (*capabilities.Capability, error) {
	return capabilities.New(
		cre.ConsensusCapability,
		capabilities.WithJobSpecFn(consensusjobs.V1JobSpecFn(chainID)),
		capabilities.WithCapabilityRegistryV1ConfigFn(consensusregistry.CapabilityV1RegistryConfigFn),
		capabilities.WithValidateFn(func(c *capabilities.Capability) error {
			if chainID == 0 {
				return fmt.Errorf("chainID is required, got %d", chainID)
			}
			return nil
		}),
	)
}

func NewV2() (*capabilities.Capability, error) {
	return capabilities.New(
		cre.ConsensusCapabilityV2,
		capabilities.WithJobSpecFn(consensusjobs.V2JobSpecFn),
		capabilities.WithCapabilityRegistryV1ConfigFn(consensusregistry.ConsensusV2CapabilityFn),
	)
}
