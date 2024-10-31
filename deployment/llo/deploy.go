package llo

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/channel_config_store"
)

type DeployLLOContractConfig struct {
	ChainsToDeploy []uint64 // Chain Selectors
}

// LLOContract covers contracts such as channel_config_store.ChannelConfigStore and fee_manager.FeeManager.
type LLOContract interface {
	// Caller:
	Owner(opts *bind.CallOpts) (common.Address, error)
	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)
	TypeAndVersion(opts *bind.CallOpts) (string, error)
	// Transactor:
	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)
	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)
}

type ContractDeploy[C LLOContract] struct {
	Address  common.Address
	Contract C
	Tx       *types.Transaction
	Tv       deployment.TypeAndVersion
	Err      error
}

var (
	ChannelConfigStore deployment.ContractType = "ChannelConfigStore"
)

func DeployChannelConfigStore(e deployment.Environment, ab deployment.AddressBook, c DeployLLOContractConfig) error {
	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	if err != nil || len(nodes) == 0 {
		e.Logger.Errorw("Failed to get node info", "err", err)
		return err
	}

	for _, chainSel := range c.ChainsToDeploy {
		chain, ok := e.Chains[chainSel]
		if !ok {
			return fmt.Errorf("Chain %d not found", chainSel)
		}
		_, err = deployChannelConfigStoreToChain(e, chain, ab)
		if err != nil {
			return err
		}
		chainAddresses, err := ab.AddressesForChain(chain.Selector)
		if err != nil {
			e.Logger.Errorw("Failed to get chain addresses", "err", err)
			return err
		}
		chainState, err := LoadChainState(chain, chainAddresses)
		if err != nil {
			e.Logger.Errorw("Failed to load chain state", "err", err)
			return err
		}
		if chainState.ChannelConfigStore == nil {
			errNoCCS := errors.New("no ChannelConfigStore on chain")
			e.Logger.Error(errNoCCS)
			return errNoCCS
		}
	}

	return nil
}

// deployChannelConfigStoreToChain deploys ChannelConfigStore to a specific chain.
//
// Note that this function modifies the given address book variable.
func deployChannelConfigStoreToChain(e deployment.Environment, chain deployment.Chain, ab deployment.AddressBook) (*ContractDeploy[*channel_config_store.ChannelConfigStore], error) {
	return deployContract(e.Logger, chain, ab, func(chain deployment.Chain) ContractDeploy[*channel_config_store.ChannelConfigStore] {
		ccsAddr, ccsTx, ccs, err := channel_config_store.DeployChannelConfigStore(
			chain.DeployerKey,
			chain.Client,
		)
		if err != nil {
			return ContractDeploy[*channel_config_store.ChannelConfigStore]{
				Err: err,
			}
		}
		return ContractDeploy[*channel_config_store.ChannelConfigStore]{
			Address:  ccsAddr,
			Contract: ccs,
			Tx:       ccsTx,
			Tv:       deployment.NewTypeAndVersion(ChannelConfigStore, deployment.Version1_0_0),
			Err:      nil,
		}
	})
}

func deployContract[C LLOContract](
	lggr logger.Logger,
	chain deployment.Chain,
	addressBook deployment.AddressBook,
	deploy func(chain deployment.Chain) ContractDeploy[C],
) (*ContractDeploy[C], error) {
	contractDeploy := deploy(chain)
	if contractDeploy.Err != nil {
		lggr.Errorw("Failed to deploy contract", "err", contractDeploy.Err)
		return nil, contractDeploy.Err
	}
	_, err := chain.Confirm(contractDeploy.Tx)
	if err != nil {
		lggr.Errorw("Failed to confirm deployment", "err", err)
		return nil, err
	}
	err = addressBook.Save(chain.Selector, contractDeploy.Address.String(), contractDeploy.Tv)
	if err != nil {
		lggr.Errorw("Failed to save contract address", "err", err)
		return nil, err
	}
	return &contractDeploy, nil
}
