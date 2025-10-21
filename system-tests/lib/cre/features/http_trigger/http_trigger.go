package httptrigger

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

const flag = cre.HTTPTriggerCapability

type HTTPTrigger struct{}

func (o *HTTPTrigger) Flag() cre.CapabilityFlag {
	return flag
}

func (o *HTTPTrigger) PreEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.DonMetadata,
	topology *cre.Topology,
	creEnv *cre.Environment,
) (*cre.PreEnvStartupOutput, error) {
	// use registry chain, because that is the chain we used when generating gateway connector part of node config (check below)
	registryChainID, chErr := chainselectors.ChainIdFromSelector(creEnv.RegistryChainSelector)
	if chErr != nil {
		return nil, errors.Wrapf(chErr, "failed to get chain ID from selector %d", creEnv.RegistryChainSelector)
	}

	// add 'http-capabilities' handler to gateway config (future jobspec)
	// add gateway connector to to node TOML config, so that node can route http trigger requests to the gateway
	handlerConfig, confErr := gateway.HandlerConfig(coregateway.HTTPCapabilityType)
	if confErr != nil {
		return nil, errors.Wrapf(confErr, "failed to get %s handler config for don %s", coregateway.HTTPCapabilityType, don.Name)
	}
	hErr := gateway.AddHandlers(*don, registryChainID, topology.GatewayJobConfigs, []config.Handler{handlerConfig})
	if hErr != nil {
		return nil, errors.Wrapf(hErr, "failed to add gateway handlers to gateway config (jobspec) for don %s ", don.Name)
	}

	cErr := gateway.AddConnectors(don, registryChainID, *topology.GatewayConnectors)
	if cErr != nil {
		return nil, errors.Wrapf(cErr, "failed to add gateway connectors to node's TOML config in for don %s", don.Name)
	}

	capabilities := []keystone_changeset.DONCapabilityWithConfig{{
		Capability: kcr.CapabilitiesRegistryCapability{
			LabelledName:   "http-trigger",
			Version:        "1.0.0-alpha",
			CapabilityType: 0, // TRIGGER
		},
		Config: &capabilitiespb.CapabilityConfig{},
	}}

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfig: capabilities,
	}, nil
}

const configTemplate = `"""
{
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

func (o *HTTPTrigger) PostEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.Don,
	dons *cre.Dons,
	creEnv *cre.Environment,
) error {
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

	var nodeSet cre.NodeSetWithCapabilityConfigs
	for _, ns := range dons.AsNodeSetWithChainCapabilities() {
		if ns.GetName() == don.Name {
			nodeSet = ns
			break
		}
	}
	if nodeSet == nil {
		return fmt.Errorf("could not find node set for Don named '%s'", don.Name)
	}

	jobSpecs, specErr := perDonJobSpecFactory.BuildJobSpec(
		flag,
		configTemplate,
		factory.NoOpExtractor,
		factory.BinaryPathBuilder,
	)(&cre.JobSpecInput{
		CreEnvironment: creEnv,
		Don:            don,
		NodeSet:        nodeSet,
	})
	if specErr != nil {
		return fmt.Errorf("failed to build job spec for http action capability: %w", specErr)
	}
	if len(jobSpecs) == 0 {
		return fmt.Errorf("no job specs created for '%s' capability, even though it is enabled", flag)
	}

	// pass all dons, since some jobs might need to be created on multiple ones
	jobErr := jobs.Create(ctx, creEnv.CldfEnvironment.Offchain, dons, jobSpecs)
	if jobErr != nil {
		return fmt.Errorf("failed to create http action jobs for don %s: %w", don.Name, jobErr)
	}

	return nil
}
