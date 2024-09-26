// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package optimism_module

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

var OptimismModuleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"blockHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"blockNumber\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"dataSize\",\"type\":\"uint256\"}],\"name\":\"getCurrentL1Fee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getGasOverhead\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"chainModuleFixedOverhead\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"chainModulePerByteOverhead\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"dataSize\",\"type\":\"uint256\"}],\"name\":\"getMaxL1Fee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506103dc806100206000396000f3fe608060405234801561001057600080fd5b50600436106100675760003560e01c80637810d12a116100505780637810d12a1461006c57806385df51fd14610098578063de9ee35e146100ab57600080fd5b8063125441401461006c57806357e871e714610092575b600080fd5b61007f61007a366004610221565b6100c2565b6040519081526020015b60405180910390f35b4361007f565b61007f6100a6366004610221565b6100d3565b6040805161ea60815261010e602082015201610089565b60006100cd82610100565b92915050565b600043821015806100ee57506101006100ec8343610269565b115b156100fb57506000919050565b504090565b60008061010e83600461027c565b67ffffffffffffffff81111561012657610126610293565b6040519080825280601f01601f191660200182016040528015610150576020820181803683370190505b50905073420000000000000000000000000000000000000f73ffffffffffffffffffffffffffffffffffffffff166349948e0e82604051806080016040528060508152602001610380605091396040516020016101ae9291906102e6565b6040516020818303038152906040526040518263ffffffff1660e01b81526004016101d99190610315565b602060405180830381865afa1580156101f6573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061021a9190610366565b9392505050565b60006020828403121561023357600080fd5b5035919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b818103818111156100cd576100cd61023a565b80820281158282048414176100cd576100cd61023a565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60005b838110156102dd5781810151838201526020016102c5565b50506000910152565b600083516102f88184602088016102c2565b83519083019061030c8183602088016102c2565b01949350505050565b60208152600082518060208401526103348160408501602087016102c2565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169190910160400192915050565b60006020828403121561037857600080fd5b505191905056feffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa164736f6c6343000813000a",
}

var OptimismModuleABI = OptimismModuleMetaData.ABI

var OptimismModuleBin = OptimismModuleMetaData.Bin

func DeployOptimismModule(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *OptimismModule, error) {
	parsed, err := OptimismModuleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OptimismModuleBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OptimismModule{address: address, abi: *parsed, OptimismModuleCaller: OptimismModuleCaller{contract: contract}, OptimismModuleTransactor: OptimismModuleTransactor{contract: contract}, OptimismModuleFilterer: OptimismModuleFilterer{contract: contract}}, nil
}

type OptimismModule struct {
	address common.Address
	abi     abi.ABI
	OptimismModuleCaller
	OptimismModuleTransactor
	OptimismModuleFilterer
}

type OptimismModuleCaller struct {
	contract *bind.BoundContract
}

type OptimismModuleTransactor struct {
	contract *bind.BoundContract
}

type OptimismModuleFilterer struct {
	contract *bind.BoundContract
}

