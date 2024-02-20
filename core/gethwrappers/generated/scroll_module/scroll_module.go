// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package scroll_module

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

var ScrollModuleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"blockHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"blockNumber\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCurrentL1Fee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getGasOverhead\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"chainModuleFixedOverhead\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"chainModulePerByteOverhead\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"dataSize\",\"type\":\"uint256\"}],\"name\":\"getMaxL1Fee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506104f6806100206000396000f3fe608060405234801561001057600080fd5b50600436106100675760003560e01c806357e871e71161005057806357e871e71461009a57806385df51fd146100a0578063de9ee35e146100b357600080fd5b8063125441401461006c57806318b8f61314610092575b600080fd5b61007f61007a3660046102e6565b6100c7565b6040519081526020015b60405180910390f35b61007f6101e8565b4361007f565b61007f6100ae3660046102e6565b6102b9565b604080516000808252602082015201610089565b6000806100d583600461032e565b67ffffffffffffffff8111156100ed576100ed61034b565b6040519080825280601f01601f191660200182016040528015610117576020820181803683370190505b50905073530000000000000000000000000000000000000273ffffffffffffffffffffffffffffffffffffffff166349948e0e826040518060a00160405280607881526020016104726078913960405160200161017592919061039e565b6040516020818303038152906040526040518263ffffffff1660e01b81526004016101a091906103cd565b602060405180830381865afa1580156101bd573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906101e1919061041e565b9392505050565b600073530000000000000000000000000000000000000273ffffffffffffffffffffffffffffffffffffffff166349948e0e6000366040518060a00160405280607881526020016104726078913960405160200161024893929190610437565b6040516020818303038152906040526040518263ffffffff1660e01b815260040161027391906103cd565b602060405180830381865afa158015610290573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102b4919061041e565b905090565b600043821015806102d457506101006102d2834361045e565b115b156102e157506000919050565b504090565b6000602082840312156102f857600080fd5b5035919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8082028115828204841417610345576103456102ff565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60005b8381101561039557818101518382015260200161037d565b50506000910152565b600083516103b081846020880161037a565b8351908301906103c481836020880161037a565b01949350505050565b60208152600082518060208401526103ec81604085016020870161037a565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169190910160400192915050565b60006020828403121561043057600080fd5b5051919050565b82848237600083820160008152835161045481836020880161037a565b0195945050505050565b81810381811115610345576103456102ff56feffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa164736f6c6343000813000a",
}

var ScrollModuleABI = ScrollModuleMetaData.ABI

var ScrollModuleBin = ScrollModuleMetaData.Bin

func DeployScrollModule(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ScrollModule, error) {
	parsed, err := ScrollModuleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ScrollModuleBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ScrollModule{address: address, abi: *parsed, ScrollModuleCaller: ScrollModuleCaller{contract: contract}, ScrollModuleTransactor: ScrollModuleTransactor{contract: contract}, ScrollModuleFilterer: ScrollModuleFilterer{contract: contract}}, nil
}

type ScrollModule struct {
	address common.Address
	abi     abi.ABI
	ScrollModuleCaller
	ScrollModuleTransactor
	ScrollModuleFilterer
}

type ScrollModuleCaller struct {
	contract *bind.BoundContract
}

type ScrollModuleTransactor struct {
	contract *bind.BoundContract
}

type ScrollModuleFilterer struct {
	contract *bind.BoundContract
}

