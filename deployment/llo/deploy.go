package llo

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/channel_config_store"
)

type (
	DeployLLOContractConfig struct {
		// ChainsToDeploy is a list of chain selectors to deploy the contract to.
		ChainsToDeploy []uint64
	}

	// LLOContract covers contracts such as channel_config_store.ChannelConfigStore and fee_manager.FeeManager.
	LLOContract interface {
		// Caller:
		Owner(opts *bind.CallOpts) (common.Address, error)
		SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)
		TypeAndVersion(opts *bind.CallOpts) (string, error)

		// Transactor:
		AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)
		TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)
	}

	ContractDeployFn[C LLOContract] func(chain deployment.Chain) *ContractDeployment[C]

	ContractDeployment[C LLOContract] struct {
		Address  common.Address
		Contract C
		Tx       *types.Transaction
		Tv       deployment.TypeAndVersion
		Err      error
	}
)

const (
	ChannelConfigStore deployment.ContractType = "ChannelConfigStore"
)

// DeployChannelConfigStore deploys ChannelConfigStore to the chains specified in the config.
//
// Note that this function modifies the given address book variable, so it should be passed by reference.
func DeployChannelConfigStore(env deployment.Environment, ab deployment.AddressBook, cc DeployLLOContractConfig) error {
	nodes, err := deployment.NodeInfo(env.NodeIDs, env.Offchain)
	if err != nil || len(nodes) == 0 {
		env.Logger.Errorw("Failed to get node info", "err", err)
		return err
	}

	for _, chainSel := range cc.ChainsToDeploy {
		chain, ok := env.Chains[chainSel]
		if !ok {
			return fmt.Errorf("Chain not found for chain selector %d", chainSel)
		}
		_, err = deployContract[*channel_config_store.ChannelConfigStore](env, ab, chain, channelConfigStoreDeployFn())
		if err != nil {
			return err
		}
		chainAddresses, err := ab.AddressesForChain(chain.Selector)
		if err != nil {
			env.Logger.Errorw("Failed to get chain addresses", "err", err)
			return err
		}
		chainState, err := LoadChainConfig(chain, chainAddresses)
		if err != nil {
			env.Logger.Errorw("Failed to load chain state", "err", err)
			return err
		}
		if chainState.ChannelConfigStores == nil || len(chainState.ChannelConfigStores[chain.Selector]) == 0 {
			errNoCCS := errors.New("no ChannelConfigStore on chain")
			env.Logger.Error(errNoCCS)
			return errNoCCS
		}
	}

	return nil
}

// channelConfigStoreDeployFn returns a function that deploys a ChannelConfigStore contract.
func channelConfigStoreDeployFn() ContractDeployFn[*channel_config_store.ChannelConfigStore] {
	return func(chain deployment.Chain) *ContractDeployment[*channel_config_store.ChannelConfigStore] {
		ccsAddr, ccsTx, ccs, err := channel_config_store.DeployChannelConfigStore(
			chain.DeployerKey,
			chain.Client,
		)
		if err != nil {
			return &ContractDeployment[*channel_config_store.ChannelConfigStore]{
				Err: err,
			}
		}
		return &ContractDeployment[*channel_config_store.ChannelConfigStore]{
			Address:  ccsAddr,
			Contract: ccs,
			Tx:       ccsTx,
			Tv:       deployment.NewTypeAndVersion(ChannelConfigStore, deployment.Version1_0_0),
			Err:      nil,
		}
	}
}

// deployContract deploys a contract and saves the address to the address book.
//
// Note that this function modifies the given address book variable, so it should be passed by reference.
func deployContract[C LLOContract](
	env deployment.Environment,
	ab deployment.AddressBook,
	chain deployment.Chain,
	deployFn ContractDeployFn[C],
) (*ContractDeployment[C], error) {
	contractDeployment := deployFn(chain)
	if contractDeployment.Err != nil {
		env.Logger.Errorw("Failed to deploy contract", "err", contractDeployment.Err)
		return nil, contractDeployment.Err
	}
	_, err := chain.Confirm(contractDeployment.Tx)
	if err != nil {
		env.Logger.Errorw("Failed to confirm deployment", "err", err)
		return nil, err
	}
	err = ab.Save(chain.Selector, contractDeployment.Address.String(), contractDeployment.Tv)
	if err != nil {
		env.Logger.Errorw("Failed to save contract address", "err", err)
		return nil, err
	}
	return contractDeployment, nil
}
