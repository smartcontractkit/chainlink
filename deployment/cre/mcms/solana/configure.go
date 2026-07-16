package solana

import (
	"errors"

	mcmschangesets "github.com/smartcontractkit/cld-changesets/legacy/mcms/changesets"

	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/cre/contracts"
)

type ConfigureSolanaMCMS struct{}

var _ cldf.ChangeSetV2[ConfigureSolanaMCMSConfig] = ConfigureSolanaMCMS{}

// ConfigureSolanaMCMSConfig is the input for updating Solana MCMS signer config on-chain.
// Mirrors cre/mcms/changeset.ConfigureChangesetInput for EVM mcms_set_config.
type ConfigureSolanaMCMSConfig struct {
	Environment string `json:"environment,omitempty" yaml:"environment,omitempty"`

	// ChainSelectors are the Solana chain selectors where MCMS contracts will be updated.
	ChainSelectors []uint64 `json:"chainSelectors" yaml:"chainSelectors"`

	// MCMSWithTimelockConfig is the signer policy to apply via SetConfigMCMSV2.
	MCMSWithTimelockConfig cldfproposalutils.MCMSWithTimelockConfig `json:"mcmsWithTimelockConfig" yaml:"mcmsWithTimelockConfig"`

	// MCMSConfig wraps the update in MCMS timelock proposals when MCMSAction is set.
	MCMSConfig contracts.MCMSConfig `json:"mcmsConfig" yaml:"mcmsConfig"`
}

func (ConfigureSolanaMCMS) VerifyPreconditions(env cldf.Environment, cfg ConfigureSolanaMCMSConfig) error {
	if len(cfg.ChainSelectors) == 0 {
		return errors.New("no chain selectors provided")
	}
	for _, sel := range cfg.ChainSelectors {
		if err := verifySelector(env, sel); err != nil {
			return err
		}
	}

	return nil
}

func (ConfigureSolanaMCMS) Apply(env cldf.Environment, cfg ConfigureSolanaMCMSConfig) (cldf.ChangesetOutput, error) {
	mcmsConfigPerChain := make(map[uint64]mcmschangesets.ConfigPerRoleV2, len(cfg.ChainSelectors))
	for _, sel := range cfg.ChainSelectors {
		mcmsConfigPerChain[sel] = mcmschangesets.ConfigPerRoleV2{
			Proposer:  &cfg.MCMSWithTimelockConfig.Proposer,
			Canceller: &cfg.MCMSWithTimelockConfig.Canceller,
			Bypasser:  &cfg.MCMSWithTimelockConfig.Bypasser,
		}
	}

	setCfg := mcmschangesets.MCMSConfigV2{
		ConfigsPerChain: mcmsConfigPerChain,
	}
	if cfg.MCMSConfig.MCMSAction != "" {
		setCfg.ProposalConfig = &cfg.MCMSConfig
	}

	return mcmschangesets.SetConfigMCMSV2(env, setCfg)
}
