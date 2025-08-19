package vault

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	vaultregistry "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilityregistry/v1/vault"
	vaulthandler "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/gateway/handlers/vault"
	vaultjobs "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/vault"
)

func New(chainID uint64) (*capabilities.Capability, error) {
	return capabilities.New(
		cre.VaultCapability,
		capabilities.WithJobSpecFn(vaultjobs.JobSpecFn(chainID)),
		capabilities.WithGatewayJobHandlerConfigFn(vaulthandler.HandlerConfigFn),
		capabilities.WithCapabilityRegistryV1ConfigFn(vaultregistry.CapabilityRegistryConfigFn),
		capabilities.WithValidateFn(func(c *capabilities.Capability) error {
			if chainID == 0 {
				return fmt.Errorf("chainID is required, got %d", chainID)
			}
			return nil
		}),
	)
}
