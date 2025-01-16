// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ccip_dummy_receiver

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

var CCIPDummyReceiverMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"ccipReceive\",\"inputs\":[{\"name\":\"message\",\"type\":\"tuple\",\"internalType\":\"structClient.Any2EVMMessage\",\"components\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destTokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structClient.EVMTokenAmount[]\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getRouter\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"MessageReceived\",\"inputs\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"data\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"InvalidRouter\",\"inputs\":[{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
	Bin: "0x60a034608757601f61068d38819003918201601f19168301916001600160401b03831184841017608c57808492602094604052833981010312608757516001600160a01b038116808203608757156071576080526040516105ea90816100a38239608051818181608b01526101380152f35b6335fdcccd60e21b600052600060045260246000fd5b600080fd5b634e487b7160e01b600052604160045260246000fdfe608080604052600436101561001357600080fd5b60003560e01c90816301ffc9a71461043d5750806385572ffb146100b45763b0f479a11461004057600080fd5b346100af5760007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126100af57602060405173ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000168152f35b600080fd5b346100af5760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126100af5760043567ffffffffffffffff81116100af5760a07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc82360301126100af5773ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016330361040f57600060405160a0810181811067ffffffffffffffff8211176103e2576040528260040135815260248301359067ffffffffffffffff821682036103de5760208101918252604484013567ffffffffffffffff81116103d6576101c6906004369187010161056c565b6040820152606484013567ffffffffffffffff81116103d6576101ef906004369187010161056c565b936060820194855260848101359067ffffffffffffffff82116103da5701366023820112156103d657600481013567ffffffffffffffff81116103a95761023b60208260051b016104f9565b91602060048185858152019360061b83010101903682116103a557602401915b818310610312575050509067ffffffffffffffff91608082015251915116925192604051918252602082015260606040820152825192836060830152825b8481106102fc57837fd1e72bffc6abde0f1ef309ad9774dafcf9d29fe059a3fe5bb322ae9f15a74423846080817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8b8785828601015201168101030190a180f35b8060208092840101516080828601015201610299565b6040833603126103a5576040516040810181811067ffffffffffffffff82111761037857604052833573ffffffffffffffffffffffffffffffffffffffff8116810361037457918160409360209352828601358382015281520192019161025b565b8880fd5b6024897f4e487b710000000000000000000000000000000000000000000000000000000081526041600452fd5b8680fd5b6024857f4e487b710000000000000000000000000000000000000000000000000000000081526041600452fd5b8380fd5b8480fd5b8280fd5b6024837f4e487b710000000000000000000000000000000000000000000000000000000081526041600452fd5b7fd7f73334000000000000000000000000000000000000000000000000000000006000523360045260246000fd5b346100af5760207ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126100af57600435907fffffffff0000000000000000000000000000000000000000000000000000000082168092036100af57817f85572ffb00000000000000000000000000000000000000000000000000000000602093149081156104cf575b5015158152f35b7f01ffc9a700000000000000000000000000000000000000000000000000000000915014836104c8565b907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f604051930116820182811067ffffffffffffffff82111761053d57604052565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b81601f820112156100af5780359067ffffffffffffffff821161053d576105ba60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f850116016104f9565b92828452602083830101116100af5781600092602080930183860137830101529056fea164736f6c634300081a000a",
}

var CCIPDummyReceiverABI = CCIPDummyReceiverMetaData.ABI

var CCIPDummyReceiverBin = CCIPDummyReceiverMetaData.Bin

func DeployCCIPDummyReceiver(auth *bind.TransactOpts, backend bind.ContractBackend, router common.Address) (common.Address, *types.Transaction, *CCIPDummyReceiver, error) {
	parsed, err := CCIPDummyReceiverMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CCIPDummyReceiverBin), backend, router)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &CCIPDummyReceiver{address: address, abi: *parsed, CCIPDummyReceiverCaller: CCIPDummyReceiverCaller{contract: contract}, CCIPDummyReceiverTransactor: CCIPDummyReceiverTransactor{contract: contract}, CCIPDummyReceiverFilterer: CCIPDummyReceiverFilterer{contract: contract}}, nil
}

type CCIPDummyReceiver struct {
	address common.Address
	abi     abi.ABI
	CCIPDummyReceiverCaller
	CCIPDummyReceiverTransactor
	CCIPDummyReceiverFilterer
}

