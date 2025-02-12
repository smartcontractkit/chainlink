// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package usdc_reader_tester

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

var USDCReaderTesterMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"emitMessageSent\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"sourceDomain\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"destinationDomain\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"recipient\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"destinationCaller\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sender\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nonce\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"messageBody\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"MessageSent\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false}]",
	Bin: "0x6080806040523460155761034b908161001b8239f35b600080fdfe608080604052600436101561001357600080fd5b60003560e01c6362826f181461002857600080fd5b346102a6576101007ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc3601126102a65760043563ffffffff811681036102a65760243563ffffffff811681036102a6576044359063ffffffff821682036102a65760c4359367ffffffffffffffff851685036102a65760e4359067ffffffffffffffff82116102a657366023830112156102a65781600401359567ffffffffffffffff87116102a65736602488850101116102a657600087601f8299017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200161011490856102ab565b8084528060208501956024018637830160200152604051948594602086019760e01b7fffffffff0000000000000000000000000000000000000000000000000000000016885260e01b7fffffffff0000000000000000000000000000000000000000000000000000000016602486015260e01b7fffffffff0000000000000000000000000000000000000000000000000000000016602885015260c01b7fffffffffffffffff00000000000000000000000000000000000000000000000016602c84015260a435603484015260643560548401526084356074840152519081609484016102009261031b565b8101036094017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08101825261023590826102ab565b60405191829160208352519081602084015281604084016102559261031b565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168101036040017f8c5261668696ce22758910d05bab8f186d6eb247ceac2af2e82c7dc17669b03691a180f35b600080fd5b90601f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0910116810190811067ffffffffffffffff8211176102ec57604052565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60005b83811061032e5750506000910152565b818101518382015260200161031e56fea164736f6c634300081a000a",
}

var USDCReaderTesterABI = USDCReaderTesterMetaData.ABI

var USDCReaderTesterBin = USDCReaderTesterMetaData.Bin

