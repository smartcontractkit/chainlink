// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package authorized_receiver

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
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated"
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
)

var AuthorizedReceiverMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"senders\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"changedBy\",\"type\":\"address\"}],\"name\":\"AuthorizedSendersChanged\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"getAuthorizedSenders\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"isAuthorizedSender\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"senders\",\"type\":\"address[]\"}],\"name\":\"setAuthorizedSenders\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

var AuthorizedReceiverABI = AuthorizedReceiverMetaData.ABI

type AuthorizedReceiver struct {
	address common.Address
	abi     abi.ABI
	AuthorizedReceiverCaller
	AuthorizedReceiverTransactor
	AuthorizedReceiverFilterer
}

type AuthorizedReceiverCaller struct {
	contract *bind.BoundContract
}

type AuthorizedReceiverTransactor struct {
	contract *bind.BoundContract
}

type AuthorizedReceiverFilterer struct {
	contract *bind.BoundContract
}

type AuthorizedReceiverSession struct {
	Contract     *AuthorizedReceiver
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type AuthorizedReceiverCallerSession struct {
	Contract *AuthorizedReceiverCaller
	CallOpts bind.CallOpts
}

type AuthorizedReceiverTransactorSession struct {
	Contract     *AuthorizedReceiverTransactor
	TransactOpts bind.TransactOpts
}

type AuthorizedReceiverRaw struct {
	Contract *AuthorizedReceiver
}

type AuthorizedReceiverCallerRaw struct {
	Contract *AuthorizedReceiverCaller
}

type AuthorizedReceiverTransactorRaw struct {
	Contract *AuthorizedReceiverTransactor
}

func NewAuthorizedReceiver(address common.Address, backend bind.ContractBackend) (*AuthorizedReceiver, error) {
	abi, err := abi.JSON(strings.NewReader(AuthorizedReceiverABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindAuthorizedReceiver(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AuthorizedReceiver{address: address, abi: abi, AuthorizedReceiverCaller: AuthorizedReceiverCaller{contract: contract}, AuthorizedReceiverTransactor: AuthorizedReceiverTransactor{contract: contract}, AuthorizedReceiverFilterer: AuthorizedReceiverFilterer{contract: contract}}, nil
}

func NewAuthorizedReceiverCaller(address common.Address, caller bind.ContractCaller) (*AuthorizedReceiverCaller, error) {
	contract, err := bindAuthorizedReceiver(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AuthorizedReceiverCaller{contract: contract}, nil
}

func NewAuthorizedReceiverTransactor(address common.Address, transactor bind.ContractTransactor) (*AuthorizedReceiverTransactor, error) {
	contract, err := bindAuthorizedReceiver(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AuthorizedReceiverTransactor{contract: contract}, nil
}

func NewAuthorizedReceiverFilterer(address common.Address, filterer bind.ContractFilterer) (*AuthorizedReceiverFilterer, error) {
	contract, err := bindAuthorizedReceiver(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AuthorizedReceiverFilterer{contract: contract}, nil
}

func bindAuthorizedReceiver(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AuthorizedReceiverABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_AuthorizedReceiver *AuthorizedReceiverRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AuthorizedReceiver.Contract.AuthorizedReceiverCaller.contract.Call(opts, result, method, params...)
}

func (_AuthorizedReceiver *AuthorizedReceiverRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AuthorizedReceiver.Contract.AuthorizedReceiverTransactor.contract.Transfer(opts)
}

func (_AuthorizedReceiver *AuthorizedReceiverRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AuthorizedReceiver.Contract.AuthorizedReceiverTransactor.contract.Transact(opts, method, params...)
}

func (_AuthorizedReceiver *AuthorizedReceiverCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AuthorizedReceiver.Contract.contract.Call(opts, result, method, params...)
}

func (_AuthorizedReceiver *AuthorizedReceiverTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AuthorizedReceiver.Contract.contract.Transfer(opts)
}

func (_AuthorizedReceiver *AuthorizedReceiverTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AuthorizedReceiver.Contract.contract.Transact(opts, method, params...)
}

func (_AuthorizedReceiver *AuthorizedReceiverCaller) GetAuthorizedSenders(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _AuthorizedReceiver.contract.Call(opts, &out, "getAuthorizedSenders")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_AuthorizedReceiver *AuthorizedReceiverSession) GetAuthorizedSenders() ([]common.Address, error) {
	return _AuthorizedReceiver.Contract.GetAuthorizedSenders(&_AuthorizedReceiver.CallOpts)
}

func (_AuthorizedReceiver *AuthorizedReceiverCallerSession) GetAuthorizedSenders() ([]common.Address, error) {
	return _AuthorizedReceiver.Contract.GetAuthorizedSenders(&_AuthorizedReceiver.CallOpts)
}

func (_AuthorizedReceiver *AuthorizedReceiverCaller) IsAuthorizedSender(opts *bind.CallOpts, sender common.Address) (bool, error) {
	var out []interface{}
	err := _AuthorizedReceiver.contract.Call(opts, &out, "isAuthorizedSender", sender)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_AuthorizedReceiver *AuthorizedReceiverSession) IsAuthorizedSender(sender common.Address) (bool, error) {
	return _AuthorizedReceiver.Contract.IsAuthorizedSender(&_AuthorizedReceiver.CallOpts, sender)
}

func (_AuthorizedReceiver *AuthorizedReceiverCallerSession) IsAuthorizedSender(sender common.Address) (bool, error) {
	return _AuthorizedReceiver.Contract.IsAuthorizedSender(&_AuthorizedReceiver.CallOpts, sender)
}

func (_AuthorizedReceiver *AuthorizedReceiverTransactor) SetAuthorizedSenders(opts *bind.TransactOpts, senders []common.Address) (*types.Transaction, error) {
	return _AuthorizedReceiver.contract.Transact(opts, "setAuthorizedSenders", senders)
}

func (_AuthorizedReceiver *AuthorizedReceiverSession) SetAuthorizedSenders(senders []common.Address) (*types.Transaction, error) {
	return _AuthorizedReceiver.Contract.SetAuthorizedSenders(&_AuthorizedReceiver.TransactOpts, senders)
}

func (_AuthorizedReceiver *AuthorizedReceiverTransactorSession) SetAuthorizedSenders(senders []common.Address) (*types.Transaction, error) {
	return _AuthorizedReceiver.Contract.SetAuthorizedSenders(&_AuthorizedReceiver.TransactOpts, senders)
}

type AuthorizedReceiverAuthorizedSendersChangedIterator struct {
	Event *AuthorizedReceiverAuthorizedSendersChanged

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *AuthorizedReceiverAuthorizedSendersChangedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AuthorizedReceiverAuthorizedSendersChanged)
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
		it.Event = new(AuthorizedReceiverAuthorizedSendersChanged)
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

func (it *AuthorizedReceiverAuthorizedSendersChangedIterator) Error() error {
	return it.fail
}

func (it *AuthorizedReceiverAuthorizedSendersChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type AuthorizedReceiverAuthorizedSendersChanged struct {
	Senders   []common.Address
	ChangedBy common.Address
	Raw       types.Log
}

func (_AuthorizedReceiver *AuthorizedReceiverFilterer) FilterAuthorizedSendersChanged(opts *bind.FilterOpts) (*AuthorizedReceiverAuthorizedSendersChangedIterator, error) {

	logs, sub, err := _AuthorizedReceiver.contract.FilterLogs(opts, "AuthorizedSendersChanged")
	if err != nil {
		return nil, err
	}
	return &AuthorizedReceiverAuthorizedSendersChangedIterator{contract: _AuthorizedReceiver.contract, event: "AuthorizedSendersChanged", logs: logs, sub: sub}, nil
}

func (_AuthorizedReceiver *AuthorizedReceiverFilterer) WatchAuthorizedSendersChanged(opts *bind.WatchOpts, sink chan<- *AuthorizedReceiverAuthorizedSendersChanged) (event.Subscription, error) {

	logs, sub, err := _AuthorizedReceiver.contract.WatchLogs(opts, "AuthorizedSendersChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(AuthorizedReceiverAuthorizedSendersChanged)
				if err := _AuthorizedReceiver.contract.UnpackLog(event, "AuthorizedSendersChanged", log); err != nil {
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

func (_AuthorizedReceiver *AuthorizedReceiverFilterer) ParseAuthorizedSendersChanged(log types.Log) (*AuthorizedReceiverAuthorizedSendersChanged, error) {
	event := new(AuthorizedReceiverAuthorizedSendersChanged)
	if err := _AuthorizedReceiver.contract.UnpackLog(event, "AuthorizedSendersChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_AuthorizedReceiver *AuthorizedReceiver) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _AuthorizedReceiver.abi.Events["AuthorizedSendersChanged"].ID:
		return _AuthorizedReceiver.ParseAuthorizedSendersChanged(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (AuthorizedReceiverAuthorizedSendersChanged) Topic() common.Hash {
	return common.HexToHash("0xf263cfb3e4298332e776194610cf9fdc09ccb3ada8b9aa39764d882e11fbf0a0")
}

func (_AuthorizedReceiver *AuthorizedReceiver) Address() common.Address {
	return _AuthorizedReceiver.address
}

type AuthorizedReceiverInterface interface {
	GetAuthorizedSenders(opts *bind.CallOpts) ([]common.Address, error)

	IsAuthorizedSender(opts *bind.CallOpts, sender common.Address) (bool, error)

	SetAuthorizedSenders(opts *bind.TransactOpts, senders []common.Address) (*types.Transaction, error)

	FilterAuthorizedSendersChanged(opts *bind.FilterOpts) (*AuthorizedReceiverAuthorizedSendersChangedIterator, error)

	WatchAuthorizedSendersChanged(opts *bind.WatchOpts, sink chan<- *AuthorizedReceiverAuthorizedSendersChanged) (event.Subscription, error)

	ParseAuthorizedSendersChanged(log types.Log) (*AuthorizedReceiverAuthorizedSendersChanged, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
