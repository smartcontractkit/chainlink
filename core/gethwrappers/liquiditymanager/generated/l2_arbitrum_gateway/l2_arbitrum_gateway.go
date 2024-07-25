// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package l2_arbitrum_gateway

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

var L2ArbitrumGatewayMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"l1Token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"DepositFinalized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"TxToL1\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"l1Token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"_l2ToL1Id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_exitNum\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"WithdrawalInitiated\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"l1ERC20\",\"type\":\"address\"}],\"name\":\"calculateL2TokenAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"counterpartGateway\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"exitNum\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"finalizeInboundTransfer\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"getOutboundCalldata\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"outboundCalldata\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_l1Token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"outboundTransfer\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_l1Token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"outboundTransfer\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"res\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"postUpgradeInit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"router\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

var L2ArbitrumGatewayABI = L2ArbitrumGatewayMetaData.ABI

type L2ArbitrumGateway struct {
	address common.Address
	abi     abi.ABI
	L2ArbitrumGatewayCaller
	L2ArbitrumGatewayTransactor
	L2ArbitrumGatewayFilterer
}

type L2ArbitrumGatewayCaller struct {
	contract *bind.BoundContract
}

type L2ArbitrumGatewayTransactor struct {
	contract *bind.BoundContract
}

type L2ArbitrumGatewayFilterer struct {
	contract *bind.BoundContract
}

