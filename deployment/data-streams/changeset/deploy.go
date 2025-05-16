package changeset

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/mcms"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	ds "github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/metadata"
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

	ContractDeployFn[C Contract] func(chain cldf.Chain) *ContractDeployment[C]

	ContractDeployment[C Contract] struct {
		Address  common.Address
		Contract C
		Tx       *types.Transaction
		Tv       cldf.TypeAndVersion
		Err      error
		Block    uint64
	}
)

type DeployOptions struct {
	ContractMetadata *metadata.SerializedContractMetadata
}

type DeploymentOutput[C Contract] struct {
	Deployments       []*ContractDeployment[C]
	TimelockProposals []mcms.TimelockProposal
}

// DeployContract deploys a contract and saves the address to datastore.
func DeployContract[C Contract](
	e cldf.Environment,
	dataStore ds.MutableDataStore[metadata.SerializedContractMetadata, ds.DefaultMetadata],
	chain cldf.Chain,
	deployFn ContractDeployFn[C],
	options *DeployOptions,
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

	// Store Address
	if err = dataStore.Addresses().Add(
		ds.AddressRef{
			ChainSelector: chain.Selector,
			Address:       contractDeployment.Address.String(),
			Type:          ds.ContractType(contractDeployment.Tv.Type),
			Version:       &contractDeployment.Tv.Version,
		},
	); err != nil {
		e.Logger.Errorw("Failed to save contract address", "err", err)
		return nil, err
	}

	if options != nil && options.ContractMetadata != nil {
		// Add a new CommonContractMetadata entry for the newly deployed contract
		if err = dataStore.ContractMetadata().Add(
			ds.ContractMetadata[metadata.SerializedContractMetadata]{
				ChainSelector: chain.Selector,
				Address:       contractDeployment.Address.String(),
				Metadata:      *options.ContractMetadata,
			},
		); err != nil {
			return nil, fmt.Errorf("failed to save contract metadata: %w", err)
		}
	}

	return contractDeployment, nil
}
