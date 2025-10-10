package webapitrigger

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
	coregateway "github.com/smartcontractkit/chainlink/v2/core/services/gateway"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/gateway"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	factory "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability/donlevel"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
)

const flag = cre.WebAPITriggerCapability

type WebAPITrigger struct{}

func (o *WebAPITrigger) Flag() cre.CapabilityFlag {
	return flag
}

func (o *WebAPITrigger) PreEnvStartup(
	testLogger zerolog.Logger,
	registryChainSelector uint64,
	cldfEnv *cldf.Environment,
	provider infra.Provider,
	topology *cre.Topology,
	blockchainOutputs []*cre.WrappedBlockchainOutput,
	capabilityConfigs cre.CapabilityConfigs, // move to Topology
	contractVersions map[string]string,
	gatewayJobConfigs map[cre.NodeUUID]*config.GatewayConfig,
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

	// add 'web-api' handler to gateway config (future jobspec)
	// add gateway connector to to node TOML config, so that node can route http requests to the gateway
	for idx, donMetadata := range donsMetadata {
		handlerConfig, confErr := gateway.HandlerConfig(coregateway.WebAPICapabilitiesType)
		if confErr != nil {
			return nil, errors.Wrapf(confErr, "failed to get %s handler config for don %s", coregateway.WebAPICapabilitiesType, donMetadata.Name)
		}
		hErr := gateway.AddHandlers(donMetadata, registryChainID, gatewayJobConfigs, []config.Handler{handlerConfig})
		if hErr != nil {
			return nil, errors.Wrapf(hErr, "failed to add gateway handlers to gateway config (jobspec) for don %s ", donMetadata.Name)
		}

		cErr := gateway.AddConnectors(donMetadata, registryChainID, topology.GatewayConnectorOutput)
		if cErr != nil {
			return nil, errors.Wrapf(cErr, "failed to add gateway connectors to node's TOML config in for don %s", donMetadata.Name)
		}

		donsMetadata[idx] = donMetadata
	}

	capabilities := make(map[uint64][]keystone_changeset.DONCapabilityWithConfig)
	for _, donMetadata := range donsMetadata {
		if capabilities[donMetadata.ID] == nil {
			capabilities[donMetadata.ID] = []keystone_changeset.DONCapabilityWithConfig{}
		}

		capabilities[donMetadata.ID] = append(capabilities[donMetadata.ID], keystone_changeset.DONCapabilityWithConfig{
			Capability: kcr.CapabilitiesRegistryCapability{
				LabelledName:   "web-api-trigger",
				Version:        "1.0.0",
				CapabilityType: 0, // TRIGGER
			},
			Config: &capabilitiespb.CapabilityConfig{},
		})
	}

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfigs: capabilities,
		GatewayJobConfigs:        gatewayJobConfigs,
	}, nil
}

const configTemplate = `""`

func (o *WebAPITrigger) PostEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	creEnv *cre.Environment,
	nodeSetOutput []*cre.WrappedNodeOutput,
	blockchainOutputs []*cre.WrappedBlockchainOutput,
	contractVersions map[string]string,
	provider infra.Provider,
	capabilityConfigs map[string]cre.CapabilityConfig,
) error {
	dons := creEnv.DonTopology.DonsWithFlag(flag)
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
		configTemplate,
		factory.NoOpExtractor, // No runtime values extraction needed
		func(_ *cre.JobSpecInput, _ cre.CapabilityConfig) (string, error) {
			return "__builtin_web-api-trigger", nil
		},
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
		// pass whole topology, since some jobs might need to be created on multiple DONs
		jobErr := jobs.Create(ctx, creEnv.CldfEnvironment.Offchain, creEnv.DonTopology, jobSpecs)
		if jobErr != nil {
			return fmt.Errorf("failed to create http action jobs for don %s: %w", don.Name, jobErr)
		}
	}

	return nil
}
