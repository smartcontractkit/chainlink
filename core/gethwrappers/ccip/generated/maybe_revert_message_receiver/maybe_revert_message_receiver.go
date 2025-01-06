// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package maybe_revert_message_receiver

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

var MaybeRevertMessageReceiverMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"toRevert\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"ccipReceive\",\"inputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structClient.Any2EVMMessage\",\"components\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destTokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structClient.EVMTokenAmount[]\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"s_toRevert\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setErr\",\"inputs\":[{\"name\":\"err\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRevert\",\"inputs\":[{\"name\":\"toRevert\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"pure\"},{\"type\":\"event\",\"name\":\"MessageReceived\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValueReceived\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"CustomError\",\"inputs\":[{\"name\":\"err\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"ReceiveRevert\",\"inputs\":[]}]",
	Bin: "0x608034607d57601f6107a838819003918201601f19168301916001600160401b03831184841017608257808492602094604052833981010312607d5751801515809103607d57600080546001600160a81b0319163360ff60a01b19161760a09290921b60ff60a01b1691909117905560405161070f90816100998239f35b600080fd5b634e487b7160e01b600052604160045260246000fdfe608080604052600436101561007e575b50361561001b57600080fd5b60ff60005460a01c16610054577fe12e3b7047ff60a2dd763cf536a43597e5ce7fe7aa7476345bd4cd079912bcef6020604051348152a1005b7f3085b8db0000000000000000000000000000000000000000000000000000000060005260046000fd5b60003560e01c90816301ffc9a7146105f3575080635100fc21146105af57806377f5b0e6146102da57806385572ffb1461014d57638fb5f171146100c2573861000f565b346101485760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261014857600435801515809103610148577fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff74ff00000000000000000000000000000000000000006000549260a01b16911617600055600080f35b600080fd5b346101485760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101485760043567ffffffffffffffff8111610148577ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc60a091360301126101485760ff60005460a01c166101ee577fd82ce31e3523f6eeb2d24317b2b4133001e8472729657f663b68624c45f8f3e8600080a1005b6040517f5a4ff6710000000000000000000000000000000000000000000000000000000081526020600482015280600060015461022a816106af565b908160248501526001811690816000146102a2575060011461024b57500390fd5b6001600090815291507fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf65b81831061028857505081010360440190fd5b805460448487010152849350602090920191600101610276565b604493507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0091501682840152151560051b8201010390fd5b346101485760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126101485760043567ffffffffffffffff8111610148573660238201121561014857806004013567ffffffffffffffff811161058057604051917fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f81601f8501160116830183811067ffffffffffffffff82111761058057604052818352366024838301011161014857816000926024602093018386013783010152805167ffffffffffffffff8111610580576103be6001546106af565b601f81116104dd575b50602091601f821160011461042357918192600092610418575b50507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8260011b9260031b1c191617600155600080f35b0151905082806103e1565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe082169260016000527fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf69160005b8581106104c55750836001951061048e575b505050811b01600155005b01517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60f88460031b161c19169055828080610483565b91926020600181928685015181550194019201610471565b6001600052601f820160051c7fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6019060208310610558575b601f0160051c7fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf601905b81811061054c57506103c7565b6000815560010161053f565b7fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf69150610515565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b346101485760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261014857602060ff60005460a01c166040519015158152f35b346101485760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc36011261014857600435907fffffffff00000000000000000000000000000000000000000000000000000000821680920361014857817f85572ffb0000000000000000000000000000000000000000000000000000000060209314908115610685575b5015158152f35b7f01ffc9a7000000000000000000000000000000000000000000000000000000009150148361067e565b90600182811c921680156106f8575b60208310146106c957565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b91607f16916106be56fea164736f6c634300081a000a",
}

var MaybeRevertMessageReceiverABI = MaybeRevertMessageReceiverMetaData.ABI

var MaybeRevertMessageReceiverBin = MaybeRevertMessageReceiverMetaData.Bin

func DeployMaybeRevertMessageReceiver(auth *bind.TransactOpts, backend bind.ContractBackend, toRevert bool) (common.Address, *types.Transaction, *MaybeRevertMessageReceiver, error) {
	parsed, err := MaybeRevertMessageReceiverMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MaybeRevertMessageReceiverBin), backend, toRevert)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MaybeRevertMessageReceiver{address: address, abi: *parsed, MaybeRevertMessageReceiverCaller: MaybeRevertMessageReceiverCaller{contract: contract}, MaybeRevertMessageReceiverTransactor: MaybeRevertMessageReceiverTransactor{contract: contract}, MaybeRevertMessageReceiverFilterer: MaybeRevertMessageReceiverFilterer{contract: contract}}, nil
}

