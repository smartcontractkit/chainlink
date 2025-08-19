package readcontract

import (
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	readcontractregistry "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilityregistry/v1/readcontract"
	factory "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
	chainlevel "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability/chainlevel"
)

const readContractConfigTemplate = `'{"chainId":{{.ChainID}},"network":"{{.NetworkFamily}}"}'`

func New() (*capabilities.Capability, error) {
	perChainJobSpecFactory := factory.NewCapabilityJobSpecFactory(
		chainlevel.IsEnabled,
		chainlevel.EnabledChains,
		chainlevel.ConfigResolver,
		chainlevel.JobName,
	)

	return capabilities.New(
		cre.ReadContractCapability,
		capabilities.WithJobSpecFn(perChainJobSpecFactory.BuildJobSpecFn(
			cre.ReadContractCapability,
			readContractConfigTemplate,
			func(chainID uint64, _ *cre.NodeMetadata) map[string]any {
				return map[string]any{
					"ChainID":       chainID,
					"NetworkFamily": "evm",
				}
			},
			factory.BinaryPathBuilder,
		)),
		capabilities.WithCapabilityRegistryV1ConfigFn(readcontractregistry.CapabilityRegistryConfigFn),
	)
}
