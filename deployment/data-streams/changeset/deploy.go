package changeset

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink/deployment"
)

type (
	// Contract covers contracts such as channel_config_store.ChannelConfigStore and fee_manager.FeeManager.
	Contract interface {
		// Caller:
		Owner(opts *bind.CallOpts) (common.Address, error)
		TypeAndVersion(opts *bind.CallOpts) (string, error)

		// Transactor:
		AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)
		TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)
	}

	ContractDeployFn[C Contract] func(chain deployment.Chain) *ContractDeployment[C]

	ContractDeployment[C Contract] struct {
		Address  common.Address
		Contract C
		Tx       *types.Transaction
		Tv       deployment.TypeAndVersion
		Err      error
	}
)

var _ deployment.ChangeSetV2[DeployChannelConfigStoreConfig] = DeployChannelConfigStore{}

// DeployContract deploys a contract and saves the address to the address book.
//
// Note that this function modifies the given address book variable, so it should be passed by reference.
func DeployContract[C Contract](
	e deployment.Environment,
	ab deployment.AddressBook,
	chain deployment.Chain,
	deployFn ContractDeployFn[C],
) (*ContractDeployment[C], error) {
	contractDeployment := deployFn(chain)
	if contractDeployment.Err != nil {
		e.Logger.Errorw("Failed to deploy contract", "err", contractDeployment.Err, "chain", chain.Selector)
		return nil, contractDeployment.Err
	}
	_, err := chain.Confirm(contractDeployment.Tx)
	if err != nil {
		e.Logger.Errorw("Failed to confirm deployment", "err", err)
		return nil, err
	}
	e.Logger.Infow("Deployed contract", "Contract", contractDeployment.Tv.String(), "addr", contractDeployment.Address.String(), "chain", chain.String())
	err = ab.Save(chain.Selector, contractDeployment.Address.String(), contractDeployment.Tv)
	if err != nil {
		e.Logger.Errorw("Failed to save contract address", "err", err)
		return nil, err
	}
	return contractDeployment, nil
}
