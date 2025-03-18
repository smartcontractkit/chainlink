package v0_5_0

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	goEthTypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/mcmsutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/txutil"
	rewardManager "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/reward_manager_v0_5_0"
)

var SetFeeManagerChangeset = deployment.CreateChangeSet(SetFeeManagerLogic, SetFeeManagerPrecondition)

type SetFeeManagerConfig struct {
	ConfigsByChain map[uint64][]SetFeeManager
	MCMSConfig     *changeset.MCMSConfig
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

func SetFeeManagerPrecondition(_ deployment.Environment, cc SetFeeManagerConfig) error {
	if err := cc.Validate(); err != nil {
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
