package config

import (
	"context"
	"fmt"
	"maps"
	"math/big"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/libocr/commontypes"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-evm/pkg/config/chaintype"
	evmconfigtoml "github.com/smartcontractkit/chainlink-evm/pkg/config/toml"
	solcfg "github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	chipingressset "github.com/smartcontractkit/chainlink-testing-framework/framework/components/dockercompose/chip_ingress_set"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"

	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	coretoml "github.com/smartcontractkit/chainlink/v2/core/config/toml"
	corechainlink "github.com/smartcontractkit/chainlink/v2/core/services/chainlink"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	creblockchains "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	aptoschain "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/aptos"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/solana"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

const TronEVMChainID = 3360022319

func PrepareNodeTOMLs(
	ctx context.Context,
	topology *cre.Topology,
	creEnv *cre.Environment,
	nodeSets []*cre.NodeSet,
	capabilities []cre.InstallableCapability, // Deprecated, use Features instead and modify node configs inside a Feature
	nodeConfigTransformerFns []cre.NodeConfigTransformerFn,
) ([]*cre.NodeSet, error) {
	bt, hasBootstrap := topology.Bootstrap()
	if !hasBootstrap {
		return nil, errors.New("no DON contains a bootstrap node, but exactly one is required")
	}

	capabilitiesPeeringData, ocrPeeringData, peeringErr := cre.PeeringCfgs(bt)
	if peeringErr != nil {
		return nil, errors.Wrap(peeringErr, "failed to find peering data")
	}

	localNodeSets := topology.NodeSets()
	chainPerSelector := make(map[uint64]creblockchains.Blockchain)
	for _, bc := range creEnv.Blockchains {
		chainPerSelector[bc.ChainSelector()] = bc
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
			return nil, fmt.Errorf("%d out of %d node specs have config overrides. Either provide overrides for all nodes or none at all", configsFound, len(localNodeSets[i].NodeSpecs))
		}

		if secretsFound != 0 && secretsFound != len(localNodeSets[i].NodeSpecs) {
			return nil, fmt.Errorf("%d out of %d node specs have secrets overrides. Either provide overrides for all nodes or none at all", secretsFound, len(localNodeSets[i].NodeSpecs))
		}

		// Allow providing only secrets, because we can decode them and use them to generate configs
		// We can't allow providing only configs, because we don't want to deal with parsing configs to set new secrets there.
		// If both are provided, we assume that the user knows what they are doing and we don't need to validate anything
		if configsFound > 0 && secretsFound == 0 {
			return nil, fmt.Errorf("nodespec config overrides are provided for DON %s, but not secrets. You need to either provide both, only secrets or nothing at all", donMetadata.Name)
		}

		configFactoryFunctions := make([]cre.NodeConfigTransformerFn, 0)
		for _, capability := range capabilities {
			configFactoryFunctions = append(configFactoryFunctions, capability.NodeConfigTransformerFn())
		}
		configFactoryFunctions = append(configFactoryFunctions, nodeConfigTransformerFns...) // allow passing custom transformers

		// generate node TOML configs only if they are not provided in the environment TOML config
		if configsFound == 0 {
			config, configErr := generateNodeTomlConfig(
				cre.GenerateConfigsInput{
					Datastore:               creEnv.CldfEnvironment.DataStore,
					ContractVersions:        creEnv.ContractVersions,
					DonMetadata:             donMetadata,
					Blockchains:             chainPerSelector,
					Flags:                   donMetadata.Flags,
					CapabilitiesPeeringData: capabilitiesPeeringData,
					OCRPeeringData:          ocrPeeringData,
					RegistryChainSelector:   creEnv.RegistryChainSelector,
					Topology:                topology,
					Provider:                creEnv.Provider,
				},
				configFactoryFunctions,
			)
			if configErr != nil {
				return nil, errors.Wrap(configErr, "failed to generate config")
			}

			// For Kubernetes, save configs to ConfigMaps
			if creEnv.Provider.IsKubernetes() {
				instanceNames := generateInstanceNames(localNodeSets[i].Name, donMetadata.NodesMetadata)
				for j, instanceName := range instanceNames {
					err := CreateConfigOverride(ctx, framework.L, CreateConfigOverrideInput{
						Namespace:    creEnv.Provider.Kubernetes.Namespace,
						InstanceName: instanceName,
						ConfigToml:   config[j],
					})
					if err != nil {
						return nil, fmt.Errorf("failed to create config override for node %s: %w", instanceName, err)
					}
				}
			}

			// Set TestConfigOverrides for all providers (needed for features to work)
			for j := range donMetadata.NodesMetadata {
				localNodeSets[i].NodeSpecs[j].Node.TestConfigOverrides = config[j]
			}
		}

		// generate node TOML secrets only if they are not provided in the environment TOML config
		if secretsFound == 0 {
			instanceNames := generateInstanceNames(localNodeSets[i].Name, donMetadata.NodesMetadata)
			for nodeIndex := range donMetadata.NodesMetadata {
				wnode := donMetadata.NodesMetadata[nodeIndex]
				nodeSecretsTOML, err := wnode.Keys.ToNodeSecretsTOML()
				if err != nil {
					return nil, errors.Wrapf(err, "failed to marshal node secrets (DON: %s, Node index: %d)", donMetadata.Name, nodeIndex)
				}

				// For Kubernetes, save secrets to Kubernetes Secrets
				if creEnv.Provider.IsKubernetes() {
					labels := map[string]string{
						"app":  "chainlink",
						"don":  donMetadata.Name,
						"node": instanceNames[nodeIndex],
					}

					err := CreateSecretsOverride(ctx, framework.L, CreateSecretsOverrideInput{
						Namespace:    creEnv.Provider.Kubernetes.Namespace,
						InstanceName: instanceNames[nodeIndex],
						SecretsToml:  nodeSecretsTOML,
						Labels:       labels,
					})
					if err != nil {
						return nil, fmt.Errorf("failed to create secrets override for node %s: %w", instanceNames[nodeIndex], err)
					}
					framework.L.Debug().Msgf("   Secret keys: P2PKey, DKGKey, EVM keys (%d chains), Solana keys (%d chains)",
						len(wnode.Keys.EVM), len(wnode.Keys.Solana))
				}

				// Set TestSecretsOverrides for all providers (needed for features and key access)
				localNodeSets[i].NodeSpecs[nodeIndex].Node.TestSecretsOverrides = nodeSecretsTOML
			}
		}
	}

	// Transform UserConfigOverrides to use platform-specific Docker host addresses.
	// This handles differences between macOS (host.docker.internal) and Linux (172.17.0.1)
	// for URLs in user-provided config overrides (e.g., AdditionalSources).
	for i := range localNodeSets {
		for j := range localNodeSets[i].NodeSpecs {
			if localNodeSets[i].NodeSpecs[j].Node.UserConfigOverrides != "" {
				localNodeSets[i].NodeSpecs[j].Node.UserConfigOverrides = transformUserConfigOverrides(
					localNodeSets[i].NodeSpecs[j].Node.UserConfigOverrides,
				)
			}
		}
	}

	return localNodeSets, nil
}

