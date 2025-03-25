package v0_5_0

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/fee_manager_v0_5_0"
)

var DeployFeeManagerChangeset = deployment.CreateChangeSet(deployFeeManagerLogic, deployFeeManagerPrecondition)

type DeployFeeManager struct {
	LinkTokenAddress     common.Address
	NativeTokenAddress   common.Address
	ProxyAddress         common.Address
	RewardManagerAddress common.Address
}

type DeployFeeManagerConfig struct {
	ChainsToDeploy map[uint64]DeployFeeManager
}

func (cc DeployFeeManagerConfig) Validate() error {
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

func deployFeeManagerLogic(e deployment.Environment, cc DeployFeeManagerConfig) (deployment.ChangesetOutput, error) {
	ab := deployment.NewMemoryAddressBook()
	err := deployFeeManager(e, ab, cc)
	if err != nil {
		e.Logger.Errorw("Failed to deploy FeeManager", "err", err, "addresses", ab)
		return deployment.ChangesetOutput{AddressBook: ab}, deployment.MaybeDataErr(err)
	}
	return deployment.ChangesetOutput{
		AddressBook: ab,
	}, nil
}

func deployFeeManagerPrecondition(_ deployment.Environment, cc DeployFeeManagerConfig) error {
	return cc.Validate()
}

func deployFeeManager(e deployment.Environment, ab deployment.AddressBook, cc DeployFeeManagerConfig) error {
	if err := cc.Validate(); err != nil {
		return fmt.Errorf("invalid DeployFeeManagerConfig: %w", err)
	}

	for chainSel := range cc.ChainsToDeploy {
		chain, ok := e.Chains[chainSel]
		if !ok {
			return fmt.Errorf("chain not found for chain selector %d", chainSel)
		}
		conf := cc.ChainsToDeploy[chainSel]
		_, err := changeset.DeployContract[*fee_manager_v0_5_0.FeeManager](e, ab, chain, FeeManagerDeployFn(conf))
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
		if len(chainState.FeeManagers) == 0 {
			errNoCCS := errors.New("no FeeManager on chain")
			e.Logger.Error(errNoCCS)
			return errNoCCS
		}
	}

	return nil
}

// FeeManagerDeployFn returns a function that deploys a FeeManager contract.
func FeeManagerDeployFn(cfg DeployFeeManager) changeset.ContractDeployFn[*fee_manager_v0_5_0.FeeManager] {
	return func(chain deployment.Chain) *changeset.ContractDeployment[*fee_manager_v0_5_0.FeeManager] {
		ccsAddr, ccsTx, ccs, err := fee_manager_v0_5_0.DeployFeeManager(
			chain.DeployerKey,
			chain.Client,
			cfg.LinkTokenAddress,
			cfg.NativeTokenAddress,
			cfg.ProxyAddress,
			cfg.RewardManagerAddress,
		)
		if err != nil {
			return &changeset.ContractDeployment[*fee_manager_v0_5_0.FeeManager]{
				Err: err,
			}
		}
		return &changeset.ContractDeployment[*fee_manager_v0_5_0.FeeManager]{
			Address:  ccsAddr,
			Contract: ccs,
			Tx:       ccsTx,
			Tv:       deployment.NewTypeAndVersion(types.FeeManager, deployment.Version0_5_0),
			Err:      nil,
		}
	}
}
