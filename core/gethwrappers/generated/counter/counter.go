// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package counter

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

var CounterMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"alwaysRevert\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"alwaysRevertWithString\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"count\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"increment\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"reset\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"error\",\"name\":\"AlwaysRevert\",\"inputs\":[]}]",
	Bin: "0x60806040526000805534801561001457600080fd5b506101af806100246000396000f3fe608060405234801561001057600080fd5b50600436106100675760003560e01c8063a7bc8cbc11610050578063a7bc8cbc14610091578063d09de08a14610099578063d826f88f146100a157600080fd5b806306661abd1461006c5780639fb3785314610087575b600080fd5b61007560005481565b60405190815260200160405180910390f35b61008f6100aa565b005b61008f6100dc565b610075610142565b61008f60008055565b6040517f8bba4aff00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f616c776179732072657665727400000000000000000000000000000000000000604482015260640160405180910390fd5b600060016000808282546101569190610163565b9091555050600054919050565b6000821982111561019d577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b50019056fea164736f6c6343000806000a",
}

var CounterABI = CounterMetaData.ABI

var CounterBin = CounterMetaData.Bin

func DeployCounter(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Counter, error) {
	parsed, err := CounterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CounterBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Counter{address: address, abi: *parsed, CounterCaller: CounterCaller{contract: contract}, CounterTransactor: CounterTransactor{contract: contract}, CounterFilterer: CounterFilterer{contract: contract}}, nil
}

type Counter struct {
	address common.Address
	abi     abi.ABI
	CounterCaller
	CounterTransactor
	CounterFilterer
}

type CounterCaller struct {
	contract *bind.BoundContract
}

type CounterTransactor struct {
	contract *bind.BoundContract
}

type CounterFilterer struct {
	contract *bind.BoundContract
}

type CounterSession struct {
	Contract     *Counter
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type CounterCallerSession struct {
	Contract *CounterCaller
	CallOpts bind.CallOpts
}

type CounterTransactorSession struct {
	Contract     *CounterTransactor
	TransactOpts bind.TransactOpts
}

type CounterRaw struct {
	Contract *Counter
}

type CounterCallerRaw struct {
	Contract *CounterCaller
}

type CounterTransactorRaw struct {
	Contract *CounterTransactor
}

func NewCounter(address common.Address, backend bind.ContractBackend) (*Counter, error) {
	abi, err := abi.JSON(strings.NewReader(CounterABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindCounter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Counter{address: address, abi: abi, CounterCaller: CounterCaller{contract: contract}, CounterTransactor: CounterTransactor{contract: contract}, CounterFilterer: CounterFilterer{contract: contract}}, nil
}

func NewCounterCaller(address common.Address, caller bind.ContractCaller) (*CounterCaller, error) {
	contract, err := bindCounter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CounterCaller{contract: contract}, nil
}

func NewCounterTransactor(address common.Address, transactor bind.ContractTransactor) (*CounterTransactor, error) {
	contract, err := bindCounter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CounterTransactor{contract: contract}, nil
}

func NewCounterFilterer(address common.Address, filterer bind.ContractFilterer) (*CounterFilterer, error) {
	contract, err := bindCounter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CounterFilterer{contract: contract}, nil
}

func bindCounter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CounterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_Counter *CounterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Counter.Contract.CounterCaller.contract.Call(opts, result, method, params...)
}

func (_Counter *CounterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Counter.Contract.CounterTransactor.contract.Transfer(opts)
}

func (_Counter *CounterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Counter.Contract.CounterTransactor.contract.Transact(opts, method, params...)
}

func (_Counter *CounterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Counter.Contract.contract.Call(opts, result, method, params...)
}

func (_Counter *CounterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Counter.Contract.contract.Transfer(opts)
}

func (_Counter *CounterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Counter.Contract.contract.Transact(opts, method, params...)
}

func (_Counter *CounterCaller) AlwaysRevert(opts *bind.CallOpts) error {
	var out []interface{}
	err := _Counter.contract.Call(opts, &out, "alwaysRevert")

	if err != nil {
		return err
	}

	return err

}

func (_Counter *CounterSession) AlwaysRevert() error {
	return _Counter.Contract.AlwaysRevert(&_Counter.CallOpts)
}

func (_Counter *CounterCallerSession) AlwaysRevert() error {
	return _Counter.Contract.AlwaysRevert(&_Counter.CallOpts)
}

func (_Counter *CounterCaller) AlwaysRevertWithString(opts *bind.CallOpts) error {
	var out []interface{}
	err := _Counter.contract.Call(opts, &out, "alwaysRevertWithString")

	if err != nil {
		return err
	}

	return err

}

func (_Counter *CounterSession) AlwaysRevertWithString() error {
	return _Counter.Contract.AlwaysRevertWithString(&_Counter.CallOpts)
}

func (_Counter *CounterCallerSession) AlwaysRevertWithString() error {
	return _Counter.Contract.AlwaysRevertWithString(&_Counter.CallOpts)
}

func (_Counter *CounterCaller) Count(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Counter.contract.Call(opts, &out, "count")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_Counter *CounterSession) Count() (*big.Int, error) {
	return _Counter.Contract.Count(&_Counter.CallOpts)
}

func (_Counter *CounterCallerSession) Count() (*big.Int, error) {
	return _Counter.Contract.Count(&_Counter.CallOpts)
}

func (_Counter *CounterTransactor) Increment(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Counter.contract.Transact(opts, "increment")
}

func (_Counter *CounterSession) Increment() (*types.Transaction, error) {
	return _Counter.Contract.Increment(&_Counter.TransactOpts)
}

func (_Counter *CounterTransactorSession) Increment() (*types.Transaction, error) {
	return _Counter.Contract.Increment(&_Counter.TransactOpts)
}

func (_Counter *CounterTransactor) Reset(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Counter.contract.Transact(opts, "reset")
}

func (_Counter *CounterSession) Reset() (*types.Transaction, error) {
	return _Counter.Contract.Reset(&_Counter.TransactOpts)
}

func (_Counter *CounterTransactorSession) Reset() (*types.Transaction, error) {
	return _Counter.Contract.Reset(&_Counter.TransactOpts)
}

func (_Counter *Counter) Address() common.Address {
	return _Counter.address
}

type CounterInterface interface {
	AlwaysRevert(opts *bind.CallOpts) error

	AlwaysRevertWithString(opts *bind.CallOpts) error

	Count(opts *bind.CallOpts) (*big.Int, error)

	Increment(opts *bind.TransactOpts) (*types.Transaction, error)

	Reset(opts *bind.TransactOpts) (*types.Transaction, error)

	Address() common.Address
}
