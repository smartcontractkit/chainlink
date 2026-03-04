package confidentialrelay

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

const flag = cre.ConfidentialRelayCapability

type ConfidentialRelay struct{}

func (o *ConfidentialRelay) Flag() cre.CapabilityFlag {
	return flag
}

func (o *ConfidentialRelay) PreEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.DonMetadata,
	topology *cre.Topology,
	creEnv *cre.Environment,
) (*cre.PreEnvStartupOutput, error) {
	registryChainID, chErr := chainselectors.ChainIdFromSelector(creEnv.RegistryChainSelector)
	if chErr != nil {
		return nil, errors.Wrapf(chErr, "failed to get chain ID from selector %d", creEnv.RegistryChainSelector)
	}

	hErr := topology.AddGatewayHandlers(*don, []string{pkg.GatewayHandlerTypeConfidentialRelay})
	if hErr != nil {
		return nil, errors.Wrapf(hErr, "failed to add gateway handlers to gateway config for don %s", don.Name)
	}

	cErr := don.ConfigureForGatewayAccess(registryChainID, *topology.GatewayConnectors)
	if cErr != nil {
		return nil, errors.Wrapf(cErr, "failed to add gateway connectors to node's TOML config for don %s", don.Name)
	}

	// Set env vars from capability config to activate the confidential relay handler on DON nodes.
	// The handler is gated by CL_CONFIDENTIAL_RELAY_TRUSTED_PCRS; without it, the subservice won't start.
	capConfig, ok := don.CapabilityConfigs[flag]
	if ok && capConfig.Values != nil {
		ns := don.MustNodeSet()
		if ns.EnvVars == nil {
			ns.EnvVars = make(map[string]string)
		}

		if v, exists := capConfig.Values["trustedPCRs"]; exists {
			ns.EnvVars["CL_CONFIDENTIAL_RELAY_TRUSTED_PCRS"] = fmt.Sprintf("%v", v)
		}
		if v, exists := capConfig.Values["caRootsPEM"]; exists {
			ns.EnvVars["CL_CONFIDENTIAL_RELAY_CA_ROOTS_PEM"] = fmt.Sprintf("%v", v)
		}
	}

	// No on-chain capability registration needed. The relay handler is a CRE subservice,
	// not a registered capability. The mock capability that runs on the relay DON is
	// registered separately via the mock flag.
	return &cre.PreEnvStartupOutput{}, nil
}

func (o *ConfidentialRelay) PostEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.Don,
	dons *cre.Dons,
	creEnv *cre.Environment,
) error {
	return nil
}
