package writeevm

import (
	"errors"
	"strings"

	corevm "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
)

const flag = cre.WriteEVMCapability

func New() (*capabilities.Capability, error) {
	return capabilities.New(
		flag,
		capabilities.WithCapabilityRegistryV1ConfigFn(registerWithV1),
		// Unfortunately, node config generation function is currently defined here system-tests/lib/cre/don/config/definitions.go
		// Instead of being passed as functional option, because we need to add the [EVM.Workflow] section in a specific place
		// of the config instead of concatenating it in any random place
	)
}

func registerWithV1(_ []string, nodeSetInput *cre.CapabilitiesAwareNodeSet) ([]keystone_changeset.DONCapabilityWithConfig, error) {
	capabilities := make([]keystone_changeset.DONCapabilityWithConfig, 0)

	if nodeSetInput == nil {
		return nil, errors.New("node set input is nil")
	}

	// it's fine if there are no chain capabilities
	if nodeSetInput.ChainCapabilities == nil {
		return nil, nil
	}

	if _, ok := nodeSetInput.ChainCapabilities[flag]; !ok {
		return nil, nil
	}

	for _, chainID := range nodeSetInput.ChainCapabilities[flag].EnabledChains {
		fullName := corevm.GenerateWriteTargetName(chainID)
		splitName := strings.Split(fullName, "@")

		capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
			Capability: kcr.CapabilitiesRegistryCapability{
				LabelledName:   splitName[0],
				Version:        splitName[1],
				CapabilityType: 3, // TARGET
				ResponseType:   1, // OBSERVATION_IDENTICAL
			},
			Config: &capabilitiespb.CapabilityConfig{},
		})
	}

	return capabilities, nil
}
