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

var PayRecipientsChangeset = cldf.CreateChangeSet(PayRecipientsLogic, PayRecipientsPrecondition)

type PayRecipientsConfig struct {
	ConfigsByChain map[uint64][]PayRecipients
	MCMSConfig     *types.MCMSConfig
}

type PayRecipients struct {
	RewardManagerAddress common.Address

	PoolID     [32]byte
	Recipients []common.Address
}

func (a PayRecipients) GetContractAddress() common.Address {
	return a.RewardManagerAddress
}

func (cfg PayRecipientsConfig) Validate() error {
	if len(cfg.ConfigsByChain) == 0 {
		return errors.New("ConfigsByChain cannot be empty")
	}
	return nil
}

func PayRecipientsPrecondition(_ deployment.Environment, cc PayRecipientsConfig) error {
	if err := cc.Validate(); err != nil {
		return fmt.Errorf("invalid PayRecipients config: %w", err)
	}
	return nil
}

func PayRecipientsLogic(e deployment.Environment, cfg PayRecipientsConfig) (deployment.ChangesetOutput, error) {
	txs, err := txutil.GetTxs(
		e,
		types.RewardManager.String(),
		cfg.ConfigsByChain,
		loadRewardManagerState,
		doPayRecipients,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed building PayRecipients txs: %w", err)
	}

	return mcmsutil.ExecuteOrPropose(e, txs, cfg.MCMSConfig, "PayRecipients proposal")
}

func doPayRecipients(vs *rewardManager.RewardManager, pr PayRecipients) (*goEthTypes.Transaction, error) {
	return vs.PayRecipients(
		deployment.SimTransactOpts(),
		pr.PoolID,
		pr.Recipients,
	)
}
