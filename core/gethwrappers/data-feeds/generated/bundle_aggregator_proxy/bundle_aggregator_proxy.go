// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bundle_aggregator_proxy

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated"
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

var BundleAggregatorProxyMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"aggregatorAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"aggregator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"bundleDecimals\",\"inputs\":[],\"outputs\":[{\"name\":\"decimals\",\"type\":\"uint8[]\",\"internalType\":\"uint8[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"confirmAggregator\",\"inputs\":[{\"name\":\"aggregatorAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"description\",\"inputs\":[],\"outputs\":[{\"name\":\"aggregatorDescription\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestBundle\",\"inputs\":[],\"outputs\":[{\"name\":\"bundle\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"latestBundleTimestamp\",\"inputs\":[],\"outputs\":[{\"name\":\"timestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proposeAggregator\",\"inputs\":[{\"name\":\"aggregatorAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"proposedAggregator\",\"inputs\":[],\"outputs\":[{\"name\":\"proposedAggregatorAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"typeAndVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"version\",\"inputs\":[],\"outputs\":[{\"name\":\"aggregatorVersion\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"AggregatorConfirmed\",\"inputs\":[{\"name\":\"previous\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"latest\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AggregatorProposed\",\"inputs\":[{\"name\":\"current\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"proposed\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AggregatorNotProposed\",\"inputs\":[{\"name\":\"aggregator\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
	Bin: "0x608060405234801561001057600080fd5b50604051610e39380380610e3983398101604081905261002f916101ae565b808060006001600160a01b03821661008e5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b03848116919091179091558116156100be576100be816100e9565b5050600280546001600160a01b0319166001600160a01b039490941693909317909255506101e19050565b336001600160a01b038216036101415760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610085565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b03811681146101a957600080fd5b919050565b600080604083850312156101c157600080fd5b6101ca83610192565b91506101d860208401610192565b90509250929050565b610c49806101f06000396000f3fe608060405234801561001057600080fd5b50600436106100df5760003560e01c80639198274f1161008c578063a928c09611610066578063a928c096146101e0578063e8c4be30146101f3578063f2fde38b14610211578063f8a2abd31461022457600080fd5b80639198274f146101bb5780639d91348d146101c3578063a3d610cc146101d857600080fd5b80637284e416116100bd5780637284e4161461018b57806379ba5097146101935780638da5cb5b1461019d57600080fd5b8063181f5a77146100e4578063245a7bfc1461013657806354fd4d5014610175575b600080fd5b6101206040518060400160405280601b81526020017f42756e646c6541676772656761746f7250726f787920312e302e30000000000081525081565b60405161012d919061098e565b60405180910390f35b60025473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161012d565b61017d610237565b60405190815260200161012d565b6101206102d0565b61019b610386565b005b60005473ffffffffffffffffffffffffffffffffffffffff16610150565b610120610488565b6101cb6104f8565b60405161012d91906109a8565b61017d6105ae565b61019b6101ee3660046109ee565b61061e565b60035473ffffffffffffffffffffffffffffffffffffffff16610150565b61019b61021f3660046109ee565b610715565b61019b6102323660046109ee565b610729565b600254604080517f54fd4d50000000000000000000000000000000000000000000000000000000008152905160009273ffffffffffffffffffffffffffffffffffffffff16916354fd4d509160048083019260209291908290030181865afa1580156102a7573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102cb9190610a24565b905090565b600254604080517f7284e416000000000000000000000000000000000000000000000000000000008152905160609273ffffffffffffffffffffffffffffffffffffffff1691637284e4169160048083019260009291908290030181865afa158015610340573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526102cb9190810190610b2c565b60015473ffffffffffffffffffffffffffffffffffffffff16331461040c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b600254604080517f9198274f000000000000000000000000000000000000000000000000000000008152905160609273ffffffffffffffffffffffffffffffffffffffff1691639198274f9160048083019260009291908290030181865afa158015610340573d6000803e3d6000fd5b600254604080517f9d91348d000000000000000000000000000000000000000000000000000000008152905160609273ffffffffffffffffffffffffffffffffffffffff1691639d91348d9160048083019260009291908290030181865afa158015610568573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526102cb9190810190610b7d565b600254604080517fa3d610cc000000000000000000000000000000000000000000000000000000008152905160009273ffffffffffffffffffffffffffffffffffffffff169163a3d610cc9160048083019260209291908290030181865afa1580156102a7573d6000803e3d6000fd5b6106266107a8565b60035473ffffffffffffffffffffffffffffffffffffffff828116911614610692576040517feb61b92800000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff82166004820152602401610403565b60028054600380547fffffffffffffffffffffffff000000000000000000000000000000000000000090811690915573ffffffffffffffffffffffffffffffffffffffff8481169183168217909355604051929091169182907f33745f67a407dcb785417f9c123dd3641479a102674b6e35c1f10975625b90e990600090a35050565b61071d6107a8565b6107268161082b565b50565b6107316107a8565b600380547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff838116918217909255600254604051919216907fc0f151710f03d713b71d9970cee0d5b11ddc9a7552abaa3f6ee818010f21600d90600090a350565b60005473ffffffffffffffffffffffffffffffffffffffff163314610829576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610403565b565b3373ffffffffffffffffffffffffffffffffffffffff8216036108aa576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610403565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b60005b8381101561093b578181015183820152602001610923565b50506000910152565b6000815180845261095c816020860160208601610920565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b6020815260006109a16020830184610944565b9392505050565b602080825282518282018190526000918401906040840190835b818110156109e357835160ff168352602093840193909201916001016109c2565b509095945050505050565b600060208284031215610a0057600080fd5b813573ffffffffffffffffffffffffffffffffffffffff811681146109a157600080fd5b600060208284031215610a3657600080fd5b5051919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610ab357610ab3610a3d565b604052919050565b60008067ffffffffffffffff841115610ad657610ad6610a3d565b50601f83017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016602001610b0981610a6c565b915050828152838383011115610b1e57600080fd5b6109a1836020830184610920565b600060208284031215610b3e57600080fd5b815167ffffffffffffffff811115610b5557600080fd5b8201601f81018413610b6657600080fd5b610b7584825160208401610abb565b949350505050565b600060208284031215610b8f57600080fd5b815167ffffffffffffffff811115610ba657600080fd5b8201601f81018413610bb757600080fd5b805167ffffffffffffffff811115610bd157610bd1610a3d565b8060051b610be160208201610a6c565b91825260208184018101929081019087841115610bfd57600080fd5b6020850194505b83851015610c31578451925060ff83168314610c1f57600080fd5b82825260209485019490910190610c04565b97965050505050505056fea164736f6c634300081a000a",
}

