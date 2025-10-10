package readcontract

import (
	"fmt"

	"github.com/pkg/errors"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	factory "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
	chainlevel "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability/chainlevel"
)

const flag = cre.ReadContractCapability
const configTemplate = `'{"chainId":{{.ChainID}},"network":"{{.NetworkFamily}}"}'`

func New() (*capabilities.Capability, error) {
	perChainJobSpecFactory, fErr := factory.NewCapabilityJobSpecFactory(
		chainlevel.CapabilityEnabler,
		chainlevel.EnabledChainsProvider,
		chainlevel.ConfigResolver,
		chainlevel.JobNamer,
	)

	if fErr != nil {
		return nil, errors.Wrap(fErr, "failed to create capability job spec factory")
	}

	return capabilities.New(
		flag,
		capabilities.WithJobSpecFn(perChainJobSpecFactory.BuildJobSpec(
			flag,
			configTemplate,
			func(chainID uint64, _ *cre.Node) map[string]any {
				return map[string]any{
					"ChainID":       chainID,
					"NetworkFamily": "evm",
				}
			},
			factory.BinaryPathBuilder,
		)),
		capabilities.WithCapabilityRegistryV1ConfigFn(registerWithV1),
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

	for _, chainID := range nodeSetInput.ChainCapabilities[cre.ReadContractCapability].EnabledChains {
		capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
			Capability: kcr.CapabilitiesRegistryCapability{
				LabelledName:   fmt.Sprintf("read-contract-evm-%d", chainID),
				Version:        "1.0.0",
				CapabilityType: 1, // ACTION
			},
			Config: &capabilitiespb.CapabilityConfig{},
		})
	}

	return capabilities, nil
}
