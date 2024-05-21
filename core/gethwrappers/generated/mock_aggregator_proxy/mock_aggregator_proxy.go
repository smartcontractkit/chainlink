// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package mock_aggregator_proxy

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

var MockAggregatorProxyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_aggregator\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"aggregator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_aggregator\",\"type\":\"address\"}],\"name\":\"updateAggregator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5060405161019138038061019183398101604081905261002f91610054565b600080546001600160a01b0319166001600160a01b0392909216919091179055610084565b60006020828403121561006657600080fd5b81516001600160a01b038116811461007d57600080fd5b9392505050565b60ff806100926000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c8063245a7bfc1460375780639fe4ee47146063575b600080fd5b6000546040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b60b5606e36600460b7565b600080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b005b60006020828403121560c857600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811460eb57600080fd5b939250505056fea164736f6c6343000806000a",
}

var MockAggregatorProxyABI = MockAggregatorProxyMetaData.ABI

var MockAggregatorProxyBin = MockAggregatorProxyMetaData.Bin

func DeployMockAggregatorProxy(auth *bind.TransactOpts, backend bind.ContractBackend, _aggregator common.Address) (common.Address, *types.Transaction, *MockAggregatorProxy, error) {
	parsed, err := MockAggregatorProxyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MockAggregatorProxyBin), backend, _aggregator)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MockAggregatorProxy{address: address, abi: *parsed, MockAggregatorProxyCaller: MockAggregatorProxyCaller{contract: contract}, MockAggregatorProxyTransactor: MockAggregatorProxyTransactor{contract: contract}, MockAggregatorProxyFilterer: MockAggregatorProxyFilterer{contract: contract}}, nil
}

type MockAggregatorProxy struct {
	address common.Address
	abi     abi.ABI
	MockAggregatorProxyCaller
	MockAggregatorProxyTransactor
	MockAggregatorProxyFilterer
}

type MockAggregatorProxyCaller struct {
	contract *bind.BoundContract
}

type MockAggregatorProxyTransactor struct {
	contract *bind.BoundContract
}

type MockAggregatorProxyFilterer struct {
	contract *bind.BoundContract
}

type MockAggregatorProxySession struct {
	Contract     *MockAggregatorProxy
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type MockAggregatorProxyCallerSession struct {
	Contract *MockAggregatorProxyCaller
	CallOpts bind.CallOpts
}

type MockAggregatorProxyTransactorSession struct {
	Contract     *MockAggregatorProxyTransactor
	TransactOpts bind.TransactOpts
}

type MockAggregatorProxyRaw struct {
	Contract *MockAggregatorProxy
}

type MockAggregatorProxyCallerRaw struct {
	Contract *MockAggregatorProxyCaller
}

type MockAggregatorProxyTransactorRaw struct {
	Contract *MockAggregatorProxyTransactor
}

func NewMockAggregatorProxy(address common.Address, backend bind.ContractBackend) (*MockAggregatorProxy, error) {
	abi, err := abi.JSON(strings.NewReader(MockAggregatorProxyABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindMockAggregatorProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MockAggregatorProxy{address: address, abi: abi, MockAggregatorProxyCaller: MockAggregatorProxyCaller{contract: contract}, MockAggregatorProxyTransactor: MockAggregatorProxyTransactor{contract: contract}, MockAggregatorProxyFilterer: MockAggregatorProxyFilterer{contract: contract}}, nil
}

func NewMockAggregatorProxyCaller(address common.Address, caller bind.ContractCaller) (*MockAggregatorProxyCaller, error) {
	contract, err := bindMockAggregatorProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MockAggregatorProxyCaller{contract: contract}, nil
}

func NewMockAggregatorProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*MockAggregatorProxyTransactor, error) {
	contract, err := bindMockAggregatorProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MockAggregatorProxyTransactor{contract: contract}, nil
}

func NewMockAggregatorProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*MockAggregatorProxyFilterer, error) {
	contract, err := bindMockAggregatorProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MockAggregatorProxyFilterer{contract: contract}, nil
}

func bindMockAggregatorProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MockAggregatorProxyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_MockAggregatorProxy *MockAggregatorProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockAggregatorProxy.Contract.MockAggregatorProxyCaller.contract.Call(opts, result, method, params...)
}

func (_MockAggregatorProxy *MockAggregatorProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockAggregatorProxy.Contract.MockAggregatorProxyTransactor.contract.Transfer(opts)
}

func (_MockAggregatorProxy *MockAggregatorProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockAggregatorProxy.Contract.MockAggregatorProxyTransactor.contract.Transact(opts, method, params...)
}

func (_MockAggregatorProxy *MockAggregatorProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockAggregatorProxy.Contract.contract.Call(opts, result, method, params...)
}

func (_MockAggregatorProxy *MockAggregatorProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockAggregatorProxy.Contract.contract.Transfer(opts)
}

func (_MockAggregatorProxy *MockAggregatorProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockAggregatorProxy.Contract.contract.Transact(opts, method, params...)
}

func (_MockAggregatorProxy *MockAggregatorProxyCaller) Aggregator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MockAggregatorProxy.contract.Call(opts, &out, "aggregator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_MockAggregatorProxy *MockAggregatorProxySession) Aggregator() (common.Address, error) {
	return _MockAggregatorProxy.Contract.Aggregator(&_MockAggregatorProxy.CallOpts)
}

func (_MockAggregatorProxy *MockAggregatorProxyCallerSession) Aggregator() (common.Address, error) {
	return _MockAggregatorProxy.Contract.Aggregator(&_MockAggregatorProxy.CallOpts)
}

func (_MockAggregatorProxy *MockAggregatorProxyTransactor) UpdateAggregator(opts *bind.TransactOpts, _aggregator common.Address) (*types.Transaction, error) {
	return _MockAggregatorProxy.contract.Transact(opts, "updateAggregator", _aggregator)
}

func (_MockAggregatorProxy *MockAggregatorProxySession) UpdateAggregator(_aggregator common.Address) (*types.Transaction, error) {
	return _MockAggregatorProxy.Contract.UpdateAggregator(&_MockAggregatorProxy.TransactOpts, _aggregator)
}

func (_MockAggregatorProxy *MockAggregatorProxyTransactorSession) UpdateAggregator(_aggregator common.Address) (*types.Transaction, error) {
	return _MockAggregatorProxy.Contract.UpdateAggregator(&_MockAggregatorProxy.TransactOpts, _aggregator)
}

func (_MockAggregatorProxy *MockAggregatorProxy) Address() common.Address {
	return _MockAggregatorProxy.address
}

type MockAggregatorProxyInterface interface {
	Aggregator(opts *bind.CallOpts) (common.Address, error)

	UpdateAggregator(opts *bind.TransactOpts, _aggregator common.Address) (*types.Transaction, error)

	Address() common.Address
}