func DeployUSDCReaderTester(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *USDCReaderTester, error) {
	parsed, err := USDCReaderTesterMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(USDCReaderTesterBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &USDCReaderTester{address: address, abi: *parsed, USDCReaderTesterCaller: USDCReaderTesterCaller{contract: contract}, USDCReaderTesterTransactor: USDCReaderTesterTransactor{contract: contract}, USDCReaderTesterFilterer: USDCReaderTesterFilterer{contract: contract}}, nil
}

type USDCReaderTester struct {
	address common.Address
	abi     abi.ABI
	USDCReaderTesterCaller
	USDCReaderTesterTransactor
	USDCReaderTesterFilterer
}

type USDCReaderTesterCaller struct {
	contract *bind.BoundContract
}

type USDCReaderTesterTransactor struct {
	contract *bind.BoundContract
}

type USDCReaderTesterFilterer struct {
	contract *bind.BoundContract
}

type USDCReaderTesterSession struct {
	Contract     *USDCReaderTester
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type USDCReaderTesterCallerSession struct {
	Contract *USDCReaderTesterCaller
	CallOpts bind.CallOpts
}

type USDCReaderTesterTransactorSession struct {
	Contract     *USDCReaderTesterTransactor
	TransactOpts bind.TransactOpts
}

type USDCReaderTesterRaw struct {
	Contract *USDCReaderTester
}

type USDCReaderTesterCallerRaw struct {
	Contract *USDCReaderTesterCaller
}

type USDCReaderTesterTransactorRaw struct {
	Contract *USDCReaderTesterTransactor
}

func NewUSDCReaderTester(address common.Address, backend bind.ContractBackend) (*USDCReaderTester, error) {
	abi, err := abi.JSON(strings.NewReader(USDCReaderTesterABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindUSDCReaderTester(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &USDCReaderTester{address: address, abi: abi, USDCReaderTesterCaller: USDCReaderTesterCaller{contract: contract}, USDCReaderTesterTransactor: USDCReaderTesterTransactor{contract: contract}, USDCReaderTesterFilterer: USDCReaderTesterFilterer{contract: contract}}, nil
}

func NewUSDCReaderTesterCaller(address common.Address, caller bind.ContractCaller) (*USDCReaderTesterCaller, error) {
	contract, err := bindUSDCReaderTester(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &USDCReaderTesterCaller{contract: contract}, nil
}

func NewUSDCReaderTesterTransactor(address common.Address, transactor bind.ContractTransactor) (*USDCReaderTesterTransactor, error) {
	contract, err := bindUSDCReaderTester(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &USDCReaderTesterTransactor{contract: contract}, nil
}

func NewUSDCReaderTesterFilterer(address common.Address, filterer bind.ContractFilterer) (*USDCReaderTesterFilterer, error) {
	contract, err := bindUSDCReaderTester(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &USDCReaderTesterFilterer{contract: contract}, nil
}

func bindUSDCReaderTester(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := USDCReaderTesterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_USDCReaderTester *USDCReaderTesterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _USDCReaderTester.Contract.USDCReaderTesterCaller.contract.Call(opts, result, method, params...)
}

func (_USDCReaderTester *USDCReaderTesterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _USDCReaderTester.Contract.USDCReaderTesterTransactor.contract.Transfer(opts)
}

func (_USDCReaderTester *USDCReaderTesterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _USDCReaderTester.Contract.USDCReaderTesterTransactor.contract.Transact(opts, method, params...)
}

func (_USDCReaderTester *USDCReaderTesterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _USDCReaderTester.Contract.contract.Call(opts, result, method, params...)
}

func (_USDCReaderTester *USDCReaderTesterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _USDCReaderTester.Contract.contract.Transfer(opts)
}

func (_USDCReaderTester *USDCReaderTesterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _USDCReaderTester.Contract.contract.Transact(opts, method, params...)
}

func (_USDCReaderTester *USDCReaderTesterTransactor) EmitMessageSent(opts *bind.TransactOpts, version uint32, sourceDomain uint32, destinationDomain uint32, recipient [32]byte, destinationCaller [32]byte, sender [32]byte, nonce uint64, messageBody []byte) (*types.Transaction, error) {
	return _USDCReaderTester.contract.Transact(opts, "emitMessageSent", version, sourceDomain, destinationDomain, recipient, destinationCaller, sender, nonce, messageBody)
}

func (_USDCReaderTester *USDCReaderTesterSession) EmitMessageSent(version uint32, sourceDomain uint32, destinationDomain uint32, recipient [32]byte, destinationCaller [32]byte, sender [32]byte, nonce uint64, messageBody []byte) (*types.Transaction, error) {
	return _USDCReaderTester.Contract.EmitMessageSent(&_USDCReaderTester.TransactOpts, version, sourceDomain, destinationDomain, recipient, destinationCaller, sender, nonce, messageBody)
}

func (_USDCReaderTester *USDCReaderTesterTransactorSession) EmitMessageSent(version uint32, sourceDomain uint32, destinationDomain uint32, recipient [32]byte, destinationCaller [32]byte, sender [32]byte, nonce uint64, messageBody []byte) (*types.Transaction, error) {
	return _USDCReaderTester.Contract.EmitMessageSent(&_USDCReaderTester.TransactOpts, version, sourceDomain, destinationDomain, recipient, destinationCaller, sender, nonce, messageBody)
}

type USDCReaderTesterMessageSentIterator struct {
	Event *USDCReaderTesterMessageSent

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *USDCReaderTesterMessageSentIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(USDCReaderTesterMessageSent)
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
		it.Event = new(USDCReaderTesterMessageSent)
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

func (it *USDCReaderTesterMessageSentIterator) Error() error {
	return it.fail
}

func (it *USDCReaderTesterMessageSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type USDCReaderTesterMessageSent struct {
	Arg0 []byte
	Raw  types.Log
}

func (_USDCReaderTester *USDCReaderTesterFilterer) FilterMessageSent(opts *bind.FilterOpts) (*USDCReaderTesterMessageSentIterator, error) {

	logs, sub, err := _USDCReaderTester.contract.FilterLogs(opts, "MessageSent")
	if err != nil {
		return nil, err
	}
	return &USDCReaderTesterMessageSentIterator{contract: _USDCReaderTester.contract, event: "MessageSent", logs: logs, sub: sub}, nil
}

func (_USDCReaderTester *USDCReaderTesterFilterer) WatchMessageSent(opts *bind.WatchOpts, sink chan<- *USDCReaderTesterMessageSent) (event.Subscription, error) {

	logs, sub, err := _USDCReaderTester.contract.WatchLogs(opts, "MessageSent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(USDCReaderTesterMessageSent)
				if err := _USDCReaderTester.contract.UnpackLog(event, "MessageSent", log); err != nil {
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

func (_USDCReaderTester *USDCReaderTesterFilterer) ParseMessageSent(log types.Log) (*USDCReaderTesterMessageSent, error) {
	event := new(USDCReaderTesterMessageSent)
	if err := _USDCReaderTester.contract.UnpackLog(event, "MessageSent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_USDCReaderTester *USDCReaderTester) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _USDCReaderTester.abi.Events["MessageSent"].ID:
		return _USDCReaderTester.ParseMessageSent(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (USDCReaderTesterMessageSent) Topic() common.Hash {
	return common.HexToHash("0x8c5261668696ce22758910d05bab8f186d6eb247ceac2af2e82c7dc17669b036")
}

func (_USDCReaderTester *USDCReaderTester) Address() common.Address {
	return _USDCReaderTester.address
}

type USDCReaderTesterInterface interface {
	EmitMessageSent(opts *bind.TransactOpts, version uint32, sourceDomain uint32, destinationDomain uint32, recipient [32]byte, destinationCaller [32]byte, sender [32]byte, nonce uint64, messageBody []byte) (*types.Transaction, error)

	FilterMessageSent(opts *bind.FilterOpts) (*USDCReaderTesterMessageSentIterator, error)

	WatchMessageSent(opts *bind.WatchOpts, sink chan<- *USDCReaderTesterMessageSent) (event.Subscription, error)

	ParseMessageSent(log types.Log) (*USDCReaderTesterMessageSent, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