var BundleAggregatorProxyABI = BundleAggregatorProxyMetaData.ABI

var BundleAggregatorProxyBin = BundleAggregatorProxyMetaData.Bin

func DeployBundleAggregatorProxy(auth *bind.TransactOpts, backend bind.ContractBackend, aggregatorAddress common.Address, owner common.Address) (common.Address, *types.Transaction, *BundleAggregatorProxy, error) {
	parsed, err := BundleAggregatorProxyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BundleAggregatorProxyBin), backend, aggregatorAddress, owner)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BundleAggregatorProxy{address: address, abi: *parsed, BundleAggregatorProxyCaller: BundleAggregatorProxyCaller{contract: contract}, BundleAggregatorProxyTransactor: BundleAggregatorProxyTransactor{contract: contract}, BundleAggregatorProxyFilterer: BundleAggregatorProxyFilterer{contract: contract}}, nil
}

type BundleAggregatorProxy struct {
	address common.Address
	abi     abi.ABI
	BundleAggregatorProxyCaller
	BundleAggregatorProxyTransactor
	BundleAggregatorProxyFilterer
}

type BundleAggregatorProxyCaller struct {
	contract *bind.BoundContract
}

type BundleAggregatorProxyTransactor struct {
	contract *bind.BoundContract
}

type BundleAggregatorProxyFilterer struct {
	contract *bind.BoundContract
}

