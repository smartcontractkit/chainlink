// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package optimism_portal

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

type TypesOutputRootProof struct {
	Version                  [32]byte
	StateRoot                [32]byte
	MessagePasserStorageRoot [32]byte
	LatestBlockhash          [32]byte
}

type TypesWithdrawalTransaction struct {
	Nonce    *big.Int
	Sender   common.Address
	Target   common.Address
	Value    *big.Int
	GasLimit *big.Int
	Data     []byte
}

var OptimismPortalMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structTypes.WithdrawalTransaction\",\"name\":\"_tx\",\"type\":\"tuple\"}],\"name\":\"finalizeWithdrawalTransaction\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structTypes.WithdrawalTransaction\",\"name\":\"_tx\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"_l2OutputIndex\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"version\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"messagePasserStorageRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"latestBlockhash\",\"type\":\"bytes32\"}],\"internalType\":\"structTypes.OutputRootProof\",\"name\":\"_outputRootProof\",\"type\":\"tuple\"},{\"internalType\":\"bytes[]\",\"name\":\"_withdrawalProof\",\"type\":\"bytes[]\"}],\"name\":\"proveWithdrawalTransaction\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

var OptimismPortalABI = OptimismPortalMetaData.ABI

type OptimismPortal struct {
	address common.Address
	abi     abi.ABI
	OptimismPortalCaller
	OptimismPortalTransactor
	OptimismPortalFilterer
}

type OptimismPortalCaller struct {
	contract *bind.BoundContract
}

type OptimismPortalTransactor struct {
	contract *bind.BoundContract
}

type OptimismPortalFilterer struct {
	contract *bind.BoundContract
}

