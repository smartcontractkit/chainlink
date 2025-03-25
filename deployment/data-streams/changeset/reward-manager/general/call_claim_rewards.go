package general

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

var ClaimRewardsChangeset = deployment.CreateChangeSet(claimRewardsLogic, claimRewardsPrecondition)

type ClaimRewardsConfig struct {
	ConfigsByChain map[uint64][]ClaimRewards
	MCMSConfig     *changeset.MCMSConfig
}

type ClaimRewards struct {
	RewardManagerAddress common.Address

	PoolIDs [][32]byte
}

func (a ClaimRewards) GetContractAddress() common.Address {
	return a.RewardManagerAddress
}

func (cfg ClaimRewardsConfig) Validate() error {
	if len(cfg.ConfigsByChain) == 0 {
		return errors.New("ConfigsByChain cannot be empty")
	}
	return nil
}

func claimRewardsPrecondition(_ deployment.Environment, cc ClaimRewardsConfig) error {
	if err := cc.Validate(); err != nil {
		return fmt.Errorf("invalid ClaimRewards config: %w", err)
	}
	return nil
}

func claimRewardsLogic(e deployment.Environment, cfg ClaimRewardsConfig) (deployment.ChangesetOutput, error) {
	txs, err := txutil.GetTxs(
		e,
		types.RewardManager.String(),
		cfg.ConfigsByChain,
		loadRewardManagerState,
		doClaimRewards,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed building ClaimRewards txs: %w", err)
	}

	return mcmsutil.ExecuteOrPropose(e, txs, cfg.MCMSConfig, "ClaimRewards proposal")
}

func doClaimRewards(vs *rewardManager.RewardManager, cr ClaimRewards) (*goEthTypes.Transaction, error) {
	return vs.ClaimRewards(
		deployment.SimTransactOpts(),
		cr.PoolIDs,
	)
}
