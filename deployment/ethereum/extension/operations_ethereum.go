package deployment_ethereum

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/deployment"
)

// Contract should be a contract binding
type DeployContractBindingFn[I, Contract any] func(
	auth *bind.TransactOpts,
	backend bind.ContractBackend,
	input I,
) (common.Address, *types.Transaction, *Contract, error)
type DeployContractBindingNoParamsFn[Contract any] func(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Contract, error)

// TODO: Add address...
type EthereumTxOutput struct {
	Hash    common.Hash
	Address common.Address

	RawReceipt any
}

type EthereumDeps struct {
	Auth    *bind.TransactOpts
	Client  deployment.OnchainClient
	Confirm func(client deployment.OnchainClient, hash common.Hash) (*types.Receipt, error)
}

// EthInput is a common interface for all Ethereum operation inputs
type EthInput interface {
	// Required to pass the parameters in the correct order
	GetOrderedParams() []any
}

type EthMethodInput interface {
	EthInput
	Address() common.Address
}

func NewEthDeployOperationFromBinding[I, C any](
	binding DeployContractBindingFn[I, C],
	version string,
) *deployment.Operation[I, EthereumTxOutput, EthereumDeps] {
	return deployment.NewOperation[I](version, "Deploy Contract Operation",
		func(ctx deployment.Context[EthereumDeps], input I) (EthereumTxOutput, error) {
			address, tx, _, err := binding(ctx.Deps.Auth, ctx.Deps.Client, input)
			if err != nil {
				return EthereumTxOutput{}, err
			}

			hash := tx.Hash()
			rec, err := ctx.Deps.Confirm(ctx.Deps.Client, hash)
			if err != nil {
				return EthereumTxOutput{}, err
			}

			return EthereumTxOutput{Hash: hash, Address: address, RawReceipt: rec}, nil
		})
}

func NewEthDeployOperationFromBindingNoParams[C any](binding DeployContractBindingNoParamsFn[C], version string) *deployment.Operation[deployment.EmptyInput, EthereumTxOutput, EthereumDeps] {
	return deployment.NewOperation(version, "Deploy Contract Operation", func(ctx deployment.Context[EthereumDeps], input deployment.EmptyInput) (EthereumTxOutput, error) {
		address, tx, _, err := binding(ctx.Deps.Auth, ctx.Deps.Client)
		if err != nil {
			return EthereumTxOutput{}, err
		}

		hash := tx.Hash()
		rec, err := ctx.Deps.Confirm(ctx.Deps.Client, hash)
		if err != nil {
			return EthereumTxOutput{}, err
		}

		return EthereumTxOutput{Hash: hash, Address: address, RawReceipt: rec}, nil
	})
}

func NewEthOperationFromBinding[I EthMethodInput](metadata *bind.MetaData, version string, method string) *deployment.Operation[I, EthereumTxOutput, EthereumDeps] {
	return deployment.NewOperation[I](version, "Contract Transactor Operation", func(ctx deployment.Context[EthereumDeps], input I) (EthereumTxOutput, error) {
		parsed, err := metadata.GetAbi()
		if err != nil {
			return EthereumTxOutput{}, err
		}

		contract, err := bind.NewBoundContract(input.Address(), *parsed, nil, ctx.Deps.Client, nil), nil
		if err != nil {
			return EthereumTxOutput{}, err
		}

		tx, err := contract.Transact(ctx.Deps.Auth, method, input.GetOrderedParams()...)
		if err != nil {
			return EthereumTxOutput{}, err
		}

		hash := tx.Hash()
		rec, err := ctx.Deps.Confirm(ctx.Deps.Client, hash)
		if err != nil {
			return EthereumTxOutput{}, err
		}

		return EthereumTxOutput{Hash: hash, Address: input.Address(), RawReceipt: rec}, nil
	})
}
