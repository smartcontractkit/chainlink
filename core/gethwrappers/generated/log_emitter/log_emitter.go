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
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"Log1\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"Log2\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"Log3\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"v\",\"type\":\"uint256[]\"}],\"name\":\"EmitLog1\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"v\",\"type\":\"uint256[]\"}],\"name\":\"EmitLog2\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"v\",\"type\":\"string[]\"}],\"name\":\"EmitLog3\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5061053e806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c8063696933c914610046578063bc253bc01461005b578063d9c21f461461006e575b600080fd5b6100596100543660046102f5565b610081565b005b6100596100693660046102f5565b6100f5565b61005961007c3660046101c7565b610159565b60005b81518110156100f1577f46692c0e59ca9cd1ad8f984a9d11715ec83424398b7eed4e05c8ce84662415a88282815181106100c0576100c06104d3565b60200260200101516040516100d791815260200190565b60405180910390a1806100e981610473565b915050610084565b5050565b60005b81518110156100f157818181518110610113576101136104d3565b60200260200101517f624fb00c2ce79f34cb543884c3af64816dce0f4cec3d32661959e49d488a7a9360405160405180910390a28061015181610473565b9150506100f8565b60005b81518110156100f1577fb94ec34dfe32a8a7170992a093976368d1e63decf8f0bc0b38a8eb89cc9f95cf828281518110610198576101986104d3565b60200260200101516040516101ad919061038d565b60405180910390a1806101bf81610473565b91505061015c565b600060208083850312156101da57600080fd5b823567ffffffffffffffff808211156101f257600080fd5b8185019150601f868184011261020757600080fd5b823561021a6102158261044f565b610400565b8082825286820191508686018a888560051b890101111561023a57600080fd5b60005b848110156102e55781358781111561025457600080fd5b8801603f81018d1361026557600080fd5b8981013560408982111561027b5761027b610502565b6102aa8c7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08b85011601610400565b8281528f828486010111156102be57600080fd5b828285018e83013760009281018d019290925250855250928801929088019060010161023d565b50909a9950505050505050505050565b6000602080838503121561030857600080fd5b823567ffffffffffffffff81111561031f57600080fd5b8301601f8101851361033057600080fd5b803561033e6102158261044f565b80828252848201915084840188868560051b870101111561035e57600080fd5b600094505b83851015610381578035835260019490940193918501918501610363565b50979650505050505050565b600060208083528351808285015260005b818110156103ba5785810183015185820160400152820161039e565b818111156103cc576000604083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016929092016040019392505050565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171561044757610447610502565b604052919050565b600067ffffffffffffffff82111561046957610469610502565b5060051b60200190565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8214156104cc577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
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
	return address, tx, &LogEmitter{LogEmitterCaller: LogEmitterCaller{contract: contract}, LogEmitterTransactor: LogEmitterTransactor{contract: contract}, LogEmitterFilterer: LogEmitterFilterer{contract: contract}}, nil
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

func (_LogEmitter *LogEmitter) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _LogEmitter.abi.Events["Log1"].ID:
		return _LogEmitter.ParseLog1(log)
	case _LogEmitter.abi.Events["Log2"].ID:
		return _LogEmitter.ParseLog2(log)
	case _LogEmitter.abi.Events["Log3"].ID:
		return _LogEmitter.ParseLog3(log)

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

func (_LogEmitter *LogEmitter) Address() common.Address {
	return _LogEmitter.address
}

type LogEmitterInterface interface {
	EmitLog1(opts *bind.TransactOpts, v []*big.Int) (*types.Transaction, error)

	EmitLog2(opts *bind.TransactOpts, v []*big.Int) (*types.Transaction, error)

	EmitLog3(opts *bind.TransactOpts, v []string) (*types.Transaction, error)

	FilterLog1(opts *bind.FilterOpts) (*LogEmitterLog1Iterator, error)

	WatchLog1(opts *bind.WatchOpts, sink chan<- *LogEmitterLog1) (event.Subscription, error)

	ParseLog1(log types.Log) (*LogEmitterLog1, error)

	FilterLog2(opts *bind.FilterOpts, arg0 []*big.Int) (*LogEmitterLog2Iterator, error)

	WatchLog2(opts *bind.WatchOpts, sink chan<- *LogEmitterLog2, arg0 []*big.Int) (event.Subscription, error)

	ParseLog2(log types.Log) (*LogEmitterLog2, error)

	FilterLog3(opts *bind.FilterOpts) (*LogEmitterLog3Iterator, error)

	WatchLog3(opts *bind.WatchOpts, sink chan<- *LogEmitterLog3) (event.Subscription, error)

	ParseLog3(log types.Log) (*LogEmitterLog3, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
