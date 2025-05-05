package reward_manager

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/mcmsutil"

	rewardManager "github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/reward_manager_v0_5_0"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
)

var DeployRewardManagerChangeset = cldf.CreateChangeSet(deployRewardManagerLogic, deployRewardManagerPrecondition)

type DeployRewardManager struct {
	LinkTokenAddress common.Address
}

type DeployRewardManagerConfig struct {
	ChainsToDeploy map[uint64]DeployRewardManager
	Ownership      types.OwnershipSettings
}

func (cc DeployRewardManagerConfig) Validate() error {
	if len(cc.ChainsToDeploy) == 0 {
		return errors.New("ChainsToDeploy is empty")
	}
	for chain := range cc.ChainsToDeploy {
		if err := deployment.IsValidChainSelector(chain); err != nil {
			return fmt.Errorf("invalid chain selector: %d - %w", chain, err)
		}
	}
	return nil
}

func deployRewardManagerLogic(e deployment.Environment, cc DeployRewardManagerConfig) (deployment.ChangesetOutput, error) {
	ab := deployment.NewMemoryAddressBook()
	err := deployRewardManager(e, ab, cc)
	if err != nil {
		e.Logger.Errorw("Failed to deploy RewardManager", "err", err, "addresses", ab)
		return deployment.ChangesetOutput{AddressBook: ab}, deployment.MaybeDataErr(err)
	}

	if cc.Ownership.ShouldTransfer && cc.Ownership.MCMSProposalConfig != nil {
		filter := deployment.NewTypeAndVersion(types.RewardManager, deployment.Version0_5_0)
		return mcmsutil.TransferToMCMSWithTimelockForTypeAndVersion(e, ab, filter, *cc.Ownership.MCMSProposalConfig)
	}

	return deployment.ChangesetOutput{
		AddressBook: ab,
	}, nil
}

func deployRewardManagerPrecondition(_ deployment.Environment, cc DeployRewardManagerConfig) error {
	if err := cc.Validate(); err != nil {
		return fmt.Errorf("invalid DeployRewardManagerConfig: %w", err)
	}

	return nil
}

func deployRewardManager(e deployment.Environment, ab deployment.AddressBook, cc DeployRewardManagerConfig) error {
	if err := cc.Validate(); err != nil {
		return fmt.Errorf("invalid DeployRewardManagerConfig: %w", err)
	}

	for chainSel := range cc.ChainsToDeploy {
		chain, ok := e.Chains[chainSel]
		if !ok {
			return fmt.Errorf("chain not found for chain selector %d", chainSel)
		}
		deployRewardManager := cc.ChainsToDeploy[chainSel]
		_, err := changeset.DeployContract(e, ab, chain, RewardManagerDeployFn(deployRewardManager.LinkTokenAddress))
		if err != nil {
			return err
		}
		chainAddresses, err := ab.AddressesForChain(chain.Selector)
		if err != nil {
			e.Logger.Errorw("Failed to get chain addresses", "err", err)
			return err
		}
		chainState, err := changeset.LoadChainState(e.Logger, chain, chainAddresses)
		if err != nil {
			e.Logger.Errorw("Failed to load chain state", "err", err)
			return err
		}
		if len(chainState.RewardManagers) == 0 {
			errNoCCS := errors.New("no RewardManager on chain")
			e.Logger.Error(errNoCCS)
			return errNoCCS
		}
	}

	return nil
}

func RewardManagerDeployFn(linkAddress common.Address) changeset.ContractDeployFn[*rewardManager.RewardManager] {
	return func(chain deployment.Chain) *changeset.ContractDeployment[*rewardManager.RewardManager] {
		ccsAddr, ccsTx, ccs, err := rewardManager.DeployRewardManager(
			chain.DeployerKey,
			chain.Client,
			linkAddress,
		)
		if err != nil {
			return &changeset.ContractDeployment[*rewardManager.RewardManager]{
				Err: err,
			}
		}
		return &changeset.ContractDeployment[*rewardManager.RewardManager]{
			Address:  ccsAddr,
			Contract: ccs,
			Tx:       ccsTx,
			Tv:       deployment.NewTypeAndVersion(types.RewardManager, deployment.Version0_5_0),
			Err:      nil,
		}
	}
}
