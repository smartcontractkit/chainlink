package reward_manager

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/reward_manager_v0_5_0"
)

// DeployRewardManagerChangeset deploys RewardManager to the chains specified in the config.
var DeployRewardManagerChangeset deployment.ChangeSetV2[DeployRewardManagerConfig] = &rewardManagerDeploy{}

type rewardManagerDeploy struct{}

type DeployRewardManagerConfig struct {
	// ChainsToDeploy is a list of chain selectors to deploy the contract to.
	ChainsToDeploy   []uint64
	LinkTokenAddress common.Address
}

func (cc DeployRewardManagerConfig) Validate() error {
	if len(cc.ChainsToDeploy) == 0 {
		return errors.New("ChainsToDeploy is empty")
	}
	for _, chain := range cc.ChainsToDeploy {
		if err := deployment.IsValidChainSelector(chain); err != nil {
			return fmt.Errorf("invalid chain selector: %d - %w", chain, err)
		}
	}
	return nil
}

func (rewardManagerDeploy) Apply(e deployment.Environment, cc DeployRewardManagerConfig) (deployment.ChangesetOutput, error) {
	ab := deployment.NewMemoryAddressBook()
	err := DeployRewardManager(e, ab, cc)
	if err != nil {
		e.Logger.Errorw("Failed to deploy RewardManager", "err", err, "addresses", ab)
		return deployment.ChangesetOutput{AddressBook: ab}, deployment.MaybeDataErr(err)
	}
	return deployment.ChangesetOutput{
		AddressBook: ab,
	}, nil
}

func (rewardManagerDeploy) VerifyPreconditions(_ deployment.Environment, cc DeployRewardManagerConfig) error {
	if err := cc.Validate(); err != nil {
		return fmt.Errorf("invalid DeployRewardManagerConfig: %w", err)
	}
	return nil
}

func DeployRewardManager(e deployment.Environment, ab deployment.AddressBook, cc DeployRewardManagerConfig) error {
	if err := cc.Validate(); err != nil {
		return fmt.Errorf("invalid DeployRewardManagerConfig: %w", err)
	}

	for _, chainSel := range cc.ChainsToDeploy {
		chain, ok := e.Chains[chainSel]
		if !ok {
			return fmt.Errorf("Chain not found for chain selector %d", chainSel)
		}
		_, err := changeset.DeployContract[*reward_manager_v0_5_0.RewardManager](e, ab, chain, RewardManagerDeployFn(cc))
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

// RewardManagerDeployFn returns a function that deploys a RewardManager contract.
func RewardManagerDeployFn(cfg DeployRewardManagerConfig) changeset.ContractDeployFn[*reward_manager_v0_5_0.RewardManager] {
	return func(chain deployment.Chain) *changeset.ContractDeployment[*reward_manager_v0_5_0.RewardManager] {
		ccsAddr, ccsTx, ccs, err := reward_manager_v0_5_0.DeployRewardManager(
			chain.DeployerKey,
			chain.Client,
			cfg.LinkTokenAddress,
		)
		if err != nil {
			return &changeset.ContractDeployment[*reward_manager_v0_5_0.RewardManager]{
				Err: err,
			}
		}
		return &changeset.ContractDeployment[*reward_manager_v0_5_0.RewardManager]{
			Address:  ccsAddr,
			Contract: ccs,
			Tx:       ccsTx,
			Tv:       deployment.NewTypeAndVersion(types.RewardManager, deployment.Version0_5_0),
			Err:      nil,
		}
	}
}