func generateNodeTomlConfig(input cre.GenerateConfigsInput, nodeConfigTransformers []cre.NodeConfigTransformerFn) (cre.NodeIndexToConfigOverride, error) {
	configOverrides := make(cre.NodeIndexToConfigOverride)

	if err := input.Validate(); err != nil {
		return nil, errors.Wrap(err, "input validation failed")
	}

	commonInputs, inputsErr := gatherCommonInputs(input)
	if inputsErr != nil {
		return nil, errors.Wrap(inputsErr, "failed to gather common inputs")
	}

	for nodeIdx, nodeMetadata := range input.DonMetadata.NodesMetadata {
		nodeConfig := baseNodeConfig(commonInputs, input.DonMetadata, nodeMetadata)
		for _, role := range nodeMetadata.Roles {
			switch role {
			case cre.BootstrapNode:
				var cErr error
				nodeConfig, cErr = addBootstrapNodeConfig(nodeConfig, input.OCRPeeringData, commonInputs)
				if cErr != nil {
					return nil, errors.Wrapf(cErr, "failed to add bootstrap node config for node at index %d in DON %s", nodeIdx, input.DonMetadata.Name)
				}
			case cre.WorkerNode:
				var cErr error
				nodeConfig, cErr = addWorkerNodeConfig(nodeConfig, input.Topology, input.OCRPeeringData, commonInputs, input.DonMetadata, nodeMetadata)
				if cErr != nil {
					return nil, errors.Wrapf(cErr, "failed to add worker node config for node at index %d in DON %s", nodeIdx, input.DonMetadata.Name)
				}
			case cre.GatewayNode:
				var cErr error
				nodeConfig, cErr = addGatewayNodeConfig(nodeConfig, input.OCRPeeringData, commonInputs, nodeMetadata)
				if cErr != nil {
					return nil, errors.Wrapf(cErr, "failed to add gateway node config for node at index %d in DON %s", nodeIdx, input.DonMetadata.Name)
				}
			default:
				supportedRoles := []string{cre.BootstrapNode, cre.WorkerNode, cre.GatewayNode}
				return nil, fmt.Errorf("unsupported node type %s found for node at index %d in DON %s. Supported roles: %s", role, nodeIdx, input.DonMetadata.Name, strings.Join(supportedRoles, ", "))
			}
		}

		marshalledConfig, mErr := toml.Marshal(nodeConfig)
		if mErr != nil {
			return nil, errors.Wrapf(mErr, "failed to marshal node config for node at index %d in DON %s", nodeIdx, input.DonMetadata.Name)
		}

		configOverrides[nodeIdx] = string(marshalledConfig)
	}

	// execute capability-provided functions that transform the node config (currently: write-evm, write-solana)
	// these functions must return whole node configs after transforming them, instead of just returning configuration parts
	// that need to be merged into the existing config
	for _, transformer := range nodeConfigTransformers {
		if transformer == nil {
			continue
		}

		modifiedConfigs, err := transformer(input, configOverrides)
		if err != nil {
			return nil, errors.Wrap(err, "failed to generate nodeset configs")
		}

		maps.Copy(configOverrides, modifiedConfigs)
	}

	return configOverrides, nil
}

