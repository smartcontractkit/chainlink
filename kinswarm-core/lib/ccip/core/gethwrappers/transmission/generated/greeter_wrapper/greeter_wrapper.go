// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package greeter_wrapper

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

var GreeterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"getGreeting\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"greeting\",\"type\":\"string\"}],\"name\":\"setGreeting\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610443806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c8063a41368621461003b578063fe50cc7214610050575b600080fd5b61004e61004936600461013f565b61006e565b005b61005861007e565b604051610065919061020e565b60405180910390f35b600061007a828261031c565b5050565b60606000805461008d9061027a565b80601f01602080910402602001604051908101604052809291908181526020018280546100b99061027a565b80156101065780601f106100db57610100808354040283529160200191610106565b820191906000526020600020905b8154815290600101906020018083116100e957829003601f168201915b5050505050905090565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60006020828403121561015157600080fd5b813567ffffffffffffffff8082111561016957600080fd5b818401915084601f83011261017d57600080fd5b81358181111561018f5761018f610110565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f011681019083821181831017156101d5576101d5610110565b816040528281528760208487010111156101ee57600080fd5b826020860160208301376000928101602001929092525095945050505050565b600060208083528351808285015260005b8181101561023b5785810183015185820160400152820161021f565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b600181811c9082168061028e57607f821691505b6020821081036102c7577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f82111561031757600081815260208120601f850160051c810160208610156102f45750805b601f850160051c820191505b8181101561031357828155600101610300565b5050505b505050565b815167ffffffffffffffff81111561033657610336610110565b61034a81610344845461027a565b846102cd565b602080601f83116001811461039d57600084156103675750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b178555610313565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b828110156103ea578886015182559484019460019091019084016103cb565b508582101561042657878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b0190555056fea164736f6c6343000813000a",
}

var GreeterABI = GreeterMetaData.ABI

var GreeterBin = GreeterMetaData.Bin

func DeployGreeter(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Greeter, error) {
	parsed, err := GreeterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(GreeterBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Greeter{address: address, abi: *parsed, GreeterCaller: GreeterCaller{contract: contract}, GreeterTransactor: GreeterTransactor{contract: contract}, GreeterFilterer: GreeterFilterer{contract: contract}}, nil
}

type Greeter struct {
	address common.Address
	abi     abi.ABI
	GreeterCaller
	GreeterTransactor
	GreeterFilterer
}

type GreeterCaller struct {
	contract *bind.BoundContract
}

type GreeterTransactor struct {
	contract *bind.BoundContract
}

type GreeterFilterer struct {
	contract *bind.BoundContract
}

type GreeterSession struct {
	Contract     *Greeter
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type GreeterCallerSession struct {
	Contract *GreeterCaller
	CallOpts bind.CallOpts
}

type GreeterTransactorSession struct {
	Contract     *GreeterTransactor
	TransactOpts bind.TransactOpts
}

type GreeterRaw struct {
	Contract *Greeter
}

type GreeterCallerRaw struct {
	Contract *GreeterCaller
}

type GreeterTransactorRaw struct {
	Contract *GreeterTransactor
}

func NewGreeter(address common.Address, backend bind.ContractBackend) (*Greeter, error) {
	abi, err := abi.JSON(strings.NewReader(GreeterABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindGreeter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Greeter{address: address, abi: abi, GreeterCaller: GreeterCaller{contract: contract}, GreeterTransactor: GreeterTransactor{contract: contract}, GreeterFilterer: GreeterFilterer{contract: contract}}, nil
}

func NewGreeterCaller(address common.Address, caller bind.ContractCaller) (*GreeterCaller, error) {
	contract, err := bindGreeter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GreeterCaller{contract: contract}, nil
}

func NewGreeterTransactor(address common.Address, transactor bind.ContractTransactor) (*GreeterTransactor, error) {
	contract, err := bindGreeter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GreeterTransactor{contract: contract}, nil
}

func NewGreeterFilterer(address common.Address, filterer bind.ContractFilterer) (*GreeterFilterer, error) {
	contract, err := bindGreeter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GreeterFilterer{contract: contract}, nil
}

func bindGreeter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := GreeterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_Greeter *GreeterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Greeter.Contract.GreeterCaller.contract.Call(opts, result, method, params...)
}

func (_Greeter *GreeterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Greeter.Contract.GreeterTransactor.contract.Transfer(opts)
}

func (_Greeter *GreeterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Greeter.Contract.GreeterTransactor.contract.Transact(opts, method, params...)
}

func (_Greeter *GreeterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Greeter.Contract.contract.Call(opts, result, method, params...)
}

func (_Greeter *GreeterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Greeter.Contract.contract.Transfer(opts)
}

func (_Greeter *GreeterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Greeter.Contract.contract.Transact(opts, method, params...)
}

func (_Greeter *GreeterCaller) GetGreeting(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Greeter.contract.Call(opts, &out, "getGreeting")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_Greeter *GreeterSession) GetGreeting() (string, error) {
	return _Greeter.Contract.GetGreeting(&_Greeter.CallOpts)
}

func (_Greeter *GreeterCallerSession) GetGreeting() (string, error) {
	return _Greeter.Contract.GetGreeting(&_Greeter.CallOpts)
}

func (_Greeter *GreeterTransactor) SetGreeting(opts *bind.TransactOpts, greeting string) (*types.Transaction, error) {
	return _Greeter.contract.Transact(opts, "setGreeting", greeting)
}

func (_Greeter *GreeterSession) SetGreeting(greeting string) (*types.Transaction, error) {
	return _Greeter.Contract.SetGreeting(&_Greeter.TransactOpts, greeting)
}

func (_Greeter *GreeterTransactorSession) SetGreeting(greeting string) (*types.Transaction, error) {
	return _Greeter.Contract.SetGreeting(&_Greeter.TransactOpts, greeting)
}

func (_Greeter *Greeter) Address() common.Address {
	return _Greeter.address
}

type GreeterInterface interface {
	GetGreeting(opts *bind.CallOpts) (string, error)

	SetGreeting(opts *bind.TransactOpts, greeting string) (*types.Transaction, error)

	Address() common.Address
}
