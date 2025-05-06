package reward_manager

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	goEthTypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/mcmsutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/txutil"

	rewardManager "github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/reward_manager_v0_5_0"
)

var SetFeeManagerChangeset = cldf.CreateChangeSet(SetFeeManagerLogic, SetFeeManagerPrecondition)

type SetFeeManagerConfig struct {
	ConfigsByChain map[uint64][]SetFeeManager
	MCMSConfig     *types.MCMSConfig
}

type SetFeeManager struct {
	FeeManagerAddress    common.Address
	RewardManagerAddress common.Address
}

func (a SetFeeManager) GetContractAddress() common.Address {
	return a.RewardManagerAddress
}

func (cfg SetFeeManagerConfig) Validate() error {
	if len(cfg.ConfigsByChain) == 0 {
		return errors.New("ConfigsByChain cannot be empty")
	}
	return nil
}

func SetFeeManagerPrecondition(e deployment.Environment, sf SetFeeManagerConfig) error {
	for chainSelector, configs := range sf.ConfigsByChain {
		for _, config := range configs {
			confState, err := loadRewardManagerState(e, chainSelector, config.RewardManagerAddress.Hex())
			if err != nil {
				return err
			}

			gotVersion, err := confState.TypeAndVersion(&bind.CallOpts{Context: e.GetContext()})
			if err != nil {
				return fmt.Errorf("failed to get RewardManager version: %w", err)
			}

			// Why v0.5.0/RewardManager.sol typeAndVersion returns 1.1.0?
			allowedVersion := deployment.NewTypeAndVersion(types.RewardManager, deployment.Version1_1_0).String()

			if gotVersion != allowedVersion {
				return fmt.Errorf("invalid RewardManager version: got %s, allowed %s", gotVersion, allowedVersion)
			}
		}
	}

	if err := sf.Validate(); err != nil {
		return fmt.Errorf("invalid SetFeeManager config: %w", err)
	}
	return nil
}

func SetFeeManagerLogic(e deployment.Environment, cfg SetFeeManagerConfig) (deployment.ChangesetOutput, error) {
	txs, err := txutil.GetTxs(
		e,
		types.FeeManager.String(),
		cfg.ConfigsByChain,
		loadRewardManagerState,
		doSetFeeManager,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed building SetFeeManager txs: %w", err)
	}

	return mcmsutil.ExecuteOrPropose(e, txs, cfg.MCMSConfig, "SetFeeManager proposal")
}

func doSetFeeManager(vs *rewardManager.RewardManager, sf SetFeeManager) (*goEthTypes.Transaction, error) {
	return vs.SetFeeManager(
		deployment.SimTransactOpts(),
		sf.FeeManagerAddress,
	)
}