type MaybeRevertMessageReceiver struct {
	address common.Address
	abi     abi.ABI
	MaybeRevertMessageReceiverCaller
	MaybeRevertMessageReceiverTransactor
	MaybeRevertMessageReceiverFilterer
}

type MaybeRevertMessageReceiverCaller struct {
	contract *bind.BoundContract
}

type MaybeRevertMessageReceiverTransactor struct {
	contract *bind.BoundContract
}

type MaybeRevertMessageReceiverFilterer struct {
	contract *bind.BoundContract
}

type MaybeRevertMessageReceiverSession struct {
	Contract     *MaybeRevertMessageReceiver
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type MaybeRevertMessageReceiverCallerSession struct {
	Contract *MaybeRevertMessageReceiverCaller
	CallOpts bind.CallOpts
}

type MaybeRevertMessageReceiverTransactorSession struct {
	Contract     *MaybeRevertMessageReceiverTransactor
	TransactOpts bind.TransactOpts
}

type MaybeRevertMessageReceiverRaw struct {
	Contract *MaybeRevertMessageReceiver
}

type MaybeRevertMessageReceiverCallerRaw struct {
	Contract *MaybeRevertMessageReceiverCaller
}

type MaybeRevertMessageReceiverTransactorRaw struct {
	Contract *MaybeRevertMessageReceiverTransactor
}

func NewMaybeRevertMessageReceiver(address common.Address, backend bind.ContractBackend) (*MaybeRevertMessageReceiver, error) {
	abi, err := abi.JSON(strings.NewReader(MaybeRevertMessageReceiverABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindMaybeRevertMessageReceiver(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MaybeRevertMessageReceiver{address: address, abi: abi, MaybeRevertMessageReceiverCaller: MaybeRevertMessageReceiverCaller{contract: contract}, MaybeRevertMessageReceiverTransactor: MaybeRevertMessageReceiverTransactor{contract: contract}, MaybeRevertMessageReceiverFilterer: MaybeRevertMessageReceiverFilterer{contract: contract}}, nil
}

func NewMaybeRevertMessageReceiverCaller(address common.Address, caller bind.ContractCaller) (*MaybeRevertMessageReceiverCaller, error) {
	contract, err := bindMaybeRevertMessageReceiver(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MaybeRevertMessageReceiverCaller{contract: contract}, nil
}

func NewMaybeRevertMessageReceiverTransactor(address common.Address, transactor bind.ContractTransactor) (*MaybeRevertMessageReceiverTransactor, error) {
	contract, err := bindMaybeRevertMessageReceiver(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MaybeRevertMessageReceiverTransactor{contract: contract}, nil
}

func NewMaybeRevertMessageReceiverFilterer(address common.Address, filterer bind.ContractFilterer) (*MaybeRevertMessageReceiverFilterer, error) {
	contract, err := bindMaybeRevertMessageReceiver(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MaybeRevertMessageReceiverFilterer{contract: contract}, nil
}

func bindMaybeRevertMessageReceiver(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MaybeRevertMessageReceiverMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MaybeRevertMessageReceiver.Contract.MaybeRevertMessageReceiverCaller.contract.Call(opts, result, method, params...)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MaybeRevertMessageReceiver.Contract.MaybeRevertMessageReceiverTransactor.contract.Transfer(opts)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MaybeRevertMessageReceiver.Contract.MaybeRevertMessageReceiverTransactor.contract.Transact(opts, method, params...)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MaybeRevertMessageReceiver.Contract.contract.Call(opts, result, method, params...)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MaybeRevertMessageReceiver.Contract.contract.Transfer(opts)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MaybeRevertMessageReceiver.Contract.contract.Transact(opts, method, params...)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverCaller) SToRevert(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _MaybeRevertMessageReceiver.contract.Call(opts, &out, "s_toRevert")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverSession) SToRevert() (bool, error) {
	return _MaybeRevertMessageReceiver.Contract.SToRevert(&_MaybeRevertMessageReceiver.CallOpts)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverCallerSession) SToRevert() (bool, error) {
	return _MaybeRevertMessageReceiver.Contract.SToRevert(&_MaybeRevertMessageReceiver.CallOpts)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _MaybeRevertMessageReceiver.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _MaybeRevertMessageReceiver.Contract.SupportsInterface(&_MaybeRevertMessageReceiver.CallOpts, interfaceId)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _MaybeRevertMessageReceiver.Contract.SupportsInterface(&_MaybeRevertMessageReceiver.CallOpts, interfaceId)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverTransactor) CcipReceive(opts *bind.TransactOpts, arg0 ClientAny2EVMMessage) (*types.Transaction, error) {
	return _MaybeRevertMessageReceiver.contract.Transact(opts, "ccipReceive", arg0)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverSession) CcipReceive(arg0 ClientAny2EVMMessage) (*types.Transaction, error) {
	return _MaybeRevertMessageReceiver.Contract.CcipReceive(&_MaybeRevertMessageReceiver.TransactOpts, arg0)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverTransactorSession) CcipReceive(arg0 ClientAny2EVMMessage) (*types.Transaction, error) {
	return _MaybeRevertMessageReceiver.Contract.CcipReceive(&_MaybeRevertMessageReceiver.TransactOpts, arg0)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverTransactor) SetErr(opts *bind.TransactOpts, err []byte) (*types.Transaction, error) {
	return _MaybeRevertMessageReceiver.contract.Transact(opts, "setErr", err)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverSession) SetErr(err []byte) (*types.Transaction, error) {
	return _MaybeRevertMessageReceiver.Contract.SetErr(&_MaybeRevertMessageReceiver.TransactOpts, err)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverTransactorSession) SetErr(err []byte) (*types.Transaction, error) {
	return _MaybeRevertMessageReceiver.Contract.SetErr(&_MaybeRevertMessageReceiver.TransactOpts, err)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverTransactor) SetRevert(opts *bind.TransactOpts, toRevert bool) (*types.Transaction, error) {
	return _MaybeRevertMessageReceiver.contract.Transact(opts, "setRevert", toRevert)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverSession) SetRevert(toRevert bool) (*types.Transaction, error) {
	return _MaybeRevertMessageReceiver.Contract.SetRevert(&_MaybeRevertMessageReceiver.TransactOpts, toRevert)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverTransactorSession) SetRevert(toRevert bool) (*types.Transaction, error) {
	return _MaybeRevertMessageReceiver.Contract.SetRevert(&_MaybeRevertMessageReceiver.TransactOpts, toRevert)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MaybeRevertMessageReceiver.contract.RawTransact(opts, nil)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverSession) Receive() (*types.Transaction, error) {
	return _MaybeRevertMessageReceiver.Contract.Receive(&_MaybeRevertMessageReceiver.TransactOpts)
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverTransactorSession) Receive() (*types.Transaction, error) {
	return _MaybeRevertMessageReceiver.Contract.Receive(&_MaybeRevertMessageReceiver.TransactOpts)
}

type MaybeRevertMessageReceiverMessageReceivedIterator struct {
	Event *MaybeRevertMessageReceiverMessageReceived

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MaybeRevertMessageReceiverMessageReceivedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MaybeRevertMessageReceiverMessageReceived)
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
		it.Event = new(MaybeRevertMessageReceiverMessageReceived)
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

func (it *MaybeRevertMessageReceiverMessageReceivedIterator) Error() error {
	return it.fail
}

func (it *MaybeRevertMessageReceiverMessageReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MaybeRevertMessageReceiverMessageReceived struct {
	Raw types.Log
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverFilterer) FilterMessageReceived(opts *bind.FilterOpts) (*MaybeRevertMessageReceiverMessageReceivedIterator, error) {

	logs, sub, err := _MaybeRevertMessageReceiver.contract.FilterLogs(opts, "MessageReceived")
	if err != nil {
		return nil, err
	}
	return &MaybeRevertMessageReceiverMessageReceivedIterator{contract: _MaybeRevertMessageReceiver.contract, event: "MessageReceived", logs: logs, sub: sub}, nil
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverFilterer) WatchMessageReceived(opts *bind.WatchOpts, sink chan<- *MaybeRevertMessageReceiverMessageReceived) (event.Subscription, error) {

	logs, sub, err := _MaybeRevertMessageReceiver.contract.WatchLogs(opts, "MessageReceived")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MaybeRevertMessageReceiverMessageReceived)
				if err := _MaybeRevertMessageReceiver.contract.UnpackLog(event, "MessageReceived", log); err != nil {
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

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverFilterer) ParseMessageReceived(log types.Log) (*MaybeRevertMessageReceiverMessageReceived, error) {
	event := new(MaybeRevertMessageReceiverMessageReceived)
	if err := _MaybeRevertMessageReceiver.contract.UnpackLog(event, "MessageReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type MaybeRevertMessageReceiverValueReceivedIterator struct {
	Event *MaybeRevertMessageReceiverValueReceived

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *MaybeRevertMessageReceiverValueReceivedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MaybeRevertMessageReceiverValueReceived)
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
		it.Event = new(MaybeRevertMessageReceiverValueReceived)
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

func (it *MaybeRevertMessageReceiverValueReceivedIterator) Error() error {
	return it.fail
}

func (it *MaybeRevertMessageReceiverValueReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type MaybeRevertMessageReceiverValueReceived struct {
	Amount *big.Int
	Raw    types.Log
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverFilterer) FilterValueReceived(opts *bind.FilterOpts) (*MaybeRevertMessageReceiverValueReceivedIterator, error) {

	logs, sub, err := _MaybeRevertMessageReceiver.contract.FilterLogs(opts, "ValueReceived")
	if err != nil {
		return nil, err
	}
	return &MaybeRevertMessageReceiverValueReceivedIterator{contract: _MaybeRevertMessageReceiver.contract, event: "ValueReceived", logs: logs, sub: sub}, nil
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverFilterer) WatchValueReceived(opts *bind.WatchOpts, sink chan<- *MaybeRevertMessageReceiverValueReceived) (event.Subscription, error) {

	logs, sub, err := _MaybeRevertMessageReceiver.contract.WatchLogs(opts, "ValueReceived")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(MaybeRevertMessageReceiverValueReceived)
				if err := _MaybeRevertMessageReceiver.contract.UnpackLog(event, "ValueReceived", log); err != nil {
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

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiverFilterer) ParseValueReceived(log types.Log) (*MaybeRevertMessageReceiverValueReceived, error) {
	event := new(MaybeRevertMessageReceiverValueReceived)
	if err := _MaybeRevertMessageReceiver.contract.UnpackLog(event, "ValueReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiver) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _MaybeRevertMessageReceiver.abi.Events["MessageReceived"].ID:
		return _MaybeRevertMessageReceiver.ParseMessageReceived(log)
	case _MaybeRevertMessageReceiver.abi.Events["ValueReceived"].ID:
		return _MaybeRevertMessageReceiver.ParseValueReceived(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (MaybeRevertMessageReceiverMessageReceived) Topic() common.Hash {
	return common.HexToHash("0xd82ce31e3523f6eeb2d24317b2b4133001e8472729657f663b68624c45f8f3e8")
}

func (MaybeRevertMessageReceiverValueReceived) Topic() common.Hash {
	return common.HexToHash("0xe12e3b7047ff60a2dd763cf536a43597e5ce7fe7aa7476345bd4cd079912bcef")
}

func (_MaybeRevertMessageReceiver *MaybeRevertMessageReceiver) Address() common.Address {
	return _MaybeRevertMessageReceiver.address
}

type MaybeRevertMessageReceiverInterface interface {
	SToRevert(opts *bind.CallOpts) (bool, error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	CcipReceive(opts *bind.TransactOpts, arg0 ClientAny2EVMMessage) (*types.Transaction, error)

	SetErr(opts *bind.TransactOpts, err []byte) (*types.Transaction, error)

	SetRevert(opts *bind.TransactOpts, toRevert bool) (*types.Transaction, error)

	Receive(opts *bind.TransactOpts) (*types.Transaction, error)

	FilterMessageReceived(opts *bind.FilterOpts) (*MaybeRevertMessageReceiverMessageReceivedIterator, error)

	WatchMessageReceived(opts *bind.WatchOpts, sink chan<- *MaybeRevertMessageReceiverMessageReceived) (event.Subscription, error)

	ParseMessageReceived(log types.Log) (*MaybeRevertMessageReceiverMessageReceived, error)

	FilterValueReceived(opts *bind.FilterOpts) (*MaybeRevertMessageReceiverValueReceivedIterator, error)

	WatchValueReceived(opts *bind.WatchOpts, sink chan<- *MaybeRevertMessageReceiverValueReceived) (event.Subscription, error)

	ParseValueReceived(log types.Log) (*MaybeRevertMessageReceiverValueReceived, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