func baseNodeConfig(commonInputs *commonInputs, donMetadata *cre.DonMetadata, nodeMetadata *cre.NodeMetadata) corechainlink.Config {
	c := corechainlink.Config{
		Core: coretoml.Core{
			InsecurePPROFHeap: ptr.Ptr(false), // Set to true to enable v2/debug/pprof/heap endpoint
			Feature: coretoml.Feature{
				LogPoller: ptr.Ptr(true),
			},
			Log: coretoml.Log{
				JSONConsole: ptr.Ptr(true),
				Level:       ptr.Ptr(coretoml.LogLevel(zapcore.DebugLevel)),
			},
			OCR2: coretoml.OCR2{
				Enabled:              ptr.Ptr(true),
				DatabaseTimeout:      commonconfig.MustNewDuration(1 * time.Second),
				ContractPollInterval: commonconfig.MustNewDuration(1 * time.Second),
			},
			CRE: coretoml.CreConfig{
				EnableDKGRecipient:   ptr.Ptr(true),
				UseLocalTimeProvider: ptr.Ptr(false),
				DebugMode:            ptr.Ptr(true),
			},
		},
	}

	if commonInputs.provider.IsDocker() {
		nodeIdentifier := donMetadata.Name + "-node-" + strconv.Itoa(nodeMetadata.Index)
		c.Telemetry = coretoml.Telemetry{
			Enabled:             ptr.Ptr(true),
			Endpoint:            ptr.Ptr(strings.TrimPrefix(framework.HostDockerInternal(), "http://") + ":4317"),
			InsecureConnection:  ptr.Ptr(true),
			LogStreamingEnabled: ptr.Ptr(true),
			TraceSampleRatio:    ptr.Ptr(0.0), // Set to > 0 to enable tracing
			ResourceAttributes: map[string]string{
				"service.name":     "chainlink-node",
				"service.instance": nodeIdentifier,
				"node.don":         donMetadata.Name,
				"node.index":       strconv.Itoa(nodeMetadata.Index),
			},
		}
		// Note: OTEL_SERVICE_NAME env var should also be set on nodes to ensure
		// the service name is applied correctly. The ResourceAttributes above may
		// not override the SDK default due to OTel resource merge behavior.
		// Add to nodeset: env_vars = { OTEL_SERVICE_NAME = "chainlink-node" }

		c.Tracing.Enabled = ptr.Ptr(false) // Set to true to enable tracing
		c.Tracing.CollectorTarget = ptr.Ptr(strings.TrimPrefix(framework.HostDockerInternal(), "http://") + ":4317")
		c.Tracing.SamplingRatio = ptr.Ptr(1.0)
		c.Tracing.Mode = ptr.Ptr("unencrypted")
		c.Tracing.NodeID = ptr.Ptr(donMetadata.Name + "-node-" + strconv.Itoa(nodeMetadata.Index))
	}

	return c
}

func addBootstrapNodeConfig(
	existingConfig corechainlink.Config,
	ocrPeeringData cre.OCRPeeringData,
	commonInputs *commonInputs,
) (corechainlink.Config, error) {
	existingConfig.OCR2 = coretoml.OCR2{
		Enabled:              ptr.Ptr(true),
		DatabaseTimeout:      commonconfig.MustNewDuration(1 * time.Second),
		ContractPollInterval: commonconfig.MustNewDuration(1 * time.Second),
	}

	ocrBoostrapperLocator, ocrBErr := commontypes.NewBootstrapperLocator(ocrPeeringData.OCRBootstraperPeerID, []string{"localhost:" + strconv.Itoa(ocrPeeringData.Port)})
	if ocrBErr != nil {
		return existingConfig, errors.Wrap(ocrBErr, "failed to create OCR bootstrapper locator")
	}

	existingConfig.P2P = coretoml.P2P{
		V2: coretoml.P2PV2{
			Enabled:              ptr.Ptr(true),
			ListenAddresses:      ptr.Ptr([]string{"0.0.0.0:" + strconv.Itoa(ocrPeeringData.Port)}),
			DefaultBootstrappers: ptr.Ptr([]commontypes.BootstrapperLocator{*ocrBoostrapperLocator}),
		},
		EnableExperimentalRageP2P: ptr.Ptr(true),
	}

	if commonInputs.provider.IsDocker() {
		existingConfig.CRE.WorkflowFetcher = &coretoml.WorkflowFetcherConfig{
			URL: ptr.Ptr("file:///home/chainlink/workflows"),
		}

		existingConfig.Telemetry.ChipIngressEndpoint = ptr.Ptr(strings.TrimPrefix(framework.HostDockerInternal(), "http://") + ":" + chipingressset.DEFAULT_CHIP_INGRESS_GRPC_PORT)
		existingConfig.Telemetry.ChipIngressInsecureConnection = ptr.Ptr(true)
		existingConfig.Telemetry.HeartbeatInterval = commonconfig.MustNewDuration(30 * time.Second)

		existingConfig.Billing = coretoml.Billing{
			URL:        ptr.Ptr("billing-platform-service:2223"),
			TLSEnabled: ptr.Ptr(false),
		}
	}

	existingConfig.Capabilities = coretoml.Capabilities{
		Peering: coretoml.P2P{
			V2: coretoml.P2PV2{
				Enabled: ptr.Ptr(false),
			},
		},
		SharedPeering: coretoml.SharedPeering{
			Enabled: ptr.Ptr(true),
		},
		Dispatcher: coretoml.Dispatcher{
			SendToSharedPeer: ptr.Ptr(true),
		},
	}

	for _, evmChain := range commonInputs.evmChains {
		appendEVMChain(&existingConfig.EVM, evmChain)
	}

	if commonInputs.solanaChain != nil {
		appendSolanaChain(&existingConfig.Solana, commonInputs.solanaChain)
	}

	for _, ac := range commonInputs.aptosChains {
		existingConfig.Aptos = append(existingConfig.Aptos, corechainlink.RawConfig{
			"ChainID": ac.ChainID,
			"Enabled": true,
			"Nodes":   []map[string]any{{"Name": "default", "URL": ac.NodeURL}},
		})
	}

	// Set external registry only (local EVM capability registry). We do not set [Capabilities.Local];
	// capabilities (e.g. cron) are registered on the on-chain capability registry via Features (e.g. Cron
	// feature PreEnvStartup), same as workflow-don-solana.toml, workflow-gateway-don.toml, workflow-don-tron.toml.
	if existingConfig.Capabilities.ExternalRegistry.Address == nil {
		existingConfig.Capabilities.ExternalRegistry = coretoml.ExternalRegistry{
			Address:         ptr.Ptr(commonInputs.capabilityRegistry.address),
			NetworkID:       ptr.Ptr("evm"),
			ChainID:         ptr.Ptr(strconv.FormatUint(commonInputs.registryChainID, 10)),
			ContractVersion: ptr.Ptr(commonInputs.capabilityRegistry.version.String()),
		}
	}

	return existingConfig, nil
}