type L2ArbitrumGatewaySession struct {
	Contract     *L2ArbitrumGateway
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type L2ArbitrumGatewayCallerSession struct {
	Contract *L2ArbitrumGatewayCaller
	CallOpts bind.CallOpts
}

type L2ArbitrumGatewayTransactorSession struct {
	Contract     *L2ArbitrumGatewayTransactor
	TransactOpts bind.TransactOpts
}

type L2ArbitrumGatewayRaw struct {
	Contract *L2ArbitrumGateway
}

type L2ArbitrumGatewayCallerRaw struct {
	Contract *L2ArbitrumGatewayCaller
}

type L2ArbitrumGatewayTransactorRaw struct {
	Contract *L2ArbitrumGatewayTransactor
}

func NewL2ArbitrumGateway(address common.Address, backend bind.ContractBackend) (*L2ArbitrumGateway, error) {
	abi, err := abi.JSON(strings.NewReader(L2ArbitrumGatewayABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindL2ArbitrumGateway(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &L2ArbitrumGateway{address: address, abi: abi, L2ArbitrumGatewayCaller: L2ArbitrumGatewayCaller{contract: contract}, L2ArbitrumGatewayTransactor: L2ArbitrumGatewayTransactor{contract: contract}, L2ArbitrumGatewayFilterer: L2ArbitrumGatewayFilterer{contract: contract}}, nil
}

func NewL2ArbitrumGatewayCaller(address common.Address, caller bind.ContractCaller) (*L2ArbitrumGatewayCaller, error) {
	contract, err := bindL2ArbitrumGateway(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &L2ArbitrumGatewayCaller{contract: contract}, nil
}

func NewL2ArbitrumGatewayTransactor(address common.Address, transactor bind.ContractTransactor) (*L2ArbitrumGatewayTransactor, error) {
	contract, err := bindL2ArbitrumGateway(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &L2ArbitrumGatewayTransactor{contract: contract}, nil
}

func NewL2ArbitrumGatewayFilterer(address common.Address, filterer bind.ContractFilterer) (*L2ArbitrumGatewayFilterer, error) {
	contract, err := bindL2ArbitrumGateway(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &L2ArbitrumGatewayFilterer{contract: contract}, nil
}

func bindL2ArbitrumGateway(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := L2ArbitrumGatewayMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_L2ArbitrumGateway *L2ArbitrumGatewayRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _L2ArbitrumGateway.Contract.L2ArbitrumGatewayCaller.contract.Call(opts, result, method, params...)
}

func (_L2ArbitrumGateway *L2ArbitrumGatewayRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L2ArbitrumGateway.Contract.L2ArbitrumGatewayTransactor.contract.Transfer(opts)
}

func (_L2ArbitrumGateway *L2ArbitrumGatewayRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _L2ArbitrumGateway.Contract.L2ArbitrumGatewayTransactor.contract.Transact(opts, method, params...)
}

func (_L2ArbitrumGateway *L2ArbitrumGatewayCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _L2ArbitrumGateway.Contract.contract.Call(opts, result, method, params...)
}

func (_L2ArbitrumGateway *L2ArbitrumGatewayTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L2ArbitrumGateway.Contract.contract.Transfer(opts)
}

func (_L2ArbitrumGateway *L2ArbitrumGatewayTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _L2ArbitrumGateway.Contract.contract.Transact(opts, method, params...)
}

func (_L2ArbitrumGateway *L2ArbitrumGatewayCaller) CalculateL2TokenAddress(opts *bind.CallOpts, l1ERC20 common.Address) (common.Address, error) {
	var out []interface{}
	err := _L2ArbitrumGateway.contract.Call(opts, &out, "calculateL2TokenAddress", l1ERC20)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_L2ArbitrumGateway *L2ArbitrumGatewaySession) CalculateL2TokenAddress(l1ERC20 common.Address) (common.Address, error) {
	return _L2ArbitrumGateway.Contract.CalculateL2TokenAddress(&_L2ArbitrumGateway.CallOpts, l1ERC20)
}

func (_L2ArbitrumGateway *L2ArbitrumGatewayCallerSession) CalculateL2TokenAddress(l1ERC20 common.Address) (common.Address, error) {
	return _L2ArbitrumGateway.Contract.CalculateL2TokenAddress(&_L2ArbitrumGateway.CallOpts, l1ERC20)
}

func (_L2ArbitrumGateway *L2ArbitrumGatewayCaller) CounterpartGateway(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L2ArbitrumGateway.contract.Call(opts, &out, "counterpartGateway")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_L2ArbitrumGateway *L2ArbitrumGatewaySession) CounterpartGateway() (common.Address, error) {
	return _L2ArbitrumGateway.Contract.CounterpartGateway(&_L2ArbitrumGateway.CallOpts)
}

func (_L2ArbitrumGateway *L2ArbitrumGatewayCallerSession) CounterpartGateway() (common.Address, error) {
	return _L2ArbitrumGateway.Contract.CounterpartGateway(&_L2ArbitrumGateway.CallOpts)
}

func (_L2ArbitrumGateway *L2ArbitrumGatewayCaller) ExitNum(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _L2ArbitrumGateway.contract.Call(opts, &out, "exitNum")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_L2ArbitrumGateway *L2ArbitrumGatewaySession) ExitNum() (*big.Int, error) {
	return _L2ArbitrumGateway.Contract.ExitNum(&_L2ArbitrumGateway.CallOpts)
}

func (_L2ArbitrumGateway *L2ArbitrumGatewayCallerSession) ExitNum() (*big.Int, error) {
	return _L2ArbitrumGateway.Contract.ExitNum(&_L2ArbitrumGateway.CallOpts)
}

func (_L2ArbitrumGateway *L2ArbitrumGatewayCaller) GetOutboundCalldata(opts *bind.CallOpts, _token common.Address, _from common.Address, _to common.Address, _amount *big.Int, _data []byte) ([]byte, error) {
	var out []interface{}
	err := _L2ArbitrumGateway.contract.Call(opts, &out, "getOutboundCalldata", _token, _from, _to, _amount, _data)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

func (_L2ArbitrumGateway *L2ArbitrumGatewaySession) GetOutboundCalldata(_token common.Address, _from common.Address, _to common.Address, _amount *big.Int, _data []byte) ([]byte, error) {
	return _L2ArbitrumGateway.Contract.GetOutboundCalldata(&_L2ArbitrumGateway.CallOpts, _token, _from, _to, _amount, _data)
}

func (_L2ArbitrumGateway *L2ArbitrumGatewayCallerSession) GetOutboundCalldata(_token common.Address, _from common.Address, _to common.Address, _amount *big.Int, _data []byte) ([]byte, error) {
	return _L2ArbitrumGateway.Contract.GetOutboundCalldata(&_L2ArbitrumGateway.CallOpts, _token, _from, _to, _amount, _data)
}

func (_L2ArbitrumGateway *L2ArbitrumGatewayCaller) Router(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L2ArbitrumGateway.contract.Call(opts, &out, "router")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_L2ArbitrumGateway *L2ArbitrumGatewaySession) Router() (common.Address, error) {
	return _L2ArbitrumGateway.Contract.Router(&_L2ArbitrumGateway.CallOpts)
}

func (_L2ArbitrumGateway *L2ArbitrumGatewayCallerSession) Router() (common.Address, error) {
	return _L2ArbitrumGateway.Contract.Router(&_L2ArbitrumGateway.CallOpts)
}

func (_L2ArbitrumGateway *L2ArbitrumGatewayTransactor) FinalizeInboundTransfer(opts *bind.TransactOpts, _token common.Address, _from common.Address, _to common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _L2ArbitrumGateway.contract.Transact(opts, "finalizeInboundTransfer", _token, _from, _to, _amount, _data)
}

func (_L2ArbitrumGateway *L2ArbitrumGatewaySession) FinalizeInboundTransfer(_token common.Address, _from common.Address, _to common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _L2ArbitrumGateway.Contract.FinalizeInboundTransfer(&_L2ArbitrumGateway.TransactOpts, _token, _from, _to, _amount, _data)
}

func (_L2ArbitrumGateway *L2ArbitrumGatewayTransactorSession) FinalizeInboundTransfer(_token common.Address, _from common.Address, _to common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _L2ArbitrumGateway.Contract.FinalizeInboundTransfer(&_L2ArbitrumGateway.TransactOpts, _token, _from, _to, _amount, _data)
}

func (_L2ArbitrumGateway *L2ArbitrumGatewayTransactor) OutboundTransfer(opts *bind.TransactOpts, _l1Token common.Address, _to common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _L2ArbitrumGateway.contract.Transact(opts, "outboundTransfer", _l1Token, _to, _amount, _data)
}

func (_L2ArbitrumGateway *L2ArbitrumGatewaySession) OutboundTransfer(_l1Token common.Address, _to common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _L2ArbitrumGateway.Contract.OutboundTransfer(&_L2ArbitrumGateway.TransactOpts, _l1Token, _to, _amount, _data)
}

func (_L2ArbitrumGateway *L2ArbitrumGatewayTransactorSession) OutboundTransfer(_l1Token common.Address, _to common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _L2ArbitrumGateway.Contract.OutboundTransfer(&_L2ArbitrumGateway.TransactOpts, _l1Token, _to, _amount, _data)
}

func (_L2ArbitrumGateway *L2ArbitrumGatewayTransactor) OutboundTransfer0(opts *bind.TransactOpts, _l1Token common.Address, _to common.Address, _amount *big.Int, arg3 *big.Int, arg4 *big.Int, _data []byte) (*types.Transaction, error) {
	return _L2ArbitrumGateway.contract.Transact(opts, "outboundTransfer0", _l1Token, _to, _amount, arg3, arg4, _data)
}

func (_L2ArbitrumGateway *L2ArbitrumGatewaySession) OutboundTransfer0(_l1Token common.Address, _to common.Address, _amount *big.Int, arg3 *big.Int, arg4 *big.Int, _data []byte) (*types.Transaction, error) {
	return _L2ArbitrumGateway.Contract.OutboundTransfer0(&_L2ArbitrumGateway.TransactOpts, _l1Token, _to, _amount, arg3, arg4, _data)
}

func (_L2ArbitrumGateway *L2ArbitrumGatewayTransactorSession) OutboundTransfer0(_l1Token common.Address, _to common.Address, _amount *big.Int, arg3 *big.Int, arg4 *big.Int, _data []byte) (*types.Transaction, error) {
	return _L2ArbitrumGateway.Contract.OutboundTransfer0(&_L2ArbitrumGateway.TransactOpts, _l1Token, _to, _amount, arg3, arg4, _data)
}

func (_L2ArbitrumGateway *L2ArbitrumGatewayTransactor) PostUpgradeInit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L2ArbitrumGateway.contract.Transact(opts, "postUpgradeInit")
}

func (_L2ArbitrumGateway *L2ArbitrumGatewaySession) PostUpgradeInit() (*types.Transaction, error) {
	return _L2ArbitrumGateway.Contract.PostUpgradeInit(&_L2ArbitrumGateway.TransactOpts)
}

func (_L2ArbitrumGateway *L2ArbitrumGatewayTransactorSession) PostUpgradeInit() (*types.Transaction, error) {
	return _L2ArbitrumGateway.Contract.PostUpgradeInit(&_L2ArbitrumGateway.TransactOpts)
}

type L2ArbitrumGatewayDepositFinalizedIterator struct {
	Event *L2ArbitrumGatewayDepositFinalized

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *L2ArbitrumGatewayDepositFinalizedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L2ArbitrumGatewayDepositFinalized)
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
		it.Event = new(L2ArbitrumGatewayDepositFinalized)
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

func (it *L2ArbitrumGatewayDepositFinalizedIterator) Error() error {
	return it.fail
}

func (it *L2ArbitrumGatewayDepositFinalizedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type L2ArbitrumGatewayDepositFinalized struct {
	L1Token common.Address
	From    common.Address
	To      common.Address
	Amount  *big.Int
	Raw     types.Log
}

func (_L2ArbitrumGateway *L2ArbitrumGatewayFilterer) FilterDepositFinalized(opts *bind.FilterOpts, l1Token []common.Address, _from []common.Address, _to []common.Address) (*L2ArbitrumGatewayDepositFinalizedIterator, error) {

	var l1TokenRule []interface{}
	for _, l1TokenItem := range l1Token {
		l1TokenRule = append(l1TokenRule, l1TokenItem)
	}
	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _L2ArbitrumGateway.contract.FilterLogs(opts, "DepositFinalized", l1TokenRule, _fromRule, _toRule)
	if err != nil {
		return nil, err
	}
	return &L2ArbitrumGatewayDepositFinalizedIterator{contract: _L2ArbitrumGateway.contract, event: "DepositFinalized", logs: logs, sub: sub}, nil
}

func (_L2ArbitrumGateway *L2ArbitrumGatewayFilterer) WatchDepositFinalized(opts *bind.WatchOpts, sink chan<- *L2ArbitrumGatewayDepositFinalized, l1Token []common.Address, _from []common.Address, _to []common.Address) (event.Subscription, error) {

	var l1TokenRule []interface{}
	for _, l1TokenItem := range l1Token {
		l1TokenRule = append(l1TokenRule, l1TokenItem)
	}
	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _L2ArbitrumGateway.contract.WatchLogs(opts, "DepositFinalized", l1TokenRule, _fromRule, _toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(L2ArbitrumGatewayDepositFinalized)
				if err := _L2ArbitrumGateway.contract.UnpackLog(event, "DepositFinalized", log); err != nil {
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

func (_L2ArbitrumGateway *L2ArbitrumGatewayFilterer) ParseDepositFinalized(log types.Log) (*L2ArbitrumGatewayDepositFinalized, error) {
	event := new(L2ArbitrumGatewayDepositFinalized)
	if err := _L2ArbitrumGateway.contract.UnpackLog(event, "DepositFinalized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type L2ArbitrumGatewayTxToL1Iterator struct {
	Event *L2ArbitrumGatewayTxToL1

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *L2ArbitrumGatewayTxToL1Iterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L2ArbitrumGatewayTxToL1)
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
		it.Event = new(L2ArbitrumGatewayTxToL1)
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

func (it *L2ArbitrumGatewayTxToL1Iterator) Error() error {
	return it.fail
}

func (it *L2ArbitrumGatewayTxToL1Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type L2ArbitrumGatewayTxToL1 struct {
	From common.Address
	To   common.Address
	Id   *big.Int
	Data []byte
	Raw  types.Log
}

func (_L2ArbitrumGateway *L2ArbitrumGatewayFilterer) FilterTxToL1(opts *bind.FilterOpts, _from []common.Address, _to []common.Address, _id []*big.Int) (*L2ArbitrumGatewayTxToL1Iterator, error) {

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

	logs, sub, err := _L2ArbitrumGateway.contract.FilterLogs(opts, "TxToL1", _fromRule, _toRule, _idRule)
	if err != nil {
		return nil, err
	}
	return &L2ArbitrumGatewayTxToL1Iterator{contract: _L2ArbitrumGateway.contract, event: "TxToL1", logs: logs, sub: sub}, nil
}

func (_L2ArbitrumGateway *L2ArbitrumGatewayFilterer) WatchTxToL1(opts *bind.WatchOpts, sink chan<- *L2ArbitrumGatewayTxToL1, _from []common.Address, _to []common.Address, _id []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _L2ArbitrumGateway.contract.WatchLogs(opts, "TxToL1", _fromRule, _toRule, _idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(L2ArbitrumGatewayTxToL1)
				if err := _L2ArbitrumGateway.contract.UnpackLog(event, "TxToL1", log); err != nil {
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

func (_L2ArbitrumGateway *L2ArbitrumGatewayFilterer) ParseTxToL1(log types.Log) (*L2ArbitrumGatewayTxToL1, error) {
	event := new(L2ArbitrumGatewayTxToL1)
	if err := _L2ArbitrumGateway.contract.UnpackLog(event, "TxToL1", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type L2ArbitrumGatewayWithdrawalInitiatedIterator struct {
	Event *L2ArbitrumGatewayWithdrawalInitiated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *L2ArbitrumGatewayWithdrawalInitiatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L2ArbitrumGatewayWithdrawalInitiated)
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
		it.Event = new(L2ArbitrumGatewayWithdrawalInitiated)
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

func (it *L2ArbitrumGatewayWithdrawalInitiatedIterator) Error() error {
	return it.fail
}

func (it *L2ArbitrumGatewayWithdrawalInitiatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type L2ArbitrumGatewayWithdrawalInitiated struct {
	L1Token  common.Address
	From     common.Address
	To       common.Address
	L2ToL1Id *big.Int
	ExitNum  *big.Int
	Amount   *big.Int
	Raw      types.Log
}

func (_L2ArbitrumGateway *L2ArbitrumGatewayFilterer) FilterWithdrawalInitiated(opts *bind.FilterOpts, _from []common.Address, _to []common.Address, _l2ToL1Id []*big.Int) (*L2ArbitrumGatewayWithdrawalInitiatedIterator, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}
	var _l2ToL1IdRule []interface{}
	for _, _l2ToL1IdItem := range _l2ToL1Id {
		_l2ToL1IdRule = append(_l2ToL1IdRule, _l2ToL1IdItem)
	}

	logs, sub, err := _L2ArbitrumGateway.contract.FilterLogs(opts, "WithdrawalInitiated", _fromRule, _toRule, _l2ToL1IdRule)
	if err != nil {
		return nil, err
	}
	return &L2ArbitrumGatewayWithdrawalInitiatedIterator{contract: _L2ArbitrumGateway.contract, event: "WithdrawalInitiated", logs: logs, sub: sub}, nil
}

func (_L2ArbitrumGateway *L2ArbitrumGatewayFilterer) WatchWithdrawalInitiated(opts *bind.WatchOpts, sink chan<- *L2ArbitrumGatewayWithdrawalInitiated, _from []common.Address, _to []common.Address, _l2ToL1Id []*big.Int) (event.Subscription, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}
	var _l2ToL1IdRule []interface{}
	for _, _l2ToL1IdItem := range _l2ToL1Id {
		_l2ToL1IdRule = append(_l2ToL1IdRule, _l2ToL1IdItem)
	}

	logs, sub, err := _L2ArbitrumGateway.contract.WatchLogs(opts, "WithdrawalInitiated", _fromRule, _toRule, _l2ToL1IdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(L2ArbitrumGatewayWithdrawalInitiated)
				if err := _L2ArbitrumGateway.contract.UnpackLog(event, "WithdrawalInitiated", log); err != nil {
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

func (_L2ArbitrumGateway *L2ArbitrumGatewayFilterer) ParseWithdrawalInitiated(log types.Log) (*L2ArbitrumGatewayWithdrawalInitiated, error) {
	event := new(L2ArbitrumGatewayWithdrawalInitiated)
	if err := _L2ArbitrumGateway.contract.UnpackLog(event, "WithdrawalInitiated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_L2ArbitrumGateway *L2ArbitrumGateway) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _L2ArbitrumGateway.abi.Events["DepositFinalized"].ID:
		return _L2ArbitrumGateway.ParseDepositFinalized(log)
	case _L2ArbitrumGateway.abi.Events["TxToL1"].ID:
		return _L2ArbitrumGateway.ParseTxToL1(log)
	case _L2ArbitrumGateway.abi.Events["WithdrawalInitiated"].ID:
		return _L2ArbitrumGateway.ParseWithdrawalInitiated(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (L2ArbitrumGatewayDepositFinalized) Topic() common.Hash {
	return common.HexToHash("0xc7f2e9c55c40a50fbc217dfc70cd39a222940dfa62145aa0ca49eb9535d4fcb2")
}

func (L2ArbitrumGatewayTxToL1) Topic() common.Hash {
	return common.HexToHash("0x2b986d32a0536b7e19baa48ab949fec7b903b7fad7730820b20632d100cc3a68")
}

func (L2ArbitrumGatewayWithdrawalInitiated) Topic() common.Hash {
	return common.HexToHash("0x3073a74ecb728d10be779fe19a74a1428e20468f5b4d167bf9c73d9067847d73")
}

func (_L2ArbitrumGateway *L2ArbitrumGateway) Address() common.Address {
	return _L2ArbitrumGateway.address
}

type L2ArbitrumGatewayInterface interface {
	CalculateL2TokenAddress(opts *bind.CallOpts, l1ERC20 common.Address) (common.Address, error)

	CounterpartGateway(opts *bind.CallOpts) (common.Address, error)

	ExitNum(opts *bind.CallOpts) (*big.Int, error)

	GetOutboundCalldata(opts *bind.CallOpts, _token common.Address, _from common.Address, _to common.Address, _amount *big.Int, _data []byte) ([]byte, error)

	Router(opts *bind.CallOpts) (common.Address, error)

	FinalizeInboundTransfer(opts *bind.TransactOpts, _token common.Address, _from common.Address, _to common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error)

	OutboundTransfer(opts *bind.TransactOpts, _l1Token common.Address, _to common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error)

	OutboundTransfer0(opts *bind.TransactOpts, _l1Token common.Address, _to common.Address, _amount *big.Int, arg3 *big.Int, arg4 *big.Int, _data []byte) (*types.Transaction, error)

	PostUpgradeInit(opts *bind.TransactOpts) (*types.Transaction, error)

	FilterDepositFinalized(opts *bind.FilterOpts, l1Token []common.Address, _from []common.Address, _to []common.Address) (*L2ArbitrumGatewayDepositFinalizedIterator, error)

	WatchDepositFinalized(opts *bind.WatchOpts, sink chan<- *L2ArbitrumGatewayDepositFinalized, l1Token []common.Address, _from []common.Address, _to []common.Address) (event.Subscription, error)

	ParseDepositFinalized(log types.Log) (*L2ArbitrumGatewayDepositFinalized, error)

	FilterTxToL1(opts *bind.FilterOpts, _from []common.Address, _to []common.Address, _id []*big.Int) (*L2ArbitrumGatewayTxToL1Iterator, error)

	WatchTxToL1(opts *bind.WatchOpts, sink chan<- *L2ArbitrumGatewayTxToL1, _from []common.Address, _to []common.Address, _id []*big.Int) (event.Subscription, error)

	ParseTxToL1(log types.Log) (*L2ArbitrumGatewayTxToL1, error)

	FilterWithdrawalInitiated(opts *bind.FilterOpts, _from []common.Address, _to []common.Address, _l2ToL1Id []*big.Int) (*L2ArbitrumGatewayWithdrawalInitiatedIterator, error)

	WatchWithdrawalInitiated(opts *bind.WatchOpts, sink chan<- *L2ArbitrumGatewayWithdrawalInitiated, _from []common.Address, _to []common.Address, _l2ToL1Id []*big.Int) (event.Subscription, error)

	ParseWithdrawalInitiated(log types.Log) (*L2ArbitrumGatewayWithdrawalInitiated, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
