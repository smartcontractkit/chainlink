package changeset

import (
	"errors"
	"fmt"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/pkg"
	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/sequences"
)

var _ cldf.ChangeSetV2[AddCapabilitiesInput] = AddCapabilities{}

// emptyQualifier is used when no specific qualifier is needed
const emptyQualifier = ""

type AddCapabilitiesInput struct {
	RegistryChainSel  uint64 `json:"registryChainSel" yaml:"registryChainSel"`
	RegistryQualifier string `json:"registryQualifier" yaml:"registryQualifier"`

	MCMSConfig        *ocr3.MCMSConfig             `json:"mcmsConfig" yaml:"mcmsConfig"`
	DonName           string                       `json:"donName" yaml:"donName"`
	CapabilityConfigs []contracts.CapabilityConfig `json:"capabilityConfigs" yaml:"capabilityConfigs"`

	// Force indicates whether to force the update even if we cannot validate that all forwarder contracts are ready to accept the new configure version.
	// This is very dangerous, and could break the whole platform if the forwarders are not ready. Be very careful with this option.
	Force bool `json:"force" yaml:"force"`
}

type AddCapabilities struct{}

func (u AddCapabilities) VerifyPreconditions(_ cldf.Environment, config AddCapabilitiesInput) error {
	if config.DonName == "" {
		return errors.New("must specify DONName")
	}
	if len(config.CapabilityConfigs) == 0 {
		return errors.New("capabilityConfigs is required")
	}
	return nil
}

func (u AddCapabilities) Apply(e cldf.Environment, config AddCapabilitiesInput) (cldf.ChangesetOutput, error) {
	var mcmsContracts *commonchangeset.MCMSWithTimelockState
	if config.MCMSConfig != nil {
		var err error
		mcmsContracts, err = strategies.GetMCMSContracts(e, config.RegistryChainSel, emptyQualifier)
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
			RegistryRef:       registryRef,
			DonName:           config.DonName,
			CapabilityConfigs: config.CapabilityConfigs,
			Force:             config.Force,
			MCMSConfig:        config.MCMSConfig,
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