func addWorkerNodeConfig(
	existingConfig corechainlink.Config,
	topology *cre.Topology,
	ocrPeeringData cre.OCRPeeringData,
	commonInputs *commonInputs,
	donMetadata *cre.DonMetadata,
	m *cre.NodeMetadata,
) (corechainlink.Config, error) {
	ocrBoostrapperLocator, ocrBErr := commontypes.NewBootstrapperLocator(ocrPeeringData.OCRBootstraperPeerID, []string{ocrPeeringData.OCRBootstraperHost + ":" + strconv.Itoa(ocrPeeringData.Port)})
	if ocrBErr != nil {
		return existingConfig, errors.Wrap(ocrBErr, "failed to create OCR bootstrapper locator")
	}

	existingConfig.OCR2 = coretoml.OCR2{
		Enabled:              ptr.Ptr(true),
		DatabaseTimeout:      commonconfig.MustNewDuration(1 * time.Second),
		ContractPollInterval: commonconfig.MustNewDuration(1 * time.Second),
	}

	existingConfig.P2P = coretoml.P2P{
		V2: coretoml.P2PV2{
			Enabled:              ptr.Ptr(true),
			ListenAddresses:      ptr.Ptr([]string{"0.0.0.0:" + strconv.Itoa(ocrPeeringData.Port)}),
			DefaultBootstrappers: ptr.Ptr([]commontypes.BootstrapperLocator{*ocrBoostrapperLocator}),
		},
		EnableExperimentalRageP2P: ptr.Ptr(true),
	}

	if commonInputs.provider.IsDocker() {
		existingConfig.CRE.WorkflowFetcher = &coretoml.WorkflowFetcherConfig{
			URL: ptr.Ptr("file:///home/chainlink/workflows"),
		}

		existingConfig.Telemetry.ChipIngressEndpoint = ptr.Ptr(strings.TrimPrefix(framework.HostDockerInternal(), "http://") + ":" + chipingressset.DEFAULT_CHIP_INGRESS_GRPC_PORT)
		existingConfig.Telemetry.ChipIngressInsecureConnection = ptr.Ptr(true)
		existingConfig.Telemetry.HeartbeatInterval = commonconfig.MustNewDuration(30 * time.Second)

		existingConfig.Billing = coretoml.Billing{
			URL:        ptr.Ptr("billing-platform-service:2223"),
			TLSEnabled: ptr.Ptr(false),
		}
	}

	// Preserve existing WorkflowRegistry config (e.g., AdditionalSourcesConfig from user_config_overrides)
	// and Local capabilities config before resetting Capabilities struct.
	existingWorkflowRegistry := existingConfig.Capabilities.WorkflowRegistry
	existingLocalCapabilities := existingConfig.Capabilities.Local

	existingConfig.Capabilities = coretoml.Capabilities{
		Peering: coretoml.P2P{
			V2: coretoml.P2PV2{
				Enabled: ptr.Ptr(false),
			},
		},
		SharedPeering: coretoml.SharedPeering{
			Enabled: ptr.Ptr(true),
		},
		Dispatcher: coretoml.Dispatcher{
			SendToSharedPeer: ptr.Ptr(true),
		},
		WorkflowRegistry: existingWorkflowRegistry,
		Local:            existingLocalCapabilities,
	}

	if len(donMetadata.RegistryBasedLaunchAllowlist) > 0 {
		existingConfig.Capabilities.Local = coretoml.LocalCapabilities{
			RegistryBasedLaunchAllowlist: donMetadata.RegistryBasedLaunchAllowlist,
		}
	}

	for _, evmChain := range commonInputs.evmChains {
		appendEVMChain(&existingConfig.EVM, evmChain)
	}

	if commonInputs.solanaChain != nil {
		appendSolanaChain(&existingConfig.Solana, commonInputs.solanaChain)
	}

	for _, ac := range commonInputs.aptosChains {
		existingConfig.Aptos = append(existingConfig.Aptos, corechainlink.RawConfig{
			"ChainID": ac.ChainID,
			"Enabled": true,
			"Nodes":   []map[string]any{{"Name": "default", "URL": ac.NodeURL}},
		})
	}

	if existingConfig.Capabilities.ExternalRegistry.Address == nil {
		existingConfig.Capabilities.ExternalRegistry = coretoml.ExternalRegistry{
			Address:         ptr.Ptr(commonInputs.capabilityRegistry.address),
			NetworkID:       ptr.Ptr("evm"),
			ChainID:         ptr.Ptr(strconv.FormatUint(commonInputs.registryChainID, 10)),
			ContractVersion: ptr.Ptr(commonInputs.capabilityRegistry.version.String()),
		}
	}

	if donMetadata.HasFlag(cre.WorkflowDON) && existingConfig.Capabilities.WorkflowRegistry.Address == nil {
		// Preserve existing AdditionalSourcesConfig when setting WorkflowRegistry fields
		// Transform URLs to use platform-specific Docker host (handles macOS vs Linux differences)
		existingAddSources := transformAdditionalSourceURLs(existingConfig.Capabilities.WorkflowRegistry.AdditionalSourcesConfig)

		existingConfig.Capabilities.WorkflowRegistry = coretoml.WorkflowRegistry{
			Address:                 ptr.Ptr(commonInputs.workflowRegistry.address),
			NetworkID:               ptr.Ptr("evm"),
			ChainID:                 ptr.Ptr(strconv.FormatUint(commonInputs.registryChainID, 10)),
			ContractVersion:         ptr.Ptr(commonInputs.workflowRegistry.version.String()),
			SyncStrategy:            ptr.Ptr("reconciliation"),
			AdditionalSourcesConfig: existingAddSources,
		}
	}

	if donMetadata.IsShardDON() {
		existingConfig.Sharding.ShardingEnabled = ptr.Ptr(true)
		existingConfig.Sharding.ShardIndex = ptr.Ptr(uint16(donMetadata.ShardIndex)) //nolint:gosec // disable G115 overflow is unrealistic
		existingConfig.Sharding.ArbiterPort = ptr.Ptr(cre.DefaultArbiterPort)
		existingConfig.Sharding.ShardOrchestratorPort = ptr.Ptr(cre.DefaultShardOrchestratorPort)

		if !donMetadata.IsShardLeader() {
			shard0, sErr := topology.DonsMetadata.ShardLeaderDON()
			if sErr != nil {
				return existingConfig, fmt.Errorf("failed to fetch shard leader DON: %w", sErr)
			}

			if m.Index >= len(shard0.NodesMetadata) {
				return existingConfig, fmt.Errorf("shard %d node index %d exceeds shard leader node count %d", donMetadata.ShardIndex, m.Index, len(shard0.NodesMetadata))
			}

			existingConfig.Sharding.ShardOrchestratorAddress = ptr.Ptr(*commonconfig.MustParseURL(shard0.NodesMetadata[m.Index].ShardOrchestratorAddressWithPort(cre.DefaultShardOrchestratorPort)))
		}
	}

	// Add only gateway connector only to workflow DON
	// Capabilities that require gateways should add gateway connector themselves
	if donMetadata.HasFlag(cre.WorkflowDON) {
		evmKey, ok := m.Keys.EVM[commonInputs.registryChainID]
		if !ok {
			return existingConfig, fmt.Errorf("failed to get EVM key (chainID %d, node index %d)", commonInputs.registryChainID, m.Index)
		}

		gateways := []coretoml.ConnectorGateway{}
		if topology != nil && topology.GatewayConnectors != nil && len(topology.GatewayConnectors.Configurations) > 0 {
			for _, gateway := range topology.GatewayConnectors.Configurations {
				gateways = append(gateways, coretoml.ConnectorGateway{
					ID: ptr.Ptr(gateway.AuthGatewayID),
					URL: ptr.Ptr(fmt.Sprintf("ws://%s:%d%s",
						gateway.Outgoing.Host,
						gateway.Outgoing.Port,
						gateway.Outgoing.Path)),
				})
			}

			existingConfig.Capabilities.GatewayConnector = coretoml.GatewayConnector{
				DonID:             ptr.Ptr(donMetadata.Name),
				ChainIDForNodeKey: ptr.Ptr(strconv.FormatUint(commonInputs.registryChainID, 10)),
				NodeAddress:       ptr.Ptr(evmKey.PublicAddress.Hex()),
				Gateways:          gateways,
			}
		}
	}
	return existingConfig, nil
}

