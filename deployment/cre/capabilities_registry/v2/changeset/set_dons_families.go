package changeset

import (
	"errors"
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/pkg"
	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/sequences"
	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
)

var _ cldf.ChangeSetV2[SetDONsFamiliesInput] = SetDONsFamilies{}

type SetDONsFamiliesInput struct {
	RegistrySelector  uint64 `json:"registrySelector" yaml:"registrySelector"`
	RegistryQualifier string `json:"registryQualifier" yaml:"registryQualifier"`

	DONsFamiliesChanges []sequences.DONFamiliesChange `json:"donsFamiliesChanges" yaml:"donsFamiliesChanges"`

	MCMSConfig *ocr3.MCMSConfig `json:"mcmsConfig,omitempty" yaml:"mcmsConfig,omitempty"`
}

type SetDONsFamilies struct{}

func (l SetDONsFamilies) VerifyPreconditions(e cldf.Environment, config SetDONsFamiliesInput) error {
	if config.RegistrySelector <= 0 {
		return errors.New("RegistrySelector must be provided")
	}
	if config.RegistryQualifier == "" {
		return errors.New("RegistryQualifier must be provided")
	}
	if len(config.DONsFamiliesChanges) == 0 {
		return errors.New("must specify at least one DON family change")
	}
	return nil
}

func (l SetDONsFamilies) Apply(e cldf.Environment, config SetDONsFamiliesInput) (cldf.ChangesetOutput, error) {
	var mcmsContracts *commonchangeset.MCMSWithTimelockState
	if config.MCMSConfig != nil {
		var err error
		mcmsContracts, err = strategies.GetMCMSContracts(e, config.RegistrySelector, emptyQualifier)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get MCMS contracts: %w", err)
		}
	}

	registryRef := pkg.GetCapRegV2AddressRefKey(config.RegistrySelector, config.RegistryQualifier)

	report, err := operations.ExecuteSequence(
		e.OperationsBundle,
		sequences.SetDONsFamilies,
		sequences.SetDONsFamiliesDeps{
			Env:           &e,
			MCMSContracts: mcmsContracts,
		},
		sequences.SetDONsFamiliesInput{
			RegistryRef: registryRef,
			DONsChanges: config.DONsFamiliesChanges,
			MCMSConfig:  config.MCMSConfig,
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute SetDONsFamilies sequence: %w", err)
	}

	return cldf.ChangesetOutput{
		Reports:               report.ExecutionReports,
		MCMSTimelockProposals: report.Output.Proposals,
	}, nil
}
