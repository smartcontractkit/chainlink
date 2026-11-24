package stellar

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	stellchain "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/stellar"
)

// Stellar is the Local CRE feature for the Stellar chain capability
type Stellar struct{}

var _ cre.Feature = (*Stellar)(nil)

const (
	flag              = cre.StellarCapability
	network           = "stellar"
	capabilityVersion = "1.0.0"
)

func (s *Stellar) Flag() cre.CapabilityFlag {
	return flag
}

func (s *Stellar) PreEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.DonMetadata,
	_ *cre.Topology,
	creEnv *cre.Environment,
) (*cre.PreEnvStartupOutput, error) {
	stellarChain := extractStellarFromEnv(creEnv)
	if stellarChain == nil {
		return &cre.PreEnvStartupOutput{}, nil
	}

	capabilityConfig, ok := don.MustNodeSet().GetCapabilityConfig(flag)
	if !ok {
		return nil, fmt.Errorf("config for '%s' capability not found for %s DON", flag, don.Name)
	}

	methodSettings, err := resolveMethodConfigSettings(capabilityConfig.Values)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve Stellar method config settings: %w", err)
	}

	// The forwarder address is part of the capability's spec_config, so deploy it first.
	if deployErr := deployStellarForwarders(ctx, testLogger, creEnv, stellarChain); deployErr != nil {
		return nil, fmt.Errorf("failed to deploy stellar CRE forwarder: %w", deployErr)
	}

	specConfig, err := buildSpecConfig(
		stellarChain.StellarChainID(),
		mustForwarderAddress(creEnv.CldfEnvironment.DataStore, stellarChain.ChainSelector()),
		methodSettings,
		don.HasFlag(cre.WorkflowDON),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build Stellar spec config: %w", err)
	}

	labelledName := CapabilityLabel(stellarChain.ChainSelector())
	capabilities := []keystone_changeset.DONCapabilityWithConfig{{
		Capability: kcr.CapabilitiesRegistryCapability{
			LabelledName:   labelledName,
			Version:        capabilityVersion,
			CapabilityType: 1,
		},
		Config: &capabilitiespb.CapabilityConfig{
			MethodConfigs: methodConfigs(methodSettings),
			LocalOnly:     don.HasOnlyLocalCapabilities(),
			SpecConfig:    specConfig,
		},
		UseCapRegOCRConfig: true,
	}}

	capabilityToExtraSignerFamilies := make(map[string][]string, len(capabilities))
	ocrConfigs := map[string]*ocr3.OracleConfig{}
	for _, capability := range capabilities {
		capabilityToExtraSignerFamilies[capability.Capability.LabelledName] = cre.OCRExtraSignerFamiliesForFamily(chainselectors.FamilyStellar)
		ocrConfigs[capability.Capability.LabelledName] = contracts.DefaultChainCapabilityOCR3Config()
	}

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfig:         capabilities,
		CapabilityToOCR3Config:          ocrConfigs,
		CapabilityToExtraSignerFamilies: capabilityToExtraSignerFamilies,
	}, nil
}

func (s *Stellar) PostEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.Don,
	dons *cre.Dons,
	creEnv *cre.Environment,
) error {
	stellarChain := extractStellarFromEnv(creEnv)
	if stellarChain == nil {
		return errors.New("no stellar blockchain found in environment")
	}

	consensusDons := dons.DonsWithFlag(cre.ConsensusCapability)
	if len(consensusDons) == 0 {
		return errors.New("no consensus DON found in environment")
	}

	if cfgErr := configureStellarForwarders(testLogger, creEnv, consensusDons[0], stellarChain); cfgErr != nil {
		return errors.Wrap(cfgErr, "failed to configure stellar CRE forwarder")
	}

	if authErr := authorizeStellarForwarderTransmitters(testLogger, creEnv, don, stellarChain); authErr != nil {
		return errors.Wrap(authErr, "failed to authorize stellar forwarder transmitters")
	}

	// The Stellar capability is launched by the node from its registry spec_config
	// (see PreEnvStartup), so no job spec is proposed here.
	return nil
}

func extractStellarFromEnv(creEnv *cre.Environment) *stellchain.Blockchain {
	for _, bcOut := range creEnv.Blockchains {
		if bcOut.IsFamily(chainselectors.FamilyStellar) {
			if sc, ok := bcOut.(*stellchain.Blockchain); ok {
				return sc
			}
		}
	}
	return nil
}
