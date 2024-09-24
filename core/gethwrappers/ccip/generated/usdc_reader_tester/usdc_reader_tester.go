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
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"MessageSent\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"version\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"sourceDomain\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"destinationDomain\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"recipient\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"destinationCaller\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"sender\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"nonce\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"messageBody\",\"type\":\"bytes\"}],\"name\":\"emitMessageSent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5061032c806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c806362826f1814610030575b600080fd5b61004361003e366004610129565b610045565b005b600061008d8a8a8a87898c8c8a8a8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506100d292505050565b90507f8c5261668696ce22758910d05bab8f186d6eb247ceac2af2e82c7dc17669b036816040516100be9190610228565b60405180910390a150505050505050505050565b606088888888888888886040516020016100f3989796959493929190610279565b604051602081830303815290604052905098975050505050505050565b803563ffffffff8116811461012457600080fd5b919050565b60008060008060008060008060006101008a8c03121561014857600080fd5b6101518a610110565b985061015f60208b01610110565b975061016d60408b01610110565b965060608a0135955060808a0135945060a08a0135935060c08a013567ffffffffffffffff80821682146101a057600080fd5b90935060e08b013590808211156101b657600080fd5b818c0191508c601f8301126101ca57600080fd5b8135818111156101d957600080fd5b8d60208285010111156101eb57600080fd5b6020830194508093505050509295985092959850929598565b60005b8381101561021f578181015183820152602001610207565b50506000910152565b6020815260008251806020840152610247816040850160208701610204565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169190910160400192915050565b60007fffffffff00000000000000000000000000000000000000000000000000000000808b60e01b168352808a60e01b166004840152808960e01b166008840152507fffffffffffffffff0000000000000000000000000000000000000000000000008760c01b16600c830152856014830152846034830152836054830152825161030b816074850160208701610204565b91909101607401999850505050505050505056fea164736f6c6343000818000a",
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
