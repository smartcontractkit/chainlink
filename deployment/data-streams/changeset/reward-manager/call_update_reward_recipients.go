package reward_manager

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	goEthTypes "github.com/ethereum/go-ethereum/core/types"

	rewardManager "github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/reward_manager_v0_5_0"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/mcmsutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/txutil"
)

var UpdateRewardRecipientsChangeset = cldf.CreateChangeSet(updateRewardRecipientsLogic, updateRewardRecipientsPrecondition)

type UpdateRewardRecipientsConfig struct {
	ConfigsByChain map[uint64][]UpdateRewardRecipients
	MCMSConfig     *types.MCMSConfig
}

type UpdateRewardRecipients struct {
	RewardManagerAddress common.Address

	PoolID                    [32]byte
	RewardRecipientAndWeights []rewardManager.CommonAddressAndWeight
}

func (a UpdateRewardRecipients) GetContractAddress() common.Address {
	return a.RewardManagerAddress
}

func (cfg UpdateRewardRecipientsConfig) Validate() error {
	if len(cfg.ConfigsByChain) == 0 {
		return errors.New("ConfigsByChain cannot be empty")
	}
	return nil
}

func updateRewardRecipientsPrecondition(_ deployment.Environment, cc UpdateRewardRecipientsConfig) error {
	if err := cc.Validate(); err != nil {
		return fmt.Errorf("invalid UpdateRewardRecipients config: %w", err)
	}
	return nil
}

func updateRewardRecipientsLogic(e deployment.Environment, cfg UpdateRewardRecipientsConfig) (deployment.ChangesetOutput, error) {
	txs, err := txutil.GetTxs(
		e,
		types.RewardManager.String(),
		cfg.ConfigsByChain,
		loadRewardManagerState,
		doUpdateRewardRecipients,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed building UpdateRewardRecipients txs: %w", err)
	}

	return mcmsutil.ExecuteOrPropose(e, txs, cfg.MCMSConfig, "UpdateRewardRecipients proposal")
}

func doUpdateRewardRecipients(vs *rewardManager.RewardManager, ur UpdateRewardRecipients) (*goEthTypes.Transaction, error) {
	return vs.UpdateRewardRecipients(
		deployment.SimTransactOpts(),
		ur.PoolID,
		ur.RewardRecipientAndWeights,
	)
}
