// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package log_emitter

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

var LogEmitterMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"Log1\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"Log2\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"Log3\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"Log4\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"v\",\"type\":\"uint256[]\"}],\"name\":\"EmitLog1\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"v\",\"type\":\"uint256[]\"}],\"name\":\"EmitLog2\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"v\",\"type\":\"string[]\"}],\"name\":\"EmitLog3\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"v\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256\",\"name\":\"w\",\"type\":\"uint256\"}],\"name\":\"EmitLog4\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5061060a806100206000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c8063696933c914610051578063a0ff277514610066578063bc253bc014610079578063d9c21f461461008c575b600080fd5b61006461005f36600461035c565b61009f565b005b610064610074366004610399565b610113565b61006461008736600461035c565b61017d565b61006461009a3660046103de565b6101e1565b60005b815181101561010f577f46692c0e59ca9cd1ad8f984a9d11715ec83424398b7eed4e05c8ce84662415a88282815181106100de576100de610503565b60200260200101516040516100f591815260200190565b60405180910390a18061010781610532565b9150506100a2565b5050565b60005b8251811015610178578183828151811061013257610132610503565b60200260200101517fba21d5b63d64546cb4ab29e370a8972bf26f78cb0c395391b4f451699fdfdc5d60405160405180910390a38061017081610532565b915050610116565b505050565b60005b815181101561010f5781818151811061019b5761019b610503565b60200260200101517f624fb00c2ce79f34cb543884c3af64816dce0f4cec3d32661959e49d488a7a9360405160405180910390a2806101d981610532565b915050610180565b60005b815181101561010f577fb94ec34dfe32a8a7170992a093976368d1e63decf8f0bc0b38a8eb89cc9f95cf82828151811061022057610220610503565b60200260200101516040516102359190610591565b60405180910390a18061024781610532565b9150506101e4565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156102c5576102c561024f565b604052919050565b600067ffffffffffffffff8211156102e7576102e761024f565b5060051b60200190565b600082601f83011261030257600080fd5b81356020610317610312836102cd565b61027e565b82815260059290921b8401810191818101908684111561033657600080fd5b8286015b84811015610351578035835291830191830161033a565b509695505050505050565b60006020828403121561036e57600080fd5b813567ffffffffffffffff81111561038557600080fd5b610391848285016102f1565b949350505050565b600080604083850312156103ac57600080fd5b823567ffffffffffffffff8111156103c357600080fd5b6103cf858286016102f1565b95602094909401359450505050565b600060208083850312156103f157600080fd5b823567ffffffffffffffff8082111561040957600080fd5b8185019150601f868184011261041e57600080fd5b823561042c610312826102cd565b81815260059190911b8401850190858101908983111561044b57600080fd5b8686015b838110156104f5578035868111156104675760008081fd5b8701603f81018c136104795760008081fd5b8881013560408882111561048f5761048f61024f565b6104be8b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08a8501160161027e565b8281528e828486010111156104d35760008081fd5b828285018d83013760009281018c01929092525084525091870191870161044f565b509998505050505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361058a577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b5060010190565b600060208083528351808285015260005b818110156105be578581018301518582016040015282016105a2565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f830116850101925050509291505056fea164736f6c6343000813000a",
}

var LogEmitterABI = LogEmitterMetaData.ABI

var LogEmitterBin = LogEmitterMetaData.Bin

