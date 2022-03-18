// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package cron_upkeep_factory_wrapper

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
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated"
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

var CronUpkeepFactoryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"upkeep\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"NewCronUpkeepCreated\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"cronDelegateAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"newCronUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

var CronUpkeepFactoryABI = CronUpkeepFactoryMetaData.ABI

type CronUpkeepFactory struct {
	address common.Address
	abi     abi.ABI
	CronUpkeepFactoryCaller
	CronUpkeepFactoryTransactor
	CronUpkeepFactoryFilterer
}

type CronUpkeepFactoryCaller struct {
	contract *bind.BoundContract
}

type CronUpkeepFactoryTransactor struct {
	contract *bind.BoundContract
}

type CronUpkeepFactoryFilterer struct {
	contract *bind.BoundContract
}

type CronUpkeepFactorySession struct {
	Contract     *CronUpkeepFactory
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type CronUpkeepFactoryCallerSession struct {
	Contract *CronUpkeepFactoryCaller
	CallOpts bind.CallOpts
}

type CronUpkeepFactoryTransactorSession struct {
	Contract     *CronUpkeepFactoryTransactor
	TransactOpts bind.TransactOpts
}

type CronUpkeepFactoryRaw struct {
	Contract *CronUpkeepFactory
}

type CronUpkeepFactoryCallerRaw struct {
	Contract *CronUpkeepFactoryCaller
}

type CronUpkeepFactoryTransactorRaw struct {
	Contract *CronUpkeepFactoryTransactor
}

func NewCronUpkeepFactory(address common.Address, backend bind.ContractBackend) (*CronUpkeepFactory, error) {
	abi, err := abi.JSON(strings.NewReader(CronUpkeepFactoryABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindCronUpkeepFactory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CronUpkeepFactory{address: address, abi: abi, CronUpkeepFactoryCaller: CronUpkeepFactoryCaller{contract: contract}, CronUpkeepFactoryTransactor: CronUpkeepFactoryTransactor{contract: contract}, CronUpkeepFactoryFilterer: CronUpkeepFactoryFilterer{contract: contract}}, nil
}

func NewCronUpkeepFactoryCaller(address common.Address, caller bind.ContractCaller) (*CronUpkeepFactoryCaller, error) {
	contract, err := bindCronUpkeepFactory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CronUpkeepFactoryCaller{contract: contract}, nil
}

func NewCronUpkeepFactoryTransactor(address common.Address, transactor bind.ContractTransactor) (*CronUpkeepFactoryTransactor, error) {
	contract, err := bindCronUpkeepFactory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CronUpkeepFactoryTransactor{contract: contract}, nil
}

func NewCronUpkeepFactoryFilterer(address common.Address, filterer bind.ContractFilterer) (*CronUpkeepFactoryFilterer, error) {
	contract, err := bindCronUpkeepFactory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CronUpkeepFactoryFilterer{contract: contract}, nil
}

func bindCronUpkeepFactory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(CronUpkeepFactoryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_CronUpkeepFactory *CronUpkeepFactoryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CronUpkeepFactory.Contract.CronUpkeepFactoryCaller.contract.Call(opts, result, method, params...)
}

func (_CronUpkeepFactory *CronUpkeepFactoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CronUpkeepFactory.Contract.CronUpkeepFactoryTransactor.contract.Transfer(opts)
}

func (_CronUpkeepFactory *CronUpkeepFactoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CronUpkeepFactory.Contract.CronUpkeepFactoryTransactor.contract.Transact(opts, method, params...)
}

func (_CronUpkeepFactory *CronUpkeepFactoryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CronUpkeepFactory.Contract.contract.Call(opts, result, method, params...)
}

func (_CronUpkeepFactory *CronUpkeepFactoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CronUpkeepFactory.Contract.contract.Transfer(opts)
}

func (_CronUpkeepFactory *CronUpkeepFactoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CronUpkeepFactory.Contract.contract.Transact(opts, method, params...)
}

func (_CronUpkeepFactory *CronUpkeepFactoryCaller) CronDelegateAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CronUpkeepFactory.contract.Call(opts, &out, "cronDelegateAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_CronUpkeepFactory *CronUpkeepFactorySession) CronDelegateAddress() (common.Address, error) {
	return _CronUpkeepFactory.Contract.CronDelegateAddress(&_CronUpkeepFactory.CallOpts)
}

func (_CronUpkeepFactory *CronUpkeepFactoryCallerSession) CronDelegateAddress() (common.Address, error) {
	return _CronUpkeepFactory.Contract.CronDelegateAddress(&_CronUpkeepFactory.CallOpts)
}

func (_CronUpkeepFactory *CronUpkeepFactoryTransactor) NewCronUpkeep(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CronUpkeepFactory.contract.Transact(opts, "newCronUpkeep")
}

func (_CronUpkeepFactory *CronUpkeepFactorySession) NewCronUpkeep() (*types.Transaction, error) {
	return _CronUpkeepFactory.Contract.NewCronUpkeep(&_CronUpkeepFactory.TransactOpts)
}

func (_CronUpkeepFactory *CronUpkeepFactoryTransactorSession) NewCronUpkeep() (*types.Transaction, error) {
	return _CronUpkeepFactory.Contract.NewCronUpkeep(&_CronUpkeepFactory.TransactOpts)
}

type CronUpkeepFactoryNewCronUpkeepCreatedIterator struct {
	Event *CronUpkeepFactoryNewCronUpkeepCreated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CronUpkeepFactoryNewCronUpkeepCreatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CronUpkeepFactoryNewCronUpkeepCreated)
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
		it.Event = new(CronUpkeepFactoryNewCronUpkeepCreated)
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

func (it *CronUpkeepFactoryNewCronUpkeepCreatedIterator) Error() error {
	return it.fail
}

func (it *CronUpkeepFactoryNewCronUpkeepCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CronUpkeepFactoryNewCronUpkeepCreated struct {
	Upkeep common.Address
	Owner  common.Address
	Raw    types.Log
}

func (_CronUpkeepFactory *CronUpkeepFactoryFilterer) FilterNewCronUpkeepCreated(opts *bind.FilterOpts) (*CronUpkeepFactoryNewCronUpkeepCreatedIterator, error) {

	logs, sub, err := _CronUpkeepFactory.contract.FilterLogs(opts, "NewCronUpkeepCreated")
	if err != nil {
		return nil, err
	}
	return &CronUpkeepFactoryNewCronUpkeepCreatedIterator{contract: _CronUpkeepFactory.contract, event: "NewCronUpkeepCreated", logs: logs, sub: sub}, nil
}

func (_CronUpkeepFactory *CronUpkeepFactoryFilterer) WatchNewCronUpkeepCreated(opts *bind.WatchOpts, sink chan<- *CronUpkeepFactoryNewCronUpkeepCreated) (event.Subscription, error) {

	logs, sub, err := _CronUpkeepFactory.contract.WatchLogs(opts, "NewCronUpkeepCreated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CronUpkeepFactoryNewCronUpkeepCreated)
				if err := _CronUpkeepFactory.contract.UnpackLog(event, "NewCronUpkeepCreated", log); err != nil {
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

func (_CronUpkeepFactory *CronUpkeepFactoryFilterer) ParseNewCronUpkeepCreated(log types.Log) (*CronUpkeepFactoryNewCronUpkeepCreated, error) {
	event := new(CronUpkeepFactoryNewCronUpkeepCreated)
	if err := _CronUpkeepFactory.contract.UnpackLog(event, "NewCronUpkeepCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_CronUpkeepFactory *CronUpkeepFactory) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _CronUpkeepFactory.abi.Events["NewCronUpkeepCreated"].ID:
		return _CronUpkeepFactory.ParseNewCronUpkeepCreated(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (CronUpkeepFactoryNewCronUpkeepCreated) Topic() common.Hash {
	return common.HexToHash("0x959d571686b1c9343b61bdc3c0459760cb9695fcd4c4c64845e3b2cdd6865ced")
}

func (_CronUpkeepFactory *CronUpkeepFactory) Address() common.Address {
	return _CronUpkeepFactory.address
}

type CronUpkeepFactoryInterface interface {
	CronDelegateAddress(opts *bind.CallOpts) (common.Address, error)

	NewCronUpkeep(opts *bind.TransactOpts) (*types.Transaction, error)

	FilterNewCronUpkeepCreated(opts *bind.FilterOpts) (*CronUpkeepFactoryNewCronUpkeepCreatedIterator, error)

	WatchNewCronUpkeepCreated(opts *bind.WatchOpts, sink chan<- *CronUpkeepFactoryNewCronUpkeepCreated) (event.Subscription, error)

	ParseNewCronUpkeepCreated(log types.Log) (*CronUpkeepFactoryNewCronUpkeepCreated, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
