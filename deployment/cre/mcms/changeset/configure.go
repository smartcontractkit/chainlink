package changeset

import (
	"errors"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/cre/contracts"
)

type CsMCMSConfigure struct{}

var _ cldf.ChangeSetV2[ConfigureChangesetInput] = CsMCMSConfigure{}

type ContractConfiguration struct {
	Config types.MCMSWithTimelockConfigV2 `json:"config,omitempty" yaml:"config,omitempty"`
}

// MCMSConfigureChangesetInput is the input for the set MCMS configuration changeset.
// The contract configuration applies this same configuration to all chain selectors.
type ConfigureChangesetInput struct {
	Environment string
	// MCMSContractConfiguration is the configuration for the MCMS contract to be updated to
	MCMSContractConfiguration ContractConfiguration `json:"contractConfigurations" yaml:"contractConfigurations"`
	// ChainSelectors are the chain selectors where the MCMS contracts will be updated
	ChainSelectors []uint64             `json:"chainSelectors" yaml:"chainSelectors"`
	MCMSConfig     contracts.MCMSConfig `json:"mcmsConfig" yaml:"mcmsConfig"`
}

// VerifyPreconditions checks if the input is valid before applying the changeset.
func (CsMCMSConfigure) VerifyPreconditions(env cldf.Environment, input ConfigureChangesetInput) error {
	if len(input.ChainSelectors) == 0 {
		return errors.New("no chain selectors provided")
	}

	return nil
}

func (CsMCMSConfigure) Apply(env cldf.Environment, input ConfigureChangesetInput) (cldf.ChangesetOutput, error) {
	mcmsConfigPerChain := make(map[uint64]changeset.ConfigPerRoleV2, len(input.ChainSelectors))
	for _, s := range input.ChainSelectors {
		c := input.MCMSContractConfiguration.Config
		mcmsConfigPerChain[s] = changeset.ConfigPerRoleV2{
			Proposer:  &c.Proposer,
			Canceller: &c.Canceller,
			Bypasser:  &c.Bypasser,
		}
	}

	cfg := changeset.MCMSConfigV2{
		ConfigsPerChain: mcmsConfigPerChain,
	}
	// MCMSConfig.validateCommon enforces that MCMSAction is set. If there isn't, presume empty.
	if input.MCMSConfig.MCMSAction != "" {
		cfg.ProposalConfig = &input.MCMSConfig
	}

	o, err := changeset.SetConfigMCMSV2(env, cfg)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	// no DataStore or reports are returned from SetConfigMCMSV2
	return o, nil
}