type CCIPDummyReceiverCaller struct {
	contract *bind.BoundContract
}

type CCIPDummyReceiverTransactor struct {
	contract *bind.BoundContract
}

type CCIPDummyReceiverFilterer struct {
	contract *bind.BoundContract
}

type CCIPDummyReceiverSession struct {
	Contract     *CCIPDummyReceiver
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type CCIPDummyReceiverCallerSession struct {
	Contract *CCIPDummyReceiverCaller
	CallOpts bind.CallOpts
}

type CCIPDummyReceiverTransactorSession struct {
	Contract     *CCIPDummyReceiverTransactor
	TransactOpts bind.TransactOpts
}

type CCIPDummyReceiverRaw struct {
	Contract *CCIPDummyReceiver
}

type CCIPDummyReceiverCallerRaw struct {
	Contract *CCIPDummyReceiverCaller
}

type CCIPDummyReceiverTransactorRaw struct {
	Contract *CCIPDummyReceiverTransactor
}

func NewCCIPDummyReceiver(address common.Address, backend bind.ContractBackend) (*CCIPDummyReceiver, error) {
	abi, err := abi.JSON(strings.NewReader(CCIPDummyReceiverABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindCCIPDummyReceiver(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CCIPDummyReceiver{address: address, abi: abi, CCIPDummyReceiverCaller: CCIPDummyReceiverCaller{contract: contract}, CCIPDummyReceiverTransactor: CCIPDummyReceiverTransactor{contract: contract}, CCIPDummyReceiverFilterer: CCIPDummyReceiverFilterer{contract: contract}}, nil
}

func NewCCIPDummyReceiverCaller(address common.Address, caller bind.ContractCaller) (*CCIPDummyReceiverCaller, error) {
	contract, err := bindCCIPDummyReceiver(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CCIPDummyReceiverCaller{contract: contract}, nil
}

func NewCCIPDummyReceiverTransactor(address common.Address, transactor bind.ContractTransactor) (*CCIPDummyReceiverTransactor, error) {
	contract, err := bindCCIPDummyReceiver(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CCIPDummyReceiverTransactor{contract: contract}, nil
}

func NewCCIPDummyReceiverFilterer(address common.Address, filterer bind.ContractFilterer) (*CCIPDummyReceiverFilterer, error) {
	contract, err := bindCCIPDummyReceiver(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CCIPDummyReceiverFilterer{contract: contract}, nil
}

func bindCCIPDummyReceiver(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CCIPDummyReceiverMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_CCIPDummyReceiver *CCIPDummyReceiverRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CCIPDummyReceiver.Contract.CCIPDummyReceiverCaller.contract.Call(opts, result, method, params...)
}

func (_CCIPDummyReceiver *CCIPDummyReceiverRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CCIPDummyReceiver.Contract.CCIPDummyReceiverTransactor.contract.Transfer(opts)
}

func (_CCIPDummyReceiver *CCIPDummyReceiverRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CCIPDummyReceiver.Contract.CCIPDummyReceiverTransactor.contract.Transact(opts, method, params...)
}

func (_CCIPDummyReceiver *CCIPDummyReceiverCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CCIPDummyReceiver.Contract.contract.Call(opts, result, method, params...)
}

func (_CCIPDummyReceiver *CCIPDummyReceiverTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CCIPDummyReceiver.Contract.contract.Transfer(opts)
}

func (_CCIPDummyReceiver *CCIPDummyReceiverTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CCIPDummyReceiver.Contract.contract.Transact(opts, method, params...)
}

func (_CCIPDummyReceiver *CCIPDummyReceiverCaller) GetRouter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CCIPDummyReceiver.contract.Call(opts, &out, "getRouter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_CCIPDummyReceiver *CCIPDummyReceiverSession) GetRouter() (common.Address, error) {
	return _CCIPDummyReceiver.Contract.GetRouter(&_CCIPDummyReceiver.CallOpts)
}

func (_CCIPDummyReceiver *CCIPDummyReceiverCallerSession) GetRouter() (common.Address, error) {
	return _CCIPDummyReceiver.Contract.GetRouter(&_CCIPDummyReceiver.CallOpts)
}

func (_CCIPDummyReceiver *CCIPDummyReceiverCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _CCIPDummyReceiver.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_CCIPDummyReceiver *CCIPDummyReceiverSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _CCIPDummyReceiver.Contract.SupportsInterface(&_CCIPDummyReceiver.CallOpts, interfaceId)
}

func (_CCIPDummyReceiver *CCIPDummyReceiverCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _CCIPDummyReceiver.Contract.SupportsInterface(&_CCIPDummyReceiver.CallOpts, interfaceId)
}

func (_CCIPDummyReceiver *CCIPDummyReceiverTransactor) CcipReceive(opts *bind.TransactOpts, message ClientAny2EVMMessage) (*types.Transaction, error) {
	return _CCIPDummyReceiver.contract.Transact(opts, "ccipReceive", message)
}

func (_CCIPDummyReceiver *CCIPDummyReceiverSession) CcipReceive(message ClientAny2EVMMessage) (*types.Transaction, error) {
	return _CCIPDummyReceiver.Contract.CcipReceive(&_CCIPDummyReceiver.TransactOpts, message)
}

func (_CCIPDummyReceiver *CCIPDummyReceiverTransactorSession) CcipReceive(message ClientAny2EVMMessage) (*types.Transaction, error) {
	return _CCIPDummyReceiver.Contract.CcipReceive(&_CCIPDummyReceiver.TransactOpts, message)
}

type CCIPDummyReceiverMessageReceivedIterator struct {
	Event *CCIPDummyReceiverMessageReceived

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *CCIPDummyReceiverMessageReceivedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CCIPDummyReceiverMessageReceived)
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
		it.Event = new(CCIPDummyReceiverMessageReceived)
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

func (it *CCIPDummyReceiverMessageReceivedIterator) Error() error {
	return it.fail
}

func (it *CCIPDummyReceiverMessageReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type CCIPDummyReceiverMessageReceived struct {
	MessageId           [32]byte
	SourceChainSelector uint64
	Data                []byte
	Raw                 types.Log
}

func (_CCIPDummyReceiver *CCIPDummyReceiverFilterer) FilterMessageReceived(opts *bind.FilterOpts) (*CCIPDummyReceiverMessageReceivedIterator, error) {

	logs, sub, err := _CCIPDummyReceiver.contract.FilterLogs(opts, "MessageReceived")
	if err != nil {
		return nil, err
	}
	return &CCIPDummyReceiverMessageReceivedIterator{contract: _CCIPDummyReceiver.contract, event: "MessageReceived", logs: logs, sub: sub}, nil
}

func (_CCIPDummyReceiver *CCIPDummyReceiverFilterer) WatchMessageReceived(opts *bind.WatchOpts, sink chan<- *CCIPDummyReceiverMessageReceived) (event.Subscription, error) {

	logs, sub, err := _CCIPDummyReceiver.contract.WatchLogs(opts, "MessageReceived")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(CCIPDummyReceiverMessageReceived)
				if err := _CCIPDummyReceiver.contract.UnpackLog(event, "MessageReceived", log); err != nil {
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

func (_CCIPDummyReceiver *CCIPDummyReceiverFilterer) ParseMessageReceived(log types.Log) (*CCIPDummyReceiverMessageReceived, error) {
	event := new(CCIPDummyReceiverMessageReceived)
	if err := _CCIPDummyReceiver.contract.UnpackLog(event, "MessageReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_CCIPDummyReceiver *CCIPDummyReceiver) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _CCIPDummyReceiver.abi.Events["MessageReceived"].ID:
		return _CCIPDummyReceiver.ParseMessageReceived(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (CCIPDummyReceiverMessageReceived) Topic() common.Hash {
	return common.HexToHash("0xd1e72bffc6abde0f1ef309ad9774dafcf9d29fe059a3fe5bb322ae9f15a74423")
}

func (_CCIPDummyReceiver *CCIPDummyReceiver) Address() common.Address {
	return _CCIPDummyReceiver.address
}

type CCIPDummyReceiverInterface interface {
	GetRouter(opts *bind.CallOpts) (common.Address, error)

	SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error)

	CcipReceive(opts *bind.TransactOpts, message ClientAny2EVMMessage) (*types.Transaction, error)

	FilterMessageReceived(opts *bind.FilterOpts) (*CCIPDummyReceiverMessageReceivedIterator, error)

	WatchMessageReceived(opts *bind.WatchOpts, sink chan<- *CCIPDummyReceiverMessageReceived) (event.Subscription, error)

	ParseMessageReceived(log types.Log) (*CCIPDummyReceiverMessageReceived, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