type ScrollModuleSession struct {
	Contract     *ScrollModule
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type ScrollModuleCallerSession struct {
	Contract *ScrollModuleCaller
	CallOpts bind.CallOpts
}

type ScrollModuleTransactorSession struct {
	Contract     *ScrollModuleTransactor
	TransactOpts bind.TransactOpts
}

type ScrollModuleRaw struct {
	Contract *ScrollModule
}

type ScrollModuleCallerRaw struct {
	Contract *ScrollModuleCaller
}

type ScrollModuleTransactorRaw struct {
	Contract *ScrollModuleTransactor
}

func NewScrollModule(address common.Address, backend bind.ContractBackend) (*ScrollModule, error) {
	abi, err := abi.JSON(strings.NewReader(ScrollModuleABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindScrollModule(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ScrollModule{address: address, abi: abi, ScrollModuleCaller: ScrollModuleCaller{contract: contract}, ScrollModuleTransactor: ScrollModuleTransactor{contract: contract}, ScrollModuleFilterer: ScrollModuleFilterer{contract: contract}}, nil
}

func NewScrollModuleCaller(address common.Address, caller bind.ContractCaller) (*ScrollModuleCaller, error) {
	contract, err := bindScrollModule(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ScrollModuleCaller{contract: contract}, nil
}

func NewScrollModuleTransactor(address common.Address, transactor bind.ContractTransactor) (*ScrollModuleTransactor, error) {
	contract, err := bindScrollModule(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ScrollModuleTransactor{contract: contract}, nil
}

func NewScrollModuleFilterer(address common.Address, filterer bind.ContractFilterer) (*ScrollModuleFilterer, error) {
	contract, err := bindScrollModule(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ScrollModuleFilterer{contract: contract}, nil
}

func bindScrollModule(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ScrollModuleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_ScrollModule *ScrollModuleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ScrollModule.Contract.ScrollModuleCaller.contract.Call(opts, result, method, params...)
}

func (_ScrollModule *ScrollModuleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ScrollModule.Contract.ScrollModuleTransactor.contract.Transfer(opts)
}

func (_ScrollModule *ScrollModuleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ScrollModule.Contract.ScrollModuleTransactor.contract.Transact(opts, method, params...)
}

func (_ScrollModule *ScrollModuleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ScrollModule.Contract.contract.Call(opts, result, method, params...)
}

func (_ScrollModule *ScrollModuleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ScrollModule.Contract.contract.Transfer(opts)
}

func (_ScrollModule *ScrollModuleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ScrollModule.Contract.contract.Transact(opts, method, params...)
}

func (_ScrollModule *ScrollModuleCaller) BlockHash(opts *bind.CallOpts, n *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _ScrollModule.contract.Call(opts, &out, "blockHash", n)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_ScrollModule *ScrollModuleSession) BlockHash(n *big.Int) ([32]byte, error) {
	return _ScrollModule.Contract.BlockHash(&_ScrollModule.CallOpts, n)
}

func (_ScrollModule *ScrollModuleCallerSession) BlockHash(n *big.Int) ([32]byte, error) {
	return _ScrollModule.Contract.BlockHash(&_ScrollModule.CallOpts, n)
}

func (_ScrollModule *ScrollModuleCaller) BlockNumber(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ScrollModule.contract.Call(opts, &out, "blockNumber")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_ScrollModule *ScrollModuleSession) BlockNumber() (*big.Int, error) {
	return _ScrollModule.Contract.BlockNumber(&_ScrollModule.CallOpts)
}

func (_ScrollModule *ScrollModuleCallerSession) BlockNumber() (*big.Int, error) {
	return _ScrollModule.Contract.BlockNumber(&_ScrollModule.CallOpts)
}

func (_ScrollModule *ScrollModuleCaller) GetCurrentL1Fee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ScrollModule.contract.Call(opts, &out, "getCurrentL1Fee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_ScrollModule *ScrollModuleSession) GetCurrentL1Fee() (*big.Int, error) {
	return _ScrollModule.Contract.GetCurrentL1Fee(&_ScrollModule.CallOpts)
}

func (_ScrollModule *ScrollModuleCallerSession) GetCurrentL1Fee() (*big.Int, error) {
	return _ScrollModule.Contract.GetCurrentL1Fee(&_ScrollModule.CallOpts)
}

func (_ScrollModule *ScrollModuleCaller) GetGasOverhead(opts *bind.CallOpts) (GetGasOverhead,

	error) {
	var out []interface{}
	err := _ScrollModule.contract.Call(opts, &out, "getGasOverhead")

	outstruct := new(GetGasOverhead)
	if err != nil {
		return *outstruct, err
	}

	outstruct.ChainModuleFixedOverhead = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.ChainModulePerByteOverhead = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

func (_ScrollModule *ScrollModuleSession) GetGasOverhead() (GetGasOverhead,

	error) {
	return _ScrollModule.Contract.GetGasOverhead(&_ScrollModule.CallOpts)
}

func (_ScrollModule *ScrollModuleCallerSession) GetGasOverhead() (GetGasOverhead,

	error) {
	return _ScrollModule.Contract.GetGasOverhead(&_ScrollModule.CallOpts)
}

func (_ScrollModule *ScrollModuleCaller) GetMaxL1Fee(opts *bind.CallOpts, dataSize *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _ScrollModule.contract.Call(opts, &out, "getMaxL1Fee", dataSize)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_ScrollModule *ScrollModuleSession) GetMaxL1Fee(dataSize *big.Int) (*big.Int, error) {
	return _ScrollModule.Contract.GetMaxL1Fee(&_ScrollModule.CallOpts, dataSize)
}

func (_ScrollModule *ScrollModuleCallerSession) GetMaxL1Fee(dataSize *big.Int) (*big.Int, error) {
	return _ScrollModule.Contract.GetMaxL1Fee(&_ScrollModule.CallOpts, dataSize)
}

type GetGasOverhead struct {
	ChainModuleFixedOverhead   *big.Int
	ChainModulePerByteOverhead *big.Int
}

func (_ScrollModule *ScrollModule) Address() common.Address {
	return _ScrollModule.address
}

type ScrollModuleInterface interface {
	BlockHash(opts *bind.CallOpts, n *big.Int) ([32]byte, error)

	BlockNumber(opts *bind.CallOpts) (*big.Int, error)

	GetCurrentL1Fee(opts *bind.CallOpts) (*big.Int, error)

	GetGasOverhead(opts *bind.CallOpts) (GetGasOverhead,

		error)

	GetMaxL1Fee(opts *bind.CallOpts, dataSize *big.Int) (*big.Int, error)

	Address() common.Address
}
