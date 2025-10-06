package environment

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	creconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/config"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

func PrepareNodeTOMLConfigurations(
	registryChainSelector uint64,
	nodeSets []*cre.CapabilitiesAwareNodeSet,
	provider infra.Provider,
	blockchainOutputs []*cre.WrappedBlockchainOutput,
	addressBook deployment.AddressBook,
	datastore datastore.DataStore,
	capabilities []cre.InstallableCapability,
	capabilityConfigs cre.CapabilityConfigs,
) (*cre.Topology, []*cre.CapabilitiesAwareNodeSet, error) {
	topology, tErr := cre.NewTopology(nodeSets, provider)
	if tErr != nil {
		return nil, nil, errors.Wrap(tErr, "failed to create topology")
	}

	bt, err := topology.BootstrapNode()
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to find bootstrap node")
	}

	capabilitiesPeeringData, ocrPeeringData, peeringErr := cre.PeeringCfgs(bt)
	if peeringErr != nil {
		return nil, nil, errors.Wrap(peeringErr, "failed to find peering data")
	}

	localNodeSets := topology.CapabilitiesAwareNodeSets()
	chainPerSelector := make(map[uint64]*cre.WrappedBlockchainOutput)
	for _, bcOut := range blockchainOutputs {
		if bcOut.SolChain != nil {
			sel := bcOut.SolChain.ChainSelector
			chainPerSelector[sel] = bcOut
			chainPerSelector[sel].ChainSelector = sel
			chainPerSelector[sel].SolChain = bcOut.SolChain
			chainPerSelector[sel].SolChain.ArtifactsDir = bcOut.SolChain.ArtifactsDir
			continue
		}
		chainPerSelector[bcOut.ChainSelector] = bcOut
	}

	for i, donMetadata := range topology.DonsMetadata.List() {
		// make sure that either all or none of the node specs have config or secrets provided in the TOML config
		configsFound := 0
		secretsFound := 0
		nodeSet := localNodeSets[i]

		for _, nodeSpec := range nodeSet.NodeSpecs {
			if nodeSpec.Node.TestConfigOverrides != "" {
				configsFound++
			}
			if nodeSpec.Node.TestSecretsOverrides != "" {
				secretsFound++
			}
		}

		if configsFound != 0 && configsFound != len(localNodeSets[i].NodeSpecs) {
			return nil, nil, fmt.Errorf("%d out of %d node specs have config overrides. Either provide overrides for all nodes or none at all", configsFound, len(localNodeSets[i].NodeSpecs))
		}

		if secretsFound != 0 && secretsFound != len(localNodeSets[i].NodeSpecs) {
			return nil, nil, fmt.Errorf("%d out of %d node specs have secrets overrides. Either provide overrides for all nodes or none at all", secretsFound, len(localNodeSets[i].NodeSpecs))
		}

		// Allow providing only secrets, because we can decode them and use them to generate configs
		// We can't allow providing only configs, because we don't want to deal with parsing configs to set new secrets there.
		// If both are provided, we assume that the user knows what they are doing and we don't need to validate anything
		if configsFound > 0 && secretsFound == 0 {
			return nil, nil, fmt.Errorf("nodespec config overrides are provided for DON %s, but not secrets. You need to either provide both, only secrets or nothing at all", donMetadata.Name)
		}

		configFactoryFunctions := make([]cre.NodeConfigTransformerFn, 0)
		for _, capability := range capabilities {
			configFactoryFunctions = append(configFactoryFunctions, capability.NodeConfigTransformerFn())
		}

		// generate node TOML configs only if they are not provided in the environment TOML config
		if configsFound == 0 {
			config, configErr := creconfig.Generate(
				cre.GenerateConfigsInput{
					AddressBook:             addressBook,
					Datastore:               datastore,
					DonMetadata:             donMetadata,
					BlockchainOutput:        chainPerSelector,
					Flags:                   donMetadata.Flags,
					CapabilitiesPeeringData: capabilitiesPeeringData,
					OCRPeeringData:          ocrPeeringData,
					HomeChainSelector:       registryChainSelector,
					GatewayConnectorOutput:  topology.GatewayConnectorOutput,
					NodeSet:                 localNodeSets[i],
					CapabilityConfigs:       capabilityConfigs,
				},
				configFactoryFunctions,
			)
			if configErr != nil {
				return nil, nil, errors.Wrap(configErr, "failed to generate config")
			}

			for j := range donMetadata.NodesMetadata {
				localNodeSets[i].NodeSpecs[j].Node.TestConfigOverrides = config[j]
			}
		}

		// generate node TOML secrets only if they are not provided in the environment TOML config
		if secretsFound == 0 {
			for nodeIndex := range donMetadata.NodesMetadata {
				wnode := donMetadata.NodesMetadata[nodeIndex]
				nodeSecretsTOML, err := wnode.Keys.ToNodeSecretsTOML()
				if err != nil {
					return nil, nil, errors.Wrapf(err, "failed to marshal node secrets (DON: %s, Node index: %d)", donMetadata.Name, nodeIndex)
				}
				localNodeSets[i].NodeSpecs[nodeIndex].Node.TestSecretsOverrides = nodeSecretsTOML
			}
		}
	}

	return topology, localNodeSets, nil
}
