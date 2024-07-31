// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package arbitrum_gateway_router

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

var ArbitrumGatewayRouterMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newDefaultGateway\",\"type\":\"address\"}],\"name\":\"DefaultGatewayUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"l1Token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"gateway\",\"type\":\"address\"}],\"name\":\"GatewaySet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_userFrom\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_userTo\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"gateway\",\"type\":\"address\"}],\"name\":\"TransferRouted\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"l1ERC20\",\"type\":\"address\"}],\"name\":\"calculateL2TokenAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"defaultGateway\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"gateway\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"finalizeInboundTransfer\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"getGateway\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"gateway\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"getOutboundCalldata\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_maxGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_gasPriceBid\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"outboundTransfer\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
}

var ArbitrumGatewayRouterABI = ArbitrumGatewayRouterMetaData.ABI

type ArbitrumGatewayRouter struct {
	address common.Address
	abi     abi.ABI
	ArbitrumGatewayRouterCaller
	ArbitrumGatewayRouterTransactor
	ArbitrumGatewayRouterFilterer
}

type ArbitrumGatewayRouterCaller struct {
	contract *bind.BoundContract
}

type ArbitrumGatewayRouterTransactor struct {
	contract *bind.BoundContract
}

type ArbitrumGatewayRouterFilterer struct {
	contract *bind.BoundContract
}