func addGatewayNodeConfig(
	existingConfig corechainlink.Config,
	ocrPeeringData cre.OCRPeeringData,
	commonInputs *commonInputs,
	m *cre.NodeMetadata,
) (corechainlink.Config, error) {
	// TODO: remove this in the future?
	// Unless node has Peering enabled it won't create capabilities registry syncer and all requests to vault handler will fail,
	// because it won't be able to find the DON with vault capability. P2P also required OCR2 to be enabled due to code requirements.
	// Having said that, this node will never receive any OCR2 or Peering traffic.
	existingConfig.OCR2 = coretoml.OCR2{
		Enabled:              ptr.Ptr(true),
		DatabaseTimeout:      commonconfig.MustNewDuration(1 * time.Second),
		ContractPollInterval: commonconfig.MustNewDuration(1 * time.Second),
	}

	if existingConfig.P2P.V2.Enabled == nil {
		existingConfig.P2P.V2.Enabled = ptr.Ptr(true)
		existingConfig.P2P.V2.ListenAddresses = ptr.Ptr([]string{"0.0.0.0:" + strconv.Itoa(ocrPeeringData.Port)})
	}

	existingConfig.Capabilities = coretoml.Capabilities{
		Peering: coretoml.P2P{
			V2: coretoml.P2PV2{
				Enabled: ptr.Ptr(false),
			},
		},
		SharedPeering: coretoml.SharedPeering{
			Enabled: ptr.Ptr(true),
		},
		Dispatcher: coretoml.Dispatcher{
			SendToSharedPeer: ptr.Ptr(true),
		},
	}

	for _, evmChain := range commonInputs.evmChains {
		appendEVMChain(&existingConfig.EVM, evmChain)
	}

	if existingConfig.Capabilities.ExternalRegistry.Address == nil {
		existingConfig.Capabilities.ExternalRegistry = coretoml.ExternalRegistry{
			Address:         ptr.Ptr(commonInputs.capabilityRegistry.address),
			NetworkID:       ptr.Ptr("evm"),
			ChainID:         ptr.Ptr(strconv.FormatUint(commonInputs.registryChainID, 10)),
			ContractVersion: ptr.Ptr(commonInputs.capabilityRegistry.version.String()),
		}
	}

	if existingConfig.Capabilities.WorkflowRegistry.Address == nil {
		existingConfig.Capabilities.WorkflowRegistry = coretoml.WorkflowRegistry{
			Address:         ptr.Ptr(commonInputs.workflowRegistry.address),
			NetworkID:       ptr.Ptr("evm"),
			ChainID:         ptr.Ptr(strconv.FormatUint(commonInputs.registryChainID, 10)),
			ContractVersion: ptr.Ptr(commonInputs.workflowRegistry.version.String()),
			SyncStrategy:    ptr.Ptr("reconciliation"),
		}
	}

	if existingConfig.Capabilities.WorkflowRegistry.WorkflowStorage.URL == nil {
		// TODO remove WorkflowStorage once it is not required on a gateway node
		existingConfig.Capabilities.WorkflowRegistry.WorkflowStorage = coretoml.WorkflowStorage{
			URL:                 ptr.Ptr("localhost"),
			TLSEnabled:          ptr.Ptr(false),
			ArtifactStorageHost: ptr.Ptr("localhost"),
		}
	}

	if len(existingConfig.Capabilities.GatewayConnector.Gateways) == 0 {
		existingConfig.Capabilities.GatewayConnector = coretoml.GatewayConnector{
			DonID:             ptr.Ptr("doesn't-matter-for-gateway-node"),
			ChainIDForNodeKey: ptr.Ptr(strconv.FormatUint(commonInputs.registryChainID, 10)),
		}
	}

	return existingConfig, nil
}

