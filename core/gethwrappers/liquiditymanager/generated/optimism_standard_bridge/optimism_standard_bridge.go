// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package optimism_standard_bridge

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

var OptimismStandardBridgeMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"localToken\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"remoteToken\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"ERC20BridgeFinalized\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_localToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_remoteToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_extraData\",\"type\":\"bytes\"}],\"name\":\"finalizeBridgeERC20\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

var OptimismStandardBridgeABI = OptimismStandardBridgeMetaData.ABI

type OptimismStandardBridge struct {
	address common.Address
	abi     abi.ABI
	OptimismStandardBridgeCaller
	OptimismStandardBridgeTransactor
	OptimismStandardBridgeFilterer
}

type OptimismStandardBridgeCaller struct {
	contract *bind.BoundContract
}

type OptimismStandardBridgeTransactor struct {
	contract *bind.BoundContract
}

type OptimismStandardBridgeFilterer struct {
	contract *bind.BoundContract
}

type OptimismStandardBridgeSession struct {
	Contract     *OptimismStandardBridge
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type OptimismStandardBridgeCallerSession struct {
	Contract *OptimismStandardBridgeCaller
	CallOpts bind.CallOpts
}

type OptimismStandardBridgeTransactorSession struct {
	Contract     *OptimismStandardBridgeTransactor
	TransactOpts bind.TransactOpts
}

type OptimismStandardBridgeRaw struct {
	Contract *OptimismStandardBridge
}

type OptimismStandardBridgeCallerRaw struct {
	Contract *OptimismStandardBridgeCaller
}

type OptimismStandardBridgeTransactorRaw struct {
	Contract *OptimismStandardBridgeTransactor
}

func NewOptimismStandardBridge(address common.Address, backend bind.ContractBackend) (*OptimismStandardBridge, error) {
	abi, err := abi.JSON(strings.NewReader(OptimismStandardBridgeABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindOptimismStandardBridge(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OptimismStandardBridge{address: address, abi: abi, OptimismStandardBridgeCaller: OptimismStandardBridgeCaller{contract: contract}, OptimismStandardBridgeTransactor: OptimismStandardBridgeTransactor{contract: contract}, OptimismStandardBridgeFilterer: OptimismStandardBridgeFilterer{contract: contract}}, nil
}

func NewOptimismStandardBridgeCaller(address common.Address, caller bind.ContractCaller) (*OptimismStandardBridgeCaller, error) {
	contract, err := bindOptimismStandardBridge(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OptimismStandardBridgeCaller{contract: contract}, nil
}

func NewOptimismStandardBridgeTransactor(address common.Address, transactor bind.ContractTransactor) (*OptimismStandardBridgeTransactor, error) {
	contract, err := bindOptimismStandardBridge(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OptimismStandardBridgeTransactor{contract: contract}, nil
}

func NewOptimismStandardBridgeFilterer(address common.Address, filterer bind.ContractFilterer) (*OptimismStandardBridgeFilterer, error) {
	contract, err := bindOptimismStandardBridge(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OptimismStandardBridgeFilterer{contract: contract}, nil
}

func bindOptimismStandardBridge(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OptimismStandardBridgeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_OptimismStandardBridge *OptimismStandardBridgeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OptimismStandardBridge.Contract.OptimismStandardBridgeCaller.contract.Call(opts, result, method, params...)
}

func (_OptimismStandardBridge *OptimismStandardBridgeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OptimismStandardBridge.Contract.OptimismStandardBridgeTransactor.contract.Transfer(opts)
}

func (_OptimismStandardBridge *OptimismStandardBridgeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OptimismStandardBridge.Contract.OptimismStandardBridgeTransactor.contract.Transact(opts, method, params...)
}

func (_OptimismStandardBridge *OptimismStandardBridgeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OptimismStandardBridge.Contract.contract.Call(opts, result, method, params...)
}

func (_OptimismStandardBridge *OptimismStandardBridgeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OptimismStandardBridge.Contract.contract.Transfer(opts)
}

func (_OptimismStandardBridge *OptimismStandardBridgeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OptimismStandardBridge.Contract.contract.Transact(opts, method, params...)
}

func (_OptimismStandardBridge *OptimismStandardBridgeTransactor) FinalizeBridgeERC20(opts *bind.TransactOpts, _localToken common.Address, _remoteToken common.Address, _from common.Address, _to common.Address, _amount *big.Int, _extraData []byte) (*types.Transaction, error) {
	return _OptimismStandardBridge.contract.Transact(opts, "finalizeBridgeERC20", _localToken, _remoteToken, _from, _to, _amount, _extraData)
}

func (_OptimismStandardBridge *OptimismStandardBridgeSession) FinalizeBridgeERC20(_localToken common.Address, _remoteToken common.Address, _from common.Address, _to common.Address, _amount *big.Int, _extraData []byte) (*types.Transaction, error) {
	return _OptimismStandardBridge.Contract.FinalizeBridgeERC20(&_OptimismStandardBridge.TransactOpts, _localToken, _remoteToken, _from, _to, _amount, _extraData)
}

func (_OptimismStandardBridge *OptimismStandardBridgeTransactorSession) FinalizeBridgeERC20(_localToken common.Address, _remoteToken common.Address, _from common.Address, _to common.Address, _amount *big.Int, _extraData []byte) (*types.Transaction, error) {
	return _OptimismStandardBridge.Contract.FinalizeBridgeERC20(&_OptimismStandardBridge.TransactOpts, _localToken, _remoteToken, _from, _to, _amount, _extraData)
}

type OptimismStandardBridgeERC20BridgeFinalizedIterator struct {
	Event *OptimismStandardBridgeERC20BridgeFinalized

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OptimismStandardBridgeERC20BridgeFinalizedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OptimismStandardBridgeERC20BridgeFinalized)
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
		it.Event = new(OptimismStandardBridgeERC20BridgeFinalized)
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

func (it *OptimismStandardBridgeERC20BridgeFinalizedIterator) Error() error {
	return it.fail
}

func (it *OptimismStandardBridgeERC20BridgeFinalizedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OptimismStandardBridgeERC20BridgeFinalized struct {
	LocalToken  common.Address
	RemoteToken common.Address
	From        common.Address
	To          common.Address
	Amount      *big.Int
	ExtraData   []byte
	Raw         types.Log
}

func (_OptimismStandardBridge *OptimismStandardBridgeFilterer) FilterERC20BridgeFinalized(opts *bind.FilterOpts, localToken []common.Address, remoteToken []common.Address, from []common.Address) (*OptimismStandardBridgeERC20BridgeFinalizedIterator, error) {

	var localTokenRule []interface{}
	for _, localTokenItem := range localToken {
		localTokenRule = append(localTokenRule, localTokenItem)
	}
	var remoteTokenRule []interface{}
	for _, remoteTokenItem := range remoteToken {
		remoteTokenRule = append(remoteTokenRule, remoteTokenItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _OptimismStandardBridge.contract.FilterLogs(opts, "ERC20BridgeFinalized", localTokenRule, remoteTokenRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &OptimismStandardBridgeERC20BridgeFinalizedIterator{contract: _OptimismStandardBridge.contract, event: "ERC20BridgeFinalized", logs: logs, sub: sub}, nil
}

func (_OptimismStandardBridge *OptimismStandardBridgeFilterer) WatchERC20BridgeFinalized(opts *bind.WatchOpts, sink chan<- *OptimismStandardBridgeERC20BridgeFinalized, localToken []common.Address, remoteToken []common.Address, from []common.Address) (event.Subscription, error) {

	var localTokenRule []interface{}
	for _, localTokenItem := range localToken {
		localTokenRule = append(localTokenRule, localTokenItem)
	}
	var remoteTokenRule []interface{}
	for _, remoteTokenItem := range remoteToken {
		remoteTokenRule = append(remoteTokenRule, remoteTokenItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _OptimismStandardBridge.contract.WatchLogs(opts, "ERC20BridgeFinalized", localTokenRule, remoteTokenRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OptimismStandardBridgeERC20BridgeFinalized)
				if err := _OptimismStandardBridge.contract.UnpackLog(event, "ERC20BridgeFinalized", log); err != nil {
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

func (_OptimismStandardBridge *OptimismStandardBridgeFilterer) ParseERC20BridgeFinalized(log types.Log) (*OptimismStandardBridgeERC20BridgeFinalized, error) {
	event := new(OptimismStandardBridgeERC20BridgeFinalized)
	if err := _OptimismStandardBridge.contract.UnpackLog(event, "ERC20BridgeFinalized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_OptimismStandardBridge *OptimismStandardBridge) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _OptimismStandardBridge.abi.Events["ERC20BridgeFinalized"].ID:
		return _OptimismStandardBridge.ParseERC20BridgeFinalized(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (OptimismStandardBridgeERC20BridgeFinalized) Topic() common.Hash {
	return common.HexToHash("0xd59c65b35445225835c83f50b6ede06a7be047d22e357073e250d9af537518cd")
}

func (_OptimismStandardBridge *OptimismStandardBridge) Address() common.Address {
	return _OptimismStandardBridge.address
}

type OptimismStandardBridgeInterface interface {
	FinalizeBridgeERC20(opts *bind.TransactOpts, _localToken common.Address, _remoteToken common.Address, _from common.Address, _to common.Address, _amount *big.Int, _extraData []byte) (*types.Transaction, error)

	FilterERC20BridgeFinalized(opts *bind.FilterOpts, localToken []common.Address, remoteToken []common.Address, from []common.Address) (*OptimismStandardBridgeERC20BridgeFinalizedIterator, error)

	WatchERC20BridgeFinalized(opts *bind.WatchOpts, sink chan<- *OptimismStandardBridgeERC20BridgeFinalized, localToken []common.Address, remoteToken []common.Address, from []common.Address) (event.Subscription, error)

	ParseERC20BridgeFinalized(log types.Log) (*OptimismStandardBridgeERC20BridgeFinalized, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
