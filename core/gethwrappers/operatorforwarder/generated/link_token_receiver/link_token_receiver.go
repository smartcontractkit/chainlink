// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package link_token_receiver

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

var LinkTokenReceiverMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"getChainlinkToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"onTokenTransfer\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"}]",
}

var LinkTokenReceiverABI = LinkTokenReceiverMetaData.ABI

type LinkTokenReceiver struct {
	address common.Address
	abi     abi.ABI
	LinkTokenReceiverCaller
	LinkTokenReceiverTransactor
	LinkTokenReceiverFilterer
}

type LinkTokenReceiverCaller struct {
	contract *bind.BoundContract
}

type LinkTokenReceiverTransactor struct {
	contract *bind.BoundContract
}

type LinkTokenReceiverFilterer struct {
	contract *bind.BoundContract
}

type LinkTokenReceiverSession struct {
	Contract     *LinkTokenReceiver
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type LinkTokenReceiverCallerSession struct {
	Contract *LinkTokenReceiverCaller
	CallOpts bind.CallOpts
}

type LinkTokenReceiverTransactorSession struct {
	Contract     *LinkTokenReceiverTransactor
	TransactOpts bind.TransactOpts
}

type LinkTokenReceiverRaw struct {
	Contract *LinkTokenReceiver
}

type LinkTokenReceiverCallerRaw struct {
	Contract *LinkTokenReceiverCaller
}

type LinkTokenReceiverTransactorRaw struct {
	Contract *LinkTokenReceiverTransactor
}

func NewLinkTokenReceiver(address common.Address, backend bind.ContractBackend) (*LinkTokenReceiver, error) {
	abi, err := abi.JSON(strings.NewReader(LinkTokenReceiverABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindLinkTokenReceiver(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LinkTokenReceiver{address: address, abi: abi, LinkTokenReceiverCaller: LinkTokenReceiverCaller{contract: contract}, LinkTokenReceiverTransactor: LinkTokenReceiverTransactor{contract: contract}, LinkTokenReceiverFilterer: LinkTokenReceiverFilterer{contract: contract}}, nil
}

func NewLinkTokenReceiverCaller(address common.Address, caller bind.ContractCaller) (*LinkTokenReceiverCaller, error) {
	contract, err := bindLinkTokenReceiver(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LinkTokenReceiverCaller{contract: contract}, nil
}

func NewLinkTokenReceiverTransactor(address common.Address, transactor bind.ContractTransactor) (*LinkTokenReceiverTransactor, error) {
	contract, err := bindLinkTokenReceiver(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LinkTokenReceiverTransactor{contract: contract}, nil
}

func NewLinkTokenReceiverFilterer(address common.Address, filterer bind.ContractFilterer) (*LinkTokenReceiverFilterer, error) {
	contract, err := bindLinkTokenReceiver(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LinkTokenReceiverFilterer{contract: contract}, nil
}

func bindLinkTokenReceiver(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := LinkTokenReceiverMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_LinkTokenReceiver *LinkTokenReceiverRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LinkTokenReceiver.Contract.LinkTokenReceiverCaller.contract.Call(opts, result, method, params...)
}

func (_LinkTokenReceiver *LinkTokenReceiverRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LinkTokenReceiver.Contract.LinkTokenReceiverTransactor.contract.Transfer(opts)
}

func (_LinkTokenReceiver *LinkTokenReceiverRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LinkTokenReceiver.Contract.LinkTokenReceiverTransactor.contract.Transact(opts, method, params...)
}

func (_LinkTokenReceiver *LinkTokenReceiverCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LinkTokenReceiver.Contract.contract.Call(opts, result, method, params...)
}

func (_LinkTokenReceiver *LinkTokenReceiverTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LinkTokenReceiver.Contract.contract.Transfer(opts)
}

func (_LinkTokenReceiver *LinkTokenReceiverTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LinkTokenReceiver.Contract.contract.Transact(opts, method, params...)
}

func (_LinkTokenReceiver *LinkTokenReceiverCaller) GetChainlinkToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LinkTokenReceiver.contract.Call(opts, &out, "getChainlinkToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_LinkTokenReceiver *LinkTokenReceiverSession) GetChainlinkToken() (common.Address, error) {
	return _LinkTokenReceiver.Contract.GetChainlinkToken(&_LinkTokenReceiver.CallOpts)
}

func (_LinkTokenReceiver *LinkTokenReceiverCallerSession) GetChainlinkToken() (common.Address, error) {
	return _LinkTokenReceiver.Contract.GetChainlinkToken(&_LinkTokenReceiver.CallOpts)
}

func (_LinkTokenReceiver *LinkTokenReceiverTransactor) OnTokenTransfer(opts *bind.TransactOpts, sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _LinkTokenReceiver.contract.Transact(opts, "onTokenTransfer", sender, amount, data)
}

func (_LinkTokenReceiver *LinkTokenReceiverSession) OnTokenTransfer(sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _LinkTokenReceiver.Contract.OnTokenTransfer(&_LinkTokenReceiver.TransactOpts, sender, amount, data)
}

func (_LinkTokenReceiver *LinkTokenReceiverTransactorSession) OnTokenTransfer(sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _LinkTokenReceiver.Contract.OnTokenTransfer(&_LinkTokenReceiver.TransactOpts, sender, amount, data)
}

func (_LinkTokenReceiver *LinkTokenReceiver) Address() common.Address {
	return _LinkTokenReceiver.address
}

type LinkTokenReceiverInterface interface {
	GetChainlinkToken(opts *bind.CallOpts) (common.Address, error)

	OnTokenTransfer(opts *bind.TransactOpts, sender common.Address, amount *big.Int, data []byte) (*types.Transaction, error)

	Address() common.Address
}
