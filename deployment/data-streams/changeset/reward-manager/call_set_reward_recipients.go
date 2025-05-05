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

var SetRewardRecipientsChangeset = cldf.CreateChangeSet(setRewardRecipientsLogic, setRewardRecipientsPrecondition)

type SetRewardRecipientsConfig struct {
	ConfigsByChain map[uint64][]SetRewardRecipients
	MCMSConfig     *types.MCMSConfig
}

type SetRewardRecipients struct {
	RewardManagerAddress common.Address

	PoolID                    [32]byte
	RewardRecipientAndWeights []rewardManager.CommonAddressAndWeight
}

func (a SetRewardRecipients) GetContractAddress() common.Address {
	return a.RewardManagerAddress
}

func (cfg SetRewardRecipientsConfig) Validate() error {
	if len(cfg.ConfigsByChain) == 0 {
		return errors.New("ConfigsByChain cannot be empty")
	}
	return nil
}

func setRewardRecipientsPrecondition(_ deployment.Environment, cc SetRewardRecipientsConfig) error {
	if err := cc.Validate(); err != nil {
		return fmt.Errorf("invalid SetRewardRecipients config: %w", err)
	}
	return nil
}

func setRewardRecipientsLogic(e deployment.Environment, cfg SetRewardRecipientsConfig) (deployment.ChangesetOutput, error) {
	txs, err := txutil.GetTxs(
		e,
		types.RewardManager.String(),
		cfg.ConfigsByChain,
		loadRewardManagerState,
		doSetRewardRecipients,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed building SetRewardRecipients txs: %w", err)
	}

	return mcmsutil.ExecuteOrPropose(e, txs, cfg.MCMSConfig, "SetRewardRecipients proposal")
}

func doSetRewardRecipients(vs *rewardManager.RewardManager, sr SetRewardRecipients) (*goEthTypes.Transaction, error) {
	return vs.SetRewardRecipients(
		deployment.SimTransactOpts(),
		sr.PoolID,
		sr.RewardRecipientAndWeights,
	)
}
