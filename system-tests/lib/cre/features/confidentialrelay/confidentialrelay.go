package confidentialrelay

import (
	"context"

	tomlser "github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	coretoml "github.com/smartcontractkit/chainlink/v2/core/config/toml"
	corechainlink "github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

const flag = cre.ConfidentialRelayCapability

type ConfidentialRelay struct {
	// TrustEnclaves makes the relay trust fake (non-Nitro) enclaves by
	// relaxing TEE attestation validation. INSECURE; test/E2E use only.
	TrustEnclaves bool
	// RequireBFTQuorum determines the required signature quorum for the relay.
	RequireBFTQuorum bool
}

func (o *ConfidentialRelay) Flag() cre.CapabilityFlag {
	return flag
}

// boolFromValues reads a bool from a capability config's values, falling back to
// def when the key is absent or not a bool. Lets a topology set the relay's knobs
// in TOML rather than requiring a Go-constructed feature.
func boolFromValues(values map[string]any, key string, def bool) bool {
	if v, ok := values[key]; ok {
		if b, isBool := v.(bool); isBool {
			return b
		}
	}

	return def
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

	// Gateway connector injection scoped to this workflow DON's don_family (see topology_don_family.go).
	cErr := don.ConfigureForGatewayAccess(registryChainID, topology.GatewayConnectorsForDonFamily(don.DonFamily))
	if cErr != nil {
		return nil, errors.Wrapf(cErr, "failed to add gateway connectors to node's TOML config for don %s", don.Name)
	}

	// Set TOML config to activate the confidential relay handler on DON nodes.
	capConfig, ok := don.CapabilityConfigs[flag]
	if ok && capConfig.Values != nil {
		ns := don.MustNodeSet()
		for i := range ns.NodeSpecs {
			currentConfig := ns.NodeSpecs[i].Node.TestConfigOverrides
			var typedConfig corechainlink.Config
			if currentConfig != "" {
				if err := tomlser.Unmarshal([]byte(currentConfig), &typedConfig); err != nil {
					return nil, errors.Wrapf(err, "failed to unmarshal node TOML config for node %d", i)
				}
			}

			enabled := true
			trustEnclaves := boolFromValues(capConfig.Values, "trustEnclaves", o.TrustEnclaves)
			requireBFTQuorum := boolFromValues(capConfig.Values, "requireBFTQuorum", o.RequireBFTQuorum)
			typedConfig.CRE.ConfidentialRelay = &coretoml.ConfidentialRelayConfig{
				Enabled:          &enabled,
				TrustEnclaves:    &trustEnclaves,
				RequireBFTQuorum: &requireBFTQuorum,
			}

			out, err := tomlser.Marshal(typedConfig)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to marshal node TOML config for node %d", i)
			}
			ns.NodeSpecs[i].Node.TestConfigOverrides = string(out)
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