type versionedAddress struct {
	address string
	version *semver.Version
}

type aptosChain struct {
	ChainID string
	NodeURL string
}

type commonInputs struct {
	registryChainID       uint64
	registryChainSelector uint64

	workflowRegistry   versionedAddress
	capabilityRegistry versionedAddress

	evmChains   []*evmChain
	solanaChain *solanaChain
	aptosChains []*aptosChain

	provider infra.Provider
}

func gatherCommonInputs(input cre.GenerateConfigsInput) (*commonInputs, error) {
	registryChainID, homeErr := chain_selectors.ChainIdFromSelector(input.RegistryChainSelector)
	if homeErr != nil {
		return nil, errors.Wrap(homeErr, "failed to get home chain ID")
	}

	evmChains := findEVMChains(input)
	solanaChain, solErr := findOneSolanaChain(input)
	if solErr != nil {
		return nil, errors.Wrap(solErr, "failed to find Solana chain in the environment configuration")
	}

	capabilitiesRegistryAddress := crecontracts.MustGetAddressFromDataStore(input.Datastore, input.RegistryChainSelector, keystone_changeset.CapabilitiesRegistry.String(), input.ContractVersions[keystone_changeset.CapabilitiesRegistry.String()], "")
	workflowRegistryAddress := crecontracts.MustGetAddressFromDataStore(input.Datastore, input.RegistryChainSelector, keystone_changeset.WorkflowRegistry.String(), input.ContractVersions[keystone_changeset.WorkflowRegistry.String()], "")

	aptosChains, aptosErr := findAptosChains(input)
	if aptosErr != nil {
		return nil, errors.Wrap(aptosErr, "failed to find Aptos chains in the environment configuration")
	}

	return &commonInputs{
		registryChainID:       registryChainID,
		registryChainSelector: input.RegistryChainSelector,
		workflowRegistry: versionedAddress{
			address: workflowRegistryAddress,
			version: input.ContractVersions[keystone_changeset.WorkflowRegistry.String()],
		},
		evmChains:   evmChains,
		solanaChain: solanaChain,
		aptosChains: aptosChains,
		capabilityRegistry: versionedAddress{
			address: capabilitiesRegistryAddress,
			version: input.ContractVersions[keystone_changeset.CapabilitiesRegistry.String()],
		},
		provider: input.Provider,
	}, nil
}