type OptimismPortalSession struct {
	Contract     *OptimismPortal
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type OptimismPortalCallerSession struct {
	Contract *OptimismPortalCaller
	CallOpts bind.CallOpts
}

type OptimismPortalTransactorSession struct {
	Contract     *OptimismPortalTransactor
	TransactOpts bind.TransactOpts
}

type OptimismPortalRaw struct {
	Contract *OptimismPortal
}

type OptimismPortalCallerRaw struct {
	Contract *OptimismPortalCaller
}

type OptimismPortalTransactorRaw struct {
	Contract *OptimismPortalTransactor
}

func NewOptimismPortal(address common.Address, backend bind.ContractBackend) (*OptimismPortal, error) {
	abi, err := abi.JSON(strings.NewReader(OptimismPortalABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindOptimismPortal(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OptimismPortal{address: address, abi: abi, OptimismPortalCaller: OptimismPortalCaller{contract: contract}, OptimismPortalTransactor: OptimismPortalTransactor{contract: contract}, OptimismPortalFilterer: OptimismPortalFilterer{contract: contract}}, nil
}

func NewOptimismPortalCaller(address common.Address, caller bind.ContractCaller) (*OptimismPortalCaller, error) {
	contract, err := bindOptimismPortal(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OptimismPortalCaller{contract: contract}, nil
}

func NewOptimismPortalTransactor(address common.Address, transactor bind.ContractTransactor) (*OptimismPortalTransactor, error) {
	contract, err := bindOptimismPortal(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OptimismPortalTransactor{contract: contract}, nil
}

func NewOptimismPortalFilterer(address common.Address, filterer bind.ContractFilterer) (*OptimismPortalFilterer, error) {
	contract, err := bindOptimismPortal(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OptimismPortalFilterer{contract: contract}, nil
}

func bindOptimismPortal(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OptimismPortalMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_OptimismPortal *OptimismPortalRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OptimismPortal.Contract.OptimismPortalCaller.contract.Call(opts, result, method, params...)
}

func (_OptimismPortal *OptimismPortalRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OptimismPortal.Contract.OptimismPortalTransactor.contract.Transfer(opts)
}

func (_OptimismPortal *OptimismPortalRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OptimismPortal.Contract.OptimismPortalTransactor.contract.Transact(opts, method, params...)
}

func (_OptimismPortal *OptimismPortalCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OptimismPortal.Contract.contract.Call(opts, result, method, params...)
}

func (_OptimismPortal *OptimismPortalTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OptimismPortal.Contract.contract.Transfer(opts)
}

func (_OptimismPortal *OptimismPortalTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OptimismPortal.Contract.contract.Transact(opts, method, params...)
}

func (_OptimismPortal *OptimismPortalCaller) Version(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _OptimismPortal.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_OptimismPortal *OptimismPortalSession) Version() (string, error) {
	return _OptimismPortal.Contract.Version(&_OptimismPortal.CallOpts)
}

func (_OptimismPortal *OptimismPortalCallerSession) Version() (string, error) {
	return _OptimismPortal.Contract.Version(&_OptimismPortal.CallOpts)
}

func (_OptimismPortal *OptimismPortalTransactor) FinalizeWithdrawalTransaction(opts *bind.TransactOpts, _tx TypesWithdrawalTransaction) (*types.Transaction, error) {
	return _OptimismPortal.contract.Transact(opts, "finalizeWithdrawalTransaction", _tx)
}

func (_OptimismPortal *OptimismPortalSession) FinalizeWithdrawalTransaction(_tx TypesWithdrawalTransaction) (*types.Transaction, error) {
	return _OptimismPortal.Contract.FinalizeWithdrawalTransaction(&_OptimismPortal.TransactOpts, _tx)
}

func (_OptimismPortal *OptimismPortalTransactorSession) FinalizeWithdrawalTransaction(_tx TypesWithdrawalTransaction) (*types.Transaction, error) {
	return _OptimismPortal.Contract.FinalizeWithdrawalTransaction(&_OptimismPortal.TransactOpts, _tx)
}

func (_OptimismPortal *OptimismPortalTransactor) ProveWithdrawalTransaction(opts *bind.TransactOpts, _tx TypesWithdrawalTransaction, _l2OutputIndex *big.Int, _outputRootProof TypesOutputRootProof, _withdrawalProof [][]byte) (*types.Transaction, error) {
	return _OptimismPortal.contract.Transact(opts, "proveWithdrawalTransaction", _tx, _l2OutputIndex, _outputRootProof, _withdrawalProof)
}

func (_OptimismPortal *OptimismPortalSession) ProveWithdrawalTransaction(_tx TypesWithdrawalTransaction, _l2OutputIndex *big.Int, _outputRootProof TypesOutputRootProof, _withdrawalProof [][]byte) (*types.Transaction, error) {
	return _OptimismPortal.Contract.ProveWithdrawalTransaction(&_OptimismPortal.TransactOpts, _tx, _l2OutputIndex, _outputRootProof, _withdrawalProof)
}

func (_OptimismPortal *OptimismPortalTransactorSession) ProveWithdrawalTransaction(_tx TypesWithdrawalTransaction, _l2OutputIndex *big.Int, _outputRootProof TypesOutputRootProof, _withdrawalProof [][]byte) (*types.Transaction, error) {
	return _OptimismPortal.Contract.ProveWithdrawalTransaction(&_OptimismPortal.TransactOpts, _tx, _l2OutputIndex, _outputRootProof, _withdrawalProof)
}

func (_OptimismPortal *OptimismPortal) Address() common.Address {
	return _OptimismPortal.address
}

type OptimismPortalInterface interface {
	Version(opts *bind.CallOpts) (string, error)

	FinalizeWithdrawalTransaction(opts *bind.TransactOpts, _tx TypesWithdrawalTransaction) (*types.Transaction, error)

	ProveWithdrawalTransaction(opts *bind.TransactOpts, _tx TypesWithdrawalTransaction, _l2OutputIndex *big.Int, _outputRootProof TypesOutputRootProof, _withdrawalProof [][]byte) (*types.Transaction, error)

	Address() common.Address
}
