package changeset

import (
	"errors"

	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/pkg"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/sequences"
)

var _ cldf.ChangeSetV2[AddCapabilitiesInput] = AddCapabilities{}

type AddCapabilitiesInput struct {
	RegistryChainSel  uint64 `json:"registry_chain_sel" yaml:"registry_chain_sel"`
	RegistryQualifier string `json:"registry_qualifier" yaml:"registry_qualifier"`
	UseMCMS           bool   `json:"use_mcms" yaml:"use_mcms"` // not implemented yet

	DonName           string                       `json:"don_name" yaml:"don_name"`
	CapabilityConfigs []contracts.CapabilityConfig `json:"capability_configs" yaml:"capability_configs"`

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
	registryRef := pkg.GetCapRegV2AddressRefKey(config.RegistryChainSel, config.RegistryQualifier)

	seqReport, err := operations.ExecuteSequence(
		e.OperationsBundle,
		sequences.AddCapabilities,
		sequences.AddCapabilitiesDeps{Env: &e},
		sequences.AddCapabilitiesInput{
			RegistryRef:       registryRef,
			DonName:           config.DonName,
			CapabilityConfigs: config.CapabilityConfigs,
			Force:             config.Force,
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	return cldf.ChangesetOutput{
		Reports: seqReport.ExecutionReports,
	}, nil
}
