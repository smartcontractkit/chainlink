package httpaction

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

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
	gatewayConfigs map[cre.NodeUUID]config.GatewayConfig,
) (*cre.PreEnvStartupOutput, error) {
	donsMetadata := topology.DonsMetadataWithFlag(o.Flag())
	if len(donsMetadata) == 0 {
		return nil, nil
	}

	chainID, chErr := chainselectors.ChainIdFromSelector(registryChainSelector)
	if chErr != nil {
		return nil, errors.Wrapf(chErr, "failed to get chain ID from selector %d", registryChainSelector)
	}

	for _, don := range donsMetadata {
		workers, wErr := don.Workers()
		if wErr != nil {
			return nil, wErr
		}
		// use registry chain, because that's chain we use when generating gateway connector part of node config
		evmKey, ok := workers[0].Keys.EVM[chainID]
		if !ok {
			return nil, fmt.Errorf("worker node at index %d does not have EVM key for chainID %d", workers[0].Index, chainID)
		}

		// for each DON, we need to add a handler config specific for this capability
		for _, gc := range gatewayConfigs {
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
					gc.Dons[donIdx].Handlers = append(gc.Dons[donIdx].Handlers, config.Handler{
						Name:        "http-capabilities",
						ServiceName: "workflows",
						// TODO - figure out correct json syntax
						Config: []byte(`
maxTriggerRequestDurationMs = 5_000
metadataPullIntervalMs = 1_000
metadataAggregationIntervalMs = 1_000
[NodeRateLimiter]
globalBurst = 10
globalRPS = 50
perSenderBurst = 10
perSenderRPS = 10`),
					})

					break
				}
			}

			// TODO make web api gateway also a feature! <--------
			if !donFound {
				members := make([]config.NodeConfig, len(workers))
				for i, worker := range workers {
					evmKey, ok := worker.Keys.EVM[registryChainSelector]
					if !ok {
						return nil, fmt.Errorf("worker node at index %d does not have EVM key for chain selector %d", worker.Index, registryChainSelector)
					}

					members[i] = config.NodeConfig{
						Address: evmKey.PublicAddress.Hex(),
						Name:    fmt.Sprintf("node-%d", worker.Index),
					}
				}
				gc.Dons = append(gc.Dons, config.DONConfig{
					DonId:   don.Name,
					F:       1,
					Members: members,
					Handlers: []config.Handler{
						{
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
						},
					},
				})
			}
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
