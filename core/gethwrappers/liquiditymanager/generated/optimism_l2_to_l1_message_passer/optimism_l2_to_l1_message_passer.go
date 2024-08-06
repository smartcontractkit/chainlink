// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package optimism_l2_to_l1_message_passer

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

var OptimismL2ToL1MessagePasserMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"withdrawalHash\",\"type\":\"bytes32\"}],\"name\":\"MessagePassed\",\"type\":\"event\"}]",
}

var OptimismL2ToL1MessagePasserABI = OptimismL2ToL1MessagePasserMetaData.ABI

type OptimismL2ToL1MessagePasser struct {
	address common.Address
	abi     abi.ABI
	OptimismL2ToL1MessagePasserCaller
	OptimismL2ToL1MessagePasserTransactor
	OptimismL2ToL1MessagePasserFilterer
}

type OptimismL2ToL1MessagePasserCaller struct {
	contract *bind.BoundContract
}

type OptimismL2ToL1MessagePasserTransactor struct {
	contract *bind.BoundContract
}

type OptimismL2ToL1MessagePasserFilterer struct {
	contract *bind.BoundContract
}

type OptimismL2ToL1MessagePasserSession struct {
	Contract     *OptimismL2ToL1MessagePasser
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type OptimismL2ToL1MessagePasserCallerSession struct {
	Contract *OptimismL2ToL1MessagePasserCaller
	CallOpts bind.CallOpts
}

type OptimismL2ToL1MessagePasserTransactorSession struct {
	Contract     *OptimismL2ToL1MessagePasserTransactor
	TransactOpts bind.TransactOpts
}

type OptimismL2ToL1MessagePasserRaw struct {
	Contract *OptimismL2ToL1MessagePasser
}

type OptimismL2ToL1MessagePasserCallerRaw struct {
	Contract *OptimismL2ToL1MessagePasserCaller
}

type OptimismL2ToL1MessagePasserTransactorRaw struct {
	Contract *OptimismL2ToL1MessagePasserTransactor
}

func NewOptimismL2ToL1MessagePasser(address common.Address, backend bind.ContractBackend) (*OptimismL2ToL1MessagePasser, error) {
	abi, err := abi.JSON(strings.NewReader(OptimismL2ToL1MessagePasserABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindOptimismL2ToL1MessagePasser(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OptimismL2ToL1MessagePasser{address: address, abi: abi, OptimismL2ToL1MessagePasserCaller: OptimismL2ToL1MessagePasserCaller{contract: contract}, OptimismL2ToL1MessagePasserTransactor: OptimismL2ToL1MessagePasserTransactor{contract: contract}, OptimismL2ToL1MessagePasserFilterer: OptimismL2ToL1MessagePasserFilterer{contract: contract}}, nil
}

func NewOptimismL2ToL1MessagePasserCaller(address common.Address, caller bind.ContractCaller) (*OptimismL2ToL1MessagePasserCaller, error) {
	contract, err := bindOptimismL2ToL1MessagePasser(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OptimismL2ToL1MessagePasserCaller{contract: contract}, nil
}

func NewOptimismL2ToL1MessagePasserTransactor(address common.Address, transactor bind.ContractTransactor) (*OptimismL2ToL1MessagePasserTransactor, error) {
	contract, err := bindOptimismL2ToL1MessagePasser(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OptimismL2ToL1MessagePasserTransactor{contract: contract}, nil
}

func NewOptimismL2ToL1MessagePasserFilterer(address common.Address, filterer bind.ContractFilterer) (*OptimismL2ToL1MessagePasserFilterer, error) {
	contract, err := bindOptimismL2ToL1MessagePasser(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OptimismL2ToL1MessagePasserFilterer{contract: contract}, nil
}

func bindOptimismL2ToL1MessagePasser(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OptimismL2ToL1MessagePasserMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_OptimismL2ToL1MessagePasser *OptimismL2ToL1MessagePasserRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OptimismL2ToL1MessagePasser.Contract.OptimismL2ToL1MessagePasserCaller.contract.Call(opts, result, method, params...)
}

func (_OptimismL2ToL1MessagePasser *OptimismL2ToL1MessagePasserRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OptimismL2ToL1MessagePasser.Contract.OptimismL2ToL1MessagePasserTransactor.contract.Transfer(opts)
}

func (_OptimismL2ToL1MessagePasser *OptimismL2ToL1MessagePasserRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OptimismL2ToL1MessagePasser.Contract.OptimismL2ToL1MessagePasserTransactor.contract.Transact(opts, method, params...)
}

func (_OptimismL2ToL1MessagePasser *OptimismL2ToL1MessagePasserCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OptimismL2ToL1MessagePasser.Contract.contract.Call(opts, result, method, params...)
}

func (_OptimismL2ToL1MessagePasser *OptimismL2ToL1MessagePasserTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OptimismL2ToL1MessagePasser.Contract.contract.Transfer(opts)
}

func (_OptimismL2ToL1MessagePasser *OptimismL2ToL1MessagePasserTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OptimismL2ToL1MessagePasser.Contract.contract.Transact(opts, method, params...)
}

type OptimismL2ToL1MessagePasserMessagePassedIterator struct {
	Event *OptimismL2ToL1MessagePasserMessagePassed

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OptimismL2ToL1MessagePasserMessagePassedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OptimismL2ToL1MessagePasserMessagePassed)
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
		it.Event = new(OptimismL2ToL1MessagePasserMessagePassed)
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

func (it *OptimismL2ToL1MessagePasserMessagePassedIterator) Error() error {
	return it.fail
}

func (it *OptimismL2ToL1MessagePasserMessagePassedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OptimismL2ToL1MessagePasserMessagePassed struct {
	Nonce          *big.Int
	Sender         common.Address
	Target         common.Address
	Value          *big.Int
	GasLimit       *big.Int
	Data           []byte
	WithdrawalHash [32]byte
	Raw            types.Log
}

func (_OptimismL2ToL1MessagePasser *OptimismL2ToL1MessagePasserFilterer) FilterMessagePassed(opts *bind.FilterOpts, nonce []*big.Int, sender []common.Address, target []common.Address) (*OptimismL2ToL1MessagePasserMessagePassedIterator, error) {

	var nonceRule []interface{}
	for _, nonceItem := range nonce {
		nonceRule = append(nonceRule, nonceItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var targetRule []interface{}
	for _, targetItem := range target {
		targetRule = append(targetRule, targetItem)
	}

	logs, sub, err := _OptimismL2ToL1MessagePasser.contract.FilterLogs(opts, "MessagePassed", nonceRule, senderRule, targetRule)
	if err != nil {
		return nil, err
	}
	return &OptimismL2ToL1MessagePasserMessagePassedIterator{contract: _OptimismL2ToL1MessagePasser.contract, event: "MessagePassed", logs: logs, sub: sub}, nil
}

func (_OptimismL2ToL1MessagePasser *OptimismL2ToL1MessagePasserFilterer) WatchMessagePassed(opts *bind.WatchOpts, sink chan<- *OptimismL2ToL1MessagePasserMessagePassed, nonce []*big.Int, sender []common.Address, target []common.Address) (event.Subscription, error) {

	var nonceRule []interface{}
	for _, nonceItem := range nonce {
		nonceRule = append(nonceRule, nonceItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var targetRule []interface{}
	for _, targetItem := range target {
		targetRule = append(targetRule, targetItem)
	}

	logs, sub, err := _OptimismL2ToL1MessagePasser.contract.WatchLogs(opts, "MessagePassed", nonceRule, senderRule, targetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OptimismL2ToL1MessagePasserMessagePassed)
				if err := _OptimismL2ToL1MessagePasser.contract.UnpackLog(event, "MessagePassed", log); err != nil {
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

func (_OptimismL2ToL1MessagePasser *OptimismL2ToL1MessagePasserFilterer) ParseMessagePassed(log types.Log) (*OptimismL2ToL1MessagePasserMessagePassed, error) {
	event := new(OptimismL2ToL1MessagePasserMessagePassed)
	if err := _OptimismL2ToL1MessagePasser.contract.UnpackLog(event, "MessagePassed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_OptimismL2ToL1MessagePasser *OptimismL2ToL1MessagePasser) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _OptimismL2ToL1MessagePasser.abi.Events["MessagePassed"].ID:
		return _OptimismL2ToL1MessagePasser.ParseMessagePassed(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (OptimismL2ToL1MessagePasserMessagePassed) Topic() common.Hash {
	return common.HexToHash("0x02a52367d10742d8032712c1bb8e0144ff1ec5ffda1ed7d70bb05a2744955054")
}

func (_OptimismL2ToL1MessagePasser *OptimismL2ToL1MessagePasser) Address() common.Address {
	return _OptimismL2ToL1MessagePasser.address
}

type OptimismL2ToL1MessagePasserInterface interface {
	FilterMessagePassed(opts *bind.FilterOpts, nonce []*big.Int, sender []common.Address, target []common.Address) (*OptimismL2ToL1MessagePasserMessagePassedIterator, error)

	WatchMessagePassed(opts *bind.WatchOpts, sink chan<- *OptimismL2ToL1MessagePasserMessagePassed, nonce []*big.Int, sender []common.Address, target []common.Address) (event.Subscription, error)

	ParseMessagePassed(log types.Log) (*OptimismL2ToL1MessagePasserMessagePassed, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
