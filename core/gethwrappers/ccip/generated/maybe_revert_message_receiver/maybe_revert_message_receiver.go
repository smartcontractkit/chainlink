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
	ABI: "[{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"toRevert\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"CustomError\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReceiveRevert\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"MessageReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"ValueReceived\",\"type\":\"event\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"messageId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sourceChainSelector\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structClient.EVMTokenAmount[]\",\"name\":\"destTokenAmounts\",\"type\":\"tuple[]\"}],\"internalType\":\"structClient.Any2EVMMessage\",\"name\":\"\",\"type\":\"tuple\"}],\"name\":\"ccipReceive\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_toRevert\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"setErr\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"toRevert\",\"type\":\"bool\"}],\"name\":\"setRevert\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x608060405234801561001057600080fd5b506040516107e73803806107e783398101604081905261002f9161005d565b600080546001600160a81b0319163360ff60a01b191617600160a01b92151592909202919091179055610086565b60006020828403121561006f57600080fd5b8151801515811461007f57600080fd5b9392505050565b610752806100956000396000f3fe60806040526004361061005e5760003560e01c806377f5b0e61161004357806377f5b0e61461015857806385572ffb1461017a5780638fb5f1711461019a57600080fd5b806301ffc9a7146100f25780635100fc211461012657600080fd5b366100ed5760005474010000000000000000000000000000000000000000900460ff16156100b8576040517f3085b8db00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040513481527fe12e3b7047ff60a2dd763cf536a43597e5ce7fe7aa7476345bd4cd079912bcef9060200160405180910390a1005b600080fd5b3480156100fe57600080fd5b5061011261010d366004610335565b6101ff565b604051901515815260200160405180910390f35b34801561013257600080fd5b506000546101129074010000000000000000000000000000000000000000900460ff1681565b34801561016457600080fd5b506101786101733660046103ad565b610298565b005b34801561018657600080fd5b5061017861019536600461047c565b6102a8565b3480156101a657600080fd5b506101786101b53660046104b7565b6000805491151574010000000000000000000000000000000000000000027fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff909216919091179055565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167f85572ffb00000000000000000000000000000000000000000000000000000000148061029257507fffffffff0000000000000000000000000000000000000000000000000000000082167f01ffc9a700000000000000000000000000000000000000000000000000000000145b92915050565b60016102a4828261057d565b5050565b60005474010000000000000000000000000000000000000000900460ff16156103095760016040517f5a4ff6710000000000000000000000000000000000000000000000000000000081526004016103009190610697565b60405180910390fd5b6040517fd82ce31e3523f6eeb2d24317b2b4133001e8472729657f663b68624c45f8f3e890600090a150565b60006020828403121561034757600080fd5b81357fffffffff000000000000000000000000000000000000000000000000000000008116811461037757600080fd5b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6000602082840312156103bf57600080fd5b813567ffffffffffffffff808211156103d757600080fd5b818401915084601f8301126103eb57600080fd5b8135818111156103fd576103fd61037e565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f011681019083821181831017156104435761044361037e565b8160405282815287602084870101111561045c57600080fd5b826020860160208301376000928101602001929092525095945050505050565b60006020828403121561048e57600080fd5b813567ffffffffffffffff8111156104a557600080fd5b820160a0818503121561037757600080fd5b6000602082840312156104c957600080fd5b8135801515811461037757600080fd5b600181811c908216806104ed57607f821691505b602082108103610526577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f821115610578576000816000526020600020601f850160051c810160208610156105555750805b601f850160051c820191505b8181101561057457828155600101610561565b5050505b505050565b815167ffffffffffffffff8111156105975761059761037e565b6105ab816105a584546104d9565b8461052c565b602080601f8311600181146105fe57600084156105c85750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b178555610574565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b8281101561064b5788860151825594840194600190910190840161062c565b508582101561068757878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b60006020808352600084546106ab816104d9565b80602087015260406001808416600081146106cd576001811461070757610737565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00851660408a0152604084151560051b8a01019550610737565b89600052602060002060005b8581101561072e5781548b8201860152908301908801610713565b8a016040019650505b50939897505050505050505056fea164736f6c6343000818000a",
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