type BundleAggregatorProxySession struct {
	Contract     *BundleAggregatorProxy
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type BundleAggregatorProxyCallerSession struct {
	Contract *BundleAggregatorProxyCaller
	CallOpts bind.CallOpts
}

type BundleAggregatorProxyTransactorSession struct {
	Contract     *BundleAggregatorProxyTransactor
	TransactOpts bind.TransactOpts
}

type BundleAggregatorProxyRaw struct {
	Contract *BundleAggregatorProxy
}

type BundleAggregatorProxyCallerRaw struct {
	Contract *BundleAggregatorProxyCaller
}

type BundleAggregatorProxyTransactorRaw struct {
	Contract *BundleAggregatorProxyTransactor
}

func NewBundleAggregatorProxy(address common.Address, backend bind.ContractBackend) (*BundleAggregatorProxy, error) {
	abi, err := abi.JSON(strings.NewReader(BundleAggregatorProxyABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindBundleAggregatorProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BundleAggregatorProxy{address: address, abi: abi, BundleAggregatorProxyCaller: BundleAggregatorProxyCaller{contract: contract}, BundleAggregatorProxyTransactor: BundleAggregatorProxyTransactor{contract: contract}, BundleAggregatorProxyFilterer: BundleAggregatorProxyFilterer{contract: contract}}, nil
}

func NewBundleAggregatorProxyCaller(address common.Address, caller bind.ContractCaller) (*BundleAggregatorProxyCaller, error) {
	contract, err := bindBundleAggregatorProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BundleAggregatorProxyCaller{contract: contract}, nil
}

func NewBundleAggregatorProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*BundleAggregatorProxyTransactor, error) {
	contract, err := bindBundleAggregatorProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BundleAggregatorProxyTransactor{contract: contract}, nil
}

func NewBundleAggregatorProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*BundleAggregatorProxyFilterer, error) {
	contract, err := bindBundleAggregatorProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BundleAggregatorProxyFilterer{contract: contract}, nil
}

func bindBundleAggregatorProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BundleAggregatorProxyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_BundleAggregatorProxy *BundleAggregatorProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BundleAggregatorProxy.Contract.BundleAggregatorProxyCaller.contract.Call(opts, result, method, params...)
}

func (_BundleAggregatorProxy *BundleAggregatorProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BundleAggregatorProxy.Contract.BundleAggregatorProxyTransactor.contract.Transfer(opts)
}

func (_BundleAggregatorProxy *BundleAggregatorProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BundleAggregatorProxy.Contract.BundleAggregatorProxyTransactor.contract.Transact(opts, method, params...)
}

func (_BundleAggregatorProxy *BundleAggregatorProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BundleAggregatorProxy.Contract.contract.Call(opts, result, method, params...)
}

func (_BundleAggregatorProxy *BundleAggregatorProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BundleAggregatorProxy.Contract.contract.Transfer(opts)
}

func (_BundleAggregatorProxy *BundleAggregatorProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BundleAggregatorProxy.Contract.contract.Transact(opts, method, params...)
}

func (_BundleAggregatorProxy *BundleAggregatorProxyCaller) Aggregator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BundleAggregatorProxy.contract.Call(opts, &out, "aggregator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BundleAggregatorProxy *BundleAggregatorProxySession) Aggregator() (common.Address, error) {
	return _BundleAggregatorProxy.Contract.Aggregator(&_BundleAggregatorProxy.CallOpts)
}

func (_BundleAggregatorProxy *BundleAggregatorProxyCallerSession) Aggregator() (common.Address, error) {
	return _BundleAggregatorProxy.Contract.Aggregator(&_BundleAggregatorProxy.CallOpts)
}

func (_BundleAggregatorProxy *BundleAggregatorProxyCaller) BundleDecimals(opts *bind.CallOpts) ([]uint8, error) {
	var out []interface{}
	err := _BundleAggregatorProxy.contract.Call(opts, &out, "bundleDecimals")

	if err != nil {
		return *new([]uint8), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint8)).(*[]uint8)

	return out0, err

}

func (_BundleAggregatorProxy *BundleAggregatorProxySession) BundleDecimals() ([]uint8, error) {
	return _BundleAggregatorProxy.Contract.BundleDecimals(&_BundleAggregatorProxy.CallOpts)
}

func (_BundleAggregatorProxy *BundleAggregatorProxyCallerSession) BundleDecimals() ([]uint8, error) {
	return _BundleAggregatorProxy.Contract.BundleDecimals(&_BundleAggregatorProxy.CallOpts)
}

func (_BundleAggregatorProxy *BundleAggregatorProxyCaller) Description(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _BundleAggregatorProxy.contract.Call(opts, &out, "description")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_BundleAggregatorProxy *BundleAggregatorProxySession) Description() (string, error) {
	return _BundleAggregatorProxy.Contract.Description(&_BundleAggregatorProxy.CallOpts)
}

func (_BundleAggregatorProxy *BundleAggregatorProxyCallerSession) Description() (string, error) {
	return _BundleAggregatorProxy.Contract.Description(&_BundleAggregatorProxy.CallOpts)
}

