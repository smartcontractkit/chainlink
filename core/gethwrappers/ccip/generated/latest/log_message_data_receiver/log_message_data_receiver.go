// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package log_message_data_receiver

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

type ClientAny2EVMMessage struct {
	MessageId           [32]byte
	SourceChainSelector uint64
	Sender              []byte
	Data                []byte
	DestTokenAmounts    []ClientEVMTokenAmount
}

type ClientEVMTokenAmount struct {
	Token  common.Address
	Amount *big.Int
}

var LogMessageDataReceiverMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"ccipReceive\",\"inputs\":[{\"name\":\"message\",\"type\":\"tuple\",\"internalType\":\"structClient.Any2EVMMessage\",\"components\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destTokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structClient.EVMTokenAmount[]\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"typeAndVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"MessageReceived\",\"inputs\":[{\"name\":\"data\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false}]",
	Bin: "0x6080806040523460155761034e908161001b8239f35b600080fdfe608080604052600436101561001357600080fd5b60003560e01c90816301ffc9a71461028557508063181f5a7714610173576385572ffb1461004057600080fd5b3461016e5760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261016e5760043567ffffffffffffffff811161016e578036039060a07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc83011261016e577fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffdd6064820135920182121561016e570160048101359067ffffffffffffffff821161016e5760240190803603821361016e5760407f4b3be2c5d6fcecc68e42c0268adeed4f145b0b4a6cbd5960dcdd39867bef682f927fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8484519586946020865281602087015286860137600085828601015201168101030190a1005b600080fd5b3461016e5760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261016e576040516040810181811067ffffffffffffffff82111761025657604052601c81527f4c6f674d65737361676544617461526563656976657220312e302e3000000000602082015260405190602082528181519182602083015260005b83811061023e5750507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f836000604080968601015201168101030190f35b602082820181015160408784010152859350016101fe565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b3461016e5760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261016e57600435907fffffffff00000000000000000000000000000000000000000000000000000000821680920361016e57817f85572ffb0000000000000000000000000000000000000000000000000000000060209314908115610317575b5015158152f35b7f01ffc9a7000000000000000000000000000000000000000000000000000000009150148361031056fea164736f6c634300081a000a",
}

var LogMessageDataReceiverABI = LogMessageDataReceiverMetaData.ABI

var LogMessageDataReceiverBin = LogMessageDataReceiverMetaData.Bin