type evmChain struct {
	Name    string
	ChainID uint64
	HTTPRPC string
	WSRPC   string
}

func findEVMChains(input cre.GenerateConfigsInput) []*evmChain {
	evmChains := make([]*evmChain, 0)
	for _, bcOut := range input.Blockchains {
		if bcOut.IsFamily(chain_selectors.FamilySolana) || bcOut.IsFamily(chain_selectors.FamilyAptos) {
			continue
		}

		// if the DON doesn't support the chain, we skip it; if slice is empty, it means that the DON supports all chains
		if len(input.DonMetadata.MustNodeSet().EVMChains()) > 0 && !slices.Contains(input.DonMetadata.MustNodeSet().EVMChains(), bcOut.ChainID()) {
			continue
		}

		evmChains = append(evmChains, &evmChain{
			Name:    fmt.Sprintf("node-%d", bcOut.ChainSelector()),
			ChainID: bcOut.ChainID(),
			HTTPRPC: bcOut.CtfOutput().Nodes[0].InternalHTTPUrl,
			WSRPC:   bcOut.CtfOutput().Nodes[0].InternalWSUrl,
		})
	}
	return evmChains
}

type solanaChain struct {
	Name    string
	ChainID string
	NodeURL string
}

func findOneSolanaChain(input cre.GenerateConfigsInput) (*solanaChain, error) {
	var solChain *solanaChain
	chainsFound := 0

	for _, bcOut := range input.Blockchains {
		if !bcOut.IsFamily(chain_selectors.FamilySolana) {
			continue
		}

		chainsFound++
		if chainsFound > 1 {
			return nil, errors.New("multiple Solana chains found, expected only one")
		}

		solBc := bcOut.(*solana.Blockchain)

		ctx, cancelFn := context.WithTimeout(context.Background(), 15*time.Second)
		chainID, err := solBc.GenesisHash(ctx)
		if err != nil {
			cancelFn()
			return nil, errors.Wrap(err, "failed to get chainID for Solana")
		}
		cancelFn()

		solChain = &solanaChain{
			Name:    fmt.Sprintf("node-%d", solBc.ChainSelector()),
			ChainID: chainID,
			NodeURL: bcOut.CtfOutput().Nodes[0].InternalHTTPUrl,
		}
	}

	return solChain, nil
}

func findAptosChains(input cre.GenerateConfigsInput) ([]*aptosChain, error) {
	capabilityChainIDs := input.DonMetadata.MustNodeSet().ChainCapabilityChainIDs()
	out := make([]*aptosChain, 0)
	for _, bcOut := range input.Blockchains {
		if !bcOut.IsFamily(chain_selectors.FamilyAptos) {
			continue
		}
		if len(capabilityChainIDs) > 0 && !slices.Contains(capabilityChainIDs, bcOut.ChainID()) {
			continue
		}

		aptosBC, ok := bcOut.(*aptoschain.Blockchain)
		if !ok {
			return nil, fmt.Errorf("expected Aptos blockchain implementation for chain %d, got %T", bcOut.ChainID(), bcOut)
		}
		nodeURL, err := aptosBC.InternalNodeURL()
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get Aptos internal node URL for chain %d", bcOut.ChainID())
		}

		out = append(out, &aptosChain{
			ChainID: strconv.FormatUint(bcOut.ChainID(), 10),
			NodeURL: nodeURL,
		})
	}
	return out, nil
}

func buildTronEVMConfig(evmChain *evmChain) evmconfigtoml.EVMConfig {
	tronRPC := strings.Replace(evmChain.HTTPRPC, "jsonrpc", "wallet", 1)
	return evmconfigtoml.EVMConfig{
		ChainID: sqlutil.New(big.NewInt(libc.MustSafeInt64(evmChain.ChainID))),
		Chain: evmconfigtoml.Chain{
			AutoCreateKey:         ptr.Ptr(false),
			ChainType:             chaintype.NewConfig("tron"),
			LogBroadcasterEnabled: ptr.Ptr(false),
			NodePool: evmconfigtoml.NodePool{
				NewHeadsPollInterval: commonconfig.MustNewDuration(10 * time.Second),
			},
		},
		Nodes: []*evmconfigtoml.Node{
			{
				Name:              ptr.Ptr(evmChain.Name),
				HTTPURL:           commonconfig.MustParseURL(evmChain.HTTPRPC),
				HTTPURLExtraWrite: commonconfig.MustParseURL(tronRPC),
			},
		},
	}
}