func (_BundleAggregatorProxy *BundleAggregatorProxyCaller) LatestBundle(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _BundleAggregatorProxy.contract.Call(opts, &out, "latestBundle")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_BundleAggregatorProxy *BundleAggregatorProxySession) LatestBundle() ([]byte, error) {
	return _BundleAggregatorProxy.Contract.LatestBundle(&_BundleAggregatorProxy.CallOpts)
}

func (_BundleAggregatorProxy *BundleAggregatorProxyCallerSession) LatestBundle() ([]byte, error) {
	return _BundleAggregatorProxy.Contract.LatestBundle(&_BundleAggregatorProxy.CallOpts)
}

func (_BundleAggregatorProxy *BundleAggregatorProxyCaller) LatestBundleTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BundleAggregatorProxy.contract.Call(opts, &out, "latestBundleTimestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_BundleAggregatorProxy *BundleAggregatorProxySession) LatestBundleTimestamp() (*big.Int, error) {
	return _BundleAggregatorProxy.Contract.LatestBundleTimestamp(&_BundleAggregatorProxy.CallOpts)
}

func (_BundleAggregatorProxy *BundleAggregatorProxyCallerSession) LatestBundleTimestamp() (*big.Int, error) {
	return _BundleAggregatorProxy.Contract.LatestBundleTimestamp(&_BundleAggregatorProxy.CallOpts)
}

func (_BundleAggregatorProxy *BundleAggregatorProxyCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BundleAggregatorProxy.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BundleAggregatorProxy *BundleAggregatorProxySession) Owner() (common.Address, error) {
	return _BundleAggregatorProxy.Contract.Owner(&_BundleAggregatorProxy.CallOpts)
}

func (_BundleAggregatorProxy *BundleAggregatorProxyCallerSession) Owner() (common.Address, error) {
	return _BundleAggregatorProxy.Contract.Owner(&_BundleAggregatorProxy.CallOpts)
}

func (_BundleAggregatorProxy *BundleAggregatorProxyCaller) ProposedAggregator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BundleAggregatorProxy.contract.Call(opts, &out, "proposedAggregator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_BundleAggregatorProxy *BundleAggregatorProxySession) ProposedAggregator() (common.Address, error) {
	return _BundleAggregatorProxy.Contract.ProposedAggregator(&_BundleAggregatorProxy.CallOpts)
}

func (_BundleAggregatorProxy *BundleAggregatorProxyCallerSession) ProposedAggregator() (common.Address, error) {
	return _BundleAggregatorProxy.Contract.ProposedAggregator(&_BundleAggregatorProxy.CallOpts)
}

func (_BundleAggregatorProxy *BundleAggregatorProxyCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _BundleAggregatorProxy.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_BundleAggregatorProxy *BundleAggregatorProxySession) TypeAndVersion() (string, error) {
	return _BundleAggregatorProxy.Contract.TypeAndVersion(&_BundleAggregatorProxy.CallOpts)
}

func (_BundleAggregatorProxy *BundleAggregatorProxyCallerSession) TypeAndVersion() (string, error) {
	return _BundleAggregatorProxy.Contract.TypeAndVersion(&_BundleAggregatorProxy.CallOpts)
}

func (_BundleAggregatorProxy *BundleAggregatorProxyCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BundleAggregatorProxy.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_BundleAggregatorProxy *BundleAggregatorProxySession) Version() (*big.Int, error) {
	return _BundleAggregatorProxy.Contract.Version(&_BundleAggregatorProxy.CallOpts)
}

func (_BundleAggregatorProxy *BundleAggregatorProxyCallerSession) Version() (*big.Int, error) {
	return _BundleAggregatorProxy.Contract.Version(&_BundleAggregatorProxy.CallOpts)
}

func (_BundleAggregatorProxy *BundleAggregatorProxyTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BundleAggregatorProxy.contract.Transact(opts, "acceptOwnership")
}

func (_BundleAggregatorProxy *BundleAggregatorProxySession) AcceptOwnership() (*types.Transaction, error) {
	return _BundleAggregatorProxy.Contract.AcceptOwnership(&_BundleAggregatorProxy.TransactOpts)
}

func (_BundleAggregatorProxy *BundleAggregatorProxyTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _BundleAggregatorProxy.Contract.AcceptOwnership(&_BundleAggregatorProxy.TransactOpts)
}

