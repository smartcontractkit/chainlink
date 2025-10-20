package customcompute

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	coregateway "github.com/smartcontractkit/chainlink/v2/core/services/gateway"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/gateway"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	factory "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability/donlevel"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
)

const flag = cre.CustomComputeCapability

type CustomCompute struct{}

func (o *CustomCompute) Flag() cre.CapabilityFlag {
	return flag
}

func (o *CustomCompute) PreEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	topology *cre.Topology,
	creEnv *cre.Environment,
) (*cre.PreEnvStartupOutput, error) {
	donsMetadata := topology.DonsMetadataWithFlag(flag)
	if len(donsMetadata) == 0 {
		return nil, nil
	}

	// use registry chain, because that is the chain we used when generating gateway connector part of node config (check below)
	registryChainID, chErr := chainselectors.ChainIdFromSelector(creEnv.RegistryChainSelector)
	if chErr != nil {
		return nil, errors.Wrapf(chErr, "failed to get chain ID from selector %d", creEnv.RegistryChainSelector)
	}

	// add 'web-api' handler to gateway config (future jobspec)
	// add gateway connector to to node TOML config, so that node can route http requests to the gateway
	for idx, donMetadata := range donsMetadata {
		handlerConfig, confErr := gateway.HandlerConfig(coregateway.WebAPICapabilitiesType)
		if confErr != nil {
			return nil, errors.Wrapf(confErr, "failed to get %s handler config for don %s", coregateway.WebAPICapabilitiesType, donMetadata.Name)
		}
		hErr := gateway.AddHandlers(donMetadata, registryChainID, topology.GatewayJobConfigs, []config.Handler{handlerConfig})
		if hErr != nil {
			return nil, errors.Wrapf(hErr, "failed to add gateway handlers to gateway config (jobspec) for don %s ", donMetadata.Name)
		}

		cErr := gateway.AddConnectors(donMetadata, registryChainID, topology.GatewayConnectors)
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
				LabelledName:   "custom-compute",
				Version:        "1.0.0",
				CapabilityType: 1, // ACTION
			},
			Config: &capabilitiespb.CapabilityConfig{},
		})
	}

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfigs: capabilities,
		GatewayJobConfigs:        topology.GatewayJobConfigs,
	}, nil
}

const configTemplate = `"""
NumWorkers = {{.NumWorkers}}
[rateLimiter]
globalRPS = {{.GlobalRPS}}
globalBurst = {{.GlobalBurst}}
perSenderRPS = {{.PerSenderRPS}}
perSenderBurst = {{.PerSenderBurst}}
"""`

func (o *CustomCompute) PostEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	donTopology *cre.DonTopology,
	creEnv *cre.Environment,
) error {
	dons := donTopology.DonsWithFlag(flag)
	if len(dons) == 0 {
		return nil
	}

	perDonJobSpecFactory, fErr := factory.NewCapabilityJobSpecFactory(
		creEnv.RegistryChainSelector,
		donlevel.CapabilityEnabler,
		donlevel.EnabledChainsProvider,
		donlevel.ConfigResolver,
		donlevel.JobNamer,
	)

	if fErr != nil {
		return errors.Wrap(fErr, "failed to create capability job spec factory")
	}

	bcOuts := make([]*blockchain.Output, len(creEnv.Blockchains))
	for i, b := range creEnv.Blockchains {
		bcOuts[i] = b.CtfOutput()
	}

	donsToJobSpecs, specErr := perDonJobSpecFactory.BuildJobSpec(
		flag,
		configTemplate,
		factory.NoOpExtractor,
		func(_ *cre.JobSpecInput, _ cre.CapabilityConfig) (string, error) {
			return "__builtin_custom-compute-action", nil
		},
	)(&cre.JobSpecInput{
		CreEnvironment: creEnv,
		DonTopology:    donTopology,
		NodeSets:       donTopology.Dons.AsNodeSetWithChainCapabilities(),
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
		jobErr := jobs.Create(ctx, creEnv.CldfEnvironment.Offchain, donTopology, jobSpecs)
		if jobErr != nil {
			return fmt.Errorf("failed to create http action jobs for don %s: %w", don.Name, jobErr)
		}
	}

	return nil
}
