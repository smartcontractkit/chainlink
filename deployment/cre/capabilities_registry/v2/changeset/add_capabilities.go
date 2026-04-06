package changeset

import (
	"errors"
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/pkg"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/sequences"
	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
	crecontracts "github.com/smartcontractkit/chainlink/deployment/cre/contracts"
)

var _ cldf.ChangeSetV2[AddCapabilitiesInput] = AddCapabilities{}

type AddCapabilitiesInput struct {
	RegistryChainSel  uint64 `json:"registryChainSel" yaml:"registryChainSel"`
	RegistryQualifier string `json:"registryQualifier" yaml:"registryQualifier"`

	MCMSConfig *crecontracts.MCMSConfig `json:"mcmsConfig" yaml:"mcmsConfig"`

	// DonCapabilityConfigs maps DON name to the list of capability configs for that DON.
	DonCapabilityConfigs map[string][]contracts.CapabilityConfig `json:"donCapabilityConfigs" yaml:"donCapabilityConfigs"`

	// Force indicates whether to force the update even if we cannot validate that all forwarder contracts are ready to accept the new configure version.
	// This is very dangerous, and could break the whole platform if the forwarders are not ready. Be very careful with this option.
	Force bool `json:"force" yaml:"force"`
}

type AddCapabilities struct{}

func (u AddCapabilities) VerifyPreconditions(_ cldf.Environment, config AddCapabilitiesInput) error {
	if len(config.DonCapabilityConfigs) == 0 {
		return errors.New("donCapabilityConfigs must contain at least one DON entry")
	}
	for donName, configs := range config.DonCapabilityConfigs {
		if donName == "" {
			return errors.New("donCapabilityConfigs keys cannot be empty strings")
		}
		if len(configs) == 0 {
			return fmt.Errorf("donCapabilityConfigs[%q] must contain at least one capability config", donName)
		}
	}
	return nil
}

func (u AddCapabilities) Apply(e cldf.Environment, config AddCapabilitiesInput) (cldf.ChangesetOutput, error) {
	var mcmsContracts *commonchangeset.MCMSWithTimelockState
	if config.MCMSConfig != nil {
		var err error
		mcmsContracts, err = strategies.GetMCMSContracts(e, config.RegistryChainSel, *config.MCMSConfig)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get MCMS contracts: %w", err)
		}
	}

	registryRef := pkg.GetCapRegV2AddressRefKey(config.RegistryChainSel, config.RegistryQualifier)

	seqReport, err := operations.ExecuteSequence(
		e.OperationsBundle,
		sequences.AddCapabilities,
		sequences.AddCapabilitiesDeps{Env: &e, MCMSContracts: mcmsContracts},
		sequences.AddCapabilitiesInput{
			RegistryRef:          registryRef,
			DonCapabilityConfigs: config.DonCapabilityConfigs,
			Force:                config.Force,
			MCMSConfig:           config.MCMSConfig,
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	return cldf.ChangesetOutput{
		Reports:               seqReport.ExecutionReports,
		MCMSTimelockProposals: seqReport.Output.Proposals,
	}, nil
}
