// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package optimism_l2_output_oracle

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

type TypesOutputProposal struct {
	OutputRoot    [32]byte
	Timestamp     *big.Int
	L2BlockNumber *big.Int
}

var OptimismL2OutputOracleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_l2OutputIndex\",\"type\":\"uint256\"}],\"name\":\"getL2Output\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"outputRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint128\",\"name\":\"timestamp\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"l2BlockNumber\",\"type\":\"uint128\"}],\"internalType\":\"structTypes.OutputProposal\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_l2BlockNumber\",\"type\":\"uint256\"}],\"name\":\"getL2OutputIndexAfter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

var OptimismL2OutputOracleABI = OptimismL2OutputOracleMetaData.ABI

type OptimismL2OutputOracle struct {
	address common.Address
	abi     abi.ABI
	OptimismL2OutputOracleCaller
	OptimismL2OutputOracleTransactor
	OptimismL2OutputOracleFilterer
}

type OptimismL2OutputOracleCaller struct {
	contract *bind.BoundContract
}

type OptimismL2OutputOracleTransactor struct {
	contract *bind.BoundContract
}

type OptimismL2OutputOracleFilterer struct {
	contract *bind.BoundContract
}

type OptimismL2OutputOracleSession struct {
	Contract     *OptimismL2OutputOracle
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type OptimismL2OutputOracleCallerSession struct {
	Contract *OptimismL2OutputOracleCaller
	CallOpts bind.CallOpts
}

type OptimismL2OutputOracleTransactorSession struct {
	Contract     *OptimismL2OutputOracleTransactor
	TransactOpts bind.TransactOpts
}

type OptimismL2OutputOracleRaw struct {
	Contract *OptimismL2OutputOracle
}

type OptimismL2OutputOracleCallerRaw struct {
	Contract *OptimismL2OutputOracleCaller
}

type OptimismL2OutputOracleTransactorRaw struct {
	Contract *OptimismL2OutputOracleTransactor
}

func NewOptimismL2OutputOracle(address common.Address, backend bind.ContractBackend) (*OptimismL2OutputOracle, error) {
	abi, err := abi.JSON(strings.NewReader(OptimismL2OutputOracleABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindOptimismL2OutputOracle(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OptimismL2OutputOracle{address: address, abi: abi, OptimismL2OutputOracleCaller: OptimismL2OutputOracleCaller{contract: contract}, OptimismL2OutputOracleTransactor: OptimismL2OutputOracleTransactor{contract: contract}, OptimismL2OutputOracleFilterer: OptimismL2OutputOracleFilterer{contract: contract}}, nil
}

func NewOptimismL2OutputOracleCaller(address common.Address, caller bind.ContractCaller) (*OptimismL2OutputOracleCaller, error) {
	contract, err := bindOptimismL2OutputOracle(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OptimismL2OutputOracleCaller{contract: contract}, nil
}

func NewOptimismL2OutputOracleTransactor(address common.Address, transactor bind.ContractTransactor) (*OptimismL2OutputOracleTransactor, error) {
	contract, err := bindOptimismL2OutputOracle(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OptimismL2OutputOracleTransactor{contract: contract}, nil
}

func NewOptimismL2OutputOracleFilterer(address common.Address, filterer bind.ContractFilterer) (*OptimismL2OutputOracleFilterer, error) {
	contract, err := bindOptimismL2OutputOracle(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OptimismL2OutputOracleFilterer{contract: contract}, nil
}

func bindOptimismL2OutputOracle(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OptimismL2OutputOracleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_OptimismL2OutputOracle *OptimismL2OutputOracleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OptimismL2OutputOracle.Contract.OptimismL2OutputOracleCaller.contract.Call(opts, result, method, params...)
}

func (_OptimismL2OutputOracle *OptimismL2OutputOracleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OptimismL2OutputOracle.Contract.OptimismL2OutputOracleTransactor.contract.Transfer(opts)
}

func (_OptimismL2OutputOracle *OptimismL2OutputOracleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OptimismL2OutputOracle.Contract.OptimismL2OutputOracleTransactor.contract.Transact(opts, method, params...)
}

func (_OptimismL2OutputOracle *OptimismL2OutputOracleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OptimismL2OutputOracle.Contract.contract.Call(opts, result, method, params...)
}

func (_OptimismL2OutputOracle *OptimismL2OutputOracleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OptimismL2OutputOracle.Contract.contract.Transfer(opts)
}

func (_OptimismL2OutputOracle *OptimismL2OutputOracleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OptimismL2OutputOracle.Contract.contract.Transact(opts, method, params...)
}

func (_OptimismL2OutputOracle *OptimismL2OutputOracleCaller) GetL2Output(opts *bind.CallOpts, _l2OutputIndex *big.Int) (TypesOutputProposal, error) {
	var out []interface{}
	err := _OptimismL2OutputOracle.contract.Call(opts, &out, "getL2Output", _l2OutputIndex)

	if err != nil {
		return *new(TypesOutputProposal), err
	}

	out0 := *abi.ConvertType(out[0], new(TypesOutputProposal)).(*TypesOutputProposal)

	return out0, err

}

func (_OptimismL2OutputOracle *OptimismL2OutputOracleSession) GetL2Output(_l2OutputIndex *big.Int) (TypesOutputProposal, error) {
	return _OptimismL2OutputOracle.Contract.GetL2Output(&_OptimismL2OutputOracle.CallOpts, _l2OutputIndex)
}

func (_OptimismL2OutputOracle *OptimismL2OutputOracleCallerSession) GetL2Output(_l2OutputIndex *big.Int) (TypesOutputProposal, error) {
	return _OptimismL2OutputOracle.Contract.GetL2Output(&_OptimismL2OutputOracle.CallOpts, _l2OutputIndex)
}

func (_OptimismL2OutputOracle *OptimismL2OutputOracleCaller) GetL2OutputIndexAfter(opts *bind.CallOpts, _l2BlockNumber *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _OptimismL2OutputOracle.contract.Call(opts, &out, "getL2OutputIndexAfter", _l2BlockNumber)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_OptimismL2OutputOracle *OptimismL2OutputOracleSession) GetL2OutputIndexAfter(_l2BlockNumber *big.Int) (*big.Int, error) {
	return _OptimismL2OutputOracle.Contract.GetL2OutputIndexAfter(&_OptimismL2OutputOracle.CallOpts, _l2BlockNumber)
}

func (_OptimismL2OutputOracle *OptimismL2OutputOracleCallerSession) GetL2OutputIndexAfter(_l2BlockNumber *big.Int) (*big.Int, error) {
	return _OptimismL2OutputOracle.Contract.GetL2OutputIndexAfter(&_OptimismL2OutputOracle.CallOpts, _l2BlockNumber)
}

func (_OptimismL2OutputOracle *OptimismL2OutputOracle) Address() common.Address {
	return _OptimismL2OutputOracle.address
}

type OptimismL2OutputOracleInterface interface {
	GetL2Output(opts *bind.CallOpts, _l2OutputIndex *big.Int) (TypesOutputProposal, error)

	GetL2OutputIndexAfter(opts *bind.CallOpts, _l2BlockNumber *big.Int) (*big.Int, error)

	Address() common.Address
}
