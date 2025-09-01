package environment

import (
	"fmt"
	"maps"
	"os"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	ctfconfig "github.com/smartcontractkit/chainlink-testing-framework/lib/config"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecapabilities "github.com/smartcontractkit/chainlink/system-tests/lib/cre/capabilities"
	libdon "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	creconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/config"
	cresecrets "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/secrets"
	creflags "github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

func BuildTopology(
	registryChainSelector uint64,
	nodeSets []*cre.CapabilitiesAwareNodeSet,
	infraInput infra.Input,
	evmChainIDs []int,
	solChainIDs []string,
	blockchainOutput map[uint64]*cre.WrappedBlockchainOutput,
	addressBook deployment.AddressBook,
	datastore datastore.DataStore,
	capabilities []cre.InstallableCapability,
	capabilityConfigs cre.CapabilityConfigs,
	copyCapabilityBinaries bool,
) (*cre.Topology, []*cre.CapabilitiesAwareNodeSet, error) {
	topologyErr := libdon.ValidateTopology(nodeSets, infraInput)
	if topologyErr != nil {
		return nil, nil, errors.Wrap(topologyErr, "failed to validate topology")
	}

	topology, err := libdon.BuildTopology(nodeSets, infraInput, registryChainSelector)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to build topology")
	}

	localNodeSets := copyCapabilityAwareNodeSets(nodeSets)

	// Generate EVM and P2P keys or read them from the config
	// That way we can pass them final configs and do away with restarting the nodes
	var keys *cre.GenerateKeysOutput

	keysOutput, keysOutputErr := cresecrets.KeysOutputFromConfig(localNodeSets)
	if keysOutputErr != nil {
		return nil, nil, errors.Wrap(keysOutputErr, "failed to generate keys output")
	}

	generateKeysInput := &cre.GenerateKeysInput{
		GenerateEVMKeysForChainIDs: evmChainIDs,
		GenerateSolKeysForChainIDs: solChainIDs,
		GenerateP2PKeys:            true,
		Topology:                   topology,
		Password:                   "", // since the test runs on private ephemeral blockchain we don't use real keys and do not care a lot about the password
		Out:                        keysOutput,
	}
	keys, keysErr := cresecrets.GenerateKeys(generateKeysInput)
	if keysErr != nil {
		return nil, nil, errors.Wrap(keysErr, "failed to generate keys")
	}

	topology, addKeysErr := cresecrets.AddKeysToTopology(topology, keys)
	if addKeysErr != nil {
		return nil, nil, errors.Wrap(addKeysErr, "failed to add keys to topology")
	}

	capabilitiesPeeringData, ocrPeeringData, peeringErr := libdon.FindPeeringData(topology)
	if peeringErr != nil {
		return nil, nil, errors.Wrap(peeringErr, "failed to find peering data")
	}

	topology.CapabilitiesPeeringData = capabilitiesPeeringData
	topology.OCRPeeringData = ocrPeeringData

	for i, donMetadata := range topology.DonsMetadata {
		configsFound := 0
		secretsFound := 0
		for _, nodeSpec := range localNodeSets[i].NodeSpecs {
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
		// We can't allow providing only configs, because we can't replace secret-related values in the configs
		// If both are provided, we assume that the user knows what they are doing and we don't need to validate anything
		// And that configs match the secrets
		if configsFound > 0 && secretsFound == 0 {
			return nil, nil, fmt.Errorf("nodese config overrides are provided for DON %d, but not secrets. You need to either provide both, only secrets or nothing at all", donMetadata.ID)
		}

		configFactoryFunctions := make([]cre.NodeConfigFn, 0)
		for _, capability := range capabilities {
			configFactoryFunctions = append(configFactoryFunctions, capability.NodeConfigFn())
		}

		// generate configs only if they are not provided
		if configsFound == 0 {
			config, configErr := creconfig.Generate(
				cre.GenerateConfigsInput{
					AddressBook:             addressBook,
					Datastore:               datastore,
					DonMetadata:             donMetadata,
					BlockchainOutput:        blockchainOutput,
					Flags:                   donMetadata.Flags,
					CapabilitiesPeeringData: capabilitiesPeeringData,
					OCRPeeringData:          ocrPeeringData,
					HomeChainSelector:       topology.HomeChainSelector,
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

		// generate secrets only if they are not provided
		if secretsFound == 0 {
			secretsInput := &cre.GenerateSecretsInput{
				DonMetadata: donMetadata,
			}

			if evmKeys, ok := keys.EVMKeys[donMetadata.ID]; ok {
				secretsInput.EVMKeys = evmKeys
			}

			if solKeys, ok := keys.SolKeys[donMetadata.ID]; ok {
				secretsInput.SolKeys = solKeys
			}

			if p2pKeys, ok := keys.P2PKeys[donMetadata.ID]; ok {
				secretsInput.P2PKeys = p2pKeys
			}

			// EVM, Solana and P2P keys will be provided to nodes as secrets
			secrets, secretsErr := cresecrets.GenerateSecrets(
				secretsInput,
			)
			if secretsErr != nil {
				return nil, nil, errors.Wrap(secretsErr, "failed to generate secrets")
			}

			for j := range donMetadata.NodesMetadata {
				localNodeSets[i].NodeSpecs[j].Node.TestSecretsOverrides = secrets[j]
			}
		}

		if !copyCapabilityBinaries {
			continue
		}

		customBinariesPaths := make(map[cre.CapabilityFlag]string)
		for flag, config := range capabilityConfigs {
			if creflags.HasFlagForAnyChain(donMetadata.Flags, flag) && config.BinaryPath != "" {
				customBinariesPaths[flag] = config.BinaryPath
			}
		}

		executableErr := crecapabilities.MakeBinariesExecutable(customBinariesPaths)
		if executableErr != nil {
			return nil, nil, errors.Wrap(executableErr, "failed to make binaries executable")
		}

		var appendErr error
		localNodeSets[i], appendErr = crecapabilities.AppendBinariesPathsNodeSpec(localNodeSets[i], donMetadata, customBinariesPaths)
		if appendErr != nil {
			return nil, nil, errors.Wrapf(appendErr, "failed to append binaries paths to node spec for DON %d", donMetadata.ID)
		}
	}

	// Add env vars, which were provided programmatically, to the node specs
	// or fail, if node specs already had some env vars set in the TOML config
	for donIdx, donMetadata := range topology.DonsMetadata {
		hasEnvVarsInTomlConfig := false
		for nodeIdx, nodeSpec := range localNodeSets[donIdx].NodeSpecs {
			if len(nodeSpec.Node.EnvVars) > 0 {
				hasEnvVarsInTomlConfig = true
				break
			}

			localNodeSets[donIdx].NodeSpecs[nodeIdx].Node.EnvVars = localNodeSets[donIdx].EnvVars
		}

		if hasEnvVarsInTomlConfig && len(localNodeSets[donIdx].EnvVars) > 0 {
			return nil, nil, fmt.Errorf("extra env vars for Chainlink Nodes are provided in the TOML config for the %s DON, but you tried to provide them programatically. Please set them only in one place", donMetadata.Name)
		}
	}

	// Deploy the DONs
	// Hack for CI that allows us to dynamically set the chainlink image and version
	// CTFv2 currently doesn't support dynamic image and version setting
	if os.Getenv("CI") == "true" {
		// Due to how we pass custom env vars to reusable workflow we need to use placeholders, so first we need to resolve what's the name of the target environment variable
		// that stores chainlink version and then we can use it to resolve the image name
		for i := range localNodeSets {
			image := fmt.Sprintf("%s:%s", os.Getenv(ctfconfig.E2E_TEST_CHAINLINK_IMAGE_ENV), ctfconfig.MustReadEnvVar_String(ctfconfig.E2E_TEST_CHAINLINK_VERSION_ENV))
			for j := range localNodeSets[i].NodeSpecs {
				localNodeSets[i].NodeSpecs[j].Node.Image = image
				// unset docker context and file path, so that we can use the image from the registry
				localNodeSets[i].NodeSpecs[j].Node.DockerContext = ""
				localNodeSets[i].NodeSpecs[j].Node.DockerFilePath = ""
			}
		}
	}

	return topology, localNodeSets, nil
}

func copyCapabilityAwareNodeSets(
	nodeSets []*cre.CapabilitiesAwareNodeSet,
) []*cre.CapabilitiesAwareNodeSet {
	copiedNodeSets := make([]*cre.CapabilitiesAwareNodeSet, len(nodeSets))
	for i, originalNs := range nodeSets {
		if originalNs == nil {
			copiedNodeSets[i] = nil
			continue
		}

		newNs := &cre.CapabilitiesAwareNodeSet{
			BootstrapNodeIndex: originalNs.BootstrapNodeIndex,
			GatewayNodeIndex:   originalNs.GatewayNodeIndex,
		}

		if originalNs.Input != nil {
			inputCopy := *originalNs.Input
			newNs.Input = &inputCopy
		}

		if originalNs.Capabilities != nil {
			newNs.Capabilities = make([]string, len(originalNs.Capabilities))
			copy(newNs.Capabilities, originalNs.Capabilities)
		}

		if originalNs.ComputedCapabilities != nil {
			newNs.ComputedCapabilities = make([]string, len(originalNs.ComputedCapabilities))
			copy(newNs.ComputedCapabilities, originalNs.ComputedCapabilities)
		}

		if originalNs.DONTypes != nil {
			newNs.DONTypes = make([]string, len(originalNs.DONTypes))
			copy(newNs.DONTypes, originalNs.DONTypes)
		}

		if originalNs.SupportedChains != nil {
			newNs.SupportedChains = make([]uint64, len(originalNs.SupportedChains))
			copy(newNs.SupportedChains, originalNs.SupportedChains)
		}

		if originalNs.EnvVars != nil {
			newNs.EnvVars = make(map[string]string, len(originalNs.EnvVars))
			maps.Copy(newNs.EnvVars, originalNs.EnvVars)
		}

		if originalNs.ChainCapabilities != nil {
			newNs.ChainCapabilities = make(map[string]*cre.ChainCapabilityConfig, len(originalNs.ChainCapabilities))
			maps.Copy(newNs.ChainCapabilities, originalNs.ChainCapabilities)
		}

		if originalNs.SupportedSolChains != nil {
			newNs.SupportedSolChains = make([]string, len(originalNs.SupportedSolChains))
			copy(newNs.SupportedSolChains, originalNs.SupportedSolChains)
		}

		copiedNodeSets[i] = newNs
	}

	return copiedNodeSets
}
