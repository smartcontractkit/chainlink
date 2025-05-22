package v0_5_0

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/mcmsutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/txutil"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/configurator"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
)

var SetStagingConfigChangeset = cldf.CreateChangeSet(setStagingConfigLogic, setStagingConfigPrecondition)

type Config interface {
	GetConfiguratorAddress() common.Address
}

type SetStagingConfigConfig struct {
	ConfigurationsByChain map[uint64][]ConfiguratorConfig
	MCMSConfig            *types.MCMSConfig
}

func setStagingConfigPrecondition(_ cldf.Environment, ss SetStagingConfigConfig) error {
	if err := ss.Validate(); err != nil {
		return fmt.Errorf("invalid SetStagingConfigConfig: %w", err)
	}

	return nil
}

func (cfg SetStagingConfigConfig) Validate() error {
	if len(cfg.ConfigurationsByChain) == 0 {
		return errors.New("ConfigurationsByChain cannot be empty")
	}
	return nil
}

func setStagingConfigLogic(e cldf.Environment, cfg SetStagingConfigConfig) (cldf.ChangesetOutput, error) {
	txs, err := txutil.GetTxs(
		e,
		types.Configurator.String(),
		cfg.ConfigurationsByChain,
		LoadConfigurator,
		doSetStagingConfig,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed building SetStagingConfig txs: %w", err)
	}

	return mcmsutil.ExecuteOrPropose(e, txs, cfg.MCMSConfig, "SetStagingConfig proposal")
}

func doSetStagingConfig(
	c *configurator.Configurator,
	cfg ConfiguratorConfig,
) (*ethTypes.Transaction, error) {
	return c.SetStagingConfig(cldf.SimTransactOpts(),
		cfg.ConfigID,
		cfg.Signers,
		cfg.OffchainTransmitters,
		cfg.F,
		cfg.OnchainConfig,
		cfg.OffchainConfigVersion,
		cfg.OffchainConfig)
}