func DeployLogEmitter(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *LogEmitter, error) {
	parsed, err := LogEmitterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(LogEmitterBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &LogEmitter{address: address, abi: *parsed, LogEmitterCaller: LogEmitterCaller{contract: contract}, LogEmitterTransactor: LogEmitterTransactor{contract: contract}, LogEmitterFilterer: LogEmitterFilterer{contract: contract}}, nil
}

type LogEmitter struct {
	address common.Address
	abi     abi.ABI
	LogEmitterCaller
	LogEmitterTransactor
	LogEmitterFilterer
}

type LogEmitterCaller struct {
	contract *bind.BoundContract
}

type LogEmitterTransactor struct {
	contract *bind.BoundContract
}

type LogEmitterFilterer struct {
	contract *bind.BoundContract
}

type LogEmitterSession struct {
	Contract     *LogEmitter
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type LogEmitterCallerSession struct {
	Contract *LogEmitterCaller
	CallOpts bind.CallOpts
}

type LogEmitterTransactorSession struct {
	Contract     *LogEmitterTransactor
	TransactOpts bind.TransactOpts
}

type LogEmitterRaw struct {
	Contract *LogEmitter
}

type LogEmitterCallerRaw struct {
	Contract *LogEmitterCaller
}

type LogEmitterTransactorRaw struct {
	Contract *LogEmitterTransactor
}

func NewLogEmitter(address common.Address, backend bind.ContractBackend) (*LogEmitter, error) {
	abi, err := abi.JSON(strings.NewReader(LogEmitterABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindLogEmitter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LogEmitter{address: address, abi: abi, LogEmitterCaller: LogEmitterCaller{contract: contract}, LogEmitterTransactor: LogEmitterTransactor{contract: contract}, LogEmitterFilterer: LogEmitterFilterer{contract: contract}}, nil
}

func NewLogEmitterCaller(address common.Address, caller bind.ContractCaller) (*LogEmitterCaller, error) {
	contract, err := bindLogEmitter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LogEmitterCaller{contract: contract}, nil
}

func NewLogEmitterTransactor(address common.Address, transactor bind.ContractTransactor) (*LogEmitterTransactor, error) {
	contract, err := bindLogEmitter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LogEmitterTransactor{contract: contract}, nil
}

func NewLogEmitterFilterer(address common.Address, filterer bind.ContractFilterer) (*LogEmitterFilterer, error) {
	contract, err := bindLogEmitter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LogEmitterFilterer{contract: contract}, nil
}

func bindLogEmitter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := LogEmitterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_LogEmitter *LogEmitterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LogEmitter.Contract.LogEmitterCaller.contract.Call(opts, result, method, params...)
}

func (_LogEmitter *LogEmitterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LogEmitter.Contract.LogEmitterTransactor.contract.Transfer(opts)
}

func (_LogEmitter *LogEmitterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LogEmitter.Contract.LogEmitterTransactor.contract.Transact(opts, method, params...)
}

func (_LogEmitter *LogEmitterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LogEmitter.Contract.contract.Call(opts, result, method, params...)
}

func (_LogEmitter *LogEmitterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LogEmitter.Contract.contract.Transfer(opts)
}

func (_LogEmitter *LogEmitterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LogEmitter.Contract.contract.Transact(opts, method, params...)
}

func (_LogEmitter *LogEmitterTransactor) EmitLog1(opts *bind.TransactOpts, v []*big.Int) (*types.Transaction, error) {
	return _LogEmitter.contract.Transact(opts, "EmitLog1", v)
}

func (_LogEmitter *LogEmitterSession) EmitLog1(v []*big.Int) (*types.Transaction, error) {
	return _LogEmitter.Contract.EmitLog1(&_LogEmitter.TransactOpts, v)
}

func (_LogEmitter *LogEmitterTransactorSession) EmitLog1(v []*big.Int) (*types.Transaction, error) {
	return _LogEmitter.Contract.EmitLog1(&_LogEmitter.TransactOpts, v)
}

func (_LogEmitter *LogEmitterTransactor) EmitLog2(opts *bind.TransactOpts, v []*big.Int) (*types.Transaction, error) {
	return _LogEmitter.contract.Transact(opts, "EmitLog2", v)
}

func (_LogEmitter *LogEmitterSession) EmitLog2(v []*big.Int) (*types.Transaction, error) {
	return _LogEmitter.Contract.EmitLog2(&_LogEmitter.TransactOpts, v)
}

func (_LogEmitter *LogEmitterTransactorSession) EmitLog2(v []*big.Int) (*types.Transaction, error) {
	return _LogEmitter.Contract.EmitLog2(&_LogEmitter.TransactOpts, v)
}

func (_LogEmitter *LogEmitterTransactor) EmitLog3(opts *bind.TransactOpts, v []string) (*types.Transaction, error) {
	return _LogEmitter.contract.Transact(opts, "EmitLog3", v)
}

func (_LogEmitter *LogEmitterSession) EmitLog3(v []string) (*types.Transaction, error) {
	return _LogEmitter.Contract.EmitLog3(&_LogEmitter.TransactOpts, v)
}

func (_LogEmitter *LogEmitterTransactorSession) EmitLog3(v []string) (*types.Transaction, error) {
	return _LogEmitter.Contract.EmitLog3(&_LogEmitter.TransactOpts, v)
}

func (_LogEmitter *LogEmitterTransactor) EmitLog4(opts *bind.TransactOpts, v []*big.Int, w *big.Int) (*types.Transaction, error) {
	return _LogEmitter.contract.Transact(opts, "EmitLog4", v, w)
}

func (_LogEmitter *LogEmitterSession) EmitLog4(v []*big.Int, w *big.Int) (*types.Transaction, error) {
	return _LogEmitter.Contract.EmitLog4(&_LogEmitter.TransactOpts, v, w)
}

func (_LogEmitter *LogEmitterTransactorSession) EmitLog4(v []*big.Int, w *big.Int) (*types.Transaction, error) {
	return _LogEmitter.Contract.EmitLog4(&_LogEmitter.TransactOpts, v, w)
}

type LogEmitterLog1Iterator struct {
	Event *LogEmitterLog1

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LogEmitterLog1Iterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LogEmitterLog1)
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
		it.Event = new(LogEmitterLog1)
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

func (it *LogEmitterLog1Iterator) Error() error {
	return it.fail
}

func (it *LogEmitterLog1Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LogEmitterLog1 struct {
	Arg0 *big.Int
	Raw  types.Log
}

func (_LogEmitter *LogEmitterFilterer) FilterLog1(opts *bind.FilterOpts) (*LogEmitterLog1Iterator, error) {

	logs, sub, err := _LogEmitter.contract.FilterLogs(opts, "Log1")
	if err != nil {
		return nil, err
	}
	return &LogEmitterLog1Iterator{contract: _LogEmitter.contract, event: "Log1", logs: logs, sub: sub}, nil
}

func (_LogEmitter *LogEmitterFilterer) WatchLog1(opts *bind.WatchOpts, sink chan<- *LogEmitterLog1) (event.Subscription, error) {

	logs, sub, err := _LogEmitter.contract.WatchLogs(opts, "Log1")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LogEmitterLog1)
				if err := _LogEmitter.contract.UnpackLog(event, "Log1", log); err != nil {
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

func (_LogEmitter *LogEmitterFilterer) ParseLog1(log types.Log) (*LogEmitterLog1, error) {
	event := new(LogEmitterLog1)
	if err := _LogEmitter.contract.UnpackLog(event, "Log1", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LogEmitterLog2Iterator struct {
	Event *LogEmitterLog2

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LogEmitterLog2Iterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LogEmitterLog2)
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
		it.Event = new(LogEmitterLog2)
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

func (it *LogEmitterLog2Iterator) Error() error {
	return it.fail
}

func (it *LogEmitterLog2Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LogEmitterLog2 struct {
	Arg0 *big.Int
	Raw  types.Log
}

func (_LogEmitter *LogEmitterFilterer) FilterLog2(opts *bind.FilterOpts, arg0 []*big.Int) (*LogEmitterLog2Iterator, error) {

	var arg0Rule []interface{}
	for _, arg0Item := range arg0 {
		arg0Rule = append(arg0Rule, arg0Item)
	}

	logs, sub, err := _LogEmitter.contract.FilterLogs(opts, "Log2", arg0Rule)
	if err != nil {
		return nil, err
	}
	return &LogEmitterLog2Iterator{contract: _LogEmitter.contract, event: "Log2", logs: logs, sub: sub}, nil
}

func (_LogEmitter *LogEmitterFilterer) WatchLog2(opts *bind.WatchOpts, sink chan<- *LogEmitterLog2, arg0 []*big.Int) (event.Subscription, error) {

	var arg0Rule []interface{}
	for _, arg0Item := range arg0 {
		arg0Rule = append(arg0Rule, arg0Item)
	}

	logs, sub, err := _LogEmitter.contract.WatchLogs(opts, "Log2", arg0Rule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LogEmitterLog2)
				if err := _LogEmitter.contract.UnpackLog(event, "Log2", log); err != nil {
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

func (_LogEmitter *LogEmitterFilterer) ParseLog2(log types.Log) (*LogEmitterLog2, error) {
	event := new(LogEmitterLog2)
	if err := _LogEmitter.contract.UnpackLog(event, "Log2", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LogEmitterLog3Iterator struct {
	Event *LogEmitterLog3

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LogEmitterLog3Iterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LogEmitterLog3)
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
		it.Event = new(LogEmitterLog3)
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

func (it *LogEmitterLog3Iterator) Error() error {
	return it.fail
}

func (it *LogEmitterLog3Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LogEmitterLog3 struct {
	Arg0 string
	Raw  types.Log
}

func (_LogEmitter *LogEmitterFilterer) FilterLog3(opts *bind.FilterOpts) (*LogEmitterLog3Iterator, error) {

	logs, sub, err := _LogEmitter.contract.FilterLogs(opts, "Log3")
	if err != nil {
		return nil, err
	}
	return &LogEmitterLog3Iterator{contract: _LogEmitter.contract, event: "Log3", logs: logs, sub: sub}, nil
}

func (_LogEmitter *LogEmitterFilterer) WatchLog3(opts *bind.WatchOpts, sink chan<- *LogEmitterLog3) (event.Subscription, error) {

	logs, sub, err := _LogEmitter.contract.WatchLogs(opts, "Log3")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LogEmitterLog3)
				if err := _LogEmitter.contract.UnpackLog(event, "Log3", log); err != nil {
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

func (_LogEmitter *LogEmitterFilterer) ParseLog3(log types.Log) (*LogEmitterLog3, error) {
	event := new(LogEmitterLog3)
	if err := _LogEmitter.contract.UnpackLog(event, "Log3", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type LogEmitterLog4Iterator struct {
	Event *LogEmitterLog4

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LogEmitterLog4Iterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LogEmitterLog4)
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
		it.Event = new(LogEmitterLog4)
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

func (it *LogEmitterLog4Iterator) Error() error {
	return it.fail
}

func (it *LogEmitterLog4Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LogEmitterLog4 struct {
	Arg0 *big.Int
	Arg1 *big.Int
	Raw  types.Log
}

func (_LogEmitter *LogEmitterFilterer) FilterLog4(opts *bind.FilterOpts, arg0 []*big.Int, arg1 []*big.Int) (*LogEmitterLog4Iterator, error) {

	var arg0Rule []interface{}
	for _, arg0Item := range arg0 {
		arg0Rule = append(arg0Rule, arg0Item)
	}
	var arg1Rule []interface{}
	for _, arg1Item := range arg1 {
		arg1Rule = append(arg1Rule, arg1Item)
	}

	logs, sub, err := _LogEmitter.contract.FilterLogs(opts, "Log4", arg0Rule, arg1Rule)
	if err != nil {
		return nil, err
	}
	return &LogEmitterLog4Iterator{contract: _LogEmitter.contract, event: "Log4", logs: logs, sub: sub}, nil
}

func (_LogEmitter *LogEmitterFilterer) WatchLog4(opts *bind.WatchOpts, sink chan<- *LogEmitterLog4, arg0 []*big.Int, arg1 []*big.Int) (event.Subscription, error) {

	var arg0Rule []interface{}
	for _, arg0Item := range arg0 {
		arg0Rule = append(arg0Rule, arg0Item)
	}
	var arg1Rule []interface{}
	for _, arg1Item := range arg1 {
		arg1Rule = append(arg1Rule, arg1Item)
	}

	logs, sub, err := _LogEmitter.contract.WatchLogs(opts, "Log4", arg0Rule, arg1Rule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LogEmitterLog4)
				if err := _LogEmitter.contract.UnpackLog(event, "Log4", log); err != nil {
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

func (_LogEmitter *LogEmitterFilterer) ParseLog4(log types.Log) (*LogEmitterLog4, error) {
	event := new(LogEmitterLog4)
	if err := _LogEmitter.contract.UnpackLog(event, "Log4", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_LogEmitter *LogEmitter) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _LogEmitter.abi.Events["Log1"].ID:
		return _LogEmitter.ParseLog1(log)
	case _LogEmitter.abi.Events["Log2"].ID:
		return _LogEmitter.ParseLog2(log)
	case _LogEmitter.abi.Events["Log3"].ID:
		return _LogEmitter.ParseLog3(log)
	case _LogEmitter.abi.Events["Log4"].ID:
		return _LogEmitter.ParseLog4(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (LogEmitterLog1) Topic() common.Hash {
	return common.HexToHash("0x46692c0e59ca9cd1ad8f984a9d11715ec83424398b7eed4e05c8ce84662415a8")
}

func (LogEmitterLog2) Topic() common.Hash {
	return common.HexToHash("0x624fb00c2ce79f34cb543884c3af64816dce0f4cec3d32661959e49d488a7a93")
}

func (LogEmitterLog3) Topic() common.Hash {
	return common.HexToHash("0xb94ec34dfe32a8a7170992a093976368d1e63decf8f0bc0b38a8eb89cc9f95cf")
}

func (LogEmitterLog4) Topic() common.Hash {
	return common.HexToHash("0xba21d5b63d64546cb4ab29e370a8972bf26f78cb0c395391b4f451699fdfdc5d")
}

func (_LogEmitter *LogEmitter) Address() common.Address {
	return _LogEmitter.address
}

type LogEmitterInterface interface {
	EmitLog1(opts *bind.TransactOpts, v []*big.Int) (*types.Transaction, error)

	EmitLog2(opts *bind.TransactOpts, v []*big.Int) (*types.Transaction, error)

	EmitLog3(opts *bind.TransactOpts, v []string) (*types.Transaction, error)

	EmitLog4(opts *bind.TransactOpts, v []*big.Int, w *big.Int) (*types.Transaction, error)

	FilterLog1(opts *bind.FilterOpts) (*LogEmitterLog1Iterator, error)

	WatchLog1(opts *bind.WatchOpts, sink chan<- *LogEmitterLog1) (event.Subscription, error)

	ParseLog1(log types.Log) (*LogEmitterLog1, error)

	FilterLog2(opts *bind.FilterOpts, arg0 []*big.Int) (*LogEmitterLog2Iterator, error)

	WatchLog2(opts *bind.WatchOpts, sink chan<- *LogEmitterLog2, arg0 []*big.Int) (event.Subscription, error)

	ParseLog2(log types.Log) (*LogEmitterLog2, error)

	FilterLog3(opts *bind.FilterOpts) (*LogEmitterLog3Iterator, error)

	WatchLog3(opts *bind.WatchOpts, sink chan<- *LogEmitterLog3) (event.Subscription, error)

	ParseLog3(log types.Log) (*LogEmitterLog3, error)

	FilterLog4(opts *bind.FilterOpts, arg0 []*big.Int, arg1 []*big.Int) (*LogEmitterLog4Iterator, error)

	WatchLog4(opts *bind.WatchOpts, sink chan<- *LogEmitterLog4, arg0 []*big.Int, arg1 []*big.Int) (event.Subscription, error)

	ParseLog4(log types.Log) (*LogEmitterLog4, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
