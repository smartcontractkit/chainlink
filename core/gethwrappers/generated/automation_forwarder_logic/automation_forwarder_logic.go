// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package automation_forwarder_logic

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

var AutomationForwarderLogicMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"getRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIAutomationRegistryConsumer\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"typeAndVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"updateRegistry\",\"inputs\":[{\"name\":\"newRegistry\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"}]",
	Bin: "0x608060405234801561001057600080fd5b506101f6806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c8063181f5a77146100465780631a5da6c8146100985780635ab1bd53146100ad575b600080fd5b6100826040518060400160405280601981526020017f4175746f6d6174696f6e466f7277617264657220312e302e300000000000000081525081565b60405161008f9190610140565b60405180910390f35b6100ab6100a63660046101ac565b6100d5565b005b60005460405173ffffffffffffffffffffffffffffffffffffffff909116815260200161008f565b60005473ffffffffffffffffffffffffffffffffffffffff1633146100f957600080fd5b600080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b600060208083528351808285015260005b8181101561016d57858101830151858201604001528201610151565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b6000602082840312156101be57600080fd5b813573ffffffffffffffffffffffffffffffffffffffff811681146101e257600080fd5b939250505056fea164736f6c6343000810000a",
}

var AutomationForwarderLogicABI = AutomationForwarderLogicMetaData.ABI

var AutomationForwarderLogicBin = AutomationForwarderLogicMetaData.Bin