func (_BundleAggregatorProxy *BundleAggregatorProxyTransactor) ConfirmAggregator(opts *bind.TransactOpts, aggregatorAddress common.Address) (*types.Transaction, error) {
	return _BundleAggregatorProxy.contract.Transact(opts, "confirmAggregator", aggregatorAddress)
}

func (_BundleAggregatorProxy *BundleAggregatorProxySession) ConfirmAggregator(aggregatorAddress common.Address) (*types.Transaction, error) {
	return _BundleAggregatorProxy.Contract.ConfirmAggregator(&_BundleAggregatorProxy.TransactOpts, aggregatorAddress)
}

func (_BundleAggregatorProxy *BundleAggregatorProxyTransactorSession) ConfirmAggregator(aggregatorAddress common.Address) (*types.Transaction, error) {
	return _BundleAggregatorProxy.Contract.ConfirmAggregator(&_BundleAggregatorProxy.TransactOpts, aggregatorAddress)
}

func (_BundleAggregatorProxy *BundleAggregatorProxyTransactor) ProposeAggregator(opts *bind.TransactOpts, aggregatorAddress common.Address) (*types.Transaction, error) {
	return _BundleAggregatorProxy.contract.Transact(opts, "proposeAggregator", aggregatorAddress)
}

func (_BundleAggregatorProxy *BundleAggregatorProxySession) ProposeAggregator(aggregatorAddress common.Address) (*types.Transaction, error) {
	return _BundleAggregatorProxy.Contract.ProposeAggregator(&_BundleAggregatorProxy.TransactOpts, aggregatorAddress)
}

func (_BundleAggregatorProxy *BundleAggregatorProxyTransactorSession) ProposeAggregator(aggregatorAddress common.Address) (*types.Transaction, error) {
	return _BundleAggregatorProxy.Contract.ProposeAggregator(&_BundleAggregatorProxy.TransactOpts, aggregatorAddress)
}

func (_BundleAggregatorProxy *BundleAggregatorProxyTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _BundleAggregatorProxy.contract.Transact(opts, "transferOwnership", to)
}

func (_BundleAggregatorProxy *BundleAggregatorProxySession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _BundleAggregatorProxy.Contract.TransferOwnership(&_BundleAggregatorProxy.TransactOpts, to)
}

func (_BundleAggregatorProxy *BundleAggregatorProxyTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _BundleAggregatorProxy.Contract.TransferOwnership(&_BundleAggregatorProxy.TransactOpts, to)
}

