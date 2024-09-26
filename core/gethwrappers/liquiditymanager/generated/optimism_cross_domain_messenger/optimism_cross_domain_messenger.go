// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package optimism_cross_domain_messenger

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

var OptimismCrossDomainMessengerMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"messageNonce\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"}],\"name\":\"SentMessage\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_nonce\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_target\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_minGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_message\",\"type\":\"bytes\"}],\"name\":\"relayMessage\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
}

var OptimismCrossDomainMessengerABI = OptimismCrossDomainMessengerMetaData.ABI

type OptimismCrossDomainMessenger struct {
	address common.Address
	abi     abi.ABI
	OptimismCrossDomainMessengerCaller
	OptimismCrossDomainMessengerTransactor
	OptimismCrossDomainMessengerFilterer
}

type OptimismCrossDomainMessengerCaller struct {
	contract *bind.BoundContract
}

type OptimismCrossDomainMessengerTransactor struct {
	contract *bind.BoundContract
}

type OptimismCrossDomainMessengerFilterer struct {
	contract *bind.BoundContract
}

type OptimismCrossDomainMessengerSession struct {
	Contract     *OptimismCrossDomainMessenger
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type OptimismCrossDomainMessengerCallerSession struct {
	Contract *OptimismCrossDomainMessengerCaller
	CallOpts bind.CallOpts
}

type OptimismCrossDomainMessengerTransactorSession struct {
	Contract     *OptimismCrossDomainMessengerTransactor
	TransactOpts bind.TransactOpts
}

type OptimismCrossDomainMessengerRaw struct {
	Contract *OptimismCrossDomainMessenger
}

type OptimismCrossDomainMessengerCallerRaw struct {
	Contract *OptimismCrossDomainMessengerCaller
}

type OptimismCrossDomainMessengerTransactorRaw struct {
	Contract *OptimismCrossDomainMessengerTransactor
}

func NewOptimismCrossDomainMessenger(address common.Address, backend bind.ContractBackend) (*OptimismCrossDomainMessenger, error) {
	abi, err := abi.JSON(strings.NewReader(OptimismCrossDomainMessengerABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindOptimismCrossDomainMessenger(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OptimismCrossDomainMessenger{address: address, abi: abi, OptimismCrossDomainMessengerCaller: OptimismCrossDomainMessengerCaller{contract: contract}, OptimismCrossDomainMessengerTransactor: OptimismCrossDomainMessengerTransactor{contract: contract}, OptimismCrossDomainMessengerFilterer: OptimismCrossDomainMessengerFilterer{contract: contract}}, nil
}

func NewOptimismCrossDomainMessengerCaller(address common.Address, caller bind.ContractCaller) (*OptimismCrossDomainMessengerCaller, error) {
	contract, err := bindOptimismCrossDomainMessenger(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OptimismCrossDomainMessengerCaller{contract: contract}, nil
}

func NewOptimismCrossDomainMessengerTransactor(address common.Address, transactor bind.ContractTransactor) (*OptimismCrossDomainMessengerTransactor, error) {
	contract, err := bindOptimismCrossDomainMessenger(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OptimismCrossDomainMessengerTransactor{contract: contract}, nil
}

func NewOptimismCrossDomainMessengerFilterer(address common.Address, filterer bind.ContractFilterer) (*OptimismCrossDomainMessengerFilterer, error) {
	contract, err := bindOptimismCrossDomainMessenger(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OptimismCrossDomainMessengerFilterer{contract: contract}, nil
}

func bindOptimismCrossDomainMessenger(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OptimismCrossDomainMessengerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_OptimismCrossDomainMessenger *OptimismCrossDomainMessengerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OptimismCrossDomainMessenger.Contract.OptimismCrossDomainMessengerCaller.contract.Call(opts, result, method, params...)
}

func (_OptimismCrossDomainMessenger *OptimismCrossDomainMessengerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OptimismCrossDomainMessenger.Contract.OptimismCrossDomainMessengerTransactor.contract.Transfer(opts)
}

func (_OptimismCrossDomainMessenger *OptimismCrossDomainMessengerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OptimismCrossDomainMessenger.Contract.OptimismCrossDomainMessengerTransactor.contract.Transact(opts, method, params...)
}

func (_OptimismCrossDomainMessenger *OptimismCrossDomainMessengerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OptimismCrossDomainMessenger.Contract.contract.Call(opts, result, method, params...)
}

func (_OptimismCrossDomainMessenger *OptimismCrossDomainMessengerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OptimismCrossDomainMessenger.Contract.contract.Transfer(opts)
}

func (_OptimismCrossDomainMessenger *OptimismCrossDomainMessengerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OptimismCrossDomainMessenger.Contract.contract.Transact(opts, method, params...)
}

func (_OptimismCrossDomainMessenger *OptimismCrossDomainMessengerTransactor) RelayMessage(opts *bind.TransactOpts, _nonce *big.Int, _sender common.Address, _target common.Address, _value *big.Int, _minGasLimit *big.Int, _message []byte) (*types.Transaction, error) {
	return _OptimismCrossDomainMessenger.contract.Transact(opts, "relayMessage", _nonce, _sender, _target, _value, _minGasLimit, _message)
}

func (_OptimismCrossDomainMessenger *OptimismCrossDomainMessengerSession) RelayMessage(_nonce *big.Int, _sender common.Address, _target common.Address, _value *big.Int, _minGasLimit *big.Int, _message []byte) (*types.Transaction, error) {
	return _OptimismCrossDomainMessenger.Contract.RelayMessage(&_OptimismCrossDomainMessenger.TransactOpts, _nonce, _sender, _target, _value, _minGasLimit, _message)
}

func (_OptimismCrossDomainMessenger *OptimismCrossDomainMessengerTransactorSession) RelayMessage(_nonce *big.Int, _sender common.Address, _target common.Address, _value *big.Int, _minGasLimit *big.Int, _message []byte) (*types.Transaction, error) {
	return _OptimismCrossDomainMessenger.Contract.RelayMessage(&_OptimismCrossDomainMessenger.TransactOpts, _nonce, _sender, _target, _value, _minGasLimit, _message)
}

type OptimismCrossDomainMessengerSentMessageIterator struct {
	Event *OptimismCrossDomainMessengerSentMessage

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OptimismCrossDomainMessengerSentMessageIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OptimismCrossDomainMessengerSentMessage)
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
		it.Event = new(OptimismCrossDomainMessengerSentMessage)
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

func (it *OptimismCrossDomainMessengerSentMessageIterator) Error() error {
	return it.fail
}

func (it *OptimismCrossDomainMessengerSentMessageIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OptimismCrossDomainMessengerSentMessage struct {
	Target       common.Address
	Sender       common.Address
	Message      []byte
	MessageNonce *big.Int
	GasLimit     *big.Int
	Raw          types.Log
}

func (_OptimismCrossDomainMessenger *OptimismCrossDomainMessengerFilterer) FilterSentMessage(opts *bind.FilterOpts, target []common.Address) (*OptimismCrossDomainMessengerSentMessageIterator, error) {

	var targetRule []interface{}
	for _, targetItem := range target {
		targetRule = append(targetRule, targetItem)
	}

	logs, sub, err := _OptimismCrossDomainMessenger.contract.FilterLogs(opts, "SentMessage", targetRule)
	if err != nil {
		return nil, err
	}
	return &OptimismCrossDomainMessengerSentMessageIterator{contract: _OptimismCrossDomainMessenger.contract, event: "SentMessage", logs: logs, sub: sub}, nil
}

func (_OptimismCrossDomainMessenger *OptimismCrossDomainMessengerFilterer) WatchSentMessage(opts *bind.WatchOpts, sink chan<- *OptimismCrossDomainMessengerSentMessage, target []common.Address) (event.Subscription, error) {

	var targetRule []interface{}
	for _, targetItem := range target {
		targetRule = append(targetRule, targetItem)
	}

	logs, sub, err := _OptimismCrossDomainMessenger.contract.WatchLogs(opts, "SentMessage", targetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OptimismCrossDomainMessengerSentMessage)
				if err := _OptimismCrossDomainMessenger.contract.UnpackLog(event, "SentMessage", log); err != nil {
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

func (_OptimismCrossDomainMessenger *OptimismCrossDomainMessengerFilterer) ParseSentMessage(log types.Log) (*OptimismCrossDomainMessengerSentMessage, error) {
	event := new(OptimismCrossDomainMessengerSentMessage)
	if err := _OptimismCrossDomainMessenger.contract.UnpackLog(event, "SentMessage", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_OptimismCrossDomainMessenger *OptimismCrossDomainMessenger) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _OptimismCrossDomainMessenger.abi.Events["SentMessage"].ID:
		return _OptimismCrossDomainMessenger.ParseSentMessage(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (OptimismCrossDomainMessengerSentMessage) Topic() common.Hash {
	return common.HexToHash("0xcb0f7ffd78f9aee47a248fae8db181db6eee833039123e026dcbff529522e52a")
}

func (_OptimismCrossDomainMessenger *OptimismCrossDomainMessenger) Address() common.Address {
	return _OptimismCrossDomainMessenger.address
}

type OptimismCrossDomainMessengerInterface interface {
	RelayMessage(opts *bind.TransactOpts, _nonce *big.Int, _sender common.Address, _target common.Address, _value *big.Int, _minGasLimit *big.Int, _message []byte) (*types.Transaction, error)

	FilterSentMessage(opts *bind.FilterOpts, target []common.Address) (*OptimismCrossDomainMessengerSentMessageIterator, error)

	WatchSentMessage(opts *bind.WatchOpts, sink chan<- *OptimismCrossDomainMessengerSentMessage, target []common.Address) (event.Subscription, error)

	ParseSentMessage(log types.Log) (*OptimismCrossDomainMessengerSentMessage, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
