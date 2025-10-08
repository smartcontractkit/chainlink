package httpaction

import (
	"context"
	"fmt"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	coretoml "github.com/smartcontractkit/chainlink/v2/core/config/toml"
	corechainlink "github.com/smartcontractkit/chainlink/v2/core/services/chainlink"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	factory "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability/donlevel"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
)

const flag = cre.HTTPActionCapability

type HTTPAction struct{}

func (o *HTTPAction) Flag() cre.CapabilityFlag {
	return flag
}

func (o *HTTPAction) PreEnvStartup(
	testLogger zerolog.Logger,
	registryChainSelector uint64,
	cldfEnv *cldf.Environment,
	provider infra.Provider,
	topology *cre.Topology,
	blockchainOutputs []*cre.WrappedBlockchainOutput,
	capabilityConfigs cre.CapabilityConfigs, // move to Topology
	contractVersions map[string]string,
	gatewayConfigs map[cre.NodeUUID]*config.GatewayConfig,
) (*cre.PreEnvStartupOutput, error) {
	donsMetadata := topology.DonsMetadataWithFlag(o.Flag())
	if len(donsMetadata) == 0 {
		return nil, nil
	}

	// use registry chain, because that is the chain we used when generating gateway connector part of node config (check below)
	registryChainID, chErr := chainselectors.ChainIdFromSelector(registryChainSelector)
	if chErr != nil {
		return nil, errors.Wrapf(chErr, "failed to get chain ID from selector %d", registryChainSelector)
	}

	// TODO export in a function, http-trigger will also need it
	// add gateway handler config to each DON that has this capability
	for _, donMetadata := range donsMetadata {
		workers, wErr := donMetadata.Workers()
		if wErr != nil {
			return nil, wErr
		}
		evmKey, ok := workers[0].Keys.EVM[registryChainID]
		if !ok {
			return nil, fmt.Errorf("worker node at index %d does not have EVM key for chainID %d", workers[0].Index, registryChainID)
		}

		// for each DON, we need to add a handler config specific for this capability
		for nodeUUID, gc := range gatewayConfigs {
			donFound := false
			for donIdx, maybeDON := range gc.Dons {
				// first we try to find DON configuration that matches current don, because it might be already present
				for _, member := range maybeDON.Members {
					// if any of the member's address matches the EVM key of the worker node, we found the right DON
					if member.Address == evmKey.PublicAddress.Hex() {
						donFound = true
						break
					}
				}

				if donFound {
					gc.Dons[donIdx].Handlers = append(gc.Dons[donIdx].Handlers, gatewayHandlerConfig())
					break
				}
			}

			// if we did not find the DON in the gateway config, we need to add it
			if !donFound {
				members := make([]config.NodeConfig, len(workers))
				for i, worker := range workers {
					evmKey, ok := worker.Keys.EVM[registryChainID]
					if !ok {
						return nil, fmt.Errorf("worker node at index %d does not have EVM key for chain ID %d", worker.Index, registryChainID)
					}

					members[i] = config.NodeConfig{
						Address: evmKey.PublicAddress.Hex(),
						Name:    fmt.Sprintf("%s-node-%d", donMetadata.Name, worker.Index),
					}
				}

				gatewayConfigs[nodeUUID].Dons = append(gatewayConfigs[nodeUUID].Dons, config.DONConfig{
					DonId:    donMetadata.Name,
					F:        1,
					Members:  members,
					Handlers: []config.Handler{gatewayHandlerConfig()},
				})
			}
		}
	}

	// TODO export in a function, http-trigger will also need it
	// add gateway connector configuration to node TOML config
	for _, donMetadata := range donsMetadata {
		workers, wErr := donMetadata.Workers()
		if wErr != nil {
			return nil, wErr
		}

		for _, workerNode := range workers {
			currentConfig := donMetadata.CapabilitiesAwareNodeSet().NodeSpecs[workerNode.Index].Node.TestConfigOverrides

			var typedConfig corechainlink.Config
			unmarshallErr := toml.Unmarshal([]byte(currentConfig), &typedConfig)
			if unmarshallErr != nil {
				return nil, errors.Wrapf(unmarshallErr, "failed to unmarshal config for node index %d", workerNode.Index)
			}

			evmKey, ok := workerNode.Keys.EVM[registryChainID]
			if !ok {
				return nil, fmt.Errorf("failed to get EVM key (chainID %d, node index %d)", registryChainID, workerNode.Index)
			}

			// if no gateways are configured, then gateway connector config is most probably also not configured
			if len(typedConfig.Capabilities.GatewayConnector.Gateways) == 0 {
				typedConfig.Capabilities.GatewayConnector = coretoml.GatewayConnector{
					DonID:             ptr.Ptr(donMetadata.Name),
					ChainIDForNodeKey: ptr.Ptr(strconv.FormatUint(registryChainID, 10)),
					NodeAddress:       ptr.Ptr(evmKey.PublicAddress.Hex()),
				}
			}

			// make sure that all other gateways are also present in the config
			for _, gatewayConnector := range topology.GatewayConnectorOutput.Configurations {
				alreadyPresent := false
				for _, existingGateway := range typedConfig.Capabilities.GatewayConnector.Gateways {
					if gatewayConnector.AuthGatewayID == *existingGateway.ID {
						alreadyPresent = true
						continue
					}
				}

				if !alreadyPresent {
					typedConfig.Capabilities.GatewayConnector.Gateways = append(typedConfig.Capabilities.GatewayConnector.Gateways, coretoml.ConnectorGateway{
						ID: ptr.Ptr(gatewayConnector.AuthGatewayID),
						URL: ptr.Ptr(fmt.Sprintf("ws://%s:%d%s",
							gatewayConnector.Outgoing.Host,
							gatewayConnector.Outgoing.Port,
							gatewayConnector.Outgoing.Path)),
					})
				}
			}

			stringifiedConfig, mErr := toml.Marshal(typedConfig)
			if mErr != nil {
				return nil, errors.Wrapf(mErr, "failed to marshal config for node index %d", workerNode.Index)
			}

			donMetadata.CapabilitiesAwareNodeSet().NodeSpecs[workerNode.Index].Node.TestConfigOverrides = string(stringifiedConfig)
		}
	}

	capabilities := make(map[int][]keystone_changeset.DONCapabilityWithConfig)
	for donIdx := range donsMetadata {
		if capabilities[donIdx] == nil {
			capabilities[donIdx] = []keystone_changeset.DONCapabilityWithConfig{}
		}

		capabilities[donIdx] = append(capabilities[donIdx], keystone_changeset.DONCapabilityWithConfig{
			Capability: kcr.CapabilitiesRegistryCapability{
				LabelledName:   "http-actions",
				Version:        "1.0.0-alpha",
				CapabilityType: 1, // ACTION
			},
			Config: &capabilitiespb.CapabilityConfig{},
		})
	}

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfigs: capabilities,
		GatewayConfigs:           gatewayConfigs,
	}, nil
}

