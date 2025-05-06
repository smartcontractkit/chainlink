package v0_5_0

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/configurator"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset"
)

var SetProductionConfigChangeset = cldf.CreateChangeSet(setProductionConfigLogic, setProductionConfigPrecondition)

// SetProductionConfigConfig contains the parameters needed to set a production config.
type SetProductionConfigConfig struct {
	ConfigurationsByChain map[uint64][]SetProductionConfig
	MCMSConfig            *changeset.MCMSConfig
}

// SetProductionConfig groups the parameters for a production config call.
type SetProductionConfig struct {
	ConfiguratorAddress   common.Address
	ConfigID              [32]byte
	Signers               [][]byte
	OffchainTransmitters  [][32]byte
	F                     uint8
	OnchainConfig         []byte
	OffchainConfigVersion uint64
	OffchainConfig        []byte
}

func (pc SetProductionConfig) GetConfiguratorAddress() common.Address {
	return pc.ConfiguratorAddress
}

func (cfg SetProductionConfigConfig) Validate() error {
	if len(cfg.ConfigurationsByChain) == 0 {
		return errors.New("ConfigurationsByChain cannot be empty")
	}
	return nil
}

func setProductionConfigPrecondition(_ deployment.Environment, cfg SetProductionConfigConfig) error {
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid DeployConfiguratorConfig: %w", err)
	}

	return nil
}

func getTransactOpts(e deployment.Environment, chainSel uint64, mcmsConfig *changeset.MCMSConfig) *bind.TransactOpts {
	if mcmsConfig == nil {
		return e.Chains[chainSel].DeployerKey
	}
	return deployment.SimTransactOpts()
}

func setProductionConfigTx(
	e deployment.Environment,
	configuratorContract *configurator.Configurator,
	prodCfg SetProductionConfig,
	opts *bind.TransactOpts,
	chain deployment.Chain,
	mcmsConfig *changeset.MCMSConfig,
) (*ethTypes.Transaction, error) {
	tx, err := configuratorContract.SetProductionConfig(opts, prodCfg.ConfigID, prodCfg.Signers, prodCfg.OffchainTransmitters, prodCfg.F, prodCfg.OnchainConfig, prodCfg.OffchainConfigVersion, prodCfg.OffchainConfig)
	if err != nil {
		return nil, fmt.Errorf("error packing setProductionConfig tx data: %w", err)
	}

	if mcmsConfig == nil {
		if _, err := deployment.ConfirmIfNoError(chain, tx, nil); err != nil {
			e.Logger.Errorw("Failed to confirm setProductionConfig tx", "chain", chain.String(), "err", err)
			return nil, err
		}
	}
	return tx, nil
}

func setProductionConfigLogic(e deployment.Environment, cfg SetProductionConfigConfig) (deployment.ChangesetOutput, error) {
	return callSetConfigCommon(e, cfg.ConfigurationsByChain, cfg.MCMSConfig, "SetProductionConfig proposal", setProductionConfigTx)
}