type OptimismModuleSession struct {
	Contract     *OptimismModule
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type OptimismModuleCallerSession struct {
	Contract *OptimismModuleCaller
	CallOpts bind.CallOpts
}

type OptimismModuleTransactorSession struct {
	Contract     *OptimismModuleTransactor
	TransactOpts bind.TransactOpts
}

type OptimismModuleRaw struct {
	Contract *OptimismModule
}

type OptimismModuleCallerRaw struct {
	Contract *OptimismModuleCaller
}

type OptimismModuleTransactorRaw struct {
	Contract *OptimismModuleTransactor
}

func NewOptimismModule(address common.Address, backend bind.ContractBackend) (*OptimismModule, error) {
	abi, err := abi.JSON(strings.NewReader(OptimismModuleABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindOptimismModule(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OptimismModule{address: address, abi: abi, OptimismModuleCaller: OptimismModuleCaller{contract: contract}, OptimismModuleTransactor: OptimismModuleTransactor{contract: contract}, OptimismModuleFilterer: OptimismModuleFilterer{contract: contract}}, nil
}

func NewOptimismModuleCaller(address common.Address, caller bind.ContractCaller) (*OptimismModuleCaller, error) {
	contract, err := bindOptimismModule(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OptimismModuleCaller{contract: contract}, nil
}

func NewOptimismModuleTransactor(address common.Address, transactor bind.ContractTransactor) (*OptimismModuleTransactor, error) {
	contract, err := bindOptimismModule(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OptimismModuleTransactor{contract: contract}, nil
}

func NewOptimismModuleFilterer(address common.Address, filterer bind.ContractFilterer) (*OptimismModuleFilterer, error) {
	contract, err := bindOptimismModule(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OptimismModuleFilterer{contract: contract}, nil
}

func bindOptimismModule(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OptimismModuleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_OptimismModule *OptimismModuleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OptimismModule.Contract.OptimismModuleCaller.contract.Call(opts, result, method, params...)
}

func (_OptimismModule *OptimismModuleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OptimismModule.Contract.OptimismModuleTransactor.contract.Transfer(opts)
}

func (_OptimismModule *OptimismModuleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OptimismModule.Contract.OptimismModuleTransactor.contract.Transact(opts, method, params...)
}

func (_OptimismModule *OptimismModuleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OptimismModule.Contract.contract.Call(opts, result, method, params...)
}

func (_OptimismModule *OptimismModuleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OptimismModule.Contract.contract.Transfer(opts)
}

func (_OptimismModule *OptimismModuleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OptimismModule.Contract.contract.Transact(opts, method, params...)
}

func (_OptimismModule *OptimismModuleCaller) BlockHash(opts *bind.CallOpts, n *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _OptimismModule.contract.Call(opts, &out, "blockHash", n)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_OptimismModule *OptimismModuleSession) BlockHash(n *big.Int) ([32]byte, error) {
	return _OptimismModule.Contract.BlockHash(&_OptimismModule.CallOpts, n)
}

func (_OptimismModule *OptimismModuleCallerSession) BlockHash(n *big.Int) ([32]byte, error) {
	return _OptimismModule.Contract.BlockHash(&_OptimismModule.CallOpts, n)
}

func (_OptimismModule *OptimismModuleCaller) BlockNumber(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OptimismModule.contract.Call(opts, &out, "blockNumber")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_OptimismModule *OptimismModuleSession) BlockNumber() (*big.Int, error) {
	return _OptimismModule.Contract.BlockNumber(&_OptimismModule.CallOpts)
}

func (_OptimismModule *OptimismModuleCallerSession) BlockNumber() (*big.Int, error) {
	return _OptimismModule.Contract.BlockNumber(&_OptimismModule.CallOpts)
}

func (_OptimismModule *OptimismModuleCaller) GetCurrentL1Fee(opts *bind.CallOpts, dataSize *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _OptimismModule.contract.Call(opts, &out, "getCurrentL1Fee", dataSize)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_OptimismModule *OptimismModuleSession) GetCurrentL1Fee(dataSize *big.Int) (*big.Int, error) {
	return _OptimismModule.Contract.GetCurrentL1Fee(&_OptimismModule.CallOpts, dataSize)
}

func (_OptimismModule *OptimismModuleCallerSession) GetCurrentL1Fee(dataSize *big.Int) (*big.Int, error) {
	return _OptimismModule.Contract.GetCurrentL1Fee(&_OptimismModule.CallOpts, dataSize)
}

func (_OptimismModule *OptimismModuleCaller) GetGasOverhead(opts *bind.CallOpts) (GetGasOverhead,

	error) {
	var out []interface{}
	err := _OptimismModule.contract.Call(opts, &out, "getGasOverhead")

	outstruct := new(GetGasOverhead)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ChainModuleFixedOverhead = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.ChainModulePerByteOverhead = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_OptimismModule *OptimismModuleSession) GetGasOverhead() (GetGasOverhead,

	error) {
	return _OptimismModule.Contract.GetGasOverhead(&_OptimismModule.CallOpts)
}

func (_OptimismModule *OptimismModuleCallerSession) GetGasOverhead() (GetGasOverhead,

	error) {
	return _OptimismModule.Contract.GetGasOverhead(&_OptimismModule.CallOpts)
}

func (_OptimismModule *OptimismModuleCaller) GetMaxL1Fee(opts *bind.CallOpts, dataSize *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _OptimismModule.contract.Call(opts, &out, "getMaxL1Fee", dataSize)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_OptimismModule *OptimismModuleSession) GetMaxL1Fee(dataSize *big.Int) (*big.Int, error) {
	return _OptimismModule.Contract.GetMaxL1Fee(&_OptimismModule.CallOpts, dataSize)
}

func (_OptimismModule *OptimismModuleCallerSession) GetMaxL1Fee(dataSize *big.Int) (*big.Int, error) {
	return _OptimismModule.Contract.GetMaxL1Fee(&_OptimismModule.CallOpts, dataSize)
}

type GetGasOverhead struct {
	ChainModuleFixedOverhead   *big.Int
	ChainModulePerByteOverhead *big.Int
}

func (_OptimismModule *OptimismModule) Address() common.Address {
	return _OptimismModule.address
}

type OptimismModuleInterface interface {
	BlockHash(opts *bind.CallOpts, n *big.Int) ([32]byte, error)

	BlockNumber(opts *bind.CallOpts) (*big.Int, error)

	GetCurrentL1Fee(opts *bind.CallOpts, dataSize *big.Int) (*big.Int, error)

	GetGasOverhead(opts *bind.CallOpts) (GetGasOverhead,

		error)

	GetMaxL1Fee(opts *bind.CallOpts, dataSize *big.Int) (*big.Int, error)

	Address() common.Address
}
