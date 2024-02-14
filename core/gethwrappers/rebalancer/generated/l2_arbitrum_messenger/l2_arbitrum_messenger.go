// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package l2_arbitrum_messenger

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

var L2ArbitrumMessengerMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"TxToL1\",\"type\":\"event\"}]",
}

var L2ArbitrumMessengerABI = L2ArbitrumMessengerMetaData.ABI

type L2ArbitrumMessenger struct {
	address common.Address
	abi     abi.ABI
	L2ArbitrumMessengerCaller
	L2ArbitrumMessengerTransactor
	L2ArbitrumMessengerFilterer
}

type L2ArbitrumMessengerCaller struct {
	contract *bind.BoundContract
}

type L2ArbitrumMessengerTransactor struct {
	contract *bind.BoundContract
}

type L2ArbitrumMessengerFilterer struct {
	contract *bind.BoundContract
}

type L2ArbitrumMessengerSession struct {
	Contract     *L2ArbitrumMessenger
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type L2ArbitrumMessengerCallerSession struct {
	Contract *L2ArbitrumMessengerCaller
	CallOpts bind.CallOpts
}

type L2ArbitrumMessengerTransactorSession struct {
	Contract     *L2ArbitrumMessengerTransactor
	TransactOpts bind.TransactOpts
}

type L2ArbitrumMessengerRaw struct {
	Contract *L2ArbitrumMessenger
}

type L2ArbitrumMessengerCallerRaw struct {
	Contract *L2ArbitrumMessengerCaller
}

type L2ArbitrumMessengerTransactorRaw struct {
	Contract *L2ArbitrumMessengerTransactor
}

func NewL2ArbitrumMessenger(address common.Address, backend bind.ContractBackend) (*L2ArbitrumMessenger, error) {
	abi, err := abi.JSON(strings.NewReader(L2ArbitrumMessengerABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindL2ArbitrumMessenger(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &L2ArbitrumMessenger{address: address, abi: abi, L2ArbitrumMessengerCaller: L2ArbitrumMessengerCaller{contract: contract}, L2ArbitrumMessengerTransactor: L2ArbitrumMessengerTransactor{contract: contract}, L2ArbitrumMessengerFilterer: L2ArbitrumMessengerFilterer{contract: contract}}, nil
}

func NewL2ArbitrumMessengerCaller(address common.Address, caller bind.ContractCaller) (*L2ArbitrumMessengerCaller, error) {
	contract, err := bindL2ArbitrumMessenger(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &L2ArbitrumMessengerCaller{contract: contract}, nil
}

func NewL2ArbitrumMessengerTransactor(address common.Address, transactor bind.ContractTransactor) (*L2ArbitrumMessengerTransactor, error) {
	contract, err := bindL2ArbitrumMessenger(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &L2ArbitrumMessengerTransactor{contract: contract}, nil
}

func NewL2ArbitrumMessengerFilterer(address common.Address, filterer bind.ContractFilterer) (*L2ArbitrumMessengerFilterer, error) {
	contract, err := bindL2ArbitrumMessenger(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &L2ArbitrumMessengerFilterer{contract: contract}, nil
}

func bindL2ArbitrumMessenger(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := L2ArbitrumMessengerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_L2ArbitrumMessenger *L2ArbitrumMessengerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _L2ArbitrumMessenger.Contract.L2ArbitrumMessengerCaller.contract.Call(opts, result, method, params...)
}

func (_L2ArbitrumMessenger *L2ArbitrumMessengerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L2ArbitrumMessenger.Contract.L2ArbitrumMessengerTransactor.contract.Transfer(opts)
}

func (_L2ArbitrumMessenger *L2ArbitrumMessengerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _L2ArbitrumMessenger.Contract.L2ArbitrumMessengerTransactor.contract.Transact(opts, method, params...)
}

func (_L2ArbitrumMessenger *L2ArbitrumMessengerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _L2ArbitrumMessenger.Contract.contract.Call(opts, result, method, params...)
}

func (_L2ArbitrumMessenger *L2ArbitrumMessengerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L2ArbitrumMessenger.Contract.contract.Transfer(opts)
}

func (_L2ArbitrumMessenger *L2ArbitrumMessengerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _L2ArbitrumMessenger.Contract.contract.Transact(opts, method, params...)
}

type L2ArbitrumMessengerTxToL1Iterator struct {
	Event *L2ArbitrumMessengerTxToL1

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *L2ArbitrumMessengerTxToL1Iterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L2ArbitrumMessengerTxToL1)
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
		it.Event = new(L2ArbitrumMessengerTxToL1)
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

func (it *L2ArbitrumMessengerTxToL1Iterator) Error() error {
	return it.fail
}

func (it *L2ArbitrumMessengerTxToL1Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type L2ArbitrumMessengerTxToL1 struct {
	From common.Address
	To   common.Address
	Id   *big.Int
	Data []byte
	Raw  types.Log
}

func (_L2ArbitrumMessenger *L2ArbitrumMessengerFilterer) FilterTxToL1(opts *bind.FilterOpts, _from []common.Address, _to []common.Address, _id []*big.Int) (*L2ArbitrumMessengerTxToL1Iterator, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}
	var _idRule []interface{}
	for _, _idItem := range _id {
		_idRule = append(_idRule, _idItem)
	}

	logs, sub, err := _L2ArbitrumMessenger.contract.FilterLogs(opts, "TxToL1", _fromRule, _toRule, _idRule)
	if err != nil {
		return nil, err
	}
	return &L2ArbitrumMessengerTxToL1Iterator{contract: _L2ArbitrumMessenger.contract, event: "TxToL1", logs: logs, sub: sub}, nil
}

func (_L2ArbitrumMessenger *L2ArbitrumMessengerFilterer) WatchTxToL1(opts *bind.WatchOpts, sink chan<- *L2ArbitrumMessengerTxToL1, _from []common.Address, _to []common.Address, _id []*big.Int) (event.Subscription, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}
	var _idRule []interface{}
	for _, _idItem := range _id {
		_idRule = append(_idRule, _idItem)
	}

	logs, sub, err := _L2ArbitrumMessenger.contract.WatchLogs(opts, "TxToL1", _fromRule, _toRule, _idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(L2ArbitrumMessengerTxToL1)
				if err := _L2ArbitrumMessenger.contract.UnpackLog(event, "TxToL1", log); err != nil {
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

func (_L2ArbitrumMessenger *L2ArbitrumMessengerFilterer) ParseTxToL1(log types.Log) (*L2ArbitrumMessengerTxToL1, error) {
	event := new(L2ArbitrumMessengerTxToL1)
	if err := _L2ArbitrumMessenger.contract.UnpackLog(event, "TxToL1", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_L2ArbitrumMessenger *L2ArbitrumMessenger) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _L2ArbitrumMessenger.abi.Events["TxToL1"].ID:
		return _L2ArbitrumMessenger.ParseTxToL1(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (L2ArbitrumMessengerTxToL1) Topic() common.Hash {
	return common.HexToHash("0x2b986d32a0536b7e19baa48ab949fec7b903b7fad7730820b20632d100cc3a68")
}

func (_L2ArbitrumMessenger *L2ArbitrumMessenger) Address() common.Address {
	return _L2ArbitrumMessenger.address
}

type L2ArbitrumMessengerInterface interface {
	FilterTxToL1(opts *bind.FilterOpts, _from []common.Address, _to []common.Address, _id []*big.Int) (*L2ArbitrumMessengerTxToL1Iterator, error)

	WatchTxToL1(opts *bind.WatchOpts, sink chan<- *L2ArbitrumMessengerTxToL1, _from []common.Address, _to []common.Address, _id []*big.Int) (event.Subscription, error)

	ParseTxToL1(log types.Log) (*L2ArbitrumMessengerTxToL1, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
