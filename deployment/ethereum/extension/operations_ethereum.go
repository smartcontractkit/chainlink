package deployment_ethereum

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/deployment"
)

// Contract should be a contract binding
type DeployContractBindingParamsFn[Contract any] func(auth *bind.TransactOpts, backend bind.ContractBackend, params ...any) (common.Address, *types.Transaction, *Contract, error)
type DeployContractBindingNoParamsFn[Contract any] func(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Contract, error)

type ContractCtorFn func(address common.Address, backend bind.ContractBackend) (*bind.BoundContract, error)

// TODO: Add address...
type EthereumTxOutput struct {
	Tx      *types.Transaction
	Address common.Address
}

type EthereumDeps struct {
	Auth    *bind.TransactOpts
	Client  bind.ContractBackend
	Confirm func(tx *types.Transaction) (uint64, error)
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

func NewEthDeployOperationFromBinding[I EthInput, C any](binding DeployContractBindingParamsFn[C], version string) *deployment.Operation[I, EthereumTxOutput, EthereumDeps] {
	return deployment.NewOperation[I](version, "Deploy Contract Operation", func(ctx deployment.Context[EthereumDeps], input I) (EthereumTxOutput, error) {
		address, tx, _, err := binding(ctx.Deps.Auth, ctx.Deps.Client, input.GetOrderedParams()...)
		if err != nil {
			return EthereumTxOutput{}, err
		}

		_, err = ctx.Deps.Confirm(tx)
		if err != nil {
			return EthereumTxOutput{}, err
		}

		return EthereumTxOutput{Tx: tx, Address: address}, nil
	})
}

func NewEthDeployOperationFromBindingNoParams[C any](binding DeployContractBindingNoParamsFn[C], version string) *deployment.Operation[deployment.EmptyInput, EthereumTxOutput, EthereumDeps] {
	return deployment.NewOperation(version, "Deploy Contract Operation", func(ctx deployment.Context[EthereumDeps], input deployment.EmptyInput) (EthereumTxOutput, error) {
		address, tx, _, err := binding(ctx.Deps.Auth, ctx.Deps.Client)
		if err != nil {
			return EthereumTxOutput{}, err
		}

		_, err = ctx.Deps.Confirm(tx)
		if err != nil {
			return EthereumTxOutput{}, err
		}

		return EthereumTxOutput{Tx: tx, Address: address}, nil
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

		_, err = ctx.Deps.Confirm(tx)
		if err != nil {
			return EthereumTxOutput{}, err
		}

		return EthereumTxOutput{Tx: tx}, nil
	})
}