func gatewayHandlerConfig() config.Handler {
	return config.Handler{
		Name:        "http-capabilities",
		ServiceName: "workflows",
		Config: []byte(`
maxTriggerRequestDurationMs = 5_000
metadataPullIntervalMs = 1_000
metadataAggregationIntervalMs = 1_000
[NodeRateLimiter]
globalBurst = 10
globalRPS = 50
perSenderBurst = 10
perSenderRPS = 10`),
	}
}

const httpActionConfigTemplate = `"""
{
	"proxyMode": "{{.ProxyMode}}",
	"incomingRateLimiter": {
		"globalBurst": {{.IncomingGlobalBurst}},
		"globalRPS": {{.IncomingGlobalRPS}},
		"perSenderBurst": {{.IncomingPerSenderBurst}},
		"perSenderRPS": {{.IncomingPerSenderRPS}}
	},
	"outgoingRateLimiter": {
		"globalBurst": {{.OutgoingGlobalBurst}},
		"globalRPS": {{.OutgoingGlobalRPS}},
		"perSenderBurst": {{.OutgoingPerSenderBurst}},
		"perSenderRPS": {{.OutgoingPerSenderRPS}}
	}
}
"""`

func (o *HTTPAction) PostEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	creEnv *cre.Environment,
	nodeSetOutput []*cre.WrappedNodeOutput,
	blockchainOutputs []*cre.WrappedBlockchainOutput,
	contractVersions map[string]string,
	provider infra.Provider,
	capabilityConfigs map[string]cre.CapabilityConfig,
) error {
	dons := creEnv.DonTopology.DonWithFlag(flag)
	if len(dons) == 0 {
		return nil
	}

	perDonJobSpecFactory, fErr := factory.NewCapabilityJobSpecFactory(
		donlevel.CapabilityEnabler,
		donlevel.EnabledChainsProvider,
		donlevel.ConfigResolver,
		donlevel.JobNamer,
	)
	if fErr != nil {
		return errors.Wrap(fErr, "failed to create capability job spec factory")
	}

	bcOuts := make([]*blockchain.Output, len(blockchainOutputs))
	for i, b := range blockchainOutputs {
		bcOuts[i] = b.BlockchainOutput
	}

	donsToJobSpecs, specErr := perDonJobSpecFactory.BuildJobSpec(
		flag,
		httpActionConfigTemplate,
		factory.NoOpExtractor,
		factory.BinaryPathBuilder,
	)(&cre.JobSpecInput{
		CldEnvironment: creEnv.CldfEnvironment,
		DonTopology:    creEnv.DonTopology,
		InfraInput:     provider,
		/// Capabilities:  // not needed,
		NodeSets:          creEnv.DonTopology.Dons.AsNodeSetWithChainCapabilities(),
		CapabilityConfigs: capabilityConfigs,
	})
	if specErr != nil {
		return fmt.Errorf("failed to build job spec for http action capability: %w", specErr)
	}

	for _, don := range dons {
		jobSpecs, ok := donsToJobSpecs[don.ID]
		if !ok {
			continue
		}
		jobErr := jobs.CreateForDON(ctx, creEnv.CldfEnvironment.Offchain, don, jobSpecs)
		if jobErr != nil {
			return fmt.Errorf("failed to create http action jobs for don %s: %w", don.Name, jobErr)
		}
	}

	return nil
}
