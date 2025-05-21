package v0_5_0

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/mcmsutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/txutil"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/configurator"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
)

var PromoteStagingConfigChangeset = cldf.CreateChangeSet(promoteStagingConfigLogic, promoteStagingConfigPrecondition)

type PromoteStagingConfigConfig struct {
	PromotionsByChain map[uint64][]PromoteStagingConfig
	MCMSConfig        *types.MCMSConfig
}

type PromoteStagingConfig struct {
	ConfiguratorAddress common.Address
	// 32-byte configId
	ConfigID [32]byte
	// Whether the current production config is considered green
	IsGreenProduction bool
}

type State struct {
	Configurator *configurator.Configurator
}

func (cfg PromoteStagingConfigConfig) Validate() error {
	if len(cfg.PromotionsByChain) == 0 {
		return errors.New("PromotionsByChain cannot be empty")
	}
	return nil
}

func (pc PromoteStagingConfig) GetContractAddress() common.Address { return pc.ConfiguratorAddress }

func promoteStagingConfigPrecondition(_ cldf.Environment, cc PromoteStagingConfigConfig) error {
	if err := cc.Validate(); err != nil {
		return fmt.Errorf("invalid DeployConfiguratorConfig: %w", err)
	}
	return nil
}

func promoteStagingConfigLogic(e cldf.Environment, cfg PromoteStagingConfigConfig) (cldf.ChangesetOutput, error) {
	txs, err := txutil.GetTxs(
		e,
		types.Configurator.String(),
		cfg.PromotionsByChain,
		LoadConfigurator,
		doPromoteStagingConfig,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed building SetProductionConfig txs: %w", err)
	}

	return mcmsutil.ExecuteOrPropose(e, txs, cfg.MCMSConfig, "SetProductionConfig proposal")
}

func doPromoteStagingConfig(
	c *configurator.Configurator,
	cfg PromoteStagingConfig,
) (*ethTypes.Transaction, error) {
	return c.PromoteStagingConfig(cldf.SimTransactOpts(), cfg.ConfigID, cfg.IsGreenProduction)
}

func LoadConfigurator(e cldf.Environment, chainSel uint64, contractAddr string) (*configurator.Configurator, error) {
	chain, ok := e.Chains[chainSel]
	if !ok {
		return nil, fmt.Errorf("chain %d not found", chainSel)
	}

	if err := utils.ValidateContract(e, chainSel, contractAddr, types.Configurator, deployment.Version0_5_0); err != nil {
		return nil, fmt.Errorf("invalid contract address %s on chain %d", contractAddr, chainSel)
	}

	conf, err := configurator.NewConfigurator(common.HexToAddress(contractAddr), chain.Client)
	if err != nil {
		return nil, fmt.Errorf("failed to load Configurator contract on chain %s (chain selector %d): %w", chain.Name(), chain.Selector, err)
	}

	return conf, nil
}