type ArbitrumGatewayRouterSession struct {
	Contract     *ArbitrumGatewayRouter
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type ArbitrumGatewayRouterCallerSession struct {
	Contract *ArbitrumGatewayRouterCaller
	CallOpts bind.CallOpts
}

type ArbitrumGatewayRouterTransactorSession struct {
	Contract     *ArbitrumGatewayRouterTransactor
	TransactOpts bind.TransactOpts
}

type ArbitrumGatewayRouterRaw struct {
	Contract *ArbitrumGatewayRouter
}

type ArbitrumGatewayRouterCallerRaw struct {
	Contract *ArbitrumGatewayRouterCaller
}

type ArbitrumGatewayRouterTransactorRaw struct {
	Contract *ArbitrumGatewayRouterTransactor
}

func NewArbitrumGatewayRouter(address common.Address, backend bind.ContractBackend) (*ArbitrumGatewayRouter, error) {
	abi, err := abi.JSON(strings.NewReader(ArbitrumGatewayRouterABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindArbitrumGatewayRouter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ArbitrumGatewayRouter{address: address, abi: abi, ArbitrumGatewayRouterCaller: ArbitrumGatewayRouterCaller{contract: contract}, ArbitrumGatewayRouterTransactor: ArbitrumGatewayRouterTransactor{contract: contract}, ArbitrumGatewayRouterFilterer: ArbitrumGatewayRouterFilterer{contract: contract}}, nil
}

func NewArbitrumGatewayRouterCaller(address common.Address, caller bind.ContractCaller) (*ArbitrumGatewayRouterCaller, error) {
	contract, err := bindArbitrumGatewayRouter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ArbitrumGatewayRouterCaller{contract: contract}, nil
}

func NewArbitrumGatewayRouterTransactor(address common.Address, transactor bind.ContractTransactor) (*ArbitrumGatewayRouterTransactor, error) {
	contract, err := bindArbitrumGatewayRouter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ArbitrumGatewayRouterTransactor{contract: contract}, nil
}

func NewArbitrumGatewayRouterFilterer(address common.Address, filterer bind.ContractFilterer) (*ArbitrumGatewayRouterFilterer, error) {
	contract, err := bindArbitrumGatewayRouter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ArbitrumGatewayRouterFilterer{contract: contract}, nil
}

func bindArbitrumGatewayRouter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ArbitrumGatewayRouterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_ArbitrumGatewayRouter *ArbitrumGatewayRouterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ArbitrumGatewayRouter.Contract.ArbitrumGatewayRouterCaller.contract.Call(opts, result, method, params...)
}

func (_ArbitrumGatewayRouter *ArbitrumGatewayRouterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ArbitrumGatewayRouter.Contract.ArbitrumGatewayRouterTransactor.contract.Transfer(opts)
}

func (_ArbitrumGatewayRouter *ArbitrumGatewayRouterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ArbitrumGatewayRouter.Contract.ArbitrumGatewayRouterTransactor.contract.Transact(opts, method, params...)
}

func (_ArbitrumGatewayRouter *ArbitrumGatewayRouterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ArbitrumGatewayRouter.Contract.contract.Call(opts, result, method, params...)
}

func (_ArbitrumGatewayRouter *ArbitrumGatewayRouterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ArbitrumGatewayRouter.Contract.contract.Transfer(opts)
}

func (_ArbitrumGatewayRouter *ArbitrumGatewayRouterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ArbitrumGatewayRouter.Contract.contract.Transact(opts, method, params...)
}

func (_ArbitrumGatewayRouter *ArbitrumGatewayRouterCaller) CalculateL2TokenAddress(opts *bind.CallOpts, l1ERC20 common.Address) (common.Address, error) {
	var out []interface{}
	err := _ArbitrumGatewayRouter.contract.Call(opts, &out, "calculateL2TokenAddress", l1ERC20)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_ArbitrumGatewayRouter *ArbitrumGatewayRouterSession) CalculateL2TokenAddress(l1ERC20 common.Address) (common.Address, error) {
	return _ArbitrumGatewayRouter.Contract.CalculateL2TokenAddress(&_ArbitrumGatewayRouter.CallOpts, l1ERC20)
}

func (_ArbitrumGatewayRouter *ArbitrumGatewayRouterCallerSession) CalculateL2TokenAddress(l1ERC20 common.Address) (common.Address, error) {
	return _ArbitrumGatewayRouter.Contract.CalculateL2TokenAddress(&_ArbitrumGatewayRouter.CallOpts, l1ERC20)
}

func (_ArbitrumGatewayRouter *ArbitrumGatewayRouterCaller) DefaultGateway(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ArbitrumGatewayRouter.contract.Call(opts, &out, "defaultGateway")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_ArbitrumGatewayRouter *ArbitrumGatewayRouterSession) DefaultGateway() (common.Address, error) {
	return _ArbitrumGatewayRouter.Contract.DefaultGateway(&_ArbitrumGatewayRouter.CallOpts)
}

func (_ArbitrumGatewayRouter *ArbitrumGatewayRouterCallerSession) DefaultGateway() (common.Address, error) {
	return _ArbitrumGatewayRouter.Contract.DefaultGateway(&_ArbitrumGatewayRouter.CallOpts)
}

func (_ArbitrumGatewayRouter *ArbitrumGatewayRouterCaller) GetGateway(opts *bind.CallOpts, _token common.Address) (common.Address, error) {
	var out []interface{}
	err := _ArbitrumGatewayRouter.contract.Call(opts, &out, "getGateway", _token)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_ArbitrumGatewayRouter *ArbitrumGatewayRouterSession) GetGateway(_token common.Address) (common.Address, error) {
	return _ArbitrumGatewayRouter.Contract.GetGateway(&_ArbitrumGatewayRouter.CallOpts, _token)
}

func (_ArbitrumGatewayRouter *ArbitrumGatewayRouterCallerSession) GetGateway(_token common.Address) (common.Address, error) {
	return _ArbitrumGatewayRouter.Contract.GetGateway(&_ArbitrumGatewayRouter.CallOpts, _token)
}

func (_ArbitrumGatewayRouter *ArbitrumGatewayRouterCaller) GetOutboundCalldata(opts *bind.CallOpts, _token common.Address, _from common.Address, _to common.Address, _amount *big.Int, _data []byte) ([]byte, error) {
	var out []interface{}
	err := _ArbitrumGatewayRouter.contract.Call(opts, &out, "getOutboundCalldata", _token, _from, _to, _amount, _data)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_ArbitrumGatewayRouter *ArbitrumGatewayRouterSession) GetOutboundCalldata(_token common.Address, _from common.Address, _to common.Address, _amount *big.Int, _data []byte) ([]byte, error) {
	return _ArbitrumGatewayRouter.Contract.GetOutboundCalldata(&_ArbitrumGatewayRouter.CallOpts, _token, _from, _to, _amount, _data)
}

func (_ArbitrumGatewayRouter *ArbitrumGatewayRouterCallerSession) GetOutboundCalldata(_token common.Address, _from common.Address, _to common.Address, _amount *big.Int, _data []byte) ([]byte, error) {
	return _ArbitrumGatewayRouter.Contract.GetOutboundCalldata(&_ArbitrumGatewayRouter.CallOpts, _token, _from, _to, _amount, _data)
}

func (_ArbitrumGatewayRouter *ArbitrumGatewayRouterTransactor) FinalizeInboundTransfer(opts *bind.TransactOpts, _token common.Address, _from common.Address, _to common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _ArbitrumGatewayRouter.contract.Transact(opts, "finalizeInboundTransfer", _token, _from, _to, _amount, _data)
}

func (_ArbitrumGatewayRouter *ArbitrumGatewayRouterSession) FinalizeInboundTransfer(_token common.Address, _from common.Address, _to common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _ArbitrumGatewayRouter.Contract.FinalizeInboundTransfer(&_ArbitrumGatewayRouter.TransactOpts, _token, _from, _to, _amount, _data)
}

func (_ArbitrumGatewayRouter *ArbitrumGatewayRouterTransactorSession) FinalizeInboundTransfer(_token common.Address, _from common.Address, _to common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _ArbitrumGatewayRouter.Contract.FinalizeInboundTransfer(&_ArbitrumGatewayRouter.TransactOpts, _token, _from, _to, _amount, _data)
}

func (_ArbitrumGatewayRouter *ArbitrumGatewayRouterTransactor) OutboundTransfer(opts *bind.TransactOpts, _token common.Address, _to common.Address, _amount *big.Int, _maxGas *big.Int, _gasPriceBid *big.Int, _data []byte) (*types.Transaction, error) {
	return _ArbitrumGatewayRouter.contract.Transact(opts, "outboundTransfer", _token, _to, _amount, _maxGas, _gasPriceBid, _data)
}

func (_ArbitrumGatewayRouter *ArbitrumGatewayRouterSession) OutboundTransfer(_token common.Address, _to common.Address, _amount *big.Int, _maxGas *big.Int, _gasPriceBid *big.Int, _data []byte) (*types.Transaction, error) {
	return _ArbitrumGatewayRouter.Contract.OutboundTransfer(&_ArbitrumGatewayRouter.TransactOpts, _token, _to, _amount, _maxGas, _gasPriceBid, _data)
}

func (_ArbitrumGatewayRouter *ArbitrumGatewayRouterTransactorSession) OutboundTransfer(_token common.Address, _to common.Address, _amount *big.Int, _maxGas *big.Int, _gasPriceBid *big.Int, _data []byte) (*types.Transaction, error) {
	return _ArbitrumGatewayRouter.Contract.OutboundTransfer(&_ArbitrumGatewayRouter.TransactOpts, _token, _to, _amount, _maxGas, _gasPriceBid, _data)
}

type ArbitrumGatewayRouterDefaultGatewayUpdatedIterator struct {
	Event *ArbitrumGatewayRouterDefaultGatewayUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ArbitrumGatewayRouterDefaultGatewayUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ArbitrumGatewayRouterDefaultGatewayUpdated)
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
		it.Event = new(ArbitrumGatewayRouterDefaultGatewayUpdated)
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

func (it *ArbitrumGatewayRouterDefaultGatewayUpdatedIterator) Error() error {
	return it.fail
}

func (it *ArbitrumGatewayRouterDefaultGatewayUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ArbitrumGatewayRouterDefaultGatewayUpdated struct {
	NewDefaultGateway common.Address
	Raw               types.Log
}

func (_ArbitrumGatewayRouter *ArbitrumGatewayRouterFilterer) FilterDefaultGatewayUpdated(opts *bind.FilterOpts) (*ArbitrumGatewayRouterDefaultGatewayUpdatedIterator, error) {

	logs, sub, err := _ArbitrumGatewayRouter.contract.FilterLogs(opts, "DefaultGatewayUpdated")
	if err != nil {
		return nil, err
	}
	return &ArbitrumGatewayRouterDefaultGatewayUpdatedIterator{contract: _ArbitrumGatewayRouter.contract, event: "DefaultGatewayUpdated", logs: logs, sub: sub}, nil
}

func (_ArbitrumGatewayRouter *ArbitrumGatewayRouterFilterer) WatchDefaultGatewayUpdated(opts *bind.WatchOpts, sink chan<- *ArbitrumGatewayRouterDefaultGatewayUpdated) (event.Subscription, error) {

	logs, sub, err := _ArbitrumGatewayRouter.contract.WatchLogs(opts, "DefaultGatewayUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ArbitrumGatewayRouterDefaultGatewayUpdated)
				if err := _ArbitrumGatewayRouter.contract.UnpackLog(event, "DefaultGatewayUpdated", log); err != nil {
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

func (_ArbitrumGatewayRouter *ArbitrumGatewayRouterFilterer) ParseDefaultGatewayUpdated(log types.Log) (*ArbitrumGatewayRouterDefaultGatewayUpdated, error) {
	event := new(ArbitrumGatewayRouterDefaultGatewayUpdated)
	if err := _ArbitrumGatewayRouter.contract.UnpackLog(event, "DefaultGatewayUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ArbitrumGatewayRouterGatewaySetIterator struct {
	Event *ArbitrumGatewayRouterGatewaySet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ArbitrumGatewayRouterGatewaySetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ArbitrumGatewayRouterGatewaySet)
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
		it.Event = new(ArbitrumGatewayRouterGatewaySet)
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

func (it *ArbitrumGatewayRouterGatewaySetIterator) Error() error {
	return it.fail
}

func (it *ArbitrumGatewayRouterGatewaySetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ArbitrumGatewayRouterGatewaySet struct {
	L1Token common.Address
	Gateway common.Address
	Raw     types.Log
}

func (_ArbitrumGatewayRouter *ArbitrumGatewayRouterFilterer) FilterGatewaySet(opts *bind.FilterOpts, l1Token []common.Address, gateway []common.Address) (*ArbitrumGatewayRouterGatewaySetIterator, error) {

	var l1TokenRule []interface{}
	for _, l1TokenItem := range l1Token {
		l1TokenRule = append(l1TokenRule, l1TokenItem)
	}
	var gatewayRule []interface{}
	for _, gatewayItem := range gateway {
		gatewayRule = append(gatewayRule, gatewayItem)
	}

	logs, sub, err := _ArbitrumGatewayRouter.contract.FilterLogs(opts, "GatewaySet", l1TokenRule, gatewayRule)
	if err != nil {
		return nil, err
	}
	return &ArbitrumGatewayRouterGatewaySetIterator{contract: _ArbitrumGatewayRouter.contract, event: "GatewaySet", logs: logs, sub: sub}, nil
}

func (_ArbitrumGatewayRouter *ArbitrumGatewayRouterFilterer) WatchGatewaySet(opts *bind.WatchOpts, sink chan<- *ArbitrumGatewayRouterGatewaySet, l1Token []common.Address, gateway []common.Address) (event.Subscription, error) {

	var l1TokenRule []interface{}
	for _, l1TokenItem := range l1Token {
		l1TokenRule = append(l1TokenRule, l1TokenItem)
	}
	var gatewayRule []interface{}
	for _, gatewayItem := range gateway {
		gatewayRule = append(gatewayRule, gatewayItem)
	}

	logs, sub, err := _ArbitrumGatewayRouter.contract.WatchLogs(opts, "GatewaySet", l1TokenRule, gatewayRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ArbitrumGatewayRouterGatewaySet)
				if err := _ArbitrumGatewayRouter.contract.UnpackLog(event, "GatewaySet", log); err != nil {
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

func (_ArbitrumGatewayRouter *ArbitrumGatewayRouterFilterer) ParseGatewaySet(log types.Log) (*ArbitrumGatewayRouterGatewaySet, error) {
	event := new(ArbitrumGatewayRouterGatewaySet)
	if err := _ArbitrumGatewayRouter.contract.UnpackLog(event, "GatewaySet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type ArbitrumGatewayRouterTransferRoutedIterator struct {
	Event *ArbitrumGatewayRouterTransferRouted

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *ArbitrumGatewayRouterTransferRoutedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ArbitrumGatewayRouterTransferRouted)
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
		it.Event = new(ArbitrumGatewayRouterTransferRouted)
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

func (it *ArbitrumGatewayRouterTransferRoutedIterator) Error() error {
	return it.fail
}

func (it *ArbitrumGatewayRouterTransferRoutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type ArbitrumGatewayRouterTransferRouted struct {
	Token    common.Address
	UserFrom common.Address
	UserTo   common.Address
	Gateway  common.Address
	Raw      types.Log
}

func (_ArbitrumGatewayRouter *ArbitrumGatewayRouterFilterer) FilterTransferRouted(opts *bind.FilterOpts, token []common.Address, _userFrom []common.Address, _userTo []common.Address) (*ArbitrumGatewayRouterTransferRoutedIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var _userFromRule []interface{}
	for _, _userFromItem := range _userFrom {
		_userFromRule = append(_userFromRule, _userFromItem)
	}
	var _userToRule []interface{}
	for _, _userToItem := range _userTo {
		_userToRule = append(_userToRule, _userToItem)
	}

	logs, sub, err := _ArbitrumGatewayRouter.contract.FilterLogs(opts, "TransferRouted", tokenRule, _userFromRule, _userToRule)
	if err != nil {
		return nil, err
	}
	return &ArbitrumGatewayRouterTransferRoutedIterator{contract: _ArbitrumGatewayRouter.contract, event: "TransferRouted", logs: logs, sub: sub}, nil
}

func (_ArbitrumGatewayRouter *ArbitrumGatewayRouterFilterer) WatchTransferRouted(opts *bind.WatchOpts, sink chan<- *ArbitrumGatewayRouterTransferRouted, token []common.Address, _userFrom []common.Address, _userTo []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var _userFromRule []interface{}
	for _, _userFromItem := range _userFrom {
		_userFromRule = append(_userFromRule, _userFromItem)
	}
	var _userToRule []interface{}
	for _, _userToItem := range _userTo {
		_userToRule = append(_userToRule, _userToItem)
	}

	logs, sub, err := _ArbitrumGatewayRouter.contract.WatchLogs(opts, "TransferRouted", tokenRule, _userFromRule, _userToRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(ArbitrumGatewayRouterTransferRouted)
				if err := _ArbitrumGatewayRouter.contract.UnpackLog(event, "TransferRouted", log); err != nil {
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

func (_ArbitrumGatewayRouter *ArbitrumGatewayRouterFilterer) ParseTransferRouted(log types.Log) (*ArbitrumGatewayRouterTransferRouted, error) {
	event := new(ArbitrumGatewayRouterTransferRouted)
	if err := _ArbitrumGatewayRouter.contract.UnpackLog(event, "TransferRouted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_ArbitrumGatewayRouter *ArbitrumGatewayRouter) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _ArbitrumGatewayRouter.abi.Events["DefaultGatewayUpdated"].ID:
		return _ArbitrumGatewayRouter.ParseDefaultGatewayUpdated(log)
	case _ArbitrumGatewayRouter.abi.Events["GatewaySet"].ID:
		return _ArbitrumGatewayRouter.ParseGatewaySet(log)
	case _ArbitrumGatewayRouter.abi.Events["TransferRouted"].ID:
		return _ArbitrumGatewayRouter.ParseTransferRouted(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (ArbitrumGatewayRouterDefaultGatewayUpdated) Topic() common.Hash {
	return common.HexToHash("0x3a8f8eb961383a94d41d193e16a3af73eaddfd5764a4c640257323a1603ac331")
}

func (ArbitrumGatewayRouterGatewaySet) Topic() common.Hash {
	return common.HexToHash("0x812ca95fe4492a9e2d1f2723c2c40c03a60a27b059581ae20ac4e4d73bfba354")
}

func (ArbitrumGatewayRouterTransferRouted) Topic() common.Hash {
	return common.HexToHash("0x85291dff2161a93c2f12c819d31889c96c63042116f5bc5a205aa701c2c429f5")
}

func (_ArbitrumGatewayRouter *ArbitrumGatewayRouter) Address() common.Address {
	return _ArbitrumGatewayRouter.address
}

type ArbitrumGatewayRouterInterface interface {
	CalculateL2TokenAddress(opts *bind.CallOpts, l1ERC20 common.Address) (common.Address, error)

	DefaultGateway(opts *bind.CallOpts) (common.Address, error)

	GetGateway(opts *bind.CallOpts, _token common.Address) (common.Address, error)

	GetOutboundCalldata(opts *bind.CallOpts, _token common.Address, _from common.Address, _to common.Address, _amount *big.Int, _data []byte) ([]byte, error)

	FinalizeInboundTransfer(opts *bind.TransactOpts, _token common.Address, _from common.Address, _to common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error)

	OutboundTransfer(opts *bind.TransactOpts, _token common.Address, _to common.Address, _amount *big.Int, _maxGas *big.Int, _gasPriceBid *big.Int, _data []byte) (*types.Transaction, error)

	FilterDefaultGatewayUpdated(opts *bind.FilterOpts) (*ArbitrumGatewayRouterDefaultGatewayUpdatedIterator, error)

	WatchDefaultGatewayUpdated(opts *bind.WatchOpts, sink chan<- *ArbitrumGatewayRouterDefaultGatewayUpdated) (event.Subscription, error)

	ParseDefaultGatewayUpdated(log types.Log) (*ArbitrumGatewayRouterDefaultGatewayUpdated, error)

	FilterGatewaySet(opts *bind.FilterOpts, l1Token []common.Address, gateway []common.Address) (*ArbitrumGatewayRouterGatewaySetIterator, error)

	WatchGatewaySet(opts *bind.WatchOpts, sink chan<- *ArbitrumGatewayRouterGatewaySet, l1Token []common.Address, gateway []common.Address) (event.Subscription, error)

	ParseGatewaySet(log types.Log) (*ArbitrumGatewayRouterGatewaySet, error)

	FilterTransferRouted(opts *bind.FilterOpts, token []common.Address, _userFrom []common.Address, _userTo []common.Address) (*ArbitrumGatewayRouterTransferRoutedIterator, error)

	WatchTransferRouted(opts *bind.WatchOpts, sink chan<- *ArbitrumGatewayRouterTransferRouted, token []common.Address, _userFrom []common.Address, _userTo []common.Address) (event.Subscription, error)

	ParseTransferRouted(log types.Log) (*ArbitrumGatewayRouterTransferRouted, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