func DeployLogMessageDataReceiver(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *LogMessageDataReceiver, error) {
	parsed, err := LogMessageDataReceiverMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(LogMessageDataReceiverBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &LogMessageDataReceiver{address: address, abi: *parsed, LogMessageDataReceiverCaller: LogMessageDataReceiverCaller{contract: contract}, LogMessageDataReceiverTransactor: LogMessageDataReceiverTransactor{contract: contract}, LogMessageDataReceiverFilterer: LogMessageDataReceiverFilterer{contract: contract}}, nil
}

type LogMessageDataReceiver struct {
	address common.Address
	abi     abi.ABI
	LogMessageDataReceiverCaller
	LogMessageDataReceiverTransactor
	LogMessageDataReceiverFilterer
}

type LogMessageDataReceiverCaller struct {
	contract *bind.BoundContract
}

type LogMessageDataReceiverTransactor struct {
	contract *bind.BoundContract
}

type LogMessageDataReceiverFilterer struct {
	contract *bind.BoundContract
}

type LogMessageDataReceiverSession struct {
	Contract     *LogMessageDataReceiver
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type LogMessageDataReceiverCallerSession struct {
	Contract *LogMessageDataReceiverCaller
	CallOpts bind.CallOpts
}

type LogMessageDataReceiverTransactorSession struct {
	Contract     *LogMessageDataReceiverTransactor
	TransactOpts bind.TransactOpts
}

type LogMessageDataReceiverRaw struct {
	Contract *LogMessageDataReceiver
}

type LogMessageDataReceiverCallerRaw struct {
	Contract *LogMessageDataReceiverCaller
}

type LogMessageDataReceiverTransactorRaw struct {
	Contract *LogMessageDataReceiverTransactor
}

func NewLogMessageDataReceiver(address common.Address, backend bind.ContractBackend) (*LogMessageDataReceiver, error) {
	abi, err := abi.JSON(strings.NewReader(LogMessageDataReceiverABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindLogMessageDataReceiver(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LogMessageDataReceiver{address: address, abi: abi, LogMessageDataReceiverCaller: LogMessageDataReceiverCaller{contract: contract}, LogMessageDataReceiverTransactor: LogMessageDataReceiverTransactor{contract: contract}, LogMessageDataReceiverFilterer: LogMessageDataReceiverFilterer{contract: contract}}, nil
}

func NewLogMessageDataReceiverCaller(address common.Address, caller bind.ContractCaller) (*LogMessageDataReceiverCaller, error) {
	contract, err := bindLogMessageDataReceiver(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LogMessageDataReceiverCaller{contract: contract}, nil
}

func NewLogMessageDataReceiverTransactor(address common.Address, transactor bind.ContractTransactor) (*LogMessageDataReceiverTransactor, error) {
	contract, err := bindLogMessageDataReceiver(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LogMessageDataReceiverTransactor{contract: contract}, nil
}

func NewLogMessageDataReceiverFilterer(address common.Address, filterer bind.ContractFilterer) (*LogMessageDataReceiverFilterer, error) {
	contract, err := bindLogMessageDataReceiver(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LogMessageDataReceiverFilterer{contract: contract}, nil
}

func bindLogMessageDataReceiver(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := LogMessageDataReceiverMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_LogMessageDataReceiver *LogMessageDataReceiverRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LogMessageDataReceiver.Contract.LogMessageDataReceiverCaller.contract.Call(opts, result, method, params...)
}

func (_LogMessageDataReceiver *LogMessageDataReceiverRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LogMessageDataReceiver.Contract.LogMessageDataReceiverTransactor.contract.Transfer(opts)
}

func (_LogMessageDataReceiver *LogMessageDataReceiverRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LogMessageDataReceiver.Contract.LogMessageDataReceiverTransactor.contract.Transact(opts, method, params...)
}

func (_LogMessageDataReceiver *LogMessageDataReceiverCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LogMessageDataReceiver.Contract.contract.Call(opts, result, method, params...)
}

func (_LogMessageDataReceiver *LogMessageDataReceiverTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LogMessageDataReceiver.Contract.contract.Transfer(opts)
}

func (_LogMessageDataReceiver *LogMessageDataReceiverTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LogMessageDataReceiver.Contract.contract.Transact(opts, method, params...)
}

func (_LogMessageDataReceiver *LogMessageDataReceiverCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _LogMessageDataReceiver.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_LogMessageDataReceiver *LogMessageDataReceiverSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _LogMessageDataReceiver.Contract.SupportsInterface(&_LogMessageDataReceiver.CallOpts, interfaceId)
}

func (_LogMessageDataReceiver *LogMessageDataReceiverCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _LogMessageDataReceiver.Contract.SupportsInterface(&_LogMessageDataReceiver.CallOpts, interfaceId)
}

func (_LogMessageDataReceiver *LogMessageDataReceiverCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _LogMessageDataReceiver.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_LogMessageDataReceiver *LogMessageDataReceiverSession) TypeAndVersion() (string, error) {
	return _LogMessageDataReceiver.Contract.TypeAndVersion(&_LogMessageDataReceiver.CallOpts)
}

func (_LogMessageDataReceiver *LogMessageDataReceiverCallerSession) TypeAndVersion() (string, error) {
	return _LogMessageDataReceiver.Contract.TypeAndVersion(&_LogMessageDataReceiver.CallOpts)
}

func (_LogMessageDataReceiver *LogMessageDataReceiverTransactor) CcipReceive(opts *bind.TransactOpts, message ClientAny2EVMMessage) (*types.Transaction, error) {
	return _LogMessageDataReceiver.contract.Transact(opts, "ccipReceive", message)
}

func (_LogMessageDataReceiver *LogMessageDataReceiverSession) CcipReceive(message ClientAny2EVMMessage) (*types.Transaction, error) {
	return _LogMessageDataReceiver.Contract.CcipReceive(&_LogMessageDataReceiver.TransactOpts, message)
}

func (_LogMessageDataReceiver *LogMessageDataReceiverTransactorSession) CcipReceive(message ClientAny2EVMMessage) (*types.Transaction, error) {
	return _LogMessageDataReceiver.Contract.CcipReceive(&_LogMessageDataReceiver.TransactOpts, message)
}

type LogMessageDataReceiverMessageReceivedIterator struct {
	Event *LogMessageDataReceiverMessageReceived

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *LogMessageDataReceiverMessageReceivedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LogMessageDataReceiverMessageReceived)
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
		it.Event = new(LogMessageDataReceiverMessageReceived)
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

func (it *LogMessageDataReceiverMessageReceivedIterator) Error() error {
	return it.fail
}

func (it *LogMessageDataReceiverMessageReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type LogMessageDataReceiverMessageReceived struct {
	Data []byte
	Raw  types.Log
}

func (_LogMessageDataReceiver *LogMessageDataReceiverFilterer) FilterMessageReceived(opts *bind.FilterOpts) (*LogMessageDataReceiverMessageReceivedIterator, error) {

	logs, sub, err := _LogMessageDataReceiver.contract.FilterLogs(opts, "MessageReceived")
	if err != nil {
		return nil, err
	}
	return &LogMessageDataReceiverMessageReceivedIterator{contract: _LogMessageDataReceiver.contract, event: "MessageReceived", logs: logs, sub: sub}, nil
}

func (_LogMessageDataReceiver *LogMessageDataReceiverFilterer) WatchMessageReceived(opts *bind.WatchOpts, sink chan<- *LogMessageDataReceiverMessageReceived) (event.Subscription, error) {

	logs, sub, err := _LogMessageDataReceiver.contract.WatchLogs(opts, "MessageReceived")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(LogMessageDataReceiverMessageReceived)
				if err := _LogMessageDataReceiver.contract.UnpackLog(event, "MessageReceived", log); err != nil {
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

func (_LogMessageDataReceiver *LogMessageDataReceiverFilterer) ParseMessageReceived(log types.Log) (*LogMessageDataReceiverMessageReceived, error) {
	event := new(LogMessageDataReceiverMessageReceived)
	if err := _LogMessageDataReceiver.contract.UnpackLog(event, "MessageReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_LogMessageDataReceiver *LogMessageDataReceiver) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _LogMessageDataReceiver.abi.Events["MessageReceived"].ID:
		return _LogMessageDataReceiver.ParseMessageReceived(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (LogMessageDataReceiverMessageReceived) Topic() common.Hash {
	return common.HexToHash("0x4b3be2c5d6fcecc68e42c0268adeed4f145b0b4a6cbd5960dcdd39867bef682f")
}

func (_LogMessageDataReceiver *LogMessageDataReceiver) Address() common.Address {
	return _LogMessageDataReceiver.address
}

type LogMessageDataReceiverInterface interface {
	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	CcipReceive(opts *bind.TransactOpts, message ClientAny2EVMMessage) (*types.Transaction, error)

	FilterMessageReceived(opts *bind.FilterOpts) (*LogMessageDataReceiverMessageReceivedIterator, error)

	WatchMessageReceived(opts *bind.WatchOpts, sink chan<- *LogMessageDataReceiverMessageReceived) (event.Subscription, error)

	ParseMessageReceived(log types.Log) (*LogMessageDataReceiverMessageReceived, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
