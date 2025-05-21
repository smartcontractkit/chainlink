package v0_5_0

import (
	"errors"
	"fmt"

	ethTypes "github.com/ethereum/go-ethereum/core/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/mcmsutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/txutil"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/configurator"
)

var SetProductionConfigChangeset = cldf.CreateChangeSet(setProductionConfigLogic, setProductionConfigPrecondition)

// SetProductionConfigConfig contains the parameters needed to set a production config.
type SetProductionConfigConfig struct {
	ConfigurationsByChain map[uint64][]ConfiguratorConfig
	MCMSConfig            *types.MCMSConfig
}

func (cfg SetProductionConfigConfig) Validate() error {
	if len(cfg.ConfigurationsByChain) == 0 {
		return errors.New("ConfigurationsByChain cannot be empty")
	}
	return nil
}

func setProductionConfigPrecondition(_ cldf.Environment, cfg SetProductionConfigConfig) error {
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid DeployConfiguratorConfig: %w", err)
	}

	return nil
}

func setProductionConfigLogic(e cldf.Environment, cfg SetProductionConfigConfig) (cldf.ChangesetOutput, error) {
	txs, err := txutil.GetTxs(
		e,
		types.Configurator.String(),
		cfg.ConfigurationsByChain,
		LoadConfigurator,
		doSetProductionConfig,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed building SetProductionConfig txs: %w", err)
	}

	return mcmsutil.ExecuteOrPropose(e, txs, cfg.MCMSConfig, "SetProductionConfig proposal")
}

func doSetProductionConfig(
	c *configurator.Configurator,
	prodCfg ConfiguratorConfig,
) (*ethTypes.Transaction, error) {
	return c.SetProductionConfig(cldf.SimTransactOpts(),
		prodCfg.ConfigID,
		prodCfg.Signers,
		prodCfg.OffchainTransmitters,
		prodCfg.F,
		prodCfg.OnchainConfig,
		prodCfg.OffchainConfigVersion,
		prodCfg.OffchainConfig)
}