func buildEVMConfig(evmChain *evmChain) evmconfigtoml.EVMConfig {
	return evmconfigtoml.EVMConfig{
		ChainID: sqlutil.New(big.NewInt(libc.MustSafeInt64(evmChain.ChainID))),
		Chain: evmconfigtoml.Chain{
			AutoCreateKey: ptr.Ptr(false),
		},
		Nodes: []*evmconfigtoml.Node{
			{
				Name:    ptr.Ptr(evmChain.Name),
				WSURL:   commonconfig.MustParseURL(evmChain.WSRPC),
				HTTPURL: commonconfig.MustParseURL(evmChain.HTTPRPC),
			},
		},
	}
}

func appendEVMChain(existingConfig *evmconfigtoml.EVMConfigs, evmChain *evmChain) {
	var cfg evmconfigtoml.EVMConfig
	if evmChain.ChainID == TronEVMChainID {
		cfg = buildTronEVMConfig(evmChain)
	} else {
		cfg = buildEVMConfig(evmChain)
	}

	// add only unconfigured chains, since other roles might have already added some chains
	for _, existingEVM := range *existingConfig {
		if existingEVM.ChainID.ToInt().Cmp(big.NewInt(libc.MustSafeInt64(evmChain.ChainID))) == 0 {
			return
		}
	}

	*existingConfig = append(*existingConfig, &cfg)
}

func appendSolanaChain(existingConfig *solcfg.TOMLConfigs, solChain *solanaChain) {
	for _, existingSol := range *existingConfig {
		if existingSol.ChainID != nil && *existingSol.ChainID == solChain.ChainID {
			return
		}
	}

	*existingConfig = append(*existingConfig, &solcfg.TOMLConfig{
		Enabled: ptr.Ptr(true),
		ChainID: ptr.Ptr(solChain.ChainID),
		Nodes: []*solcfg.Node{
			{
				Name: &solChain.Name,
				URL:  commonconfig.MustParseURL(solChain.NodeURL),
			},
		},
	})
}

// transformAdditionalSourceURLs transforms URLs in AdditionalSourcesConfig to use
// platform-specific Docker host addresses. This handles differences between macOS
// (host.docker.internal) and Linux (172.17.0.1 or similar) Docker host resolution.
func transformAdditionalSourceURLs(sources []coretoml.AdditionalWorkflowSource) []coretoml.AdditionalWorkflowSource {
	if len(sources) == 0 {
		return sources
	}

	// Get the platform-specific Docker host (e.g., "http://host.docker.internal" on macOS,
	// "http://172.17.0.1" on Linux)
	dockerHost := strings.TrimPrefix(framework.HostDockerInternal(), "http://")

	transformed := make([]coretoml.AdditionalWorkflowSource, len(sources))
	for i, src := range sources {
		transformed[i] = src
		if src.URL != nil {
			// Replace "host.docker.internal" with the platform-specific host
			url := *src.URL
			url = strings.Replace(url, "host.docker.internal", dockerHost, 1)
			transformed[i].URL = &url
		}
	}

	return transformed
}

// transformUserConfigOverrides transforms URLs in a user config overrides string to use
// platform-specific Docker host addresses. This handles differences between macOS
// (host.docker.internal) and Linux (172.17.0.1 or similar) Docker host resolution.
// This is necessary because UserConfigOverrides is passed directly to containers as a
// separate config file, bypassing the structured config transformation.
func transformUserConfigOverrides(userConfig string) string {
	if userConfig == "" {
		return userConfig
	}

	// Get the platform-specific Docker host (e.g., "http://host.docker.internal" on macOS,
	// "http://172.17.0.1" on Linux)
	dockerHost := strings.TrimPrefix(framework.HostDockerInternal(), "http://")

	// Replace all occurrences of "host.docker.internal" with the platform-specific host
	return strings.ReplaceAll(userConfig, "host.docker.internal", dockerHost)
}

// generateInstanceNames creates Kubernetes-compatible instance names for nodes
// Bootstrap nodes get names like "workflow-bt-0", plugin nodes get "workflow-0", "workflow-1", etc.
// This is a wrapper around infra.GenerateNodeInstanceNames that converts NodeMetadata to bool roles
func generateInstanceNames(nodeSetName string, nodesMetadata []*cre.NodeMetadata) []string {
	// Convert NodeMetadata to bool slice indicating bootstrap roles
	nodeMetadataRoles := make([]bool, len(nodesMetadata))
	for i, nodeMetadata := range nodesMetadata {
		nodeMetadataRoles[i] = nodeMetadata.HasRole(cre.BootstrapNode)
	}

	return infra.GenerateNodeInstanceNames(nodeSetName, nodeMetadataRoles)
}
