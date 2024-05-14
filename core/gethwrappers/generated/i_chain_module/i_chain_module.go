// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package i_chain_module

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

var IChainModuleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"blockHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"blockNumber\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCurrentL1Fee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getGasOverhead\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"chainModuleFixedOverhead\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"chainModulePerByteOverhead\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"dataSize\",\"type\":\"uint256\"}],\"name\":\"getMaxL1Fee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

var IChainModuleABI = IChainModuleMetaData.ABI

type IChainModule struct {
	address common.Address
	abi     abi.ABI
	IChainModuleCaller
	IChainModuleTransactor
	IChainModuleFilterer
}

type IChainModuleCaller struct {
	contract *bind.BoundContract
}

type IChainModuleTransactor struct {
	contract *bind.BoundContract
}

type IChainModuleFilterer struct {
	contract *bind.BoundContract
}

type IChainModuleSession struct {
	Contract     *IChainModule
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type IChainModuleCallerSession struct {
	Contract *IChainModuleCaller
	CallOpts bind.CallOpts
}

type IChainModuleTransactorSession struct {
	Contract     *IChainModuleTransactor
	TransactOpts bind.TransactOpts
}

type IChainModuleRaw struct {
	Contract *IChainModule
}

type IChainModuleCallerRaw struct {
	Contract *IChainModuleCaller
}

type IChainModuleTransactorRaw struct {
	Contract *IChainModuleTransactor
}

func NewIChainModule(address common.Address, backend bind.ContractBackend) (*IChainModule, error) {
	abi, err := abi.JSON(strings.NewReader(IChainModuleABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindIChainModule(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IChainModule{address: address, abi: abi, IChainModuleCaller: IChainModuleCaller{contract: contract}, IChainModuleTransactor: IChainModuleTransactor{contract: contract}, IChainModuleFilterer: IChainModuleFilterer{contract: contract}}, nil
}

func NewIChainModuleCaller(address common.Address, caller bind.ContractCaller) (*IChainModuleCaller, error) {
	contract, err := bindIChainModule(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IChainModuleCaller{contract: contract}, nil
}

func NewIChainModuleTransactor(address common.Address, transactor bind.ContractTransactor) (*IChainModuleTransactor, error) {
	contract, err := bindIChainModule(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IChainModuleTransactor{contract: contract}, nil
}

func NewIChainModuleFilterer(address common.Address, filterer bind.ContractFilterer) (*IChainModuleFilterer, error) {
	contract, err := bindIChainModule(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IChainModuleFilterer{contract: contract}, nil
}

func bindIChainModule(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IChainModuleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_IChainModule *IChainModuleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IChainModule.Contract.IChainModuleCaller.contract.Call(opts, result, method, params...)
}

func (_IChainModule *IChainModuleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IChainModule.Contract.IChainModuleTransactor.contract.Transfer(opts)
}

func (_IChainModule *IChainModuleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IChainModule.Contract.IChainModuleTransactor.contract.Transact(opts, method, params...)
}

func (_IChainModule *IChainModuleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IChainModule.Contract.contract.Call(opts, result, method, params...)
}

func (_IChainModule *IChainModuleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IChainModule.Contract.contract.Transfer(opts)
}

func (_IChainModule *IChainModuleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IChainModule.Contract.contract.Transact(opts, method, params...)
}

func (_IChainModule *IChainModuleCaller) BlockHash(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _IChainModule.contract.Call(opts, &out, "blockHash", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_IChainModule *IChainModuleSession) BlockHash(arg0 *big.Int) ([32]byte, error) {
	return _IChainModule.Contract.BlockHash(&_IChainModule.CallOpts, arg0)
}

func (_IChainModule *IChainModuleCallerSession) BlockHash(arg0 *big.Int) ([32]byte, error) {
	return _IChainModule.Contract.BlockHash(&_IChainModule.CallOpts, arg0)
}

func (_IChainModule *IChainModuleCaller) BlockNumber(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IChainModule.contract.Call(opts, &out, "blockNumber")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IChainModule *IChainModuleSession) BlockNumber() (*big.Int, error) {
	return _IChainModule.Contract.BlockNumber(&_IChainModule.CallOpts)
}

func (_IChainModule *IChainModuleCallerSession) BlockNumber() (*big.Int, error) {
	return _IChainModule.Contract.BlockNumber(&_IChainModule.CallOpts)
}

func (_IChainModule *IChainModuleCaller) GetCurrentL1Fee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IChainModule.contract.Call(opts, &out, "getCurrentL1Fee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IChainModule *IChainModuleSession) GetCurrentL1Fee() (*big.Int, error) {
	return _IChainModule.Contract.GetCurrentL1Fee(&_IChainModule.CallOpts)
}

func (_IChainModule *IChainModuleCallerSession) GetCurrentL1Fee() (*big.Int, error) {
	return _IChainModule.Contract.GetCurrentL1Fee(&_IChainModule.CallOpts)
}

func (_IChainModule *IChainModuleCaller) GetGasOverhead(opts *bind.CallOpts) (GetGasOverhead,

	error) {
	var out []interface{}
	err := _IChainModule.contract.Call(opts, &out, "getGasOverhead")

	outstruct := new(GetGasOverhead)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ChainModuleFixedOverhead = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.ChainModulePerByteOverhead = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_IChainModule *IChainModuleSession) GetGasOverhead() (GetGasOverhead,

	error) {
	return _IChainModule.Contract.GetGasOverhead(&_IChainModule.CallOpts)
}

func (_IChainModule *IChainModuleCallerSession) GetGasOverhead() (GetGasOverhead,

	error) {
	return _IChainModule.Contract.GetGasOverhead(&_IChainModule.CallOpts)
}

func (_IChainModule *IChainModuleCaller) GetMaxL1Fee(opts *bind.CallOpts, dataSize *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _IChainModule.contract.Call(opts, &out, "getMaxL1Fee", dataSize)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_IChainModule *IChainModuleSession) GetMaxL1Fee(dataSize *big.Int) (*big.Int, error) {
	return _IChainModule.Contract.GetMaxL1Fee(&_IChainModule.CallOpts, dataSize)
}

func (_IChainModule *IChainModuleCallerSession) GetMaxL1Fee(dataSize *big.Int) (*big.Int, error) {
	return _IChainModule.Contract.GetMaxL1Fee(&_IChainModule.CallOpts, dataSize)
}

type GetGasOverhead struct {
	ChainModuleFixedOverhead   *big.Int
	ChainModulePerByteOverhead *big.Int
}

func (_IChainModule *IChainModule) Address() common.Address {
	return _IChainModule.address
}

type IChainModuleInterface interface {
	BlockHash(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error)

	BlockNumber(opts *bind.CallOpts) (*big.Int, error)

	GetCurrentL1Fee(opts *bind.CallOpts) (*big.Int, error)

	GetGasOverhead(opts *bind.CallOpts) (GetGasOverhead,

		error)

	GetMaxL1Fee(opts *bind.CallOpts, dataSize *big.Int) (*big.Int, error)

	Address() common.Address
}
