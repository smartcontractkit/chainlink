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
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"blockHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"blockNumber\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCurrentL1Fee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getGasOverhead\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"chainModuleFixedOverhead\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"chainModulePerByteOverhead\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"dataSize\",\"type\":\"uint256\"}],\"name\":\"getMaxL1Fee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506104d1806100206000396000f3fe608060405234801561001057600080fd5b50600436106100675760003560e01c806357e871e71161005057806357e871e71461009a57806385df51fd146100a0578063de9ee35e146100b357600080fd5b8063125441401461006c57806318b8f61314610092575b600080fd5b61007f61007a3660046102e9565b6100ca565b6040519081526020015b60405180910390f35b61007f6101eb565b4361007f565b61007f6100ae3660046102e9565b6102bc565b6040805161ea60815261010e602082015201610089565b6000806100d8836004610331565b67ffffffffffffffff8111156100f0576100f061034e565b6040519080825280601f01601f19166020018201604052801561011a576020820181803683370190505b50905073420000000000000000000000000000000000000f73ffffffffffffffffffffffffffffffffffffffff166349948e0e82604051806080016040528060508152602001610475605091396040516020016101789291906103a1565b6040516020818303038152906040526040518263ffffffff1660e01b81526004016101a391906103d0565b602060405180830381865afa1580156101c0573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906101e49190610421565b9392505050565b600073420000000000000000000000000000000000000f73ffffffffffffffffffffffffffffffffffffffff166349948e0e6000366040518060800160405280605081526020016104756050913960405160200161024b9392919061043a565b6040516020818303038152906040526040518263ffffffff1660e01b815260040161027691906103d0565b602060405180830381865afa158015610293573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102b79190610421565b905090565b600043821015806102d757506101006102d58343610461565b115b156102e457506000919050565b504090565b6000602082840312156102fb57600080fd5b5035919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b808202811582820484141761034857610348610302565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60005b83811015610398578181015183820152602001610380565b50506000910152565b600083516103b381846020880161037d565b8351908301906103c781836020880161037d565b01949350505050565b60208152600082518060208401526103ef81604085016020870161037d565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169190910160400192915050565b60006020828403121561043357600080fd5b5051919050565b82848237600083820160008152835161045781836020880161037d565b0195945050505050565b818103818111156103485761034861030256feffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa164736f6c6343000813000a",
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

func (_OptimismModule *OptimismModuleCaller) GetCurrentL1Fee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OptimismModule.contract.Call(opts, &out, "getCurrentL1Fee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_OptimismModule *OptimismModuleSession) GetCurrentL1Fee() (*big.Int, error) {
	return _OptimismModule.Contract.GetCurrentL1Fee(&_OptimismModule.CallOpts)
}

func (_OptimismModule *OptimismModuleCallerSession) GetCurrentL1Fee() (*big.Int, error) {
	return _OptimismModule.Contract.GetCurrentL1Fee(&_OptimismModule.CallOpts)
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

	GetCurrentL1Fee(opts *bind.CallOpts) (*big.Int, error)

	GetGasOverhead(opts *bind.CallOpts) (GetGasOverhead,

		error)

	GetMaxL1Fee(opts *bind.CallOpts, dataSize *big.Int) (*big.Int, error)

	Address() common.Address
}