type BundleAggregatorProxyAggregatorConfirmedIterator struct {
	Event *BundleAggregatorProxyAggregatorConfirmed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BundleAggregatorProxyAggregatorConfirmedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BundleAggregatorProxyAggregatorConfirmed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(BundleAggregatorProxyAggregatorConfirmed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *BundleAggregatorProxyAggregatorConfirmedIterator) Error() error {
	return it.fail
}

func (it *BundleAggregatorProxyAggregatorConfirmedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BundleAggregatorProxyAggregatorConfirmed struct {
	Previous common.Address
	Latest   common.Address
	Raw      types.Log
}

func (_BundleAggregatorProxy *BundleAggregatorProxyFilterer) FilterAggregatorConfirmed(opts *bind.FilterOpts, previous []common.Address, latest []common.Address) (*BundleAggregatorProxyAggregatorConfirmedIterator, error) {

	var previousRule []interface{}
	for _, previousItem := range previous {
		previousRule = append(previousRule, previousItem)
	}
	var latestRule []interface{}
	for _, latestItem := range latest {
		latestRule = append(latestRule, latestItem)
	}

	logs, sub, err := _BundleAggregatorProxy.contract.FilterLogs(opts, "AggregatorConfirmed", previousRule, latestRule)
	if err != nil {
		return nil, err
	}
	return &BundleAggregatorProxyAggregatorConfirmedIterator{contract: _BundleAggregatorProxy.contract, event: "AggregatorConfirmed", logs: logs, sub: sub}, nil
}

func (_BundleAggregatorProxy *BundleAggregatorProxyFilterer) WatchAggregatorConfirmed(opts *bind.WatchOpts, sink chan<- *BundleAggregatorProxyAggregatorConfirmed, previous []common.Address, latest []common.Address) (event.Subscription, error) {

	var previousRule []interface{}
	for _, previousItem := range previous {
		previousRule = append(previousRule, previousItem)
	}
	var latestRule []interface{}
	for _, latestItem := range latest {
		latestRule = append(latestRule, latestItem)
	}

	logs, sub, err := _BundleAggregatorProxy.contract.WatchLogs(opts, "AggregatorConfirmed", previousRule, latestRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BundleAggregatorProxyAggregatorConfirmed)
				if err := _BundleAggregatorProxy.contract.UnpackLog(event, "AggregatorConfirmed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_BundleAggregatorProxy *BundleAggregatorProxyFilterer) ParseAggregatorConfirmed(log types.Log) (*BundleAggregatorProxyAggregatorConfirmed, error) {
	event := new(BundleAggregatorProxyAggregatorConfirmed)
	if err := _BundleAggregatorProxy.contract.UnpackLog(event, "AggregatorConfirmed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BundleAggregatorProxyAggregatorProposedIterator struct {
	Event *BundleAggregatorProxyAggregatorProposed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BundleAggregatorProxyAggregatorProposedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BundleAggregatorProxyAggregatorProposed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(BundleAggregatorProxyAggregatorProposed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *BundleAggregatorProxyAggregatorProposedIterator) Error() error {
	return it.fail
}

func (it *BundleAggregatorProxyAggregatorProposedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BundleAggregatorProxyAggregatorProposed struct {
	Current  common.Address
	Proposed common.Address
	Raw      types.Log
}

func (_BundleAggregatorProxy *BundleAggregatorProxyFilterer) FilterAggregatorProposed(opts *bind.FilterOpts, current []common.Address, proposed []common.Address) (*BundleAggregatorProxyAggregatorProposedIterator, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var proposedRule []interface{}
	for _, proposedItem := range proposed {
		proposedRule = append(proposedRule, proposedItem)
	}

	logs, sub, err := _BundleAggregatorProxy.contract.FilterLogs(opts, "AggregatorProposed", currentRule, proposedRule)
	if err != nil {
		return nil, err
	}
	return &BundleAggregatorProxyAggregatorProposedIterator{contract: _BundleAggregatorProxy.contract, event: "AggregatorProposed", logs: logs, sub: sub}, nil
}

func (_BundleAggregatorProxy *BundleAggregatorProxyFilterer) WatchAggregatorProposed(opts *bind.WatchOpts, sink chan<- *BundleAggregatorProxyAggregatorProposed, current []common.Address, proposed []common.Address) (event.Subscription, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var proposedRule []interface{}
	for _, proposedItem := range proposed {
		proposedRule = append(proposedRule, proposedItem)
	}

	logs, sub, err := _BundleAggregatorProxy.contract.WatchLogs(opts, "AggregatorProposed", currentRule, proposedRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BundleAggregatorProxyAggregatorProposed)
				if err := _BundleAggregatorProxy.contract.UnpackLog(event, "AggregatorProposed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_BundleAggregatorProxy *BundleAggregatorProxyFilterer) ParseAggregatorProposed(log types.Log) (*BundleAggregatorProxyAggregatorProposed, error) {
	event := new(BundleAggregatorProxyAggregatorProposed)
	if err := _BundleAggregatorProxy.contract.UnpackLog(event, "AggregatorProposed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BundleAggregatorProxyOwnershipTransferRequestedIterator struct {
	Event *BundleAggregatorProxyOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BundleAggregatorProxyOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BundleAggregatorProxyOwnershipTransferRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(BundleAggregatorProxyOwnershipTransferRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *BundleAggregatorProxyOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *BundleAggregatorProxyOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BundleAggregatorProxyOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_BundleAggregatorProxy *BundleAggregatorProxyFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BundleAggregatorProxyOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BundleAggregatorProxy.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &BundleAggregatorProxyOwnershipTransferRequestedIterator{contract: _BundleAggregatorProxy.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_BundleAggregatorProxy *BundleAggregatorProxyFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *BundleAggregatorProxyOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BundleAggregatorProxy.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BundleAggregatorProxyOwnershipTransferRequested)
				if err := _BundleAggregatorProxy.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_BundleAggregatorProxy *BundleAggregatorProxyFilterer) ParseOwnershipTransferRequested(log types.Log) (*BundleAggregatorProxyOwnershipTransferRequested, error) {
	event := new(BundleAggregatorProxyOwnershipTransferRequested)
	if err := _BundleAggregatorProxy.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type BundleAggregatorProxyOwnershipTransferredIterator struct {
	Event *BundleAggregatorProxyOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *BundleAggregatorProxyOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BundleAggregatorProxyOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(BundleAggregatorProxyOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *BundleAggregatorProxyOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *BundleAggregatorProxyOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type BundleAggregatorProxyOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_BundleAggregatorProxy *BundleAggregatorProxyFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BundleAggregatorProxyOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BundleAggregatorProxy.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &BundleAggregatorProxyOwnershipTransferredIterator{contract: _BundleAggregatorProxy.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_BundleAggregatorProxy *BundleAggregatorProxyFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BundleAggregatorProxyOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BundleAggregatorProxy.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(BundleAggregatorProxyOwnershipTransferred)
				if err := _BundleAggregatorProxy.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_BundleAggregatorProxy *BundleAggregatorProxyFilterer) ParseOwnershipTransferred(log types.Log) (*BundleAggregatorProxyOwnershipTransferred, error) {
	event := new(BundleAggregatorProxyOwnershipTransferred)
	if err := _BundleAggregatorProxy.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_BundleAggregatorProxy *BundleAggregatorProxy) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _BundleAggregatorProxy.abi.Events["AggregatorConfirmed"].ID:
		return _BundleAggregatorProxy.ParseAggregatorConfirmed(log)
	case _BundleAggregatorProxy.abi.Events["AggregatorProposed"].ID:
		return _BundleAggregatorProxy.ParseAggregatorProposed(log)
	case _BundleAggregatorProxy.abi.Events["OwnershipTransferRequested"].ID:
		return _BundleAggregatorProxy.ParseOwnershipTransferRequested(log)
	case _BundleAggregatorProxy.abi.Events["OwnershipTransferred"].ID:
		return _BundleAggregatorProxy.ParseOwnershipTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (BundleAggregatorProxyAggregatorConfirmed) Topic() common.Hash {
	return common.HexToHash("0x33745f67a407dcb785417f9c123dd3641479a102674b6e35c1f10975625b90e9")
}

func (BundleAggregatorProxyAggregatorProposed) Topic() common.Hash {
	return common.HexToHash("0xc0f151710f03d713b71d9970cee0d5b11ddc9a7552abaa3f6ee818010f21600d")
}

func (BundleAggregatorProxyOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (BundleAggregatorProxyOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (_BundleAggregatorProxy *BundleAggregatorProxy) Address() common.Address {
	return _BundleAggregatorProxy.address
}

type BundleAggregatorProxyInterface interface {
	Aggregator(opts *bind.CallOpts) (common.Address, error)

	BundleDecimals(opts *bind.CallOpts) ([]uint8, error)

	Description(opts *bind.CallOpts) (string, error)

	LatestBundle(opts *bind.CallOpts) ([]byte, error)

	LatestBundleTimestamp(opts *bind.CallOpts) (*big.Int, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	ProposedAggregator(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	Version(opts *bind.CallOpts) (*big.Int, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	ConfirmAggregator(opts *bind.TransactOpts, aggregatorAddress common.Address) (*types.Transaction, error)

	ProposeAggregator(opts *bind.TransactOpts, aggregatorAddress common.Address) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterAggregatorConfirmed(opts *bind.FilterOpts, previous []common.Address, latest []common.Address) (*BundleAggregatorProxyAggregatorConfirmedIterator, error)

	WatchAggregatorConfirmed(opts *bind.WatchOpts, sink chan<- *BundleAggregatorProxyAggregatorConfirmed, previous []common.Address, latest []common.Address) (event.Subscription, error)

	ParseAggregatorConfirmed(log types.Log) (*BundleAggregatorProxyAggregatorConfirmed, error)

	FilterAggregatorProposed(opts *bind.FilterOpts, current []common.Address, proposed []common.Address) (*BundleAggregatorProxyAggregatorProposedIterator, error)

	WatchAggregatorProposed(opts *bind.WatchOpts, sink chan<- *BundleAggregatorProxyAggregatorProposed, current []common.Address, proposed []common.Address) (event.Subscription, error)

	ParseAggregatorProposed(log types.Log) (*BundleAggregatorProxyAggregatorProposed, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BundleAggregatorProxyOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *BundleAggregatorProxyOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*BundleAggregatorProxyOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*BundleAggregatorProxyOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BundleAggregatorProxyOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*BundleAggregatorProxyOwnershipTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