func DeployAutomationForwarderLogic(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *AutomationForwarderLogic, error) {
	parsed, err := AutomationForwarderLogicMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(AutomationForwarderLogicBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AutomationForwarderLogic{address: address, abi: *parsed, AutomationForwarderLogicCaller: AutomationForwarderLogicCaller{contract: contract}, AutomationForwarderLogicTransactor: AutomationForwarderLogicTransactor{contract: contract}, AutomationForwarderLogicFilterer: AutomationForwarderLogicFilterer{contract: contract}}, nil
}

type AutomationForwarderLogic struct {
	address common.Address
	abi     abi.ABI
	AutomationForwarderLogicCaller
	AutomationForwarderLogicTransactor
	AutomationForwarderLogicFilterer
}

type AutomationForwarderLogicCaller struct {
	contract *bind.BoundContract
}

type AutomationForwarderLogicTransactor struct {
	contract *bind.BoundContract
}

type AutomationForwarderLogicFilterer struct {
	contract *bind.BoundContract
}

type AutomationForwarderLogicSession struct {
	Contract     *AutomationForwarderLogic
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type AutomationForwarderLogicCallerSession struct {
	Contract *AutomationForwarderLogicCaller
	CallOpts bind.CallOpts
}

type AutomationForwarderLogicTransactorSession struct {
	Contract     *AutomationForwarderLogicTransactor
	TransactOpts bind.TransactOpts
}

type AutomationForwarderLogicRaw struct {
	Contract *AutomationForwarderLogic
}

type AutomationForwarderLogicCallerRaw struct {
	Contract *AutomationForwarderLogicCaller
}

type AutomationForwarderLogicTransactorRaw struct {
	Contract *AutomationForwarderLogicTransactor
}

func NewAutomationForwarderLogic(address common.Address, backend bind.ContractBackend) (*AutomationForwarderLogic, error) {
	abi, err := abi.JSON(strings.NewReader(AutomationForwarderLogicABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindAutomationForwarderLogic(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AutomationForwarderLogic{address: address, abi: abi, AutomationForwarderLogicCaller: AutomationForwarderLogicCaller{contract: contract}, AutomationForwarderLogicTransactor: AutomationForwarderLogicTransactor{contract: contract}, AutomationForwarderLogicFilterer: AutomationForwarderLogicFilterer{contract: contract}}, nil
}

func NewAutomationForwarderLogicCaller(address common.Address, caller bind.ContractCaller) (*AutomationForwarderLogicCaller, error) {
	contract, err := bindAutomationForwarderLogic(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AutomationForwarderLogicCaller{contract: contract}, nil
}

func NewAutomationForwarderLogicTransactor(address common.Address, transactor bind.ContractTransactor) (*AutomationForwarderLogicTransactor, error) {
	contract, err := bindAutomationForwarderLogic(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AutomationForwarderLogicTransactor{contract: contract}, nil
}

func NewAutomationForwarderLogicFilterer(address common.Address, filterer bind.ContractFilterer) (*AutomationForwarderLogicFilterer, error) {
	contract, err := bindAutomationForwarderLogic(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AutomationForwarderLogicFilterer{contract: contract}, nil
}

func bindAutomationForwarderLogic(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AutomationForwarderLogicMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_AutomationForwarderLogic *AutomationForwarderLogicRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AutomationForwarderLogic.Contract.AutomationForwarderLogicCaller.contract.Call(opts, result, method, params...)
}

func (_AutomationForwarderLogic *AutomationForwarderLogicRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationForwarderLogic.Contract.AutomationForwarderLogicTransactor.contract.Transfer(opts)
}

func (_AutomationForwarderLogic *AutomationForwarderLogicRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AutomationForwarderLogic.Contract.AutomationForwarderLogicTransactor.contract.Transact(opts, method, params...)
}

func (_AutomationForwarderLogic *AutomationForwarderLogicCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AutomationForwarderLogic.Contract.contract.Call(opts, result, method, params...)
}

func (_AutomationForwarderLogic *AutomationForwarderLogicTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AutomationForwarderLogic.Contract.contract.Transfer(opts)
}

func (_AutomationForwarderLogic *AutomationForwarderLogicTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AutomationForwarderLogic.Contract.contract.Transact(opts, method, params...)
}

func (_AutomationForwarderLogic *AutomationForwarderLogicCaller) GetRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AutomationForwarderLogic.contract.Call(opts, &out, "getRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_AutomationForwarderLogic *AutomationForwarderLogicSession) GetRegistry() (common.Address, error) {
	return _AutomationForwarderLogic.Contract.GetRegistry(&_AutomationForwarderLogic.CallOpts)
}

func (_AutomationForwarderLogic *AutomationForwarderLogicCallerSession) GetRegistry() (common.Address, error) {
	return _AutomationForwarderLogic.Contract.GetRegistry(&_AutomationForwarderLogic.CallOpts)
}

func (_AutomationForwarderLogic *AutomationForwarderLogicCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AutomationForwarderLogic.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_AutomationForwarderLogic *AutomationForwarderLogicSession) TypeAndVersion() (string, error) {
	return _AutomationForwarderLogic.Contract.TypeAndVersion(&_AutomationForwarderLogic.CallOpts)
}

func (_AutomationForwarderLogic *AutomationForwarderLogicCallerSession) TypeAndVersion() (string, error) {
	return _AutomationForwarderLogic.Contract.TypeAndVersion(&_AutomationForwarderLogic.CallOpts)
}

func (_AutomationForwarderLogic *AutomationForwarderLogicTransactor) UpdateRegistry(opts *bind.TransactOpts, newRegistry common.Address) (*types.Transaction, error) {
	return _AutomationForwarderLogic.contract.Transact(opts, "updateRegistry", newRegistry)
}

func (_AutomationForwarderLogic *AutomationForwarderLogicSession) UpdateRegistry(newRegistry common.Address) (*types.Transaction, error) {
	return _AutomationForwarderLogic.Contract.UpdateRegistry(&_AutomationForwarderLogic.TransactOpts, newRegistry)
}

func (_AutomationForwarderLogic *AutomationForwarderLogicTransactorSession) UpdateRegistry(newRegistry common.Address) (*types.Transaction, error) {
	return _AutomationForwarderLogic.Contract.UpdateRegistry(&_AutomationForwarderLogic.TransactOpts, newRegistry)
}

func (_AutomationForwarderLogic *AutomationForwarderLogic) Address() common.Address {
	return _AutomationForwarderLogic.address
}

type AutomationForwarderLogicInterface interface {
	GetRegistry(opts *bind.CallOpts) (common.Address, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	UpdateRegistry(opts *bind.TransactOpts, newRegistry common.Address) (*types.Transaction, error)

	Address() common.Address
}
